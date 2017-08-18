package AlmacenControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modelos/OperacionModel"
	"../../Modelos/ProductoModel"
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
var limitePorPagina = 5

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Almacen
func IndexGet(ctx *iris.Context) {

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

	var Cabecera, Cuerpo string
	numeroRegistros = AlmacenModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Almacens := AlmacenModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Almacens {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(Almacens[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(Almacens[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("AlmacenIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Almacen
func IndexPost(ctx *iris.Context) {

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

	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Almacen.EVARIABLEAlmacen.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := AlmacenModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("AlmacenIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Almacen
func AltaGet(ctx *iris.Context) {

	var SAlmacenes AlmacenModel.SAlmacen
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SAlmacenes.SSesion.Name = NameUsrLoged
	SAlmacenes.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SAlmacenes.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SAlmacenes.SEstado = false
		SAlmacenes.SMsj = errSes.Error()
		ctx.Render("ZError.html", SAlmacenes)
		return
	}

	SAlmacenes.SEstado = true
	SAlmacenes.ETipoAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(132, ""))
	// fmt.Println(CargaCombos.CargaComboCatalogo(133, ""))
	SAlmacenes.EClasificacionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(133, ""))
	SAlmacenes.EConexionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboConexiones(""))
	SAlmacenes.EEstatusAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(134, "58e56961e75770120c60befc"))

	// var almacen AlmacenModel.AlmacenMgo
	// almacen.ID = bson.NewObjectId()
	// fmt.Println("almacen.ID", almacen.ID)
	// almacen.GetAncestorsIDs()
	// almacen.GetDescendantsIDs()

	// ObjetoPredecesorSucesor := almacen.AvailableOthersIDs()

	// // SAlmacenes.EPredecesorAlmacen.Ihtml = template.HTML(GetTemplateAvailableOthersIds(ObjetoPredecesorSucesor, almacen.ID.Hex()))
	// SAlmacenes.ESucesorAlmacen.Ihtml = template.HTML(GetTemplateAvailableOthersIds(ObjetoPredecesorSucesor, almacen.ID.Hex()))

	ctx.Render("AlmacenAlta.html", SAlmacenes)

}

//AltaPost regresa la petición post que se hizo desde el alta de Almacen
func AltaPost(ctx *iris.Context) {
	var SAlmacen AlmacenModel.SAlmacen

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SAlmacen.SSesion.Name = NameUsrLoged
	SAlmacen.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SAlmacen.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SAlmacen.SEstado = false
		SAlmacen.SMsj = errSes.Error()
		ctx.Render("ZError.html", SAlmacen)
		return
	}

	//######### LEE TU OBJETO DEL FORMULARIO #########
	tipoFisico := bson.ObjectIdHex("58e56906e75770120c60bef2")
	var Almacen AlmacenModel.AlmacenMgo
	Almacen.Nombre = ctx.FormValue("Nombre")

	// predecesor := MoGeneral.EliminarEspaciosInicioFinal(ctx.FormValue("Predecesor"))

	// if !MoGeneral.EstaVacio(predecesor) {
	// 	Almacen.Predecesor = bson.ObjectIdHex(predecesor)
	// }

	Almacen.Tipo = bson.ObjectIdHex(ctx.FormValue("Tipo"))
	Almacen.Clasificacion = bson.ObjectIdHex(ctx.FormValue("Clasificacion"))
	Almacen.Estatus = bson.ObjectIdHex("58e56961e75770120c60befc")
	// descendientes := ctx.Request.Form["Sucesor"]

	if Almacen.Tipo == tipoFisico {
		Almacen.Conexion = bson.ObjectIdHex(ctx.FormValue("Conexion"))
	}
	EstatusPeticion := false //True indica que hay un error

	Almacen.ID = bson.NewObjectId()

	// for _, item := range descendientes {
	// 	aux := AlmacenModel.GetOne(bson.ObjectIdHex(item))
	// 	if Almacen.Clasificacion == aux.Clasificacion && Almacen.ID != aux.ID && Almacen.ID != Almacen.Predecesor {
	// 		Almacen.Sucesor = append(Almacen.Sucesor, bson.ObjectIdHex(item))
	// 	}
	// 	fmt.Println("Nuevo Objeto Almacen - Sucesores:", item)
	// }

	if EstatusPeticion {
		SAlmacen.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SAlmacen.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("AlmacenAlta.html", SAlmacen)
	} else {
		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		Almacen.FechaHora = time.Now()
		if Almacen.InsertaMgo() {
			if Almacen.InsertaElastic() {
				if Almacen.Tipo == tipoFisico {
					AlmacenModel.CreaTablasPostgres(Almacen.ID.Hex(), Almacen.Conexion)
				}
				SAlmacen.SEstado = true
				SAlmacen.SMsj = "Se ha realizado una inserción exitosa"
				SAlmacen.Almacen.ID = Almacen.ID
				SAlmacen.Almacen.ENombreAlmacen.Nombre = Almacen.Nombre
				SAlmacen.Almacen.ETipoAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(132, Almacen.Tipo.Hex()))
				SAlmacen.Almacen.EClasificacionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(133, Almacen.Clasificacion.Hex()))
				SAlmacen.Almacen.EEstatusAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(134, Almacen.Estatus.Hex()))
				SAlmacen.Almacen.EConexionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboConexiones(Almacen.Conexion.Hex()))
				ctx.Render("AlmacenDetalle.html", SAlmacen)
			} else {
				if Almacen.EliminaByIDMgo() {
					SAlmacen.SEstado = false
					SAlmacen.SMsj = "Falló la inserción del Almacen, se ha eliminado los registros"
					SAlmacen.Almacen.ID = Almacen.ID
					SAlmacen.Almacen.ENombreAlmacen.Nombre = Almacen.Nombre
					SAlmacen.Almacen.ETipoAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(132, Almacen.Tipo.Hex()))
					SAlmacen.Almacen.EClasificacionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(133, Almacen.Clasificacion.Hex()))
					SAlmacen.Almacen.EEstatusAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(134, Almacen.Estatus.Hex()))
					SAlmacen.Almacen.EConexionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboConexiones(Almacen.Conexion.Hex()))
					ctx.Render("AlmacenAlta.html", SAlmacen)
				} else {
					SAlmacen.SEstado = false
					SAlmacen.SMsj = "No se inserto en elastic, no se pudo eliminar el registro en MongoDb"
					ctx.Render("AlmacenAlta.html", SAlmacen)
				}
			}
		} else {
			SAlmacen.SEstado = false
			SAlmacen.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("AlmacenAlta.html", SAlmacen)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Almacen
func EditaGet(ctx *iris.Context) {

	var SAlmacenes AlmacenModel.SAlmacen

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SAlmacenes.SSesion.Name = NameUsrLoged
	SAlmacenes.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SAlmacenes.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SAlmacenes.SEstado = false
		SAlmacenes.SMsj = errSes.Error()
		ctx.Render("ZError.html", SAlmacenes)
		return
	}

	id := ctx.Param("ID")
	if bson.IsObjectIdHex(id) {
		Almacen := AlmacenModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Almacen) {

			SAlmacenes.SEstado = true
			SAlmacenes.ID = Almacen.ID
			SAlmacenes.ENombreAlmacen.Nombre = Almacen.Nombre
			SAlmacenes.ETipoAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(132, Almacen.Tipo.Hex()))
			SAlmacenes.EClasificacionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(133, Almacen.Clasificacion.Hex()))
			SAlmacenes.EConexionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboConexiones(Almacen.Conexion.Hex()))
			SAlmacenes.EEstatusAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(134, Almacen.Estatus.Hex()))

			ctx.Render("AlmacenEdita.html", SAlmacenes)

		} else {
			SAlmacenes.SEstado = false
			SAlmacenes.SMsj = "El almacen no se ha encontrado, intente de nuevo."
			ctx.Render("AlmacenIndex.html", SAlmacenes)
		}

	} else {
		fmt.Println("No es un id")
		SAlmacenes.SEstado = false
		SAlmacenes.SMsj = "No se ha recibido un parámetro adecuado para poder editar El Almacen, intente de nuevo."
		ctx.Render("AlmacenIndex.html", SAlmacenes)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Almacen
func EditaPost(ctx *iris.Context) {

	var SAlmacenes AlmacenModel.SAlmacen

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SAlmacenes.SSesion.Name = NameUsrLoged
	SAlmacenes.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SAlmacenes.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SAlmacenes.SEstado = false
		SAlmacenes.SMsj = errSes.Error()
		ctx.Render("ZError.html", SAlmacenes)
		return
	}

	id := ctx.Param("ID")
	if bson.IsObjectIdHex(id) {
		var AlmacenNuevo AlmacenModel.AlmacenMgo
		AlmacenExistente := AlmacenModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(AlmacenExistente) {
			AlmacenNuevo.ID = AlmacenExistente.ID
			AlmacenNuevo.Nombre = ctx.FormValue("Nombre")
			AlmacenNuevo.Tipo = AlmacenExistente.Tipo
			clasificacionAlmacen := bson.ObjectIdHex(ctx.FormValue("Clasificacion"))
			AlmacenNuevo.Clasificacion = clasificacionAlmacen
			AlmacenNuevo.Estatus = bson.ObjectIdHex(ctx.FormValue("Estatus"))
			AlmacenNuevo.Conexion = AlmacenExistente.Conexion

			SAlmacenes.SEstado = true
			SAlmacenes.Almacen.ID = AlmacenExistente.ID
			SAlmacenes.Almacen.ENombreAlmacen.Nombre = AlmacenNuevo.Nombre
			SAlmacenes.Almacen.ETipoAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(132, AlmacenNuevo.Tipo.Hex()))
			SAlmacenes.Almacen.EClasificacionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(133, AlmacenNuevo.Clasificacion.Hex()))
			SAlmacenes.Almacen.EEstatusAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(134, AlmacenNuevo.Estatus.Hex()))
			SAlmacenes.Almacen.EConexionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboConexiones(AlmacenNuevo.Conexion.Hex()))
			// SAlmacenes.Almacen.EPredecesorAlmacen.Ihtml = template.HTML(CargaCoboGrupoSelected(AlmacenNuevo.AvailableOthersIDsWithSameType(), predeceso))
			// SAlmacenes.Almacen.ESucesorAlmacen.Ihtml = template.HTML(CargaCoboGrupoAlmacenes(AlmacenNuevo.AvailableOthersIDsWithSameType(), AlmacenNuevo.Sucesor))

			if AlmacenNuevo.ReemplazaMgo() {
				if AlmacenNuevo.ReemplazaElastic() {
					SAlmacenes.SMsj = "La Actualización del almacen fue exitosa." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
					ctx.Render("AlmacenDetalle.html", SAlmacenes)
				} else {
					AlmacenNuevo.EliminaByIDMgo()
					SAlmacenes.SEstado = false
					SAlmacenes.SMsj = "Ocurrió un error al registrar su almacen: ES."
					ctx.Render("AlmacenEdita.html", SAlmacenes)
				}
			} else {
				SAlmacenes.SEstado = false
				SAlmacenes.SMsj = "No se pudo actualizar el Almacen"
				ctx.Render("AlmacenEdita.html", SAlmacenes)
			}
		} else {
			SAlmacenes.SEstado = false
			SAlmacenes.SMsj = "El almacen no se a encontrado, intente de nuevo."
			ctx.Render("AlmacenIndex.html", SAlmacenes)
		}
	} else {
		fmt.Println("No es un id")
		SAlmacenes.SEstado = false
		SAlmacenes.SMsj = "No se ha recibido un parámetro adecuado para poder editar El Almacen, intente de nuevo."
		ctx.Render("AlmacenIndex.html", SAlmacenes)
	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {

	var SAlmacenes AlmacenModel.SAlmacen

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SAlmacenes.SSesion.Name = NameUsrLoged
	SAlmacenes.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SAlmacenes.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SAlmacenes.SEstado = false
		SAlmacenes.SMsj = errSes.Error()
		ctx.Render("ZError.html", SAlmacenes)
		return
	}

	id := ctx.Param("ID")
	if bson.IsObjectIdHex(id) {
		Almacen := AlmacenModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Almacen) {
			var predecesor []bson.ObjectId
			predecesor = append(predecesor, Almacen.ID)
			SAlmacenes.SEstado = true
			SAlmacenes.ID = Almacen.ID
			SAlmacenes.ENombreAlmacen.Nombre = Almacen.Nombre
			SAlmacenes.ETipoAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(132, Almacen.Tipo.Hex()))
			SAlmacenes.EClasificacionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(133, Almacen.Clasificacion.Hex()))
			SAlmacenes.EEstatusAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(134, Almacen.Estatus.Hex()))
			SAlmacenes.EConexionAlmacen.Ihtml = template.HTML(CargaCombos.CargaComboConexiones(Almacen.Conexion.Hex()))
			// SAlmacenes.Almacen.EPredecesorAlmacen.Ihtml = template.HTML(CargaCoboGrupoSelected(Almacen.AvailableOthersIDsWithSameType(), predecesor))
			// SAlmacenes.Almacen.ESucesorAlmacen.Ihtml = template.HTML(CargaCoboGrupoSelected(Almacen.AvailableOthersIDsWithSameType(), Almacen.Sucesor))
			ctx.Render("AlmacenDetalle.html", SAlmacenes)

		} else {
			SAlmacenes.SEstado = false
			SAlmacenes.SMsj = "El almacen no se a encontrado, intente de nuevo."
			ctx.Render("AlmacenIndex.html", SAlmacenes)
		}

	} else {
		fmt.Println("No es un id")
		SAlmacenes.SEstado = false
		SAlmacenes.SMsj = "No se ha recibido un parámetro adecuado para poder editar El Almacen, intente de nuevo."
		ctx.Render("AlmacenIndex.html", SAlmacenes)
	}

}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
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

	//###### TU CÓDIGO AQUÍ PROGRAMADOR

	ctx.Render("AlmacenDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//ObtenerIDAlmacenes Realiza una busqueda de almacenes de tipo cliente en la base de datos mongo
//Regresa un arreglo de identificadores encontrados
func ObtenerIDAlmacenes(ctx *iris.Context) {
	almacenesPropios := AlmacenModel.GetEspecificsByTag("Clasificacion", bson.ObjectIdHex("58e5692ee75770120c60befa"))
	var almacenes []string
	for _, value := range almacenesPropios {
		almacenes = append(almacenes, value.ID.Hex())
		almacenes = append(almacenes, value.Nombre)
	}
	ctx.JSON(200, almacenes)
}

//ConsultarAlmacenesPostgres funcion que consulta la base de datos de postgres
//Param requiere ids de los almacenes a consultar
func ConsultarAlmacenesPostgres(ctx *iris.Context) {
	idproducto := ctx.FormValue("idproducto")
	identificadorAlmacen := ctx.FormValue("identificadorAlmacen")
	var inventario OperacionModel.InventarioPostgres
	inventario.IDProducto = idproducto
	informacion, err := inventario.ConsultaInventarioPostgres(bson.ObjectIdHex(identificadorAlmacen))
	if err != nil {
		fmt.Println(err)
	}
	if informacion.Encontrado == false {
		fmt.Println("Producto no encontrado en los almacenes")
	}
	if err == nil || informacion.Encontrado == true {
		ctx.JSON(200, informacion)
	}
}

//ConsultarExistenciaInventario funcion que consulta la base de datos de postgres
//Param requiere ids de los almacenes a consultar
func ConsultarExistenciaInventario(ctx *iris.Context) {
	idproducto := ctx.FormValue("idproducto")
	identificadorAlmacen := ctx.FormValue("identificadorAlmacen")
	existe, err := ConsultasSql.ConsultaExistenciaProductoActivo(identificadorAlmacen, idproducto)
	if err != nil {
		fmt.Println("Error en la consulta: ", err)
	}
	var mensaje = "NO EXISTE"
	if existe {
		mensaje = "EXISTE"
	}
	_, err = ctx.Writef(mensaje)
	if err != nil {
		fmt.Println("Error al mostrar la existencia: ", err)
	}
}

//MovimientoAjusteGet Renderea a la vista Get de Ajustes y Traslados
func MovimientoAjusteGet(ctx *iris.Context) {
	fmt.Println("Entrando al GET de Movimientos")
	// var almacenesExistentes []AlmacenModel.AlmacenMgo
	// almacenesExistentes = AlmacenModel.GetAll()
	//Crear una estuctura anidada SSesion en la V.
	//AjusteTraslado
	fmt.Println("Entrando MovimientoAjusteGet")
	ctx.Render("OperacionAjuste.html", nil)
}

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send AlmacenModel.SAlmacen

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

		Cabecera, Cuerpo := AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrToMongo))
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

