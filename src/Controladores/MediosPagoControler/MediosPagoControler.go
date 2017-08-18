package MediosPagoControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modulos/Session"

	"../../Modelos/CatalogoModel"
	"../../Modelos/MediosPagoModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//##########< Variables Generales > ############
var catalogoNombresPagos = 156
var catalogoTipos = 157
var catalogoEstatus = 158

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

//IndexGet renderea al index de MediosPago
func IndexGet(ctx *iris.Context) {

	var Send MediosPagoModel.SMediosPago

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
	numeroRegistros = MediosPagoModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	MediosPagos := MediosPagoModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range MediosPagos {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagos[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagos[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("MediosPagoIndex.html", Send)
}

//IndexPost regresa la peticon post que se hizo desde el index de MediosPago
func IndexPost(ctx *iris.Context) {
	var Send MediosPagoModel.SMediosPago

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
	//Send.MediosPago.EVARIABLEMediosPago.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := MediosPagoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("MediosPagoIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de MediosPago
func AltaGet(ctx *iris.Context) {
	var Send MediosPagoModel.SMediosPago

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

	Send.MediosPago.ENombreMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoNombresPagos, ""))
	Send.MediosPago.ETipoMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipos, ""))
	ctx.Render("MediosPagoAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de MediosPago
func AltaPost(ctx *iris.Context) {

	var SMediosPago MediosPagoModel.SMediosPago

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SMediosPago.SSesion.Name = NameUsrLoged
	SMediosPago.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SMediosPago.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SMediosPago.SEstado = false
		SMediosPago.SMsj = errSes.Error()
		ctx.Render("ZError.html", SMediosPago)
		return
	}

	EstatusPeticion := false //True indica que hay un error
	var MediosPago MediosPagoModel.MediosPagoMgo
	MediosPago.FechaHora = time.Now()
	MediosPago.ID = bson.NewObjectId()

	nombre := ctx.FormValue("Nombre")
	fmt.Println("nombre:", nombre)
	SMediosPago.MediosPago.ENombreMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoNombresPagos, nombre))
	MediosPago.Nombre = bson.ObjectIdHex(nombre)
	if nombre == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.ENombreMediosPago.IEstatus = true
		SMediosPago.MediosPago.ENombreMediosPago.IMsj = "Seleccione un  Medio de Pago "
	}

	descripcion := ctx.FormValue("Descripcion")
	fmt.Println("descripcion:", descripcion)
	SMediosPago.MediosPago.EDescripcionMediosPago.Descripcion = descripcion
	MediosPago.Descripcion = descripcion
	if descripcion == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.EDescripcionMediosPago.IEstatus = true
		SMediosPago.MediosPago.EDescripcionMediosPago.IMsj = "El campo Descripción no debe estar vacio"
	}

	//Campo Codigo Sat
	NombreSAT := CatalogoModel.GetEspecificByFields("Valores._id", bson.ObjectIdHex(nombre))
	for _, val := range NombreSAT.Valores {
		if nombre == val.ID.Hex() {
			MediosPago.CodigoSat = val.Clave
		}
	}

	tipo := ctx.FormValue("Tipo")
	fmt.Println("tipo:", tipo)
	SMediosPago.MediosPago.ETipoMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipos, tipo))
	MediosPago.Tipo = bson.ObjectIdHex(tipo)
	if tipo == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.ETipoMediosPago.IEstatus = true
		SMediosPago.MediosPago.ETipoMediosPago.IMsj = "Seleccione un Tipo de Comisión"
	}

	valorCom := ctx.FormValue("Comision")
	valor64, e := strconv.ParseFloat(valorCom, 64)
	if e != nil {
		fmt.Println(e)
	}
	SMediosPago.MediosPago.EComisionMediosPago.Comision = float64(valor64)
	MediosPago.Comision = float64(valor64)
	if valorCom == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.EComisionMediosPago.IEstatus = true
		SMediosPago.MediosPago.EComisionMediosPago.IMsj = "Debe agregar una comision mayor o igual a 0"
	}

	cambio := ctx.FormValue("Cambio")
	fmt.Println("cambio", cambio)
	MediosPago.Cambio = false
	if cambio != "" {
		SMediosPago.MediosPago.ECambioMediosPago.Ihtml = template.HTML("checked")
		MediosPago.Cambio = true
	}

	MediosPago.Estatus = CatalogoModel.RegresaIDEstatusActivo(catalogoEstatus)
	//	SMediosPago.MediosPago = MediosPago //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.
	fmt.Println(MediosPago)
	if EstatusPeticion {
		SMediosPago.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SMediosPago.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("MediosPagoAlta.html", SMediosPago)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if MediosPago.InsertaMgo() {
			if MediosPago.InsertaElastic() {
				SMediosPago.SEstado = true
				SMediosPago.SMsj = "Se ha realizado una inserción exitosa"
				ctx.Redirect("/MediosPagos/detalle/"+MediosPago.ID.Hex(), 301)
			} else {
				SMediosPago.SEstado = false
				SMediosPago.SMsj = "Ocurrió un error al insertar en el servidor de búsqueda, intente más tarde"
				ctx.Render("MediosPagoAlta.html", SMediosPago)
			}

		} else {
			SMediosPago.SEstado = false
			SMediosPago.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("MediosPagoAlta.html", SMediosPago)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de MediosPago
func EditaGet(ctx *iris.Context) {

	var SMediosPago MediosPagoModel.SMediosPago

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SMediosPago.SSesion.Name = NameUsrLoged
	SMediosPago.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SMediosPago.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SMediosPago.SEstado = false
		SMediosPago.SMsj = errSes.Error()
		ctx.Render("ZError.html", SMediosPago)
		return
	}

	id := ctx.Param("ID")
	mediopago := MediosPagoModel.GetOne(bson.ObjectIdHex(id))

	SMediosPago.MediosPago.ID = mediopago.ID
	SMediosPago.MediosPago.ENombreMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoNombresPagos, mediopago.Nombre.Hex()))
	SMediosPago.MediosPago.EDescripcionMediosPago.Descripcion = mediopago.Descripcion
	SMediosPago.MediosPago.ETipoMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipos, mediopago.Tipo.Hex()))
	SMediosPago.MediosPago.EComisionMediosPago.Comision = mediopago.Comision
	SMediosPago.MediosPago.EEstatusMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoEstatus, mediopago.Estatus.Hex()))
	if mediopago.Cambio {
		SMediosPago.MediosPago.ECambioMediosPago.Ihtml = template.HTML("checked")
	}
	SMediosPago.MediosPago.ECodigoSatMediosPago.CodigoSat = mediopago.CodigoSat
	ctx.Render("MediosPagoEdita.html", SMediosPago)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de MediosPago
