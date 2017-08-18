package SesionesController

import (
	"encoding/json"
	"html/template"
	"strconv"

	"../../Modelos/SesionModel"
	"../../Modelos/UsuarioModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/IndexUsuarios"

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

//limitePorPagina limite de registros a mostrar en la pagina
var limitePorPagina = 10

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//IndexGet Index de Sesiones
func IndexGet(ctx *iris.Context) {

	var Send SesionModel.SSesiones
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
	numeroRegistros = UsuarioModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Usuarios := UsuarioModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Usuarios {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusquedaSesiones(Usuarios[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusquedaSesiones(Usuarios[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true
	ctx.Render("SesionesIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Usuario
func IndexPost(ctx *iris.Context) {

	var Send SesionModel.SSesiones
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
	//Send.Usuario.EVARIABLEUsuario.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := UsuarioModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := IndexUsuarios.GeneraTemplatesBusquedaSesiones(UsuarioModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusquedaSesiones(UsuarioModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusquedaSesiones(UsuarioModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("SesionesIndex.html", Send)

}

//DetalleGet renderea al alta de Usuario
func DetalleGet(ctx *iris.Context) {
	var Send SesionModel.SSesiones
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

	usuarioses := ctx.Param("ID")
	if usuarioses != "" {
		Send.Sesion.Nombre = usuarioses
		sesiones, err := Session.GetInfoSession(usuarioses)
		if err != nil {
			Send.Sesion.IEstatus = true
			Send.Sesion.IMsj = err.Error()
		}
		Send.Sesion.Ihtml = template.HTML(CreaTablaSesionesDetalle(sesiones))

	} else {
		Send.SEstado = false
		Send.SMsj = "Debe Introducir un parametro valido"
	}

	ctx.Render("SesionesDetalle.html", Send)
}

//DetallePost renderea al alta de Usuario
func DetallePost(ctx *iris.Context) {

	var Send SesionModel.SSesiones
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

	ctx.Render("SesionesDetalle.html", Send)
}

//EditaGet renderea al alta de Usuario
func EditaGet(ctx *iris.Context) {
	var Send SesionModel.SSesiones
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

	usuarioses := ctx.Param("ID")
	if usuarioses != "" {
		Send.Sesion.Nombre = usuarioses
		sesiones, err := Session.GetInfoSession(usuarioses)
		if err != nil {
			Send.Sesion.IEstatus = true
			Send.Sesion.IMsj = err.Error()
		}
		if len(sesiones) > 0 {
			html := `<div class="text-right"> <button type="button" class="btn btn-danger closeAll" onclick="closeAll()"><span class="glyphicon glyphicon-remove">Cerrar Todas las Sesiones</span></button>  </div> </br>`
			html += CreaTablaSesionesEdita(sesiones)
			Send.Sesion.Ihtml = template.HTML(html)
		} else {
			Send.Sesion.Ihtml = template.HTML(CreaTablaSesionesEdita(sesiones))
		}

	} else {
		Send.SEstado = false
		Send.SMsj = "Debe Introducir un parametro valido"
	}

	ctx.Render("SesionesEdita.html", Send)
}

//SesionesTotal carga la vista con todas las sesiones iniciadas en el sistema
func SesionesTotal(ctx *iris.Context) {

	var Send SesionModel.SSesiones
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

	SesionesActivas, err := Session.RegresaSesionesActivas()
	if err != nil {
		Send.Sesion.IEstatus = true
		Send.Sesion.IMsj = "Error al consultar las sesiones(Redis)"
	}

	Send.Sesion.Ihtml = template.HTML(CreaTablaSesionesEdita(SesionesActivas))

	ctx.Render("SesionesTotal.html", Send)
}

// EliminaByID elimina una sesion por un ID especificado peticion AJAX
func EliminaByID(ctx *iris.Context) {
	var Send SesionModel.SAjaxSesiones
	SesionActiva, funcion, _, errSes := Session.GetDataSessionAJAX(ctx) //Retorna los datos de la session

	if errSes != nil {
		Send.SEstado = false
		Send.SMsj = errSes.Error()
	}

	if !SesionActiva {
		Send.SEstado = false
		Send.SMsj = "Necesitas iniciar Sesion"
		Send.SFuncion = funcion
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	//Verificacion de Permisos a URL
	permiso := true
	if !permiso {
		Send.SEstado = false
		Send.SMsj = "Sin Autorizacion para realizar esta accion"
	} else {
		usuarioses := ctx.FormValue("ID")
		err := Session.EliminaSesionByID(usuarioses)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Existe un problema, intentelo de nuevo"
		} else {
			Send.SEstado = true
			Send.SMsj = "Sesion Finalizada"
		}

	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)

}

// EliminaByName elimina todas las sesiones de un usuarios especificado - peticion AJAX
func EliminaByName(ctx *iris.Context) {
	var Send SesionModel.SAjaxSesiones
	SesionActiva, funcion, _, errSes := Session.GetDataSessionAJAX(ctx) //Retorna los datos de la session

	if errSes != nil {
		Send.SEstado = false
		Send.SMsj = errSes.Error()
	}

	if !SesionActiva {
		Send.SEstado = false
		Send.SMsj = "Necesitas iniciar Sesion"
		Send.SFuncion = funcion
		jData, _ := json.Marshal(Send)
		ctx.Header().Set("Content-Type", "application/json")
		ctx.Write(jData)
		return
	}

	//Verificacion de Permisos a URL
	permiso := true
	if !permiso {
		Send.SEstado = false
		Send.SMsj = "Sin Autorizacion para realizar esta accion"
	} else {
		usuarioses := ctx.FormValue("ID")
		err := Session.EliminarSesionesUsr(usuarioses)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Existe un problema, intentelo de nuevo"
		} else {
			Send.SEstado = true
			Send.SMsj = "Sesiones Finalizadas"
			Send.SFuncion = "location.reload();"
		}

	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)

}

// CreaTablaSesionesDetalle crea la tabla de detalles de sesiones
func CreaTablaSesionesDetalle(sesiones []map[string]string) string {
	html := ""
	html += `<table class="table table-sm" style="background: black;color: lawngreen;">
				<thead class="thead-inverse">
				<tr>
				<th>ID</th>
				<th>UserAgent</th>
				<th>Ubicacion Actual</th>
				<th>IP</th>
				<th>ServerHost</th>
				<th>Fecha</th>
				<th>VirtualHostName</th>
				</tr>
			</thead>
			<tbody>`

	for _, val := range sesiones {
		html += `<tr>`
		html += `<td>` + val["IDSession"] + `</td>`
		html += `<td>` + val["UserAgent"] + `</td>`
		html += `<td>` + val["Ubicacion"] + `</td>`
		html += `<td>` + val["IP"] + `</td>`
		html += `<td>` + val["ServerHost"] + `</td>`
		html += `<td>` + val["Time"] + `</td>`
		html += `<td>` + val["VirtualHostName"] + `</td>`
		html += `</tr>`
	}
	html += `</tbody>
			</table>`
	return html
}

// CreaTablaSesionesEdita crea la tabla de detalles de sesiones
func CreaTablaSesionesEdita(sesiones []map[string]string) string {
	html := ""

	html += `<table class="table table-sm" style="background: black;color: lawngreen;">
				<thead class="thead-inverse">
				<tr>
				<th>ID</th>
				<th>UserAgent</th>
				<th>Ubicacion Actual</th>
				<th>IP</th>
				<th>ServerHost</th>
				<th>Fecha</th>
				<th>VirtualHostName</th>
				<th>Eliminar</th>
				</tr>
			</thead>
			<tbody>`

	for _, val := range sesiones {
		html += `<tr>`
		html += `<td>` + val["IDSession"] + `</td>`
		html += `<td>` + val["UserAgent"] + `</td>`
		html += `<td>` + val["Ubicacion"] + `</td>`
		html += `<td>` + val["IP"] + `</td>`
		html += `<td>` + val["ServerHost"] + `</td>`
		html += `<td>` + val["Time"] + `</td>`
		html += `<td>` + val["VirtualHostName"] + `</td>`
		html += `<td><button type="button" class="btn btn-danger deleteButton" ID="` + val["IDSession"] + `"><span class="glyphicon glyphicon-remove"></span></button> </td>`
		html += `</tr>`
	}
	html += `</tbody>
			</table>`
	return html
}
