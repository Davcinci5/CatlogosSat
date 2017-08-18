package EmpresaControler

import (
	"fmt"
	"html/template"
	"regexp"
	"strconv"

	"../../Modelos/EmpresaModel"
	"../../Modulos/CargaCombos"
	"../../Modulos/Conexiones"

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
var result []EmpresaModel.Empresa
var resultPage []EmpresaModel.Empresa
var templatePaginacion = ``

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de Empresa
func IndexGet(ctx *iris.Context) {
	ctx.Render("EmpresaIndex.html", nil)
}

//IndexPost regresa la peticon post que se hizo desde el index de Empresa
func IndexPost(ctx *iris.Context) {

	templatePaginacion = ``

	var resultados []EmpresaModel.EmpresaMgo
	var IDToObjID bson.ObjectId
	var arrObjIds []bson.ObjectId
	var arrToMongo []bson.ObjectId

	cadenaBusqueda = ctx.FormValue("searchbox")
	buscarEn = ctx.FormValue("buscaren")

	if cadenaBusqueda != "" {

		docs := EmpresaModel.BuscarEnElastic(cadenaBusqueda)

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

			resultados = EmpresaModel.GetEspecifics(arrToMongo)

			MoConexion.FlushElastic()

		}

	}

	templatePaginacion = ConstruirPaginacion()

	ctx.Render("EmpresaIndex.html", map[string]interface{}{
		"result":          resultados,
		"cadena_busqueda": cadenaBusqueda,
		"PaginacionT":     template.HTML(templatePaginacion),
	})

}

//###########################< ALTA >################################

//AltaGet renderea al alta de Empresa
func AltaGet(ctx *iris.Context) {
	ctx.Render("EmpresaAlta.html", nil)
}

//AltaPost regresa la petición post que se hizo desde el alta de Empresa
func AltaPost(ctx *iris.Context) {

	//######### LEE TU OBJETO DEL FORMULARIO #########
	var Empresa EmpresaModel.EmpresaMgo
	ctx.ReadForm(&Empresa)

	//######### VALIDA TU OBJETO #########
	EstatusPeticion := true //True indica que hay un error
	//##### TERMINA TU VALIDACION ########

	//########## Asigna vairables a la estructura que enviarás a la vista
	Empresa.ID = bson.NewObjectId()

	//######### ENVIA TUS RESULTADOS #########
	var SEmpresa EmpresaModel.SEmpresa

	//	SEmpresa.Empresa = Empresa //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.

	if EstatusPeticion {
		SEmpresa.SEstado = false                                                           //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SEmpresa.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("EmpresaAlta.html", SEmpresa)
	} else {

		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if Empresa.InsertaMgo() {
			SEmpresa.SEstado = true
			SEmpresa.SMsj = "Se ha realizado una inserción exitosa"

			//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
			ctx.Render("EmpresaDetalle.html", SEmpresa)

		} else {
			SEmpresa.SEstado = false
			SEmpresa.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde"
			ctx.Render("EmpresaAlta.html", SEmpresa)
		}

	}

}

//###########################< EDICION >###############################

