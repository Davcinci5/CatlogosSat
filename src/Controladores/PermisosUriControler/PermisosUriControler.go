package PermisosUriControler

import (
	"encoding/json"
	"html/template"
	"strconv"

	"../../Modelos/GrupoPersonaModel"
	"../../Modelos/PermisosUriModel"
	"../../Modelos/UsuarioModel"

	"../../Modulos/CargaCombos"
	"../../Modulos/General"
	"../../Modulos/Redis"
	"../../Modulos/Session"

	"strings"

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

//PermisosUris guarda todas las estructuras de usuarios y grupos par amostrar en el index.
var PermisosUris []PermisosUriModel.PermisoAux

var arrIDMgo []bson.ObjectId
var arrID []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrIDElasticG []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de PermisosUri
func IndexGet(ctx *iris.Context) {

	var Send PermisosUriModel.SPermisosUri

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
	numUsuarios := UsuarioModel.CountAll()
	numGrupos := GrupoPersonaModel.CountAll()

	numeroRegistros = numUsuarios + numGrupos

	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)

	Usuarios := UsuarioModel.GetAll()
	Grupos := GrupoPersonaModel.GetAll()

	Permisos := PermisosUriModel.JoinGroupUser(&Usuarios, &Grupos)
	PermisosUris = *Permisos

	arrIDMgo = []bson.ObjectId{}
	for _, v := range PermisosUris {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrID = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("PermisosUriIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de PermisosUri
func IndexPost(ctx *iris.Context) {

	var Send PermisosUriModel.SPermisosUri
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
	//Send.PermisosUri.EVARIABLEPermisosUri.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := PermisosUriModel.BuscarEnElastic(cadenaBusqueda)
		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}
			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}
		}

		docs2 := PermisosUriModel.BuscarEnElasticG(cadenaBusqueda)
		if docs2.Hits.TotalHits > 0 {
			arrIDElasticG = []bson.ObjectId{}
			for _, item := range docs2.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElasticG = append(arrIDElasticG, IDElastic)
			}
		}

		arrID = arrIDElastic
		for _, v := range arrIDElasticG {
			arrID = append(arrID, v)
		}

		numeroRegistros = len(arrID)

		Usuarios := UsuarioModel.GetEspecifics(arrIDElastic)
		Grupos := GrupoPersonaModel.GetEspecifics(arrIDElasticG)

		Permisos := PermisosUriModel.JoinGroupUser(&Usuarios, &Grupos)
		PermisosUris = *Permisos

		if numeroRegistros > 0 {
			Cabecera, Cuerpo := PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris)
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:numeroRegistros])
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:limitePorPagina])
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
	ctx.Render("PermisosUriIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de PermisosUri
func AltaGet(ctx *iris.Context) {

	var Send PermisosUriModel.SPermisosUri
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
	ID := ctx.Param("ID")
	if bson.IsObjectIdHex(ID) {
		Send.PermisosUri.ID = bson.ObjectIdHex(ID)
		datoUsr := UsuarioModel.GetOne(bson.ObjectIdHex(ID))
		if MoGeneral.EstaVacio(datoUsr) {
			datosGpo := GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToUpper(datosGpo.Nombre)
		} else {
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToLower(datoUsr.Usuario)
		}

		ides, err := Redis.ObtenerMiembrosdeGrupo(ID)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Error al Consultar ACL (Redis)"
		}
		UrisPermitias := PermisosUriModel.GeneraEtiquetasEspecificPermisosURIS(ides)
		UrisNoPermitidas := PermisosUriModel.GeneraEtiquetasEspecificNoPermisosURIS(ides)
		Send.EPermisoAceptadoPermisosUri.Ihtml = template.HTML(UrisPermitias)
		Send.EPermisoNegadoPermisosUri.Ihtml = template.HTML(UrisNoPermitidas)

	} else {
		Send.SEstado = false
		Send.SMsj = "El parametro recibido no es correcto"
	}

	ctx.Render("PermisosUriAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de PermisosUri
func AltaPost(ctx *iris.Context) {

	var Send PermisosUriModel.SPermisosUri
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
	ID := ctx.FormValue("ID")
	ides := ctx.Request.Form["PermisosAceptados"]
	if bson.IsObjectIdHex(ID) {
		Send.PermisosUri.ID = bson.ObjectIdHex(ID)
		datoUsr := UsuarioModel.GetOne(bson.ObjectIdHex(ID))
		if MoGeneral.EstaVacio(datoUsr) {
			datosGpo := GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToUpper(datosGpo.Nombre)
		} else {
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToLower(datoUsr.Usuario)
		}

		err1 := Redis.EliminarConjunto(ID)
		err2 := Redis.InsertaRedis(ID, ides)
		if err1 != nil || err2 != nil {

			UrisPermitias := PermisosUriModel.GeneraEtiquetasEspecificPermisosURIS(ides)
			UrisNoPermitidas := PermisosUriModel.GeneraEtiquetasEspecificNoPermisosURIS(ides)
			Send.EPermisoAceptadoPermisosUri.Ihtml = template.HTML(UrisPermitias)
			Send.EPermisoNegadoPermisosUri.Ihtml = template.HTML(UrisNoPermitidas)
			Send.SEstado = false
			Send.SMsj = "Error al realizar accion en Redis"

		} else {
			ctx.Redirect("/PermisosUris/detalle/"+ID, 301)
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "El parametro recibido no es correcto"
	}
	ctx.Render("PermisosUriAlta.html", Send)

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de PermisosUri
func EditaGet(ctx *iris.Context) {

	var Send PermisosUriModel.SPermisosUri
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
	ID := ctx.Param("ID")
	if bson.IsObjectIdHex(ID) {
		datoUsr := UsuarioModel.GetOne(bson.ObjectIdHex(ID))
		if MoGeneral.EstaVacio(datoUsr) {
			datosGpo := GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToUpper(datosGpo.Nombre)
		} else {
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToLower(datoUsr.Usuario)
		}
		Send.PermisosUri.ID = bson.ObjectIdHex(ID)

		ides, err := Redis.ObtenerMiembrosdeGrupo(ID)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Error al Consultar ACL (Redis)"
		}
		UrisPermitias := PermisosUriModel.GeneraEtiquetasEspecificPermisosURIS(ides)
		UrisNoPermitidas := PermisosUriModel.GeneraEtiquetasEspecificNoPermisosURIS(ides)
		Send.EPermisoAceptadoPermisosUri.Ihtml = template.HTML(UrisPermitias)
		Send.EPermisoNegadoPermisosUri.Ihtml = template.HTML(UrisNoPermitidas)

	} else {
		Send.SEstado = false
		Send.SMsj = "El parametro recibido no es correcto"
	}

	ctx.Render("PermisosUriEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de PermisosUri
func EditaPost(ctx *iris.Context) {

	var Send PermisosUriModel.SPermisosUri
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
	ID := ctx.FormValue("ID")
	ides := ctx.Request.Form["PermisosAceptados"]
	if bson.IsObjectIdHex(ID) {
		datoUsr := UsuarioModel.GetOne(bson.ObjectIdHex(ID))
		if MoGeneral.EstaVacio(datoUsr) {
			datosGpo := GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToUpper(datosGpo.Nombre)
		} else {
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToLower(datoUsr.Usuario)
		}
		Send.PermisosUri.ID = bson.ObjectIdHex(ID)

		err1 := Redis.EliminarConjunto(ID)
		err2 := Redis.InsertaRedis(ID, ides)
		if err1 != nil || err2 != nil {

			UrisPermitias := PermisosUriModel.GeneraEtiquetasEspecificPermisosURIS(ides)
			UrisNoPermitidas := PermisosUriModel.GeneraEtiquetasEspecificNoPermisosURIS(ides)
			Send.EPermisoAceptadoPermisosUri.Ihtml = template.HTML(UrisPermitias)
			Send.EPermisoNegadoPermisosUri.Ihtml = template.HTML(UrisNoPermitidas)
			Send.SEstado = false
			Send.SMsj = "Error al realizar accion en Redis"

		} else {
			ctx.Redirect("/PermisosUris/detalle/"+ID, 301)
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "El parametro recibido no es correcto"
	}

	ctx.Render("PermisosUriEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send PermisosUriModel.SPermisosUri

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
	ID := ctx.Param("ID")
	if bson.IsObjectIdHex(ID) {
		datoUsr := UsuarioModel.GetOne(bson.ObjectIdHex(ID))
		if MoGeneral.EstaVacio(datoUsr) {
			datosGpo := GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToUpper(datosGpo.Nombre)
		} else {
			Send.PermisosUri.EGrupoPermisosUri.Grupo = strings.ToLower(datoUsr.Usuario)
		}
		Send.PermisosUri.ID = bson.ObjectIdHex(ID)

		ides, err := Redis.ObtenerMiembrosdeGrupo(ID)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Error al Consultar ACL (Redis)"
		}

		UrisPermitias := PermisosUriModel.GeneraEtiquetasEspecificPermisosURIS(ides)
		Send.EPermisoAceptadoPermisosUri.Ihtml = template.HTML(UrisPermitias)
	} else {
		Send.SEstado = false
		Send.SMsj = "El parametro recibido no es correcto"
	}
	ctx.Render("PermisosUriDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send PermisosUriModel.SPermisosUri
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

	ctx.Render("PermisosUriDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send PermisosUriModel.SPermisosUri
	var Cabecera, Cuerpo string
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
				Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[skip:limite])
			} else {
				Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[skip : skip+final])
			}

		} else {
			Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[skip:limite])
		}

		//	Cabecera, Cuerpo := PermisosUriModel.GeneraTemplatesBusqueda(PermisosUriModel.GetEspecifics(arrToMongo))
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
	var Send PermisosUriModel.SPermisosUri
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.PermisosUri.ENombrePermisosUri.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := PermisosUriModel.BuscarEnElastic(cadenaBusqueda)
		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}
			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}
		}

		docs2 := PermisosUriModel.BuscarEnElasticG(cadenaBusqueda)
		if docs2.Hits.TotalHits > 0 {
			arrIDElasticG = []bson.ObjectId{}
			for _, item := range docs2.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElasticG = append(arrIDElasticG, IDElastic)
			}
		}

		arrID = arrIDElastic
		for _, v := range arrIDElasticG {
			arrID = append(arrID, v)
		}

		numeroRegistros = len(arrID)

		Usuarios := UsuarioModel.GetEspecifics(arrIDElastic)
		Grupos := GrupoPersonaModel.GetEspecifics(arrIDElasticG)

		Permisos := PermisosUriModel.JoinGroupUser(&Usuarios, &Grupos)
		PermisosUris = *Permisos

		if numeroRegistros > 0 {
			Cabecera, Cuerpo := PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris)
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:numeroRegistros])
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:limitePorPagina])
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

		if numeroRegistros <= limitePorPagina {
			Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:numeroRegistros])
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = PermisosUriModel.GeneraTemplatesBusqueda2(PermisosUris[0:limitePorPagina])
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
