package TrasladoAjusteControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modelos/CatalogoModel"
	"../../Modelos/OperacionModel"
	"../../Modelos/ProductoModel"
	"../../Modelos/UnidadModel"

	"../../Modulos/CargaCombos"
	"../../Modulos/General"
	"../../Modulos/Imagenes"
	"../../Modulos/Session"

	_ "github.com/bmizerany/pq"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

var catalogoTipo = 162
var catalogoEstatusProd = 161

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

//MovimientoAjusteGet Renderea a la vista Get de Ajustes y Traslados
func MovimientoAjusteGet(ctx *iris.Context) {
	// var almacenesExistentes []AlmacenModel.AlmacenMgo
	// almacenesExistentes = AlmacenModel.GetAll()
	//Crear una estuctura anidada SSesion en la V.
	//AjusteTraslado
	ctx.Render("OperacionAjuste.html", nil)
}

//MovimientoTrasladoGet Renderea a la vista Post de Ajustes y Traslados
func MovimientoTrasladoGet(ctx *iris.Context) {
	ctx.Render("OperacionTraslado.html", nil)
}

//movimiento

// MovimientosGet renderea la seleccion de movimientos para ajustes y traslados.
func MovimientosGet(ctx *iris.Context) {
	var SMovimientoAlmacen AlmacenModel.SMovimientoAlmacen
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SMovimientoAlmacen.SSesion.Name = NameUsrLoged
	SMovimientoAlmacen.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SMovimientoAlmacen.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SMovimientoAlmacen.SEstado = false
		SMovimientoAlmacen.SMsj = errSes.Error()
		ctx.Render("ZError.html", SMovimientoAlmacen)
		return
	}

	var Almacenes []AlmacenModel.AlmacenMgo
	var AllAlmacenes AlmacenModel.MovimientoAlmacenes

	Almacenes, err := AlmacenModel.GetEspecificsByTagAndTestConexion("Clasificacion", bson.ObjectIdHex("58e5692ee75770120c60befa"))
	var cadena string
	for _, value := range err {
		cadena += value + "<br>"
	}

	if len(Almacenes) > 0 {
		if len(err) > 0 {
			SMovimientoAlmacen.SEstado = false
			SMovimientoAlmacen.SMsj = "Se encontraron errores de conexion en los siguients almacenes	" + cadena
		} else {
			SMovimientoAlmacen.SEstado = true
		}
		//Almacenes = AlmacenModel.AlmacenFields()
		AllAlmacenes.Almacenes = Almacenes

		exito := ctx.URLParam("exito")
		if exito == "1" {
			SMovimientoAlmacen.SEstado = true
			SMovimientoAlmacen.SMsj = "Movimiento Realizado con Exito"
		}
		var SAlmacenes []AlmacenModel.Almacen
		var SAlmacen AlmacenModel.Almacen
		var templateAlmacenes string

		for _, v := range Almacenes {
			SAlmacen.ENombreAlmacen.Nombre = v.Nombre
			SAlmacen.ID = v.ID
			SAlmacen.ETipoAlmacen.Tipo = v.Tipo
			SAlmacenes = append(SAlmacenes, SAlmacen)
			var tipoString string
			tipoString = CargaCombos.CargaCatalogoByID(132, v.Tipo)
			// fmt.Printf("\n Id -> %v    Tipo -> %v \n ", v.Tipo, tipoString)
			if tipoString != "" {
				templateAlmacenes += fmt.Sprintf(`<option value="%v">%v</option>`, v.ID.Hex(), v.Nombre)
			}
		}

		SMovimientoAlmacen.Almacenes = SAlmacenes
		SMovimientoAlmacen.Ihtml = template.HTML(templateAlmacenes)
	} else {
		SMovimientoAlmacen.SMsj = "No se ha podido cargar los almacenes..." + cadena
		SMovimientoAlmacen.SEstado = false
	}
	ctx.Render("TrasladoAjuste.html", SMovimientoAlmacen)
}