//EditaGet renderea a la edición de Empresa
func EditaGet(ctx *iris.Context) {
	Sempresa := EmpresaModel.SEmpresa{}
	Sempresas := EmpresaModel.GetAll()
	id := ""
	if len(Sempresas) == 0 {
		fmt.Println("No existe una empresa")
		tmpl := CargaCombos.CargaComboCatalogo(99, "")
		tipoCon := CargaCombos.CargaComboCatalogo(143, "")
		MetodoC := CargaCombos.CargaComboCatalogo(144, "")
		Sempresa.Empresa.EDatosComercialesEmpresa.DatosComerciales.Domicilio.EEstadoDomicilio.Ihtml = template.HTML(tmpl)
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EEstadoDomicilio.Ihtml = template.HTML(tmpl)
		Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.ETipoCorreoYEnvio.Ihtml = template.HTML(tipoCon)
		Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.ECifradoCorreoYEnvio.Ihtml = template.HTML(MetodoC)
		Sempresa.SEstado = false
		Sempresa.SMsj = "No existe empresa registrada."
		id = ""
	} else if len(Sempresas) == 1 {
		fmt.Println("Existe una empresa: ", Sempresas[0].ID)
		Sempresa.Empresa.EDatosComercialesEmpresa.DatosComerciales.ENombreDatosCommerciales.Nombre = Sempresas[0].DatosComerciales.Nombre
		Sempresa.Empresa.EDatosComercialesEmpresa.DatosComerciales.Domicilio.ECalleDomicilio.Calle = Sempresas[0].DatosComerciales.Domicilio.Calle
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ENumExteriorDomicilio.NumExterior = Sempresas[0].DatosComerciales.Domicilio.NumExterior
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ENumInteriorDomicilio.NumInterior = Sempresas[0].DatosComerciales.Domicilio.NumInterior
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EColoniaDomicilio.Colonia = Sempresas[0].DatosComerciales.Domicilio.Colonia
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EMunicipioDomicilio.Municipio = Sempresas[0].DatosComerciales.Domicilio.Municipio
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EEstadoDomicilio.Estado = Sempresas[0].DatosComerciales.Domicilio.Estado
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EEstadoDomicilio.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(99, Sempresas[0].DatosComerciales.Domicilio.Estado))
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EPaisDomicilio.Pais = Sempresas[0].DatosComerciales.Domicilio.Pais
		Sempresa.Empresa.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ECPDomicilio.CP = Sempresas[0].DatosComerciales.Domicilio.CP

		Sempresa.Empresa.DatosComerciales.EContactosDatosCommerciales.Contactos.EAliasContactos.Alias = Sempresas[0].DatosComerciales.Contactos.Alias
		Sempresa.Empresa.DatosComerciales.EContactosDatosCommerciales.Contactos.EEmailContactos.Email = Sempresas[0].DatosComerciales.Contactos.Email
		Sempresa.Empresa.DatosComerciales.EContactosDatosCommerciales.Contactos.ETelefonoContactos.Telefono = Sempresas[0].DatosComerciales.Contactos.Telefono
		Sempresa.Empresa.DatosComerciales.EContactosDatosCommerciales.Contactos.EMovilContactos.Movil = Sempresas[0].DatosComerciales.Contactos.Movil

		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.ERazonSocialDatosFiscales.RazonSocial = Sempresas[0].DatosFiscales.RazonSocial
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.ERFCDatosFiscales.RFC = Sempresas[0].DatosFiscales.RFC
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECalleDomicilio.Calle = Sempresas[0].DatosFiscales.Domicilio.Calle
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ENumExteriorDomicilio.NumExterior = Sempresas[0].DatosFiscales.Domicilio.NumExterior
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ENumInteriorDomicilio.NumInterior = Sempresas[0].DatosFiscales.Domicilio.NumInterior
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EColoniaDomicilio.Colonia = Sempresas[0].DatosFiscales.Domicilio.Colonia
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EMunicipioDomicilio.Municipio = Sempresas[0].DatosFiscales.Domicilio.Municipio
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EEstadoDomicilio.Estado = Sempresas[0].DatosFiscales.Domicilio.Estado
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EEstadoDomicilio.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(99, Sempresas[0].DatosFiscales.Domicilio.Estado))
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EPaisDomicilio.Pais = Sempresas[0].DatosFiscales.Domicilio.Pais
		Sempresa.Empresa.EDatosFiscalesEmpresa.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECPDomicilio.CP = Sempresas[0].DatosFiscales.Domicilio.CP

		Sempresa.Empresa.EDatosFacturaEmpresa.DatosFactura.EKeyDatosFactura.Key = Sempresas[0].DatosFactura.Key
		Sempresa.Empresa.EDatosFacturaEmpresa.DatosFactura.ECerDatosFactura.Cer = Sempresas[0].DatosFactura.Cer
		Sempresa.Empresa.EDatosFacturaEmpresa.DatosFactura.EPemDatosFactura.Pem = Sempresas[0].DatosFactura.Pem
		tipoCon := CargaCombos.CargaComboCatalogo(143, Sempresas[0].CorreoYEnvio.Tipo)
		MetodoC := CargaCombos.CargaComboCatalogo(144, Sempresas[0].CorreoYEnvio.Cifrado)
		Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.ECorreoCorreoYEnvio.Correo = Sempresas[0].CorreoYEnvio.Correo
		Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.EPassCorreoYEnvio.Pass = Sempresas[0].CorreoYEnvio.Pass
		Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.ETipoCorreoYEnvio.Ihtml = template.HTML(tipoCon)
		Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.EPuertoCorreoYEnvio.Puerto = Sempresas[0].CorreoYEnvio.Puerto
		Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.ECifradoCorreoYEnvio.Ihtml = template.HTML(MetodoC)
		Sempresa.SEstado = false
		Sempresa.Empresa.ID = Sempresas[0].ID
		id = Sempresa.Empresa.ID.Hex()
		Sempresa.SMsj = "Una empresa encontrada."
	} else {
		fmt.Println("Los dato son Inconsistentes")
		Sempresa.SEstado = true
		Sempresa.SMsj = "Datos inconsistentes"
		// Sempresa.Empresa.ID = ""
		id = ""
	}
	fmt.Println(id)
	ctx.Render("EmpresaEdita.html", Sempresa)
}

