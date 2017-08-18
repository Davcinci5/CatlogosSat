package GrupoControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/GrupoModel"
	"../../Modelos/GrupoPersonaModel"

	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
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

var tipoConsultaBase = 1

//limitePorPagina limite de registros a mostrar en la pagina
var limitePorPagina = 10

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Grupo
func IndexGet(ctx *iris.Context) {

	var Send GrupoModel.SGrupo

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
	numeroRegistros = GrupoModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Grupos := GrupoModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Grupos {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(Grupos[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(Grupos[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("GrupoIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Grupo
func IndexPost(ctx *iris.Context) {

	var Send GrupoModel.SGrupo

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
	//Send.Grupo.EVARIABLEGrupo.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := GrupoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("GrupoIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Grupo
func AltaGet(ctx *iris.Context) {

	var Send GrupoModel.SGrupo

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
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Send.Grupo.ETipoGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(179, ""))
	//####   TÚ CÓDIGO PARA CARGAR DATOS A LA VISTA DE ALTA----> PROGRAMADOR

	ctx.Render("GrupoAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Grupo
func AltaPost(ctx *iris.Context) {

	var Send GrupoModel.SGrupo
	var SendMgo GrupoModel.GrupoMgo
	EstatusPeticion := false
	var Cabecera, Cuerpo, Type string
	var Encontrados int

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

	Nombre := ctx.FormValue("Nombre")
	Send.ENombreGrupo.Nombre = Nombre
	SendMgo.Nombre = Nombre
	if Nombre == "" {
		EstatusPeticion = true
		Send.ENombreGrupo.IEstatus = true
		Send.ENombreGrupo.IMsj = "El campo Nombre es obligatorio."
	}

	if SendMgo.ConsultaExistenciaByFieldMgo("Nombre", Nombre) {
		EstatusPeticion = true
		Send.ENombreGrupo.IEstatus = true
		Send.ENombreGrupo.IMsj = "El Nombre de este grupo ya existe, favor de intentar con otro."
	}

	Descripcion := ctx.FormValue("Descripcion")
	Send.EDescripcionGrupo.Descripcion = Descripcion
	SendMgo.Descripcion = Descripcion
	if Descripcion == "" {
		EstatusPeticion = true
		Send.EDescripcionGrupo.IEstatus = true
		Send.EDescripcionGrupo.IMsj = "El campo Descripción es obligatorio."
	}

	Tipo := ctx.FormValue("Tipo")
	Send.Grupo.ETipoGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(179, Tipo))
	if !bson.IsObjectIdHex(Tipo) {
		EstatusPeticion = true
		Send.ETipoGrupo.IEstatus = true
		Send.ETipoGrupo.IMsj = "El campo Tipo es obligatorio."
	} else {
		SendMgo.Tipo = bson.ObjectIdHex(Tipo)
		Type = CatalogoModel.ObtenerValoresCatalogoPorValor(bson.ObjectIdHex(Tipo))
		if Type == "" {
			EstatusPeticion = true
			Send.ETipoGrupo.IEstatus = true
			Send.ETipoGrupo.IMsj = "El Tipo seleccionado no coincide con los valores actuales de la base de datos, favor de recargar su página."
		}
	}

	Formulario := ctx.Request.Form
	Seleccionados := Formulario["Objeto"]
	fmt.Println("Seleccionados: ", Seleccionados)
	if len(Seleccionados) > 0 {
		Ids := GrupoPersonaModel.ConvierteAObjectIDS(Seleccionados)
		fmt.Println("Ids: ", Ids)
		if len(Seleccionados) != len(Ids) {
			EstatusPeticion = true
			Send.EMiembrosGrupo.IEstatus = true
			Send.EMiembrosGrupo.IMsj = "Hay insonsistencias en los objetos seleccionados, no son los mismos los que se asignan a los que existen en la base de Datos, favor de recargar su página."
		} else {
			SendMgo.Miembros = Ids

			Cabecera, Cuerpo, Encontrados = GrupoModel.GeneraTemplatesBusquedaObjetos(SendMgo.Miembros, Type, true)

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			paginasTotales = MoGeneral.Totalpaginas(Encontrados, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)
		}
	} else {
		EstatusPeticion = true
		Send.EMiembrosGrupo.IEstatus = true
		Send.EMiembrosGrupo.IMsj = "Debe seleccionar al menos un objeto a guardar en este Grupo."
	}

	SendMgo.Estatus = CargaCombos.CargaEstatusActivoEnAlta(146)
	SendMgo.FechaHora = time.Now()

	if EstatusPeticion {
		Send.SEstado = false                                                              //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		Send.SMsj = "La validación indica que el objeto capturado no puede darse de alta" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
	} else {
		SendMgo.ID = bson.NewObjectId()
		if SendMgo.InsertaMgo() {
			if SendMgo.InsertaElastic() {
				Send.SEstado = true
				Send.SMsj = "Se ha realizado una inserción exitosa"
				ctx.Redirect("/Grupos/detalle/"+SendMgo.ID.Hex(), 301)
			} else {
				if !SendMgo.EliminaByIDMgo() {
					Send.SEstado = false
					Send.SMsj = "No se inserto en elastic, no se pudo eliminar el registro en MongoDb"
				}
			}
		} else {
			Send.SEstado = false
			Send.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
		}
	}
	ctx.Render("GrupoAlta.html", Send)
}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Grupo
func EditaGet(ctx *iris.Context) {

	var Send GrupoModel.SGrupo
	var Cabecera, Cuerpo string
	var Encontrados int

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

	ID := ctx.Param("ID")
	if bson.IsObjectIdHex(ID) {
		Grupo := GrupoModel.GetOne(bson.ObjectIdHex(ID))
		if Grupo.ID.Hex() != "" {
			Send.Grupo.ID = Grupo.ID
			Send.Grupo.ENombreGrupo.Nombre = Grupo.Nombre
			Send.Grupo.EDescripcionGrupo.Descripcion = Grupo.Descripcion
			Send.Grupo.ETipoGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(179, Grupo.Tipo.Hex()))
			Type := CatalogoModel.ObtenerValoresCatalogoPorValor(Grupo.Tipo)

			Cabecera, Cuerpo, Encontrados = GrupoModel.GeneraTemplatesBusquedaObjetos(Grupo.Miembros, Type, true)

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			paginasTotales = MoGeneral.Totalpaginas(Encontrados, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)
			Send.Grupo.EEstatusGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(146, Grupo.Estatus.Hex()))

			Send.Grupo.EFechaHoraGrupo.FechaHora = Grupo.FechaHora
			Send.SEstado = true
		} else {
			Send.SEstado = false
			Send.SMsj = "Grupo no encontrado"
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "Error en la referencia al Grupo"
	}

	ctx.Render("GrupoEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Grupo
func EditaPost(ctx *iris.Context) {

	var Send GrupoModel.SGrupo
	var SendMgo GrupoModel.GrupoMgo
	EstatusPeticion := false
	var Cabecera, Cuerpo string
	var Encontrados int
	Type := ``

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

	Nombre := ctx.FormValue("Nombre")
	Send.ENombreGrupo.Nombre = Nombre
	SendMgo.Nombre = Nombre
	if Nombre == "" {
		EstatusPeticion = true
		Send.ENombreGrupo.IEstatus = true
		Send.ENombreGrupo.IMsj = "El campo Nombre es obligatorio."
	}

	Descripcion := ctx.FormValue("Descripcion")
	Send.EDescripcionGrupo.Descripcion = Descripcion
	SendMgo.Descripcion = Descripcion
	if Descripcion == "" {
		EstatusPeticion = true
		Send.EDescripcionGrupo.IEstatus = true
		Send.EDescripcionGrupo.IMsj = "El campo Descripción es obligatorio."
	}

	Tipo := ctx.FormValue("Tipo")
	Send.Grupo.ETipoGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(179, Tipo))
	if !bson.IsObjectIdHex(Tipo) {
		EstatusPeticion = true
		Send.ETipoGrupo.IEstatus = true
		Send.ETipoGrupo.IMsj = "El campo Tipo es obligatorio."
	} else {
		SendMgo.Tipo = bson.ObjectIdHex(Tipo)
		Type = CatalogoModel.ObtenerValoresCatalogoPorValor(bson.ObjectIdHex(Tipo))
		if Type == "" {
			EstatusPeticion = true
			Send.ETipoGrupo.IEstatus = true
			Send.ETipoGrupo.IMsj = "El Tipo seleccionado no coincide con los valores actuales de la base de datos, favor de recargar su página."
		}
	}
	Formulario := ctx.Request.Form
	Seleccionados := Formulario["Objeto"]
	fmt.Println("Seleccionados: ", Seleccionados)
	if len(Seleccionados) > 0 {
		Ids := GrupoPersonaModel.ConvierteAObjectIDS(Seleccionados)
		fmt.Println("Ids: ", Ids)
		if len(Seleccionados) != len(Ids) {
			EstatusPeticion = true
			Send.EMiembrosGrupo.IEstatus = true
			Send.EMiembrosGrupo.IMsj = "Hay insonsistencias en los objetos seleccionados, no son los mismos los que se asignan a los que existen en la base de Datos, favor de recargar su página."
		} else {
			SendMgo.Miembros = Ids

			Cabecera, Cuerpo, Encontrados = GrupoModel.GeneraTemplatesBusquedaObjetos(SendMgo.Miembros, Type, true)

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			paginasTotales = MoGeneral.Totalpaginas(Encontrados, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		}
	} else {
		EstatusPeticion = true
		Send.EMiembrosGrupo.IEstatus = true
		Send.EMiembrosGrupo.IMsj = "Debe seleccionar al menos un objeto a guardar en este Grupo."
	}

	Estatus := ctx.FormValue("Estatus")
	Send.Grupo.EEstatusGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(146, Estatus))
	if !bson.IsObjectIdHex(Tipo) {
		EstatusPeticion = true
		Send.EEstatusGrupo.IEstatus = true
		Send.EEstatusGrupo.IMsj = "El campo Estatus es obligatorio."
	} else {
		SendMgo.Estatus = bson.ObjectIdHex(Estatus)
		Status := CatalogoModel.ObtenerValoresCatalogoPorValor(bson.ObjectIdHex(Estatus))
		if Status == "" {
			EstatusPeticion = true
			Send.EEstatusGrupo.IEstatus = true
			Send.EEstatusGrupo.IMsj = "El Estatus seleccionado no coincide con los valores actuales de la base de datos, favor de recargar su página."
		}
	}

	ide := ctx.FormValue("ID")

	if bson.IsObjectIdHex(ide) {
		SendMgo.ID = bson.ObjectIdHex(ide)
		Send.Grupo.ID = bson.ObjectIdHex(ide)
		Grupo := GrupoModel.GetOne(bson.ObjectIdHex(ide))
		if !MoGeneral.EstaVacio(Grupo) {
			if Grupo.Nombre != Nombre {
				if SendMgo.ConsultaExistenciaByFieldMgo("Nombre", Nombre) {
					EstatusPeticion = true
					Send.ENombreGrupo.IEstatus = true
					Send.ENombreGrupo.IMsj = "El Nombre del catálogo ya existe, favor de intentar con otro."
				}
			}
			SendMgo.FechaEdicion = append(Grupo.FechaEdicion, time.Now())
		} else {
			Send.SEstado = false
			Send.SMsj = "Ha ocurrido un problema con la referencia del catálogo vuelva a seleccionar su catálogo a editar, disculpe la molestia."
			ctx.Render("GrupoEdita.html", Send)
			return
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "Ha ocurrido un problema con la referencia del catálogo vuelva a seleccionar su catálogo a editar, disculpe la molestia."
		ctx.Render("GrupoEdita.html", Send)
		return
	}

	if EstatusPeticion {
		Send.SEstado = false
		Send.SMsj = "La validación indica que el objeto capturado no puede aceptar la edición."
	} else {

		if SendMgo.ActualizaMgo([]string{"Nombre", "Descripcion", "Tipo", "Miembros", "Estatus", "FechaEdicion"}, []interface{}{SendMgo.Nombre, SendMgo.Descripcion, SendMgo.Tipo, SendMgo.Miembros, SendMgo.Estatus, SendMgo.FechaEdicion}) {
			if SendMgo.ActualizaElastic() {
				Send.SEstado = true
				Send.SMsj = "Se ha realizado una inserción exitosa."
				ctx.Redirect("/Grupos/detalle/"+SendMgo.ID.Hex(), 301)
			} else {
				Send.SEstado = false
				Send.SMsj = "Se Actualizó en Mongo, pero no en Elastic."
			}
		} else {
			Send.SEstado = false
			Send.SMsj = "Ocurrió un error al actualizar el Objeto, intente más tarde."
		}

	}

	ctx.Render("GrupoEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send GrupoModel.SGrupo
	var Cabecera, Cuerpo string
	var Encontrados int

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

	ID := ctx.Param("ID")
	if bson.IsObjectIdHex(ID) {
		Grupo := GrupoModel.GetOne(bson.ObjectIdHex(ID))
		if Grupo.ID.Hex() != "" {
			Send.Grupo.ID = Grupo.ID
			Send.Grupo.ENombreGrupo.Nombre = Grupo.Nombre
			Send.Grupo.EDescripcionGrupo.Descripcion = Grupo.Descripcion
			Send.Grupo.ETipoGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(179, Grupo.Tipo.Hex()))
			Type := CatalogoModel.ObtenerValoresCatalogoPorValor(Grupo.Tipo)

			Cabecera, Cuerpo, Encontrados = GrupoModel.GeneraTemplatesBusquedaObjetosDetalle(Grupo.Miembros, Type)

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			paginasTotales = MoGeneral.Totalpaginas(Encontrados, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)
			Send.Grupo.EEstatusGrupo.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(146, Grupo.Estatus.Hex()))

			Send.Grupo.EFechaHoraGrupo.FechaHora = Grupo.FechaHora
			// Send.Grupo.EFechaEdicionGrupo.Ihtml = template.HTML(Grupo.FechaEdicion.Format("2006-01-02 15:04:05"))
			Send.SEstado = true
		} else {
			Send.SEstado = false
			Send.SMsj = "Grupo no encontrado"
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "Error en la referencia al Grupo"
	}
	ctx.Render("GrupoDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send GrupoModel.SGrupo

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

	ctx.Render("GrupoDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send GrupoModel.SGrupo

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

		Cabecera, Cuerpo := GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrToMongo))
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
	var Send GrupoModel.SGrupo
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Grupo.ENombreGrupo.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := GrupoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = GrupoModel.GeneraTemplatesBusqueda(GrupoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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

//ConsultaElasticBase recupera la cedena de texto a buscar y regresa los datos con la paginación
func ConsultaElasticBase(ctx *iris.Context) {
	var Send GrupoModel.SGrupo
	var Cabecera, Cuerpo string
	var Encontrados int
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

	grupo := ctx.FormValue("GrupoBase")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	Tipo := ctx.FormValue("Tipo")
	if !bson.IsObjectIdHex(Tipo) {
		Send.SEstado = false
		Send.SMsj = "No se seleccionó un tipo adecuado, favor de seleccionarlo 1."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	Type := CatalogoModel.ObtenerValoresCatalogoPorValor(bson.ObjectIdHex(Tipo))
	if Type == "" {
		Send.SEstado = false
		Send.SMsj = "No se seleccionó un tipo adecuado, favor de seleccionarlo 2."
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	cadena := ctx.FormValue("Cadena")
	if cadena != "" {
		docs := GrupoModel.BuscaObjetosEnElastic(cadena, Type)
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
			Cabecera, Cuerpo, Encontrados = GrupoModel.GeneraTemplatesBusquedaObjetos(arrToMongo, Type, false)

			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			paginasTotales = MoGeneral.Totalpaginas(Encontrados, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)
			Send.SEstado = true
		} else {
			Send.SEstado = false
			Send.SMsj = "No se encontraron resultados para mostrar."
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//MuestraIndexPorGrupoB regresa template de busqueda y paginacion de acuerdo a la agrupacion solicitada
func MuestraIndexPorGrupoB(ctx *iris.Context) {

	var Send GrupoModel.SGrupo

	// filtro := ctx.FormValue("Filtro")
	// if filtro != "" {
	// 	fil, _ := strconv.Atoi(filtro)
	// 	tipoConsultaBase = fil
	// }

	// grupo := ctx.FormValue("GrupoBase")
	// if grupo != "" {
	// 	gru, _ := strconv.Atoi(grupo)
	// 	limitePorPaginaBase = gru
	// }

	// cadenaBusqueda = ctx.FormValue("Cadena")
	// busqueda := ctx.FormValue("Busqueda")

	// index := ctx.GetCookie("IDUsuario")
	// if index == "" {
	// 	Send.SEstado = false
	// 	Send.SMsj = "No se encuentran datos de sesion."
	// } else {
	// 	tipo := "productoservicio"

	// 	if cadenaBusqueda != "" {

	// 		var arrProductoElastic2 *[]ListaModel.ProductosBase
	// 		if busqueda == "" {
	// 			arrProductoElastic2 = ListaModel.BuscarEnElasticDefault(index, tipo, cadenaBusqueda, tipoConsultaBase)
	// 		} else {
	// 			arrProductoElastic2 = ListaModel.BuscarEnElasticAvanzada(index, tipo, cadenaBusqueda, tipoConsultaBase)
	// 		}
	// 		numeroRegistros = len(*arrProductoElastic2)

	// 		if numeroRegistros == 0 {
	// 			Send.SEstado = false
	// 			Send.SMsj = "No se encontró ningún registro para mostrar."
	// 			jData, _ := json.Marshal(Send)
	// 			ctx.Header().Set("Content-Type", "application/json")
	// 			ctx.Write(jData)
	// 			return
	// 		}
	// 		arrProductoElastic = *arrProductoElastic2

	// 		arrProductoVista = []ListaModel.ProductosBase{}
	// 		if numeroRegistros <= limitePorPaginaBase {
	// 			for _, v := range arrProductoElastic[0:numeroRegistros] {
	// 				arrProductoVista = append(arrProductoVista, v)
	// 			}
	// 		} else if numeroRegistros >= limitePorPaginaBase {
	// 			for _, v := range arrProductoElastic[0:limitePorPaginaBase] {
	// 				arrProductoVista = append(arrProductoVista, v)
	// 			}
	// 		}

	// 		Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaBase(arrProductoVista, tipoConsultaBase)
	// 		Send.SCabecera = template.HTML(Cabecera)
	// 		Send.SBody = template.HTML(Cuerpo)

	// 		paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPaginaBase)
	// 		Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	// 		Send.SPaginacion = template.HTML(Paginacion)

	// 		Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaBase))

	// 		Send.SEstado = true

	// 	} else {
	// 		if len(arrProductoElastic) > 0 {

	// 			arrProductoVista = []ListaModel.ProductosBase{}
	// 			if numeroRegistros <= limitePorPaginaBase {
	// 				for _, v := range arrProductoElastic[0:numeroRegistros] {
	// 					arrProductoVista = append(arrProductoVista, v)
	// 				}
	// 			} else if numeroRegistros >= limitePorPaginaBase {
	// 				for _, v := range arrProductoElastic[0:limitePorPaginaBase] {
	// 					arrProductoVista = append(arrProductoVista, v)
	// 				}
	// 			}

	// 			if tipoConsultaBase == 2 {
	// 				arrProductoVista = *ListaModel.CheckConCodigo(&arrProductoVista)
	// 			} else if tipoConsultaBase == 1 {
	// 				arrProductoVista = *ListaModel.CheckSinCodigo(&arrProductoVista)
	// 			}

	// 			Cabecera, Cuerpo := ListaModel.GeneraTemplatesBusquedaBase(arrProductoVista, tipoConsultaBase)
	// 			Send.SCabecera = template.HTML(Cabecera)
	// 			Send.SBody = template.HTML(Cuerpo)

	// 			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPaginaBase)
	// 			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	// 			Send.SPaginacion = template.HTML(Paginacion)

	// 			Send.SGrupo = template.HTML(MoGeneral.CargaComboMostrarEnIndex(limitePorPaginaBase))

	// 		}
	// 	}
	// }

	Send.SEstado = true
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}
