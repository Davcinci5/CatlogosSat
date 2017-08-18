package ProductoControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"time"

	"../../Modulos/Session"

	"../../Modelos/CatalogoModel"
	"../../Modelos/ProductoModel"
	"../../Modelos/UnidadModel"

	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Imagenes"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//##########< Variables Generales > ############
var catalogoTipo = 162
var catalogoEstatusProd = 161

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

//IndexGet renderea al index de Producto
func IndexGet(ctx *iris.Context) {

	var Send ProductoModel.SProducto
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
	numeroRegistros = ProductoModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Productos := ProductoModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Productos {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(Productos[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(Productos[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("ProductoIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Producto
func IndexPost(ctx *iris.Context) {

	var Send ProductoModel.SProducto
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
	//Send.Producto.EVARIABLEProducto.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := ProductoModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("ProductoIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Producto
func AltaGet(ctx *iris.Context) {

	var Send ProductoModel.SProducto
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

	Send.Producto.ETipoProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipo, ""))
	Send.Producto.EUnidadProducto.Ihtml = template.HTML(CargaCombos.CargaComboUnidades(""))
	Send.Producto.EMmvProducto.Mmv = "1"
	ctx.Render("ProductoAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Producto
func AltaPost(ctx *iris.Context) {
	var SProducto ProductoModel.SProducto
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SProducto.SSesion.Name = NameUsrLoged
	SProducto.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SProducto.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SProducto.SEstado = false
		SProducto.SMsj = errSes.Error()
		ctx.Render("ZError.html", SProducto)
		return
	}

	EstatusPeticion := false
	var Producto ProductoModel.ProductoMgo
	Producto.ID = bson.NewObjectId()
	Producto.FechaHora = time.Now()

	//get img form
	file, header, err := ctx.FormFile("Imagenes")
	if err != nil {
		SProducto.Producto.EImagenesProducto.IEstatus = true
		SProducto.Producto.EImagenesProducto.IMsj = "Error al Subir la imagen"
	} else {
		// Insertar la imagen en mongo
		idsImg, err := MoConexion.InsertarImagen(file, header)
		if err != nil {
			EstatusPeticion = true
			SProducto.Producto.EImagenesProducto.IEstatus = true
			SProducto.Producto.EImagenesProducto.IMsj = idsImg
		}
		Producto.Imagenes = append(Producto.Imagenes, bson.ObjectIdHex(idsImg))
		SProducto.Producto.EImagenesProducto.Imagenes = append(Producto.Imagenes, bson.ObjectIdHex(idsImg))
	}
	htmlImagenes, err := Imagenes.CargaTemplateImagenes(Producto.Imagenes)
	fmt.Println(htmlImagenes)
	SProducto.Producto.EImagenesProducto.Ihtml = htmlImagenes

	nombre := ctx.FormValue("Nombre")
	SProducto.Producto.ENombreProducto.Nombre = nombre
	Producto.Nombre = nombre

	if nombre == "" {
		EstatusPeticion = true
		SProducto.Producto.ENombreProducto.IEstatus = true
		SProducto.Producto.ENombreProducto.IMsj = "Campo Descrpción es obligatorio"

	}

	codigos := ctx.Request.Form["Codigos"]
	codigosval := ctx.Request.Form["Valcodigos"]
	if len(codigos) > 0 {
		SProducto.Producto.ECodigosProducto.Ihtml = template.HTML(creaHTMLCodigos(codigos, codigosval))
		Producto.Codigos.Claves = codigos
		Producto.Codigos.Valores = codigosval
	} else {
		EstatusPeticion = true
		SProducto.Producto.ECodigosProducto.IEstatus = true
		SProducto.Producto.ECodigosProducto.IMsj = "Debe contenerse almenos un codigo y valor"
	}

	tipo := ctx.FormValue("Tipo")
	SProducto.Producto.ETipoProducto.Tipo = bson.ObjectIdHex(tipo)
	SProducto.Producto.ETipoProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipo, tipo))
	Producto.Tipo = bson.ObjectIdHex(tipo)
	if tipo == "" {
		EstatusPeticion = true
		SProducto.Producto.ETipoProducto.IEstatus = true
		SProducto.Producto.ETipoProducto.IMsj = "El campo Tipo es obligatorio"
	}

	Mmv := ctx.FormValue("Mmv")
	SProducto.Producto.EMmvProducto.Mmv = Mmv
	m, err := strconv.ParseFloat(Mmv, 64)
	if err != nil {
		EstatusPeticion = true
		Producto.Mmv = 0
		SProducto.Producto.ETipoProducto.IEstatus = true
		SProducto.Producto.ETipoProducto.IMsj = "El campo Tipo es obligatorio"
	}
	Producto.Mmv = m

	// -------- falta imagenes ------

	unidad := ctx.FormValue("Unidad")
	SProducto.Producto.EUnidadProducto.Unidad = bson.ObjectIdHex(unidad)
	SProducto.Producto.EUnidadProducto.Ihtml = template.HTML(CargaCombos.CargaComboUnidades(unidad))
	Producto.Unidad = bson.ObjectIdHex(unidad)
	if unidad == "" {
		EstatusPeticion = true
		SProducto.Producto.EUnidadProducto.IEstatus = true
		SProducto.Producto.EUnidadProducto.IMsj = "Debe seleccionar una unidad"
	}

	ventaFrac := ctx.FormValue("VentaFraccion")
	if ventaFrac != "" {
		SProducto.Producto.EVentaFraccionProducto.Ihtml = template.HTML("checked")
		Producto.VentaFraccion = true
	}

	etiquetas := ctx.Request.Form["Etiquetas"]
	if len(etiquetas) > 0 {
		SProducto.Producto.EEtiquetasProducto.Ihtml = template.HTML(creaHTMLEtiquetas(etiquetas))
		Producto.Etiquetas = etiquetas
	} else {
		EstatusPeticion = true
		SProducto.Producto.EEtiquetasProducto.IEstatus = true
		SProducto.Producto.EEtiquetasProducto.IMsj = "Debe contener al menos una etiqueta"
	}

	SProducto.Producto.EEstatusProducto.Estatus = CatalogoModel.RegresaIDEstatusActivo(catalogoEstatusProd)
	Producto.Estatus = CatalogoModel.RegresaIDEstatusActivo(catalogoEstatusProd)
	fmt.Println(CatalogoModel.RegresaIDEstatusActivo(catalogoEstatusProd))

	// SProducto.Producto.EVentaFraccionProducto.VentaFraccion = ventaFrac
	fmt.Println(Producto)

	if EstatusPeticion {
		SProducto.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SProducto.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("ProductoAlta.html", SProducto)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if Producto.InsertaMgo() {

			if Producto.InsertaElastic() {
				SProducto.SEstado = true
				SProducto.SMsj = "Se ha realizado una inserción exitosa"
				//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
				ctx.Redirect("/Productos/detalle/"+Producto.ID.Hex(), 301)
			} else {
				SProducto.SEstado = false
				SProducto.SMsj = "Ocurrió un error al insertar en elasticSearch"
				ctx.Render("ProductoAlta.html", SProducto)
			}

		} else {
			SProducto.SEstado = false
			SProducto.SMsj = "Ocurrió un error al insertar en MongoDb"
			ctx.Render("ProductoAlta.html", SProducto)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Producto
func EditaGet(ctx *iris.Context) {
	var SProducto ProductoModel.SProducto
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SProducto.SSesion.Name = NameUsrLoged
	SProducto.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SProducto.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SProducto.SEstado = false
		SProducto.SMsj = errSes.Error()
		ctx.Render("ZError.html", SProducto)
		return
	}

	var producto ProductoModel.ProductoMgo
	id := ctx.Param("ID")

	if bson.IsObjectIdHex(id) {
		producto = ProductoModel.GetOne(bson.ObjectIdHex(id))

		if producto.ID.Hex() == "" {
			SProducto.SEstado = false
			SProducto.SMsj = "No se encontr el objeto en la base, verifique su conexion e intente de nuevo"
			ctx.Render("ProductoIndex.html", SProducto)
			return
		}
	} else {
		SProducto.SEstado = false
		SProducto.SMsj = "El parametro que se recibio no es adecuado, intente de nuevo"
		ctx.Render("ProductoIndex.html", SProducto)
		return
	}

	SProducto.Producto.ID = producto.ID
	SProducto.Producto.ENombreProducto.Nombre = producto.Nombre
	SProducto.Producto.ECodigosProducto.Ihtml = template.HTML(creaHTMLCodigos(producto.Codigos.Claves, producto.Codigos.Valores))
	SProducto.Producto.ETipoProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipo, producto.Tipo.Hex()))
	// SProducto.Producto.EImagenesProducto = "nada por ahora"
	SProducto.Producto.EUnidadProducto.Ihtml = template.HTML(CargaCombos.CargaComboUnidades(producto.Unidad.Hex()))
	if producto.VentaFraccion == true {
		SProducto.Producto.EVentaFraccionProducto.Ihtml = template.HTML("checked")
	}
	SProducto.Producto.EEtiquetasProducto.Ihtml = template.HTML(creaHTMLEtiquetas(producto.Etiquetas))
	SProducto.Producto.EEstatusProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoEstatusProd, producto.Estatus.Hex()))
	SProducto.Producto.EFechaHoraProducto.FechaHora = producto.FechaHora

	SProducto.Producto.EMmvProducto.Mmv = strconv.FormatFloat(producto.Mmv, 'f', 0, 64)

	//imagenes
	htmlImagenes, _ := Imagenes.CargaTemplateImagenes(producto.Imagenes)
	SProducto.Producto.EImagenesProducto.Ihtml = htmlImagenes
	ctx.Render("ProductoEdita.html", SProducto)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Producto
func EditaPost(ctx *iris.Context) {

	var SProducto ProductoModel.SProducto
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SProducto.SSesion.Name = NameUsrLoged
	SProducto.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SProducto.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SProducto.SEstado = false
		SProducto.SMsj = errSes.Error()
		ctx.Render("ZError.html", SProducto)
		return
	}

	errorID := ""
	EstatusPeticion := false
	var Producto ProductoModel.ProductoMgo
	var ProductoAnt ProductoModel.ProductoMgo
	//proexistente := ProductoModel.ProductoMgo{}
	Producto.FechaHora = time.Now()

	id := ctx.FormValue("IDname")
	if bson.IsObjectIdHex(id) {
		Producto = ProductoModel.GetOne(bson.ObjectIdHex(id))
		if Producto.ID.Hex() == "" {
			EstatusPeticion = true
			errorID = "Error al buscar el Producto, intente mas tarde"
		}
		ProductoAnt = Producto
		/*
			Producto.ID = bson.ObjectIdHex(id)
			SProducto.Producto.ID = bson.ObjectIdHex(id)
			proexistente = ProductoModel.GetOne(bson.ObjectIdHex(id))
			if proexistente.ID.Hex() == "" {
				EstatusPeticion = true
				errorID = ", Error al buscar el Producto, intente mas tarde"
			}
		*/
	} else {
		EstatusPeticion = true
		errorID = ", Error en la referencia del Producto"
	}
	//Imagenes//get img form
	file, header, err := ctx.FormFile("Imagenes")
	if err != nil {
		SProducto.Producto.EImagenesProducto.IEstatus = true
		SProducto.Producto.EImagenesProducto.IMsj = "Error al Subir la imagen"
	} else {
		// Insertar la imagen en mongo
		idsImg, err := MoConexion.InsertarImagen(file, header)
		if err != nil {
			EstatusPeticion = true
			SProducto.Producto.EImagenesProducto.IEstatus = true
			SProducto.Producto.EImagenesProducto.IMsj = idsImg
		}
		Producto.Imagenes = append(Producto.Imagenes, bson.ObjectIdHex(idsImg))
		SProducto.Producto.EImagenesProducto.Imagenes = append(Producto.Imagenes, bson.ObjectIdHex(idsImg))
	}
	htmlImagenes, err := Imagenes.CargaTemplateImagenes(Producto.Imagenes)
	SProducto.Producto.EImagenesProducto.Ihtml = htmlImagenes

	//	SProducto.Producto = Producto //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.
	nombre := ctx.FormValue("Nombre")
	SProducto.Producto.ENombreProducto.Nombre = nombre
	Producto.Nombre = nombre

	if nombre == "" {
		EstatusPeticion = true
		SProducto.Producto.ENombreProducto.IEstatus = true
		SProducto.Producto.ENombreProducto.IMsj = "Campo Descrpción es obligatorio"

	}

	codigos := ctx.Request.Form["Codigos"]
	codigosval := ctx.Request.Form["Valcodigos"]
	if len(codigos) > 0 {
		SProducto.Producto.ECodigosProducto.Ihtml = template.HTML(creaHTMLCodigos(codigos, codigosval))
		Producto.Codigos.Claves = codigos
		Producto.Codigos.Valores = codigosval
	} else {
		EstatusPeticion = true
		SProducto.Producto.ECodigosProducto.IEstatus = true
		SProducto.Producto.ECodigosProducto.IMsj = "Debe contenerse almenos un codigo y valor"
	}

	tipo := ctx.FormValue("Tipo")
	SProducto.Producto.ETipoProducto.Tipo = bson.ObjectIdHex(tipo)
	SProducto.Producto.ETipoProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipo, tipo))
	Producto.Tipo = bson.ObjectIdHex(tipo)
	if tipo == "" {
		EstatusPeticion = true
		SProducto.Producto.ETipoProducto.IEstatus = true
		SProducto.Producto.ETipoProducto.IMsj = "El campo Tipo es obligatorio"
	}

	unidad := ctx.FormValue("Unidad")
	SProducto.Producto.EUnidadProducto.Unidad = bson.ObjectIdHex(unidad)
	SProducto.Producto.EUnidadProducto.Ihtml = template.HTML(CargaCombos.CargaComboUnidades(unidad))
	Producto.Unidad = bson.ObjectIdHex(unidad)
	if unidad == "" {
		EstatusPeticion = true
		SProducto.Producto.EUnidadProducto.IEstatus = true
		SProducto.Producto.EUnidadProducto.IMsj = "Debe seleccionar una unidad"
	}

	Mmv := ctx.FormValue("Mmv")
	SProducto.Producto.EMmvProducto.Mmv = Mmv
	m, err := strconv.ParseFloat(Mmv, 64)
	if err != nil {
		EstatusPeticion = true
		Producto.Mmv = 0
		SProducto.Producto.ETipoProducto.IEstatus = true
		SProducto.Producto.ETipoProducto.IMsj = "El campo Tipo es obligatorio"
	}
	Producto.Mmv = m

	ventaFrac := ctx.FormValue("VentaFraccion")
	if ventaFrac != "" {
		SProducto.Producto.EVentaFraccionProducto.Ihtml = template.HTML("checked")
		Producto.VentaFraccion = true
	}

	etiquetas := ctx.Request.Form["Etiquetas"]
	if len(etiquetas) > 0 {
		SProducto.Producto.EEtiquetasProducto.Ihtml = template.HTML(creaHTMLEtiquetas(etiquetas))
		Producto.Etiquetas = etiquetas
	} else {
		EstatusPeticion = true
		SProducto.Producto.EEtiquetasProducto.IEstatus = true
		SProducto.Producto.EEtiquetasProducto.IMsj = "Debe contener al menos una etiqueta"
	}
	estatus := ctx.FormValue("Estatus")
	SProducto.Producto.EEstatusProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoEstatusProd, estatus))
	Producto.Estatus = bson.ObjectIdHex(estatus)

	if EstatusPeticion {
		SProducto.SEstado = false                                                                     //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SProducto.SMsj = "La validación indica que el objeto capturado no puede procesarse" + errorID //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("ProductoEdita.html", SProducto)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if Producto.ReemplazaMgo() {
			errupd := Producto.ActualizaElastic()
			if errupd == nil {
				SProducto.SEstado = true
				SProducto.SMsj = "Se ha realizado una inserción exitosa"
				//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
				ctx.Redirect("/Productos/detalle/"+Producto.ID.Hex(), 301)
			} else {
				if ProductoAnt.ReemplazaMgo() {
					SProducto.SEstado = false
					SProducto.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + "). Se ha reestablecido la informacion"
					ctx.Render("ProductoEdita.html", SProducto)
				} else {
					SProducto.SEstado = false
					SProducto.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + ") No se pudo reestablecer la informacion"
					ctx.Render("ProductoEdita.html", SProducto)
				}
			}
		} else {
			SProducto.SEstado = false
			SProducto.SMsj = "Ocurrió un error al Actualizar el Objeto, intente más tarde"
			ctx.Render("ProductoEdita.html", SProducto)
		}

	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {

	var SProducto ProductoModel.SProducto
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SProducto.SSesion.Name = NameUsrLoged
	SProducto.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SProducto.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SProducto.SEstado = false
		SProducto.SMsj = errSes.Error()
		ctx.Render("ZError.html", SProducto)
		return
	}

	id := ctx.Param("ID")
	if bson.IsObjectIdHex(id) {
		producto := ProductoModel.GetOne(bson.ObjectIdHex(id))
		fmt.Println(producto)
		if !MoGeneral.EstaVacio(producto) {
			SProducto.Producto.ID = producto.ID
			SProducto.Producto.ENombreProducto.Nombre = producto.Nombre
			SProducto.Producto.ECodigosProducto.Ihtml = template.HTML(creaHTMLCodigosDetalle(producto.Codigos.Claves, producto.Codigos.Valores))
			SProducto.Producto.ETipoProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoTipo, producto.Tipo.Hex()))
			// SProducto.Producto.EImagenesProducto = "nada por ahora"
			SProducto.Producto.EUnidadProducto.Ihtml = template.HTML(CargaCombos.CargaComboUnidades(producto.Unidad.Hex()))
			if producto.VentaFraccion == true {
				SProducto.Producto.EVentaFraccionProducto.Ihtml = template.HTML("checked")
			}
			SProducto.Producto.EEtiquetasProducto.Ihtml = template.HTML(creaHTMLEtiquetasDetalle(producto.Etiquetas))
			SProducto.Producto.EEstatusProducto.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(catalogoEstatusProd, producto.Estatus.Hex()))
			SProducto.Producto.EFechaHoraProducto.FechaHora = producto.FechaHora
			SProducto.Producto.EMmvProducto.Mmv = strconv.FormatFloat(producto.Mmv, 'f', 0, 64)
			htmlImagenes, _ := Imagenes.CargaTemplateImagenes(producto.Imagenes)
			SProducto.Producto.EImagenesProducto.Ihtml = htmlImagenes
		} else {
			SProducto.SEstado = false
			SProducto.SMsj = "No se encontro el Producto, posiblemente existe un error con la conexion"
		}

	} else {
		SProducto.SEstado = false
		SProducto.SMsj = "La referencia del Producto no es correcta"
	}

	ctx.Render("ProductoDetalle.html", SProducto)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send ProductoModel.SProducto
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

	ctx.Render("ProductoDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

func creaHTMLCodigosDetalle(codigos, codigosval []string) string {
	htmlcodigos := ``
	for i := 0; i < len(codigos); i++ {
		htmlcodigos += `<tr><td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Codigos" value="` + codigos[i] + `" readonly></td>
					<td><input type="text" class="form-control" name="Valcodigos" value="` + codigosval[i] + `" readonly></td>
					</tr>`
	}
	return htmlcodigos
}

//creaHTMLEtiquetas devuelve los trs dormados de las etiquetas
func creaHTMLEtiquetas(etiquetas []string) string {
	htmle := ``
	for i := 0; i < len(etiquetas); i++ {
		htmle += `<tr>
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Etiquetas" value="` + etiquetas[i] + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
	}
	return htmle
}
func creaHTMLEtiquetasDetalle(etiquetas []string) string {
	htmle := ``
	for i := 0; i < len(etiquetas); i++ {
		htmle += `<tr>
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Etiquetas" value="` + etiquetas[i] + `" readonly></td>
				</tr>`
	}
	return htmle
}

//creaHTMLCodigos devuelve los trs formados por las llaves valor
func creaHTMLCodigos(codigos, codigosval []string) string {
	htmlcodigos := ``
	for i := 0; i < len(codigos); i++ {
		htmlcodigos += `<tr><td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Codigos" value="` + codigos[i] + `" readonly></td>
					<td><input type="text" class="form-control" name="Valcodigos" value="` + codigosval[i] + `" readonly></td>
					<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
					</tr>`
	}
	return htmlcodigos
}

/*
//creaHTMLEtiquetas devuelve los trs dormados de las etiquetas
func creaHTMLEtiquetas(etiquetas []string) string {
	htmle := ``
	for i := 0; i < len(etiquetas); i++ {
		htmle += `<tr>
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Etiquetas" value="` + etiquetas[i] + `" readonly></td>
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
				</tr>`
	}
	return htmle
}

func creaHTMLEtiquetasDetalle(etiquetas []string) string {
	htmle := ``
	for i := 0; i < len(etiquetas); i++ {
		htmle += `<tr>
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Etiquetas" value="` + etiquetas[i] + `" readonly></td>
				</tr>`
	}
	return htmle
}

*/

/*
//ConsultarProductos Consulta productos con sus respectivos impuestos
//Para armar la vista que conforma los productos, los precios e impuestos dinamicamente
func ConsultarProductos(ctx *iris.Context) {
	nombreProducto := ctx.FormValue("nombreProducto")

	productos := ProductoModel.BusquedasProductosNombre(nombreProducto)

	//Comienza a crear la fila para mostrar el producto
	etiquetaFila := `<tr id="filaProducto` + productos.ID.Hex() + `">`
	//Etiqueta para la imaben del Producto
	htmlImagenes, _ := Imagenes.CargaTemplateImagenesCompras(productos.Imagenes)
	Imagenes := `<td>` + htmlImagenes + `</td>`

	//Etiqueta para el nombre del producto
	etiquetaNombre := `<td><label>` + productos.Nombre + `</td>`
	//Se forman las etiquetas para los codigos de los productos, de tal forma que se muestren agrupados en forma vertical
	etiquetaCodigos := `<td>`
	for _, value := range productos.Codigos.Valores {
		etiquetaCodigos += `<label>` + value + `</label>`
	}
	etiquetaCodigos += `</td>`
	//Busca el tipo del producto, pasandole como parametro el objectId del tipo
	tipo := CatalogoModel.RegresaNombreSubCatalogo(productos.Tipo)
	if tipo == "" {
		tipo = "Tipo no encontrado"
	}
	//Etiqueta para el tipo de producto
	etiquetaTipo := `<td><label>` + tipo + `</label></td>`
	//Busca el nombre de la unidad correspondiente del producto , pasandole como parametro el objectId de la unidad
	unidad := UnidadModel.RegresaNombreUnidad(productos.Unidad)
	if unidad == "" {
		unidad = "Unidad no encontrado"
	}
	//Etiqueta para la unidad del producto
	etiquetaUnidad := `<td><label>` + unidad + `</label></td>`

	//Busca el estatus del producto, pasandole como parametro el ObjectId del estatus en particular
	estatus := CatalogoModel.RegresaNombreSubCatalogo(productos.Estatus)
	if estatus == "" {
		estatus = "Estatus desconocido"
	}
	//Etiqueta para el estatus del producto
	etiquetaEstatus := `<td><label>` + estatus + `</label></td>`

	//Se forman las etiquetas extras de los productos, para mostrarlo agrupado en un solo campo
	etiquetaExtras := `<td>`
	for _, value := range productos.Etiquetas {
		etiquetaExtras += `<label>` + value + `</label>`
	}
	etiquetaExtras += `</td>`
	etiquetaBoton := `<td><button id="` + productos.ID.Hex() + `" type="button" class="btn btn-info" onClick='leerFilaProducto(this.id)'>
                                <span class="glyphicon glyphicon-ok-sign">Seleccionar</span>
                              </button></td>`
	htmlProducto := etiquetaNombre + etiquetaCodigos + etiquetaTipo + etiquetaUnidad + etiquetaEstatus + etiquetaExtras + Imagenes + etiquetaBoton
	ctx.Writef(etiquetaFila + htmlProducto)
}
*/

//SDataProducto estructura de Productos para la vista
type SDataProducto struct {
	SEstado bool
	SMsj    string
	SIhtml  template.HTML
}

//ConsultarProductos Consulta productos con sus respectivos impuestos
//Para armar la vista que conforma los productos, los precios e impuestos dinamicamente
func ConsultarProductos(ctx *iris.Context) {
	cadenaBusqueda = ctx.FormValue("nombreProducto")
	Send := SDataProducto{}

	if cadenaBusqueda != "" {
		docs, err := ProductoModel.BusquedaElastic(cadenaBusqueda)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Ocurrieron errores: " + err.Error()
		}
		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)
			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
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
			productos := (ProductoModel.GetEspecifics(arrToMongo))
			tabla := GeneraTabla(productos)
			Send.SEstado = true
			Send.SIhtml = template.HTML(tabla)

		} else {
			Send.SEstado = false
			Send.SMsj = "No se encontraron resultados."
		}

	} else {
		Send.SEstado = false
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
	}
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
}

