package BugControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modulos/Session"

	"../../Modelos/BugModel"
	"../../Modelos/CatalogoModel"
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

//CatalogoDeEstatusDeBugs para leer los estatus de los bugs
var CatalogoDeEstatusDeBugs = 181

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Bug
func IndexGet(ctx *iris.Context) {

	var Send BugModel.SBug

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
	numeroRegistros = BugModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Bugs := BugModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Bugs {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(Bugs[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(Bugs[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("BugIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Bug
func IndexPost(ctx *iris.Context) {

	var Send BugModel.SBug

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
	//Send.Bug.EVARIABLEBug.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := BugModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("BugIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Bug
func AltaGet(ctx *iris.Context) {
	var Send BugModel.SBug
	var BugMgo BugModel.BugMgo

	SesionActiva, funcion, UsuarioLogeado, errSes := Session.GetDataSessionAJAX(ctx)

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

	Send.SEstado = true
	Titulo := ctx.FormValue("Titulo")
	BugMgo.Titulo = Titulo

	if Titulo == "" {
		Send.SEstado = false
		Send.SMsj = "Debes agregar un TITULO para identificar tu problema"
	}

	Descripcion := ctx.FormValue("Contenido")
	BugMgo.Descripcion = Descripcion
	if Titulo == "" {
		Send.SEstado = false
		Send.SMsj = "Digas mam.. si agregas un problema, tienes que describir que ocurrio"
	}

	BugMgo.Usuario = UsuarioLogeado
	BugMgo.EsAjax = true

	Ruta := ctx.FormValue("Ruta")
	BugMgo.Ruta = Ruta
	if Ruta == "" {
		Send.SEstado = false
		Send.SMsj = "No se recibio la RUTA de manera correcta"
	}

	BugMgo.FechaHora = time.Now()

	// BugMgo.Estatus  //Falta crear catalogo

	if Send.SEstado {
		if BugMgo.InsertaMgo() {
			Send.SMsj = "Tu problema ha sido guardado, y sera atendido lo antes posible."
			jData, _ := json.Marshal(Send)
			ctx.Header().Set("Content-Type", "application/json")
			ctx.Write(jData)
			return
		}
		Send.SMsj = "Existio un problema al guardar el problema en Mongo"
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//AltaPost regresa la petición post que se hizo desde el alta de Bug
func AltaPost(ctx *iris.Context) {

	var Send BugModel.SBug
	var BugMgo BugModel.BugMgo

	SesionActiva, funcion, UsuarioLogeado, errSes := Session.GetDataSessionAJAX(ctx)

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

	Send.SEstado = true
	Titulo := ctx.FormValue("Titulo")
	BugMgo.Titulo = Titulo

	if Titulo == "" {
		Send.SEstado = false
		Send.SMsj = "Debes agregar un TITULO para identificar tu problema"
	}

	Descripcion := ctx.FormValue("Contenido")
	BugMgo.Descripcion = Descripcion
	if Titulo == "" {
		Send.SEstado = false
		Send.SMsj = "Digas mam.. si agregas un problema, tienes que describir que ocurrio"
	}

	BugMgo.Usuario = UsuarioLogeado
	BugMgo.EsAjax = true
	EstatusBug := CatalogoModel.GetValorDefaultCatalogoBugs(CatalogoDeEstatusDeBugs)
	if EstatusBug != "" {
		BugMgo.Estatus = bson.ObjectIdHex(EstatusBug)
	} else {
		Send.SEstado = false
		Send.SMsj = "Error al establecer el estatis PENDIENTE, porfavor verifique sus catalogos"
		fmt.Println("No se esncontro el estatus por defecto para el altaBug en mongo")
	}

	Ruta := ctx.FormValue("Ruta")
	BugMgo.Ruta = Ruta
	if Ruta == "" {
		Send.SEstado = false
		Send.SMsj = "No se recibio la RUTA de manera correcta"
	}

	BugMgo.FechaHora = time.Now()

	// BugMgo.Estatus  //Falta crear catalogo

	if Send.SEstado {
		if BugMgo.InsertaMgo() {
			if BugMgo.InsertaElastic() {
				Send.SMsj = "Tu problema ha sido guardado, y sera atendido lo antes posible."
				jData, _ := json.Marshal(Send)
				ctx.Header().Set("Content-Type", "application/json")
				ctx.Write(jData)
				return
			}
			Send.SEstado = false
			Send.SMsj = "Existio un problema al guardar el problema en Elastic"
		}
		Send.SEstado = false
		Send.SMsj = "Existio un problema al guardar el problema en Mongo"
	}

	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Bug
func EditaGet(ctx *iris.Context) {

	var Send BugModel.SBug

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

	ctx.Render("BugEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Bug
func EditaPost(ctx *iris.Context) {

	var Send BugModel.SBug

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

	ctx.Render("BugEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send BugModel.SBug

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

	ctx.Render("BugDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send BugModel.SBug

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

	ctx.Render("BugDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send BugModel.SBug

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

		Cabecera, Cuerpo := BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrToMongo))
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
	var Send BugModel.SBug
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Bug.ENombreBug.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := BugModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = BugModel.GeneraTemplatesBusqueda(BugModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
