package FacturacionModel

import "encoding/xml"

//XMLCFDI Estructura que representa el cuerpo de un Comprobante Fiscal Digital.
type XMLCFDI struct {
	XMLName                xml.Name  `xml:"cfdi:Comprobante"`
	SchemaLocation         string    `xml:"xsi:schemaLocation,attr"`  //"http://www.sat.gob.mx/cfd/3 http://www.sat.gob.mx/sitio_internet/cfd/3/cfdv32.xsd http://www.sat.gob.mx/TimbreFiscalDigital http://www.sat.gob.mx/sitio_internet/TimbreFiscalDigital/TimbreFiscalDigital.xsd http://www.sat.gob.mx/detallista http://www.sat.gob.mx/sitio_internet/cfd/detallista/detallista.xsd http://www.sat.gob.mx/implocal http://www.sat.gob.mx/sitio_internet/cfd/implocal/implocal.xsd http://www.sat.gob.mx/iedu http://www.sat.gob.mx/sitio_internet/cfd/iedu/iedu.xsd http://www.sat.gob.mx/pfic http://www.sat.gob.mx/sitio_internet/cfd/pfic/pfic.xsd"
	SchemaLocationcfdi     string    `xml:"xmlns:cfdi,attr"`          //"http://www.sat.gob.mx/cfd/3"
	SchemaLocationxsi      string    `xml:"xmlns:xsi,attr,omitempty"` //"http://www.w3.org/2001/XMLSchema-instance"
	SchemaLocationxs       string    `xml:"xmlns:xs,attr,omitempty"`  //"http://www.w3.org/2001/XMLSchema"
	SchemaLocationImplocal string    `xml:"xmlns:implocal,attr"`      //"http://www.sat.gob.mx/implocal"
	Version                string    `xml:"version,attr"`
	FechaXML               string    `xml:"fecha,attr"`
	TipoDeComprobante      string    `xml:"tipoDeComprobante,attr"`
	FormaDePago            string    `xml:"formaDePago,attr"`
	SubTotal               float64   `xml:"subTotal,attr"`
	Descuento              float64   `xml:"descuento,attr,omitempty"`
	Moneda                 string    `xml:"Moneda,attr,omitempty"`
	Total                  float64   `xml:"total,attr,"`
	MetodoDePago           string    `xml:"metodoDePago,attr"`
	LugarExpedicion        string    `xml:"LugarExpedicion,attr"`
	NumCtaPago             string    `xml:"NumCtaPago,attr,omitempty"`
	Serie                  string    `xml:"serie,attr,omitempty"`
	Folio                  string    `xml:"folio,attr,omitempty"`
	Certificado            string    `xml:"certificado,attr"`
	NoCertificado          string    `xml:"noCertificado,attr"`
	Sello                  string    `xml:"sello,attr"`
	CFDIEmisor             Emisor    `xml:"cfdi:Emisor"`    //nodo hijo que contiene la informacion fiscal del emisor
	CFDIReceptor           Receptor  `xml:"cfdi:Receptor"`  //nodo hijo que contiene la informacion fiscal del receptor
	CFDIConceptos          Conceptos `xml:"cfdi:Conceptos"` //nodo hijo que contiene la informacion de productos y servicios
	CFDIImpuestos          Impuestos `xml:"cfdi:Impuestos"` //nodo hijo que contiene la informacion referente a impuestos
}

//Emisor Estructura que permite generar el nodo Emisor de un CFD.
type Emisor struct {
	XMLName         xml.Name            `xml:"cfdi:Emisor"`
	Nombre          string              `xml:"nombre,attr,omitempty"`
	RFC             string              `xml:"rfc,attr"`
	DomicilioFiscal CFDIDomicilioFiscal `xml:"cfdi:DomicilioFiscal"`
	ExpedidoEn      CFDIExpedidoEn      `xml:"cfdi:ExpedidoEn,omitempty"`
	CFDIRegimen     Regimen             `xml:"cfdi:RegimenFiscal"`
}

//Receptor Estructura que permite generar el nodo Receptor de un CFD.
type Receptor struct {
	XMLName   xml.Name      `xml:"cfdi:Receptor"`
	Nombre    string        `xml:"nombre,attr,omitempty"`
	RFC       string        `xml:"rfc,attr"`
	Domicilio CFDIDomicilio `xml:"cfdi:Domicilio,omitempty"`
}