/*
// MovimientosPost renderea la seleccion de movimientos para ajustes y traslados.
func MovimientosPost(ctx *iris.Context) {

	var almOrigen string
	var almDestino string
	var tipo string
	var template = ``
	var isAjuste bool
	var whoIs string

	ctx.Request.ParseForm()

	for i, v := range ctx.Request.Form {
		switch i {
		case "origen":
			for _, valor := range v {
				almOrigen = valor
			}
			break
		case "destino":
			for _, valor := range v {
				almDestino = valor
			}
			break
		case "tipo":
			for _, valor := range v {
				tipo = valor
			}
			break
		}
	}

	fmt.Println("Origen ->", almOrigen)
	fmt.Println("Destino ->", almDestino)
	fmt.Println("Tipo ->", tipo)

	if tipo == "ajuste" {
		isAjuste = true
		whoIs = "ajuste"
	} else {
		whoIs = "traslado"

	}

	template += `	<div class="page-header">`
	if isAjuste {
		template += `	<h3 class="text-center">Ajuste : </h3>`
	} else {
		template += `	<h3 class="text-center">Traslado : </h3>`
	}
	template += fmt.Sprintf(`
					</div>


					<div class="col-sm-12">
						<div class="checkbox">
							<label>
								<input id="agregarcantidad" name="agregarcantidad" type="checkbox">
								Agregar cantidad al insertar codigo de barra
							</label>
						</div>
					</div>

					<div class="input-group input-group-md">
							<span class="input-group-addon">Buscar Articulo:</span>
							<input type="text" name="elarticulo" onKeydown="Javascript: if (event.keyCode==13) buscarProducto();"  id="elarticulo" class="form-control selectpicker" autofocus>
							<input type="text" hidden value="%v" id="tipomovimiento" name="tipomovimiento">

							<span class="input-group-addon">
								<buttont type="button" class="btn btn-primary" onClick="buscarProducto();" > <span class="glyphicon glyphicon-search"></span>  Buscar</button>
							</span>
					</div>`, whoIs)
	fmt.Println("Xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Println(template)
	fmt.Println("Xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Fprintf(ctx.ResponseWriter, template)

}
*/

