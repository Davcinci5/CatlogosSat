package PromocionModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombrePromocion Estructura de campo de Promocion
type ENombrePromocion struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionPromocion Estructura de campo de Promocion
type EDescripcionPromocion struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EPorcentajeDescPromocion Estructura de campo de Promocion
type EPorcentajeDescPromocion struct {
	PorcentajeDesc Oferta
	IEstatus       bool
	IMsj           string
	Ihtml          template.HTML
}

//EPrecioOfertaPromocion Estructura de campo de Promocion
type EPrecioOfertaPromocion struct {
	PrecioOferta Oferta
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//EOfertaMontoPromocion Estructura de campo de Promocion
type EOfertaMontoPromocion struct {
	OfertaMonto Oferta
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EOfertaPiezaPiezaPromocion Estructura de campo de Promocion
type EOfertaPiezaPiezaPromocion struct {
	OfertaPiezaPieza Oferta
	IEstatus         bool
	IMsj             string
	Ihtml            template.HTML
}

//EOfertaPiezaPorcentajePromocion Estructura de campo de Promocion
type EOfertaPiezaPorcentajePromocion struct {
	OfertaPiezaPorcentaje Oferta
	IEstatus              bool
	IMsj                  string
	Ihtml                 template.HTML
}

//EEstatusPromocion Estructura de campo de Promocion
type EEstatusPromocion struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaInicioPromocion Estructura de campo de Promocion
type EFechaInicioPromocion struct {
	FechaInicio time.Time
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EFechaFinPromocion Estructura de campo de Promocion
type EFechaFinPromocion struct {
	FechaFin time.Time
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Promocion estructura de Promocions mongo
type Promocion struct {
	ID bson.ObjectId
	ENombrePromocion
	EDescripcionPromocion
	EPorcentajeDescPromocion
	EPrecioOfertaPromocion
	EOfertaMontoPromocion
	EOfertaPiezaPiezaPromocion
	EOfertaPiezaPorcentajePromocion
	EEstatusPromocion
	EFechaInicioPromocion
	EFechaFinPromocion
}

//SSesion estructura de variables de sesion de Usuarios del sistema
type SSesion struct {
	Name          string
	MenuPrincipal template.HTML
	MenuUsr       template.HTML
}

//SIndex estructura de variables de index
type SIndex struct {
	SResultados bool
	SRMsj       string
	SCabecera   template.HTML
	SBody       template.HTML
	SPaginacion template.HTML
	SGrupo      template.HTML
}

//SPromocion estructura de Promociones para la vista
type SPromocion struct {
	SEstado bool
	SMsj    string
	Promocion
	SIndex
	SSesion
}

//ECantidadPorcentajeDesc Estructura de campo de PorcentajeDesc
type ECantidadPorcentajeDesc struct {
	Cantidad float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EValorPorcentajeDesc Estructura de campo de PorcentajeDesc
type EValorPorcentajeDesc struct {
	Valor    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Oferta subestructura de Promocion
type Oferta struct {
	ECantidadPorcentajeDesc
	EValorPorcentajeDesc
}
