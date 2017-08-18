package FacturacionModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EOperacionFacturacion Estructura de campo de Facturacion
type EOperacionFacturacion struct {
	Operacion bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EDetalleFacturacion Estructura de campo de Facturacion
type EDetalleFacturacion struct {
	Detalle  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EClienteFacturacion Estructura de campo de Facturacion
type EClienteFacturacion struct {
	Cliente  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDomicilioFiscalFacturacion Estructura de campo de Facturacion
type EDomicilioFiscalFacturacion struct {
	DomicilioFiscal Direccion
	IEstatus        bool
	IMsj            string
	Ihtml           template.HTML
}

//EEstatusFacturacion Estructura de campo de Facturacion
type EEstatusFacturacion struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraFacturacion Estructura de campo de Facturacion
type EFechaHoraFacturacion struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Facturacion estructura de Facturacions mongo
type Facturacion struct {
	ID bson.ObjectId
	EOperacionFacturacion
	EDetalleFacturacion
	EClienteFacturacion
	EDomicilioFiscalFacturacion
	EEstatusFacturacion
	EFechaHoraFacturacion
}

//SFacturacion estructura de Facturaciones para la vista
type SFacturacion struct {
	SEstado bool
	SMsj    string
	Facturacion
}

//ECalleDomicilios Estructura de campo de Domicilios
type ECalleDomicilios struct {
	Calle    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ENumInteriorDomicilios Estructura de campo de Domicilios
type ENumInteriorDomicilios struct {
	NumInterior string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//ENumExteriorDomicilios Estructura de campo de Domicilios
type ENumExteriorDomicilios struct {
	NumExterior string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EColoniaDomicilios Estructura de campo de Domicilios
type EColoniaDomicilios struct {
	Colonia  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMunicipioDomicilios Estructura de campo de Domicilios
type EMunicipioDomicilios struct {
	Municipio bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EEstadoDomicilios Estructura de campo de Domicilios
type EEstadoDomicilios struct {
	Estado   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPaisDomicilios Estructura de campo de Domicilios
type EPaisDomicilios struct {
	Pais     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECPDomicilios Estructura de campo de Domicilios
type ECPDomicilios struct {
	CP       string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusDomicilios Estructura de campo de Domicilios
type EEstatusDomicilios struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Direccion subestructura de Facturacion
type Direccion struct {
	ID bson.ObjectId
	ECalleDomicilios
	ENumInteriorDomicilios
	ENumExteriorDomicilios
	EColoniaDomicilios
	EMunicipioDomicilios
	EEstadoDomicilios
	EPaisDomicilios
	ECPDomicilios
	EEstatusDomicilios
}
