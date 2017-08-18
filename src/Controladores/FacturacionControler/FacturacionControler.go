package FacturacionControler

import (
	"html/template"
	"strconv"

	"../../Modelos/FacturacionModel"
	"../../Modulos/Conexiones"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"

	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"

	_ "github.com/nakagami/firebirdsql"
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
var result []FacturacionModel.Facturacion
var resultPage []FacturacionModel.Facturacion
var templatePaginacion = ``

//wizard steps, valida completitud de proceso
var PasoWizard = 0

//0, sin iniciar
//1, id movimiento buscado
//2, Datos de usuario modifiicado
//3, Confirmacion
//4, Detalles Timbrado

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Facturacion
func IndexGet(ctx *iris.Context) {
	ctx.Render("FacturacionIndex.html", nil)
}

//IndexPost regresa la peticon post que se hizo desde el index de Facturacion
func IndexPost(ctx *iris.Context) {

	templatePaginacion = ``

	var resultados []FacturacionModel.FacturacionMgo
	var IDToObjID bson.ObjectId
	var arrObjIds []bson.ObjectId
	var arrToMongo []bson.ObjectId

	cadenaBusqueda = ctx.FormValue("searchbox")
	buscarEn = ctx.FormValue("buscaren")

	if cadenaBusqueda != "" {

		docs := FacturacionModel.BuscarEnElastic(cadenaBusqueda)

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

			resultados = FacturacionModel.GetEspecifics(arrToMongo)

			MoConexion.FlushElastic()

		}

	}

	templatePaginacion = ConstruirPaginacion()

	ctx.Render("FacturacionIndex.html", map[string]interface{}{
		"result":          resultados,
		"cadena_busqueda": cadenaBusqueda,
		"PaginacionT":     template.HTML(templatePaginacion),
	})

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Facturacion
func AltaGet(ctx *iris.Context) {
	ctx.Render("FacturacionAlta.html", nil)
}

//AltaPost regresa la petición post que se hizo desde el alta de Facturacion
func AltaPost(ctx *iris.Context) {
	//######### LEE TU OBJETO DEL FORMULARIO #########
	var Facturacion FacturacionModel.FacturacionMgo
	ctx.ReadForm(&Facturacion)

	//######### VALIDA TU OBJETO #########
	EstatusPeticion := true //True indica que hay un error
	//##### TERMINA TU VALIDACION ########

	//########## Asigna vairables a la estructura que enviarás a la vista
	Facturacion.ID = bson.NewObjectId()

	//######### ENVIA TUS RESULTADOS #########
	var SFacturacion FacturacionModel.SFacturacion

	//	SFacturacion.Facturacion = Facturacion //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.

	if EstatusPeticion {
		SFacturacion.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SFacturacion.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("FacturacionAlta.html", SFacturacion)
	} else {

		// //Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		// if Facturacion.InsertaMgo() {
		// 	SFacturacion.SEstado = true
		// 	SFacturacion.SMsj = "Se ha realizado una inserción exitosa"

		// 	//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
		// 	ctx.Render("FacturacionDetalle.html", SFacturacion)

		// } else {
		// 	SFacturacion.SEstado = false
		// 	SFacturacion.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
		// 	ctx.Render("FacturacionAlta.html", SFacturacion)
		// }

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Facturacion
func EditaGet(ctx *iris.Context) {
	ctx.Render("FacturacionEdita.html", nil)
}

//EditaPost regresa el resultado de la petición post generada desde la edición de Facturacion
func EditaPost(ctx *iris.Context) {
	ctx.Render("FacturacionEdita.html", nil)
}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	ctx.Render("FacturacionDetalle.html", nil)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	ctx.Render("FacturacionDetalle.html", nil)
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
				<a href="/Facturacions/1" aria-label="Primera">
				<span aria-hidden="true">&laquo;</span>
				</a>
			</li>`

	templateP += ``
	for i := 0; i <= paginasTotales; i++ {
		if i == 1 {

			templateP += `<li class="active"><a href="/Facturacions/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		} else if i > 1 && i < 11 {
			templateP += `<li><a href="/Facturacions/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`

		} else if i > 11 && i == paginasTotales {
			templateP += `<li><span aria-hidden="true">...</span></li><li><a href="/Facturacions/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		}
	}
	templateP += `<li><a href="/Facturacions/` + strconv.Itoa(paginasTotales) + `" aria-label="Ultima"><span aria-hidden="true">&raquo;</span></a></li></ul></nav>`
	return templateP
}

//####################< GenerarCFDI >##########################
//Obtencion de datos xml

//EliminarEspaciosInicioFinal Elimina los espacios en blanco Al inicio y final de una cadena:
//recibe cadena, regresa cadena limpia de espacios al inicio o final o "" si solo contiene espacios
func EliminarEspaciosInicioFinal(cadena string) string {
	var cadenalimpia string
	cadenalimpia = cadena
	re := regexp.MustCompile("(^\\s+|\\s+$)")
	cadenalimpia = re.ReplaceAllString(cadenalimpia, "")
	return cadenalimpia
}

//EliminarMultiplesEspaciosIntermedios Elimina los espacios en blanco de una cadena:
//recibe cadena, regresa cadena limpia  si solo contiene espacios
func EliminarMultiplesEspaciosIntermedios(cadena string) string {
	var cadenalimpia string
	cadenalimpia = cadena
	re := regexp.MustCompile("[\\s]+")
	cadenalimpia = re.ReplaceAllString(cadenalimpia, " ")
	return cadenalimpia
}

//LimpiarCadena Elimina los espacios en blanco de una cadena:
//recibe cadena, regresa cadena limpia o "" si solo contiene espacios
func LimpiarCadena(cadena string) string {
	var cadenalimpia string
	cadenalimpia = EliminarMultiplesEspaciosIntermedios(cadena)
	cadenalimpia = EliminarEspaciosInicioFinal(cadenalimpia)
	return cadenalimpia
}

// check Función usada para verificar error.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//GetSQLHeader obtener los datos principales delmovimiento a facturar
func GetSQLHeader(NumMov string, conn *sql.DB) FacturacionModel.CabeceraSQL {
	var cab FacturacionModel.CabeceraSQL
	consulta := fmt.Sprintf(`
	SELECT CAB.fcapmov as fecha,CAB.tipo_comp as tipodecomprobante,CFD.forpag2 AS formadepago
		,CAB.impmov + CAB.dctmov as subtotal,CAB.dctmov as descuento,CAB.nummon as moneda,CAB.impmov+CAB.ivatmov as total
		,CFD.forpag AS metododepago,CAB.pobmov AS lugarexpedicion,ATT.valattr AS numctapago
		,ctemov.numcte,CAB.numalm
	FROM maemovca02 CAB
		inner join movcte ctemov on (ctemov.nummov = cab.nummov)
		inner join maemovref REFCFD on (REFCFD.nummov =CAB.nummov)
		inner join maemovcfd CFD on (CFD.numcfd = REFCFD.numcfd)
		inner join maeattributes att on (CFD.numcfd = att.numcfd)
	WHERE CAB.nummov = %v
	`, NumMov)
	conn.QueryRow(consulta).Scan(&cab.Fecha, &cab.TipoComprobante, &cab.FormaPago, &cab.Subtotal,
		&cab.Descuento, &cab.Moneda, &cab.Total, &cab.MetodoPago, &cab.LugarExpedicion, &cab.NumCtaPago, &cab.NumCte, &cab.NumAlmacen)
	return cab
}

// SetCFDIHeader Funcion implementada para cargar los datos de la cabecera en la estructura de un CFD
func SetCFDIHeader(cab FacturacionModel.CabeceraSQL) *FacturacionModel.XMLCFDI {
	v := &FacturacionModel.XMLCFDI{}
	v.SchemaLocation = "http://www.sat.gob.mx/cfd/3 http://www.sat.gob.mx/sitio_internet/cfd/3/cfdv32.xsd http://www.sat.gob.mx/TimbreFiscalDigital http://www.sat.gob.mx/sitio_internet/TimbreFiscalDigital/TimbreFiscalDigital.xsd http://www.sat.gob.mx/detallista http://www.sat.gob.mx/sitio_internet/cfd/detallista/detallista.xsd http://www.sat.gob.mx/implocal http://www.sat.gob.mx/sitio_internet/cfd/implocal/implocal.xsd http://www.sat.gob.mx/iedu http://www.sat.gob.mx/sitio_internet/cfd/iedu/iedu.xsd http://www.sat.gob.mx/pfic http://www.sat.gob.mx/sitio_internet/cfd/pfic/pfic.xsd"
	v.SchemaLocationcfdi = "http://www.sat.gob.mx/cfd/3"
	v.SchemaLocationxsi = "http://www.w3.org/2001/XMLSchema-instance"
	v.SchemaLocationxs = "http://www.w3.org/2001/XMLSchema"
	v.SchemaLocationImplocal = "http://www.sat.gob.mx/implocal"
	v.Version = "3.2"
	strlayout := "2006-01-02T15:04:05"
	layoutF := cab.Fecha.Format(strlayout)
	v.FechaXML = layoutF
	v.TipoDeComprobante = LimpiarCadena(cab.TipoComprobante)
	v.FormaDePago = LimpiarCadena(cab.FormaPago)
	v.SubTotal = cab.Subtotal
	v.Descuento = cab.Descuento
	v.Moneda = LimpiarCadena(cab.Moneda)
	v.Total = cab.Total
	v.MetodoDePago = LimpiarCadena(cab.MetodoPago)
	v.LugarExpedicion = LimpiarCadena(cab.LugarExpedicion)
	v.Certificado = "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"
	v.NoCertificado = "00001000000403900497"
	v.Sello = "i9AZ4qeDXWStQROWdrAWsYNFaTm1GZD+/jUy6hnol3HUNY7a4pG1vNZK50x0A8oQsAnTiICw6ZkrFdRjTU76HfWVADIUQb+KPjC9SbrQ7dVS0hRjvcTMBPg7gnGBTH8oxMP0SKy9fuvStJaE7+/fflV05xUqvojaCy7bR7/bENKByPRN0EoGYOBdMnzbWxsTyGCkX0B5lAriVqKcM1TE1bT6SYYLMV2QoSNmxMZZPGkkPR/untUBaaGUcPKRIdFJrfcmF77vWm9gIrclbm/scS6nxk+oUWdEZO3j//RHj9KbFI1H8+Lfr3Otx9wv/sFVEtl8sUImgpMlfT9WYd/EJw=="
	return v
}

// GetConceptos Función realizada para cargar los datos de la lista de conceptos de un CFD
func GetConceptos(movimiento string, conn *sql.DB) []FacturacionModel.CFDIConcepto {
	consulta := fmt.Sprintf(`
	select 	DET.candtm as cantidad,DET.cveuni as unidad,DET.numart as noidentificacion,DET.nomdtm
			,DET.prudtm as valorunitario,DET.candtm * DET.prudtm as importe
	FROM maedtma02 DET
	WHERE DET.nummov=%v
	`, movimiento)

	rows, err := conn.Query(consulta)
	check(err)

	defer rows.Close()
	// var qty float64
	concepto := FacturacionModel.CFDIConcepto{}
	conceptoSql := FacturacionModel.ConceptosSQL{}
	var retorno []FacturacionModel.CFDIConcepto
	for rows.Next() {
		if err := rows.Scan(&conceptoSql.Cantidad, &conceptoSql.Unidad, &conceptoSql.NoIdentificacion, &conceptoSql.NomDTM,
			&conceptoSql.ValorUnitario, &conceptoSql.Importe); err != nil {
			panic(err)
		}
		i := fmt.Sprintf("%.4f", conceptoSql.ValorUnitario)
		conceptoSql.ValorUnitario, _ = strconv.ParseFloat(i, 4)
		i = fmt.Sprintf("%.4f", conceptoSql.Importe)
		conceptoSql.Importe, _ = strconv.ParseFloat(i, 4)
		i = fmt.Sprintf("%.4f", conceptoSql.Cantidad)
		conceptoSql.Cantidad, _ = strconv.ParseFloat(i, 4)
		concepto.Cantidad = conceptoSql.Cantidad
		concepto.Unidad = conceptoSql.Unidad
		concepto.NoIdentificacion = conceptoSql.NoIdentificacion
		concepto.Descripcion = conceptoSql.NomDTM
		concepto.ValorUnitario = conceptoSql.ValorUnitario
		concepto.Importe = conceptoSql.Importe
		retorno = append(retorno, concepto)
	}
	check(err)
	return retorno
}

// GetReceptor Funcion orientada a obtener los datos del receptor del CFD, regresa receptor:Receptor
func GetReceptor(NumCte string, conn *sql.DB) FacturacionModel.Receptor {
	var receptor FacturacionModel.Receptor
	consultaCte := fmt.Sprintf(`
	select 
		CTE.rfccte AS rfc,CTE.nomcte as nombre,CTE.dircte as calle,CTE.nlecte as noexterior, 
		COALESCE(CTE.nlicte,'') as nointerior
		,CTE.colcte as colonia
		,CTE.pobcte as localidad,CTE.muncte as municipio,CTE.estado,CTE.pais,CTE.cpcte as codigopostal
	FROM MAECTE CTE
    WHERE CTE.numcte=%v
	`, NumCte)
	err := conn.QueryRow(consultaCte).Scan(&receptor.RFC, &receptor.Nombre, &receptor.Domicilio.Calle, &receptor.Domicilio.NoExterior,
		&receptor.Domicilio.NoInterior, &receptor.Domicilio.Colonia, &receptor.Domicilio.Localidad,
		&receptor.Domicilio.Municipio, &receptor.Domicilio.Estado, &receptor.Domicilio.Pais,
		&receptor.Domicilio.CodigoPostal)
	check(err)
	return receptor
}

// GetEmisor Funcion orientada a obtener los datos del Emisor del CFD, regresa emisor:Emisor
func GetEmisor(NumAlmacen string, conn *sql.DB) FacturacionModel.Emisor {
	var emisor FacturacionModel.Emisor
	consultaEmpresa := fmt.Sprintf(`
    SELECT EMP.rfcempr AS RFC,EMP.nombempr AS NOMBRE,EMP.domempr AS CALLE,EMP.noextempr AS NOEXTERIOR,
        coalesce(EMP.nointempr,'') as NOINTERIOR,
        EMP.colempr AS COLONIA,EMP.pobempr AS LOCALIDAD,EMP.muniempr AS MUNICIPIO,EMP.edoempr AS ESTADO,
        EMP.paisempr AS PAIS ,EMP.cpempr AS CODIGOPOSTAL,(SELECT VARVAL
            FROM MAEVAR
            WHERE VARNAME = 'REGIMEN') AS REGIMENFISCAL
    FROM maeempr EMP
	`)
	err := conn.QueryRow(consultaEmpresa).Scan(&emisor.RFC, &emisor.Nombre, &emisor.DomicilioFiscal.Calle, &emisor.DomicilioFiscal.NoExterior,
		&emisor.DomicilioFiscal.NoInterior, &emisor.DomicilioFiscal.Colonia, &emisor.DomicilioFiscal.Localidad,
		&emisor.DomicilioFiscal.Municipio, &emisor.DomicilioFiscal.Estado, &emisor.DomicilioFiscal.Pais,
		&emisor.DomicilioFiscal.CodigoPostal, &emisor.CFDIRegimen.Regimen)
	check(err)
	consultaExpedidoEn := fmt.Sprintf(`
	SELECT ALM.diralm, ALM.noextalm,
		coalesce(ALM.nointalm ,'') as nointerior
		,ALM.colalm,ALM.pobalm,ALM.munialm,ALM.edoalm,ALM.paisalm,ALM.cpalm
	FROM
		almacen ALM
	WHERE ALM.numalm='%v'
	`, NumAlmacen)
	// fmt.Println("Consulta Cliente: ", consultaExpedidoEn)

	err = conn.QueryRow(consultaExpedidoEn).Scan(&emisor.ExpedidoEn.Calle, &emisor.ExpedidoEn.NoExterior,
		&emisor.ExpedidoEn.NoInterior, &emisor.ExpedidoEn.Colonia, &emisor.ExpedidoEn.Localidad,
		&emisor.ExpedidoEn.Municipio, &emisor.ExpedidoEn.Estado, &emisor.ExpedidoEn.Pais,
		&emisor.ExpedidoEn.CodigoPostal)
	check(err)

	return emisor
}

// GetTraslados Funcion que devuelve los impuestos trasladados de un comprobante
func GetTraslados(Movimiento string, conn *sql.DB) []FacturacionModel.CFDITraslado {
	consultaImpuestos := fmt.Sprintf(`
		SELECT 
			impu.numimpu as impuesto,impu.porimpu,imp.impimpu as importe
		FROM 
			dtmaimpu imp
		inner join 
			maeimpu impu 
		on 
			(impu.numimpu = imp.numimpu)
		WHERE imp.nummov =%v
		`, Movimiento)
	rows, err := conn.Query(consultaImpuestos)
	check(err)

	defer rows.Close()

	impuesto := FacturacionModel.CFDITraslado{}
	var traslados []FacturacionModel.CFDITraslado
	for rows.Next() {
		if err := rows.Scan(&impuesto.Impuesto, &impuesto.Tasa, &impuesto.Importe); err != nil {
			panic(err)
		}
		i := fmt.Sprintf("%.4f", impuesto.Importe)
		impuesto.Importe, _ = strconv.ParseFloat(i, 4)
		i = fmt.Sprintf("%.4f", impuesto.Tasa)
		impuesto.Tasa, _ = strconv.ParseFloat(i, 4)
		if impuesto.Impuesto == "IVA16" { //Si aparece iva, cambiar a expresion regular
			impuesto.Impuesto = "IVA"
		}
		traslados = append(traslados, impuesto)
	}
	return traslados
}

// GenerarArchivoXML Genera el XML producto del CFD, regresa erron en caso de fallo o nil en éxito
func GenerarArchivoXML(NumMov string, v FacturacionModel.XMLCFDI) error {
	filename := fmt.Sprintf("%v.xml", NumMov)
	f, err := os.Create(filename)
	check(err)
	xmlWriter := io.Writer(f)
	especsXML := `<?xml version="1.0" encoding="utf-8"?>
	`
	var bytesXMLEnc []byte
	bytesXMLEnc = []byte(especsXML)
	n, err := f.Write(bytesXMLEnc)
	fmt.Printf("Encabezado: %v bytes escritos.\n", n)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	err = enc.Encode(v)
	return err
}

// GenerarArchivoCadenaOriginal Funcion que genera la cadena original utilizando xsltproc del sistema unix/Linux
func GenerarArchivoCadenaOriginal(NumMov string) {
	fmt.Println("Ingreso a ")
	filename := fmt.Sprintf("Cadena-%v.txt", NumMov)
	f, err := os.Create(filename)
	check(err)
	// OriginalStringWriter := io.Writer(f)

	//Generacion de cadena original
	patron := "cadenaoriginal_3_2.xslt"
	archivoBase := fmt.Sprintf("%v.xml", NumMov)

	binary, lookErr := exec.LookPath("xsltproc")
	check(lookErr)
	fmt.Println(binary)
	args := []string{"xsltproc", patron, archivoBase}
	data, execErr := exec.Command(args[0], args[1], args[2]).Output()
	check(execErr)
	n, err := f.Write(data)
	check(err)
	fmt.Printf("Cadena Original: %v bytes escritos.\n", n)
}
