package ClienteControler

import (
	"encoding/json"
	"html/template"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modelos/GrupoPersonaModel"
	"../../Modelos/PersonaModel"
	"../../Modulos/Session"

	"fmt"

	"../../Modelos/CatalogoModel"
	"../../Modelos/ClienteModel"
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
var limitePorPagina = 5

//IDElastic id obtenido de Elastic
var IDElastic bson.ObjectId
var arrIDMgo []bson.ObjectId
var arrIDElastic []bson.ObjectId
var arrToMongo []bson.ObjectId

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Cliente
func IndexGet(ctx *iris.Context) {

	var Send ClienteModel.SCliente

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
	numeroRegistros = ClienteModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Clientes := ClienteModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Clientes {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(Clientes[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(Clientes[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("ClienteIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Cliente
func IndexPost(ctx *iris.Context) {

	var Send ClienteModel.SCliente

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
	//Send.Cliente.EVARIABLECliente.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := ClienteModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("ClienteIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Cliente
func AltaGet(ctx *iris.Context) {

	var Send ClienteModel.SCliente
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
	Send.Cliente.ETipoCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(175, "5936efac8c649f1b8839e48d"))
	Send.Cliente.EIDPersonaCliente.ESexo.Ihtml = template.HTML(CargaCombos.CargaComboSexo(""))
	Send.Cliente.EDireccionesCliente.Direcciones.EPaisDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboxPaises("5931ab6b565925984547673a"))
	Send.Cliente.EDireccionesCliente.Direcciones.EEstadoDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboEstados(""))
	Send.Cliente.EDireccionesCliente.Direcciones.ETipoDireccion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(174, ""))
	Send.Cliente.EIDPersonaCliente.ETipoPersona.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(159, ""))
	// Send.Cliente.EIDPersonaCliente.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasMulti(""))
	Send.Cliente.EIDPersonaCliente.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(""))
	Send.Cliente.EAlmacenesCliente.Ihtml = template.HTML(AlmacenModel.CargaComboGrupoAlmacenesMulti(""))

	Send.Cliente.EDireccionesCliente.NumDirecciones = 0

	ctx.Render("ClienteAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Cliente
func AltaPost(ctx *iris.Context) {
	EstatusPeticion := false
	var Send ClienteModel.SCliente
	var PersonaMongo PersonaModel.PersonaMgo
	var ClienteMongo ClienteModel.ClienteMgo

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
	ClienteMongo.ID = bson.NewObjectId()

	Tipo := ctx.FormValue("Tipo")
	Send.Cliente.ETipoCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(175, Tipo))
	if bson.IsObjectIdHex(Tipo) {
		ClienteMongo.TipoCliente = bson.ObjectIdHex(Tipo)
	}
	Send.Cliente.ETipoCliente.TipoCliente = bson.ObjectIdHex(Tipo)
	if Tipo == "5936efac8c649f1b8839e48d" {
		Sexo := ctx.FormValue("Sexo")
		PersonaMongo.Sexo = Sexo
		Send.Cliente.EIDPersonaCliente.ESexo.Ihtml = template.HTML(CargaCombos.CargaComboSexo(Sexo))

		FechaNacimiento := ctx.FormValue("FechaNacimiento")
		if FechaNacimiento != "" {
			t, _ := time.Parse("02-01-2006", FechaNacimiento)
			PersonaMongo.FechaNacimiento = t
			Send.Cliente.EIDPersonaCliente.EFechaNacimiento.FechaNacimiento = FechaNacimiento
		}
	} else {

		Send.Cliente.EIDPersonaCliente.ESexo.Ihtml = template.HTML(CargaCombos.CargaComboSexo(""))
	}

	Nombre := ctx.FormValue("Nombre")
	Send.Cliente.EIDPersonaCliente.ENombrePersona.Nombre = Nombre
	PersonaMongo.Nombre = Nombre
	if Nombre == "" {
		EstatusPeticion = true
		Send.Cliente.EIDPersonaCliente.ENombrePersona.IEstatus = true
		Send.Cliente.EIDPersonaCliente.ENombrePersona.IMsj = "El campo Nombre no debe estar vacio"
	}

	var arrTipo []bson.ObjectId
	arrTipo = append(arrTipo, bson.ObjectIdHex("58e56616e75770120c60bec4"))
	PersonaMongo.Tipo = arrTipo

	// Grupos := ctx.Request.Form["Grupos"]
	// Send.Cliente.EIDPersonaCliente.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(Grupos))
	// var arrGpo []bson.ObjectId
	// for _, val := range Grupos {
	// 	arrGpo = append(arrGpo, bson.ObjectIdHex(val))
	// }
	// PersonaMongo.Grupos = arrGpo

	Predecesor := ctx.FormValue("Predecesor")
	Send.Cliente.EIDPersonaCliente.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(Predecesor))
	if bson.IsObjectIdHex(Predecesor) {
		Send.Cliente.EIDPersonaCliente.EPredecesorPersona.Predecesor = bson.ObjectIdHex(Predecesor)
		PersonaMongo.Predecesor = bson.ObjectIdHex(Predecesor)
	}

	PersonaMongo.ID = bson.NewObjectId()
	ClienteMongo.IDPersona = PersonaMongo.ID

	RFC := ctx.FormValue("RFC")
	Send.Cliente.ERFCCliente.RFC = RFC
	ClienteMongo.RFC = RFC
	if Nombre == "" {
		EstatusPeticion = true
		Send.Cliente.ERFCCliente.IEstatus = true
		Send.Cliente.ERFCCliente.IMsj = "El campo RFC no debe estar vacio"
	}

	Almacenes := ctx.Request.Form["Almacenes"]
	Send.Cliente.EAlmacenesCliente.Ihtml = template.HTML(AlmacenModel.CargaComboAlamcenesArray(Almacenes))
	var arrAlm []bson.ObjectId
	for _, val := range Almacenes {
		arrAlm = append(arrAlm, bson.ObjectIdHex(val))
	}
	ClienteMongo.Almacenes = arrAlm

	Send.Cliente.EDireccionesCliente.Direcciones.EEstadoDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboEstados(""))
	Send.Cliente.EDireccionesCliente.Direcciones.EPaisDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboxPaises("5931ab6b565925984547673a"))

	sufijosDirCte := ctx.Request.Form["NumDirCliente"]
	var Direccion ClienteModel.DireccionMgo
	var arrDirecciones []ClienteModel.DireccionMgo
	fmt.Println(ctx.FormValues())
	var tr string //variable que contiene las filas de tablas de direciones
	numTr := 0    //variable para el numero de  direcciones del cliente
	for i, val := range sufijosDirCte {

		Direccion.Pais = "Mexico"
		Direccion.Estado = bson.ObjectIdHex(ctx.FormValue("Estador" + val))
		Direccion.Municipio = bson.ObjectIdHex(ctx.FormValue("Municipior" + val))
		Direccion.Colonia = bson.ObjectIdHex(ctx.FormValue("Coloniar" + val))
		Direccion.Calle = ctx.FormValue("Caller" + val)
		Direccion.NumExterior = ctx.FormValue("NumExteriorr" + val)
		Direccion.NumInterior = ctx.FormValue("NumInteriorr" + val)
		Direccion.CP = ctx.FormValue("cpr" + val)
		Direccion.TipoDireccion = bson.ObjectIdHex(ctx.FormValue("TipoDireccionr" + val))

		Idtr := ctx.FormValue("IDr" + val)
		if Idtr == "" {
			Direccion.ID = bson.NewObjectId()
		} else {
			Direccion.ID = bson.ObjectIdHex(Idtr)
		}

		tr += `<tr class='direccionCliente' data-id-direccion-cliente='` + Direccion.ID.Hex() + `' data-num-direccion-cliente="` + strconv.Itoa(i+1) + `" >
				<td><input type='hidden' name='IDr` + strconv.Itoa(i+1) + `' value='` + Direccion.ID.Hex() + `'>
				    <input type='hidden' name='Estador` + strconv.Itoa(i+1) + `' value='` + Direccion.Estado.Hex() + `'>` + CatalogoModel.GetNameEstado(Direccion.Estado) + `</td>
				<td><input type='hidden' name='Municipior` + strconv.Itoa(i+1) + `' value='` + Direccion.Municipio.Hex() + `'>` + CatalogoModel.GetNameMunicipio(Direccion.Municipio) + `</td>
				<td><input type='hidden' name='Coloniar` + strconv.Itoa(i+1) + `' value='` + Direccion.Colonia.Hex() + `'>` + CatalogoModel.GetNameColonia(Direccion.Colonia) + `</td>
				<td><input type='hidden' name='cpr` + strconv.Itoa(i+1) + `' value='` + Direccion.CP + `'>` + Direccion.CP + `</td>
				<td><input type='hidden' name='Caller` + strconv.Itoa(i+1) + `' value='` + Direccion.Calle + `'>` + Direccion.Calle + `</td>
				<td><input type='hidden' name='NumExteriorr` + strconv.Itoa(i+1) + `' value='` + Direccion.NumExterior + `'>` + Direccion.NumExterior + `</td>
				<td><input type='hidden' name='NumInteriorr` + strconv.Itoa(i+1) + `' value='` + Direccion.NumInterior + `'>` + Direccion.NumInterior + `</td>
				<td><input type='hidden' name='TipoDireccionr` + strconv.Itoa(i+1) + `' value='` + Direccion.TipoDireccion.Hex() + `'>` + CatalogoModel.RegresaNombreSubCatalogo(Direccion.TipoDireccion) + `</td>
				<td><button type="button" class="btn btn-danger deleteDirCP"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		numTr = i + 1
		arrDirecciones = append(arrDirecciones, Direccion)
	}

	ClienteMongo.Direcciones = arrDirecciones
	Send.Cliente.EDireccionesCliente.NumDirecciones = numTr
	Send.Cliente.EDireccionesCliente.Ihtml = template.HTML(tr)

	CorreoPrincipal := ctx.FormValue("CorreoPrincipal")
	ClienteMongo.MediosDeContacto.Correos.Principal = CorreoPrincipal
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.Principal = CorreoPrincipal
	if CorreoPrincipal == "" {
		EstatusPeticion = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.IEstatus = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.IMsj = "Debe seleccionar un correo como principal"
	}

	Correos := ctx.Request.Form["Correos"]

	ClienteMongo.MediosDeContacto.Correos.Correos = Correos
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.Ihtml = template.HTML(GeneraTemplateCorreos(Correos, CorreoPrincipal))
	if len(Correos) == 0 {
		EstatusPeticion = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.IEstatus = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.IMsj = "Debe agregar al menos un correo"
	}

	TelefonoPrincipal := ctx.FormValue("TelefonoPrincipal")
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.Principal = TelefonoPrincipal
	ClienteMongo.MediosDeContacto.Telefonos.Principal = TelefonoPrincipal
	if TelefonoPrincipal == "" {
		EstatusPeticion = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.IEstatus = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.IMsj = "Debe seleccionar un telefono como principal"
	}

	Telefonos := ctx.Request.Form["Telefonos"]
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.Ihtml = template.HTML(GeneraTemplateTelefonos(Telefonos, TelefonoPrincipal))
	ClienteMongo.MediosDeContacto.Telefonos.Telefonos = Telefonos
	if len(TelefonoPrincipal) == 0 {
		EstatusPeticion = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.IEstatus = true
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.IMsj = "Debe agregar al menos un telefono"
	}

	Otros := ctx.Request.Form["Otros"]
	//Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.OtrosMedios = Otros
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.EOtrosMediosDeContacto.Ihtml = template.HTML(GeneraTemplateOtros(Otros))
	ClienteMongo.MediosDeContacto.Otros = Otros

	totalPersonas := ctx.Request.Form["NumPer"]
	fmt.Println("totalPersonas:", totalPersonas)
	var DireccionContPer ClienteModel.DireccionMgo
	var PersonaContacto ClienteModel.PersonaContactoMgo
	var arrPersonaContacto []ClienteModel.PersonaContactoMgo
	var templatePerCont string
	for i, val := range totalPersonas {

		nombre := ctx.FormValue("Nombre" + val)
		PersonaContacto.Nombre = nombre
		var arrDireccionesPC []ClienteModel.DireccionMgo
		sufijosDirPerCont := ctx.Request.Form["NumDirPerCont"+val]

		for _, numDirPerCont := range sufijosDirPerCont {
			DireccionContPer.Pais = "Mexico"
			DireccionContPer.Estado = bson.ObjectIdHex(ctx.FormValue("EstadoPC" + numDirPerCont))
			DireccionContPer.Municipio = bson.ObjectIdHex(ctx.FormValue("MunicipioPC" + numDirPerCont))
			DireccionContPer.Colonia = bson.ObjectIdHex(ctx.FormValue("ColoniaPC" + numDirPerCont))
			DireccionContPer.Calle = ctx.FormValue("CallePC" + numDirPerCont)
			DireccionContPer.NumExterior = ctx.FormValue("NumExteriorPC" + numDirPerCont)
			DireccionContPer.NumInterior = ctx.FormValue("NumInteriorPC" + numDirPerCont)
			DireccionContPer.CP = ctx.FormValue("cpPC" + numDirPerCont)
			DireccionContPer.TipoDireccion = bson.ObjectIdHex(ctx.FormValue("TipoDireccionPC" + numDirPerCont))

			Idtr := ctx.FormValue("IDr" + numDirPerCont)
			if Idtr == "" {
				DireccionContPer.ID = bson.NewObjectId()
			} else {
				DireccionContPer.ID = bson.ObjectIdHex(Idtr)
			}

			arrDireccionesPC = append(arrDireccionesPC, DireccionContPer)
		}
		PersonaContacto.Direcciones = arrDireccionesPC
		CorreoPrincipalCp := ctx.FormValue("CorreoPrincipal" + val)
		PersonaContacto.MediosDeContacto.Correos.Principal = CorreoPrincipalCp

		Correos := ctx.Request.Form["Correos"+val]
		PersonaContacto.MediosDeContacto.Correos.Correos = Correos

		TelefonoPrincipal := ctx.FormValue("TelefonoPrincipal" + val)
		PersonaContacto.MediosDeContacto.Telefonos.Principal = TelefonoPrincipal

		Telefonos := ctx.Request.Form["Telefonos"+val]
		PersonaContacto.MediosDeContacto.Telefonos.Telefonos = Telefonos

		Otros := ctx.Request.Form["Otros"+val]
		PersonaContacto.MediosDeContacto.Otros = Otros

		arrPersonaContacto = append(arrPersonaContacto, PersonaContacto)

		templatePerCont += GeneraGeneraTemplatePersonasContacto(strconv.Itoa(i+1), PersonaContacto)

	}

	ClienteMongo.PersonasContacto = arrPersonaContacto
	Send.Cliente.EPersonasContactoCliente.NumPerCont = len(arrPersonaContacto)
	Send.Cliente.EPersonasContactoCliente.Ihtml = template.HTML(templatePerCont)

	ClienteMongo.Estatus = bson.ObjectIdHex("58e57696e75770120c60bf06")
	Send.Cliente.EEstatusCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(137, "58e57696e75770120c60bf06"))

	if EstatusPeticion {
		Send.SMsj = "La validacion Indica que algo ha salido mal"
		ctx.Render("ClienteAlta.html", Send)

	} else {

		if PersonaMongo.InsertaMgo() {
			if ClienteMongo.InsertaMgo() {
				if PersonaMongo.InsertaElastic() {
					if ClienteMongo.InsertaElastic() {
						Send.SEstado = true
						Send.SMsj = "Se ha realizado una inserción exitosa"
						ctx.Redirect("/Clientes/detalle/"+ClienteMongo.ID.Hex(), 301)
					} else {
						Send.SEstado = false
						Send.SMsj = "Ocurrió un error al insertar Cliente en elasticSearch"
					}
				} else {
					Send.SEstado = false
					Send.SMsj = "Ocurrió un error al insertar Persona en elasticSearch"
				}
			} else {
				eliminado := ClienteMongo.EliminaByIDMgo()
				Send.SEstado = false
				Send.SMsj = "Ocurrió un error al insertar Cliente en MongoDb"
				if !eliminado {
					Send.SMsj += "No se ha podido eliminar el objeto"
				}
			}
		} else {
			Send.SEstado = false
			Send.SMsj = "Ocurrió un error al insertar Persona en MongoDb"
		}
	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Cliente
func EditaGet(ctx *iris.Context) {

	var Send ClienteModel.SCliente

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
	id := ctx.Param("ID")
	Cliente := ClienteModel.GetOne(bson.ObjectIdHex(id))
	Persona := PersonaModel.GetOne(Cliente.IDPersona)

	Send.Cliente.ETipoCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(175, Cliente.TipoCliente.Hex()))
	Send.Cliente.ETipoCliente.TipoCliente = Cliente.TipoCliente
	Send.Cliente.EIDPersonaCliente.ESexo.Ihtml = template.HTML(CargaCombos.CargaComboSexo(Persona.Sexo))
	if Cliente.TipoCliente.Hex() == "5936efac8c649f1b8839e48d" {
		if Persona.FechaNacimiento.Format("02-01-2006") != "01-01-0001" {
			Send.Cliente.EIDPersonaCliente.EFechaNacimiento.FechaNacimiento = Persona.FechaNacimiento.Format("02-01-2006")
		}
	}

	Send.Cliente.ID = Cliente.ID
	Send.Cliente.EIDPersonaCliente.ID = Persona.ID
	Send.Cliente.EIDPersonaCliente.ENombrePersona.Nombre = Persona.Nombre
	// Send.Cliente.EIDPersonaCliente.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(ConvertArrayObjectIDToArrayString(Persona.Grupos)))
	Send.Cliente.EIDPersonaCliente.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(Persona.Predecesor.Hex()))
	Send.Cliente.EAlmacenesCliente.Ihtml = template.HTML(AlmacenModel.CargaComboAlamcenesArray(ConvertArrayObjectIDToArrayString(Cliente.Almacenes)))
	Send.Cliente.EEstatusCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(137, "58e57696e75770120c60bf06"))
	Send.Cliente.EDireccionesCliente.Direcciones.EEstadoDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboEstados(""))
	Send.Cliente.EDireccionesCliente.Direcciones.EPaisDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboxPaises("5931ab6b565925984547673a"))
	Send.Cliente.EDireccionesCliente.Direcciones.ETipoDireccion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(174, ""))
	Send.Cliente.ERFCCliente.RFC = Cliente.RFC

	var tr string //variable que contiene las filas de tablas de direciones
	numTr := 0    //variable para el numero de  direcciones del cliente
	for i, Val := range Cliente.Direcciones {

		tr += `<tr class='direccionCliente' data-id-direccion-cliente='` + Val.ID.Hex() + `' data-num-direccion-cliente="` + strconv.Itoa(i+1) + `" >
				<td><input type='hidden' name='IDr` + strconv.Itoa(i+1) + `' value='` + Val.ID.Hex() + `'>
				    <input type='hidden' name='Estador` + strconv.Itoa(i+1) + `' value='` + Val.Estado.Hex() + `'>` + CatalogoModel.GetNameEstado(Val.Estado) + `</td>
				<td><input type='hidden' name='Municipior` + strconv.Itoa(i+1) + `' value='` + Val.Municipio.Hex() + `'>` + CatalogoModel.GetNameMunicipio(Val.Municipio) + `</td>
				<td><input type='hidden' name='Coloniar` + strconv.Itoa(i+1) + `' value='` + Val.Colonia.Hex() + `'>` + CatalogoModel.GetNameColonia(Val.Colonia) + `</td>
				<td><input type='hidden' name='cpr` + strconv.Itoa(i+1) + `' value='` + Val.CP + `'>` + Val.CP + `</td>
				<td><input type='hidden' name='Caller` + strconv.Itoa(i+1) + `' value='` + Val.Calle + `'>` + Val.Calle + `</td>
				<td><input type='hidden' name='NumExteriorr` + strconv.Itoa(i+1) + `' value='` + Val.NumExterior + `'>` + Val.NumExterior + `</td>
				<td><input type='hidden' name='NumInteriorr` + strconv.Itoa(i+1) + `' value='` + Val.NumInterior + `'>` + Val.NumInterior + `</td>
				<td><input type='hidden' name='TipoDireccionr` + strconv.Itoa(i+1) + `' value='` + Val.TipoDireccion.Hex() + `'>` + CatalogoModel.RegresaNombreSubCatalogo(Val.TipoDireccion) + `</td>
				<td><button type="button" class="btn btn-danger deleteDirCP"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		numTr = i + 1

	}
	Send.Cliente.EDireccionesCliente.NumDirecciones = numTr
	Send.Cliente.EDireccionesCliente.Ihtml = template.HTML(tr)

	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.Principal = Cliente.MediosDeContacto.Correos.Principal
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.Ihtml = template.HTML(GeneraTemplateCorreos(Cliente.MediosDeContacto.Correos.Correos, Cliente.MediosDeContacto.Correos.Principal))

	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.Principal = Cliente.MediosDeContacto.Telefonos.Principal
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.Ihtml = template.HTML(GeneraTemplateTelefonos(Cliente.MediosDeContacto.Telefonos.Telefonos, Cliente.MediosDeContacto.Telefonos.Principal))

	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.EOtrosMediosDeContacto.Ihtml = template.HTML(GeneraTemplateOtros(Cliente.MediosDeContacto.Otros))

	var templatePerCont string
	for i, Val := range Cliente.PersonasContacto {
		templatePerCont += GeneraGeneraTemplatePersonasContacto(strconv.Itoa(i+1), Val)
	}

	Send.Cliente.EPersonasContactoCliente.NumPerCont = len(Cliente.PersonasContacto)
	Send.Cliente.EPersonasContactoCliente.Ihtml = template.HTML(templatePerCont)
	ctx.Render("ClienteEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Cliente
func EditaPost(ctx *iris.Context) {

	var Send ClienteModel.SCliente
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
	id := ctx.Param("ID")
	if bson.IsObjectIdHex(id) {

		Cliente := ClienteModel.GetOne(bson.ObjectIdHex(id))
		Persona := PersonaModel.GetOne(Cliente.IDPersona)
		Nombre := ctx.FormValue("Nombre")
		Send.Cliente.ID = Cliente.ID
		Send.Cliente.EIDPersonaCliente.ENombrePersona.Nombre = Nombre
		Persona.Nombre = Nombre
		if Nombre == "" {
			EstatusPeticion = true
			Send.Cliente.EIDPersonaCliente.ENombrePersona.IEstatus = true
			Send.Cliente.EIDPersonaCliente.ENombrePersona.IMsj = "El campo Nombre no debe estar vacio"
		}

		Tipo := ctx.FormValue("Tipo")
		Send.Cliente.ETipoCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(175, Tipo))
		Send.Cliente.ETipoCliente.TipoCliente = bson.ObjectIdHex(Tipo)
		if bson.IsObjectIdHex(Tipo) {
			Cliente.TipoCliente = bson.ObjectIdHex(Tipo)
		}

		if Tipo == "5936efac8c649f1b8839e48d" {

			Sexo := ctx.FormValue("Sexo")
			Persona.Sexo = Sexo
			Send.Cliente.EIDPersonaCliente.ESexo.Ihtml = template.HTML(CargaCombos.CargaComboSexo(Sexo))

			FechaNacimiento := ctx.FormValue("FechaNacimiento")
			if FechaNacimiento != "" {
				t, _ := time.Parse("02-01-2006", FechaNacimiento)
				Persona.FechaNacimiento = t
				Send.Cliente.EIDPersonaCliente.EFechaNacimiento.FechaNacimiento = FechaNacimiento
				fmt.Println("hola  Fecha: ", Persona.FechaNacimiento)
			}
		} else {

			Send.Cliente.EIDPersonaCliente.ESexo.Ihtml = template.HTML(CargaCombos.CargaComboSexo(""))
		}

		var arrTipo []bson.ObjectId
		arrTipo = append(arrTipo, bson.ObjectIdHex("58e56616e75770120c60bec4"))
		Persona.Tipo = arrTipo

		// Grupos := ctx.Request.Form["Grupos"]
		// Send.Cliente.EIDPersonaCliente.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(Grupos))
		// var arrGpo []bson.ObjectId
		// for _, val := range Grupos {
		// 	arrGpo = append(arrGpo, bson.ObjectIdHex(val))
		// }
		// Persona.Grupos = arrGpo

		Predecesor := ctx.FormValue("Predecesor")
		Send.Cliente.EIDPersonaCliente.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(Predecesor))
		if !bson.IsObjectIdHex(Predecesor) {
			Predecesor = ""
			//Send.Cliente.EIDPersonaCliente.EPredecesorPersona.Predecesor = bson.ObjectIdHex(Predecesor)
			//Persona.Predecesor = bson.ObjectIdHex(Predecesor)
		}

		RFC := ctx.FormValue("RFC")
		Send.Cliente.ERFCCliente.RFC = RFC
		Cliente.RFC = RFC
		if RFC == "" {
			EstatusPeticion = true
			Send.Cliente.ERFCCliente.IEstatus = true
			Send.Cliente.ERFCCliente.IMsj = "El campo RFC no debe estar vacio"
		}

		Almacenes := ctx.Request.Form["Almacenes"]
		Send.Cliente.EAlmacenesCliente.Ihtml = template.HTML(AlmacenModel.CargaComboAlamcenesArray(Almacenes))
		var arrAlm []bson.ObjectId
		for _, val := range Almacenes {
			arrAlm = append(arrAlm, bson.ObjectIdHex(val))
		}
		Cliente.Almacenes = arrAlm

		Estatus := ctx.FormValue("Estatus")
		Send.Cliente.EEstatusCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(137, Estatus))
		if bson.IsObjectIdHex(Estatus) {
			//Send.Cliente.EEstatusCliente.Estatus = bson.ObjectIdHex(Estatus)
			Cliente.Estatus = bson.ObjectIdHex(Estatus)
		} else {
			EstatusPeticion = true
			Send.Cliente.EEstatusCliente.IEstatus = true
			Send.Cliente.EEstatusCliente.IMsj = "Se debe seleccionar el estaus del Cliente"
		}

		Send.Cliente.EDireccionesCliente.Direcciones.EEstadoDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboEstados(""))
		Send.Cliente.EDireccionesCliente.Direcciones.EPaisDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboxPaises("5931ab6b565925984547673a"))
		Send.Cliente.EDireccionesCliente.Direcciones.ETipoDireccion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(174, ""))

		sufijosDirCte := ctx.Request.Form["NumDirCliente"]
		var Direccion ClienteModel.DireccionMgo
		var arrDirecciones []ClienteModel.DireccionMgo
		var tr string //variable que contiene las filas de tablas de direciones
		numTr := 0    //variable para el numero de  direcciones del cliente

		for i, val := range sufijosDirCte {
			Direccion.Pais = "Mexico"
			Direccion.Estado = bson.ObjectIdHex(ctx.FormValue("Estador" + val))
			Direccion.Municipio = bson.ObjectIdHex(ctx.FormValue("Municipior" + val))
			Direccion.Colonia = bson.ObjectIdHex(ctx.FormValue("Coloniar" + val))
			Direccion.Calle = ctx.FormValue("Caller" + val)
			Direccion.NumExterior = ctx.FormValue("NumExteriorr" + val)
			Direccion.NumInterior = ctx.FormValue("NumInteriorr" + val)
			Direccion.CP = ctx.FormValue("cpr" + val)
			Direccion.TipoDireccion = bson.ObjectIdHex(ctx.FormValue("TipoDireccionr" + val))

			Idtr := ctx.FormValue("IDr" + val)
			if Idtr == "" {
				Direccion.ID = bson.NewObjectId()
			} else {
				Direccion.ID = bson.ObjectIdHex(Idtr)
			}

			tr += `<tr class='direccionCliente' data-id-direccion-cliente='` + Direccion.ID.Hex() + `' data-num-direccion-cliente="` + strconv.Itoa(i+1) + `" >
				<td><input type='hidden' name='IDr` + strconv.Itoa(i+1) + `' value='` + Direccion.ID.Hex() + `'>
				    <input type='hidden' name='Estador` + strconv.Itoa(i+1) + `' value='` + Direccion.Estado.Hex() + `'>` + CatalogoModel.GetNameEstado(Direccion.Estado) + `</td>
				<td><input type='hidden' name='Municipior` + strconv.Itoa(i+1) + `' value='` + Direccion.Municipio.Hex() + `'>` + CatalogoModel.GetNameMunicipio(Direccion.Municipio) + `</td>
				<td><input type='hidden' name='Coloniar` + strconv.Itoa(i+1) + `' value='` + Direccion.Colonia.Hex() + `'>` + CatalogoModel.GetNameColonia(Direccion.Colonia) + `</td>
				<td><input type='hidden' name='cpr` + strconv.Itoa(i+1) + `' value='` + Direccion.CP + `'>` + Direccion.CP + `</td>
				<td><input type='hidden' name='Caller` + strconv.Itoa(i+1) + `' value='` + Direccion.Calle + `'>` + Direccion.Calle + `</td>
				<td><input type='hidden' name='NumExteriorr` + strconv.Itoa(i+1) + `' value='` + Direccion.NumExterior + `'>` + Direccion.NumExterior + `</td>
				<td><input type='hidden' name='NumInteriorr` + strconv.Itoa(i+1) + `' value='` + Direccion.NumInterior + `'>` + Direccion.NumInterior + `</td>
				<td><input type='hidden' name='TipoDireccionr` + strconv.Itoa(i+1) + `' value='` + Direccion.TipoDireccion.Hex() + `'>` + CatalogoModel.RegresaNombreSubCatalogo(Direccion.TipoDireccion) + `</td>
				<td><button type="button" class="btn btn-danger deleteDirCP"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
			numTr = i + 1
			arrDirecciones = append(arrDirecciones, Direccion)
		}

		Cliente.Direcciones = arrDirecciones
		Send.Cliente.EDireccionesCliente.NumDirecciones = numTr
		Send.Cliente.EDireccionesCliente.Ihtml = template.HTML(tr)

		CorreoPrincipal := ctx.FormValue("CorreoPrincipal")
		Cliente.MediosDeContacto.Correos.Principal = CorreoPrincipal
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.Principal = CorreoPrincipal
		if CorreoPrincipal == "" {
			EstatusPeticion = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.IEstatus = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.IMsj = "Debe seleccionar un correo como principal"
		}

		Correos := ctx.Request.Form["Correos"]
		Cliente.MediosDeContacto.Correos.Correos = Correos
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.Ihtml = template.HTML(GeneraTemplateCorreos(Correos, CorreoPrincipal))
		if len(Correos) == 0 {
			EstatusPeticion = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.IEstatus = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.IMsj = "Debe agregar al menos un correo"
		}

		TelefonoPrincipal := ctx.FormValue("TelefonoPrincipal")
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.Principal = TelefonoPrincipal
		Cliente.MediosDeContacto.Telefonos.Principal = TelefonoPrincipal
		if TelefonoPrincipal == "" {
			EstatusPeticion = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.IEstatus = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.IMsj = "Debe seleccionar un telefono como principal"
		}

		Telefonos := ctx.Request.Form["Telefonos"]
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.Ihtml = template.HTML(GeneraTemplateTelefonos(Telefonos, TelefonoPrincipal))
		Cliente.MediosDeContacto.Telefonos.Telefonos = Telefonos
		if len(TelefonoPrincipal) == 0 {
			EstatusPeticion = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.IEstatus = true
			Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.IMsj = "Debe agregar al menos un telefono"
		}

		Otros := ctx.Request.Form["Otros"]
		//Send.Usuario.EMediosDeContactoUsuario.EOtrosMedios.OtrosMedios = Otros
		Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.EOtrosMediosDeContacto.Ihtml = template.HTML(GeneraTemplateOtros(Otros))
		Cliente.MediosDeContacto.Otros = Otros

		totalPersonas := ctx.Request.Form["NumPer"]
		var DireccionContPer ClienteModel.DireccionMgo
		var PersonaContacto ClienteModel.PersonaContactoMgo
		var arrPersonaContacto []ClienteModel.PersonaContactoMgo
		var templatePerCont string
		for i, val := range totalPersonas {

			nombre := ctx.FormValue("Nombre" + val)
			PersonaContacto.Nombre = nombre
			var arrDireccionesPC []ClienteModel.DireccionMgo
			sufijosDirPerCont := ctx.Request.Form["NumDirPerCont"+val]

			for _, numDirPerCont := range sufijosDirPerCont {
				DireccionContPer.Pais = "Mexico"
				DireccionContPer.Estado = bson.ObjectIdHex(ctx.FormValue("EstadoPC" + numDirPerCont))
				DireccionContPer.Municipio = bson.ObjectIdHex(ctx.FormValue("MunicipioPC" + numDirPerCont))
				DireccionContPer.Colonia = bson.ObjectIdHex(ctx.FormValue("ColoniaPC" + numDirPerCont))
				DireccionContPer.Calle = ctx.FormValue("CallePC" + numDirPerCont)
				DireccionContPer.NumExterior = ctx.FormValue("NumExteriorPC" + numDirPerCont)
				DireccionContPer.NumInterior = ctx.FormValue("NumInteriorPC" + numDirPerCont)
				DireccionContPer.CP = ctx.FormValue("cpPC" + numDirPerCont)
				DireccionContPer.TipoDireccion = bson.ObjectIdHex(ctx.FormValue("TipoDireccionPC" + numDirPerCont))

				Idtr := ctx.FormValue("IDr" + numDirPerCont)
				if Idtr == "" {
					DireccionContPer.ID = bson.NewObjectId()
				} else {
					DireccionContPer.ID = bson.ObjectIdHex(Idtr)
				}

				arrDireccionesPC = append(arrDireccionesPC, DireccionContPer)
			}
			PersonaContacto.Direcciones = arrDireccionesPC

			CorreoPrincipalCp := ctx.FormValue("CorreoPrincipal" + val)
			PersonaContacto.MediosDeContacto.Correos.Principal = CorreoPrincipalCp

			Correos := ctx.Request.Form["Correos"+val]
			PersonaContacto.MediosDeContacto.Correos.Correos = Correos

			TelefonoPrincipal := ctx.FormValue("TelefonoPrincipal" + val)
			PersonaContacto.MediosDeContacto.Telefonos.Principal = TelefonoPrincipal

			Telefonos := ctx.Request.Form["Telefonos"+val]
			PersonaContacto.MediosDeContacto.Telefonos.Telefonos = Telefonos

			Otros := ctx.Request.Form["Otros"+val]
			PersonaContacto.MediosDeContacto.Otros = Otros

			arrPersonaContacto = append(arrPersonaContacto, PersonaContacto)

			templatePerCont += GeneraGeneraTemplatePersonasContacto(strconv.Itoa(i+1), PersonaContacto)

		}

		Send.Cliente.EPersonasContactoCliente.NumPerCont = len(arrPersonaContacto)
		Send.Cliente.EPersonasContactoCliente.Ihtml = template.HTML(templatePerCont)

		Send.Cliente.EIDPersonaCliente.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(ConvertArrayObjectIDToArrayString(Persona.Grupos)))

		ActualizadoPersona := Persona.ActualizaMgo([]string{"Nombre", "Predecesor", "Sexo", "FechaNacimiento"}, []interface{}{Persona.Nombre, Predecesor, Persona.Sexo, Persona.FechaNacimiento})
		ActualizadoCliente := Cliente.ActualizaMgo([]string{"RFC", "Direcciones", "MediosDeContacto", "PersonasContacto", "Almacenes", "Estatus", "TipoCliente"}, []interface{}{Cliente.RFC, Cliente.Direcciones, Cliente.MediosDeContacto, arrPersonaContacto, Cliente.Almacenes, Cliente.Estatus, Cliente.TipoCliente})
		if EstatusPeticion {
			Send.SEstado = true
			Send.SMsj = "Cliente Editado con exito"
			ctx.Render("ClienteEdita.html", Send)
			// ctx.Redirect("ClienteEdita.html/"+id, 301)
		} else {

			if ActualizadoPersona && ActualizadoCliente {
				ctx.Render("ClienteDetalle.html", Send)
				// ctx.Redirect("ClienteDetalle.html/"+id, 301)
			} else {
				Send.SEstado = false
				Send.SMsj = "Error al editar Al Cliente"
				// ctx.Render("ClienteDetalle.html", Send)
				ctx.Render("ClienteEdita.html", Send)
			}
		}

	} else {
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar El Almacen, intente de nuevo."
		ctx.Render("AlmacenIndex.html", Send)
	}
	//ctx.Render("ClienteEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send ClienteModel.SCliente

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
	id := ctx.Param("ID")
	Cliente := ClienteModel.GetOne(bson.ObjectIdHex(id))
	Persona := PersonaModel.GetOne(Cliente.IDPersona)

	Send.Cliente.ETipoCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(175, Cliente.TipoCliente.Hex()))
	Send.Cliente.ETipoCliente.TipoCliente = Cliente.TipoCliente
	Send.Cliente.EIDPersonaCliente.ESexo.Ihtml = template.HTML(CargaCombos.CargaComboSexo(Persona.Sexo))
	if Cliente.TipoCliente.Hex() == "5936efac8c649f1b8839e48d" {
		if Persona.FechaNacimiento.Format("02-01-2006") != "01-01-0001" {
			Send.Cliente.EIDPersonaCliente.EFechaNacimiento.FechaNacimiento = Persona.FechaNacimiento.Format("02-01-2006")
		}
	}
	Send.Cliente.ID = Cliente.ID
	Send.Cliente.EIDPersonaCliente.ID = Persona.ID
	Send.Cliente.EIDPersonaCliente.ENombrePersona.Nombre = Persona.Nombre
	Send.Cliente.EIDPersonaCliente.EGruposPersona.Ihtml = template.HTML(GrupoPersonaModel.CargaComboGrupoPersonasArray(ConvertArrayObjectIDToArrayString(Persona.Grupos)))
	Send.Cliente.EIDPersonaCliente.EPredecesorPersona.Ihtml = template.HTML(PersonaModel.ConstruirComboUsuarioPredecesores(Persona.Predecesor.Hex()))
	Send.Cliente.EAlmacenesCliente.Ihtml = template.HTML(AlmacenModel.CargaComboAlamcenesArray(ConvertArrayObjectIDToArrayString(Cliente.Almacenes)))
	Send.Cliente.EEstatusCliente.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(137, "58e57696e75770120c60bf06"))
	Send.Cliente.EDireccionesCliente.Direcciones.EEstadoDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboEstados(""))
	Send.Cliente.EDireccionesCliente.Direcciones.EPaisDirecciones.Ihtml = template.HTML(CargaCombos.CargaComboxPaises("5931ab6b565925984547673a"))
	Send.Cliente.EDireccionesCliente.Direcciones.ETipoDireccion.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(174, ""))
	Send.Cliente.ERFCCliente.RFC = Cliente.RFC

	var tr string //variable que contiene las filas de tablas de direciones
	numTr := 0    //variable para el numero de  direcciones del cliente
	for i, Val := range Cliente.Direcciones {

		tr += `<tr class='direccionCliente' data-id-direccion-cliente='` + Val.ID.Hex() + `' data-num-direccion-cliente="` + strconv.Itoa(i+1) + `" >
				<td><input type='hidden' name='IDr` + strconv.Itoa(i+1) + `' value='` + Val.ID.Hex() + `'>
				    <input type='hidden' name='Estador` + strconv.Itoa(i+1) + `' value='` + Val.Estado.Hex() + `'>` + CatalogoModel.GetNameEstado(Val.Estado) + `</td>
				<td><input type='hidden' name='Municipior` + strconv.Itoa(i+1) + `' value='` + Val.Municipio.Hex() + `'>` + CatalogoModel.GetNameMunicipio(Val.Municipio) + `</td>
				<td><input type='hidden' name='Coloniar` + strconv.Itoa(i+1) + `' value='` + Val.Colonia.Hex() + `'>` + CatalogoModel.GetNameColonia(Val.Colonia) + `</td>
				<td><input type='hidden' name='cpr` + strconv.Itoa(i+1) + `' value='` + Val.CP + `'>` + Val.CP + `</td>
				<td><input type='hidden' name='Caller` + strconv.Itoa(i+1) + `' value='` + Val.Calle + `'>` + Val.Calle + `</td>
				<td><input type='hidden' name='NumExteriorr` + strconv.Itoa(i+1) + `' value='` + Val.NumExterior + `'>` + Val.NumExterior + `</td>
				<td><input type='hidden' name='NumInteriorr` + strconv.Itoa(i+1) + `' value='` + Val.NumInterior + `'>` + Val.NumInterior + `</td>
				<td><input type='hidden' name='TipoDireccionr` + strconv.Itoa(i+1) + `' value='` + Val.TipoDireccion.Hex() + `'>` + CatalogoModel.RegresaNombreSubCatalogo(Val.TipoDireccion) + `</td>
				<td><button type="button" class="btn btn-danger deleteDirCP"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		numTr = i + 1

	}
	Send.Cliente.EDireccionesCliente.NumDirecciones = numTr
	Send.Cliente.EDireccionesCliente.Ihtml = template.HTML(tr)

	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.EPrincipalCorreos.Principal = Cliente.MediosDeContacto.Correos.Principal
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ECorreosMediosDeContacto.Correos.ECorreosCorreos.Ihtml = template.HTML(GeneraTemplateCorreos(Cliente.MediosDeContacto.Correos.Correos, Cliente.MediosDeContacto.Correos.Principal))

	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.EPrincipalTelefonos.Principal = Cliente.MediosDeContacto.Telefonos.Principal
	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.ETelefonosMediosDeContacto.Telefonos.ETelefonosTelefonos.Ihtml = template.HTML(GeneraTemplateTelefonos(Cliente.MediosDeContacto.Telefonos.Telefonos, Cliente.MediosDeContacto.Telefonos.Principal))

	Send.Cliente.EMediosDeContactoCliente.MediosDeContacto.EOtrosMediosDeContacto.Ihtml = template.HTML(GeneraTemplateOtros(Cliente.MediosDeContacto.Otros))

	var templatePerCont string
	for i, Val := range Cliente.PersonasContacto {
		templatePerCont += GeneraGeneraTemplatePersonasContacto(strconv.Itoa(i+1), Val)
	}

	Send.Cliente.EPersonasContactoCliente.Ihtml = template.HTML(templatePerCont)
	ctx.Render("ClienteDetalle.html", Send)

}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send ClienteModel.SCliente

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

	ctx.Render("ClienteDetalle.html", Send)
}

