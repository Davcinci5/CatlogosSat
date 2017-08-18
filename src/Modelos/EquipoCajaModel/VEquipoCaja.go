package EquipoCajaModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreEquipoCaja Estructura de campo de EquipoCaja
type ENombreEquipoCaja struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionEquipoCaja Estructura de campo de EquipoCaja
type EDescripcionEquipoCaja struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EUsuarioEquipoCaja Estructura de campo de EquipoCaja
type EUsuarioEquipoCaja struct {
	Usuario  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDispositivoEquipoCaja Estructura de campo de EquipoCaja
type EDispositivoEquipoCaja struct {
	Dispositivo bson.ObjectId
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EEstatusEquipoCaja Estructura de campo de EquipoCaja
type EEstatusEquipoCaja struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraEquipoCaja Estructura de campo de EquipoCaja
type EFechaHoraEquipoCaja struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EquipoCaja estructura de EquipoCajas mongo
type EquipoCaja struct {
	ID bson.ObjectId
	ENombreEquipoCaja
	EDescripcionEquipoCaja
	EUsuarioEquipoCaja
	EDispositivoEquipoCaja
	EEstatusEquipoCaja
	EFechaHoraEquipoCaja
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

//SEquipoCaja estructura de EquipoCajas para la vista
type SEquipoCaja struct {
	SEstado bool
	SMsj    string
	EquipoCaja
	SIndex
	SSesion
}