//Regimen Nodo hijo de Emisor requerido de un CFD que indica Regimen Fiscal del Contribuyente.
type Regimen struct {
	XMLName xml.Name `xml:"cfdi:RegimenFiscal"`
	Regimen string   `xml:"Regimen,attr"`
}

//CFDIDomicilio Nodo hijo de Receptor opcional para la definición de la ubicación del domicilio del receptor del comprobante fiscal.
type CFDIDomicilio struct {
	XMLName      xml.Name `xml:"cfdi:Domicilio"`
	Calle        string   `xml:"calle,attr"`
	NoExterior   string   `xml:"noExterior,attr"`
	NoInterior   string   `xml:"noInterior,attr,omitempty"`
	Colonia      string   `xml:"colonia,attr"`
	Localidad    string   `xml:"localidad,attr"`
	Municipio    string   `xml:"municipio,attr"`
	Estado       string   `xml:"estado,attr"`
	Pais         string   `xml:"pais,attr"`
	CodigoPostal string   `xml:"codigoPostal,attr"`
}

//CFDIExpedidoEn Nodo hijo de Emisor opcional para la definición de la ubicación donde se expide el CFD.
type CFDIExpedidoEn struct {
	XMLName      xml.Name `xml:"cfdi:ExpedidoEn"`
	Calle        string   `xml:"calle,attr"`
	NoExterior   string   `xml:"noExterior,attr"`
	NoInterior   string   `xml:"noInterior,attr,omitempty"`
	Colonia      string   `xml:"colonia,attr"`
	Localidad    string   `xml:"localidad,attr"`
	Municipio    string   `xml:"municipio,attr"`
	Estado       string   `xml:"estado,attr"`
	Pais         string   `xml:"pais,attr"`
	CodigoPostal string   `xml:"codigoPostal,attr"`
}

// CFDIDomicilioFiscal Nodo hijo de Emisor opcional para la definición del domicilio fiscal del emisor del CFD.
type CFDIDomicilioFiscal struct {
	XMLName      xml.Name `xml:"cfdi:DomicilioFiscal"`
	Calle        string   `xml:"calle,attr"`
	NoExterior   string   `xml:"noExterior,attr"`
	NoInterior   string   `xml:"noInterior,attr,omitempty"`
	Colonia      string   `xml:"colonia,attr"`
	Localidad    string   `xml:"localidad,attr"`
	Municipio    string   `xml:"municipio,attr"`
	Estado       string   `xml:"estado,attr"`
	Pais         string   `xml:"pais,attr"`
	CodigoPostal string   `xml:"codigoPostal,attr"`
}

//Conceptos Nodo hijo de XMLCFDI que representa la lista de cometarios
type Conceptos struct {
	XMLName   xml.Name       `xml:"cfdi:Conceptos"`
	Conceptos []CFDIConcepto `xml:"cfdi:Concepto"`
}
type CFDIConcepto struct {
	XMLName          xml.Name `xml:"cfdi:Concepto"`
	Cantidad         float64  `xml:"cantidad,attr"`
	Unidad           string   `xml:"unidad,attr"`
	NoIdentificacion string   `xml:"noIdentificacion,attr,omitempty"`
	Descripcion      string   `xml:"descripcion,attr"`
	ValorUnitario    float64  `xml:"valorUnitario,attr"`
	Importe          float64  `xml:"importe,attr"`
}
type Impuestos struct {
	XMLName   xml.Name      `xml:"cfdi:Impuestos"`
	Traslados CFDITraslados `xml:"cfdi:Traslados"`
}
type CFDITraslados struct {
	XMLName  xml.Name       `xml:"cfdi:Traslados"`
	Traslado []CFDITraslado `xml:"cfdi:Traslado"`
}
type CFDITraslado struct {
	XMLName  xml.Name `xml:"cfdi:Traslado"`
	Impuesto string   `xml:"impuesto,attr"`
	Tasa     float64  `xml:"tasa,attr"`
	Importe  float64  `xml:"importe,attr"`
}
