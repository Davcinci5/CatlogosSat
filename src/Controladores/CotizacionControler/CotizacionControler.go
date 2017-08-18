package CotizacionControler

import (
	"encoding/json"
	"html/template"
	"strconv"
	"time"

	"../../Modulos/ConsultasSql"

	"../../Modelos/ProductoModel"

	"../../Modulos/Session"

	"fmt"

	"../../Modelos/ClienteModel"
	"../../Modelos/CotizacionModel"
	"../../Modelos/PersonaModel"
	"../../Modelos/UsuarioModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
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

//Numeros de Catalogos
var CatalogoFormasDePago = 168
var CatalogoDeFormaDeEnvio = 180

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Cotizacion
func IndexGet(ctx *iris.Context) {

	var Send CotizacionModel.SCotizacion

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
	numeroRegistros = CotizacionModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Cotizacions := CotizacionModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Cotizacions {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(Cotizacions[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(Cotizacions[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("CotizacionIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Cotizacion
func IndexPost(ctx *iris.Context) {

	var Send CotizacionModel.SCotizacion

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
	//Send.Cotizacion.EVARIABLECotizacion.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := CotizacionModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("CotizacionIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Cotizacion
func AltaGet(ctx *iris.Context) {

	var Send CotizacionModel.SCotizacion

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

	//####   TÚ CÓDIGO PARA CARGAR DATOS A LA VISTA DE ALTA----> PROGRAMADOR
	Send.Cotizacion.EFormaDePagoCotizacion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoFormasDePago, ""))
	Send.Cotizacion.EFormaDeEnvíoCotizacion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoDeFormaDeEnvio, ""))
	Send.Cotizacion.EUsuarioCotizacion.Usuario = UsuarioModel.GetIDByField("Usuario", NameUsrLoged)
	Send.Cotizacion.ID = bson.NewObjectId()
	ctx.Render("CotizacionAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Cotizacion
func AltaPost(ctx *iris.Context) {

	var Send CotizacionModel.SCotizacion
	var CotizacionMgo CotizacionModel.CotizacionMgo
	EstatusPeticion := false

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

	//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR
	// Send.Cotizacion.ID
	// Send.Cotizacion.EOperacionCotizacion
	// Send.Cotizacion.EUsuarioCotizacion.Usuario
	Cotizacion := ctx.FormValue("ID")
	Send.Cotizacion.ID = bson.ObjectIdHex(Cotizacion)
	if bson.IsObjectIdHex(Cotizacion) {
		CotizacionMgo = CotizacionModel.GetOne(bson.ObjectIdHex(Cotizacion))
		carrito, resumen := CotizacionModel.ConsultaDatosOperacion(Cotizacion)
		Send.Cotizacion.ECarritoCotizacion.Ihtml = template.HTML(carrito)
		Send.Cotizacion.EResumenCotizacion.Ihtml = template.HTML(resumen)
		if MoGeneral.EstaVacio(CotizacionMgo) {
			EstatusPeticion = true
			Send.SEstado = false
			Send.SMsj = "Agrega Productos al Carrito, Si el problema persiste, porfavor recarga la pagina"
		}
	} else {
		EstatusPeticion = true
		Send.SEstado = false
		Send.SMsj = "La referencia a la Cotizacion es incorrecta"
	}
	Cliente := ctx.FormValue("ClienteId")
	Nombre := ctx.FormValue("Nombre")
	Telefono := ctx.FormValue("Telefono")
	Correo := ctx.FormValue("Correo")
	Send.Cotizacion.ENombreCotizacion.Nombre = Nombre
	Send.Cotizacion.ETelefonoCotizacion.Telefono = Telefono
	Send.Cotizacion.ECorreoCotizacion.Correo = Correo

	if Cliente == "" {
		if Nombre == "" {
			EstatusPeticion = true
			Send.Cotizacion.ENombreCotizacion.IEstatus = true
			Send.Cotizacion.ENombreCotizacion.IMsj = "Campo obligatorio"
			Send.SEstado = false
			Send.SMsj = "Debe seleccionar un cliente, o llenar los campos minimos requeridos"
		} else {
			CotizacionMgo.Nombre = Nombre
		}

		if Telefono == "" && Correo == "" {
			EstatusPeticion = true
			Send.Cotizacion.ETelefonoCotizacion.IEstatus = true
			Send.Cotizacion.ETelefonoCotizacion.IMsj = "Campo obligatorio"
			Send.Cotizacion.ECorreoCotizacion.IEstatus = true
			Send.Cotizacion.ECorreoCotizacion.IMsj = "Campo obligatorio"
			Send.SEstado = false
			Send.SMsj = "Debe seleccionar un cliente, o llenar los campos minimos requeridos"
		} else {
			CotizacionMgo.Correo = Correo
			CotizacionMgo.Telefono = Telefono
		}

	} else {
		//Buscar el cliente y validar que exista
		if bson.IsObjectIdHex(Cliente) {
			ClienteEncontrado := ClienteModel.GetOne(bson.ObjectIdHex(Cliente))
			if !MoGeneral.EstaVacio(ClienteEncontrado) {
				busqueda := CotizacionModel.GeneraTemplateBusquedaClienteEspecifico(ClienteEncontrado)
				Send.Cotizacion.EClienteCotizacion.Ihtml = template.HTML(busqueda)
				CotizacionMgo.Cliente = bson.ObjectIdHex(Cliente)
			} else {
				EstatusPeticion = true
				Send.SEstado = false
				Send.SMsj = "Cliente No encontrado"
			}

		}

	}

	Cantidad := ctx.Request.Form["Producto"]
	if len(Cantidad) == 0 {
		EstatusPeticion = true
		Send.Cotizacion.EBuscarCotizacion.IEstatus = true
		Send.Cotizacion.EBuscarCotizacion.IMsj = "Debe agregar al menos un producto al carrito"
	}

	FormaPago := ctx.FormValue("FormaDePago")
	if bson.IsObjectIdHex(FormaPago) {
		Send.Cotizacion.EFormaDePagoCotizacion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoFormasDePago, FormaPago))
		CotizacionMgo.FormaDePago = bson.ObjectIdHex(FormaPago)
	} else {
		EstatusPeticion = true
		Send.Cotizacion.EFormaDePagoCotizacion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoFormasDePago, ""))
		Send.Cotizacion.EFormaDePagoCotizacion.IEstatus = true
		Send.Cotizacion.EFormaDePagoCotizacion.IMsj = "la referencia de la forma de pago es incorrecta"
	}

	FormaEnvio := ctx.FormValue("FormaDeEnvío")
	if bson.IsObjectIdHex(FormaPago) {
		Send.Cotizacion.EFormaDeEnvíoCotizacion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoDeFormaDeEnvio, FormaEnvio))
		CotizacionMgo.FormaDeEnvío = bson.ObjectIdHex(FormaPago)
	} else {
		EstatusPeticion = true
		Send.Cotizacion.EFormaDeEnvíoCotizacion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoDeFormaDeEnvio, ""))
		Send.Cotizacion.EFormaDeEnvíoCotizacion.IEstatus = true
		Send.Cotizacion.EFormaDeEnvíoCotizacion.IMsj = "la referencia de la forma de envío es incorrecta"
	}
	if !EstatusPeticion {
		var campos []string
		var valores []interface{}
		if bson.IsObjectIdHex(CotizacionMgo.Cliente.Hex()) {
			campos = []string{"Cliente", "Nombre", "Telefono", "Correo"}
			valores = []interface{}{CotizacionMgo.Cliente, CotizacionMgo.Nombre, CotizacionMgo.Telefono, CotizacionMgo.Correo}
		} else {
			campos = []string{"Nombre", "Telefono", "Correo"}
			valores = []interface{}{CotizacionMgo.Nombre, CotizacionMgo.Telefono, CotizacionMgo.Correo}
		}
		if CotizacionMgo.ActualizaMgo(campos, valores) {
			if CotizacionMgo.InsertaElastic() {
				Send.SEstado = true
				Send.SMsj = "Todo ha salido bien"
				ctx.Redirect("/Cotizacions", iris.StatusFound)
			}
			Send.SEstado = false
			Send.SMsj = "Existio un problema al guardar la Cotizacion en Elastic"
		} else {
			Send.SEstado = false
			Send.SMsj = "Existio un problema al guarar la Cotizacion en MogoDB"
		}
	}

	if Send.SMsj == "" {
		Send.SMsj = "LA validacion indica que algo salio mal"
	}
	fmt.Println("Errores")
	ctx.Render("CotizacionAlta.html", Send)

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Cotizacion
func EditaGet(ctx *iris.Context) {

	var Send CotizacionModel.SCotizacion

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

	ctx.Render("CotizacionEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Cotizacion
func EditaPost(ctx *iris.Context) {

	var Send CotizacionModel.SCotizacion

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

	//####   TÚ CÓDIGO PARA PROCESAR DATOS DE LA VISTA DE ALTA Y GUARDARLOS O REGRESARLOS----> PROGRAMADOR

	ctx.Render("CotizacionEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send CotizacionModel.SCotizacion

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

	//###### TU CÓDIGO AQUÍ PROGRAMADOR

	ctx.Render("CotizacionDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send CotizacionModel.SCotizacion

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

	//###### TU CÓDIGO AQUÍ PROGRAMADOR

	ctx.Render("CotizacionDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send CotizacionModel.SCotizacion

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

		Cabecera, Cuerpo := CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrToMongo))
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
	var Send CotizacionModel.SCotizacion
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Cotizacion.ENombreCotizacion.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := CotizacionModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = CotizacionModel.GeneraTemplatesBusqueda(CotizacionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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

//funciones Extras

// TraerClientes busca un cliente y retorna un template para la vista
func TraerClientes(ctx *iris.Context) {
	var Send CotizacionModel.SDataCliente

	//Aqui!  --> Validacion de permisos AJAx

	Cliente := ctx.FormValue("Cliente")
	fmt.Println(Cliente)

	docs := PersonaModel.BuscarEnElastic(Cliente + " +CLIENTE")
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
		busqueda := CotizacionModel.GeneraTemplateBusquedaClientes(PersonaModel.GetEspecifics(arrToMongo))
		if busqueda == "" {
			Send.SEstado = false
			Send.SMsj = "No se encontraron clientes"
		} else {
			Send.SIhtml = template.HTML(busqueda)
			Send.SEstado = true
			Send.SMsj = ""
		}

		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return

	}

	Send.SEstado = false
	Send.SMsj = "El Cliente No Existe en la base de Datos."
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//TraerCliente regresa un template para la tabla de cotizaciones, solo un cliente --AJAX
func TraerCliente(ctx *iris.Context) {
	var Send CotizacionModel.SDataCliente

	//Aqui!  --> Validacion de permisos AJAx

	Cliente := ctx.FormValue("Cliente")
	Send.SEstado = false
	Send.SMsj = "El Producto No Existe en la base de Datos."
	if Cliente != "" {
		if bson.IsObjectIdHex(Cliente) {
			ClienteEncontrado := ClienteModel.GetOne(bson.ObjectIdHex(Cliente))
			if !MoGeneral.EstaVacio(ClienteEncontrado) {
				busqueda := CotizacionModel.GeneraTemplateBusquedaClienteEspecifico(ClienteEncontrado)
				Send.SEstado = true
				Send.SMsj = ""
				Send.SIhtml = template.HTML(busqueda)

				jData, _ := json.Marshal(Send)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return

			}
			Send.SMsj = "Existe un error Con el Cliente, No encontrado en la Base de Datos"

		} else {
			Send.SMsj = "La referencia al objeto no es valida"
		}

	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

// TraerProducto Busca productos para la cotizacion AJAX
func TraerProducto(ctx *iris.Context) {
	var Send CotizacionModel.SDataProducto

	Cotizacion := ctx.FormValue("Operacion")
	if !bson.IsObjectIdHex(Cotizacion) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia a la operacion solicitada, porfavor recarga la pagina y vuelve a intentarlo."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	if MoGeneral.EstaVacio(CotizacionModel.GetOne(bson.ObjectIdHex(Cotizacion))) {
		Send.SEstado = false
		Send.SMsj = "No se pudo crear la Operación asociada a la cotización, verifique su conexión e intente de nuevo más tarde o actualice."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Cantidad := ctx.FormValue("Cantidad")
	ValorNuevo, err := strconv.ParseFloat(Cantidad, 64)
	if err != nil {
		Send.SEstado = false
		Send.SMsj = "La Cantidad especificada no es válida."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Almacen := ctx.FormValue("Almacen")
	if !bson.IsObjectIdHex(Almacen) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia al almacén."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Producto := ctx.FormValue("Producto")
	if !bson.IsObjectIdHex(Producto) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia al almacén."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	if ValorNuevo > 0 {
		sisi, err := ConsultasSql.AgregaACarritoCotizacion(Cotizacion, Almacen, Producto, ValorNuevo)
		if err != nil {
			if err != nil {
				Send.SEstado = false
				Send.SMsj = "Ocurrió un error al agregar Producto a Carrito."
				jData, _ := json.Marshal(Send)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}
		}
		if sisi {
			carrito, resumen := CotizacionModel.ConsultaDatosOperacion(Cotizacion)

			Send.ID = Cotizacion
			Send.SEstado = true
			Send.SMsj = "Actualizacion a inventario satisfactoria"
			Send.SIhtml = template.HTML(carrito)
			Send.SCalculadora = template.HTML(resumen)

			jData, _ := json.Marshal(Send)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "No puedes agregar 0."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

}

// TraerProductos Busca productos para la cotizacion AJAX
func TraerProductos(ctx *iris.Context) {
	var Send CotizacionModel.SDataProducto
	var Cotizacion CotizacionModel.CotizacionMgo

	_, _, Usuario, _ := Session.GetDataSessionAJAX(ctx)

	CotizacionID := ctx.FormValue("CotizacionID")
	Cliente := ctx.FormValue("Cliente")

	if bson.IsObjectIdHex(CotizacionID) {
		Operacion := CotizacionModel.GetOne(bson.ObjectIdHex(CotizacionID))
		if MoGeneral.EstaVacio(Operacion) { //Si aun no tiene nada
			Cotizacion.ID = bson.ObjectIdHex(CotizacionID)
			UsuarioOrigen := UsuarioModel.GetIDByField("Usuario", Usuario)
			Cotizacion.Usuario = UsuarioOrigen

			if bson.IsObjectIdHex(Cliente) {
				Cotizacion.Cliente = bson.ObjectIdHex(Cliente)
			}

			Cotizacion.FechaInicio = time.Now()

			if !Cotizacion.InsertaMgo() {
				Send.SEstado = false
				Send.SMsj = "Error al Almacenar Cotizacion"
			}

		}
	} else {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia a la operacion solicitada, porfavor recarga la pagina y vuelve a intentarlo."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	//Si el movimiento no es nuevo
	Producto := ctx.FormValue("Producto")

	if Producto != "" {
		// ProductoSol := ProductoModel.GetEspecificByFields("Codigos.Valores", Producto)
		// if !MoGeneral.EstaVacio(ProductoSol) {
		// 	Send.SElastic = false

		// 	existe, err := ConsultasSql.ConsultaExistenciaProductoEnAlmacenCotizacion(CotizacionID, ProductoSol.ID.Hex())
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		Send.SEstado = false
		// 		Send.SMsj = "Error al consultar en postgres"
		// 	} else {
		// 		if !existe {
		// 			//		ConsultasSql.InsertaProductoEnCotizacion(Cotizacion.ID, "592f1012e757701c5075c192", Producto)
		// 		}
		// 	}

		// 	jData, _ := json.Marshal(Send)
		// 	ctx.Header().Set("Content-Type", "application/json")
		// 	ctx.Write(jData)
		// 	return
		// }
		//Buscar en elastic
		docs := ProductoModel.BuscarEnElastic(Producto)
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
			busqueda := CotizacionModel.GeneraTemplateBusquedaDeProductos(ProductoModel.GetEspecifics(arrToMongo))
			Send.SBusqueda = template.HTML(busqueda)
			Send.SElastic = true
			Send.SMsj = "El producto no se encontro en algun almacen"
			Send.SEstado = true

			jData, _ := json.Marshal(Send)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}

	}
	Send.SEstado = false
	Send.SMsj = "No se recibio ningun nombre o código de Productos"
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

// ActualizaProductoCarrito actualiza un producto seleccionado del carrito mediante AJAX
func ActualizaProductoCarrito(ctx *iris.Context) {
	var Send CotizacionModel.SDataProducto

	Cotizacion := ctx.FormValue("Operacion")
	if !bson.IsObjectIdHex(Cotizacion) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia a la operacion solicitada, porfavor recarga la pagina y vuelve a intentarlo."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	if MoGeneral.EstaVacio(CotizacionModel.GetOne(bson.ObjectIdHex(Cotizacion))) {
		Send.SEstado = false
		Send.SMsj = "No se pudo crear la Operación asociada a la cotización, verifique su conexión e intente de nuevo más tarde o actualice."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Cantidad := ctx.FormValue("Cantidad")
	ValorNuevo, err := strconv.ParseFloat(Cantidad, 64)
	if err != nil {
		Send.SEstado = false
		Send.SMsj = "La Cantidad especificada no es válida."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Almacen := ctx.FormValue("Almacen")
	if !bson.IsObjectIdHex(Almacen) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia al almacén."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Producto := ctx.FormValue("Producto")
	if !bson.IsObjectIdHex(Producto) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia al almacén."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	fmt.Println(Producto, Almacen, Cotizacion, ValorNuevo)
	if ValorNuevo > 0 {
		SePudo, err := ConsultasSql.ConsultaPrecioExistenciaYActualizaProductoEnCotizacion(Cotizacion, Cotizacion, Almacen, Producto, ValorNuevo)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Ocurrio un error al actualizar la cotizacion (PSQL)."
			jData, _ := json.Marshal(Send)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}
		if SePudo {
			carrito, resumen := CotizacionModel.ConsultaDatosOperacion(Cotizacion)
			Send.ID = Cotizacion
			Send.SEstado = true
			Send.SMsj = "Actualizacion a inventario satisfactoria"
			Send.SIhtml = template.HTML(carrito)
			Send.SCalculadora = template.HTML(resumen)

			jData, _ := json.Marshal(Send)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}
		Send.SEstado = false
		Send.SMsj = "No se actualizo el carrito."

	} else {
		Send.SEstado = false
		Send.SMsj = "Valor invalido, debe ser mayor que 0."
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//QuitarProducto quita un producto del carrito de cotizacion
func QuitarProducto(ctx *iris.Context) {
	var Send CotizacionModel.SDataProducto

	Cotizacion := ctx.FormValue("Operacion")
	if !bson.IsObjectIdHex(Cotizacion) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia a la operacion solicitada, porfavor recarga la pagina y vuelve a intentarlo."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	if MoGeneral.EstaVacio(CotizacionModel.GetOne(bson.ObjectIdHex(Cotizacion))) {
		Send.SEstado = false
		Send.SMsj = "No se pudo crear la Operación asociada a la cotización, verifique su conexión e intente de nuevo más tarde o actualice."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Almacen := ctx.FormValue("Almacen")
	if !bson.IsObjectIdHex(Almacen) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia al almacén."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Producto := ctx.FormValue("Producto")
	if !bson.IsObjectIdHex(Producto) {
		Send.SEstado = false
		Send.SMsj = "No se pudo hacer referencia al almacén."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	fmt.Println(Producto, Almacen, Cotizacion)
	SePudo, err := ConsultasSql.EliminaProductoCarritoYActualizaInventarioCotizacion(Cotizacion, Cotizacion, Almacen, Producto)
	if err != nil {
		Send.SEstado = false
		Send.SMsj = "Ocurrio un error al Eliminar el producto de la cotizacion (PSQL)."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}
	if SePudo {
		carrito, resumen := CotizacionModel.ConsultaDatosOperacion(Cotizacion)
		Send.ID = Cotizacion
		Send.SEstado = true
		Send.SMsj = "Actualizacion a inventario satisfactoria"
		Send.SIhtml = template.HTML(carrito)
		Send.SCalculadora = template.HTML(resumen)

		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}
	Send.SEstado = false
	Send.SMsj = "No se actualizo el carrito."
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}
