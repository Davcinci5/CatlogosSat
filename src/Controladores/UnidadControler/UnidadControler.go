package UnidadControler

import (
	"fmt"
	"html/template"

	"../../Modelos/CatalogoModel"
	"../../Modelos/UnidadModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/General"
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
var result []UnidadModel.Unidad
var resultPage []UnidadModel.Unidad
var templatePaginacion = ``

//####################< INDEX (BUSQUEDA) >###########################

//Index regresa la peticon post que se hizo desde el index de Unidad
func Index(ctx *iris.Context) {
	var Send UnidadModel.SUnidad
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

	Send.SEstado = true
	Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, ""))
	ctx.Render("UnidadIndex.html", Send)
}

//###########################< ALTA >################################

//AltaGet renderea al alta de Unidad
func AltaGet(ctx *iris.Context) {
	var Send UnidadModel.SUnidad
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

	Send.SEstado = true
	Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, ""))
	ctx.Render("UnidadAlta.html", Send)
}

//AltaPost regresa la petición post que se hizo desde el alta de Unidad
func AltaPost(ctx *iris.Context) {
	var Send UnidadModel.SUnidad
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

	Send.SEstado = true
	Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, ""))
	ctx.Render("UnidadAlta.html", Send)
}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Unidad
func EditaGet(ctx *iris.Context) {
	var Send UnidadModel.SUnidad
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

	id := ctx.Param("ID")
	fmt.Println(id)
	if bson.IsObjectIdHex(id) {
		Unidad := UnidadModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Unidad) {
			Send.Unidad.ID = Unidad.ID
			Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
			Send.Unidad.EDescripcionUnidad.Descripcion = Unidad.Descripcion

			valorestmpl := ``
			for _, v := range Unidad.Datos {
				valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="` + v.ID.Hex() + `">
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Nombres" value="` + v.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Abreviaturas" value="` + v.Abreviatura + `" readonly>
									</td>
									<td>
										<button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button>	
									</td>
								</tr>`
			}

			Send.Unidad.EDatosUnidad.Ihtml = template.HTML(valorestmpl)
			ctx.Render("UnidadEdita.html", Send)

		} else {
			Send.Unidad.ID = Unidad.ID
			Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
			ctx.Render("UnidadEdita.html", Send)
		}
	} else {
		Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar la Magnitud, intente de nuevo."
		ctx.Render("UnidadIndex.html", Send)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Unidad
func EditaPost(ctx *iris.Context) {
	EstatusPeticion := false //Comienza suponiendo que no hay error

	var Send UnidadModel.SUnidad
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

	var AuxUnidad UnidadModel.Unidad
	var UnidadM UnidadModel.UnidadMgo

	id := ctx.Param("ID")
	if bson.IsObjectIdHex(id) {
		UnidadG := UnidadModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(UnidadG) {
			// EDITAMOS y utilizamos UnidadG
			AuxUnidad.ID = UnidadG.ID
			AuxUnidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))

			descripcion := ctx.FormValue("Descripcion")
			AuxUnidad.EDescripcionUnidad.Descripcion = descripcion
			UnidadG.Descripcion = descripcion

			Nombres := ctx.Request.Form["Nombres"]
			Abreviaturas := ctx.Request.Form["Abreviaturas"]
			DataIds := ctx.Request.Form["DatosIds"]

			if Nombres == nil {
				EstatusPeticion = true
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IMsj = "Debe crear al menos una unidad para esta magnitud si desea guardar."
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IEstatus = true
			}

			if len(Nombres) == len(Abreviaturas) {
				var val UnidadModel.DataUnidadMgo
				var vals []UnidadModel.DataUnidadMgo

				valorestmpl := ``

				for i, v := range Nombres {
					val = UnidadModel.DataUnidadMgo{}

					if DataIds[i] != "" {
						val.ID = bson.ObjectIdHex(DataIds[i])
						valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="` + val.ID.Hex() + `">
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Nombres" value="` + val.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" class="form-control" name="Abreviaturas" value="` + val.Abreviatura + `" readonly>
									</td>
									<td>
										<button type="button" onblur="ValidaCampo2(this)" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button>	
									</td>
								</tr>`
					} else {
						val.ID = bson.NewObjectId()

						valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="">
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Nombres" value="` + val.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Abreviaturas" value="` + val.Abreviatura + `" readonly>
									</td>
									<td>
										<button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button>	
										<button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"> </span></button>
									</td>
								</tr>`
					}

					val.Nombre = v
					val.Abreviatura = Abreviaturas[i]

					vals = append(vals, val)
				}

				AuxUnidad.EDatosUnidad.Ihtml = template.HTML(valorestmpl)
				UnidadG.Datos = vals

				if EstatusPeticion {
					Send.SEstado = false
					Send.SMsj = "Existen errores en su captura de catálogo, no se puede procesar su solicitud de Alta."
					ctx.Render("UnidadEdita.html", Send)
				} else {
					if UnidadG.ActualizaMgo([]string{"Magnitud", "Descripcion", "Datos"}, []interface{}{UnidadG.Magnitud, UnidadG.Descripcion, UnidadG.Datos}) {

						Send.SEstado = true
						Send.SMsj = "Actualización de Unidad exitosa."
						Unidad := UnidadModel.GetOne(bson.ObjectIdHex(id))
						Send.Unidad.ID = Unidad.ID
						Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
						Send.Unidad.EDescripcionUnidad.Descripcion = Unidad.Descripcion

						valorestmpl := ``
						for _, v := range Unidad.Datos {
							valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="` + v.ID.Hex() + `">
										<input type="text" class="form-control" name="Nombres" value="` + v.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" class="form-control" name="Abreviaturas" value="` + v.Abreviatura + `" readonly>
									</td>
								</tr>`
						}

						Send.Unidad.EDatosUnidad.Ihtml = template.HTML(valorestmpl)
						ctx.Render("UnidadDetalle.html", Send)
					} else {
						Send.SEstado = false
						Send.SMsj = "Ocurrió un problema al dar de alta su catálogo, intente de nuevo más tarde." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
						ctx.Render("UnidadEdita.html", Send)
					}
				}

			} else {
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IMsj = "Ocurrió un error al leer las unidades y procesarlas del lado del servidor, intente de nuevo, disculpe la molestia."
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IEstatus = true
				Send.SEstado = false
				Send.SMsj = "Existen errores en su captura de catálogo, no se puede procesar su solicitud de Alta." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
				ctx.Render("UnidadEdita.html", Send)
			}

		} else {
			//DAMOS DE ALTA y utilizamos UnidadM

			AuxUnidad.ID = bson.ObjectIdHex(id)
			UnidadM.ID = bson.ObjectIdHex(id)
			nombremagnitud := UnidadModel.GetNameMagnitud(id)
			AuxUnidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
			UnidadM.Magnitud = nombremagnitud

			descripcion := ctx.FormValue("Descripcion")
			AuxUnidad.EDescripcionUnidad.Descripcion = descripcion
			UnidadM.Descripcion = descripcion

			//###Validamos etiquetas de unidades######

			Nombres := ctx.Request.Form["Nombres"]
			Abreviaturas := ctx.Request.Form["Abreviaturas"]

			if Nombres == nil {
				EstatusPeticion = true
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IMsj = "Debe crear al menos una unidad para esta magnitud si desea guardar."
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IEstatus = true
			}

			if len(Nombres) == len(Abreviaturas) {
				var val UnidadModel.DataUnidadMgo
				var vals []UnidadModel.DataUnidadMgo

				valorestmpl := ``

				for i, v := range Nombres {
					val = UnidadModel.DataUnidadMgo{}

					val.ID = bson.NewObjectId()
					val.Nombre = v
					val.Abreviatura = Abreviaturas[i]

					vals = append(vals, val)

					valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="">
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Nombres" value="` + val.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" onblur="ValidaCampo2(this)" class="form-control" name="Abreviaturas" value="` + val.Abreviatura + `" readonly>
									</td>
									<td>
										<button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button>	
										<button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"> </span></button>
									</td>
								</tr>`
				}

				AuxUnidad.EDatosUnidad.Ihtml = template.HTML(valorestmpl)
				UnidadM.Datos = vals

				Send.Unidad = AuxUnidad

				if EstatusPeticion {
					Send.SEstado = false
					Send.SMsj = "Existen errores en su captura de catálogo, no se puede procesar su solicitud de Alta." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
					ctx.Render("UnidadEdita.html", Send)
				} else {
					if UnidadM.InsertaMgo() {
						var Send UnidadModel.SUnidad
						Send.SEstado = true
						Send.SMsj = "Creación de Unidades Exitosa."
						Unidad := UnidadModel.GetOne(bson.ObjectIdHex(id))
						Send.Unidad.ID = Unidad.ID
						Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
						Send.Unidad.EDescripcionUnidad.Descripcion = Unidad.Descripcion

						valorestmpl := ``
						for _, v := range Unidad.Datos {
							valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="` + v.ID.Hex() + `">
										<input type="text" class="form-control" name="Nombres" value="` + v.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" class="form-control" name="Abreviaturas" value="` + v.Abreviatura + `" readonly>
									</td>
								</tr>`
						}

						Send.Unidad.EDatosUnidad.Ihtml = template.HTML(valorestmpl)
						ctx.Render("UnidadDetalle.html", Send)
					} else {
						Send.SEstado = false
						Send.SMsj = "Ocurrió un problema al dar de alta su catálogo, intente de nuevo más tarde." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
						ctx.Render("UnidadEdita.html", Send)
					}
				}

			} else {
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IMsj = "Ocurrió un error al leer las unidades y procesarlas del lado del servidor, intente de nuevo, disculpe la molestia."
				AuxUnidad.EDatosUnidad.Datos.ENombreDatos.IEstatus = true
				Send.SEstado = false
				Send.SMsj = "Existen errores en su captura de catálogo, no se puede procesar su solicitud de Alta." //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
				ctx.Render("UnidadEdita.html", Send)
			}
		}
	}

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send UnidadModel.SUnidad
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

	id := ctx.Param("ID")
	if bson.IsObjectIdHex(id) {
		Unidad := UnidadModel.GetOne(bson.ObjectIdHex(id))
		if !MoGeneral.EstaVacio(Unidad) {
			Send.Unidad.ID = Unidad.ID
			Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
			Send.Unidad.EDescripcionUnidad.Descripcion = Unidad.Descripcion

			valorestmpl := ``
			for _, v := range Unidad.Datos {
				valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="` + v.ID.Hex() + `">
										<input type="text" class="form-control" name="Nombres" value="` + v.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" class="form-control" name="Abreviaturas" value="` + v.Abreviatura + `" readonly>
									</td>
								</tr>`
			}

			Send.Unidad.EDatosUnidad.Ihtml = template.HTML(valorestmpl)
			ctx.Render("UnidadDetalle.html", Send)

		} else {
			Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
			Send.SEstado = false
			Send.SMsj = "No se encontró el registro en la base de datos, edite la magnitud para agregar unidades o intente más tarde."
			ctx.Render("UnidadIndex.html", Send)
		}
	} else {
		Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, id))
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar la Magnitud, intente de nuevo."
		ctx.Render("UnidadDetalle.html", Send)
	}
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send UnidadModel.SUnidad
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

	ctx.Render("UnidadDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//ConsultaUnidadesDeMagnitud consulta Las unidades de determinada
//magnitud y les da formato para mostrar en la vista
func ConsultaUnidadesDeMagnitud(ctx *iris.Context) {
	var Send UnidadModel.SUnidad
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

	id := ctx.FormValue("magnitud")

	if bson.IsObjectIdHex(id) {

		Unidad := UnidadModel.GetSubByField("Magnitud", "_id", bson.ObjectIdHex(id))

		if !MoGeneral.EstaVacio(Unidad) {
			Send.Unidad.ID = Unidad.ID
			Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(UnidadModel.CargaComboUnidades(id))
			Send.Unidad.EDescripcionUnidad.Descripcion = Unidad.Descripcion

			valorestmpl := ``
			for _, v := range Unidad.Datos {
				valorestmpl += `<tr>
									<td>
										<input type="hidden" class="form-control" name="DatosIds" value="` + v.ID.Hex() + `">
										<input type="text" class="form-control" name="Nombres" value="` + v.Nombre + `" readonly>
									</td>
									<td>
										<input type="text" class="form-control" name="Abreviaturas" value="` + v.Abreviatura + `" readonly>
									</td>
									<td>
									    <button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button>	
										<button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"> </span></button>
									</td>
								</tr>`
			}

			Send.Unidad.EDatosUnidad.Ihtml = template.HTML(valorestmpl)
			ctx.Writef(valorestmpl)
		}

	} else {

	}

}

//AgregaMagnitudDesdeUnidades esta funcion recibe el
//nombre de la nueva Magnitud a agregar en el catálogo de magnitudes
func AgregaMagnitudDesdeUnidades(ctx *iris.Context) {
	fmt.Println(ctx.Request.Method)
	var Send UnidadModel.SUnidad
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

	nombre := ctx.FormValue("ModalMagnitud")

	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(166))

	if nombre == "" {
		Send.SEstado = false
		Send.SMsj = "No existe un parámetro adecuado para el alta de la nueva magnitud, favor de intentar de nuevo."
		Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, ""))
		ctx.Render("UnidadIndex.html", Send)

	} else if Catalogo.ConsultaExistenciaByFieldMgo("Valores.Valor", nombre) {
		Send.SEstado = false
		Send.SMsj = "El nombre de la Magnitud ya existe, No se pudo dar de alta."
		Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, ""))
		ctx.Render("UnidadIndex.html", Send)

	} else {

		var valor CatalogoModel.ValoresMgo
		valor.ID = bson.NewObjectId()
		valor.Valor = nombre
		Catalogo.Valores = append(Catalogo.Valores, valor)

		if Catalogo.ActualizaMgo([]string{"Valores"}, []interface{}{Catalogo.Valores}) {
			Send.SEstado = true
			Send.SMsj = "Alta de Magnitud exitosa."
			Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, ""))
			ctx.Render("UnidadIndex.html", Send)
		} else {
			Send.SEstado = false
			Send.SMsj = "Falló el intento por dar de alta la nueva magnitud, favor de intenter más tarde."
			Send.Unidad.EMagnitudUnidad.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(166, ""))
			ctx.Render("UnidadIndex.html", Send)
		}
	}

}
