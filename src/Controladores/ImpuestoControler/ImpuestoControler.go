package ImpuestoControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modulos/Session"

	"../../Modelos/CatalogoModel"
	"../../Modelos/ImpuestoModel"
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

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Impuesto
func IndexGet(ctx *iris.Context) {
	var Send ImpuestoModel.SImpuesto
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
	numeroRegistros = ImpuestoModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Impuestos := ImpuestoModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Impuestos {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(Impuestos[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(Impuestos[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("ImpuestoIndex.html", Send)
}

//IndexPost regresa la peticon post que se hizo desde el index de Impsuesto
func IndexPost(ctx *iris.Context) {

	var Send ImpuestoModel.SImpuesto
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
	//Send.Impuesto.EVARIABLEImpuesto.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := ImpuestoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("ImpuestoIndex.html", Send)
}

//###########################< ALTA >################################

//AltaGet renderea al alta de Impuesto
func AltaGet(ctx *iris.Context) {
	var Send ImpuestoModel.SImpuesto
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

	Send.Impuesto.ENombreImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(149, ""))
	Send.Impuesto.EClasificacionImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(148, ""))
	Send.Impuesto.ESubClasificacionImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(176, ""))
	Send.Impuesto.EImpuestosImpuesto.Impuestos.ETipoImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(150, ""))
	Send.Impuesto.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(151, ""))
	//template.HTML(CargaCombos.CargaComboCatalogo(151, ""))

	ctx.Render("ImpuestoAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Impuesto
func AltaPost(ctx *iris.Context) {

	var Send ImpuestoModel.SImpuesto
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

	EstatusPeticion := false

	var ImpuestoVista ImpuestoModel.Impuesto
	var ImpuestoMongo ImpuestoModel.ImpuestoMgo

	grupo := ctx.FormValue("Abreviatura")
	ImpuestoMongo.TipoImpuesto = bson.ObjectIdHex(grupo)
	ImpuestoVista.ENombreImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboTipoDeImpuestos(grupo))

	if grupo == "" {
		EstatusPeticion = true
		ImpuestoVista.ENombreImpuesto.IEstatus = true
		ImpuestoVista.ENombreImpuesto.IMsj = "Debe Seleccionar un grupo de Impuestos para dar de alta"
	}

	clasificacion := ctx.FormValue("Clasificacion")
	ImpuestoMongo.Clasificacion = bson.ObjectIdHex(clasificacion)
	ImpuestoVista.EClasificacionImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(148, clasificacion))

	if clasificacion == "" {
		EstatusPeticion = true
		ImpuestoVista.EClasificacionImpuesto.IEstatus = true
		ImpuestoVista.EClasificacionImpuesto.IMsj = "Debe Seleccionar un grupo de Impuestos para dar de alta"
	}

	nombre := ctx.FormValue("Nombre")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.Nombre = nombre
	ImpuestoMongo.Datos.Nombre = nombre

	if nombre == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.IMsj = "Debe capturar un nombre para este impuesto"
	}

	valor := ctx.FormValue("Valor")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.Valor, _ = strconv.ParseFloat(valor, 64)
	ImpuestoMongo.Datos.Max, _ = strconv.ParseFloat(valor, 64)

	if valor == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.IMsj = "Debe capturar un valor para este impuesto"
	}

	tipo := ctx.FormValue("Tipo")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.ETipoImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(150, tipo))
	ImpuestoMongo.Datos.TipoFactor = bson.ObjectIdHex(tipo)

	if tipo == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ETipoImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ETipoImpuestos.IMsj = "Debe capturar un valor para este impuesto"
	}

	unidad := ctx.FormValue("Unidad")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(151, unidad))
	ImpuestoMongo.Datos.Unidad = bson.ObjectIdHex(unidad)

	if unidad == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.IMsj = "Debe capturar un valor para este impuesto"
	}

	ImpuestoMongo.ID = bson.NewObjectId()
	ImpuestoMongo.Estatus = CatalogoModel.RegresaIDEstatusActivo(136)
	ImpuestoMongo.FechaHora = time.Now()

	Send.Impuesto = ImpuestoVista

	if EstatusPeticion {
		Send.SEstado = false
		Send.SMsj = "La validación indica que el objeto capturado no puede procesarse"
		ctx.Render("ImpuestoAlta.html", Send)
	} else {
		insertaElastic := ImpuestoMongo.InsertaElastic()
		fmt.Println("Insertado en elastic", insertaElastic)

		if ImpuestoMongo.InsertaMgo() {
			if ImpuestoMongo.InsertaElastic() {
				Send.SEstado = true
				Send.SMsj = "Se ha realizado una inserción exitosa"
				Send.Impuesto = ImpuestoModel.Impuesto{}
				Send.Impuesto.ENombreImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(149, ""))
				Send.Impuesto.EClasificacionImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(174, ""))
				Send.Impuesto.EImpuestosImpuesto.Impuestos.ETipoImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(150, ""))
				Send.Impuesto.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(151, ""))

				ctx.Render("ImpuestoAlta.html", Send)
			} else {
				Send.SEstado = false
				Send.SMsj = "Ocurrió un error al insertar en el servidor de busqueda, intente mas tarde"
				ctx.Render("ImpuestoAlta.html", Send)
			}

		} else {
			Send.SEstado = false
			Send.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("ImpuestoAlta.html", Send)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Impuesto
func EditaGet(ctx *iris.Context) {
	var SImp ImpuestoModel.SImpuesto

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SImp.SSesion.Name = NameUsrLoged
	SImp.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SImp.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SImp.SEstado = false
		SImp.SMsj = errSes.Error()
		ctx.Render("ZError.html", SImp)
		return
	}

	var aux ImpuestoModel.Impuesto

	id := ctx.Param("ID")

	if bson.IsObjectIdHex(id) {
		Impuesto := ImpuestoModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Impuesto) {
			aux.ID = Impuesto.ID
			aux.EClasificacionImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(148, Impuesto.Clasificacion.Hex()))
			aux.ENombreImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(149, Impuesto.TipoImpuesto.Hex()))
			aux.EImpuestosImpuesto.Impuestos.ETipoImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(150, Impuesto.Datos.TipoFactor.Hex()))
			aux.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(151, Impuesto.Datos.Unidad.Hex()))
			aux.EImpuestosImpuesto.Impuestos.ENombreImpuestos.Nombre = Impuesto.Datos.Nombre
			aux.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.Valor = Impuesto.Datos.Max
			aux.EEstatusImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(147, Impuesto.Estatus.Hex()))

			SImp.Impuesto = aux

			ctx.Render("ImpuestoEdita.html", SImp)

		} else {

			SImp.SEstado = false
			SImp.SMsj = "No se encontró el impuesto en la base, es posible que la conexión halla fallado, vuelva a intentar."
			ctx.Render("ImpuestoIndex.html", SImp)
		}
	} else {

		SImp.SEstado = false
		SImp.SMsj = "No se recibió un parámetro adcuado para mostrar el detalle del impuesto, favor intente de nuevo."
		ctx.Render("ImpuestoIndex.html", SImp)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Impuesto
func EditaPost(ctx *iris.Context) {
	var Send ImpuestoModel.SImpuesto

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

	EstatusPeticion := false

	var ImpuestoVista ImpuestoModel.Impuesto
	var ImpuestoMongo ImpuestoModel.ImpuestoMgo

	id := ctx.FormValue("ID")

	ImpuestoMongo.ID = bson.ObjectIdHex(id)
	ImpuestoMongo.FechaHora = time.Now()
	ImpuestoVista.ID = ImpuestoMongo.ID

	if !bson.IsObjectIdHex(id) {
		EstatusPeticion = true
		Send.SEstado = false
		Send.SMsj = "No se obtuvo una referencia adecuada para poder continuar con la edición del Impuesto, favor de intenta más tarde, disculpe la molestia."
		ctx.Render("ImpuestoIndex.html", Send)
	}

	ImpuestoAnterior := ImpuestoModel.GetOne(bson.ObjectIdHex(id))
	if MoGeneral.EstaVacio(ImpuestoAnterior) {
		EstatusPeticion = true
		Send.SEstado = false
		Send.SMsj = "El Objeto no existe o no se pudo hacer referencia a él para editarlo, verifique su conexión a la base de datos."
		ctx.Render("ImpuestoIndex.html", Send)
	}

	grupo := ctx.FormValue("Abreviatura")
	ImpuestoMongo.TipoImpuesto = bson.ObjectIdHex(grupo)
	ImpuestoVista.ENombreImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(149, grupo))

	if grupo == "" {
		EstatusPeticion = true
		ImpuestoVista.ENombreImpuesto.IEstatus = true
		ImpuestoVista.ENombreImpuesto.IMsj = "Debe Seleccionar un grupo de Impuestos para dar de alta."
	}

	clasificacion := ctx.FormValue("Clasificacion")
	ImpuestoMongo.Clasificacion = bson.ObjectIdHex(clasificacion)
	ImpuestoVista.EClasificacionImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(148, clasificacion))

	if clasificacion == "" {
		EstatusPeticion = true
		ImpuestoVista.EClasificacionImpuesto.IEstatus = true
		ImpuestoVista.EClasificacionImpuesto.IMsj = "Debe Seleccionar un grupo de Impuestos para dar de alta."
	}

	nombre := ctx.FormValue("Nombre")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.Nombre = nombre
	ImpuestoMongo.Datos.Nombre = nombre

	if nombre == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.IMsj = "Debe capturar un nombre para este impuesto."
	} else if nombre != ImpuestoAnterior.Datos.Nombre {
		if ImpuestoMongo.ConsultaExistenciaByFieldMgo("Impuestos.Nombre", nombre) {
			EstatusPeticion = true
			ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.IEstatus = true
			ImpuestoVista.EImpuestosImpuesto.Impuestos.ENombreImpuestos.IMsj = "El nombre de este impuesto ya existe, favor de utilizar otro."
		}
	}

	valor := ctx.FormValue("Valor")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.Valor, _ = strconv.ParseFloat(valor, 64)
	ImpuestoMongo.Datos.Max, _ = strconv.ParseFloat(valor, 64)

	if valor == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.IMsj = "Debe capturar un valor para este impuesto"
	}

	tipo := ctx.FormValue("Tipo")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.ETipoImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(150, tipo))
	ImpuestoMongo.Datos.TipoFactor = bson.ObjectIdHex(tipo)

	if tipo == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ETipoImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.ETipoImpuestos.IMsj = "Debe capturar un valor para este impuesto."
	}

	unidad := ctx.FormValue("Unidad")
	ImpuestoVista.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(151, unidad))
	ImpuestoMongo.Datos.Unidad = bson.ObjectIdHex(unidad)

	if unidad == "" {
		EstatusPeticion = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.IEstatus = true
		ImpuestoVista.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.IMsj = "Debe capturar un valor para este impuesto."
	}

	estatus := ctx.FormValue("Estatus")
	ImpuestoMongo.Estatus = bson.ObjectIdHex(estatus)
	ImpuestoVista.EEstatusImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(147, estatus))

	if estatus == "" {
		EstatusPeticion = true
		ImpuestoVista.EEstatusImpuesto.IEstatus = true
		ImpuestoVista.EEstatusImpuesto.IMsj = "Debe Seleccionar un estatus."
	}

	Send.Impuesto = ImpuestoVista

	if EstatusPeticion {
		Send.SEstado = false
		Send.SMsj = "La validación indica que el objeto capturado no puede procesarse"
		ctx.Render("ImpuestoEdita.html", Send)
	} else {

		if ImpuestoMongo.ReemplazaMgo() {

			Send.SEstado = true
			Send.SMsj = "Se ha realizado una actualización exitosa"
			ctx.Render("ImpuestoDetalle.html", Send)

		} else {
			Send.SEstado = false
			Send.SMsj = "Lamentablemente no se pudo Actualizar el Objeto, revise su conexión e intente de nuevo, disculpe la moestia."
			ctx.Render("ImpuestoEdita.html", Send)
		}

	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {

	var SImp ImpuestoModel.SImpuesto
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SImp.SSesion.Name = NameUsrLoged
	SImp.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SImp.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SImp.SEstado = false
		SImp.SMsj = errSes.Error()
		ctx.Render("ZError.html", SImp)
		return
	}

	var aux ImpuestoModel.Impuesto

	id := ctx.Param("ID")

	if bson.IsObjectIdHex(id) {
		Impuesto := ImpuestoModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Impuesto) {
			aux.ID = Impuesto.ID
			aux.EClasificacionImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(148, Impuesto.Clasificacion.Hex()))
			aux.ENombreImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(149, Impuesto.TipoImpuesto.Hex()))
			aux.EImpuestosImpuesto.Impuestos.ETipoImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(150, Impuesto.Datos.TipoFactor.Hex()))
			aux.EImpuestosImpuesto.Impuestos.EUnidadImpuestos.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(151, Impuesto.Datos.Unidad.Hex()))
			aux.EImpuestosImpuesto.Impuestos.ENombreImpuestos.Nombre = Impuesto.Datos.Nombre
			aux.EImpuestosImpuesto.Impuestos.EValorMaxImpuestos.Valor = Impuesto.Datos.Max
			aux.EEstatusImpuesto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(147, Impuesto.Estatus.Hex()))

			SImp.Impuesto = aux

			ctx.Render("ImpuestoDetalle.html", SImp)

		} else {

			SImp.SEstado = false
			SImp.SMsj = "No se encontró el impuesto en la base, es posible que la conexión halla fallado, vuelva a intentar."
			ctx.Render("ImpuestoIndex.html", SImp)
		}
	} else {

		SImp.SEstado = false
		SImp.SMsj = "No se recibió un parámetro adcuado para mostrar el detalle del impuesto, favor intente de nuevo."
		ctx.Render("ImpuestoIndex.html", SImp)
	}

}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send ImpuestoModel.SImpuesto
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

	ctx.Render("ImpuestoDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send ImpuestoModel.SImpuesto

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

		Cabecera, Cuerpo := ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrToMongo))
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
	var Send ImpuestoModel.SImpuesto
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Impuesto.ENombreImpuesto.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := ImpuestoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = ImpuestoModel.GeneraTemplatesBusqueda(ImpuestoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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

//ConsultarImpuestoPorGrupo recibe un parámetro de grupo de impuesto y regeresa todos los impuestos asociados a este impuesto
func ConsultarImpuestoPorGrupo(ctx *iris.Context) {
	var Send ImpuestoModel.SImpuesto

	tipo := ctx.FormValue("Tipo")

	if bson.IsObjectIdHex(tipo) {
		Send.SEstado = true
		Send.SIhtml = template.HTML(CargaCombos.CargaComboImpuestos(tipo, ""))
	} else {
		Send.SMsj = "No se recibió un parámetro para procesar la consulta de Impuestos."
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}