//EditaPost regresa el resultado de la petición post generada desde la edición de Empresa
func EditaPost(ctx *iris.Context) {
	Sempresa := EmpresaModel.SEmpresa{}
	//var EmpresaDt EmpresaModel.EmpresaMgo
	Sempresa.SEstado = false
	result := EmpresaModel.EmpresaMgo{}

	MyDatosComerciales := EmpresaModel.EDatosComercialesEmpresa{}

	Nombre := ctx.FormValue("Nombre")
	Nombre = LimpiarCadena(Nombre)
	// EmpresaDt.DatosComerciales.Nombre = Nombre
	if CadenaVacia(Nombre) {
		MyDatosComerciales.DatosComerciales.ENombreDatosCommerciales.IEstatus = true
		MyDatosComerciales.DatosComerciales.ENombreDatosCommerciales.IMsj = "El campo Nombre de la seccion Comercial se encuentra vacío."
		Sempresa.SMsj = "El campo Nombre de la seccion Comercial se encuentra vacío."
		Sempresa.SEstado = true
	}
	CalleComercial := ctx.FormValue("CalleComercial")
	CalleComercial = LimpiarCadena(CalleComercial)
	// EmpresaDt.DatosComerciales.Domicilio.Calle = CalleComercial
	if CadenaVacia(CalleComercial) {
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ECalleDomicilio.IEstatus = true
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ECalleDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "El campo Calle de la seccion Comercial se encuentra vacío."
		Sempresa.SEstado = true
	}

	NumExtComercial := ctx.FormValue("NumExtComercial")
	NumExtComercial = LimpiarCadena(NumExtComercial)
	// EmpresaDt.DatosComerciales.Domicilio.NumExterior = NumExtComercial
	if CadenaVacia(NumExtComercial) {
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ENumExteriorDomicilio.IEstatus = true
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ENumExteriorDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "El campo Num. Exterior de la seccion Comercial se encuentra vacío."
		Sempresa.SEstado = true
	}

	NumIntComercial := ctx.FormValue("NumIntComercial")
	NumIntComercial = LimpiarCadena(NumIntComercial)
	// EmpresaDt.DatosComerciales.Domicilio.NumInterior = NumIntComercial

	ColoniaLocalidadComercial := ctx.FormValue("ColoniaLocalidadComercial")
	ColoniaLocalidadComercial = LimpiarCadena(ColoniaLocalidadComercial)
	// EmpresaDt.DatosComerciales.Domicilio.Colonia = ColoniaLocalidadComercial
	if CadenaVacia(ColoniaLocalidadComercial) {
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EColoniaDomicilio.IEstatus = true
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EColoniaDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "El campo Colonia de la seccion Comercial se encuentra vacío."
		Sempresa.SEstado = true
	}

	MunicipioComercial := ctx.FormValue("MunicipioComercial")
	MunicipioComercial = LimpiarCadena(MunicipioComercial)
	// EmpresaDt.DatosComerciales.Domicilio.Municipio = MunicipioComercial
	if CadenaVacia(MunicipioComercial) {
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EMunicipioDomicilio.IEstatus = true
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EMunicipioDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "El campo Municipio de la seccion Comercial se encuentra vacío."
		Sempresa.SEstado = true
	}

	EstadoComercial := ctx.FormValue("EstadoComercial")
	EstadoComercial = LimpiarCadena(EstadoComercial)
	// EmpresaDt.DatosComerciales.Domicilio.Estado = EstadoComercial

	PaisComercial := ctx.FormValue("PaisComercial")
	PaisComercial = LimpiarCadena(PaisComercial)
	// EmpresaDt.DatosComerciales.Domicilio.Pais = PaisComercial
	if CadenaVacia(PaisComercial) {
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EPaisDomicilio.IEstatus = true
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EPaisDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "El campo País de la seccion Comercial se encuentra vacío."
		Sempresa.SEstado = true
	}
	CPComercial := ctx.FormValue("CPComercial")
	CPComercial = LimpiarCadena(CPComercial)
	// EmpresaDt.DatosComerciales.Domicilio.CP = CPComercial
	if CadenaVacia(CPComercial) {
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EPaisDomicilio.IEstatus = true
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EPaisDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "El campo Código Postal de la seccion Comercial se encuentra vacío."
	} else if !CPValido(CPComercial) {
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ECPDomicilio.IEstatus = true
		MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ECPDomicilio.IMsj = "No corresponde a un C.P. valido (cinco digitos)."
		Sempresa.SMsj = "El campo Código Postal de la seccion Comercial no es válido."
		Sempresa.SEstado = true
	}

	MyDatosComerciales.DatosComerciales.ENombreDatosCommerciales.Nombre = Nombre
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ECalleDomicilio.Calle = CalleComercial
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ENumExteriorDomicilio.NumExterior = NumExtComercial
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ENumInteriorDomicilio.NumInterior = NumIntComercial
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EColoniaDomicilio.Colonia = ColoniaLocalidadComercial
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EMunicipioDomicilio.Municipio = MunicipioComercial
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EEstadoDomicilio.Estado = EstadoComercial
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EEstadoDomicilio.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(99, EstadoComercial))
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.EPaisDomicilio.Pais = PaisComercial
	MyDatosComerciales.DatosComerciales.EDomicilioDatosCommerciales.Domicilio.ECPDomicilio.CP = CPComercial

	// Indice:  Alias
	// Indice:  Email
	// Indice:  Movil
	// Indice:  Telefono
	contacto := EmpresaModel.ContactoMgo{}
	contactoVista := EmpresaModel.Contacto{}
	Alias := ctx.FormValue("Alias")
	Alias = LimpiarCadena(Alias)
	Email := ctx.FormValue("Email")
	Email = LimpiarCadena(Email)
	Movil := ctx.FormValue("Movil")
	Movil = LimpiarCadena(Movil)
	Telefono := ctx.FormValue("Telefono")
	Telefono = LimpiarCadena(Telefono)

	// EmpresaDt.DatosComerciales.Contactos
	if CadenaVacia(Alias) {
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EAliasContactos.IEstatus = true
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EAliasContactos.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "Alias de la sección Contacto Comercial vacío"
		Sempresa.SEstado = true
	}

	if CadenaVacia(Email) {
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EEmailContactos.IEstatus = true
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EEmailContactos.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "Email de la sección Contacto Comercial vacío"
		Sempresa.SEstado = true
	}
	if CadenaVacia(Movil) {
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EMovilContactos.IEstatus = true
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EMovilContactos.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "Movil de la sección Contacto Comercial vacío"
		Sempresa.SEstado = true
	}

	if CadenaVacia(Telefono) {
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.ETelefonoContactos.IEstatus = true
		MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.ETelefonoContactos.IMsj = "El campo se encuentra vacío."
		Sempresa.SMsj = "Telefono de la sección Contacto Comercial vacío"
		Sempresa.SEstado = true
	}
	contacto.Alias = Alias
	contacto.Email = Email
	contacto.Telefono = Telefono
	contacto.Movil = Movil
	contactoVista.EAliasContactos.Alias = Alias
	contactoVista.EEmailContactos.Email = Email
	contactoVista.ETelefonoContactos.Telefono = Telefono
	contactoVista.EMovilContactos.Movil = Movil
	// EmpresaDt.DatosComerciales.DatosComerciales = append(EmpresaDt.DatosComerciales.DatosComerciales, contacto)
	MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos = contactoVista
	// EAliasContactos.Alias = Alias
	// MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EEmailContactos.Email = Email
	// MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.ETelefonoContactos.Telefono = Telefono
	// MyDatosComerciales.DatosComerciales.EContactosDatosCommerciales.Contactos.EMovilContactos.Movil = Movil

	// Indice:  RFC
	// Indice:  RazonSocial
	// Indice:  CalleFiscal
	// Indice:  NumExtFiscal
	// Indice:  NumIntFiscal
	// Indice:  ColoniaLocalidadFiscal
	// Indice:  MunicipioFiscal
	// Indice:  EstadoFiscal
	// Indice:  PaisFiscal
	// Indice:  CPFiscal
	MyDatosFiscales := EmpresaModel.EDatosFiscalesEmpresa{}

	RFC := ctx.FormValue("RFC")
	RFC = LimpiarCadena(RFC)
	if CadenaVacia(RFC) {
		MyDatosFiscales.DatosFiscales.ERFCDatosFiscales.IEstatus = true
		MyDatosFiscales.DatosFiscales.ERFCDatosFiscales.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
		Sempresa.SMsj = "RFC de la sección Dato Fiscales vacío"
	} else if !RFCValido(RFC) {
		MyDatosFiscales.DatosFiscales.ERazonSocialDatosFiscales.IEstatus = true
		MyDatosFiscales.DatosFiscales.ERazonSocialDatosFiscales.IMsj = "El Valor del campo no es un RFC."
		Sempresa.SEstado = true
	}
	RazonSocial := ctx.FormValue("RazonSocial")
	RazonSocial = LimpiarCadena(RazonSocial)
	if CadenaVacia(RazonSocial) {
		MyDatosFiscales.DatosFiscales.ERFCDatosFiscales.IEstatus = true
		MyDatosFiscales.DatosFiscales.ERFCDatosFiscales.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
	}

	CalleFiscal := ctx.FormValue("CalleFiscal")
	CalleFiscal = LimpiarCadena(CalleFiscal)
	if CadenaVacia(CalleFiscal) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECalleDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECalleDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true

	}

	NumExtFiscal := ctx.FormValue("NumExtFiscal")
	NumExtFiscal = LimpiarCadena(NumExtFiscal)
	if CadenaVacia(NumExtFiscal) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ENumExteriorDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ENumExteriorDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true

	}

	NumIntFiscal := ctx.FormValue("NumIntFiscal")
	NumIntFiscal = LimpiarCadena(NumIntFiscal)

	ColoniaLocalidadFiscal := ctx.FormValue("ColoniaLocalidadFiscal")
	ColoniaLocalidadFiscal = LimpiarCadena(ColoniaLocalidadFiscal)
	if CadenaVacia(ColoniaLocalidadFiscal) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EColoniaDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EColoniaDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
	}
	MunicipioFiscal := ctx.FormValue("MunicipioFiscal")
	MunicipioFiscal = LimpiarCadena(MunicipioFiscal)
	if CadenaVacia(MunicipioFiscal) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EMunicipioDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EMunicipioDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
	}

	EstadoFiscal := ctx.FormValue("EstadoFiscal")
	EstadoFiscal = LimpiarCadena(EstadoFiscal)

	PaisFiscal := ctx.FormValue("PaisFiscal")
	PaisFiscal = LimpiarCadena(PaisFiscal)
	if CadenaVacia(PaisFiscal) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EPaisDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EPaisDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
	}

	CPFiscal := ctx.FormValue("CPFiscal")
	CPFiscal = LimpiarCadena(CPFiscal)
	if CadenaVacia(CPFiscal) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EPaisDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EPaisDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
	} else if !CPValido(CPFiscal) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECPDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECPDomicilio.IMsj = "No corresponde a un C.P. valido (cinco digitos)."
		Sempresa.SEstado = true
	}
	MyDatosFiscales.DatosFiscales.ERFCDatosFiscales.RFC = RFC
	MyDatosFiscales.DatosFiscales.ERazonSocialDatosFiscales.RazonSocial = RazonSocial
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECalleDomicilio.Calle = CalleFiscal
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ENumExteriorDomicilio.NumExterior = NumExtFiscal
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ENumInteriorDomicilio.NumInterior = NumIntFiscal
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EColoniaDomicilio.Colonia = ColoniaLocalidadFiscal
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EMunicipioDomicilio.Municipio = MunicipioFiscal
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EEstadoDomicilio.Estado = EstadoFiscal
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EEstadoDomicilio.Ihtml = template.HTML(CargaCombos.CargaComboCatalogo(99, EstadoFiscal))

	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EPaisDomicilio.Pais = PaisFiscal
	MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.ECPDomicilio.CP = CPFiscal

	// Indice:  Key
	// Indice:  Cer
	// Indice:  Pem
	MyDatosFactura := EmpresaModel.EDatosFacturaEmpresa{}
	Key := ctx.FormValue("Key")
	Cer := ctx.FormValue("Cer")
	Pem := ctx.FormValue("Pem")
	MyDatosFactura.DatosFactura.EKeyDatosFactura.Key = Key
	MyDatosFactura.DatosFactura.ECerDatosFactura.Cer = Cer
	MyDatosFactura.DatosFactura.EPemDatosFactura.Pem = Pem

	MyDatosCorreoEnvio := EmpresaModel.ECorreoYEnvioEmpresa{}
	Correo := ctx.FormValue("Correo")
	Pass := ctx.FormValue("Pass")
	Cifrado := ctx.FormValue("Cifrado")
	Puerto := ctx.FormValue("Puerto")
	Tipo := ctx.FormValue("Tipo")
	MyDatosCorreoEnvio.CorreoYEnvio.ECorreoCorreoYEnvio.Correo = Correo
	MyDatosCorreoEnvio.CorreoYEnvio.EPassCorreoYEnvio.Pass = Pass
	MyDatosCorreoEnvio.CorreoYEnvio.ECifradoCorreoYEnvio.Cifrado = Cifrado
	MyDatosCorreoEnvio.CorreoYEnvio.EPuertoCorreoYEnvio.Puerto = Puerto
	MyDatosCorreoEnvio.CorreoYEnvio.ETipoCorreoYEnvio.Tipo = Tipo

	Tipo = LimpiarCadena(Tipo)
	if CadenaVacia(Tipo) {
		MyDatosCorreoEnvio.CorreoYEnvio.ETipoCorreoYEnvio.IEstatus = true
		MyDatosCorreoEnvio.CorreoYEnvio.ETipoCorreoYEnvio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
	}
	Cifrado = LimpiarCadena(Cifrado)
	if CadenaVacia(Cifrado) {
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EColoniaDomicilio.IEstatus = true
		MyDatosFiscales.DatosFiscales.EDomicilioDatosFiscales.Domicilio.EColoniaDomicilio.IMsj = "El campo se encuentra vacío."
		Sempresa.SEstado = true
	}

	Sempresa.Empresa.EDatosComercialesEmpresa = MyDatosComerciales
	Sempresa.Empresa.EDatosFiscalesEmpresa = MyDatosFiscales
	Sempresa.Empresa.EDatosFacturaEmpresa = MyDatosFactura
	Sempresa.Empresa.ECorreoYEnvioEmpresa = MyDatosCorreoEnvio

	result.DatosComerciales.Nombre = Nombre
	result.DatosComerciales.Domicilio.Calle = CalleComercial
	result.DatosComerciales.Domicilio.NumExterior = NumExtComercial
	result.DatosComerciales.Domicilio.NumInterior = NumIntComercial
	result.DatosComerciales.Domicilio.Colonia = ColoniaLocalidadComercial
	result.DatosComerciales.Domicilio.Municipio = MunicipioComercial
	result.DatosComerciales.Domicilio.Estado = EstadoComercial
	result.DatosComerciales.Domicilio.Pais = PaisComercial
	result.DatosComerciales.Domicilio.CP = CPComercial
	result.DatosComerciales.Contactos = contacto
	// result.DatosComerciales.Contactos = append(result.DatosComerciales.Contactos, ContactoMgo{Alias: Alias, Email: Email, Telefono: Telefono, Movil: Movil})

	result.DatosFiscales.RFC = RFC
	result.DatosFiscales.RazonSocial = RazonSocial
	result.DatosFiscales.Domicilio.Calle = CalleFiscal
	result.DatosFiscales.Domicilio.NumExterior = NumExtFiscal
	result.DatosFiscales.Domicilio.NumInterior = NumIntFiscal
	result.DatosFiscales.Domicilio.Colonia = ColoniaLocalidadFiscal
	result.DatosFiscales.Domicilio.Municipio = MunicipioFiscal
	result.DatosFiscales.Domicilio.Estado = EstadoFiscal
	result.DatosFiscales.Domicilio.Pais = PaisFiscal
	result.DatosFiscales.Domicilio.CP = CPFiscal

	result.DatosFactura.Cer = Key
	result.DatosFactura.Key = Cer
	result.DatosFactura.Pem = Pem

	result.CorreoYEnvio.Correo = Correo
	result.CorreoYEnvio.Pass = Pass
	result.CorreoYEnvio.Puerto = Puerto
	result.CorreoYEnvio.Tipo = Tipo
	result.CorreoYEnvio.Cifrado = Cifrado

	tipoCon := CargaCombos.CargaComboCatalogo(143, Tipo)
	MetodoC := CargaCombos.CargaComboCatalogo(144, Cifrado)

	Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.ETipoCorreoYEnvio.Ihtml = template.HTML(tipoCon)
	Sempresa.Empresa.ECorreoYEnvioEmpresa.CorreoYEnvio.ECifradoCorreoYEnvio.Ihtml = template.HTML(MetodoC)

	fmt.Println(result)
	ID := ctx.FormValue("idEmpresa")
	if !Sempresa.SEstado {
		if ID != "" {
			objeto := EmpresaModel.GetOne(bson.ObjectIdHex(ID))
			if objeto.ID.Hex() == ID {
				objeto.DatosComerciales = result.DatosComerciales
				objeto.DatosFiscales = result.DatosFiscales
				objeto.DatosFactura = result.DatosFactura
				objeto.CorreoYEnvio = result.CorreoYEnvio
				objeto.ReemplazaMgo()
				fmt.Println("Intenta Actualizar")
				Sempresa.Empresa.ID = objeto.ID
				Sempresa.SEstado = false
				Sempresa.SMsj = "Se ha realizado la actualización."
			} else {
				fmt.Println("Se recibio un id de empresa diferente")
				Sempresa.SEstado = true
				Sempresa.SMsj = "No coincide o no encuentra el OID al actualizar"
			}
		} else {
			NewID := bson.NewObjectId()
			result.ID = NewID
			result.InsertaMgo()
			Sempresa.Empresa.ID = result.ID
			fmt.Println("Intenta insertar")
			Sempresa.SEstado = false
			Sempresa.SMsj = "Se ha realizado la inserción."
		}
	}
	ctx.Render("EmpresaEdita.html", Sempresa)

}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
	ctx.Render("EmpresaDetalle.html", nil)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
	ctx.Render("EmpresaDetalle.html", nil)
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
				<a href="/Empresas/1" aria-label="Primera">
				<span aria-hidden="true">&laquo;</span>
				</a>
			</li>`

	templateP += ``
	for i := 0; i <= paginasTotales; i++ {
		if i == 1 {

			templateP += `<li class="active"><a href="/Empresas/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		} else if i > 1 && i < 11 {
			templateP += `<li><a href="/Empresas/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`

		} else if i > 11 && i == paginasTotales {
			templateP += `<li><span aria-hidden="true">...</span></li><li><a href="/Empresas/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		}
	}
	templateP += `<li><a href="/Empresas/` + strconv.Itoa(paginasTotales) + `" aria-label="Ultima"><span aria-hidden="true">&raquo;</span></a></li></ul></nav>`
	return templateP
}

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

//RFCValido Sirve para reconocer si un RFC es Válido:
//recibe cadena, regresa true si es un RFC válido, false en otro caso
func RFCValido(rfc string) bool {
	re := regexp.MustCompile("^([a-zA-Z]{3}|[a-zA-Z]{4})\\d{6}[a-zA-Z0-9]{3}$")
	return re.MatchString(rfc)
}

//CPValido Sirve para reconocer si un CP es Válido:
//recibe cadena, regresa true si es un CP válido, false en otro caso
func CPValido(CP string) bool {
	re := regexp.MustCompile("^[0-9]{5}$")
	return re.MatchString(CP)
}

//TelOCelValido Sirve para reconocer si un Telefono o Celular es Válido:
//recibe cadena, regresa true si es un el valor recibido es un conjunto de 10 digitos, false en otro caso
func TelOCelValido(CP string) bool {
	re := regexp.MustCompile("^([0-9]{10}$")
	return re.MatchString(CP)
}

//CadenaVacia Sirve para reconocer si un RFC es Válido:
//recibe cadena, regresa true si es un RFC válido, false en otro caso
func CadenaVacia(cadena string) bool {
	if cadena == "" {
		return true
	}
	return false
}
