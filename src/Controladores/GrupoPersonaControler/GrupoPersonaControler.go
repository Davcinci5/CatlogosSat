package GrupoPersonaControler

import (
	"html/template"
	"strconv"
	"time"

	"fmt"

	"../../Modelos/GrupoPersonaModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Redis"
	"../../Modulos/Session"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//##########< Variables Generales > ############

var cadenaBusqueda string
var buscarEn string
var numeroRegistros int64
var paginasTotales int

//NumPagina especifica ***************
var NumPagina float32

//limitePorPagina especifica ***************
var limitePorPagina = 10
var result []GrupoPersonaModel.GrupoPersona
var resultPage []GrupoPersonaModel.GrupoPersona
var templatePaginacion = ``

//Numero de Catalogos
var CatalogoEstatusGruposPer = 146
var CatalogoUris = 177

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de GrupoPersona
func IndexGet(ctx *iris.Context) {
	var Send GrupoPersonaModel.SGrupoPersona
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

	IDURI := CargaCombos.ExisteEnCatalogo(CatalogoUris, MoGeneral.MiURI(ctx.Request.RequestURI, ctx.Param("ID")))
	IDUSR := GrupoPersonaModel.GetIDUsuario(NameUsrLoged)
	Permiso, err := PermisosUsuario(IDUSR, IDURI)

	if err != nil { //Ocurrio algun Error?
		Send.SEstado = false
		Send.SMsj = err.Error()
		ctx.Render("ZError.html", Send)
	} else {
		if !Permiso { //No Tiene Permitito acceder?
			ctx.EmitError(iris.StatusForbidden)
		} else { //Si tiene el permiso, aqui se tratan los datos
			ctx.Render("GrupoPersonaIndex.html", Send)
		} //Fin del else{} de permisos
	}
}

//IndexPost regresa la peticon post que se hizo desde el index de GrupoPersona
func IndexPost(ctx *iris.Context) {

	templatePaginacion = ``

	var resultados []GrupoPersonaModel.GrupoPersonaMgo
	var IDToObjID bson.ObjectId
	var arrObjIds []bson.ObjectId
	var arrToMongo []bson.ObjectId

	cadenaBusqueda = ctx.FormValue("searchbox")
	buscarEn = ctx.FormValue("buscaren")

	if cadenaBusqueda != "" {

		docs := GrupoPersonaModel.BuscarEnElastic(cadenaBusqueda)

		if docs.Hits.TotalHits > 0 {
			numeroRegistros = docs.Hits.TotalHits

			paginasTotales = Totalpaginas()

			for _, item := range docs.Hits.Hits {
				IDToObjID = bson.ObjectIdHex(item.Id)
				arrObjIds = append(arrObjIds, IDToObjID)
			}

			if numeroRegistros <= int64(limitePorPagina) {
				for _, v := range arrObjIds[0:numeroRegistros] {
					arrToMongo = append(arrToMongo, v)
				}
			} else if numeroRegistros >= int64(limitePorPagina) {
				for _, v := range arrObjIds[0:limitePorPagina] {
					arrToMongo = append(arrToMongo, v)
				}
			}

			resultados = GrupoPersonaModel.GetEspecifics(arrToMongo)

			MoConexion.FlushElastic()

		}

	}

	templatePaginacion = ConstruirPaginacion()

	ctx.Render("GrupoPersonaIndex.html", map[string]interface{}{
		"result":          resultados,
		"cadena_busqueda": cadenaBusqueda,
		"PaginacionT":     template.HTML(templatePaginacion),
	})

}

//###########################< ALTA >################################

//AltaGet renderea al alta de GrupoPersona
func AltaGet(ctx *iris.Context) {
	var Send GrupoPersonaModel.SGrupoPersona

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

	IDURI := CargaCombos.ExisteEnCatalogo(CatalogoUris, MoGeneral.MiURI(ctx.Request.RequestURI, ctx.Param("ID")))
	IDUSR := GrupoPersonaModel.GetIDUsuario(NameUsrLoged)
	Permiso, err := PermisosUsuario(IDUSR, IDURI)

	if err != nil { //Ocurrio algun Error?
		Send.SEstado = false
		Send.SMsj = err.Error()
		ctx.Render("ZError.html", Send)
	} else {
		if !Permiso { //No Tiene Permitito acceder?
			ctx.EmitError(iris.StatusForbidden)
		} else { //Si tiene el permiso, aqui se tratan los datos
			usuariosgrupos := GrupoPersonaModel.GeneraEtiquetasPersonas()
			Send.GrupoPersona.ENoMiembrosGrupoPersona.Ihtml = template.HTML(usuariosgrupos)
			ctx.Render("GrupoPersonaAlta.html", Send)
		} //Fin del else{} de permisos
	}

}

//AltaPost regresa la petición post que se hizo desde el alta de GrupoPersona
func AltaPost(ctx *iris.Context) {
	EstatusPeticion := false

	//######### ENVIA TUS RESULTADOS #########
	var Send GrupoPersonaModel.SGrupoPersona
	var SendMgo GrupoPersonaModel.GrupoPersonaMgo

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

	IDURI := CargaCombos.ExisteEnCatalogo(CatalogoUris, MoGeneral.MiURI(ctx.Request.RequestURI, ctx.Param("ID")))
	IDUSR := GrupoPersonaModel.GetIDUsuario(NameUsrLoged)
	Permiso, err := PermisosUsuario(IDUSR, IDURI)

	if err != nil { //Ocurrio algun Error?
		Send.SEstado = false
		Send.SMsj = err.Error()
		ctx.Render("ZError.html", Send)
	} else {
		if !Permiso { //No Tiene Permitito acceder?
			ctx.EmitError(iris.StatusForbidden)
		} else { //Si tiene el permiso, aqui se tratan los datos

			Nombre := ctx.FormValue("Nombre")
			Send.ENombreGrupoPersona.Nombre = Nombre
			SendMgo.Nombre = Nombre
			if Nombre == "" {
				EstatusPeticion = true
				Send.ENombreGrupoPersona.IEstatus = true
				Send.ENombreGrupoPersona.IMsj = "El campo Nombre es obligatorio"
			}

			Descripcion := ctx.FormValue("Descripcion")
			Send.EDescripcionGrupoPersona.Descripcion = Descripcion
			SendMgo.Descripcion = Descripcion
			if Descripcion == "" {
				EstatusPeticion = true
				Send.EDescripcionGrupoPersona.IEstatus = true
				Send.EDescripcionGrupoPersona.IMsj = "El campo Descripción es obligatorio"
			}

			Miembros := ctx.Request.Form["Miembros"]

			if len(Miembros) == 0 {
				EstatusPeticion = true
				usuariosgrupos := GrupoPersonaModel.GeneraEtiquetasPersonas()
				Send.GrupoPersona.ENoMiembrosGrupoPersona.Ihtml = template.HTML(usuariosgrupos)
				Send.EMiembrosGrupoPersona.IEstatus = true
				Send.EMiembrosGrupoPersona.IMsj = "El grupo debe tener al menos un Miembro"
			} else {
				MiembrosIDS := GrupoPersonaModel.ConvierteAObjectIDS(Miembros)
				SendMgo.Miembros = MiembrosIDS
				UsuariosDentro := GrupoPersonaModel.GeneraEtiquetasPersonasEnUnGpo(MiembrosIDS)
				Send.GrupoPersona.EMiembrosGrupoPersona.Ihtml = template.HTML(UsuariosDentro)

				OtrosPersonaGpos := ctx.Request.Form["PersonaGpo"]
				IDSNoMiembros := GrupoPersonaModel.ConvierteAObjectIDS(OtrosPersonaGpos)
				UsuariosFuera := GrupoPersonaModel.GeneraEtiquetasPersonasFueraDeUnGpo(IDSNoMiembros)
				Send.GrupoPersona.ENoMiembrosGrupoPersona.Ihtml = template.HTML(UsuariosFuera)
			}

			SendMgo.Estatus = CargaCombos.CargaEstatusActivoEnAlta(CatalogoEstatusGruposPer)
			SendMgo.FechaHora = time.Now()

			if EstatusPeticion {
				Send.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
				Send.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
			} else {
				SendMgo.ID = bson.NewObjectId()
				//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
				if SendMgo.InsertaMgo() {
					Send.SEstado = true
					Send.SMsj = "Se ha realizado una inserción exitosa"
					ctx.Redirect("/GrupoPersonas/detalle/"+SendMgo.ID.Hex(), 301)
				} else {
					Send.SEstado = false
					Send.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
				}

			}

			ctx.Render("GrupoPersonaAlta.html", Send)

		} //Fin del else{} de permisos
	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de GrupoPersona
func EditaGet(ctx *iris.Context) {
	// name, nivel, id := Session.GetUserName(ctx.Request)

	var Send GrupoPersonaModel.SGrupoPersona
	var Grupo GrupoPersonaModel.GrupoPersonaMgo

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

	IDURI := CargaCombos.ExisteEnCatalogo(CatalogoUris, MoGeneral.MiURI(ctx.Request.RequestURI, ctx.Param("ID")))
	IDUSR := GrupoPersonaModel.GetIDUsuario(NameUsrLoged)
	Permiso, err := PermisosUsuario(IDUSR, IDURI)

	if err != nil { //Ocurrio algun Error?
		Send.SEstado = false
		Send.SMsj = err.Error()
		ctx.Render("ZError.html", Send)
	} else {
		if !Permiso { //No Tiene Permitito acceder?
			ctx.EmitError(iris.StatusForbidden)
		} else { //Si tiene el permiso, aqui se tratan los datos
			ID := ctx.Param("ID")
			if bson.IsObjectIdHex(ID) {
				Grupo = GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
				if Grupo.ID.Hex() != "" {
					Send.GrupoPersona.ID = Grupo.ID
					Send.GrupoPersona.ENombreGrupoPersona.Nombre = Grupo.Nombre
					Send.GrupoPersona.EDescripcionGrupoPersona.Descripcion = Grupo.Descripcion
					Send.GrupoPersona.EEstatusGrupoPersona.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoEstatusGruposPer, Grupo.Estatus.Hex()))
					Send.GrupoPersona.EFechaHoraGrupoPersona.FechaHora = Grupo.FechaHora
					UsuariosDentro := GrupoPersonaModel.GeneraEtiquetasPersonasEnUnGpo(Grupo.Miembros)
					Send.GrupoPersona.EMiembrosGrupoPersona.Ihtml = template.HTML(UsuariosDentro)
					UsuariosFuera := GrupoPersonaModel.GeneraEtiquetasPersonasFueraDeUnGpo(Grupo.Miembros)
					Send.GrupoPersona.ENoMiembrosGrupoPersona.Ihtml = template.HTML(UsuariosFuera)
					// Encontrar diferentes de... mongo
					// "_id":{$ne: ObjectId("5919d5e8d2b2132204a793fa")}
				} else {
					Send.SEstado = false
					Send.SMsj = "Grupo no encontrado"
				}
			} else {
				Send.SEstado = false
				Send.SMsj = "Error en la referencia al Grupo"
			}
			ctx.Render("GrupoPersonaEdita.html", Send)

		} //Fin del else{} de permisos
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de GrupoPersona
func EditaPost(ctx *iris.Context) {
	EstatusPeticion := false
	var Send GrupoPersonaModel.SGrupoPersona
	var Grupo GrupoPersonaModel.GrupoPersonaMgo

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

	IDURI := CargaCombos.ExisteEnCatalogo(CatalogoUris, MoGeneral.MiURI(ctx.Request.RequestURI, ctx.Param("ID")))
	IDUSR := GrupoPersonaModel.GetIDUsuario(NameUsrLoged)
	Permiso, err := PermisosUsuario(IDUSR, IDURI)

	if err != nil { //Ocurrio algun Error?
		Send.SEstado = false
		Send.SMsj = err.Error()
		ctx.Render("ZError.html", Send)
	} else {

		if !Permiso { //No Tiene Permitito acceder?
			ctx.EmitError(iris.StatusForbidden)
		} else { //Si tiene el permiso, aqui se tratan los datos
			ID := ctx.FormValue("IDname")
			if bson.IsObjectIdHex(ID) {
				Grupo = GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
				if Grupo.ID.Hex() != "" {
					Send.GrupoPersona.ID = Grupo.ID
					Nombre := ctx.FormValue("Nombre")
					Send.ENombreGrupoPersona.Nombre = Nombre
					Grupo.Nombre = Nombre
					if Nombre == "" {
						EstatusPeticion = true
						Send.ENombreGrupoPersona.IEstatus = true
						Send.ENombreGrupoPersona.IMsj = "El campo Nombre es obligatorio"
					}

					Descripcion := ctx.FormValue("Descripcion")
					Send.EDescripcionGrupoPersona.Descripcion = Descripcion
					Grupo.Descripcion = Descripcion
					if Descripcion == "" {
						EstatusPeticion = true
						Send.EDescripcionGrupoPersona.IEstatus = true
						Send.EDescripcionGrupoPersona.IMsj = "El campo Descripción es obligatorio"
					}

					Miembros := ctx.Request.Form["Miembros"]
					if len(Miembros) == 0 {
						EstatusPeticion = true
						usuariosgrupos := GrupoPersonaModel.GeneraEtiquetasPersonas()
						Send.GrupoPersona.ENoMiembrosGrupoPersona.Ihtml = template.HTML(usuariosgrupos)
						Send.EMiembrosGrupoPersona.IEstatus = true
						Send.EMiembrosGrupoPersona.IMsj = "El grupo debe tener almenos un Miembro"
					} else {
						MiembrosIDS := GrupoPersonaModel.ConvierteAObjectIDS(Miembros)
						Grupo.Miembros = MiembrosIDS
						UsuariosDentro := GrupoPersonaModel.GeneraEtiquetasPersonasEnUnGpo(MiembrosIDS)
						Send.GrupoPersona.EMiembrosGrupoPersona.Ihtml = template.HTML(UsuariosDentro)
						UsuariosFuera := GrupoPersonaModel.GeneraEtiquetasPersonasFueraDeUnGpo(MiembrosIDS)
						Send.GrupoPersona.ENoMiembrosGrupoPersona.Ihtml = template.HTML(UsuariosFuera)
					}

					Estatus := ctx.FormValue("Estatus")
					if bson.IsObjectIdHex(Estatus) {
						Grupo.Estatus = bson.ObjectIdHex(Estatus)
						Send.GrupoPersona.EEstatusGrupoPersona.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoEstatusGruposPer, Estatus))
					} else {
						EstatusPeticion = true
						Send.GrupoPersona.EEstatusGrupoPersona.IEstatus = true
						Send.GrupoPersona.EEstatusGrupoPersona.IMsj = "Valor no valido para el campo Estatus"
					}

					if EstatusPeticion {
						Send.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
						Send.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
					} else {
						//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
						if Grupo.ReemplazaMgo() {
							Send.SEstado = true
							Send.SMsj = "Se ha realizado una Actualización exitosa"
							ctx.Redirect("/GrupoPersonas/detalle/"+Grupo.ID.Hex(), 301)
						} else {
							Send.SEstado = false
							Send.SMsj = "Ocurrió un error al Actualizar el Objeto, intente más tarde"
						}

					}

				} else {
					Send.SEstado = false
					Send.SMsj = "Grupo no encontrado"
				}
			} else {
				Send.SEstado = false
				Send.SMsj = "Error en la referencia al Grupo"
			}
			ctx.Render("GrupoPersonaEdita.html", Send)

		} //Fin del else{} de permisos
	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send GrupoPersonaModel.SGrupoPersona
	var Grupo GrupoPersonaModel.GrupoPersonaMgo

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

	IDURI := CargaCombos.ExisteEnCatalogo(CatalogoUris, MoGeneral.MiURI(ctx.Request.RequestURI, ctx.Param("ID")))
	IDUSR := GrupoPersonaModel.GetIDUsuario(NameUsrLoged)
	Permiso, err := PermisosUsuario(IDUSR, IDURI)

	if err != nil { //Ocurrio algun Error?
		Send.SEstado = false
		Send.SMsj = err.Error()
		ctx.Render("ZError.html", Send)
	} else {
		if !Permiso { //No Tiene Permitito acceder?
			ctx.EmitError(iris.StatusForbidden)
		} else { //Si tiene el permiso, aqui se tratan los datos
			ID := ctx.Param("ID")
			if bson.IsObjectIdHex(ID) {
				Grupo = GrupoPersonaModel.GetOne(bson.ObjectIdHex(ID))
				if Grupo.ID.Hex() != "" {
					Send.GrupoPersona.ID = Grupo.ID
					Send.GrupoPersona.ENombreGrupoPersona.Nombre = Grupo.Nombre
					Send.GrupoPersona.EDescripcionGrupoPersona.Descripcion = Grupo.Descripcion
					Send.GrupoPersona.EEstatusGrupoPersona.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(CatalogoEstatusGruposPer, Grupo.Estatus.Hex()))
					Send.GrupoPersona.EFechaHoraGrupoPersona.FechaHora = Grupo.FechaHora
					UsuariosDentro := GrupoPersonaModel.GeneraEtiquetasPersonasEnUnGpo(Grupo.Miembros)
					Send.GrupoPersona.EMiembrosGrupoPersona.Ihtml = template.HTML(UsuariosDentro)
				} else {
					Send.SEstado = false
					Send.SMsj = "Grupo no encontrado"
				}
			} else {
				Send.SEstado = false
				Send.SMsj = "Error en la referencia al Grupo"
			}
			ctx.Render("GrupoPersonaDetalle.html", Send)
		} //Fin del else{} de permisos
	}

}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	ctx.Render("GrupoPersonaDetalle.html", nil)
}

