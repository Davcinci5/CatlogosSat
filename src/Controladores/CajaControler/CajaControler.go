package CajaControler

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"time"

	"../../Modelos/PuntoVentaModel"
	"../../Modulos/Session"
	"github.com/jung-kurt/gofpdf"
	//"github.com/leekchan/accounting"

	"../../Modelos/EquipoCajaModel"

	"../../Modelos/MediosPagoModel"

	"../../Modelos/CajaModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	iris "gopkg.in/kataras/iris.v6"
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

//IndexGet renderea al index de Caja
func IndexGet(ctx *iris.Context) {

	var Send CajaModel.SCaja

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
	numeroRegistros = CajaModel.CountAll()
	paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
	Cajas := CajaModel.GetAll()

	arrIDMgo = []bson.ObjectId{}
	for _, v := range Cajas {
		arrIDMgo = append(arrIDMgo, v.ID)
	}
	arrIDElastic = arrIDMgo

	if numeroRegistros <= limitePorPagina {
		Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(Cajas[0:numeroRegistros])
	} else if numeroRegistros >= limitePorPagina {
		Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(Cajas[0:limitePorPagina])
	}

	Send.SIndex.SCabecera = template.HTML(Cabecera)
	Send.SIndex.SBody = template.HTML(Cuerpo)
	Send.SIndex.SGrupo = template.HTML(CargaCombos.CargaComboMostrarEnIndex(limitePorPagina))
	Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
	Send.SIndex.SPaginacion = template.HTML(Paginacion)
	Send.SIndex.SResultados = true

	ctx.Render("CajaIndex.html", Send)

}

//IndexPost regresa la peticon post que se hizo desde el index de Caja
func IndexPost(ctx *iris.Context) {

	var Send CajaModel.SCaja

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
	//Send.Caja.EVARIABLECaja.VARIABLE = cadenaBusqueda    //Variable a autilizar para regresar la cadena de búsqueda.

	if cadenaBusqueda != "" {

		docs := CajaModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo := CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {
			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
	ctx.Render("CajaIndex.html", Send)

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Caja
func AltaGet(ctx *iris.Context) {
	var scaja CajaModel.SCaja

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	scaja.SSesion.Name = NameUsrLoged
	scaja.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	scaja.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		scaja.SEstado = false
		scaja.SMsj = errSes.Error()
		ctx.Render("ZError.html", scaja)
		return
	}

	//Cargo combo de estatus de caja
	var cajac CajaModel.SCaja
	cajac.Caja.EEstatusCaja.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(135, ""))
	//Leo usuario logueado para buscar si tiene caja abierta.
	//Id de caja a buscar para prueba
	//ID := bson.ObjectIdHex("58e65fdca137ae32dc4e317c")
	//Caja a buscar
	//cajaID := bson.ObjectIdHex("58e65fdca137ae32dc4e3179")
	//Usuario a buscar
	usuarioID := bson.ObjectIdHex("58e64924a137ae2a3461d995")
	//Buscamos cajas para este usuario
	//caja := CajaModel.GetOne(ID)
	caja := CajaModel.GetCajaAbiertaByUsuario(usuarioID)
	vcaja := CajaModel.Caja{}
	vcaja.EUsuarioCaja.Usuario = caja.Usuario
	vcaja.ECajaCaja.Caja = caja.Caja
	vcaja.ESaldoCaja.Saldo = caja.Saldo
	scaja.Caja = vcaja
	scaja.Caja.EEstatusCaja.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(135, ""))
	//fmt.Println("Hola,", cajaID, usuarioID, ID)
	//Voy por las operaciones pendientes de pago, estatus 0, estructura temporal que hice
	//operaciones := CajaModel.GetOperacionesAll()
	//Estructura real de modelo POS
	operaciones := PuntoVentaModel.GetAll()
	//Voy a postgres por los movimientos, este metodo es del POS
	html, productos := PuntoVentaModel.ConsultaDatosOperacionAll()
	//Validamos si la estructura viene con datos o vacía.
	vacio := MoGeneral.EstaVacio(caja)
	if vacio == false {
		scaja.SEstado = true
		scaja.SMsj = "Tu tienes una caja abierta."
		fmt.Println("Existe caja abierta, mando vista de operaciones.")
		//ctx.Render("CajaOperaciones.html", scaja)
		ctx.Render("CajaOperaciones.html", map[string]interface{}{
			"scaja":       scaja,
			"operaciones": operaciones,
			"productos":   productos,
			"html":        template.HTML(html),
		})
	} else {
		scaja.SEstado = false
		scaja.SMsj = "Debes ABRIR tu caja."
		fmt.Println("NO existe caja abierta, mando vista para abrir caja.")
		ctx.Render("CajaAlta.html", scaja)
	}

}

