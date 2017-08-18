package ConexionControler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"../../Modulos/Session"

	"../../Modelos/ConexionModel"
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

//IndexGet renderea al index de Conexion
func IndexGet(ctx *iris.Context) {

	var Send ConexionModel.SConexion

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
	numeroRegistros = ConexionModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Conexions := ConexionModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Conexions {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(Conexions[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(Conexions[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("ConexionIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Conexion
func IndexPost(ctx *iris.Context) {

	var Send ConexionModel.SConexion
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
	//Send.Conexion.EVARIABLEConexion.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := ConexionModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("ConexionIndex.html", Send)

}

//###########################< TEST >#############################

//TestConexion sirve para aser un test de conexion
func TestConexion(ctx *iris.Context) {
	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = ctx.FormValue("Servidor")
	paramConex.Usuario = ctx.FormValue("UsuarioBD")
	paramConex.Pass = ctx.FormValue("PassBD")
	paramConex.NombreBase = ctx.FormValue("NombreBD")
	Conexionn, err := MoConexion.ConexioServidorAlmacenPing(paramConex)

	var Send ConexionModel.TestConReturn

	if Conexionn == false {

		Send.Estatus = false
		Send.Mensage = "Error al realizar la conexion"
		Send.MensageError = err.Error()
		err := ctx.JSON(200, Send)
		if err != nil {
			fmt.Println(err)
		}
		// jData, _ := json.Marshal(Send)
		// ctx.Header().Set("Content-Type", "application/json")
		// ctx.Write(jData)
	} else {
		Send.Estatus = true
		Send.Mensage = "conexion realizada con exito"
		err := ctx.JSON(200, Send)
		if err != nil {
			fmt.Println(err)
		}
		//jData, errjSON := json.Marshal(Send)
		//ctx.Header().Set("Content-Type", "application/json")
		//ctx.Write(jData)
	}
}

//###########################< ALTA >################################

//AltaGet renderea al alta de Conexion
func AltaGet(ctx *iris.Context) {

	var Send ConexionModel.SConexion

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

	ctx.Render("ConexionAlta.html", Send)

}

//AltaPost regresa la petición post que se hizo desde el alta de Conexion
func AltaPost(ctx *iris.Context) {

	var Send ConexionModel.SConexion
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
	var Conexion ConexionModel.ConexionMgo
	Conexion.Nombre = ctx.FormValue("Nombre")
	Conexion.Servidor = ctx.FormValue("Servidor")
	Conexion.NombreBD = ctx.FormValue("NombreBD")
	Conexion.UsuarioBD = ctx.FormValue("UsuarioBD")
	Conexion.PassBD = ctx.FormValue("PassBD")
	ID := bson.NewObjectId()
	Conexion.ID = ID

	Send.Conexion.ID = Conexion.ID
	Send.Conexion.ENombreConexion.Nombre = Conexion.Nombre
	Send.Conexion.EServidorConexion.Servidor = Conexion.Servidor
	Send.Conexion.ENombreBDConexion.NombreBD = Conexion.NombreBD
	Send.Conexion.EUsuarioBDConexion.UsuarioBD = Conexion.UsuarioBD
	Send.Conexion.EPassBDConexion.PassBD = Conexion.PassBD

	if Conexion.InsertaMgo() {
		if Conexion.InsertaElastic() {
			Send.SEstado = true
			Send.SMsj = "Se ha realizado una inserción exitosa"
			ctx.Render("ConexionDetalle.html", Send)
		} else {
			if Conexion.EliminaByIDMgo() {
				Send.SEstado = false
				Send.SMsj = "Falló la inserción de la conexion, se ha eliminado los registros"
				ctx.Render("ConexionAlta.html", Send)
			}
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "Ocurrió un error al insertar La Conexion, intente más tarde"
		ctx.Render("ConexionAlta.html", Send)
	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Conexion
func EditaGet(ctx *iris.Context) {

	var Send ConexionModel.SConexion
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
	idCo := ctx.Param("ID")
	if bson.IsObjectIdHex(idCo) {
		Conexion := ConexionModel.GetOne(bson.ObjectIdHex(idCo))
		if !MoGeneral.EstaVacio(Conexion) {

			Send.SEstado = true
			//Send.SMsj = "Se ha realizado
			Send.Conexion.ID = Conexion.ID
			Send.Conexion.ENombreConexion.Nombre = Conexion.Nombre
			Send.Conexion.EServidorConexion.Servidor = Conexion.Servidor
			Send.Conexion.ENombreBDConexion.NombreBD = Conexion.NombreBD
			Send.Conexion.EUsuarioBDConexion.UsuarioBD = Conexion.UsuarioBD
			Send.Conexion.EPassBDConexion.PassBD = Conexion.PassBD
			ctx.Render("ConexionEdita.html", Send)
		} else {
			Send.SEstado = false
			Send.SMsj = "La conexion no se a encontrado, intente de nuevo."
			ctx.Render("ConexionIndex.html", Send)
		}
	} else {
		fmt.Println("No es un id")
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar La conexion, intente de nuevo."
		ctx.Render("ConexionIndex.html", Send)
	}

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Conexion
func EditaPost(ctx *iris.Context) {

	var Send ConexionModel.SConexion

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
	idCo := ctx.Param("ID")
	if bson.IsObjectIdHex(idCo) {
		Conexion := ConexionModel.GetOne(bson.ObjectIdHex(idCo))
		ConexionExis := Conexion
		if !MoGeneral.EstaVacio(Conexion) {
			Conexion.Nombre = ctx.FormValue("Nombre")
			Conexion.Servidor = ctx.FormValue("Servidor")
			Conexion.NombreBD = ctx.FormValue("NombreBD")
			Conexion.UsuarioBD = ctx.FormValue("UsuarioBD")
			Conexion.PassBD = ctx.FormValue("PassBD")

			Send.Conexion.ID = Conexion.ID
			Send.Conexion.ENombreConexion.Nombre = Conexion.Nombre
			Send.Conexion.EServidorConexion.Servidor = Conexion.Servidor
			Send.Conexion.ENombreBDConexion.NombreBD = Conexion.NombreBD
			Send.Conexion.EUsuarioBDConexion.UsuarioBD = Conexion.UsuarioBD
			Send.Conexion.EPassBDConexion.PassBD = Conexion.PassBD

			actualizado := Conexion.ActualizaMgo([]string{"Nombre", "Servidor", "NombreBD", "UsuarioBD", "PassBD"}, []interface{}{Conexion.Nombre, Conexion.Servidor, Conexion.NombreBD, Conexion.UsuarioBD, Conexion.PassBD})
			if actualizado {
				errupd := Conexion.ActualizaElastic()
				if errupd == nil {
					Send.SEstado = true
					Send.SMsj = "Se ha actualizado la conexion"
					ctx.Render("ConexionDetalle.html", Send)
				} else {
					if ConexionExis.ReemplazaMgo() {
						Send.SEstado = false
						Send.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + "). Se ha reestablecido la informacion"
						ctx.Render("ConexionEdita.html", Send)
					} else {
						Send.SEstado = false
						Send.SMsj = "Ocurrió el siguiente error al actualizar su catálogo: (" + errupd.Error() + ") No se pudo reestablecer la informacion"
						ctx.Render("ConexionEdita.html", Send)
					}
				}
			}
		} else {
			Send.SEstado = false
			Send.SMsj = "El Equipo Caja no se a encontrado, intente de nuevo."
			ctx.Render("ConexionIndex.html", Send)
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
	var Send ConexionModel.SConexion

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
	idCo := ctx.Param("ID")
	if bson.IsObjectIdHex(idCo) {
		Conexion := ConexionModel.GetOne(bson.ObjectIdHex(idCo))
		if !MoGeneral.EstaVacio(Conexion) {
			Send.SEstado = true
			//Send.SMsj = "Se ha realizado
			Send.Conexion.ID = Conexion.ID
			Send.Conexion.ENombreConexion.Nombre = Conexion.Nombre
			Send.Conexion.EServidorConexion.Servidor = Conexion.Servidor
			Send.Conexion.ENombreBDConexion.NombreBD = Conexion.NombreBD
			Send.Conexion.EUsuarioBDConexion.UsuarioBD = Conexion.UsuarioBD
			Send.Conexion.EPassBDConexion.PassBD = Conexion.PassBD
			ctx.Render("ConexionDetalle.html", Send)
		} else {
			Send.SEstado = false
			Send.SMsj = "La conexion no se a encontrado, intente de nuevo."
			ctx.Render("ConexionDetalle.html", Send)
		}
	} else {
		Send.SEstado = false
		Send.SMsj = "No se ha recibido un parámetro adecuado para poder editar La conexion, intente de nuevo."
		ctx.Render("ConexionIndex.html", Send)
	}
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send ConexionModel.SConexion
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

	ctx.Render("ConexionDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send ConexionModel.SConexion

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

		Cabecera, Cuerpo := ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrToMongo))
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
	var Send ConexionModel.SConexion
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Conexion.ENombreConexion.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := ConexionModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = ConexionModel.GeneraTemplatesBusqueda(ConexionModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
