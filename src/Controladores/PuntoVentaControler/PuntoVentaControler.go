package PuntoVentaControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/OperacionModel"
	"../../Modelos/ProductoModel"
	"../../Modelos/PuntoVentaModel"
	"../../Modelos/UsuarioModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
	"../../Modulos/ConsultasSql"
	"../../Modulos/General"
	"../../Modulos/Session"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//##########< Variables Generales > ############

var cadenaBusqueda string
var numeroRegistros int
var paginasTotales int

//NumPagina especifica el numero de página en la que se cargarán los registros
var NumPagina int

//limitePorPagina limite de registros a mostrar en la pagina
var limitePorPagina = 10

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de PuntoVenta
func IndexGet(ctx *iris.Context) {

	var Send PuntoVentaModel.SDataProducto
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

	var Cabecera, Cuerpo string

	PuntoVentas := OperacionModel.GetEspecificsByFields2([]string{"TipoOperacion", "Estatus"}, []interface{}{bson.ObjectIdHex("58efbf8bd2b2131778e9c929"), bson.ObjectIdHex("58efc5c5d2b2131778e9c931")})

	numeroRegistros = len(PuntoVentas)
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)

	arrIDMgo = []bson.ObjectId{}
	for _, v := range PuntoVentas {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(PuntoVentas[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(PuntoVentas[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("PuntoVentaIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de PuntoVenta
func IndexPost(ctx *iris.Context) {
	var Send PuntoVentaModel.SDataProducto
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

	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.PuntoVenta.EVARIABLEPuntoVenta.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := OperacionModel.BuscarEnElastic(cadenaBusqueda)

		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)

			arrToMongo = []bson.ObjectId{}
			if numeroRegistros <= limitePorPagina {
				for _, v := range arrIDElastic[0:numeroRegistros] {
					arrToMongo = append(arrToMongo, v)
				}
			} else if numeroRegistros >= limitePorPagina {
				for _, v := range arrIDElastic[0:limitePorPagina] {
					arrToMongo = append(arrToMongo, v)
				}
			}

			MoConexion.FlushElastic()

			Cabecera, Cuerpo := OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
			}

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

			Send.SIndex.SRMsj = "No se encontraron resultados para: " + cadenaBusqueda + " ."
		}

		Send.SEstado = true

	} else {
		Send.SEstado = false
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
		Send.SResultados = false
	}
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	ctx.Render("PuntoVentaIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de PuntoVenta
func AltaGet(ctx *iris.Context) {

	var Send PuntoVentaModel.SPuntoVenta
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

	//Carga Combo de Almacenes Propios
	Send.PuntoVenta.EOperacionPuntoVenta.Ihtml = template.HTML(CargaCombos.CargaComboAlmacenes(bson.ObjectIdHex("58e5692ee75770120c60befa")))

	Send.PuntoVenta.EOperacionPuntoVenta.Operacion = bson.NewObjectId()
	fmt.Println("Nueva Venta--> :", Send.PuntoVenta.EOperacionPuntoVenta.Operacion.Hex())
	ctx.Render("PuntoVentaAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de PuntoVenta
func AltaPost(ctx *iris.Context) {

	var Send PuntoVentaModel.SDataProducto
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

	operacion := ctx.FormValue("Operacion")
	if PuntoVentaModel.AplicaOperacionParaPago(operacion) {
		Send.SEstado = true
		Send.SMsj = "Registro de Pago en Kardex Exitosa."
	} else {
		Send.SEstado = false
		Send.SMsj = "Registro de Pago en Kardex con errores."
	}
	Send.SIhtml = template.HTML(PuntoVentaModel.GeneraBotonesPago(operacion))

	ctx.Render("PuntoVentaDetalle.html", Send)

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de PuntoVenta
func EditaGet(ctx *iris.Context) {
	var SDP PuntoVentaModel.SPuntoVenta
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SDP.SSesion.Name = NameUsrLoged
	SDP.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SDP.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SDP.SEstado = false
		SDP.SMsj = errSes.Error()
		ctx.Render("ZError.html", SDP)
		return
	}

	idOperacion := ctx.Param("ID")

	if bson.IsObjectIdHex(idOperacion) { //Es un ObjectId
		SDP.PuntoVenta.EOperacionPuntoVenta.Operacion = bson.ObjectIdHex(idOperacion)

		operacion := OperacionModel.GetOne(bson.ObjectIdHex(idOperacion))
		if !MoGeneral.EstaVacio(operacion) {
			estatus := CatalogoModel.GetSubEspecificByFields("Valores._id", operacion.Estatus)
			if estatus.Valor == "PENDIENTE" { //La operacion debe estar aun PENDIENTE
				carrito, calculadora := PuntoVentaModel.ConsultaDatosOperacion(idOperacion)
				SDP.SEstado = true
				SDP.SMsj = "Operación de Venta Encontrada"
				SDP.PuntoVenta.ECarritoPuntoVenta.Ihtml = template.HTML(carrito)
				SDP.PuntoVenta.EResumenPuntoVenta.Ihtml = template.HTML(calculadora)
			}
		} else {
			SDP.SEstado = false
			SDP.SMsj = "Posiblemente hay un error con su Conexión, favor de intentar de nuevo."
		}
	} else {
		SDP.SEstado = false
		SDP.SMsj = "No se encuentra la referencia a la Operacion solicitada."
	}
	//Carga Combo de Almacenes Propios
	SDP.PuntoVenta.EOperacionPuntoVenta.Ihtml = template.HTML(CargaCombos.CargaComboAlmacenes(bson.ObjectIdHex("58e5692ee75770120c60befa")))

	ctx.Render("PuntoVentaEdita.html", SDP)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de PuntoVenta
func EditaPost(ctx *iris.Context) {

	var Send PuntoVentaModel.SDataProducto
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

	operacion := ctx.FormValue("Operacion")

	// if PuntoVentaModel.ActualizaOperacionParaPago(operacion) {
	// 	Send.SEstado = true
	// 	Send.SMsj = "Actualizacion de Registro de Pago en Kardex Exitosa."
	// } else {
	// 	Send.SEstado = false
	// 	Send.SMsj = "Actualizacion de Registro de Pago en Kardex con Errores."
	// }

	Send.SEstado = true
	Send.SMsj = `Proceso de Pago de la Operación:  ` + `"` + operacion + `"`

	Send.SIhtml = template.HTML(PuntoVentaModel.GeneraBotonesPago(operacion))

	ctx.Render("PuntoVentaDetalle.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	return
}

//DetallePost solicita info para genenerar ticket
func DetallePost(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto
	id := ctx.FormValue("operacion")
	SDP.SIhtml = template.HTML(PuntoVentaModel.GetHTMLMontoOperacion(id))
	SDP.ID = id
	SDP.SEstado = true
	SDP.SMsj = "Ticket bien."
	jData, _ := json.Marshal(SDP)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//####################< RUTINAS ADICIONALES >##########################

//TraePrimerDatoDeProducto recibe la peticion de un producto mediante ajax y devuelve un template html
func TraePrimerDatoDeProducto(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto
	var Operacion OperacionModel.OperacionMgo

	cliente := ctx.FormValue("Cliente")
	movimiento := ""
	operacion := ctx.FormValue("Operacion")
	almacen := ctx.FormValue("Almacen")

	if almacen == "" {
		almacen = "592f1012e757701c5075c192"
	}

	if bson.IsObjectIdHex(operacion) {
		Operacion = OperacionModel.GetOne(bson.ObjectIdHex(operacion))
		if MoGeneral.EstaVacio(Operacion) {

			Operacion.ID = bson.ObjectIdHex(operacion)
			id, _, _, _ := Session.GetDataSession(ctx)

			if bson.IsObjectIdHex(id) {
				UsuarioOrigen := UsuarioModel.GetOne(bson.ObjectIdHex(id))
				Operacion.UsuarioOrigen = UsuarioOrigen.ID
			} else {
				ctx.Redirect("/", 403)
			}

			if !bson.IsObjectIdHex(cliente) {
				Operacion.UsuarioDestino = bson.ObjectIdHex("59139e69e757703c3400d6f9")
			} else {
				Operacion.UsuarioDestino = bson.ObjectIdHex(cliente)
			}

			Operacion.TipoOperacion = bson.ObjectIdHex("58efbf8bd2b2131778e9c929")
			Operacion.Estatus = bson.ObjectIdHex("58efc5c5d2b2131778e9c931")
			Operacion.FechaHoraRegistro = time.Now()

			//Sólo existe un movimiento al momento de dar de alta la operación
			Movimiento := OperacionModel.MovimientoMgo{}
			Transaccion := OperacionModel.TransaccionMgo{}

			Movimiento.IDMovimiento = bson.NewObjectId()
			Movimiento.AlmacenOrigen = bson.ObjectIdHex(almacen)

			if !bson.IsObjectIdHex(cliente) {
				Movimiento.AlmacenDestino = bson.ObjectIdHex("58ed23518c649f28f0445013")
			} else {
				Movimiento.AlmacenDestino = bson.ObjectIdHex(cliente)
			}

			Movimiento.Estatus = bson.ObjectIdHex("58efc5c5d2b2131778e9c931")

			//Sólo existe una transaccion, porque sólo existe una ruta.
			Transaccion.IDTransaccion = bson.NewObjectId()

			Transaccion.AlmacenOrigen = bson.ObjectIdHex(almacen)

			if !bson.IsObjectIdHex(cliente) {
				Transaccion.AlmacenDestino = bson.ObjectIdHex("58ed23518c649f28f0445013")
			} else {
				Transaccion.AlmacenDestino = bson.ObjectIdHex(cliente)
			}
			Transaccion.Estatus = bson.ObjectIdHex("58efc5c5d2b2131778e9c931")

			Movimiento.Transacciones = append(Movimiento.Transacciones, Transaccion)
			Operacion.Movimientos = append(Operacion.Movimientos, Movimiento)

			movimiento = Movimiento.IDMovimiento.Hex()

			if !Operacion.InsertaMgo() {
				SDP.ID = operacion
				SDP.SEstado = false
				SDP.SMsj = "No se pudo crear la operación, recargue la página y vuelva a intentarlo, disculpe las molestia."

				jData, _ := json.Marshal(SDP)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}

		}
	} else {
		SDP.ID = operacion
		SDP.SEstado = false
		SDP.SMsj = "No se pudo hacer referencia a la operación correspondiente, favor de recargar la página y vuelva a intentarlo. Disculpe las molestias."

		jData, _ := json.Marshal(SDP)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return

	}

	codigo := ctx.FormValue("Codigo")

	if codigo != "" {
		producto := ProductoModel.GetEspecificByFields("Codigos.Valores", codigo)

		if !MoGeneral.EstaVacio(producto) {

			//Aquí debo discernir si el producto corresponde al mismo movimiento o si hay uno nuevo
			encontrado := false
			if movimiento == "" { //Implica que no es el primer producto y tenemos que buscar si corresponde a un movimiento adecuado
				for _, v := range Operacion.Movimientos {
					encontrado = false
					if v.AlmacenOrigen.Hex() == almacen { //Implica que el movimiento ya existe y ponemos una bandera y obtenemos el id del movimiento que encontramos
						movimiento = v.IDMovimiento.Hex()
						encontrado = true
					}
				}

				if !encontrado { //implica que tendremos que crear un nuevo movimiento en la operacion
					Movimiento := OperacionModel.MovimientoMgo{}
					Transaccion := OperacionModel.TransaccionMgo{}

					Movimiento.IDMovimiento = bson.NewObjectId()
					Movimiento.AlmacenOrigen = bson.ObjectIdHex(almacen)

					if !bson.IsObjectIdHex(cliente) {
						Movimiento.AlmacenDestino = bson.ObjectIdHex("58ed23518c649f28f0445013")
					} else {
						Movimiento.AlmacenDestino = bson.ObjectIdHex(cliente)
					}

					Movimiento.Estatus = bson.ObjectIdHex("58efc5c5d2b2131778e9c931")

					//Sólo existe una transaccion, porque sólo existe una ruta.
					Transaccion.IDTransaccion = bson.NewObjectId()
					Transaccion.AlmacenOrigen = bson.ObjectIdHex(almacen)

					if !bson.IsObjectIdHex(cliente) {
						Transaccion.AlmacenDestino = bson.ObjectIdHex("58ed23518c649f28f0445013")
					} else {
						Transaccion.AlmacenDestino = bson.ObjectIdHex(cliente)
					}
					Transaccion.Estatus = bson.ObjectIdHex("58efc5c5d2b2131778e9c931")

					Movimiento.Transacciones = append(Movimiento.Transacciones, Transaccion)
					Operacion.Movimientos = append(Operacion.Movimientos, Movimiento)
					movimiento = Movimiento.IDMovimiento.Hex()

					if !Operacion.ActualizaMgo([]string{"Movimientos"}, []interface{}{Operacion.Movimientos}) {
						SDP.ID = operacion
						SDP.SEstado = false
						SDP.SMsj = "No se pudo actualizar la Operación, por favor revisa tu conexión."

						jData, _ := json.Marshal(SDP)
						ctx.Header().Set("Content-Type", "application/json")
						ctx.Write(jData)
						return
					}
				}
			}

			actualiza, _, _, _, _, err := ConsultasSql.ConsultaPrecioExistenciaYActualizaProductoEnAlmacen(operacion, movimiento, almacen, producto.ID.Hex(), float64(0), float64(1))
			if err != nil {
				fmt.Println(err)
				SDP.ID = operacion
				SDP.SEstado = false
				SDP.SMsj = "Ocurrió un problema al consultar bases de datos, verifique su conexión y vuelva a intentarlo."

				jData, _ := json.Marshal(SDP)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}

			if actualiza {
				html, calculadora := PuntoVentaModel.ConsultaDatosOperacion(operacion)

				SDP.ID = operacion
				SDP.SEstado = true
				SDP.SMsj = "Actualizacion a inventario satisfactoria"
				SDP.SIhtml = template.HTML(html)
				SDP.SCalculadora = template.HTML(calculadora)

				jData, _ := json.Marshal(SDP)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}

			SDP.SEstado = false
			SDP.SMsj = "No es posible surtir el producto."
			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return

		}
		//Busca en Elastic

		docs := ProductoModel.BuscarEnElastic(codigo)

		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)

			arrToMongo = []bson.ObjectId{}
			if numeroRegistros <= limitePorPagina {
				for _, v := range arrIDElastic[0:numeroRegistros] {
					arrToMongo = append(arrToMongo, v)
				}
			} else if numeroRegistros >= limitePorPagina {
				for _, v := range arrIDElastic[0:limitePorPagina] {
					arrToMongo = append(arrToMongo, v)
				}
			}

			MoConexion.FlushElastic()
			busqueda := PuntoVentaModel.GeneraTemplateBusquedaDeProductosEnElasticYAlmacenesParaPuntoVenta(ProductoModel.GetEspecifics(arrToMongo))
			SDP.SElastic = true
			SDP.ID = operacion
			SDP.SEstado = true
			SDP.SMsj = ""

			SDP.SBusqueda = template.HTML(busqueda)

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return

		}

		SDP.SEstado = false
		SDP.SMsj = "El Producto No Existe en la base de Datos."
		jData, _ := json.Marshal(SDP)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return

	}

	SDP.SEstado = false
	SDP.SMsj = "No se recibió un código, favor de verificar su conexión, e intente más tarde."
	jData, _ := json.Marshal(SDP)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//QuitarProducto recibe un id de un producto que existe en la vista para hace el rollback
func QuitarProducto(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto
	idOperacion := ctx.FormValue("Operacion")
	idProducto := ctx.FormValue("Producto")

	exito, err := ConsultasSql.EliminaProductoCarritoYActualizaInventarioAlmacen(idOperacion, idProducto)
	if err != nil {
		fmt.Println(err)
		SDP.ID = idProducto
		SDP.SEstado = false
		SDP.SMsj = "Ocurrió un problema al consultar bases de datos, verifique su conexión y vuelva a intentarlo."
		jData, _ := json.Marshal(SDP)

		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	if exito {
		html, calculadora := PuntoVentaModel.ConsultaDatosOperacion(idOperacion)
		SDP.ID = idOperacion
		SDP.SEstado = true
		SDP.SMsj = "Actualizacion a inventario satisfactoria"
		SDP.SIhtml = template.HTML(html)
		SDP.SCalculadora = template.HTML(calculadora)

		jData, _ := json.Marshal(SDP)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

}

//ModificarCantidad actualiza la tabla ventaTemporal con la cantidad deseada
func ModificarCantidad(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto
	var Operacion OperacionModel.OperacionMgo

	idOperacion := ctx.FormValue("Operacion")
	idProducto := ctx.FormValue("Producto")
	ValorNuevo := ctx.FormValue("Cantidad")
	ValorPrevio := ctx.FormValue("Previo")
	almacen := ctx.FormValue("Almacen")
	movimiento := ctx.FormValue("Movimiento")

	if bson.IsObjectIdHex(idOperacion) {

		Operacion = OperacionModel.GetOne(bson.ObjectIdHex(idOperacion))
		if MoGeneral.EstaVacio(Operacion) {
			SDP.ID = idOperacion
			SDP.SEstado = false
			SDP.SMsj = "No se pudo hacer referencia a la operación correspondiente, favor de recargar la página y vuelva a intentarlo. Disculpe las molestias."

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		VP, _ := strconv.ParseFloat(ValorPrevio, 64)
		VN, _ := strconv.ParseFloat(ValorNuevo, 64)

		actualiza, _, _, _, _, err := ConsultasSql.ConsultaPrecioExistenciaYActualizaProductoEnAlmacen(idOperacion, movimiento, almacen, idProducto, VP, VN)
		if err != nil {
			fmt.Println(err)
			SDP.ID = idOperacion
			SDP.SEstado = false
			SDP.SMsj = "Ocurrió un problema al consultar bases de datos, verifique su conexión y vuelva a intentarlo."

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		if actualiza {
			html, calculadora := PuntoVentaModel.ConsultaDatosOperacion(idOperacion)

			SDP.ID = idOperacion
			SDP.SEstado = true
			SDP.SMsj = "Actualizacion a inventario satisfactoria"
			SDP.SIhtml = template.HTML(html)
			SDP.SCalculadora = template.HTML(calculadora)

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		SDP.SEstado = false
		SDP.SMsj = "No es posible surtir el producto."
		jData, _ := json.Marshal(SDP)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return

	}

	SDP.ID = idOperacion
	SDP.SEstado = false
	SDP.SMsj = "No se pudo hacer referencia a la operación correspondiente, favor de recargar la página y vuelva a intentarlo. Disculpe las molestias."

	jData, _ := json.Marshal(SDP)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//ModificarImpuesto actualiza la tabla ventaTemporal con el nuevo impuesto
func ModificarImpuesto(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto

	idOperacion := ctx.FormValue("Operacion")
	Producto := ctx.FormValue("Producto")
	ValorNuevo := ctx.FormValue("ValorNuevo")
	ValorPrevio := ctx.FormValue("ValorPrevio")
	Almacen := ctx.FormValue("Almacen")
	Movimiento := ctx.FormValue("Movimiento")
	Price := ctx.FormValue("Precio")
	Factor := ctx.FormValue("Factor")
	Tipo := ctx.FormValue("Tipo")

	VP, _ := strconv.ParseFloat(ValorPrevio, 64)
	VN, _ := strconv.ParseFloat(ValorNuevo, 64)
	Precio, _ := strconv.ParseFloat(Price, 64)

	err := ConsultasSql.ActualizaImpuestoDeProductoEnAlmacenPsql(idOperacion, Movimiento, Factor, Tipo, Almacen, Producto, VP, VN, Precio)
	if err != nil {
		fmt.Println(err)
		SDP.SEstado = false
		SDP.SMsj = "Actualización de Impuesto Con Errores: " + err.Error()
	} else {
		err = ConsultasSql.ActualizaVentaTemporalPorImpuestos(idOperacion, Movimiento, Almacen, Producto, VN, Precio)
		if err != nil {
			SDP.SEstado = false
			SDP.SMsj = "Se Actualizó Impuesto, pero no se pudo actualizar Venta debido a que: " + err.Error()
		} else {
			SDP.SEstado = true
			SDP.SMsj = "Actualización de Impuesto Satisfactoria."
		}
	}

	SDP.ID = idOperacion
	html, calculadora := PuntoVentaModel.ConsultaDatosOperacion(idOperacion)
	SDP.SIhtml = template.HTML(html)
	SDP.SCalculadora = template.HTML(calculadora)

	jData, _ := json.Marshal(SDP)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//ModificarCantidadModal solicita una cantidad variada de artículos desde el punto de venta
func ModificarCantidadModal(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto
	var Operacion OperacionModel.OperacionMgo

	idOperacion := ctx.FormValue("Operacion")
	cliente := ctx.FormValue("Cliente")
	Productos := ctx.Request.Form["Productos[]"]
	Almacenes := ctx.Request.Form["Almacenes[]"]
	Cantidades := ctx.Request.Form["Cantidades[]"]

	var Movimientos []string

	if bson.IsObjectIdHex(idOperacion) {
		Operacion = OperacionModel.GetOne(bson.ObjectIdHex(idOperacion))
		if MoGeneral.EstaVacio(Operacion) {
			SDP.ID = idOperacion
			SDP.SEstado = false
			SDP.SMsj = "No se pudo hacer referencia a la operación correspondiente, favor de recargar la página y vuelva a intentarlo. Disculpe las molestias."

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		encontrado := false
		movimiento := ""

		for _, a := range Almacenes {
			encontrado = false
			for _, v := range Operacion.Movimientos {
				if v.AlmacenOrigen.Hex() == a { //Implica que el movimiento ya existe y ponemos una bandera y obtenemos el id del movimiento que encontramos
					movimiento = v.IDMovimiento.Hex()
					encontrado = true
				}
			}

			if !encontrado { //implica que tendremos que crear un nuevo movimiento en la operacion
				Movimiento := OperacionModel.MovimientoMgo{}
				Transaccion := OperacionModel.TransaccionMgo{}

				Movimiento.IDMovimiento = bson.NewObjectId()
				Movimiento.AlmacenOrigen = bson.ObjectIdHex(a)

				if !bson.IsObjectIdHex(cliente) {
					Movimiento.AlmacenDestino = bson.ObjectIdHex("58ed23518c649f28f0445013")
				} else {
					Movimiento.AlmacenDestino = bson.ObjectIdHex(cliente)
				}

				Movimiento.Estatus = bson.ObjectIdHex("58efc5c5d2b2131778e9c931")

				//Sólo existe una transaccion, porque sólo existe una ruta.
				Transaccion.IDTransaccion = bson.NewObjectId()
				Transaccion.AlmacenOrigen = bson.ObjectIdHex(a)

				if !bson.IsObjectIdHex(cliente) {
					Transaccion.AlmacenDestino = bson.ObjectIdHex("58ed23518c649f28f0445013")
				} else {
					Transaccion.AlmacenDestino = bson.ObjectIdHex(cliente)
				}
				Transaccion.Estatus = bson.ObjectIdHex("58efc5c5d2b2131778e9c931")

				Movimiento.Transacciones = append(Movimiento.Transacciones, Transaccion)
				Operacion.Movimientos = append(Operacion.Movimientos, Movimiento)
				movimiento = Movimiento.IDMovimiento.Hex()

				if !Operacion.ActualizaMgo([]string{"Movimientos"}, []interface{}{Operacion.Movimientos}) {
					SDP.ID = idOperacion
					SDP.SEstado = false
					SDP.SMsj = "No se pudo actualizar la Operación, por favor revisa tu conexión."

					jData, _ := json.Marshal(SDP)
					ctx.Header().Set("Content-Type", "application/json")
					ctx.Write(jData)
					return
				}
			}

			Movimientos = append(Movimientos, movimiento)
		}

		err := ConsultasSql.ConsultaPrecioExistenciaYActualizaProductoEnAlmacenModal(idOperacion, Movimientos, Almacenes, Productos, Cantidades)
		SDP.ID = idOperacion
		SDP.SEstado = true
		html, calculadora := PuntoVentaModel.ConsultaDatosOperacion(idOperacion)
		SDP.SIhtml = template.HTML(html)
		SDP.SCalculadora = template.HTML(calculadora)

		if err == nil {
			SDP.SMsj = "Actualizacion a inventario satisfactoria"
		} else {
			SDP.SMsj = "Ocurrieron errores al Actualizar las bases de datos"
		}
		jData, _ := json.Marshal(SDP)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return

	}

	SDP.ID = idOperacion
	SDP.SEstado = false
	SDP.SMsj = "No se pudo hacer referencia a la operación correspondiente, favor de recargar la página y vuelva a intentarlo. Disculpe las molestias."

	jData, _ := json.Marshal(SDP)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//ModificarCantidadModal2 actualiza la tabla ventaTemporal con la cantidad deseada
func ModificarCantidadModal2(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto
	var Operacion OperacionModel.OperacionMgo

	idOperacion := ctx.FormValue("Operacion")
	idProducto := ctx.FormValue("Producto")
	ValorNuevo := ctx.FormValue("Cantidad")
	ValorPrevio := ctx.FormValue("Previo")
	almacen := "5915e0a5e757702f60a21b78"
	movimiento := ""

	if bson.IsObjectIdHex(idOperacion) {

		Operacion = OperacionModel.GetOne(bson.ObjectIdHex(idOperacion))
		if !MoGeneral.EstaVacio(Operacion) {
			movimiento = Operacion.Movimientos[0].IDMovimiento.Hex()
		} else {
			SDP.ID = idOperacion
			SDP.SEstado = false
			SDP.SMsj = "No se pudo hacer referencia a la operación correspondiente, favor de recargar la página y vuelva a intentarlo. Disculpe las molestias."

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		VP, _ := strconv.ParseFloat(ValorPrevio, 64)
		VN, _ := strconv.ParseFloat(ValorNuevo, 64)

		actualiza, _, _, _, _, err := ConsultasSql.ConsultaPrecioExistenciaYActualizaProductoEnAlmacen(idOperacion, movimiento, almacen, idProducto, VP, VN)
		if err != nil {
			fmt.Println(err)
			SDP.ID = idOperacion
			SDP.SEstado = false
			SDP.SMsj = "Ocurrió un problema al consultar bases de datos, verifique su conexión y vuelva a intentarlo."

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		if actualiza {
			html, calculadora := PuntoVentaModel.ConsultaDatosOperacion(idOperacion)

			SDP.ID = idOperacion
			SDP.SEstado = true
			SDP.SMsj = "Actualizacion a inventario satisfactoria"
			SDP.SIhtml = template.HTML(html)
			SDP.SCalculadora = template.HTML(calculadora)

			jData, _ := json.Marshal(SDP)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

		SDP.SEstado = false
		SDP.SMsj = "No es posible surtir el producto."
		jData, _ := json.Marshal(SDP)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return

	}

	SDP.ID = idOperacion
	SDP.SEstado = false
	SDP.SMsj = "No se pudo hacer referencia a la operación correspondiente, favor de recargar la página y vuelva a intentarlo. Disculpe las molestias."

	jData, _ := json.Marshal(SDP)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//TraeOperacion Carga los datos de la operacion del id recibido
func TraeOperacion(ctx *iris.Context) {
	var SDP PuntoVentaModel.SDataProducto
	idOperacion := ctx.FormValue("Operacion")

	if bson.IsObjectIdHex(idOperacion) { //Es un ObjectId
		operacion := OperacionModel.GetOne(bson.ObjectIdHex(idOperacion))
		if operacion.ID.Hex() != "" { //No existe el ObjectId
			estatus := CatalogoModel.GetSubEspecificByFields("Valores._id", operacion.Estatus)
			if estatus.Valor == "PENDIENTE" { //La operacion debe estar aun PENDIENTE
				html, calculadora := PuntoVentaModel.ConsultaDatosOperacion(idOperacion)
				SDP.ID = idOperacion
				SDP.SEstado = true
				SDP.SMsj = "Actualizacion a inventario satisfactoria"
				SDP.SIhtml = template.HTML(html)
				SDP.SCalculadora = template.HTML(calculadora)

				jData, _ := json.Marshal(SDP)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}
			SDP.ID = idOperacion
			SDP.SEstado = false
			SDP.SMsj = "Esta Operacion ya se ha concretado."
			jData, _ := json.Marshal(SDP)

			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return

		}

	}
	SDP.ID = idOperacion
	SDP.SEstado = false
	SDP.SMsj = "No se encuentra la referencia a la Operacion solicitada."
	jData, _ := json.Marshal(SDP)

	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send PuntoVentaModel.SPuntoVenta

	Pagina := ctx.FormValue("Pag")
	if Pagina != "" {
		num, _ := strconv.Atoi(Pagina)
		if num == 0 {
			num = 1
		}
		NumPagina = num
		skip := limitePorPagina * (NumPagina - 1)
		limite := skip + limitePorPagina

		arrToMongo = []bson.ObjectId{}
		if NumPagina == paginasTotales {
			final := int(numeroRegistros) % limitePorPagina
			if final == 0 {
				for _, v := range arrIDElastic[skip:limite] {
					arrToMongo = append(arrToMongo, v)
				}
			} else {
				for _, v := range arrIDElastic[skip : skip+final] {
					arrToMongo = append(arrToMongo, v)
				}
			}

		} else {
			for _, v := range arrIDElastic[skip:limite] {
				arrToMongo = append(arrToMongo, v)
			}
		}

		Cabecera, Cuerpo := OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrToMongo))
		Send.SIndex.SCabecera = template.HTML(Cabecera)
		Send.SIndex.SBody = template.HTML(Cuerpo)

		Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, NumPagina)
		Send.SIndex.SPaginacion = template.HTML(Paginacion)

	} else {
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."

	}

	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Send.SEstado = true

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//MuestraIndexPorGrupo regresa template de busqueda y paginacion de acuerdo a la agrupacion solicitada
func MuestraIndexPorGrupo(ctx *iris.Context) {
	var Send PuntoVentaModel.SPuntoVenta
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.PuntoVenta.ENombrePuntoVenta.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := PuntoVentaModel.BuscarEnElastic(cadenaBusqueda)

		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)

			arrToMongo = []bson.ObjectId{}
			if numeroRegistros <= limitePorPagina {
				for _, v := range arrIDElastic[0:numeroRegistros] {
					arrToMongo = append(arrToMongo, v)
				}
			} else if numeroRegistros >= limitePorPagina {
				for _, v := range arrIDElastic[0:limitePorPagina] {
					arrToMongo = append(arrToMongo, v)
				}
			}

			Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
			}

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

			Send.SIndex.SRMsj = "No se encontraron resultados para: " + cadenaBusqueda + " ."
		}

	} else {

		if numeroRegistros <= limitePorPagina {
			Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = OperacionModel.GeneraTemplatesBusquedaParaPuntoDeVenta(OperacionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
		}

		Send.SIndex.SCabecera = template.HTML(Cabecera)
		Send.SIndex.SBody = template.HTML(Cuerpo)

		paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
		Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
		Send.SIndex.SPaginacion = template.HTML(Paginacion)

		Send.SIndex.SRMsj = "No se encontraron resultados para: " + cadenaBusqueda + " ."
	}
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Send.SEstado = true

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}
