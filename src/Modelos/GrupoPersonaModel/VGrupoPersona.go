package GrupoPersonaModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreGrupoPersona Estructura de campo de GrupoPersona
type ENombreGrupoPersona struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionGrupoPersona Estructura de campo de GrupoPersona
type EDescripcionGrupoPersona struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EEstatusGrupoPersona Estructura de campo de GrupoPersona
type EEstatusGrupoPersona struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraGrupoPersona Estructura de campo de GrupoPersona
type EFechaHoraGrupoPersona struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

// EMiembrosGrupoPersona estructura del Campo Miembros
type EMiembrosGrupoPersona struct {
	Miembros []string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

// ENoMiembrosGrupoPersona estructura del Campo Miembros
type ENoMiembrosGrupoPersona struct {
	Miembros []string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//GrupoPersona estructura de GrupoPersonas mongo
type GrupoPersona struct {
	ID bson.ObjectId
	ENombreGrupoPersona
	EDescripcionGrupoPersona
	EEstatusGrupoPersona
	EFechaHoraGrupoPersona
	EMiembrosGrupoPersona
	ENoMiembrosGrupoPersona
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

//SGrupoPersona estructura de GrupoPersonas para la vista
type SGrupoPersona struct {
	SEstado bool
	SMsj    string
	GrupoPersona
	SIndex
	SSesion
}
