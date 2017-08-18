package FacturacionModel

import "time"

//CabeceraSQL obtiene los datos del movimiento desde sql
type CabeceraSQL struct {
	Fecha           time.Time
	TipoComprobante string
	FormaPago       string
	Subtotal        float64
	Descuento       float64
	Moneda          string
	Total           float64
	MetodoPago      string
	LugarExpedicion string
	NumCtaPago      string
	NumCte          string
	NumAlmacen      string
}

//ConceptosSQL obtiene los conceptos del movimiento desde sql
type ConceptosSQL struct {
	Cantidad         float64
	Unidad           string
	NoIdentificacion string
	NomDTM           string
	ValorUnitario    float64
	Importe          float64
}

// ClienteSQL utilizado para obtener los datos de cliente sql
type ClienteSQL struct {
	RFC       string
	Nombre    string
	Domicilio DomicilioSQL
}

//DomicilioSQL utilizado para obtener los datos de una direccion o domicilio Fiscal sql
type DomicilioSQL struct {
	Calle        string
	NoExterior   string
	NoInterior   string
	Colonia      string
	Localidad    string
	Municipio    string
	Estado       string
	Pais         string
	CodigoPostal string
}
