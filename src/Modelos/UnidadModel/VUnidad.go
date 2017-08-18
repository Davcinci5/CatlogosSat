package UnidadModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EMagnitudUnidad Estructura de campo de Unidad
type EMagnitudUnidad struct {
	Magnitud string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionUnidad Estructura de campo de Unidad
type EDescripcionUnidad struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EDatosUnidad Estructura de campo de Unidad
type EDatosUnidad struct {
	Datos    DataUnidad
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusUnidad Estructura de campo de Unidad
type EEstatusUnidad struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraUnidad Estructura de campo de Unidad
type EFechaHoraUnidad struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Unidad estructura de Unidads mongo
type Unidad struct {
	ID bson.ObjectId
	EMagnitudUnidad
	EDescripcionUnidad
	EDatosUnidad
	EEstatusUnidad
	EFechaHoraUnidad
}

//SUnidad estructura de Unidades para la vista
type SUnidad struct {
	SEstado bool
	SMsj    string
	Unidad
	SSesion
}

//ENombreDatos Estructura de campo de Datos
type ENombreDatos struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EAbreviaturaDatos Estructura de campo de Datos
type EAbreviaturaDatos struct {
	Abreviatura string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//SSesion estructura de variables de sesion de Usuarios del sistema
type SSesion struct {
	Name          string
	MenuPrincipal template.HTML
	MenuUsr       template.HTML
}

//DataUnidad subestructura de Unidad
type DataUnidad struct {
	ID bson.ObjectId
	ENombreDatos
	EAbreviaturaDatos
}
