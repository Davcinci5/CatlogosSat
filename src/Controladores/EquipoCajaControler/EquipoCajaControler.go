package EquipoCajaControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"../../Modelos/EquipoCajaModel"

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

//limitePorPagina limite de registros a mostrar en la pagina
var limitePorPagina = 5

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de EquipoCaja
func IndexGet(ctx *iris.Context) {
	var Send EquipoCajaModel.SEquipoCaja

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
	numeroRegistros = EquipoCajaModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	EquipoCajas := EquipoCajaModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range EquipoCajas {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajas[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajas[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("EquipoCajaIndex.html", Send)
}

//IndexPost regresa la peticon post que se hizo desde el index de EquipoCaja
func IndexPost(ctx *iris.Context) {

	var Send EquipoCajaModel.SEquipoCaja
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
	//Send.EquipoCaja.EVARIABLEEquipoCaja.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := EquipoCajaModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("EquipoCajaIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de EquipoCaja
func AltaGet(ctx *iris.Context) {

	var Send EquipoCajaModel.SEquipoCaja
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

	ctx.Render("EquipoCajaAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de EquipoCaja
func AltaPost(ctx *iris.Context) {

	var SEquipoCaja EquipoCajaModel.SEquipoCaja

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SEquipoCaja.SSesion.Name = NameUsrLoged
	SEquipoCaja.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SEquipoCaja.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SEquipoCaja.SEstado = false
		SEquipoCaja.SMsj = errSes.Error()
		ctx.Render("ZError.html", SEquipoCaja)
		return
	}

	//######### LEE TU OBJETO DEL FORMULARIO #########
	var EquipoCaja EquipoCajaModel.EquipoCajaMgo
	ctx.ReadForm(&EquipoCaja)
	fmt.Println(EquipoCaja)
	//Asigno de manera fija el estatus, el objectID pertenece a "ACTIVO"
	IDEstatus := bson.ObjectIdHex("58e577d7e75770120c60bf26")
	EquipoCaja.Estatus = IDEstatus
	fmt.Println(EquipoCaja)
	//######### VALIDA TU OBJETO #########
	EstatusPeticion := false //True indica que hay un error
	//##### TERMINA TU VALIDACION ########

	//########## Asigna vairables a la estructura que enviarás a la vista
	ID := bson.NewObjectId()
	EquipoCaja.ID = ID
	//EquipoCaja.Nombre = ctx.FormValue("Nombre")
	//######### ENVIA TUS RESULTADOS #########

	SEquipoCaja.EquipoCaja.ID = ID
	SEquipoCaja.EquipoCaja.ENombreEquipoCaja.Nombre = ctx.FormValue("Nombre")
	SEquipoCaja.EquipoCaja.EDescripcionEquipoCaja.Descripcion = ctx.FormValue("Descripcion")
	SEquipoCaja.EquipoCaja.EEstatusEquipoCaja.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(145, IDEstatus.Hex()))
	//	SEquipoCaja.EquipoCaja = EquipoCaja //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.

	if EstatusPeticion {
		SEquipoCaja.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SEquipoCaja.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("EquipoCajaAlta.html", SEquipoCaja)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if EquipoCaja.InsertaMgo() {
			SEquipoCaja.SEstado = true
			SEquipoCaja.SMsj = "Se ha realizado una inserción exitosa"

			//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
			ctx.Render("EquipoCajaDetalle.html", SEquipoCaja)

		} else {
			SEquipoCaja.SEstado = false
			SEquipoCaja.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("EquipoCajaAlta.html", SEquipoCaja)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de EquipoCaja
func EditaGet(ctx *iris.Context) {
	var Send EquipoCajaModel.SEquipoCaja

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
	idEc := ctx.Param("ID")
	if bson.IsObjectIdHex(idEc) {
		Caja := EquipoCajaModel.GetOne(bson.ObjectIdHex(idEc))
		if !MoGeneral.EstaVacio(Caja) {

			Send.EquipoCaja.ID = Caja.ID
			Send.EquipoCaja.ENombreEquipoCaja.Nombre = Caja.Nombre
			Send.EquipoCaja.EDescripcionEquipoCaja.Descripcion = Caja.Descripcion
			Send.EquipoCaja.EEstatusEquipoCaja.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(145, Caja.Estatus.Hex()))
			Send.SEstado = true
			//Send.SMsj = "Se ha realizado
			ctx.Render("EquipoCajaEdita.html", Send)
		} else {
			Send.SEstado = false
			Send.SMsj = "El Equipo Caja no se a encontrado, intente de nuevo."
			ctx.Render("EquipoCajaIndex.html", Send)
		}
	} else {
		fmt.Println("No es un id")
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar El Equipo Caja, intente de nuevo."
		ctx.Render("EquipoCajaIndex.html", Send)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de EquipoCaja
func EditaPost(ctx *iris.Context) {
	var Send EquipoCajaModel.SEquipoCaja
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
	idEc := ctx.Param("ID")
	if bson.IsObjectIdHex(idEc) {
		Caja := EquipoCajaModel.GetOne(bson.ObjectIdHex(idEc))
		if !MoGeneral.EstaVacio(Caja) {
			Caja.Nombre = ctx.FormValue("Nombre")
			Caja.Descripcion = ctx.FormValue("Descripcion")
			Caja.Estatus = bson.ObjectIdHex(ctx.FormValue("Estatus"))
			actualizado := Caja.ActualizaMgo([]string{"Nombre", "Descripcion", "Estatus"}, []interface{}{Caja.Nombre, Caja.Descripcion, Caja.Estatus})
			if !actualizado {
				Send.SEstado = false
				Send.SMsj = "No se pudo actualizar el Equipo Caja"
				ctx.Render("AlmacenIndex.html", Send)
			}
			Send.EquipoCaja.ID = Caja.ID
			Send.EquipoCaja.ENombreEquipoCaja.Nombre = Caja.Nombre
			Send.EquipoCaja.EDescripcionEquipoCaja.Descripcion = Caja.Descripcion
			Send.EquipoCaja.EEstatusEquipoCaja.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(145, Caja.Estatus.Hex()))
			Send.SEstado = true
			ctx.Render("EquipoCajaDetalle.html", Send)
		} else {
			Send.SEstado = false
			Send.SMsj = "El Equipo Caja no se a encontrado, intente de nuevo."
			ctx.Render("EquipoCajaIndex.html", Send)
		}

	} else {
		fmt.Println("No es un id")
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar El Equipo Caja, intente de nuevo."
		ctx.Render("EquipoCajaIndex.html", Send)
	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send EquipoCajaModel.SEquipoCaja

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
	idEc := ctx.Param("ID")
	if bson.IsObjectIdHex(idEc) {
		Caja := EquipoCajaModel.GetOne(bson.ObjectIdHex(idEc))
		if !MoGeneral.EstaVacio(Caja) {
			Send.EquipoCaja.ID = Caja.ID
			Send.EquipoCaja.ENombreEquipoCaja.Nombre = Caja.Nombre
			Send.EquipoCaja.EDescripcionEquipoCaja.Descripcion = Caja.Descripcion
			Send.EquipoCaja.EEstatusEquipoCaja.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(145, Caja.Estatus.Hex()))
			Send.SEstado = true
			ctx.Render("EquipoCajaDetalle.html", Send)
		} else {
			Send.SEstado = false
			Send.SMsj = "El Equipo Caja no se a encontrado, intente de nuevo."
			ctx.Render("EquipoCajaIndex.html", Send)
		}

	} else {
		fmt.Println("No es un id")
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar El Equipo Caja, intente de nuevo."
		ctx.Render("EquipoCajaIndex.html", Send)
	}

}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send EquipoCajaModel.SEquipoCaja
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

	ctx.Render("EquipoCajaDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send EquipoCajaModel.SEquipoCaja

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

		Cabecera, Cuerpo := EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrToMongo))
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
	var Send EquipoCajaModel.SEquipoCaja
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.EquipoCaja.ENombreEquipoCaja.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := EquipoCajaModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = EquipoCajaModel.GeneraTemplatesBusqueda(EquipoCajaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
