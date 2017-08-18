package CatalogoModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EClaveCatalogo Estructura de campo de Catalogo
type EClaveCatalogo struct {
	Clave    int64
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ENombreCatalogo Estructura de campo de Catalogo
type ENombreCatalogo struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionCatalogo Estructura de campo de Catalogo
type EDescripcionCatalogo struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EValoresCatalogo Estructura de campo de Catalogo
type EValoresCatalogo struct {
	Valores  Valores
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusCatalogo Estructura de campo de Catalogo
type EEstatusCatalogo struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraCatalogo Estructura de campo de Catalogo
type EFechaHoraCatalogo struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Catalogo estructura de Catalogos mongo
type Catalogo struct {
	ID bson.ObjectId
	EClaveCatalogo
	EDescripcionCatalogo
	ENombreCatalogo
	EValoresCatalogo
	EEstatusCatalogo
	EFechaHoraCatalogo
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

//SCatalogo estructura de Catalogos para la vista
type SCatalogo struct {
	SEstado bool
	SMsj    string
	Catalogo
	SIndex
	SSesion
}

//EValorValores Estructura de campo de Valores
type EValorValores struct {
	Valor    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Valores subestructura de Catalogo
type Valores struct {
	ID bson.ObjectId
	EValorValores
}
