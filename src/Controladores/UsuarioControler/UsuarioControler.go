package UsuarioControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"../../Modelos/CatalogoModel"

	"../../Modulos/Session"

	"../../Modelos/EquipoCajaModel"
	"../../Modelos/GrupoPersonaModel"
	"../../Modelos/PersonaModel"
	"../../Modelos/UsuarioModel"

	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/IndexUsuarios"

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

var CatalogotipoPersonas = 159

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Usuario
func IndexGet(ctx *iris.Context) {

	var Send UsuarioModel.SUsuario
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
		Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(Usuarios[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(Usuarios[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true
	ctx.Render("UsuarioIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Usuario
func IndexPost(ctx *iris.Context) {

	var Send UsuarioModel.SUsuario
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

			Cabecera, Cuerpo := IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("UsuarioIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Usuario
func AltaGet(ctx *iris.Context) {

	var Send UsuarioModel.SUsuario
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

	Send.Usuario.EPersonaUsuario.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasMulti(""))
	Send.Usuario.EPersonaUsuario.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(""))

	Send.Usuario.ECajasUsuario.Ihtml = template.HTML(EquipoCajaModel.CargaComboCajasMulti())
	ctx.Render("UsuarioAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Usuario
func AltaPost(ctx *iris.Context) {
	EstatusPeticion := false
	var Send UsuarioModel.SUsuario
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

	//############ TU CÓDIGO AQUÍ
	var SendMgoPer PersonaModel.PersonaMgo
	var SendMgoUsr UsuarioModel.UsuarioMgo

	Nombre := ctx.FormValue("Nombre")
	Send.Usuario.EPersonaUsuario.ENombrePersona.Nombre = Nombre
	SendMgoPer.Nombre = Nombre
	if Nombre == "" {
		EstatusPeticion = true
		Send.Usuario.EPersonaUsuario.ENombrePersona.IEstatus = true
		Send.Usuario.EPersonaUsuario.ENombrePersona.IMsj = "El campo Nombre no debe estar vacio"
	}
	SendMgoPer.Tipo = append(SendMgoPer.Tipo, CargaCombos.CargaIDTipoPersonaPorDefecto(CatalogotipoPersonas))

	Grupos := ctx.Request.Form["Grupos"]
	Send.Usuario.EPersonaUsuario.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(Grupos))
	var arrGpo []bson.ObjectId
	for _, val := range Grupos {
		arrGpo = append(arrGpo, bson.ObjectIdHex(val))
	}
	SendMgoPer.Grupos = arrGpo

	Predecesor := ctx.FormValue("Predecesor")
	Send.Usuario.EPersonaUsuario.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(Predecesor))
	if Predecesor != "" {
		if bson.IsObjectIdHex(Predecesor) {
			Send.Usuario.EPersonaUsuario.EPredecesorPersona.Predecesor = bson.ObjectIdHex(Predecesor)
			SendMgoPer.Predecesor = bson.ObjectIdHex(Predecesor)
		} else {
			EstatusPeticion = true
			Send.Usuario.EPersonaUsuario.EPredecesorPersona.IEstatus = true
			Send.Usuario.EPersonaUsuario.EPredecesorPersona.IMsj = "Valor incorrecto para predecesor"
		}
	}
	Usuario := ctx.FormValue("Usuario")
	Send.Usuario.EUsuarioUsuario.Usuario = Usuario
	SendMgoUsr.Usuario = Usuario
	UsuarioSys := UsuarioModel.GetEspecificByFields("Usuario", Usuario)
	if UsuarioSys.Usuario == Usuario {
		EstatusPeticion = true
		Send.Usuario.EUsuarioUsuario.IEstatus = true
		Send.Usuario.EUsuarioUsuario.IMsj = "El nombre de Usuario ya existe"
	}

	Contraseña := ctx.FormValue("Contraseña")
	Send.Usuario.ECredencialesUsuario.Contraseña.Contraseña = Contraseña
	SendMgoUsr.Credenciales.Contraseña = Contraseña

	Pin := ctx.FormValue("Pin")
	Send.Usuario.ECredencialesUsuario.Pin.Pin = Pin
	SendMgoUsr.Credenciales.Pin = Pin

	CodigoBarra := ctx.FormValue("CodigoBarra")
	Send.Usuario.ECredencialesUsuario.CodigoBarra.CodigoBarra = CodigoBarra
	SendMgoUsr.Credenciales.CodigoBarra = CodigoBarra

	Huella := ctx.FormValue("Huella")
	Send.Usuario.ECredencialesUsuario.Huella.Huella = Huella
	SendMgoUsr.Credenciales.Huella = Huella

	Cajas := ctx.Request.Form["Cajas"]
	Send.Usuario.ECajasUsuario.Ihtml = template.HTML(EquipoCajaModel.CargaComboCajasMultiArrayObjID(Cajas))
	var arrCajas []bson.ObjectId
	for _, val := range Cajas {
		arrCajas = append(arrCajas, bson.ObjectIdHex(val))
	}
	SendMgoUsr.Cajas = arrCajas

	CorreoPrincipal := ctx.FormValue("CorreoPrincipal")
	Send.Usuario.EMediosDeContactoUsuario.ECorreos.Principal = CorreoPrincipal
	SendMgoUsr.MediosDeContacto.Correos.Principal = CorreoPrincipal
	if CorreoPrincipal == "" {
		EstatusPeticion = true
		Send.Usuario.EMediosDeContactoUsuario.ECorreos.IEstatus = true
		Send.Usuario.EMediosDeContactoUsuario.ECorreos.IMsj = "Debe seleccionar un correo como principal"
	}

	Correos := ctx.Request.Form["Correos"]
	Send.Usuario.EMediosDeContactoUsuario.ECorreos.Ihtml = template.HTML(GeneraTemplateCorreos(Correos, CorreoPrincipal))
	SendMgoUsr.MediosDeContacto.Correos.Correos = Correos

	if len(Correos) == 0 {
		EstatusPeticion = true
		Send.Usuario.EMediosDeContactoUsuario.ECorreos.IEstatus = true
		Send.Usuario.EMediosDeContactoUsuario.ECorreos.IMsj = "Debe agregar al menos un correo"
	}

	TelefonoPrincipal := ctx.FormValue("TelefonoPrincipal")
	Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Principal = TelefonoPrincipal
	SendMgoUsr.MediosDeContacto.Telefonos.Principal = TelefonoPrincipal

	Telefonos := ctx.Request.Form["Telefonos"]
	Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Ihtml = template.HTML(GeneraTemplateTelefonos(Telefonos, TelefonoPrincipal))
	SendMgoUsr.MediosDeContacto.Telefonos.Telefonos = Telefonos

	Otros := ctx.Request.Form["Otros"]
	Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.OtrosMedios = Otros
	Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.Ihtml = template.HTML(GeneraTemplateOtros(Otros))
	SendMgoUsr.MediosDeContacto.Otros = Otros

	SendMgoPer.Estatus = CatalogoModel.RegresaIDEstatusActivo(160)
	SendMgoUsr.Estatus = CatalogoModel.RegresaIDEstatusActivo(167)

	if EstatusPeticion {
		Send.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		Send.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.

	} else {
		SendMgoPer.ID = bson.NewObjectId()
		SendMgoUsr.ID = bson.NewObjectId()
		SendMgoUsr.IDPersona = SendMgoPer.ID
		fmt.Println("Persona ID", SendMgoPer.ID)
		fmt.Println("Usuario ID", SendMgoUsr.ID)
		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		/*
		   ############################################################
		   		var insertado bool

		   		if SendMgoPer.InsertaMgo() {
		   			insertado = true
		   		}

		   		if !insertado{
		   			if SendMgoUsr.InsertaMgo() {
		   				roolback(SendMgoPer)
		   			}
		   		}

		   		if SendMgoPer.InsertaElastic()
		   			if !insertado{
		   				roolback(SendMgoPer)
		   				roolback(SendMgoUsr)
		   			}
		   		}
		   ############################################################
		*/

		if SendMgoPer.InsertaMgo() {
			if SendMgoUsr.InsertaMgo() {
				if SendMgoPer.InsertaElastic() {
					if SendMgoUsr.InsertaElastic() {
						Send.SEstado = true
						Send.SMsj = "Se ha realizado una inserción exitosa"
						ctx.Redirect("/Usuarios/detalle/"+SendMgoUsr.ID.Hex(), 301)
					} else {
						Send.SEstado = false
						Send.SMsj = "Ocurrió un error al insertar Usuario en elasticSearch"

					}

				} else {
					Send.SEstado = false
					Send.SMsj = "Ocurrió un error al insertar Persona en elasticSearch"
				}

			} else {
				eliminado := SendMgoPer.EliminaByIDMgo()
				Send.SEstado = false
				Send.SMsj = "Ocurrió un error al insertar Usuario en MongoDb"
				if !eliminado {
					Send.SMsj += "No se ha podido eliminar el objeto"
				}
			}

		} else {
			Send.SEstado = false
			Send.SMsj = "Ocurrió un error al insertar Persona en MongoDb"

		}
	}
	ctx.Render("UsuarioAlta.html", Send)

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Usuario
func EditaGet(ctx *iris.Context) {

	var Send UsuarioModel.SUsuario
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
	idusr := ctx.Param("ID")
	fmt.Println("iduser:", idusr)
	var usuario UsuarioModel.UsuarioMgo
	var persona PersonaModel.PersonaMgo
	if bson.IsObjectIdHex(idusr) {
		usuario = UsuarioModel.GetOne(bson.ObjectIdHex(idusr))
		persona = PersonaModel.GetOne(usuario.IDPersona)
		if usuario.ID.Hex() != "" {
			if persona.ID.Hex() != "" {
				Send.Usuario.ID = usuario.ID
				Send.Usuario.EPersonaUsuario.ENombrePersona.Nombre = persona.Nombre
				var tiposPersona []string
				for _, val := range persona.Tipo {
					tiposPersona = append(tiposPersona, val.Hex())
				}
				Send.Usuario.EPersonaUsuario.ETipoPersona.Ihtml = template.HTML(CargaCombos.CargaComboTiposPersonas(CatalogotipoPersonas, tiposPersona))
				var gruposPersona []string
				for _, val := range persona.Grupos {
					gruposPersona = append(gruposPersona, val.Hex())
				}
				Send.Usuario.EPersonaUsuario.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(gruposPersona))
				Send.Usuario.EPersonaUsuario.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(persona.Predecesor.Hex()))
				Send.Usuario.EUsuarioUsuario.Usuario = usuario.Usuario

				var cajasPersona []string
				for _, val := range usuario.Cajas {
					cajasPersona = append(cajasPersona, val.Hex())
				}
				Send.Usuario.ECajasUsuario.Ihtml = template.HTML(EquipoCajaModel.CargaComboCajasMultiArrayObjID(cajasPersona))
				Send.Usuario.EMediosDeContactoUsuario.ECorreos.Principal = usuario.MediosDeContacto.Correos.Principal
				Send.Usuario.EMediosDeContactoUsuario.ECorreos.Ihtml = template.HTML(GeneraTemplateCorreos(usuario.MediosDeContacto.Correos.Correos, usuario.MediosDeContacto.Correos.Principal))
				Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Principal = usuario.MediosDeContacto.Telefonos.Principal
				Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Ihtml = template.HTML(GeneraTemplateTelefonos(usuario.MediosDeContacto.Telefonos.Telefonos, usuario.MediosDeContacto.Telefonos.Principal))
				Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.Ihtml = template.HTML(GeneraTemplateOtros(usuario.MediosDeContacto.Otros))
			} else {
				Send.SEstado = false
				Send.SMsj = "El Usuario no se ha encontrado"
			}

		} else {
			Send.SEstado = false
			Send.SMsj = "El Usuario no se ha encontrado"
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "Error en la referencia al Usuario"
	}
	ctx.Render("UsuarioEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Usuario
func EditaPost(ctx *iris.Context) {
	var Send UsuarioModel.SUsuario
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

	idusr := ctx.FormValue("IDname")
	fmt.Println("iduser:", idusr)
	if bson.IsObjectIdHex(idusr) {
		var SendMgoUsrExis, SendMgoUsrNew UsuarioModel.UsuarioMgo
		var SendMgoPerExis, SendMgoPerNew PersonaModel.PersonaMgo
		SendMgoUsrExis = UsuarioModel.GetOne(bson.ObjectIdHex(idusr))
		SendMgoPerExis = PersonaModel.GetOne(SendMgoUsrExis.IDPersona)
		SendMgoUsrNew = SendMgoUsrExis
		SendMgoPerNew = SendMgoPerExis
		if SendMgoUsrExis.ID.Hex() != "" {
			if SendMgoPerExis.ID.Hex() != "" {
				Send.Usuario.ID = SendMgoPerNew.ID

				Nombre := ctx.FormValue("Nombre")
				Send.Usuario.EPersonaUsuario.ENombrePersona.Nombre = Nombre
				SendMgoPerNew.Nombre = Nombre
				if Nombre == "" {
					EstatusPeticion = true
					Send.Usuario.EPersonaUsuario.ENombrePersona.IEstatus = true
					Send.Usuario.EPersonaUsuario.ENombrePersona.IMsj = "El campo Nombre no debe estar vacio"
				}

				Tipo := ctx.Request.Form["Tipo"]
				Send.Usuario.EPersonaUsuario.ETipoPersona.Ihtml = template.HTML(CargaCombos.CargaComboTiposPersonas(CatalogotipoPersonas, Tipo))
				var arrTipo []bson.ObjectId
				for _, val := range Tipo {
					arrTipo = append(arrTipo, bson.ObjectIdHex(val))
				}
				SendMgoPerNew.Tipo = arrTipo
				if len(Tipo) == 0 {
					EstatusPeticion = true
					Send.Usuario.EPersonaUsuario.ETipoPersona.IEstatus = true
					Send.Usuario.EPersonaUsuario.ETipoPersona.IMsj = "Debe seleccionar al menos un Tipo"
				}

				Grupos := ctx.Request.Form["Grupos"]
				Send.Usuario.EPersonaUsuario.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(Grupos))
				var arrGpo []bson.ObjectId
				for _, val := range Grupos {
					arrGpo = append(arrGpo, bson.ObjectIdHex(val))
				}
				SendMgoPerNew.Grupos = arrGpo

				Predecesor := ctx.FormValue("Predecesor")
				Send.Usuario.EPersonaUsuario.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(Predecesor))
				if Predecesor != "" {
					if bson.IsObjectIdHex(Predecesor) {
						Send.Usuario.EPersonaUsuario.EPredecesorPersona.Predecesor = bson.ObjectIdHex(Predecesor)
						SendMgoPerNew.Predecesor = bson.ObjectIdHex(Predecesor)
					} else {
						EstatusPeticion = true
						Send.Usuario.EPersonaUsuario.EPredecesorPersona.IEstatus = true
						Send.Usuario.EPersonaUsuario.EPredecesorPersona.IMsj = "Valor incorrecto para predecesor"
					}
				}

				Usuario := ctx.FormValue("Usuario")
				UsuarioSys := UsuarioModel.GetEspecificByFields("Usuario", Usuario)
				if UsuarioSys.Usuario == Usuario {
					if SendMgoUsrNew.Usuario != UsuarioSys.Usuario {
						EstatusPeticion = true
						Send.Usuario.EUsuarioUsuario.IEstatus = true
						Send.Usuario.EUsuarioUsuario.IMsj = "El nombre de Usuario ya existe"
					}
				}

				Send.Usuario.EUsuarioUsuario.Usuario = Usuario
				SendMgoUsrNew.Usuario = Usuario

				Cajas := ctx.Request.Form["Cajas"]
				Send.Usuario.ECajasUsuario.Ihtml = template.HTML(EquipoCajaModel.CargaComboCajasMultiArrayObjID(Cajas))
				var arrCajas []bson.ObjectId
				for _, val := range Cajas {
					arrCajas = append(arrCajas, bson.ObjectIdHex(val))
				}
				SendMgoUsrNew.Cajas = arrCajas

				CorreoPrincipal := ctx.FormValue("CorreoPrincipal")
				Send.Usuario.EMediosDeContactoUsuario.ECorreos.Principal = CorreoPrincipal
				SendMgoUsrNew.MediosDeContacto.Correos.Principal = CorreoPrincipal
				if CorreoPrincipal == "" {
					EstatusPeticion = true
					Send.Usuario.EMediosDeContactoUsuario.ECorreos.IEstatus = true
					Send.Usuario.EMediosDeContactoUsuario.ECorreos.IMsj = "Debe seleccionar un correo como principal"
				}

				Correos := ctx.Request.Form["Correos"]
				Send.Usuario.EMediosDeContactoUsuario.ECorreos.Ihtml = template.HTML(GeneraTemplateCorreos(Correos, CorreoPrincipal))
				SendMgoUsrNew.MediosDeContacto.Correos.Correos = Correos

				if len(Correos) == 0 {
					EstatusPeticion = true
					Send.Usuario.EMediosDeContactoUsuario.ECorreos.IEstatus = true
					Send.Usuario.EMediosDeContactoUsuario.ECorreos.IMsj = "Debe agregar al menos un correo"
				}

				TelefonoPrincipal := ctx.FormValue("TelefonoPrincipal")
				Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Principal = TelefonoPrincipal
				SendMgoUsrNew.MediosDeContacto.Telefonos.Principal = TelefonoPrincipal

				Telefonos := ctx.Request.Form["Telefonos"]
				Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Ihtml = template.HTML(GeneraTemplateTelefonos(Telefonos, TelefonoPrincipal))
				SendMgoUsrNew.MediosDeContacto.Telefonos.Telefonos = Telefonos

				Otros := ctx.Request.Form["Otros"]
				Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.OtrosMedios = Otros
				Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.Ihtml = template.HTML(GeneraTemplateOtros(Otros))
				SendMgoUsrNew.MediosDeContacto.Otros = Otros

				SendMgoPerNew.Estatus = CatalogoModel.RegresaIDEstatusActivo(160)
				SendMgoUsrNew.Estatus = CatalogoModel.RegresaIDEstatusActivo(167)

				if EstatusPeticion {
					Send.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
					Send.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
				} else {

					//Actualizacion de persona en MOngoDb Y ElasticSearch.
					var personaActualizada bool
					if SendMgoPerNew.ReemplazaMgo() {
						errupd := SendMgoPerNew.ActualizaElastic()
						if errupd != nil {
							if SendMgoPerExis.ReemplazaMgo() {
								Send.SEstado = false
								Send.SMsj = "Ocurrió el siguiente error al actualizar la persona: (" + errupd.Error() + "). Se ha reestablecido la informacion"
							} else {
								Send.SEstado = false
								Send.SMsj = "Ocurrió el siguiente error al actualizar su persona: (" + errupd.Error() + ") No se pudo reestablecer la informacion"
							}
						} else {
							personaActualizada = true
						}
					} else {
						Send.SEstado = false
						Send.SMsj = "Ocurrió un error al Actualizar Persona en MongoDb"
					}

					//Actualizacion de Usuario en MOngoDb Y ElasticSearch.
					if personaActualizada {
						if SendMgoUsrNew.ReemplazaMgo() {
							errupd := SendMgoUsrNew.ActualizaElastic()
							if errupd == nil {
								Send.SEstado = true
								Send.SMsj = "Se ha realizado una Actualización exitosa"
								ctx.Redirect("/Usuarios/detalle/"+SendMgoUsrNew.ID.Hex(), 301)
							} else {
								if SendMgoUsrExis.ReemplazaMgo() {
									Send.SEstado = false
									Send.SMsj = "Ocurrió el siguiente error al actualizar la persona: (" + errupd.Error() + "). Se ha reestablecido la informacion"
								} else {
									Send.SEstado = false
									Send.SMsj = "Ocurrió el siguiente error al actualizar su persona: (" + errupd.Error() + ") No se pudo reestablecer la informacion"
								}
							}
						} else {
							Send.SEstado = false
							Send.SMsj = "Ocurrió un error al Actualizar Usuario en MongoDb, (La persona quedará registrada)"
						}
					}
				}
			} else {
				Send.SEstado = false
				Send.SMsj = "Usuario no Encontrado"
			}
		} else {
			Send.SEstado = false
			Send.SMsj = "Usuario no Encontrado"
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "Error en la referencia al Usuario"
	}

	ctx.Render("UsuarioEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send UsuarioModel.SUsuario
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
	idusr := ctx.Param("ID")
	fmt.Println("iduser:", idusr)
	var usuario UsuarioModel.UsuarioMgo
	var persona PersonaModel.PersonaMgo
	if bson.IsObjectIdHex(idusr) {
		usuario = UsuarioModel.GetOne(bson.ObjectIdHex(idusr))
		persona = PersonaModel.GetOne(usuario.IDPersona)
		if usuario.ID.Hex() != "" {
			if persona.ID.Hex() != "" {
				Send.Usuario.ID = usuario.ID
				Send.Usuario.EPersonaUsuario.ENombrePersona.Nombre = persona.Nombre
				var tiposPersona []string
				for _, val := range persona.Tipo {
					tiposPersona = append(tiposPersona, val.Hex())
				}
				Send.Usuario.EPersonaUsuario.ETipoPersona.Ihtml = template.HTML(CargaCombos.CargaComboTiposPersonas(CatalogotipoPersonas, tiposPersona))
				var gruposPersona []string
				for _, val := range persona.Grupos {
					gruposPersona = append(gruposPersona, val.Hex())
				}
				Send.Usuario.EPersonaUsuario.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(gruposPersona))
				Send.Usuario.EPersonaUsuario.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(persona.Predecesor.Hex()))
				Send.Usuario.EUsuarioUsuario.Usuario = usuario.Usuario

				var cajasPersona []string
				for _, val := range usuario.Cajas {
					cajasPersona = append(cajasPersona, val.Hex())
				}
				Send.Usuario.ECajasUsuario.Ihtml = template.HTML(EquipoCajaModel.CargaComboCajasMultiArrayObjID(cajasPersona))
				Send.Usuario.EMediosDeContactoUsuario.ECorreos.Principal = usuario.MediosDeContacto.Correos.Principal
				Send.Usuario.EMediosDeContactoUsuario.ECorreos.Ihtml = template.HTML(GeneraTemplateCorreosDetalle(usuario.MediosDeContacto.Correos.Correos, usuario.MediosDeContacto.Correos.Principal))
				Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Principal = usuario.MediosDeContacto.Telefonos.Principal
				Send.Usuario.EMediosDeContactoUsuario.ETelefonos.Ihtml = template.HTML(GeneraTemplateTelefonosDetalle(usuario.MediosDeContacto.Telefonos.Telefonos, usuario.MediosDeContacto.Telefonos.Principal))
				Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.Ihtml = template.HTML(GeneraTemplateOtrosDetalle(usuario.MediosDeContacto.Otros))
			} else {
				Send.SEstado = false
				Send.SMsj = "El Usuario no se ha encontrado"
			}

		} else {
			Send.SEstado = false
			Send.SMsj = "El Usuario no se ha encontrado"
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "Error en la referencia al Usuario"
	}

	ctx.Render("UsuarioDetalle.html", Send)

}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send UsuarioModel.SUsuario
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

	ctx.Render("UsuarioDetalle.html", Send)
}

//EditaPropioGet renderea a la edicion del propio usuario
func EditaPropioGet(ctx *iris.Context) {
	fmt.Println("ModificaCuenta GET")
	ctx.Render("ModificaCuenta.html", nil)
}

//####################< RUTINAS ADICIONALES >##########################

//Perfil Renderea a la Vista del Perfil de Usuario
func Perfil(ctx *iris.Context) {
	fmt.Println("Perfil GET")
	var Send UsuarioModel.SUsuario
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

	ctx.Render("Perfil.html", Send)

}

//AdminUsers Renderea la vista de Administracion de Usuarios
func AdminUsers(ctx *iris.Context) {
	fmt.Println("Administrar Usuarios GET")
	var Send UsuarioModel.SUsuario
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

	ctx.Render("AdministrarUsuarios.html", Send)

}

func NotificacionesDeUsuario(ctx *iris.Context) {
	fmt.Println("Administrar Usuarios GET")
	var Send UsuarioModel.SUsuario
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

	ctx.Render("Notificaciones.html", Send)

}

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send UsuarioModel.SUsuario

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

		Cabecera, Cuerpo := IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrToMongo))
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
	var Send UsuarioModel.SUsuario
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Usuario.ENombreUsuario.Nombre = cadenaBusqueda

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

			Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = IndexUsuarios.GeneraTemplatesBusqueda(UsuarioModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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

//GeneraTemplateCorreos Genera la tabla dinamica de los correos ingresados
func GeneraTemplateCorreos(correos []string, principal string) string {
	html := ``
	for _, val := range correos {
		if val == principal {
			html += `<tr>
					<td><input type="radio" name="CorreosPrincipal" value="` + val + `" checked></td>
					<td><input type="text" class="form-control" name="Correos" value="` + val + `" readonly></td>
					<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		} else {
			html += `<tr>
					<td><input type="radio" name="CorreosPrincipal" value="` + val + `"></td>
					<td><input type="text" class="form-control" name="Correos" value="` + val + `" readonly></td>
					<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		}
	}
	return html
}

//GeneraTemplateCorreosDetalle Genera la tabla dinamica de los correos ingresados
func GeneraTemplateCorreosDetalle(correos []string, principal string) string {
	html := ``
	for _, val := range correos {
		if val == principal {
			html += `<tr>
					<td><input type="radio" name="CorreosPrincipal" value="` + val + `" checked disabled></td>
					<td><input type="text" class="form-control" name="Correos" value="` + val + `" readonly></td>
				</tr>`
		} else {
			html += `<tr>
					<td><input type="radio" name="CorreosPrincipal" value="` + val + `" disabled></td>
					<td><input type="text" class="form-control" name="Correos" value="` + val + `" readonly></td>
				</tr>`
		}
	}
	return html
}

//GeneraTemplateTelefonos Genera la tabla dinamica de los correos ingresados
func GeneraTemplateTelefonos(telefonos []string, principal string) string {
	html := ``
	for _, val := range telefonos {
		if val == principal {
			html += `<tr>
				<td><input type="radio" name="TelefonosPrincipal" value="` + val + `" checked></td>
				<td><input type="text" class="form-control" name="Telefonos" value="` + val + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		} else {
			html += `<tr>
				<td><input type="radio" name="TelefonosPrincipal" value="` + val + `"></td>
				<td><input type="text" class="form-control" name="Telefonos" value="` + val + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		}
	}
	return html
}

//GeneraTemplateTelefonosDetalle Genera la tabla dinamica de los correos ingresados
func GeneraTemplateTelefonosDetalle(telefonos []string, principal string) string {
	html := ``
	for _, val := range telefonos {
		if val == principal {
			html += `<tr>
				<td><input type="radio" name="TelefonosPrincipal" value="` + val + `" checked disabled></td>
				<td><input type="text" class="form-control" name="Telefonos" value="` + val + `" readonly></td>
				</tr>`
		} else {
			html += `<tr>
				<td><input type="radio" name="TelefonosPrincipal" value="` + val + `" disabled></td>
				<td><input type="text" class="form-control" name="Telefonos" value="` + val + `" readonly></td>
				</tr>`
		}
	}
	return html
}

//GeneraTemplateOtros Genera la tabla dinamica de los correos ingresados
func GeneraTemplateOtros(otros []string) string {
	html := ``
	for _, val := range otros {
		html += `<tr>
				<td><input type="text" class="form-control" name="Otros" value="` + val + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`

	}
	return html
}

//GeneraTemplateOtrosDetalle Genera la tabla dinamica de los correos ingresados
func GeneraTemplateOtrosDetalle(otros []string) string {
	html := ``
	for _, val := range otros {
		html += `<tr>
				<td><input type="text" class="form-control" name="Otros" value="` + val + `" readonly></td>
				</tr>`

	}
	return html
}
