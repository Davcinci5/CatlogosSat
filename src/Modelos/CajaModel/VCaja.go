package CajaModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EUsuarioCaja Estructura de campo de Caja
type EUsuarioCaja struct {
	Usuario  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECajaCaja Estructura de campo de Caja
type ECajaCaja struct {
	Caja     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECargoCaja Estructura de campo de Caja
type ECargoCaja struct {
	Cargo    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EAbonoCaja Estructura de campo de Caja
type EAbonoCaja struct {
	Abono    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ESaldoCaja Estructura de campo de Caja
type ESaldoCaja struct {
	Saldo    float64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EOperacionCaja Estructura de campo de Caja
type EOperacionCaja struct {
	Operacion bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EEstatusCaja Estructura de campo de Caja
type EEstatusCaja struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraCaja Estructura de campo de Caja
type EFechaHoraCaja struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Caja estructura de Cajas mongo
type Caja struct {
	ID bson.ObjectId
	EUsuarioCaja
	ECajaCaja
	ECargoCaja
	EAbonoCaja
	ESaldoCaja
	EOperacionCaja
	EEstatusCaja
	EFechaHoraCaja
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

//SCaja estructura de Cajas para la vista
type SCaja struct {
	SEstado bool
	SMsj    string
	Caja
	SIndex
	SSesion
}