//MovimientosAjustePost funcion  ase
func MovimientosAjustePost(ctx *iris.Context) {
	var almOrigen string
	var almDestino string
	var tipoMovimiento string
	var numeroArticulos string
	var cantidadOn bool

	almOrigen = ctx.FormValue("origen")
	almDestino = ctx.FormValue("destino")
	tipoMovimiento = ctx.FormValue("tipomovimiento")
	numeroArticulos = ctx.FormValue("articulos_agregados")
	cantidadOn, _ = strconv.ParseBool(ctx.FormValue("agregarcantidad"))
	ID := ctx.FormValue("ID")

	var Producto ProductoModel.ProductoMgo
	var Codigos ProductoModel.CodigoMgo

	Producto = ProductoModel.GetOne(bson.ObjectIdHex(ID)) // ProductoModel.GetProductoMongo(codigoBarra)

	Codigos = Producto.Codigos

	if Producto.ID == "" {
		fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Articulo no encontrado</h3>`)
	} else {

		var ProductoOrigen ProductoModel.ProductoPgrs
		var ProductoDestino ProductoModel.ProductoPgrs

		ProductoOrigen, ProductoDestino = ProductoModel.GetProductoPostgres(Producto.ID.Hex(), almOrigen, almDestino)
		ExistenciaOrigenString := strconv.FormatFloat(ProductoOrigen.Existencia, 'f', -1, 64)
		ExistenciaDestinoString := strconv.FormatFloat(ProductoDestino.Existencia, 'f', -1, 64)
		PrecioOrigenString := strconv.FormatFloat(ProductoOrigen.Precio, 'f', -1, 64)
		PrecioDestinoString := strconv.FormatFloat(ProductoDestino.Precio, 'f', -1, 64)
		///codigo que forma la fila para ajuste
		tAjuste := `<tr class="renglon " id="row-` + Producto.ID.Hex() + `">
							<td id="codigo_b` + numeroArticulos + `" hidden>` + Producto.ID.Hex() + `</td>`

		for _, v := range Codigos.Valores {
			tAjuste += `<td class="codigosdearticulos" id="` + Producto.ID.Hex() + `:input" hidden>` + v + `</td>`
		}

		tAjuste += `		<td id="desc_b` + numeroArticulos + `">` + Producto.Nombre + `</td>
							<td id="precio_b` + numeroArticulos + `">` + PrecioOrigenString + `</td> 
							<td id="origen_b` + numeroArticulos + `">` + ExistenciaOrigenString + `</td>
							
							<td id="operacion_b` + numeroArticulos + `"> 

							<label class="radio-inline"><input type="radio" name="operacion` + Producto.ID.Hex() + `" id="operacion` + Producto.ID.Hex() + `" value ="sumar" checked>Suma</label>
							<label class="radio-inline"><input type="radio" name="operacion` + Producto.ID.Hex() + `" id="operacion` + Producto.ID.Hex() + `" value ="restar">Resta</label>

							</td>`

		if cantidadOn {
			tAjuste += `<td id="cantidad_b` + numeroArticulos + `"><input type="text" name="cantidad` + Producto.ID.Hex() + `" id="` + Producto.ID.Hex() + `" value="0" requerided  pattern="\d*" onKeydown="Javascript: if (event.keyCode==13) FocusInputCodigos();" autofocus="true"></td>
						<td hidden> <script> $("#` + Producto.ID.Hex() + `").select();  $("#` + Producto.ID.Hex() + `").focus(); </script></td>`
		} else {

			tAjuste += `	<td id="cantidad_b` + numeroArticulos + `"><input type="text" name="cantidad` + Producto.ID.Hex() + `" id="` + Producto.ID.Hex() + `" value="0" requerided  pattern="\d*" onKeydown="Javascript: if (event.keyCode==13) FocusInputCodigos();"></td>`
		}

		tAjuste += `    <td><span class="btn btn-danger eliminar glyphicon glyphicon-remove"></span></td>
						</tr>`

		//codigo que forma la fila para traslado
		tTraslado := `<tr class="renglon"  id="row-` + Producto.ID.Hex() + `">
							<td id="codigo_b` + numeroArticulos + `" hidden>` + Producto.ID.Hex() + `</td>`

		for _, v := range Codigos.Valores {
			tTraslado += `<td class="codigosdearticulos" id="` + Producto.ID.Hex() + `:input" hidden>` + v + `</td>`
		}

		tTraslado += `	<td id="desc_b` + numeroArticulos + `">` + Producto.Nombre + `</td>
							<td id="precio_a` + numeroArticulos + `">` + PrecioOrigenString + `</td> 
							<td id="precio_b` + numeroArticulos + `">` + PrecioDestinoString + `</td> 
							<td id="origen_b` + numeroArticulos + `">` + ExistenciaOrigenString + `</td>
							<td id="destino_b` + numeroArticulos + `">` + ExistenciaDestinoString + `</td>`
		if cantidadOn {
			tTraslado += `<td id="cantidad_b` + numeroArticulos + `"><input type="text" name="cantidad` + Producto.ID.Hex() + `" id="` + Producto.ID.Hex() + `" value="0" requerided  pattern="\d*" onKeydown="Javascript: if (event.keyCode==13) FocusInputCodigos();"></td>
						  <td hidden> <script> $("#` + Producto.ID.Hex() + `").focus(); </script></td>`

		} else {

			tTraslado += `	<td id="cantidad_b` + numeroArticulos + `"><input type="text" name="cantidad` + Producto.ID.Hex() + `" id="` + Producto.ID.Hex() + `" value="0" requerided  pattern="\d*" onKeydown="Javascript: if (event.keyCode==13) FocusInputCodigos();"></td>`
		}

		tTraslado += `		<td><span class="btn btn-danger eliminar glyphicon glyphicon-remove"></span></td>
						</tr>`
		switch {
		//1
		case ProductoDestino.Estatus == "ACTIVO":
			switch {
			case ProductoOrigen.Estatus == "ACTIVO":
				switch {
				case tipoMovimiento == "ajuste":
					fmt.Fprintf(ctx.ResponseWriter, tAjuste)
					break
				case tipoMovimiento == "traslado":
					fmt.Fprintf(ctx.ResponseWriter, tTraslado)
					break
				}
				break
			case ProductoOrigen.Estatus == "INACTIVO":
				fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Origen: Inactivo     Destino: Activo</h3>`)
				break
			case ProductoOrigen.Estatus == "BLOQUEADO":
				fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Origen: Bloqueado      Destino: Activo</h3>`)
				break
			case ProductoOrigen.Estatus == "":
				fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Origen: No Existe     Destino: Activo</h3>`)

			}
			break
		//2
		case ProductoOrigen.Estatus == "ACTIVO":
			switch {
			case ProductoDestino.Estatus == "ACTIVO":
				switch {
				case tipoMovimiento == "ajuste":
					fmt.Fprintf(ctx.ResponseWriter, tAjuste)
					break
				case tipoMovimiento == "traslado":
					fmt.Fprintf(ctx.ResponseWriter, tTraslado)
					break
				}
				break
			case ProductoDestino.Estatus == "INACTIVO":
				fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Origen: Activo      Destino: Inactivo</h3>`)
				break
			case ProductoDestino.Estatus == "BLOQUEADO":
				fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Origen: Activo      Destino: Bloqueado</h3>`)
				break
			case ProductoDestino.Estatus == "":
				switch {
				case tipoMovimiento == "ajuste":
					fmt.Fprintf(ctx.ResponseWriter, tAjuste)
					break
				case tipoMovimiento == "traslado":
					fmt.Fprintf(ctx.ResponseWriter, tTraslado)
					break
				}
				break
			}
			break
		//3
		case ProductoDestino.Estatus == "":
			switch {
			case ProductoOrigen.Estatus == "":
				fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Origen: No Existe      Destino: No Existe</h3>`)
				break
			}
			break
		//4
		case ProductoOrigen.Estatus == "":
			switch {
			case ProductoDestino.Estatus == "":
				fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Origen: No Existe      Destino: No Existe</h3>`)
				break
			}
			break
		}
	}
}

