package CompraControler

import (
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modelos/CatalogoModel"
	//"../../Modelos/ImpuestoModel"
	"../../Modelos/ListaPrecioModel"
	"../../Modelos/OperacionModel"
	"../../Modelos/ProductoModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/Session"
	_ "github.com/bmizerany/pq"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//##########< Variables Generales > ############

var cadenaBusqueda string
var buscarEn string
var numeroRegistros int64
var paginasTotales int

//NumPagina especifica ***************
var NumPagina float32

//limitePorPagina especifica ***************
var limitePorPagina = 10
var result []ListaPrecioModel.ListaPrecio
var resultPage []ListaPrecioModel.ListaPrecio
var templatePaginacion = ``

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Compras
func IndexGet(ctx *iris.Context) {

	var Send OperacionModel.SOperacion
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	Send.SSesion.Name = NameUsrLoged
	Send.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	Send.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		Send.SEstado = false
		Send.SMsj = errSes.Error()
		ctx.Render("ZError.html", Send)
		return
	}

	ctx.Render("ComprasIndex.html", Send)
}

//IndexPost regresa la peticon post que se hizo desde el index de Compras
func IndexPost(ctx *iris.Context) {

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Compra
func AltaGet(ctx *iris.Context) {
	var Send AlmacenModel.SAlmacen

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	Send.SSesion.Name = NameUsrLoged
	Send.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	Send.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		Send.SEstado = false
		Send.SMsj = errSes.Error()
		ctx.Render("ZError.html", Send)
		return
	}

	almacenes, err := AlmacenModel.GetEspecificsByTagAndTestConexion("Clasificacion", bson.ObjectIdHex("58e5692ee75770120c60befa"))
	var cadena string
	for _, value := range err {
		cadena += value + "<br>"
	}

	Send.AlmacenesDisponibles = len(almacenes)
	if len(almacenes) > 0 {
		if len(err) > 0 {
			Send.SEstado = false
			Send.SMsj = "Se encontraron errores de conexion en los siguients almacenes<br>" + cadena
		} else {
			Send.SEstado = true
		}
		Send.Almacen.ENombreAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboAlmacenes(bson.ObjectIdHex("58e5692ee75770120c60befa")))

		Send.Almacen.ETipoAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo2(149, ""))

		Send.Almacen.EClasificacionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboImpuesto(bson.ObjectIdHex("58e43a2ae757702ffce31c48"), bson.ObjectIdHex("58e43a2ae757702ffce31c48")))
	} else {
		Send.SMsj = "No se ha podido cargar los almacenes..." + cadena
		Send.SEstado = false
	}
	ctx.Render("ComprasAlta.html", Send)
}