//GeneraTabla genera una tabla utilizada en el modulo de ventas
func GeneraTabla(productos []ProductoModel.ProductoMgo) string {
	tabla := ``
	for _, producto := range productos {
		if bson.IsObjectIdHex(producto.ID.Hex()) {

			htmlProducto := `<tr id="filaProducto` + producto.ID.Hex() + `">`
			htmlProducto += `<td><label>` + producto.Nombre + `</td>`
			htmlProducto += `<td>`
			for _, value := range producto.Codigos.Valores {
				htmlProducto += `<label>` + value + `</label>`
			}
			htmlProducto += `</td>`
			tipo := CatalogoModel.RegresaNombreSubCatalogo(producto.Tipo)
			if tipo == "" {
				tipo = "Tipo no encontrado"
			}
			htmlProducto += `<td><label>` + tipo + `</label></td>`
			unidad := UnidadModel.RegresaNombreUnidad(producto.Unidad)
			if unidad == "" {
				unidad = "Unidad no encontrado"
			}

			htmlProducto += `<td><label>` + unidad + `</label></td>`
			estatus := CatalogoModel.RegresaNombreSubCatalogo(producto.Estatus)
			if estatus == "" {
				estatus = "Estatus desconocido"
			}
			htmlProducto += `<td><label>` + estatus + `</label></td>`
			htmlProducto += `<td>`
			for _, value := range producto.Etiquetas {
				htmlProducto += `<label>` + value + `</label>`
			}
			htmlProducto += `</td>`
			htmlImagenes, _ := Imagenes.CargaTemplateImagenesCompras(producto.Imagenes)

			htmlProducto += `<td>` + htmlImagenes + `</td>`
			htmlProducto += `<td><button id="` + producto.ID.Hex() + `" type="button" class="btn btn-info" onClick='leerFilaProducto(this.id)'>`
			htmlProducto += `<span class="glyphicon glyphicon-ok-sign">Seleccionar</span>`
			htmlProducto += `</button></td>`
			htmlProducto += `</tr>`
			tabla += htmlProducto
		}
	}
	return tabla
}

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send ProductoModel.SProducto

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

		Cabecera, Cuerpo := ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrToMongo))
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
	var Send ProductoModel.SProducto
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}
	cadenaBusqueda = ctx.FormValue("searchbox")
	if cadenaBusqueda != "" {
		docs := ProductoModel.BuscarEnElastic(cadenaBusqueda)
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

			Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = ProductoModel.GeneraTemplatesBusqueda(ProductoModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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

//ConsultarProductosConImpuestosDeAlmacen Consulta productos con sus respectivos impuestos
//Para armar la vista que conforma los productos, los precios e impuestos dinamicamente
func ConsultarProductosConImpuestosDeAlmacen(ctx *iris.Context) {
	cadenaBusqueda = ctx.FormValue("nombreProducto")
	almacen := ctx.FormValue("Almacen")
	Send := ProductoModel.SDataProducto{}

	if cadenaBusqueda != "" {
		docs, err := ProductoModel.BusquedaElastic(cadenaBusqueda)
		if err != nil {
			Send.SEstado = false
			Send.SMsj = "Ocurrieron errores: " + err.Error()
		}
		if docs.Hits.TotalHits > 0 {
			arrIDElastic = []bson.ObjectId{}

			for _, item := range docs.Hits.Hits {
				IDElastic = bson.ObjectIdHex(item.Id)
				arrIDElastic = append(arrIDElastic, IDElastic)
			}

			numeroRegistros = len(arrIDElastic)
			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
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
			productos := (ProductoModel.GetEspecifics(arrToMongo))
			tabla := ProductoModel.GeneraTablaConImpuestos(productos, almacen)
			Send.SEstado = true
			Send.SIhtml = template.HTML(tabla)

		} else {
			Send.SEstado = false
			Send.SMsj = "No se encontraron resultados."
		}

	} else {
		Send.SEstado = false
		Send.SMsj = "No se recibió una cadena de consulta, favor de escribirla."
	}
	jData, _ := json.Marshal(Send)
	ctx.Header().Set("Content-Type", "application/json")
	ctx.Write(jData)
}
