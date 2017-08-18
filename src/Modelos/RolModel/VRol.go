package RolModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreRol Estructura de campo de Rol
type ENombreRol struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionRol Estructura de campo de Rol
type EDescripcionRol struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EPermisosRol Estructura de campo de Rol
type EPermisosRol struct {
	Permisos []bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusRol Estructura de campo de Rol
type EEstatusRol struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraRol Estructura de campo de Rol
type EFechaHoraRol struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Rol estructura de Rols mongo
type Rol struct {
	ID bson.ObjectId
	ENombreRol
	EDescripcionRol
	EPermisosRol
	EEstatusRol
	EFechaHoraRol
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

//SRol estructura de Roles para la vista
type SRol struct {
	SEstado bool
	SMsj    string
	Rol
	SIndex
	SSesion
}
