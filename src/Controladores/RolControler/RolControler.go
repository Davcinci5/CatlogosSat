package RolControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modulos/Session"

	"../../Modelos/RolModel"
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

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Rol
func IndexGet(ctx *iris.Context) {

	var Send RolModel.SRol
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
	numeroRegistros = RolModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Rols := RolModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Rols {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(Rols[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(Rols[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("RolIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Rol
func IndexPost(ctx *iris.Context) {

	var Send RolModel.SRol
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
	//Send.Rol.EVARIABLERol.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := RolModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("RolIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Rol
func AltaGet(ctx *iris.Context) {
	var Send RolModel.SRol
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

	Send.EPermisosRol.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(163, ""))

	ctx.Render("RolAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Rol
func AltaPost(ctx *iris.Context) {
	var SRol RolModel.SRol
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SRol.SSesion.Name = NameUsrLoged
	SRol.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SRol.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SRol.SEstado = false
		SRol.SMsj = errSes.Error()
		ctx.Render("ZError.html", SRol)
		return
	}

	//######### LEE TU OBJETO DEL FORMULARIO #########

	fmt.Println("=================ALTA DE ROL========================")

	var Rol RolModel.RolMgo
	ctx.ReadForm(&Rol)

	EstatusPeticion := false //True indica que hay un error

	//Campo ID
	Rol.ID = bson.NewObjectId()
	SRol.ID = Rol.ID

	//Campo Nombre
	if Rol.Nombre == "" {
		EstatusPeticion = true
		SRol.ENombreRol.IEstatus = true
		SRol.ENombreRol.IMsj = "El nombre no puede estar vacio"
	} else {
		SRol.ENombreRol.IEstatus = false
		SRol.ENombreRol.Nombre = Rol.Nombre
		SRol.ENombreRol.Ihtml = template.HTML(Rol.Nombre)
	}

	//Campo Descripcion
	if Rol.Descripcion == "" {
		SRol.EDescripcionRol.IEstatus = false
		SRol.EDescripcionRol.Descripcion = Rol.Descripcion

	} else {
		SRol.EDescripcionRol.IEstatus = false
		SRol.EDescripcionRol.Descripcion = Rol.Descripcion
	}

	//Campo Estatus
	Rol.Estatus = CargaCombos.CargaEstatusActivoEnAlta(164)
	SRol.EEstatusRol.Estatus = Rol.Estatus
	SRol.EEstatusRol.Ihtml = `<option>ACTIVO</option>`

	//Campo FechaHora
	Rol.FechaHora = time.Now()
	SRol.EFechaHoraRol.FechaHora = Rol.FechaHora

	//Campo Permisos
	if Rol.Permisos == nil {
		SRol.EPermisosRol.IEstatus = true
		EstatusPeticion = true
		SRol.EPermisosRol.IMsj = "Debes seleccionar un permiso"
		SRol.EPermisosRol.Ihtml = template.HTML(CargaCombos.CargaComboCatalogoMulti(163, ""))
	} else {
		SRol.EPermisosRol.Permisos = Rol.Permisos
		SRol.EPermisosRol.Ihtml = template.HTML(RolModel.CargaComboCatalogoArrayID(163, Rol))
	}

	//Debo construir un multiselect para que atrapemos varios objects ids para los permisos que tenga un rol.

	if EstatusPeticion {
		SRol.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SRol.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("RolAlta.html", SRol)
	} else {
		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if Rol.InsertaMgo() {
			if Rol.InsertaElastic() {
				SRol.SEstado = true
				SRol.SMsj = "Se ha realizado una inserción exitosa"
				ctx.Render("RolDetalle.html", SRol)
			} else {
				SRol.SEstado = false
				SRol.SMsj = "Ocurrió un error al insertar en el servidor de busqueda, inténtelo mas tarde"
				ctx.Render("RolAlta.html", SRol)
			}

		} else {
			SRol.SEstado = false
			SRol.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("RolAlta.html", SRol)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Rol
func EditaGet(ctx *iris.Context) {
	var Send RolModel.SRol
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

	ctx.Render("RolEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Rol
func EditaPost(ctx *iris.Context) {
	var Send RolModel.SRol
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

	ctx.Render("RolEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send RolModel.SRol
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

	ctx.Render("RolDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send RolModel.SRol
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

	ctx.Render("RolDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send RolModel.SRol

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

		Cabecera, Cuerpo := RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrToMongo))
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
	var Send RolModel.SRol
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Rol.ENombreRol.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := RolModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = RolModel.GeneraTemplatesBusqueda(RolModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