//AltaPost regresa la petición post que se hizo desde el alta de Caja
func AltaPost(ctx *iris.Context) {
	var SCaja CajaModel.SCaja
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SCaja.SSesion.Name = NameUsrLoged
	SCaja.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SCaja.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SCaja.SEstado = false
		SCaja.SMsj = errSes.Error()
		ctx.Render("ZError.html", SCaja)
		return
	}

	//######### LEE TU OBJETO DEL FORMULARIO #########
	var Caja CajaModel.CajaMgo
	//Dejamos de leer todo el formulario
	//ctx.ReadForm(&Caja)

	//Y leemos cada campo por separado.
	usuario := bson.ObjectIdHex(ctx.FormValue("Usuario")) //bson.NewObjectId() //ctx.FormValue("Usuario")
	cajax := bson.ObjectIdHex(ctx.FormValue("Caja"))      //ctx.FormValue("Caja")
	cargo, err := strconv.ParseFloat(ctx.FormValue("Cargo"), 64)
	abono := 0.00
	saldo := cargo
	fmt.Println(usuario, cajax, cargo)
	fmt.Println(err)
	if err == nil {
		fmt.Println("Información correcta.")
	} else {
		fmt.Println("Dato inválido.")
	}
	operacion := bson.NewObjectId()
	estatus := bson.ObjectIdHex(ctx.FormValue("Estatus"))
	fechaHora := time.Now()

	//SCaja.Caja.ECargoCaja.Cargo = cargo
	//Lo asignamos a la estructura para guardarlo
	Caja.Usuario = usuario
	Caja.Caja = cajax
	Caja.Cargo = cargo
	Caja.Abono = abono
	Caja.Saldo = saldo
	Caja.Operacion = operacion
	Caja.Estatus = estatus
	Caja.FechaHora = fechaHora

	//fmt.Println(usuario, caja, cargo, operacion, estatus, fechaHora)

	//Caja = Caja{usuario, caja, cargo, cabono, saldo, fechaHora}
	//Lee como arreglo
	//data = ctx.FormValues()
	//fmt.Println(data["Cargo"])
	//######### VALIDA TU OBJETO #########
	EstatusPeticion := false //True indica que hay un error
	//##### TERMINA TU VALIDACION ########

	//########## Asigna vairables a la estructura que enviarás a la vista
	ID := bson.NewObjectId()
	Caja.ID = ID

	//######### ENVIA TUS RESULTADOS #########
	//Asigno el valor introducido para pintarlo en la vista
	SCaja.Caja.ECargoCaja.Cargo = cargo
	SCaja.Caja.EAbonoCaja.Abono = abono
	SCaja.Caja.ESaldoCaja.Saldo = saldo
	SCaja.Caja.EUsuarioCaja.Usuario = usuario
	SCaja.Caja.ECajaCaja.Caja = cajax
	SCaja.Caja.EOperacionCaja.Operacion = operacion

	//SCaja.Caja.EEstatusCaja.Estatus = estatus
	SCaja.Caja.EEstatusCaja.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(135, ctx.FormValue("Estatus")))
	SCaja.Caja.EFechaHoraCaja.FechaHora = fechaHora
	//Valido, si el campo CARGO viene vacio
	if err != nil {

		EstatusPeticion = true
		SCaja.Caja.ECargoCaja.IEstatus = true
		SCaja.Caja.ECargoCaja.IMsj = "Campo DOTACION no puede ser nulo"

	}
	//Valido Usuario
	//Valido Caja
	fmt.Println("cajx:", cajax)
	//	SCaja.Caja = Caja //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.

	if EstatusPeticion {
		SCaja.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SCaja.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("CajaAlta.html", SCaja)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if Caja.InsertaMgo() {
			SCaja.SEstado = true
			SCaja.SMsj = "Se ha realizado una inserción exitosa"
			x := (ID.Hex())
			Ticketpdf(x)
			//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
			ctx.Render("CajaDetalle.html", SCaja)

		} else {
			SCaja.SEstado = false
			SCaja.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("CajaAlta.html", SCaja)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Caja
func EditaGet(ctx *iris.Context) {
	var Send CajaModel.SCaja
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

	ctx.Render("CajaEdita.html", Send)

}

//EditaPost regresa el resultado de la petición post generada desde la edición de Caja
func EditaPost(ctx *iris.Context) {

	var Send CajaModel.SCaja

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

	ctx.Render("CajaEdita.html", Send)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	var Send CajaModel.SCaja
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

	ctx.Render("CajaDetalle.html", Send)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	var Send CajaModel.SCaja
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

	ctx.Render("CajaDetalle.html", Send)
}

//####################< RUTINAS ADICIONALES >##########################

//CutPrice funcion para formatear un numero.
func CutPrice(num float64) float64 {
	str := strconv.FormatFloat(num, 'f', 2, 64)
	i, _ := strconv.ParseFloat(str, 64)
	return i
}

//BuscaDocumentoPost busca un documento
func BuscaDocumentoPost(ctx *iris.Context) {
	fmt.Println("Voy por el documento")

	ID := ctx.FormValue("ID") //"58e67cdda137ae2600adaea2"
	fmt.Println("ID: ", ID)
	if bson.IsObjectIdHex(ID) {
		fmt.Println("ID valido")
	} else {
		fmt.Println("ID NO valido")
	}
	//Busquea anterior
	caja := CajaModel.GetOperacionesByID(bson.ObjectIdHex(ID))
	html, calculadora, total, idOperacion := PuntoVentaModel.ConsultaDatosOperacionByID(ID)
	formasdepago := MediosPagoModel.GetAll() //CajaModel.GetAllFormas()
	//fmt.Println(formasdepago)
	//var cortar_precio = template.FuncMap{"CutPrice": CutPrice}
	ctx.Render("DetalleDocumento.html", map[string]interface{}{
		"caja":         caja,
		"formasdepago": formasdepago,
		"html":         template.HTML(html),
		"calculadora":  template.HTML(calculadora),
		"total":        total,
		"idOperacion":  idOperacion,
	})
	//ctx.Render("CajaDetalle.html", nil)
}

//InsertaDocumento inserta un movimiento en pagos
func InsertaDocumento(ctx *iris.Context) {
	var SCaja CajaModel.SCaja
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	SCaja.SSesion.Name = NameUsrLoged
	SCaja.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	SCaja.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		SCaja.SEstado = false
		SCaja.SMsj = errSes.Error()
		ctx.Render("ZError.html", SCaja)
		return
	}

	//fmt.Println("hola")
	//######### LEE TU OBJETO DEL FORMULARIO #########
	var Caja CajaModel.CajaMgo
	//Dejamos de leer todo el formulario
	//ctx.ReadForm(&Caja)
	//Y leemos cada campo por separado.
	usuario := bson.ObjectIdHex("58e64924a137ae2a3461d995") //bson.NewObjectId() //ctx.FormValue("Usuario")
	cajax := bson.ObjectIdHex(ctx.FormValue("cajax"))       //ctx.FormValue("Caja")
	cargo, err := strconv.ParseFloat(ctx.FormValue("cargo"), 64)
	abono := 0.00
	saldo := cargo
	fmt.Println(err)
	fmt.Println(cargo)
	fmt.Println(usuario, cajax, cargo)
	if err == nil {
		fmt.Println("Información correcta.")
	} else {
		fmt.Println("Dato inválido.")
	}

	//abono := ctx.FormValue("Abono")
	//saldo := ctx.FormValue("Saldo")
	operacion := bson.ObjectIdHex(ctx.FormValue("operacion"))
	estatus := bson.NewObjectId()
	fechaHora := time.Now()

	//SCaja.Caja.ECargoCaja.Cargo = cargo
	//Lo asignamos a la estructura para guardado
	Caja.Usuario = usuario
	Caja.Caja = cajax
	Caja.Cargo = cargo
	Caja.Saldo = cargo
	Caja.Operacion = operacion
	Caja.Estatus = estatus
	Caja.FechaHora = fechaHora

	//fmt.Println(usuario, caja, cargo, operacion, estatus, fechaHora)

	//Caja = Caja{usuario, caja, cargo, cabono, saldo, fechaHora}
	//Lee como arreglo
	//data = ctx.FormValues()
	//fmt.Println(data["Cargo"])
	//######### VALIDA TU OBJETO #########
	EstatusPeticion := false //True indica que hay un error
	//##### TERMINA TU VALIDACION ########

	//########## Asigna vairables a la estructura que enviarás a la vista

	ID := bson.NewObjectId()
	Caja.ID = ID

	//######### ENVIA TUS RESULTADOS #########

	//Asigno el valor introducido para pintarlo en la vista
	SCaja.Caja.ID = ID
	SCaja.Caja.ECargoCaja.Cargo = cargo
	SCaja.Caja.EAbonoCaja.Abono = abono
	SCaja.Caja.ESaldoCaja.Saldo = saldo
	SCaja.Caja.EUsuarioCaja.Usuario = usuario
	SCaja.Caja.ECajaCaja.Caja = cajax
	SCaja.Caja.EOperacionCaja.Operacion = operacion
	SCaja.Caja.EFechaHoraCaja.FechaHora = fechaHora
	//Valido, si el campo CARGO viene vacio
	if err != nil {
		EstatusPeticion = true
		SCaja.Caja.ECargoCaja.IEstatus = true
		SCaja.Caja.ECargoCaja.IMsj = "Campo DOTACION no puede ser nulo"
	}
	//	SCaja.Caja = Caja //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.

	if EstatusPeticion {
		SCaja.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SCaja.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("CajaAlta.html", SCaja)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if Caja.InsertaMgo() {
			SCaja.SEstado = true
			SCaja.SMsj = "Se ha realizado una inserción exitosa"

			//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
			//ctx.Render("CajaDetalle.html", SCaja)
			x := (ID.Hex())
			Ticketpdf(x)
			ctx.Render("Comprobante.html", SCaja)
		} else {
			SCaja.SEstado = false
			SCaja.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("CajaAlta.html", SCaja)
		}

	}

}

//Ticketpdf genera PDF de operacion de caja.
func Ticketpdf(ID string) {
	pdf := gofpdf.New("P", "mm", "A6", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "COMPROBANTE DE PAGO")
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(40, 10, "Operacion: "+ID)
	pdf.Ln(5)
	pdf.Cell(40, 10, "Fecha: "+time.Now().String())
	err := pdf.OutputFileAndClose("PDF/" + ID + ".pdf")
	if err != nil {
		fmt.Println(err)
	}
}

//Emisor estructura provisional para generacion de XML
type Emisor struct {
	XMLName   xml.Name `xml:"cfdi:Receptor"`
	RFC       string   `xml:"RFC,attr"`
	FirstName string   `xml:"cfdi:DomicilioFiscal"`
	LastName  string   `xml:"cfdi:ExpedidoEn"`
	UserName  string   `xml:"cfdi:RegimenFiscal"`
}

//Receptor estructura provisional para generacion de XML
type Receptor struct {
	XMLName   xml.Name `xml:"cfdi:Receptor"`
	RFC       string   `xml:"RFC,attr"`
	FirstName string   `xml:"cfdi:Domicilio"`
}

//Item estructura provisional para generacion de XML
type Item struct {
	XMLName     xml.Name `xml:"cfdi:Conceptos"`
	Quantity    int      `xml:"cantidad"`
	Unit        string   `xml:"unidad"`
	Description string   `xml:"descripcion"`
	Value       float64  `xml:"valorUnitario"`
	Total       float64  `xml:"importe"`
}

//Cfdi estructura provisional para generacion de XML
type Cfdi struct {
	XMLName  xml.Name   `xml:"cfdi:Comprobante"`
	Xmlns    string     `xml:"xmlns:cfdi,attr"`
	Emisor   []Emisor   `xml:"emisor"`
	Receptor []Receptor `xml:"receptor"`
	Item     []Item     `xml:"conceptos"`
}

//CierraCaja cierra la caja en turno de usuario logueado
func CierraCaja(ctx *iris.Context) {
	cajaID := ctx.FormValue("caja")
	//Inicia XML
	fmt.Println("Genera XML")
	v := &Cfdi{Xmlns: "http://www.sat.gob.mx/cfd/3"}
	// add two staff details
	v.Emisor = append(v.Emisor, Emisor{RFC: "MERM830502IG8", FirstName: "Melchor", LastName: "Mendoza", UserName: "melchor"})
	v.Receptor = append(v.Receptor, Receptor{RFC: "MERM830502IG8", FirstName: "Las Palmas 405, San Sebastian Tutla, Oaxaca, Mx."})
	v.Item = append(v.Item, Item{Quantity: 2, Unit: "PZA", Description: "MOTO BOMBA HP", Value: 800.00, Total: 1600.00})
	filename := "XML/" + cajaID + ".xml"
	file, _ := os.Create(filename)
	xmlWriter := io.Writer(file)
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	Pdf(cajaID)
	//Actualizo el equipo caja en mongo
	//Voy por los datos de la caja
	eCaja := EquipoCajaModel.GetOne(bson.ObjectIdHex(cajaID))
	fmt.Println(eCaja)
	if !MoGeneral.EstaVacio(eCaja) {
		//Asigno valores, de manera fija mando el estatus CLOSE, revisar despues como resolver esto
		eCaja.Estatus = bson.ObjectIdHex("58fa87b4a137ae27b84ebc2f")
		//Ejecuto la acutalizacion
		up := eCaja.ActualizaMgo([]string{"Estatus"}, []interface{}{eCaja.Estatus})
		if up {
			fmt.Println("Actualizado: ", up)
		} else {
			fmt.Println("Actualizado: ", up)
		}
	}
}

//Pdf metodo imprime PDF
func Pdf(cajaID string) {
	//var Caja CajaModel.CajaMgo
	caja := CajaModel.GetAll()
	//fmt.Println(caja[0].Cargo)
	fmt.Println("Genera PDF")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "CORTE DE CAJA")
	pdf.Ln(10)
	pdf.Cell(40, 10, "CAJA: "+cajaID)
	t1 := time.Now()
	pdf.Ln(10)
	pdf.Cell(40, 10, "Fecha: "+t1.String())
	pdf.Ln(20)
	pdf.Cell(40, 10, "R E S U M E N")
	pdf.Ln(10)
	var lineas int
	var sumacargo float64
	var sumaabono float64
	for i, row := range caja {
		pdf.Ln(7)
		pdf.Cell(12, 6, strconv.Itoa(i))
		//pdf.CellFormat(12, 6, strconv.Itoa(i), "1", 0, "", false, 0, "")
		//pdf.Cell(45, 10, row.ID.Hex())
		cargo := strconv.FormatFloat(row.Cargo, 'f', 2, 64)
		pdf.Cell(30, 6, cargo)
		abono := strconv.FormatFloat(row.Abono, 'f', 2, 64)
		pdf.Cell(30, 6, abono)
		//pdf.Cell(80, 10, row.Operacion.Hex())
		pdf.Cell(60, 6, (row.FechaHora.String()))
		lineas = lineas + 1
		sumacargo = sumacargo + row.Cargo
		sumaabono = sumaabono + row.Abono
	}
	//
	pdf.Ln(10)
	pdf.Cell(40, 10, "T O T A L E S")
	pdf.Ln(15)
	pdf.Cell(12, 6, strconv.Itoa(lineas))
	sumacargox := strconv.FormatFloat(sumacargo, 'f', 2, 64)
	pdf.Cell(30, 6, sumacargox)
	sumaabonox := strconv.FormatFloat(sumaabono, 'f', 2, 64)
	pdf.Cell(30, 6, sumaabonox)
	pdf.Cell(60, 6, "")
	//
	pdf.Ln(10)
	pdf.Cell(60, 10, "SALDO EN CAJAS")
	saldo := strconv.FormatFloat(sumacargo-sumaabono, 'f', 2, 64)
	pdf.Cell(60, 10, saldo)
	//Pie
	//
	err := pdf.OutputFileAndClose("PDF/" + cajaID + ".pdf")
	if err != nil {
		fmt.Println(err)
	}
}