//####################< RUTINAS ADICIONALES >##########################

//Totalpaginas calcula el número de paginaciones de acuerdo al número
// de resultados encontrados y los que se quieren mostrar en la página.
func Totalpaginas() int {

	NumPagina = float32(numeroRegistros) / float32(limitePorPagina)
	NumPagina2 := int(NumPagina)
	if NumPagina > float32(NumPagina2) {
		NumPagina2++
	}
	totalpaginas := NumPagina2
	return totalpaginas

}

//ConstruirPaginacion construtye la paginación en formato html para usarse en la página
func ConstruirPaginacion() string {
	var templateP string
	templateP += `
	<nav aria-label="Page navigation">
		<ul class="pagination">
			<li>
				<a href="/GrupoPersonas/1" aria-label="Primera">
				<span aria-hidden="true">&laquo;</span>
				</a>
			</li>`

	templateP += ``
	for i := 0; i <= paginasTotales; i++ {
		if i == 1 {

			templateP += `<li class="active"><a href="/GrupoPersonas/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		} else if i > 1 && i < 11 {
			templateP += `<li><a href="/GrupoPersonas/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`

		} else if i > 11 && i == paginasTotales {
			templateP += `<li><span aria-hidden="true">...</span></li><li><a href="/GrupoPersonas/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		}
	}
	templateP += `<li><a href="/GrupoPersonas/` + strconv.Itoa(paginasTotales) + `" aria-label="Ultima"><span aria-hidden="true">&raquo;</span></a></li></ul></nav>`
	return templateP
}

//PermisosUsuario verifica los permisos de un usuario a una URI
func PermisosUsuario(ID string, URI string) (bool, error) {
	permisos := GrupoPersonaModel.RegresaGruposDeUsuario(bson.ObjectIdHex(ID))
	permisos = append(permisos, ID)
	exist := false
	var err error
	for _, val := range permisos {
		p, err := Redis.ObtenerMiemmbroDeGrupo(val, URI)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		if p {
			exist = true
		}
	}
	return exist, err
}