//MovimientoTrasladoGet Renderea a la vista Post de Ajustes y Traslados
func MovimientoTrasladoGet(ctx *iris.Context) {
	fmt.Println("Entrando MovimientoTrasladoGet")
	ctx.Render("OperacionTraslado.html", nil)
}

//MuestraIndexPorGrupo regresa template de busqueda y paginacion de acuerdo a la agrupacion solicitada
func MuestraIndexPorGrupo(ctx *iris.Context) {
	var Send AlmacenModel.SAlmacen
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Almacen.ENombreAlmacen.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := AlmacenModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = AlmacenModel.GeneraTemplatesBusqueda(AlmacenModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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

	fmt.Println("Entrando MovimientosGet")
	var Almacenes []AlmacenModel.AlmacenMgo
	var AllAlmacenes AlmacenModel.MovimientoAlmacenes

	Almacenes = AlmacenModel.AlmacenFields()
	AllAlmacenes.Almacenes = Almacenes

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

	SMovimientoAlmacen.SEstado = true //Estado de la Peticion ------FALSE INDICA ERROR-----, ----TRUE INDICA EXITO-----
	SMovimientoAlmacen.SMsj = "Mensaje de Error :D"
	SMovimientoAlmacen.Almacenes = SAlmacenes
	SMovimientoAlmacen.Ihtml = template.HTML(templateAlmacenes)

	// for _, v := range SAlmacenes {
	// 	var tipoString string
	// 	tipoString = CargaCombos.CargaCatalogoByID(132, v.ETipoAlmacen.Tipo)
	// 	fmt.Printf("\n Id -> %v    Tipo -> %v ", v.ETipoAlmacen.Tipo, tipoString)

	// }

	// fmt.Println("ALMACENES MOVIMIENTOS GET", Almacenes)
	ctx.Render("Movimientos.html", SMovimientoAlmacen)
}

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
							<span class="input-group-addon">Codigo de Barra:</span>
							<input type="text" name="elarticulo" onKeydown="Javascript: if (event.keyCode==13) GetArticulo();"  id="elarticulo" class="form-control selectpicker" autofocus>
							<input type="text" hidden value="%v" id="tipomovimiento" name="tipomovimiento">
							<label class="input-group-addon" for="elarticulo">
							<span  class="glyphicon glyphicon-circle-arrow-down"></span>
							</label>													
							
							<span class="input-group-addon">
								<buttont type="button" class="btn btn-primary" data-toggle="modal" data-backdrop="static" data-keyboard="false" data-target=".bd-example-modal-lg">Buscar</button>										
							</span>																
					</div>`, whoIs)

	fmt.Fprintf(ctx.ResponseWriter, template)

}

//MovimientosAjustePost rgergergerger
func MovimientosAjustePost(ctx *iris.Context) {

	var almOrigen string
	var almDestino string
	var tipoMovimiento string
	var numeroArticulos string
	var codigoBarra string
	var cantidadOn bool

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
		case "tipomovimiento":
			for _, valor := range v {
				tipoMovimiento = valor
			}
			break
		case "cod_b_art":
			for _, valor := range v {
				codigoBarra = valor
			}
			break
		case "articulos_agregados":
			for _, valor := range v {
				numeroArticulos = valor
			}
			break
		case "agregarcantidad":
			for _, valor := range v {
				cantidadOn, _ = strconv.ParseBool(valor)

			}
		}
	}

	// fmt.Println("-> ", almOrigen)
	// fmt.Println("-> ", almDestino)
	// fmt.Println("-> ", tipoMovimiento)
	fmt.Println("-> ", numeroArticulos)
	fmt.Println("-> ", cantidadOn)
	// fmt.Println("-> ", codigoBarra)

	var Producto ProductoModel.ProductoMgo
	var Codigos ProductoModel.CodigoMgo

	Producto = ProductoModel.GetProductoMongo(codigoBarra)

	Codigos = Producto.Codigos

	if Producto.ID == "" {
		fmt.Fprintf(ctx.ResponseWriter, `<h3 class="text-center">Articulo no encontrado</h3>`)
		fmt.Println("NO ENCONTRADO")
	} else {

		var ProductoOrigen ProductoModel.ProductoPgrs
		var ProductoDestino ProductoModel.ProductoPgrs

		ProductoOrigen, ProductoDestino = ProductoModel.GetProductoPostgres(Producto.ID.Hex(), almOrigen, almDestino)

		ExistenciaOrigenString := strconv.FormatFloat(ProductoOrigen.Existencia, 'f', -1, 64)
		ExistenciaDestinoString := strconv.FormatFloat(ProductoDestino.Existencia, 'f', -1, 64)
		PrecioOrigenString := strconv.FormatFloat(ProductoOrigen.Precio, 'f', -1, 64)
		PrecioDestinoString := strconv.FormatFloat(ProductoDestino.Precio, 'f', -1, 64)

		fmt.Println("Existe Origen", ExistenciaOrigenString)
		fmt.Println("Existe Destino", ExistenciaDestinoString)

		tAjuste := `	<tr class="renglon">
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

		tAjuste += `    <td><span class="btn btn-danger eliminar">-</span></td>
						</tr>`

		tTraslado := `<tr class="renglon">
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

		tTraslado += `		<td><span class="btn btn-danger eliminar">-</span></td>
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

//RealizarMovimiento rtgrtgtrgtgs
func RealizarMovimiento(ctx *iris.Context) {

	ctx.Request.ParseForm()
	fmt.Println("===================================================================================")

	var Operacion OperacionModel.OperacionMgo

	var Movimiento OperacionModel.MovimientoMgo
	// var Movimientos []OperacionModel.MovimientoMgo

	// var Transaccion OperacionModel.TransaccionMgo
	// var Transacciones []OperacionModel.TransaccionMgo

	// var KardexPosgres OperacionModel.KardexPostgres

	var AlmacenOrigen string
	var AlmacenDestino string

	Operacion.ID = bson.NewObjectId()
	Operacion.UsuarioOrigen = bson.NewObjectId()
	Operacion.UsuarioDestino = bson.NewObjectId()
	Operacion.FechaHoraRegistro = time.Now()

	Movimiento.IDMovimiento = bson.NewObjectId()

	AlmacenOrigen = ctx.Request.FormValue("origen")
	AlmacenDestino = ctx.Request.FormValue("destino")

	var codigos = []string{}
	var nombres = []string{}
	// var precios = []string{}
	var existencias = []string{}
	var existenciasd = []string{}
	var operaciones = []string{}
	var cantidades = []string{}

	for k, v := range ctx.Request.PostForm {

		if k == "codigos[]" {
			codigos = v
		}
		if k == "nombres[]" {
			nombres = v
		}
		// if k == "precios[]" {
		// 	precios = v
		// }
		if k == "existencias[]" {
			existencias = v
		}
		if k == "operaciones[]" && AlmacenOrigen == AlmacenDestino {
			operaciones = v
		}
		if k == "cantidades[]" {
			cantidades = v
		}
		if k == "existenciasd[]" && AlmacenOrigen != AlmacenDestino {
			existenciasd = v
		}
	}

	fmt.Println("codigos  ->", codigos)
	fmt.Println("nombres  ->", nombres)
	fmt.Println("existencias  ->", existencias)
	fmt.Println("operaciones  ->", operaciones)
	fmt.Println("cantidades  ->", cantidades)
	fmt.Println("existenciasd  ->", existenciasd)

	lim := len(codigos)

	fmt.Println("limite  ->", lim)

	var TotalAlmacenOrigen float64
	var TotalAlmacenDestino float64

	// var movimiento_ajuste Movimiento
	// var producto_ajuste Producto
	// var productos_ajuste []Producto

	// var movimiento_traslado MovimientoA
	// var producto_traslado ProductoA
	// var productos_traslado []ProductoA

	if AlmacenDestino == AlmacenOrigen {

		fmt.Println("A J U S T E S")

		for i := 0; i < lim; i++ {

			nombresito := nombres[i]
			IDProducto := codigos[i]
			cantidad := cantidades[i]
			existe := existencias[i]
			operacion := operaciones[i]

			fmt.Println(nombresito)

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

			// 	producto_ajuste = Producto{ID: object_articulo, CodigoDeBarra: codigodebarra, Nombre: nombresito, CantidadA: existe_float, CantidadB: cantidad_float, CantidadC: total, Operacion: operacion}
			// 	productos_ajuste = append(productos_ajuste, producto_ajuste)
			// 	ajustes_m.RealizarAjuste(origen, articulo_mongo_id_string, total)

		}

		// movimiento_ajuste = Movimiento{Tipo: "Ajuste", Date: time.Now(), AlmacenOrigen: object_origen, AlmacenDestino: object_destino, Usuario: usuario_mongo.Id, Productos: productos_ajuste}
		// fmt.Println(movimiento_ajuste)
		// insertado := ajustes_m.GuardarMovimiento(movimiento_ajuste)
		// if insertado == true {
		// 	fmt.Println("Movimiento de ajuste realizado")
		// }
	} else {

		fmt.Println("T R A S L A D O S")

		for i := 0; i < lim; i++ {

			nombresito := nombres[i]
			IDProducto := codigos[i]
			cantidad := cantidades[i]
			existe := existencias[i]
			existed := existenciasd[i]
			ExistenciaOrigen, _ := strconv.ParseFloat(existe, 64)
			ExistenciaDestino, _ := strconv.ParseFloat(existed, 64)
			cantidadSolicitada, _ := strconv.ParseFloat(cantidad, 64)

			ProductoOrigen, ProductoDestino := ProductoModel.GetProductoPostgres(IDProducto, AlmacenOrigen, AlmacenDestino)

			TotalAlmacenDestino = ProductoDestino.Existencia + cantidadSolicitada
			TotalAlmacenOrigen = ProductoOrigen.Existencia - cantidadSolicitada

			fmt.Println("Nombresito  ", nombresito)
			fmt.Println("IDProducto  ", IDProducto)
			fmt.Println("cantidad  ", cantidad)
			fmt.Println("existe  ", existe)
			fmt.Println("existed  ", existed)

			fmt.Println("ExistenciaOrigen  ", ExistenciaOrigen)
			fmt.Println("ExistenciaDestino  ", ExistenciaDestino)
			fmt.Println("cantidadSolicitada  ", cantidadSolicitada)

			fmt.Println("TotalAlmacenDestino  ", TotalAlmacenDestino)
			fmt.Println("TotalAlmacenOrigen  ", TotalAlmacenOrigen)

			// 	producto_traslado = ProductoA{ID: object_articulo, CodigoDeBarra: codigodebarra, Nombre: nombresito, CantidadA: existe_float, CantidadB: existed_float, CantidadC: cantidad_float, CantidadD: total, CantidadE: totald}

			// 	productos_traslado = append(productos_traslado, producto_traslado)

			if TotalAlmacenOrigen < 0 {
				fmt.Fprintf(ctx.ResponseWriter, "no")
			} else {
				opExitosa, errores := ProductoModel.RealizarTraslado(AlmacenOrigen, AlmacenDestino, IDProducto, TotalAlmacenOrigen, TotalAlmacenDestino)

				fmt.Println("Operacion Exitosa", opExitosa)
				fmt.Println(errores)
			}

		}
		// movimiento_traslado = MovimientoA{Tipo: "Traslado", Date: time.Now(), AlmacenOrigen: object_origen, AlmacenDestino: object_destino, Usuario: usuario_mongo.Id, Productos: productos_traslado}
		// fmt.Println(movimiento_traslado)
		// insertado := ajustes_m.GuardarMovimiento(movimiento_traslado)
		// if insertado == true {
		// 	fmt.Println("Movimiento de traslado realizado")
		// }
	}
}

// GetTemplateAvailableOthersIds Genera el template para un permitir a un nuevo elemento elegir antecesor y sucesor
func GetTemplateAvailableOthersIds(objetos []AlmacenModel.AlmacenMgo, selectedID string) string {
	OptionsString := `<option value ="" >---Selecciona una Opcion---</option>`
	for _, opcion := range objetos {
		if selectedID != opcion.ID.Hex() {
			OptionsString = OptionsString + `<option value ="` + opcion.ID.Hex() + `" >` + opcion.Nombre + `</option>`
		} else {
			OptionsString = OptionsString + `<option value ="` + opcion.ID.Hex() + `" selected>` + opcion.Nombre + `</option>`
		}
	}
	return OptionsString
}

// CargaCoboGrupoAlmacenes Obtiene el combo de almacenes y marca los seleccionados
func CargaCoboGrupoAlmacenes(objectiveIDs []AlmacenModel.AlmacenMgo, selectedIDs []bson.ObjectId) string {
	OptionsString := `<option value ="" >---Selecciona una Opcion---</option>`
	if len(selectedIDs) > 0 {
		for _, objective := range objectiveIDs {
			if IndexOfID(selectedIDs, objective.ID) != -1 {
				OptionsString = OptionsString + `<option value ="` + objective.ID.Hex() + `" selected >` + objective.Nombre + `</option>`

			} else {
				OptionsString = OptionsString + `<option value ="` + objective.ID.Hex() + `"  >` + objective.Nombre + `</option>`

			}
		}
	} else {
		for _, objective := range objectiveIDs {
			OptionsString = OptionsString + `<option value ="` + objective.ID.Hex() + `"  >` + objective.Nombre + `</option>`
		}
	}
	return OptionsString
}

// CargaCoboGrupoSelected Obtiene el combo de almacenes  los seleccionados
func CargaCoboGrupoSelected(objectiveIDs []AlmacenModel.AlmacenMgo, selectedIDs []bson.ObjectId) string {
	OptionsString := ``
	if len(selectedIDs) > 0 {
		for _, objective := range objectiveIDs {
			if IndexOfID(selectedIDs, objective.ID) != -1 {
				OptionsString = OptionsString + `<option value ="` + objective.ID.Hex() + `" selected >` + objective.Nombre + `</option>`
			}
		}
	}
	return OptionsString
}

// IndexOfID regresa el indice del elemento TestingObject en los ids de almacen seleccionados SelectedAlmacens
func IndexOfID(SelectedIDs []bson.ObjectId, TestingObject bson.ObjectId) int {
	for i, selected := range SelectedIDs {
		if selected == TestingObject {
			return i
		}
	}
	return -1
}