//AltaPost regresa la petición post que se hizo desde el alta de compra
func AltaPost(ctx *iris.Context) {
	var Send AlmacenModel.SAlmacen
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	Send.SSesion.Name = NameUsrLoged
	Send.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	Send.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		Send.SEstado = false
		Send.SMsj = errSes.Error()
		ctx.Render("ZError.html", Send)
		return
	}
	//######### Leyendo los objetos de la compra #########
	CantidadGeneral := ctx.FormValue("CantidadComprada")

	//Invertí costos y precio para poder solventar el problema eventualmente.

	PrecioGeneral := ctx.FormValue("CostoGeneral")
	CostoGeneral := ctx.FormValue("PrecioGeneral")

	IDProducto := ctx.FormValue("idProducto")
	AlmacenDefecto := ctx.FormValue("AlmacenDefecto")
	Almacenes := ctx.Request.Form["Almacenes"]
	Cantidades := ctx.Request.Form["Cantidades"]

	//########## Obsoletos
	// Impuesto := ctx.FormValue("Impuesto")
	// Factor := ctx.FormValue("Factor")
	// Tratamiento := ctx.FormValue("Tratamiento")
	//##########

	TiposDeImpuestos := ctx.Request.Form["TipoImpuestoLista"]
	ValoresDeImpuestos := ctx.Request.Form["ValorImpuestoLista"]
	FactoresDeImpuestos := ctx.Request.Form["FactorImpuestoLista"]
	TratamientosDeImpuestos := ctx.Request.Form["TratamientoImpuestoLista"]

	//Se crean las operaciones para almacenar en MOngo
	var operacion OperacionModel.OperacionMgo
	operacion.ID = bson.NewObjectId()
	operacion.UsuarioOrigen = bson.ObjectIdHex("58e1474bc4ae08e9ccc4d73d")
	operacion.UsuarioDestino = bson.ObjectIdHex("58ed29f924a59d3a683d14cf")
	operacion.TipoOperacion = bson.ObjectIdHex("58efbf8bd2b2131778e9c928")
	fechaHoraRegistro := time.Now()
	operacion.FechaHoraRegistro = fechaHoraRegistro
	operacion.Estatus = bson.ObjectIdHex("58e43924e757703988714a93")

	//Conversion de datos numericos
	cantidad, err := strconv.ParseFloat(CantidadGeneral, 64)
	if err != nil {
		fmt.Println("Ocurrio un error al momento de convertir la cadena: ", CantidadGeneral, err)
	}

	costo, err := strconv.ParseFloat(CostoGeneral, 64)
	if err != nil {
		fmt.Println("Ocurrio un error al momento de convertir la cadena: ", CostoGeneral, err)
	}
	precio, err := strconv.ParseFloat(PrecioGeneral, 64)
	if err != nil {
		fmt.Println("Ocurrio un error al momento de convertir la cadena: ", PrecioGeneral, err)
	}

	var totalimp float64
	for k, v := range ValoresDeImpuestos {
		imp, err := strconv.ParseFloat(v, 64)
		if err != nil {
			fmt.Println("Ocurrio un error al momento de convertir la cadena: ", v, err)
		}

		if FactoresDeImpuestos[k] == "Tasa" {
			totalimp += precio * (imp / 100)
		} else if FactoresDeImpuestos[k] == "Cuota" {
			totalimp += imp
		} else if FactoresDeImpuestos[k] == "Exento" {
			totalimp += 0
		} else {
			totalimp += 0
		}
	}
	fmt.Println("Al comprar se calcula un total de impuesto por producto de: $", totalimp)
	//Se llena la estructura kardex para posteriormente insertar los datos en postgres
	var kardex OperacionModel.KardexPostgres
	kardex.IDOperacion = operacion.ID.Hex()
	idMovimiento := bson.NewObjectId()
	kardex.IDMovimiento = idMovimiento.Hex()
	kardex.IDProducto = IDProducto
	kardex.Cantidad = cantidad
	kardex.Costo = costo
	kardex.Precio = precio
	kardex.ImpuestoTotal = totalimp * cantidad
	kardex.DescuentoTotal = 0
	kardex.TipoOperacion = operacion.TipoOperacion.Hex()
	kardex.Existencia = 0
	kardex.FechaHora = fechaHoraRegistro

	//Se creaan los Movimientos
	movimientos := OperacionModel.MovimientoMgo{}
	movimientos.IDMovimiento = idMovimiento
	movimientos.AlmacenDestino = bson.ObjectIdHex(AlmacenDefecto)
	movimientos.Estatus = bson.ObjectIdHex("58e43924e757703988714a93")
	operacion.Movimientos = append(operacion.Movimientos, movimientos)

	//Se inserta la operacion en mongoDB, si ocurrió un error,muestra un mensaje en pantalla
	if !operacion.InsertaMgo() {
		fmt.Println("Ocurrio un detalle al insertar", operacion)
	}

	//Crear un arreglo de kardex para almacenar todos los movimientos generados en los distintos almacenes
	var kardexTotal []OperacionModel.Kardex

	//Lllenar la estructura Kardex con el almacen por defecto
	var kardexLoc OperacionModel.Kardex
	almacen := AlmacenModel.GetOne(bson.ObjectIdHex(AlmacenDefecto))
	kardexLoc.NombreAlmacen = almacen.Nombre
	kardexLoc.IDMovimiento = kardex.IDMovimiento
	producto := ProductoModel.GetOne(bson.ObjectIdHex(kardex.IDProducto))
	kardexLoc.IDProducto = producto.Nombre
	kardexLoc.Cantidad = kardex.Cantidad
	kardexLoc.Costo = kardex.Costo
	kardexLoc.Precio = kardex.Precio
	kardexLoc.ImpuestoTotal = kardex.ImpuestoTotal
	kardexLoc.DescuentoTotal = kardex.DescuentoTotal
	nombreOperacion := CatalogoModel.RegresaNombreSubCatalogo(bson.ObjectIdHex(kardex.TipoOperacion))
	kardexLoc.TipoOperacion = nombreOperacion
	kardexLoc.Existencia = kardex.Existencia
	kardexLoc.FechaHora = kardex.FechaHora
	kardexTotal = append(kardexTotal, kardexLoc)

	/*
		SKardex.IDMovimiento = kardex.IDMovimiento

		SKardex.IDProducto =
		SKardex.Cantidad = kardex.Cantidad
		SKardex.Costo = kardex.Costo
		SKardex.Precio = kardex.Precio
		nombreOperacion := CatalogoModel.RegresaNombreSubCatalogo(bson.ObjectIdHex(kardex.TipoOperacion))
		SKardex.TipoOperacion = nombreOperacion
		SKardex.Existencia = kardex.Existencia
		SKardex.FechaHora = kardex.FechaHora
	*/

	//OBSOLETO

	// var ImpuestoToInsert ImpuestoModel.ImpuestoMgo
	// idImp := bson.NewObjectId()
	// ImpuestoToInsert.ID = idImp
	// ImpuestoToInsert.Descripcion = "IMPUESTO DADO DE ALTA EN COMPRA PARA EL PRODUCTO: " + producto.Nombre + "CON ID: " + IDProducto
	// ImpuestoToInsert.TipoImpuesto = bson.ObjectIdHex("58e43a2ae757702ffce31c48")

	// var Datos ImpuestoModel.DataImpuestoMgo

	// Datos.Nombre = "IVA DE " + Impuesto + " Para: " + producto.Nombre
	// Datos.Max = imp

	// if Tratamiento == "Retenido" {
	// 	Datos.Retencion = true
	// } else {
	// 	Datos.Traslado = true
	// }

	// if Factor == "Tasa" {
	// 	Datos.TipoFactor = bson.ObjectIdHex("58e43203e757702b203de258")
	// } else if Factor == "Cuota" {
	// 	Datos.TipoFactor = bson.ObjectIdHex("58e43203e757702b203de259")
	// } else if Factor == "Exento" {
	// 	Datos.TipoFactor = bson.ObjectIdHex("58e43203e757702b203de25a")
	// } else {
	// 	Datos.TipoFactor = bson.ObjectIdHex("58e43203e757702b203de259")
	// }

	// Datos.Rango = false
	// Datos.Unidad = bson.ObjectIdHex("58e4329fe757702b203de25c")
	// Datos.FechaHora = fechaHoraRegistro
	// ImpuestoToInsert.Datos = Datos
	// ImpuestoToInsert.Editable = true
	// ImpuestoToInsert.Estatus = bson.ObjectIdHex("58e43924e757703988714a90")
	// ImpuestoToInsert.FechaHora = fechaHoraRegistro
	// flag1 := false

	// if ImpuestoToInsert.InsertaMgo() {
	// 	flag1 = true
	// 	if ImpuestoToInsert.InsertaElastic() {
	// 		fmt.Println("SE INSERTÓ el impuseto en Mongo y en elastic")
	// 		Send.SEstado = true
	// 	}
	// } else {
	// 	fmt.Println("IMPUESTOS ERROR: No se pudo insertar el impuesto en Mongo")
	// 	Send.SEstado = false
	// 	Send.SMsj += template.HTML("No se pudo insertar en el impuesto en Mongo<br>")
	// }

	// var impuesto OperacionModel.ImpuestoPostgres
	// impuesto.IDMovimiento = idMovimiento.Hex()
	// impuesto.IDProducto = IDProducto
	// impuesto.Valor = imp
	// impuesto.FechaHora = fechaHoraRegistro
	// if flag1 {
	// 	impuesto.IDImpuesto = idImp.Hex()
	// }

	//##################

	var ImpuestosPsql []OperacionModel.ImpuestoPostgres
	var ImpuestoPsql OperacionModel.ImpuestoPostgres

	for k := range TiposDeImpuestos {
		ImpuestoPsql = OperacionModel.ImpuestoPostgres{}
		ImpuestoPsql.IDMovimiento = idMovimiento.Hex()
		ImpuestoPsql.IDProducto = IDProducto
		ImpuestoPsql.TipoDeImpuesto = TiposDeImpuestos[k]
		ImpuestoPsql.Valor = ValoresDeImpuestos[k]
		ImpuestoPsql.Factor = FactoresDeImpuestos[k]
		ImpuestoPsql.Tratamiento = TratamientosDeImpuestos[k]
		ImpuestoPsql.FechaHora = time.Now()
		ImpuestosPsql = append(ImpuestosPsql, ImpuestoPsql)
	}

	//Se insertarn los mismos elementos en los almacenes restantes
	for key := range Almacenes {
		cant, _ := strconv.ParseFloat(Cantidades[key], 64)
		if cant > 0 {
			var idAlmacen = bson.ObjectIdHex(Almacenes[key])
			var inventario OperacionModel.InventarioPostgres

			//Inserta los datos correspondientes al inventario en la estructura inventario
			inventario.IDProducto = IDProducto
			cant, err := strconv.ParseFloat(Cantidades[key], 64)
			if err != nil {
				fmt.Println("Ocurrio un error al momento de convertir la cadena: ", CantidadGeneral, err)
			}
			inventario.Existencia = cant
			inventario.Estatus = "ACTIVO"
			inventario.Costo = costo
			inventario.Precio = precio

			mje, insert, err := kardex.InsertaKardexYActualizaInventario(idAlmacen)
			if !insert {
				//Inserta en inventario y kardex de postgres
				mensaje, encontrado, err := kardex.InsertaKardexYActualizaInventario(idAlmacen)

				fmt.Println("YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY", mje, mensaje, encontrado, err)
				producto.ActualizaAlmacenes(idAlmacen)
			} else {
				producto.ActualizaAlmacenes(idAlmacen)
				//Crear un kardex por cada movimiento en los almacenes extras
				//Posteriormente se agrega al arreglo de kardex para mostrarlo en la viista correspondiente
				var kardexLoc OperacionModel.Kardex
				almacen := AlmacenModel.GetOne(idAlmacen)
				kardexLoc.NombreAlmacen = almacen.Nombre
				kardexLoc.IDMovimiento = kardex.IDMovimiento
				kardexLoc.IDProducto = producto.Nombre
				kardexLoc.Cantidad = cant
				kardexLoc.Costo = kardex.Costo
				kardexLoc.Precio = kardex.Precio
				kardexLoc.ImpuestoTotal = kardex.ImpuestoTotal
				kardexLoc.DescuentoTotal = kardex.DescuentoTotal
				nombreOperacion := CatalogoModel.RegresaNombreSubCatalogo(bson.ObjectIdHex(kardex.TipoOperacion))
				kardexLoc.TipoOperacion = nombreOperacion
				kardexLoc.Existencia = kardex.Existencia
				kardexLoc.FechaHora = kardex.FechaHora

				kardexTotal = append(kardexTotal, kardexLoc)
			}

			for _, v := range ImpuestosPsql {
				msj, estado, err2 := v.InsertaImpuestoEnAlmacenComprasPsql(idAlmacen)
				if !estado {
					fmt.Println("IMPUESTOS ERROR (!estado): ", msj, err2)
					Send.SEstado = false
					Send.SMsj += "No se pudo insertar en el almacen seleccionado El impuesto : " + err2.Error() + "   " + msj
				}
				msj, estado, err2 = v.InsertaImpuestoEnAlmacenVentasPsql(idAlmacen)
				if !estado {
					fmt.Println("IMPUESTOS ERROR (estado): ", msj, err2)
					Send.SEstado = false
					Send.SMsj += "No se pudo insertar en el almacen seleccionado El impuesto : " + err2.Error() + "   " + msj
				}

			}
		}

	}

	//Se insertara el idproducto, cantidad, costo, precio en la tabla kardex del almacen correspondiente
	//nombreTablaKardex := "Kardex_" + AlmacenDefecto
	idAlmacen := bson.ObjectIdHex(AlmacenDefecto)
	mensaje, encontrado, err := kardex.InsertaKardexYActualizaInventario(idAlmacen)

	for _, v := range ImpuestosPsql {
		msj, estado, err2 := v.InsertaImpuestoEnAlmacenComprasPsql(idAlmacen)
		if !estado {
			fmt.Println("IMPUESTOS ERROR (!estado:default): ", msj, err2)
			Send.SEstado = false
			Send.SMsj += "No se pudo insertar en el almacen seleccionado El impuesto :	" + err2.Error() + "	" + msj
		}
		msj, estado, err2 = v.InsertaImpuestoEnAlmacenVentasPsql(idAlmacen)
		if !estado {
			fmt.Println("IMPUESTOS ERROR: (=estado:default)", msj, err2)
			Send.SEstado = false
			Send.SMsj += "No se pudo insertar en el almacen seleccionado El impuesto :	" + err2.Error() + "	" + msj
		}
	}

	if err == nil {
		if encontrado {
			producto.ActualizaAlmacenes(idAlmacen)
			SKardex := EstructuraKardexvista(&kardexTotal, NameUsrLoged)
			//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
			ctx.Render("ComprasDetalle.html", SKardex)
		} else {
			var inventario OperacionModel.InventarioPostgres
			inventario.IDProducto = IDProducto
			inventario.Existencia = cantidad
			inventario.Estatus = "ACTIVO"
			inventario.Costo = costo
			inventario.Precio = precio
			mensaje, encontrado, err := kardex.InsertaKardexInsertaInventario(idAlmacen, inventario)
			if err == nil {
				if encontrado {
					producto.ActualizaAlmacenes(idAlmacen)
					SKardex := EstructuraKardexvista(&kardexTotal, NameUsrLoged)
					fmt.Println(SKardex)
					/*
						var campoAlmacen []string
						campoAlmacen = append(campoAlmacen, "Almacenes")
						vals := make([]interface{}, 1)
						vals[0] = idAlmacen
						producto.ActualizaMgo(campoAlmacen, vals)
					*/
					//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
					ctx.Render("ComprasDetalle.html", SKardex)
				}
			}
			fmt.Println(mensaje, encontrado, err)
		}
	} else {
		Send.Almacen.ENombreAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboAlmacenes(bson.ObjectIdHex("58e5692ee75770120c60befa")))
		Send.AlmacenesDisponibles = 1
		Send.SMsj += "No se pudo insertar en el almacen seleccionado:	" + err.Error() + "		" + mensaje
		Send.SEstado = false
		ctx.Render("ComprasAlta.html", Send)
	}

}

