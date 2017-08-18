//Package CatalogoControler es el controlador de catálogos del sistema
package CatalogoControler

import (
	"encoding/json"
	"html/template"
	"strconv"
	"time"

	"../../Modulos/Session"

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

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Catalogo
func IndexGet(ctx *iris.Context) {

	var Send CatalogoModel.SCatalogo

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
	numeroRegistros = CatalogoModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Catalogos := CatalogoModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Catalogos {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(Catalogos[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(Catalogos[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("CatalogoIndex.html", Send)

}

//IndexPost renderea al index de Catalogo
func IndexPost(ctx *iris.Context) {

	var Send CatalogoModel.SCatalogo

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
	//Send.Catalogo.EVARIABLECatalogo.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := CatalogoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
		Send.SIndex.SResultados = false
	}
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	ctx.Render("CatalogoIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Catalogo
func AltaGet(ctx *iris.Context) {

	var Send CatalogoModel.SCatalogo
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

	ctx.Render("CatalogoAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Catalogo
func AltaPost(ctx *iris.Context) {
	var CataE CatalogoModel.SCatalogo
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	CataE.SSesion.Name = NameUsrLoged
	CataE.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	CataE.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		CataE.SEstado = false
		CataE.SMsj = errSes.Error()
		ctx.Render("ZError.html", CataE)
		return
	}

	EstatusPeticion := false //Comienza suponiendo que no hay error

	var Cata CatalogoModel.Catalogo

	var Catalogo CatalogoModel.CatalogoMgo

	descripcion := ctx.FormValue("Descripcion")
	Cata.EDescripcionCatalogo.Descripcion = descripcion
	Catalogo.Descripcion = descripcion

	if descripcion != "" {
		if len(descripcion) < 20 || len(descripcion) > 250 {
			EstatusPeticion = true
			Cata.EDescripcionCatalogo.IEstatus = true
			Cata.EDescripcionCatalogo.IMsj = "La descripción del catálogo debe tener una longitud entre [20 y 250] caracteres."
		}
	}

	nombre := ctx.FormValue("Nombre")

	Cata.ENombreCatalogo.Nombre = nombre
	Catalogo.Nombre = nombre

	if nombre == "" {
		EstatusPeticion = true
		Cata.ENombreCatalogo.IEstatus = true
		Cata.ENombreCatalogo.IMsj = "El Nombre del catálogo es requerido."
	} else if len(nombre) < 5 || len(nombre) > 50 {
		EstatusPeticion = true
		Cata.ENombreCatalogo.IEstatus = true
		Cata.ENombreCatalogo.IMsj = "El Nombre del catálogo debe tener una longitud entre [5 y 50]."
	} else if Catalogo.ConsultaExistenciaByFieldMgo("Nombre", nombre) {
		EstatusPeticion = true
		Cata.ENombreCatalogo.IEstatus = true
		Cata.ENombreCatalogo.IMsj = "El Nombre del catálogo ya existe."
	}

	Valores := ctx.Request.Form["Valores"]

	if Valores == nil {
		EstatusPeticion = true
		Cata.EValoresCatalogo.Valores.EValorValores.IEstatus = true
		Cata.EValoresCatalogo.Valores.EValorValores.IMsj += "El catálogo debe contener al menos un valor, recuerde que debe agregarse con el botón Agregar."
	}

	var val CatalogoModel.ValoresMgo
	var vals []CatalogoModel.ValoresMgo

	valorestmpl := ``

	for _, v := range Valores {
		val = CatalogoModel.ValoresMgo{}
		val.ID = bson.NewObjectId()
		val.Valor = v
		vals = append(vals, val)

		valorestmpl += `<tr>
							<td>
								<input type="hidden" class="form-control" name="ValoresIds" value="` + val.ID.Hex() + `">
								<input type="text" class="form-control" name="Valores" value="` + v + `" readonly>
							</td>
							<td>
								<button type="button" class="btn btn-danger deleteButton">
								<span class="glyphicon glyphicon-trash btn-xs"> </span></button>
							</td>
						</tr>`
	}

	Cata.EValoresCatalogo.Valores.EValorValores.Ihtml = template.HTML(valorestmpl)
	Catalogo.Valores = vals
	Catalogo.ID = bson.NewObjectId()
	Catalogo.Estatus = CatalogoModel.RegresaIDEstatusActivo(132)
	Catalogo.FechaHora = time.Now()
	Catalogo.Clave = CatalogoModel.SiguienteClaveDisponible()
	CataE.Catalogo = Cata

	if EstatusPeticion {
		CataE.SEstado = false
		CataE.SMsj = "Existen errores en su captura de catálogo, no se puede procesar su solicitud de Alta." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("CatalogoAlta.html", CataE)
	} else {

		if Catalogo.InsertaMgo() {

			if Catalogo.InsertaElastic() {
				CataE.SEstado = true
				CataE.SMsj = "El Alta de su catálogo fue exitosa." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
				Cata = CatalogoModel.Catalogo{}
				CataE.Catalogo = Cata
				ctx.Render("CatalogoAlta.html", CataE)
			} else {
				CataE.SEstado = false
				CataE.SMsj = "Ocurrió un error al insertar en elasticSearch" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
				ctx.Render("CatalogoAlta.html", CataE)
			}
		} else {
			CataE.SEstado = false
			CataE.SMsj = "Ocurrió un error al insertar en MongoDb" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
			ctx.Render("CatalogoAlta.html", CataE)
		}
	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Catalogo
func EditaGet(ctx *iris.Context) {

	var SCatalogos CatalogoModel.SCatalogo

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SCatalogos.SSesion.Name = NameUsrLoged
	SCatalogos.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SCatalogos.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SCatalogos.SEstado = false
		SCatalogos.SMsj = errSes.Error()
		ctx.Render("ZError.html", SCatalogos)
		return
	}

	id := ctx.Param("ID")

	if bson.IsObjectIdHex(id) {
		Catalogo := CatalogoModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Catalogo) {
			SCatalogos.Catalogo.ID = bson.ObjectIdHex(id)
			SCatalogos.Catalogo.ENombreCatalogo.Nombre = Catalogo.Nombre
			SCatalogos.Catalogo.EDescripcionCatalogo.Descripcion = Catalogo.Descripcion
			SCatalogos.Catalogo.EClaveCatalogo.Clave = Catalogo.Clave

			valorestmpl := ``
			for _, v := range Catalogo.Valores {
				valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="ValoresIds" value="` + v.ID.Hex() + `">
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Valores" value="` + v.Valor + `" readonly>
									</td>
									<td>
										<button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button>										
									</td>
								</tr>`
			}
			SCatalogos.Catalogo.EValoresCatalogo.Valores.EValorValores.Ihtml = template.HTML(valorestmpl)
			ctx.Render("CatalogoEdita.html", SCatalogos)
		} else {
			Catalogos := CargaCombos.CargaComboCatalogo(0, "")
			SCatalogos.Catalogo.ENombreCatalogo.Ihtml = template.HTML(Catalogos)
			SCatalogos.SEstado = false
			SCatalogos.SMsj = "No se encontró el catálogo en la base, es posible que la conexión halla fallado, vuelva a intentar."
			ctx.Render("CatalogoIndex.html", SCatalogos)
		}
		SCatalogos.SEstado = true
	} else {
		Catalogos := CargaCombos.CargaComboCatalogo(0, "")
		SCatalogos.Catalogo.ENombreCatalogo.Ihtml = template.HTML(Catalogos)
		SCatalogos.SEstado = false
		SCatalogos.SMsj = "No se ha recibido un parámetro adecuado de catálogo para editar, vuelva a intentar"
		ctx.Render("CatalogoIndex.html", SCatalogos)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Catalogo
func EditaPost(ctx *iris.Context) {

	var CataE CatalogoModel.SCatalogo
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	CataE.SSesion.Name = NameUsrLoged
	CataE.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	CataE.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		CataE.SEstado = false
		CataE.SMsj = errSes.Error()
		ctx.Render("ZError.html", CataE)
		return
	}

	EstatusPeticion := false //Comienza suponiendo que hay error

	var Cata CatalogoModel.Catalogo

	var Catalogo CatalogoModel.CatalogoMgo

	clave, _ := strconv.Atoi(ctx.FormValue("Clave"))
	Cata.EClaveCatalogo.Clave = int64(clave)

	descripcion := ctx.FormValue("Descripcion")
	Cata.EDescripcionCatalogo.Descripcion = descripcion
	Catalogo.Descripcion = descripcion

	if descripcion != "" {
		if len(descripcion) < 20 || len(descripcion) > 250 {
			EstatusPeticion = true
			Cata.EDescripcionCatalogo.IEstatus = true
			Cata.EDescripcionCatalogo.IMsj = "La descripció del catálogo debe tener una longitud entre [20 y 250] caracteres."
		}
	}

	nombre := ctx.FormValue("Nombre")

	Cata.ENombreCatalogo.Nombre = nombre
	Catalogo.Nombre = nombre

	if nombre == "" {
		EstatusPeticion = true
		Cata.ENombreCatalogo.IEstatus = true
		Cata.ENombreCatalogo.IMsj = "El Nombre del catálogo es requerido."
	} else if len(nombre) < 5 || len(nombre) > 50 {
		EstatusPeticion = true
		Cata.ENombreCatalogo.IEstatus = true
		Cata.ENombreCatalogo.IMsj = "El Nombre del catálogo debe tener una longitud entre [5 y 50]."
	}

	Valores := ctx.Request.Form["Valores"]
	ValoresIds := ctx.Request.Form["ValoresIds"]

	if Valores == nil {
		EstatusPeticion = true
		Cata.EValoresCatalogo.Valores.EValorValores.IEstatus = true
		Cata.EValoresCatalogo.Valores.EValorValores.IMsj += "El catálogo debe contener al menos un valor."
	}

	var val CatalogoModel.ValoresMgo
	var vals []CatalogoModel.ValoresMgo

	valorestmpl := ``

	for i, v := range Valores {
		val = CatalogoModel.ValoresMgo{}

		if ValoresIds[i] != "" {
			val.ID = bson.ObjectIdHex(ValoresIds[i])
		} else {
			val.ID = bson.NewObjectId()
		}

		val.Valor = v
		vals = append(vals, val)
		valorestmpl += `<tr>
						<td>
							<input type="hidden" class="form-control" name="ValoresIds" value="` + val.ID.Hex() + `">
							<input type="text" class="form-control" name="Valores" value="` + v + `" readonly>
						</td>
							<td>
								<button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button>
							</td>
						</tr>`

	}

	Cata.EValoresCatalogo.Valores.EValorValores.Ihtml = template.HTML(valorestmpl)
	Catalogo.Valores = vals

	CataE.Catalogo = Cata

	if EstatusPeticion {
		CataE.Catalogo.ID = bson.ObjectIdHex(ctx.FormValue("ID"))
		CataE.SEstado = false
		CataE.SMsj = "Existen errores en su captura de catálogo, no se puede procesar su solicitud de Edición."
		ctx.Render("CatalogoEdita.html", CataE)

	} else {

		ide := ctx.FormValue("ID")

		if bson.IsObjectIdHex(ide) {
			Catalogo.ID = bson.ObjectIdHex(ide)
			CataE.Catalogo.ID = bson.ObjectIdHex(ide)

			CataNom := CatalogoModel.GetOne(bson.ObjectIdHex(ide))

			if !MoGeneral.EstaVacio(CataNom) {
				if CataNom.Nombre != nombre {
					if Catalogo.ConsultaExistenciaByFieldMgo("Nombre", nombre) {
						EstatusPeticion = true
						CataE.Catalogo.ENombreCatalogo.IEstatus = true
						CataE.Catalogo.ENombreCatalogo.IMsj = "El Nombre del catálogo ya existe."
						CataE.SEstado = false
						CataE.SMsj = "Existen errores en su captura de catálogo, no se puede procesar su solicitud de Edición."
						ctx.Render("CatalogoEdita.html", CataE)
					}
				}
			} else {
				CataE.SEstado = false
				CataE.SMsj = "Ha ocurrido un problema con la referencia del catálogo vuelva a seleccionar su catálogo a editar, disculpe la molestia."
				Catalogos := CargaCombos.CargaComboCatalogo(0, "")
				CataE.Catalogo.ENombreCatalogo.Ihtml = template.HTML(Catalogos)
				ctx.Render("CatalogoIndex.html", CataE)
			}
			if !EstatusPeticion {
				if Catalogo.ActualizaMgo([]string{"Nombre", "Descripcion", "Valores"}, []interface{}{Catalogo.Nombre, Catalogo.Descripcion, Catalogo.Valores}) {
					errupd := Catalogo.ActualizaElastic()
					if errupd == nil {
						CataE.SEstado = true
						CataE.SMsj = "La Actualización de su catálogo fue exitosa."

						valorestmpl = ``
						for i, v := range Valores {
							val = CatalogoModel.ValoresMgo{}

							if ValoresIds[i] != "" {
								val.ID = bson.ObjectIdHex(ValoresIds[i])
							} else {
								val.ID = bson.NewObjectId()
							}

							val.Valor = v
							vals = append(vals, val)
							valorestmpl += `<tr>
											<td>
												<input type="hidden" class="form-control" name="ValoresIds" value="` + val.ID.Hex() + `">
												<input type="text" class="form-control" name="Valores" value="` + v + `" readonly>
											</td>
										</tr>`
						}

						CataE.Catalogo.EValoresCatalogo.Valores.EValorValores.Ihtml = template.HTML(valorestmpl)
						ctx.MustRender("CatalogoDetalle.html", CataE)
					} else {
						if CataNom.ReemplazaMgo() {
							CataE.SEstado = false
							CataE.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + "). Se ha reestablecido la informacion"
							ctx.Render("CatalogoEdita.html", CataE)
						} else {
							CataE.SEstado = false
							CataE.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + ") No se pudo reestablecer la informacion"
							ctx.Render("CatalogoEdita.html", CataE)
						}
					}
				} else {
					CataE.SEstado = false
					CataE.SMsj = "Ocurrió un problema al actualizar su catálogo en la base de datos, intente de nuevo más tarde."
					ctx.Render("CatalogoEdita.html", CataE)
				}
			}

		} else {
			CataE.SEstado = false
			CataE.SMsj = "Ha ocurrido un problema con la referencia del catálogo vuelva a seleccionar su catálogo a editar, disculpe la molestia."
			Catalogos := CargaCombos.CargaComboCatalogo(0, "")
			CataE.Catalogo.ENombreCatalogo.Ihtml = template.HTML(Catalogos)
			ctx.Render("CatalogoIndex.html", CataE)
		}

	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {

	var SCatalogos CatalogoModel.SCatalogo
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SCatalogos.SSesion.Name = NameUsrLoged
	SCatalogos.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SCatalogos.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SCatalogos.SEstado = false
		SCatalogos.SMsj = errSes.Error()
		ctx.Render("ZError.html", SCatalogos)
		return
	}

	id := ctx.Param("ID")

	if bson.IsObjectIdHex(id) {
		Catalogo := CatalogoModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Catalogo) {
			SCatalogos.Catalogo.ID = bson.ObjectIdHex(id)
			SCatalogos.Catalogo.ENombreCatalogo.Nombre = Catalogo.Nombre
			SCatalogos.Catalogo.EDescripcionCatalogo.Descripcion = Catalogo.Descripcion
			SCatalogos.Catalogo.EClaveCatalogo.Clave = Catalogo.Clave

			valorestmpl := ``
			for _, v := range Catalogo.Valores {
				valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="ValoresIds" value="` + v.ID.Hex() + `">
										<input type="text" class="form-control" name="Valores" value="` + v.Valor + `" readonly>
									</td>									
								</tr>`
			}
			SCatalogos.Catalogo.EValoresCatalogo.Valores.EValorValores.Ihtml = template.HTML(valorestmpl)
			ctx.Render("CatalogoDetalle.html", SCatalogos)
		} else {
			Catalogos := CargaCombos.CargaComboCatalogo(0, "")
			SCatalogos.Catalogo.ENombreCatalogo.Ihtml = template.HTML(Catalogos)
			SCatalogos.SEstado = false
			SCatalogos.SMsj = "No se encontró el catálogo en la base, es posible que la conexión halla fallado, vuelva a intentar."
			ctx.Render("CatalogoIndex.html", SCatalogos)
		}
		SCatalogos.SEstado = true
	} else {
		Catalogos := CargaCombos.CargaComboCatalogo(0, "")
		SCatalogos.Catalogo.ENombreCatalogo.Ihtml = template.HTML(Catalogos)
		SCatalogos.SEstado = false
		SCatalogos.SMsj = "No se ha recibido un parámetro adecuado de catálogo para editar, vuelva a intentar"
		ctx.Render("CatalogoIndex.html", SCatalogos)
	}

}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send CatalogoModel.SCatalogo
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

	ctx.Render("CatalogoDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send CatalogoModel.SCatalogo

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

		Cabecera, Cuerpo := CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrToMongo))
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
	var Send CatalogoModel.SCatalogo
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Catalogo.ENombreCatalogo.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := CatalogoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = CatalogoModel.GeneraTemplatesBusqueda(CatalogoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
