package ImpuestoModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreImpuesto Estructura de campo de Impuesto
type ENombreImpuesto struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EClasificacionImpuesto Estructura de campo de Impuesto
type EClasificacionImpuesto struct {
	Clasificacion bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//ESubClasificacionImpuesto Estructura de campo de Impuesto
type ESubClasificacionImpuesto struct {
	Clasificacion bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EImpuestosImpuesto Estructura de campo de Impuesto
type EImpuestosImpuesto struct {
	Impuestos DataImpuestos
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EEstatusImpuesto Estructura de campo de Impuesto
type EEstatusImpuesto struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraImpuesto Estructura de campo de Impuesto
type EFechaHoraImpuesto struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Impuesto estructura de Impuestos mongo
type Impuesto struct {
	ID bson.ObjectId
	ENombreImpuesto
	EClasificacionImpuesto
	ESubClasificacionImpuesto
	EImpuestosImpuesto
	EEstatusImpuesto
	EFechaHoraImpuesto
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

//SImpuesto estructura de Impuestos para la vista
type SImpuesto struct {
	SEstado bool
	SMsj    string
	SIhtml  template.HTML
	Impuesto
	SIndex
	SSesion
}

//ENombreImpuestos Estructura de campo de Impuestos
type ENombreImpuestos struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ERetenidoImpuestos Estructura de campo de Impuestos
type ERetenidoImpuestos struct {
	Valor    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETrasladadoImpuestos Estructura de campo de Impuestos
type ETrasladadoImpuestos struct {
	Valor    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ERangoImpuestos Estructura de campo de Impuestos
type ERangoImpuestos struct {
	Rango    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EValorMaxImpuestos Estructura de campo de Impuestos
type EValorMaxImpuestos struct {
	Valor    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EValorMinImpuestos Estructura de campo de Impuestos
type EValorMinImpuestos struct {
	Valor    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoImpuestos Estructura de campo de Impuestos
type ETipoImpuestos struct {
	Tipo     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EUnidadImpuestos Estructura de campo de Impuestos
type EUnidadImpuestos struct {
	Unidad   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraImpuestos Estructura de campo de Impuestos
type EFechaHoraImpuestos struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//DataImpuestos subestructura de Impuesto
type DataImpuestos struct {
	ID bson.ObjectId
	ENombreImpuestos
	ERetenidoImpuestos
	ETrasladadoImpuestos
	ERangoImpuestos
	EValorMaxImpuestos
	EValorMinImpuestos
	ETipoImpuestos
	EUnidadImpuestos
	EFechaHoraImpuestos
}