//RealizarMovimiento funccion que realiza el movimiento
func RealizarMovimiento(ctx *iris.Context) {
	ctx.Request.ParseForm()

	var Operacion OperacionModel.OperacionMgo
	Operacion.ID = bson.NewObjectId()
	Operacion.UsuarioOrigen = bson.NewObjectId()
	Operacion.UsuarioDestino = bson.NewObjectId()
	Operacion.FechaHoraRegistro = time.Now()

	var Movimiento OperacionModel.MovimientoMgo
	Movimiento.IDMovimiento = bson.NewObjectId()

	AlmacenOrigen := ctx.Request.FormValue("origen")
	AlmacenDestino := ctx.Request.FormValue("destino")
	var operaciones = []string{}
	codigos := ctx.Request.PostForm["codigos[]"]
	existencias := ctx.Request.PostForm["existencias[]"]
	if AlmacenOrigen == AlmacenDestino {
		operaciones = ctx.Request.PostForm["operaciones[]"]
	}
	cantidades := ctx.Request.PostForm["cantidades[]"]

	lim := len(codigos)

	var TotalAlmacenOrigen float64
	var TotalAlmacenDestino float64
	if AlmacenDestino == AlmacenOrigen {
		for i := 0; i < lim; i++ {

			//nombresito := nombres[i]
			IDProducto := codigos[i]
			cantidad := cantidades[i]
			existe := existencias[i]
			operacion := operaciones[i]

			ExistenciaOrigen, _ := strconv.ParseFloat(existe, 64)
			cantidadSolicitada, _ := strconv.ParseFloat(cantidad, 64)

			if operacion == "sumar" {
				TotalAlmacenOrigen = cantidadSolicitada + ExistenciaOrigen

			} else if operacion == "restar" {
				TotalAlmacenOrigen = ExistenciaOrigen - cantidadSolicitada
			}

			if TotalAlmacenOrigen < 0 {
				fmt.Fprintf(ctx.ResponseWriter, "no")
			} else {
				opExitosa, errores := ProductoModel.RealizarAjuste(AlmacenOrigen, IDProducto, TotalAlmacenOrigen)
				fmt.Println("Operacion Exitosa", opExitosa)
				fmt.Println(errores)
			}
		}
	} else {
		for i := 0; i < lim; i++ {
			IDProducto := codigos[i]
			cantidad := cantidades[i]
			cantidadSolicitada, err := strconv.ParseFloat(cantidad, 64)
			if err != nil {
				fmt.Println("Error en la conversion: ", err)
			}

			ProductoOrigen, ProductoDestino := ProductoModel.GetProductoPostgres(IDProducto, AlmacenOrigen, AlmacenDestino)

			TotalAlmacenDestino = ProductoDestino.Existencia + cantidadSolicitada
			TotalAlmacenOrigen = ProductoOrigen.Existencia - cantidadSolicitada
			if TotalAlmacenOrigen < 0 {
				fmt.Fprintf(ctx.ResponseWriter, "no")
			} else {
				opExitosa, errores := ProductoModel.RealizarTraslado(AlmacenOrigen, AlmacenDestino, IDProducto, TotalAlmacenOrigen, TotalAlmacenDestino)
				if !opExitosa {
					fmt.Println(errores)
				}
			}

		}
	}
}