//GetSaldoCaja devuelve el saldo de la caja para mostrar en el cierre
func GetSaldoCaja(ctx *iris.Context) {
	templatePaginacion := ``
	//Trae documentos
	cajaID := ctx.FormValue("caja")
	fmt.Println(cajaID)
	caja := CajaModel.GetAll()
	fmt.Println("cajaSaldos:", caja)
	var lineas int
	var sumacargo float64
	var sumaabono float64
	for i, row := range caja {
		fmt.Println(i)
		lineas = lineas + 1
		sumacargo = sumacargo + row.Cargo
		sumaabono = sumaabono + row.Abono
	}
	//x := template.HTML(template)
	//ctx.Render("CajaSaldo.html", lineas)
	ctx.Render("CajaSaldo.html", map[string]interface{}{
		"caja":        caja,
		"PaginacionT": template.HTML(templatePaginacion),
	})
}

//CobrarDesdeVentaGet recibe un parámetro de referencia a la operación que se cobrará
func CobrarDesdeVentaGet(ctx *iris.Context) {
	var scaja CajaModel.SCaja
	var operaciones []PuntoVentaModel.PuntoVentaMgo
	var html string
	var productos []PuntoVentaModel.DatosVentaTemporal

	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	scaja.SSesion.Name = NameUsrLoged
	scaja.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	scaja.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		scaja.SEstado = false
		scaja.SMsj = errSes.Error()
		ctx.Render("ZError.html", scaja)
		return
	}

	usuarioID := bson.ObjectIdHex("58e64924a137ae2a3461d995")
	//Buscamos cajas para este usuario
	//caja := CajaModel.GetOne(ID)
	caja := CajaModel.GetCajaAbiertaByUsuario(usuarioID)
	vcaja := CajaModel.Caja{}
	vcaja.EUsuarioCaja.Usuario = caja.Usuario
	vcaja.ECajaCaja.Caja = caja.Caja
	vcaja.ESaldoCaja.Saldo = caja.Saldo
	scaja.Caja = vcaja

	ID := ctx.Param("ID") //"58e67cdda137ae2600adaea2"
	if bson.IsObjectIdHex(ID) {
		caja := CajaModel.GetOperacionesByID(bson.ObjectIdHex(ID))
		if !MoGeneral.EstaVacio(caja) {
			html, calculadora, total, idOperacion := PuntoVentaModel.ConsultaDatosOperacionByID(ID)
			formasdepago := MediosPagoModel.GetAll()

			cobro := fmt.Sprintf(`<div class="container">
	<div class="row">
    <div class="col col-lg-6">
		<table class="table table-hover table-sm table-bordered">
		<thead>
			<tr>
			<th>Código</th>
			<th>Descripción</th>
			<th>Cantidad</th>
			<th>Precio</th>
			<th>Unidad</th>
			<th>Importe</th>
			</tr>
		</thead>		
		<tbody>
			<input type="hidden" class="form-control" id="documentoIDx" value="%v">
			<input type="hidden" class="form-control" id="importeDocumento" value="%v">`, idOperacion, total)

			for _, v := range caja {
				cobro += fmt.Sprintf(`					
					<tr>					
						<th scope="row">%v<input type="text" class="form-control" id="documentoIDx" value="%v"></th>
						<td>%v</td>
						<td>%v</td>
						<td>%v<input type="text" class="form-control" id="importeDocumento" value="%v"></td>
					</tr>
					`, v.ID.Hex(), v.ID.Hex(), v.Tipo, v.Concepto, v.Monto, v.Monto)
			}

			cobro += fmt.Sprintf(`
			%v
		</tbody>
		<tfoot>
			<div class="pull-left">
			%v
			</div>
		</tfoot>
		</table>
		<div class="col col-lg-6">
		</div>
		<div class="col col-lg-6">
			<!--<p><div style="border-style: dotted; font-size:09pt; height:30vh" id="div1" ondrop="drop(event)" ondragover="allowDrop(event)"></div></p>-->
			<!--
			<ul style="font-size:09pt;" class="source connected">
				<li class="list-group-item active">+ Formas de pago</li>
				<li class="list-group-item"></li>			
			</ul>
			-->
			<div class="input-group">
			<input type="text" id="cambio" class="form-control input-lg" placeholder="Su cambio" readonly>
			<span class="input-group-btn">
				<!--<button class="btn btn-secondary" type="button">Go!</button>-->
				<button type="button" class="btn btn-primary btn-lg btn-block" onclick="insertaDocumento()">Aplicar</button>
			</span>
			</div>
		</div>
    </div>
	<div class="col col-lg-4">
		<!--
		<ul class="list-group">
			<li class="list-group-item active">Formas de pago</li>`, html, calculadora)

			for _, v := range formasdepago {
				cobro += fmt.Sprintf(`			
			<li style="cursor:move;" id="drag1" class="list-group-item" draggable="true" ondragstart="drag(event)"><p class="mb-1">%v</p></li>
			`, v.Nombre)
			}
			cobro += fmt.Sprintf(`
		</ul>
		-->
		<ul class="list-group target compra connected" style="font-size:09pt;">
			<li class="list-group-item active">Formas de pago</li>`)

			for _, v := range formasdepago {
				cobro += fmt.Sprintf(`	
				<li style="cursor:move" title="arrastre y suelte" class="list-group-item"></span> %v <div class="input-group"><span class="input-group-addon">Monto</span><input onkeyup="validaMonto('%v','%v')" id="%v" type="text" class="form-control fpagos" placeholder="0.00" aria-describedby="basic-addon1"></div>
				`, v.Descripcion, v.ID.Hex(), v.Cambio, v.ID.Hex())
				if v.Cambio {
					cobro += fmt.Sprintf(`<span class="label label-info">ACEPTA CAMBIO</span>`)
				} else {
					cobro += fmt.Sprintf(`<span class="label label-info">NO ACEPTA CAMBIO</span>`)
				}
				cobro += fmt.Sprintf(`</li>
					`)
			}
			cobro += fmt.Sprintf(`
					</ul>

				</div>
				</div>
			</div>`)

			ctx.Render("CajaOperaciones.html", map[string]interface{}{
				"scaja":  scaja,
				"EsPago": true,
				"Cobro":  template.HTML(cobro),
			})
			return
		}
		scaja.SEstado = false
		scaja.SMsj = "La Operación a la que hace referencia no se ha encontrado, favor de revisar su conexión e intentar de nuevo."

		operaciones = PuntoVentaModel.GetAll()
		//Voy a postgres por los movimientos, este metodo es del POS
		html, productos = PuntoVentaModel.ConsultaDatosOperacionAll()

		if !MoGeneral.EstaVacio(caja) {
			ctx.Render("CajaOperaciones.html", map[string]interface{}{
				"scaja":       scaja,
				"EsPago":      false,
				"operaciones": operaciones,
				"productos":   productos,
				"html":        template.HTML(html),
			})
		} else {
			scaja.SEstado = false
			scaja.SMsj = "Debes ABRIR tu caja."
			fmt.Println("NO existe caja abierta, mando vista para abrir caja.")
			ctx.Render("CajaAlta.html", scaja)
		}

	}
	scaja.SEstado = false
	scaja.SMsj = "Error al leer referencia, El parámetro no es adecuado."

	operaciones = PuntoVentaModel.GetAll()
	//Voy a postgres por los movimientos, este metodo es del POS
	html, productos = PuntoVentaModel.ConsultaDatosOperacionAll()
	if !MoGeneral.EstaVacio(caja) {
		ctx.Render("CajaOperaciones.html", map[string]interface{}{
			"scaja":       scaja,
			"EsPago":      false,
			"operaciones": operaciones,
			"productos":   productos,
			"html":        template.HTML(html),
		})

	} else {
		scaja.SEstado = false
		scaja.SMsj = "Debes ABRIR tu caja."
		fmt.Println("NO existe caja abierta, mando vista para abrir caja.")
		ctx.Render("CajaAlta.html", scaja)
	}

}

