package MediosPagoModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreMediosPago Estructura de campo de MediosPago
type ENombreMediosPago struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionMediosPago Estructura de campo de MediosPago
type EDescripcionMediosPago struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//ECodigoSatMediosPago Estructura de campo de MediosPago
type ECodigoSatMediosPago struct {
	CodigoSat string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//ETipoMediosPago Estructura de campo de MediosPago
type ETipoMediosPago struct {
	Tipo     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EComisionMediosPago Estructura de campo de MediosPago
type EComisionMediosPago struct {
	Comision float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECambioMediosPago Estructura de campo de MediosPago
type ECambioMediosPago struct {
	Cambio   bool
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusMediosPago Estructura de campo de MediosPago
type EEstatusMediosPago struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraMediosPago Estructura de campo de MediosPago
type EFechaHoraMediosPago struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//MediosPago estructura de MediosPagos mongo
type MediosPago struct {
	ID bson.ObjectId
	ENombreMediosPago
	EDescripcionMediosPago
	ECodigoSatMediosPago
	ETipoMediosPago
	EComisionMediosPago
	ECambioMediosPago
	EEstatusMediosPago
	EFechaHoraMediosPago
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

//SMediosPago estructura de MediosPago para la vista
type SMediosPago struct {
	SEstado bool
	SMsj    string
	MediosPago
	SIndex
	SSesion
}