//EstructuraKardexvista crea los datos para mostrar en la vista al finalizar la compra
func EstructuraKardexvista(kardex *[]OperacionModel.Kardex, name string) OperacionModel.SKardex {

	var SKardex OperacionModel.SKardex

	SKardex.SSesion.Name = name

	SKardex.SEstado = true
	SKardex.SMsj = "La compra se ha registrado satisfactoriamente"

	SKardex.Kardex = *kardex
	return SKardex
}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de ListaPrecio
func EditaGet(ctx *iris.Context) {
	ctx.Render("ListaPrecioEdita.html", nil)
}

//EditaPost regresa el resultado de la petición post generada desde la edición de ListaPrecio
func EditaPost(ctx *iris.Context) {
	ctx.Render("ListaPrecioEdita.html", nil)
}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	ctx.Render("ComprasDetalle.html", nil)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	ctx.Render("ComprasDetalle.html", nil)
}

//####################< RUTINAS ADICIONALES >##########################

//Totalpaginas calcula el número de paginaciones de acuerdo al número
// de resultados encontrados y los que se quieren mostrar en la página.
func Totalpaginas() int {

	NumPagina = float32(numeroRegistros) / float32(limitePorPagina)
	NumPagina2 := int(NumPagina)
	if NumPagina > float32(NumPagina2) {
		NumPagina2++
	}
	totalpaginas := NumPagina2
	return totalpaginas

}

//ConstruirPaginacion construtye la paginación en formato html para usarse en la página
func ConstruirPaginacion() string {
	var templateP string
	templateP += `
	<nav aria-label="Page navigation">
		<ul class="pagination">
			<li>
				<a href="/ListaPrecios/1" aria-label="Primera">
				<span aria-hidden="true">&laquo;</span>
				</a>
			</li>`

	templateP += ``
	for i := 0; i <= paginasTotales; i++ {
		if i == 1 {

			templateP += `<li class="active"><a href="/ListaPrecios/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		} else if i > 1 && i < 11 {
			templateP += `<li><a href="/ListaPrecios/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`

		} else if i > 11 && i == paginasTotales {
			templateP += `<li><span aria-hidden="true">...</span></li><li><a href="/ListaPrecios/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		}
	}
	templateP += `<li><a href="/ListaPrecios/` + strconv.Itoa(paginasTotales) + `" aria-label="Ultima"><span aria-hidden="true">&raquo;</span></a></li></ul></nav>`
	return templateP
}
