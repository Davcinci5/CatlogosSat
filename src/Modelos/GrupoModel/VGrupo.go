package GrupoModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreGrupo Estructura de campo de Grupo
type ENombreGrupo struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionGrupo Estructura de campo de Grupo
type EDescripcionGrupo struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EPermiteVenderGrupo Estructura de campo de Grupo
type EPermiteVenderGrupo struct {
	PermiteVender bool
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//ETipoGrupo Estructura de campo de Grupo
type ETipoGrupo struct {
	Tipo     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMiembrosGrupo Estructura de campo de Grupo
type EMiembrosGrupo struct {
	Miembros []bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusGrupo Estructura de campo de Grupo
type EEstatusGrupo struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraGrupo Estructura de campo de Grupo
type EFechaHoraGrupo struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EFechaEdicionGrupo Estructura de campo de Grupo
type EFechaEdicionGrupo struct {
	FechaEdicion time.Time
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//Grupo estructura de Grupos mongo
type Grupo struct {
	ID bson.ObjectId
	ENombreGrupo
	EDescripcionGrupo
	EPermiteVenderGrupo
	ETipoGrupo
	EMiembrosGrupo
	EEstatusGrupo
	EFechaHoraGrupo
	EFechaEdicionGrupo
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

//SGrupo estructura de Grupos para la vista
type SGrupo struct {
	SEstado bool
	SMsj    string
	Grupo
	SIndex
	SSesion
}