//####################< RUTINAS PARA DIRECCIONES >##########################
//GetPaisesForSelect  funcion que obtiene los paises
func GetPaisesForSelect(ctx *iris.Context) {
	Combox := CargaCombos.CargaComboxPaises("")
	var Send ClienteModel.SCliente
	Send.SEstado = true
	Send.SMsj = Combox
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//GetEstadosForSelect  funcion que obtiene los Estados
func GetEstadosForSelect(ctx *iris.Context) {

	Combox := CargaCombos.CargaComboEstados("")
	var Send ClienteModel.SCliente
	Send.SEstado = true
	Send.SMsj = Combox
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//GetMunicipiosForClaveEstado  funcion que retorna los municipios deacuerdo a una clave  de estado
func GetMunicipiosForClaveEstado(ctx *iris.Context) {
	ID := ctx.FormValue("ID")
	Clave := CatalogoModel.GetClaveEstado(bson.ObjectIdHex(ID))
	Combox := CargaCombos.CargaComboMunicipiosForClaveEstado(Clave, "")
	var Send ClienteModel.SCliente
	Send.SEstado = true
	Send.SMsj = Combox
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//GetColoniasForClaveMunicipio funcion que obtiene las colonias por la clave de el municipio
func GetColoniasForClaveMunicipio(ctx *iris.Context) {

	ID := ctx.FormValue("ID")
	Clave := CatalogoModel.GetClaveMunicipio(bson.ObjectIdHex(ID))
	Combox := CargaCombos.CargaComboColoniasForClaveMunicipio(Clave, "")
	var Send ClienteModel.SCliente
	Send.SEstado = true
	Send.SMsj = Combox
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//GetCPForColonia funcion que obtiene el codigo postal por colonia
func GetCPForColonia(ctx *iris.Context) {
	ID := ctx.FormValue("ID")
	CP := CatalogoModel.GetCPForColonia(bson.ObjectIdHex(ID))
	var Send ClienteModel.SCliente
	Send.SEstado = true
	Send.SMsj = CP
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//GetTipoDireciones funcion que obtiene las Tipos de direcciones
func GetTipoDireciones(ctx *iris.Context) {
	TP := CargaCombos.CargaComboCatalogo(174, "")
	var Send ClienteModel.SCliente
	Send.SEstado = true
	Send.SMsj = TP
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
	return
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send ClienteModel.SCliente

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

		Cabecera, Cuerpo := ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrToMongo))
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
	var Send ClienteModel.SCliente
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Cliente.ENombreCliente.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := ClienteModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = ClienteModel.GeneraTemplatesBusqueda(ClienteModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

//GeneraGeneraTemplatePersonasContacto funcion que genera
func GeneraGeneraTemplatePersonasContacto(ID string, PContacto ClienteModel.PersonaContactoMgo) string {
	var template = `<div id="divPerCon` + ID + `" class="divPerContacto" data-persona-contacto="` + ID + `">
					<div class="CabContactoPersona" data-Num-Contacto-Persona="` + ID + `">
						<div class="titlePerCon">
							<span class="glyphicon glyphicon-user" ></span>
							Persona Contacto
							<div class="botonDeleteCP"> <span class="glyphicon glyphicon-trash" ></span> Eliminar</div>
						</div>
					</div>
					<div id="BodyContactoPersona` + ID + `" ><br/>	
						<div class="form-group">
							<label class="col-sm-4 control-label" for="Nombre` + ID + `">Nombre:</label>
							<div class="col-sm-5">
								<input type="text" name="Nombre` + ID + `" id="Nombre` + ID + `" class="form-control" value="` + PContacto.Nombre + `" >
							</div>
						</div>
						<div class="CollapseDetalle" >
						<div class="form-group">
							<label class="col-sm-4 control-label" for="Estado` + ID + `">Estado:</label>
							<div class="col-sm-5">
								<select id="Estado` + ID + `" name="Estado` + ID + `" class="form-control selectpicker SelectEstado" data-estado-select="` + ID + `">
								` + CargaCombos.CargaComboEstados("") + `
								</select>
							</div>
						</div>
						<div class="form-group">
							<label class="col-sm-4 control-label" for="Municipio` + ID + `">Municipio:</label>
							<div class="col-sm-5">
								<select id="Municipio` + ID + `" name="Municipio` + ID + `" class="form-control selectpicker SelectMunicipio" data-municipo-select="` + ID + `">
								</select>
							</div>
						</div>
						<div class="form-group">
							<label class="col-sm-4 control-label" for="Colonia` + ID + `">Colonia:</label>
							<div class="col-sm-5">
								<select id="Colonia` + ID + `" name="Colonia` + ID + `" class="form-control selectpicker SelectColonia" data-colonia-select="` + ID + `">
								</select>
							</div>
						</div>
						<div class="form-group">
							<label class="col-sm-4 control-label" for="CP` + ID + `">CP:</label>
							<div class="col-sm-5">
								<input type="text" name="CP` + ID + `" id="CP` + ID + `" class="form-control" value="">
							</div>
						</div>
						<div class="form-group">
							<label class="col-sm-4 control-label" for="Calle` + ID + `">Calle:</label>
							<div class="col-sm-5">
								<input type="text" name="Calle` + ID + `" id="Calle` + ID + `" class="form-control" value="">
							</div>
						</div>
						<div class="form-group">
							<label class="col-sm-4 control-label" for="NumExterior` + ID + `">NumExterior:</label>
							<div class="col-sm-5">
								<input type="text" name="NumExterior` + ID + `" id="NumExterior` + ID + `" class="form-control" value="">
							</div>
						</div>
						<div class="form-group">
							<label class="col-sm-4 control-label" for="NumInterior` + ID + `">NumInterior:</label>
							<div class="col-sm-5">
								<input type="text" name="NumInterior` + ID + `" id="NumInterior` + ID + `" class="form-control" value="">
							</div>
						</div>
						<div class="form-group">
							<label class="col-sm-4 control-label" for="TipoDireccion` + ID + `">Tipo Direccion:</label>
							<div class="col-sm-5">
								<select id="TipoDireccion` + ID + `" name="TipoDireccion` + ID + `" class="form-control selectpicker">	
								` + CargaCombos.CargaComboCatalogo(174, "") + `
								</select>
							</div>
						</div>
						<div class="form-group">
							<div class="col-md-9">
								<div id="SufijosDirPersona` + ID + `"></div>
								<button id="AddNewDireccion` + ID + `" type="button" class="btn btn-primary pull-right btn-lg addNewDireccionPersona" data-NewDirect-button="` + ID + `"><span class="glyphicon glyphicon-plus " ></span>Agregar Direccion</button>
								<input type="hidden" name="NumDirContPer` + ID + `"  id="NumDirContPer` + ID + `" value="` + strconv.Itoa(len(PContacto.Direcciones)) + `">
							</div>
						</div>
						</div>
						<h3 class="text-center"><strong>Direcciones</strong></h3>
						<table class="table" id="tableDir` + ID + `">
							<thead>
								<tr>
									<th>Estado</th>
									<th>Municipio</th>
									<th>Colonia</th>
									<th>Cp</th>
									<th>Calle</th>
									<th>Num. Ext</th>
									<th>Num. Int</th>
									<th>Tipo</th>
									<th>Eliminar</th>
								</tr>
							</thead>
							<tbody id="cuerpoDirecciones` + ID + `">
							` + GeneraTRDireccionContactoPersona(ID, PContacto.Direcciones) + `
							</tbody>
						</table>
						<div class="row">
							<div class="col-sm-12">
								<div class="col-md-4">
									<div class="CollapseDetalle" >
									<div class="row text-center" style="padding-bottom:10px;">
										<button id="AgregaEmail` + ID + `" name="AgregaEmail` + ID + `" value="AgregaEmail` + ID + `" data-new-mail="` + ID + `" type="button" class="btn btn-success btn-lg col-md-8 addNewMailPer" style="float:right;margin-right:10%;"><span class="glyphicon glyphicon-plus"></span>Agregar Email</button>
									</div>
									<div class="form-group">
										<label class="col-sm-4 control-label" for="Email` + ID + `">Email:</label>
										<div class="col-sm-8">
											<input type="hidden" name="CorreoPrincipal` + ID + `" id="CorreoPrincipal` + ID + `" value="` + PContacto.MediosDeContacto.Correos.Principal + `" >
											<input type="text" name="Email` + ID + `" id="Email` + ID + `" class="form-control" >
										</div>
									</div>
									</div>
									<h3 class="text-center"><strong>Correos</strong></h3>
									<div class="col-sm-12 table-responsive container" id="div_tabla_correos">
										<table class="table table-condensed table-hover">
											<thead class="thead-inverse">
												<tr>
													<th>Principal</th>
													<th>Correos</th>
													<th>Eliminar</th>
												</tr>
											</thead>
											<tbody id="tbody_etiquetas_correos` + ID + `">
											` + GeneraTemplateCorreosPC(ID, PContacto.MediosDeContacto.Correos.Correos, PContacto.MediosDeContacto.Correos.Principal) + `
											</tbody>
										</table>
									</div>
								</div>
								<div class="col-md-4">
									<div class="CollapseDetalle" >
									<div class="row text-center" style="padding-bottom:10px;">
										<button id="AgregaTelefono` + ID + `" name="AgregaTelefono` + ID + `" type="button" data-new-Tel="` + ID + `" class="btn btn-success btn-lg col-md-8 addNewTelefonoPer" style="float:right;margin-right:10%;"><span class="glyphicon glyphicon-plus"></span>Agregar Telefono</button>
									</div>
									<div class="form-group">
										<label class="col-sm-4 control-label" for="Telefono` + ID + `">Telefono:</label>
										<div class="col-sm-8">
											<input type="hidden" name="TelefonoPrincipal` + ID + `" id="TelefonoPrincipal` + ID + `" value="` + PContacto.MediosDeContacto.Telefonos.Principal + `" >
											<input type="text" name="Telefono` + ID + `" id="Telefono` + ID + `" class="form-control">
										</div>
									</div>
									</div>
									<h3 class="text-center"><strong>Telefonos</strong></h3>
									<div class="col-sm-12 table-responsive container" id="div_tabla_telefonos">
										<table class="table table-condensed table-hover">
											<thead class="thead-inverse">
												<tr>
													<th>Principal</th>
													<th>Telefonos</th>
													<th>Eliminar</th>
												</tr>
											</thead>
											<tbody id="tbody_etiquetas_telefonos` + ID + `">
											` + GeneraTemplateTelefonosCP(ID, PContacto.MediosDeContacto.Telefonos.Telefonos, PContacto.MediosDeContacto.Telefonos.Principal) + `
											</tbody>
										</table>
									</div>
								</div>
								<div class="col-md-4">
									<div class="CollapseDetalle" >
									<div class="row text-center" style="padding-bottom:10px;">
										<button id="AgregaOtro` + ID + `" name="AgregaOtro` + ID + `" type="button" data-new-otro="` + ID + `" class="btn btn-success btn-lg col-md-8 addNewOtroContPer" style="float:right;margin-right:10%;"><span class="glyphicon glyphicon-plus"></span>Agregar Otro Medio</button>
									</div>
									<div class="form-group">
										<label class="col-sm-4 control-label" for="Otro` + ID + `">Otros:</label>
										<div class="col-sm-8">
											<input type="text" name="Otro` + ID + `" id="Otro` + ID + `" class="form-control">
										</div>
									</div>
									</div>
									<h3 class="text-center"><strong>Otros</strong></h3>
									<div class="col-sm-12 table-responsive container" id="div_tabla_otros">
										<table class="table table-condensed table-hover">
											<thead class="thead-inverse">
												<tr>
													<th>Otros</th>
													<th>Eliminar</th>
												</tr>
											</thead>
											<tbody id="tbody_etiquetas_otros` + ID + `">
											` + GeneraTemplateOtrosCP(ID, PContacto.MediosDeContacto.Otros) + ` 
											</tbody>
										</table>
									</div>
								</div>
							</div>
						</div>
					</div>
					</div>`
	return template
}

//GeneraTRDireccionContactoPersona funcion que genera
func GeneraTRDireccionContactoPersona(ID string, Direcciones []ClienteModel.DireccionMgo) string {
	tr := ""
	for i, Direccion := range Direcciones {
		tr += `<tr class='direccionTr` + ID + `'  data-num-dir-contper="` + strconv.Itoa(i+1) + `">
				<td>
					<input type='hidden' name='IDr` + strconv.Itoa(i+1) + `' value='` + Direccion.ID.Hex() + `'>
					<input type='hidden' name='EstadoPC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.Estado.Hex() + `'>` + CatalogoModel.GetNameEstado(Direccion.Estado) + `</td>
				<td><input type='hidden' name='MunicipioPC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.Municipio.Hex() + `'>` + CatalogoModel.GetNameMunicipio(Direccion.Municipio) + `</td>
				<td><input type='hidden' name='ColoniaPC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.Colonia.Hex() + `'>` + CatalogoModel.GetNameColonia(Direccion.Colonia) + `</td>
				<td><input type='hidden' name='cpPC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.CP + `'>` + Direccion.CP + `</td>
				<td><input type='hidden' name='CallePC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.Calle + `'>` + Direccion.Calle + `</td>
				<td><input type='hidden' name='NumExteriorPC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.NumExterior + `'>` + Direccion.NumExterior + `</td>
				<td><input type='hidden' name='NumInteriorPC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.NumInterior + `'>` + Direccion.NumInterior + `</td>
				<td><input type='hidden' name='TipoDireccionPC` + ID + `-` + strconv.Itoa(i+1) + `' value='` + Direccion.TipoDireccion.Hex() + `'>` + CatalogoModel.RegresaNombreSubCatalogo(Direccion.TipoDireccion) + `</td>
				<td><button type="button" class="btn btn-danger deleteDirCP"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
	}
	return tr

}

//GeneraTRDireccionPersonaContacto funcion que genera
func GeneraTRDireccionPersonaContacto(ID string, Estado []string, Municipio []string, Colonia []string, CP []string, Calle []string, NumExterior []string, NumInterior []string, TipoDireccion []string) string {
	fmt.Println("Dierccion::::", Estado, Municipio, Colonia, CP, Calle, NumExterior, NumInterior, TipoDireccion)
	tr := ""
	for i, _ := range Estado {
		tr = tr + `<tr class='direccionTr' >
	 			<td><input type='hidden' name='EstadoPC` + ID + `' value='` + Estado[i] + `'>` + CatalogoModel.GetNameEstado(bson.ObjectIdHex(Estado[i])) + `</td>
	 			<td><input type='hidden' name='MunicipioPC` + ID + `' value='` + Municipio[i] + `'>` + CatalogoModel.GetNameMunicipio(bson.ObjectIdHex(Municipio[i])) + `</td>
	 			<td><input type='hidden' name='ColoniaPC` + ID + `' value='` + Colonia[i] + `'>` + CatalogoModel.GetNameColonia(bson.ObjectIdHex(Colonia[i])) + `</td>
	 			<td><input type='hidden' name='cpPC` + ID + `' value='` + CP[i] + `'>` + CP[i] + `</td>
	 			<td><input type='hidden' name='CallePC` + ID + `' value='` + Calle[i] + `'>` + Calle[i] + `</td>
	 			<td><input type='hidden' name='NumExteriorPC` + ID + `' value='` + NumExterior[i] + `'>` + NumExterior[i] + `</td>
	 			<td><input type='hidden' name='NumInteriorPC` + ID + `' value='` + NumInterior[i] + `'>` + NumInterior[i] + `</td>
	 			<td><input type='hidden' name='TipoDireccionPC` + ID + `' value='` + TipoDireccion[i] + `'>` + CatalogoModel.RegresaNombreSubCatalogo(bson.ObjectIdHex(TipoDireccion[i])) + `</td>
	 			</tr>`
	}
	return tr
}

//GeneraTemplateCorreosPC metodo que genera las tr personas contacto
func GeneraTemplateCorreosPC(ID string, correos []string, principal string) string {
	html := ``
	for _, val := range correos {
		if val == principal {
			html += `<tr>
					<td><input type="radio" name="CorreosPrincipal` + ID + `" value="` + val + `" checked></td>
					<td><input type="text" class="form-control" name="Correos` + ID + `" value="` + val + `" readonly></td>
					<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		} else {
			html += `<tr>
					<td><input type="radio" name="CorreosPrincipal` + ID + `" value="` + val + `"></td>
					<td><input type="text" class="form-control" name="Correos` + ID + `" value="` + val + `" readonly></td>
					<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		}
	}
	return html
}

//GeneraTemplateTelefonosCP Genera la tabla dinamica de los correos ingresados
func GeneraTemplateTelefonosCP(ID string, telefonos []string, principal string) string {
	html := ``
	for _, val := range telefonos {
		if val == principal {
			html += `<tr>
				<td><input type="radio" name="TelefonosPrincipal` + ID + `" value="` + val + `" checked></td>
				<td><input type="text" class="form-control" name="Telefonos` + ID + `" value="` + val + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		} else {
			html += `<tr>
				<td><input type="radio" name="TelefonosPrincipal` + ID + `" value="` + val + `"></td>
				<td><input type="text" class="form-control" name="Telefonos` + ID + `" value="` + val + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
		}
	}
	return html
}

//GeneraTemplateOtrosCP Genera la tabla dinamica de los correos ingresados
func GeneraTemplateOtrosCP(ID string, otros []string) string {
	html := ``
	for _, val := range otros {
		html += `<tr>
				<td><input type="text" class="form-control" name="Otros` + ID + `" value="` + val + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`

	}
	return html
}

//ConvertArrayObjectIDToArrayString funcion que comvierte un arreglo de objetsId a array de String
func ConvertArrayObjectIDToArrayString(ObjetsIds []bson.ObjectId) []string {
	var ArrayString []string
	for _, Val := range ObjetsIds {
		ArrayString = append(ArrayString, Val.Hex())
	}
	return ArrayString
}