//BuscaPagina regresa la tabla de busqueda y su paginacion en el momento de especificar página
func BuscaPagina(ctx *iris.Context) {
	var Send CajaModel.SCaja

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

		Cabecera, Cuerpo := CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrToMongo))
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
	var Send CajaModel.SCaja
	var Cabecera, Cuerpo string

	grupo := ctx.FormValue("Grupox")
	if grupo != "" {
		gru, _ := strconv.Atoi(grupo)
		limitePorPagina = gru
	}

	cadenaBusqueda = ctx.FormValue("searchbox")
	//Send.Caja.ENombreCaja.Nombre = cadenaBusqueda

	if cadenaBusqueda != "" {

		docs := CajaModel.BuscarEnElastic(cadenaBusqueda)

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

			Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrToMongo))
			Send.SIndex.SCabecera = template.HTML(Cabecera)
			Send.SIndex.SBody = template.HTML(Cuerpo)
			MoConexion.FlushElastic()

			paginasTotales = MoGeneral.Totalpaginas(numeroRegistros, limitePorPagina)
			Paginacion := MoGeneral.ConstruirPaginacion(paginasTotales, 1)
			Send.SIndex.SPaginacion = template.HTML(Paginacion)

		} else {

			if numeroRegistros <= limitePorPagina {
				Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
			} else if numeroRegistros >= limitePorPagina {
				Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
			Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrIDMgo[0:numeroRegistros]))
		} else if numeroRegistros >= limitePorPagina {
			Cabecera, Cuerpo = CajaModel.GeneraTemplatesBusqueda(CajaModel.GetEspecifics(arrIDMgo[0:limitePorPagina]))
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
