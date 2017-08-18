package PermisosUriModel

import (
	"html/template"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EGrupoPermisosUri Estructura de campo de PermisosUri
type EGrupoPermisosUri struct {
	Grupo    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPermisoNegadoPermisosUri Estructura de campo de PermisosUri
type EPermisoNegadoPermisosUri struct {
	PermisoNegado string
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EPermisoAceptadoPermisosUri Estructura de campo de PermisosUri
type EPermisoAceptadoPermisosUri struct {
	PermisoAceptado string
	IEstatus        bool
	IMsj            string
	Ihtml           template.HTML
}

//PermisosUri estructura de PermisosUris mongo
type PermisosUri struct {
	ID bson.ObjectId
	EGrupoPermisosUri
	EPermisoNegadoPermisosUri
	EPermisoAceptadoPermisosUri
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

//SPermisosUri estructura de PermisosUris para la vista
type SPermisosUri struct {
	SEstado bool
	SMsj    string
	PermisosUri
	SIndex
	SSesion
}