func EditaPost(ctx *iris.Context) {
	var SMediosPago MediosPagoModel.SMediosPago

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SMediosPago.SSesion.Name = NameUsrLoged
	SMediosPago.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SMediosPago.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SMediosPago.SEstado = false
		SMediosPago.SMsj = errSes.Error()
		ctx.Render("ZError.html", SMediosPago)
		return
	}

	errorID := ""
	EstatusPeticion := false
	var MediosPago MediosPagoModel.MediosPagoMgo
	medioexistente := MediosPagoModel.MediosPagoMgo{}

	id := ctx.FormValue("IDname")
	if bson.IsObjectIdHex(id) {
		MediosPago.ID = bson.ObjectIdHex(id)
		SMediosPago.MediosPago.ID = bson.ObjectIdHex(id)
		medioexistente = MediosPagoModel.GetOne(bson.ObjectIdHex(id))
		if medioexistente.ID.Hex() == "" {
			EstatusPeticion = true
			errorID = ", Error al buscar el el Medio de Pago, intente mas tarde"
		}
	} else {
		EstatusPeticion = true
		errorID = ", Error en la referencia del Medio de Pago"
	}

	nombre := ctx.FormValue("Nombre")
	fmt.Println("nombre:", nombre)
	SMediosPago.MediosPago.ENombreMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoNombresPagos, nombre))
	MediosPago.Nombre = bson.ObjectIdHex(nombre)
	if nombre == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.ENombreMediosPago.IEstatus = true
		SMediosPago.MediosPago.ENombreMediosPago.IMsj = "Seleccione un  Medio de Pago "
	}

	descripcion := ctx.FormValue("Descripcion")
	fmt.Println("descripcion:", descripcion)
	SMediosPago.MediosPago.EDescripcionMediosPago.Descripcion = descripcion
	MediosPago.Descripcion = descripcion
	if descripcion == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.EDescripcionMediosPago.IEstatus = true
		SMediosPago.MediosPago.EDescripcionMediosPago.IMsj = "El campo Descripción no debe estar vacio"
	}

	//Campo Codigo Sat
	NombreSAT := CatalogoModel.GetEspecificByFields("Valores._id", bson.ObjectIdHex(nombre))
	for _, val := range NombreSAT.Valores {
		if nombre == val.ID.Hex() {
			MediosPago.CodigoSat = val.Clave
			SMediosPago.MediosPago.ECodigoSatMediosPago.CodigoSat = val.Clave
		}
	}

	tipo := ctx.FormValue("Tipo")
	fmt.Println("tipo:", tipo)
	SMediosPago.MediosPago.ETipoMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipos, tipo))
	MediosPago.Tipo = bson.ObjectIdHex(tipo)
	if tipo == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.ETipoMediosPago.IEstatus = true
		SMediosPago.MediosPago.ETipoMediosPago.IMsj = "Seleccione un Tipo de Comisión"
	}

	valorCom := ctx.FormValue("Comision")
	valor64, e := strconv.ParseFloat(valorCom, 64)
	if e != nil {
		fmt.Println(e)
	}
	SMediosPago.MediosPago.EComisionMediosPago.Comision = float64(valor64)
	MediosPago.Comision = float64(valor64)
	if valorCom == "" {
		EstatusPeticion = true
		SMediosPago.MediosPago.EComisionMediosPago.IEstatus = true
		SMediosPago.MediosPago.EComisionMediosPago.IMsj = "Debe agregar una comision mayor o igual a 0"
	}

	cambio := ctx.FormValue("Cambio")
	fmt.Println("cambio", cambio)
	MediosPago.Cambio = false
	if cambio != "" {
		SMediosPago.MediosPago.ECambioMediosPago.Ihtml = template.HTML("checked")
		MediosPago.Cambio = true
	}

	estatus := ctx.FormValue("Estatus")
	MediosPago.Estatus = bson.ObjectIdHex(estatus)
	SMediosPago.MediosPago.EEstatusMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoEstatus, estatus))
	if EstatusPeticion {
		SMediosPago.SEstado = false                                                                     //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SMediosPago.SMsj = "La validación indica que el objeto capturado no puede procesarse" + errorID //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("MediosPagoEdita.html", SMediosPago)
	} else {
		MediosPago.FechaHora = time.Now()
		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if MediosPago.ReemplazaMgo() {
			errupd := MediosPago.ActualizaElastic()
			if errupd == nil {
				SMediosPago.SEstado = true
				SMediosPago.SMsj = "Se ha realizado una inserción exitosa"
				ctx.Redirect("/MediosPagos/detalle/"+MediosPago.ID.Hex(), 301)
			} else {
				if medioexistente.ReemplazaMgo() {
					SMediosPago.SEstado = false
					SMediosPago.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + "). Se ha reestablecido la informacion"
					ctx.Render("MediosPagoDetalle.html", SMediosPago)
				} else {
					SMediosPago.SEstado = false
					SMediosPago.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + ") No se pudo reestablecer la informacion"
					ctx.Render("MediosPagoEdita.html", SMediosPago)
				}
			}

		} else {
			SMediosPago.SEstado = false
			SMediosPago.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("MediosPagoEdita.html", SMediosPago)
		}

	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {

	var SMediosPago MediosPagoModel.SMediosPago
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SMediosPago.SSesion.Name = NameUsrLoged
	SMediosPago.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SMediosPago.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SMediosPago.SEstado = false
		SMediosPago.SMsj = errSes.Error()
		ctx.Render("ZError.html", SMediosPago)
		return
	}

	id := ctx.Param("ID")
	mediopago := MediosPagoModel.GetOne(bson.ObjectIdHex(id))

	SMediosPago.MediosPago.ID = bson.ObjectIdHex(id)
	SMediosPago.MediosPago.ENombreMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoNombresPagos, mediopago.Nombre.Hex()))
	SMediosPago.MediosPago.EDescripcionMediosPago.Descripcion = mediopago.Descripcion
	SMediosPago.MediosPago.ETipoMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipos, mediopago.Tipo.Hex()))
	SMediosPago.MediosPago.EComisionMediosPago.Comision = mediopago.Comision
	SMediosPago.MediosPago.EEstatusMediosPago.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoEstatus, mediopago.Estatus.Hex()))
	if mediopago.Cambio {
		SMediosPago.MediosPago.ECambioMediosPago.Ihtml = template.HTML("checked")
	}

	SMediosPago.MediosPago.ECodigoSatMediosPago.CodigoSat = mediopago.CodigoSat
	SMediosPago.MediosPago.EFechaHoraMediosPago.FechaHora = mediopago.FechaHora

	ctx.Render("MediosPagoDetalle.html", SMediosPago)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send MediosPagoModel.SMediosPago

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

	ctx.Render("MediosPagoDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send MediosPagoModel.SMediosPago

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

		Cabecera, Cuerpo := MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrToMongo))
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
	var Send MediosPagoModel.SMediosPago
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.MediosPago.ENombreMediosPago.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := MediosPagoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = MediosPagoModel.GeneraTemplatesBusqueda(MediosPagoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
