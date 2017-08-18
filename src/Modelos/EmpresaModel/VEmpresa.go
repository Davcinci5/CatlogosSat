package EmpresaModel

import (
	"html/template"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EDatosComercialesEmpresa Estructura de campo de Empresa
type EDatosComercialesEmpresa struct {
	DatosComerciales Comercial
	IEstatus         bool
	IMsj             string
	Ihtml            template.HTML
}

//EDatosFiscalesEmpresa Estructura de campo de Empresa
type EDatosFiscalesEmpresa struct {
	DatosFiscales Fiscal
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EDatosFacturaEmpresa Estructura de campo de Empresa
type EDatosFacturaEmpresa struct {
	DatosFactura Factura
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//ECorreoYEnvioEmpresa Estructura de campo de Empresa
type ECorreoYEnvioEmpresa struct {
	CorreoYEnvio ConfiguracionCorreo
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//Empresa estructura de Empresas mongo
type Empresa struct {
	ID bson.ObjectId
	EDatosComercialesEmpresa
	EDatosFiscalesEmpresa
	EDatosFacturaEmpresa
	ECorreoYEnvioEmpresa
}

//SEmpresa estructura de Empresas para la vista
type SEmpresa struct {
	SEstado bool
	SMsj    string
	Empresa
}

//ENombreDatosCommerciales Estructura de campo de DatosCommerciales
type ENombreDatosCommerciales struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDomicilioDatosCommerciales Estructura de campo de DatosCommerciales
type EDomicilioDatosCommerciales struct {
	Domicilio Direccion
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EContactosDatosCommerciales Estructura de campo de DatosCommerciales
type EContactosDatosCommerciales struct {
	Contactos Contacto
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Comercial subestructura de Empresa
type Comercial struct {
	ENombreDatosCommerciales
	EDomicilioDatosCommerciales
	EContactosDatosCommerciales
}

//ERazonSocialDatosFiscales Estructura de campo de DatosFiscales
type ERazonSocialDatosFiscales struct {
	RazonSocial string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//ERFCDatosFiscales Estructura de campo de DatosFiscales
type ERFCDatosFiscales struct {
	RFC      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDomicilioDatosFiscales Estructura de campo de DatosFiscales
type EDomicilioDatosFiscales struct {
	Domicilio Direccion
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EContactosDatosFiscales Estructura de campo de DatosFiscales
type EContactosDatosFiscales struct {
	Contactos Contacto
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Fiscal subestructura de Empresa
type Fiscal struct {
	ERazonSocialDatosFiscales
	ERFCDatosFiscales
	EDomicilioDatosFiscales
	EContactosDatosFiscales
}

//EKeyDatosFactura Estructura de campo de DatosFactura
type EKeyDatosFactura struct {
	Key      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECerDatosFactura Estructura de campo de DatosFactura
type ECerDatosFactura struct {
	Cer      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPemDatosFactura Estructura de campo de DatosFactura
type EPemDatosFactura struct {
	Pem      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Factura subestructura de Empresa
type Factura struct {
	EKeyDatosFactura
	ECerDatosFactura
	EPemDatosFactura
}

//ECorreoCorreoYEnvio Estructura de campo de CorreoYEnvio
type ECorreoCorreoYEnvio struct {
	Correo   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPassCorreoYEnvio Estructura de campo de CorreoYEnvio
type EPassCorreoYEnvio struct {
	Pass     string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoCorreoYEnvio Estructura de campo de CorreoYEnvio
type ETipoCorreoYEnvio struct {
	Tipo     string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPuertoCorreoYEnvio Estructura de campo de CorreoYEnvio
type EPuertoCorreoYEnvio struct {
	Puerto   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECifradoCorreoYEnvio Estructura de campo de CorreoYEnvio
type ECifradoCorreoYEnvio struct {
	Cifrado  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ConfiguracionCorreo subestructura de Empresa
type ConfiguracionCorreo struct {
	ID bson.ObjectId
	ECorreoCorreoYEnvio
	EPassCorreoYEnvio
	ETipoCorreoYEnvio
	EPuertoCorreoYEnvio
	ECifradoCorreoYEnvio
}

//ECalleDomicilio Estructura de campo de Domicilio
type ECalleDomicilio struct {
	Calle    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ENumInteriorDomicilio Estructura de campo de Domicilio
type ENumInteriorDomicilio struct {
	NumInterior string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//ENumExteriorDomicilio Estructura de campo de Domicilio
type ENumExteriorDomicilio struct {
	NumExterior string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EColoniaDomicilio Estructura de campo de Domicilio
type EColoniaDomicilio struct {
	Colonia  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMunicipioDomicilio Estructura de campo de Domicilio
type EMunicipioDomicilio struct {
	Municipio string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EEstadoDomicilio Estructura de campo de Domicilio
type EEstadoDomicilio struct {
	Estado   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPaisDomicilio Estructura de campo de Domicilio
type EPaisDomicilio struct {
	Pais     string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECPDomicilio Estructura de campo de Domicilio
type ECPDomicilio struct {
	CP       string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Direccion subestructura de Empresa
type Direccion struct {
	ECalleDomicilio
	ENumInteriorDomicilio
	ENumExteriorDomicilio
	EColoniaDomicilio
	EMunicipioDomicilio
	EEstadoDomicilio
	EPaisDomicilio
	ECPDomicilio
}

//EAliasContactos Estructura de campo de Contactos
type EAliasContactos struct {
	Alias    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEmailContactos Estructura de campo de Contactos
type EEmailContactos struct {
	Email    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETelefonoContactos Estructura de campo de Contactos
type ETelefonoContactos struct {
	Telefono string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMovilContactos Estructura de campo de Contactos
type EMovilContactos struct {
	Movil    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDomicilioContactos Estructura de campo de Contactos
type EDomicilioContactos struct {
	Domicilio Direccion
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Contacto subestructura de Empresa
type Contacto struct {
	ID bson.ObjectId
	EAliasContactos
	EEmailContactos
	ETelefonoContactos
	EMovilContactos
	EDomicilioContactos
}