//SDataProducto estructura de Productos para la vista
type SDataProducto struct {
	SEstado bool
	SMsj    string
	SIhtml  template.HTML
}

//ConsultarProductos Consulta productos con sus respectivos impuestos
//Para armar la vista que conforma los productos, los precios e impuestos dinamicamente
func ConsultarProductos(ctx *iris.Context) {
	cadenaBusqueda = ctx.FormValue("nombreProducto")
	Send := SDataProducto{}

	if cadenaBusqueda != "" {
		docs := ProductoModel.BuscarEnElastic(cadenaBusqueda)
		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)
			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)

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

			productos := (ProductoModel.GetEspecifics(arrToMongo))

			tabla := GeneraTabla(productos)
			Send.SEstado = true
			Send.SIhtml = template.HTML(tabla)

		} else {
			Send.SEstado = false
			Send.SMsj = "No se encontraron resultados."
		}

	} else {
		Send.SEstado = false
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
	}
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
}

//GeneraTabla genera una tabla utilizada en el modulo de ventas
func GeneraTabla(productos []ProductoModel.ProductoMgo) string {
	tabla := ``
	for _, producto := range productos {
		if bson.IsObjectIdHex(producto.ID.Hex()) {

			htmlProducto := `<tr id="filaProducto-` + producto.ID.Hex() + `" class="rowsPro">`
			htmlProducto += `<td><label>` + producto.Nombre + `</label></td>`
			htmlProducto += `<td>`
			for _, value := range producto.Codigos.Valores {
				htmlProducto += `<label>` + value + `</label>`
			}
			htmlProducto += `</td>`
			tipo := CatalogoModel.RegresaNombreSubCatalogo(producto.Tipo)
			if tipo == "" {
				tipo = "Tipo no encontrado"
			}
			htmlProducto += `<td><label>` + tipo + `</label></td>`
			unidad := UnidadModel.RegresaNombreUnidad(producto.Unidad)
			if unidad == "" {
				unidad = "Unidad no encontrado"
			}
			htmlProducto += `<td><label>` + unidad + `</label></td>`
			estatus := CatalogoModel.RegresaNombreSubCatalogo(producto.Estatus)
			if estatus == "" {
				estatus = "Estatus desconocido"
			}
			htmlProducto += `<td><label>` + estatus + `</label></td>`
			htmlProducto += `<td>`
			for _, value := range producto.Etiquetas {
				htmlProducto += `<label>` + value + `</label>`
			}
			htmlProducto += `</td>`
			htmlImagenes, _ := Imagenes.CargaTemplateImagenesCompras(producto.Imagenes)
			htmlProducto += `<td>` + htmlImagenes + `</td>`
			htmlProducto += `<td><button type="button" class="btn btn-info" onClick='AgregarArticulo("` + producto.ID.Hex() + `")'>`
			htmlProducto += `<span class="glyphicon glyphicon-ok-sign">Seleccionar</span>`
			htmlProducto += `</button></td>`
			htmlProducto += `</tr>`
			tabla += htmlProducto
		}
	}
	return tabla
}
