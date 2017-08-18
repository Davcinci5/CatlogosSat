package BugModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ETipoBug Estructura de campo de Bug
type ETipoBug struct {
	Tipo     string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETituloBug Estructura de campo de Bug
type ETituloBug struct {
	Titulo   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionBug Estructura de campo de Bug
type EDescripcionBug struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EUsuarioBug Estructura de campo de Bug
type EUsuarioBug struct {
	Usuario  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMetodoBug Estructura de campo de Bug
type EMetodoBug struct {
	Metodo   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEsAjaxBug Estructura de campo de Bug
type EEsAjaxBug struct {
	EsAjax   bool
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusPeticionBug Estructura de campo de Bug
type EEstatusPeticionBug struct {
	EstatusPeticion string
	IEstatus        bool
	IMsj            string
	Ihtml           template.HTML
}

//EEstatusBug Estructura de campo de Bug
type EEstatusBug struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ERutaBug Estructura de campo de Bug
type ERutaBug struct {
	Ruta     string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraBug Estructura de campo de Bug
type EFechaHoraBug struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Bug estructura de Bugs mongo
type Bug struct {
	ID bson.ObjectId
	ETipoBug
	ETituloBug
	EDescripcionBug
	EUsuarioBug
	EMetodoBug
	EEsAjaxBug
	EEstatusPeticionBug
	EEstatusBug
	ERutaBug
	EFechaHoraBug
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

//SBug estructura de Bugs para la vista
type SBug struct {
	SEstado  bool
	SMsj     string
	SFuncion string
	Bug
	SIndex
	SSesion
}
