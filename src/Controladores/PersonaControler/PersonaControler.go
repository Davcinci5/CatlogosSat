package PersonaControler

import (
	"encoding/json"
	"html/template"
	"strconv"
	"time"

	"../../Modelos/GrupoPersonaModel"
	"../../Modelos/PersonaModel"

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
var limitePorPagina = 10

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Persona
func IndexGet(ctx *iris.Context) {

	var Send PersonaModel.SPersona
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
	numeroRegistros = PersonaModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Personas := PersonaModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Personas {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(Personas[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(Personas[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("PersonaIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Persona
func IndexPost(ctx *iris.Context) {

	var Send PersonaModel.SPersona
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
	//Send.Persona.EVARIABLEPersona.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := PersonaModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("PersonaIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Persona
func AltaGet(ctx *iris.Context) {

	var Send PersonaModel.SPersona
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

	Send.Persona.ETipoPersona.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(159, ""))
	Send.Persona.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasMulti(""))
	Send.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(""))

	ctx.Render("PersonaAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Persona
func AltaPost(ctx *iris.Context) {
	var SPersona PersonaModel.SPersona
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SPersona.SSesion.Name = NameUsrLoged
	SPersona.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SPersona.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SPersona.SEstado = false
		SPersona.SMsj = errSes.Error()
		ctx.Render("ZError.html", SPersona)
		return
	}

	//######### LEE TU OBJETO DEL FORMULARIO #########
	var Persona PersonaModel.PersonaMgo

	ctx.ReadForm(&Persona)
	grupos := ctx.Request.PostForm["Grupos"]
	Persona.Grupos = CargaCombos.ArrayStringToObjectID(grupos)

	EstatusPeticion := false //True indica que hay un error

	Persona.ID = bson.NewObjectId()
	Persona.FechaHora = time.Now()

	//Rellenar Campo Notificaciones
	obj1 := bson.NewObjectId()
	obj2 := bson.NewObjectId()

	Persona.Notificacion = append(Persona.Notificacion, obj1)
	Persona.Notificacion = append(Persona.Notificacion, obj2)

	Persona.Estatus = CargaCombos.CargaEstatusActivoEnAlta(160)

	SPersona.EEstatusPersona.Estatus = Persona.Estatus
	SPersona.EEstatusPersona.IEstatus = false
	SPersona.EEstatusPersona.Ihtml = `<input type="text" name="Estatus" id="Estatus" class="form-control" value="ACTIVO" readonly>`

	//Campo Nombre
	if Persona.Nombre == "" {
		//Si viene vacio el nombre es ERROR por lo tanto SPersona.ENombrePersona.IEstatus = true
		SPersona.ENombrePersona.IEstatus = true
		SPersona.ENombrePersona.IMsj = "El nombre esta vacio"
		EstatusPeticion = true

	} else {
		SPersona.ENombrePersona.Nombre = Persona.Nombre
		SPersona.ENombrePersona.IEstatus = false
	}

	//Campo Tipos
	if Persona.Tipo == nil {
		SPersona.ETipoPersona.IEstatus = true
		SPersona.ETipoPersona.IMsj = "El campo Tipo es obligatorio"
		SPersona.ETipoPersona.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(159, ""))
		EstatusPeticion = true
	} else {
		SPersona.ETipoPersona.Tipo = Persona.Tipo
		SPersona.ETipoPersona.IEstatus = false
		SPersona.ETipoPersona.Ihtml = template.HTML(PersonaModel.CargaComboCatalogoArrayID(159, Persona))
	}

	//Campo Grupos
	if Persona.Grupos == nil {
		SPersona.EGruposPersona.IEstatus = false
		SPersona.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(grupos))
	} else {
		SPersona.EGruposPersona.Grupos = Persona.Grupos
		SPersona.EGruposPersona.IEstatus = false
		SPersona.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(grupos))
	}

	//Campo Predecesor
	if Persona.Predecesor == "" {
		SPersona.EPredecesorPersona.IMsj = "Sin predecesor"
		SPersona.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(""))

	} else {
		SPersona.EPredecesorPersona.Predecesor = Persona.Predecesor
		SPersona.EPredecesorPersona.IEstatus = false
		SPersona.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.CargaNombrePredecesor(Persona.Predecesor))

	}

	//Campo Notificaciones
	if Persona.Notificacion == nil {
		SPersona.ENotificacionPersona.IEstatus = true
		SPersona.ENotificacionPersona.IMsj = "Sin notificaciones"

	} else {
		SPersona.ENotificacionPersona.IEstatus = false
		SPersona.ENotificacionPersona.Notificacion = Persona.Notificacion
		SPersona.ENotificacionPersona.Ihtml = `<option>Si tiene Notificaciones</option>`
	}

	SPersona.ID = Persona.ID
	SPersona.EFechaHoraPersona.FechaHora = Persona.FechaHora
	SPersona.EFechaHoraPersona.Ihtml = template.HTML(`<input type="text" name="FechaHora" id="FechaHora" class="form-control" value="` + Persona.FechaHora.Format("Mon Jan _2 15:04:05 2006") + `" readonly>`)

	if EstatusPeticion {
		SPersona.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SPersona.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("PersonaAlta.html", SPersona)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if Persona.InsertaMgo() {
			SPersona.SEstado = true
			SPersona.SMsj = "Se ha realizado una inserción exitosa"

			//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
			ctx.Render("PersonaDetalle.html", SPersona)

		} else {
			SPersona.SEstado = false
			SPersona.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde =)!"
			ctx.Render("PersonaAlta.html", SPersona)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Persona
func EditaGet(ctx *iris.Context) {
	var Send PersonaModel.SPersona
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

	ctx.Render("PersonaEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Persona
func EditaPost(ctx *iris.Context) {

	var Send PersonaModel.SPersona
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

	ctx.Render("PersonaEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send PersonaModel.SPersona
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

	ctx.Render("PersonaDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send PersonaModel.SPersona
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

	ctx.Render("PersonaDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send PersonaModel.SPersona

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

		Cabecera, Cuerpo := PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrToMongo))
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
	var Send PersonaModel.SPersona
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Persona.ENombrePersona.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := PersonaModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = PersonaModel.GeneraTemplatesBusqueda(PersonaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
