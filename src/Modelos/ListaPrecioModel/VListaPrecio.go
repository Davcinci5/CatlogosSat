package ListaPrecioModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreListaPrecio Estructura de campo de ListaPrecio
type ENombreListaPrecio struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionListaPrecio Estructura de campo de ListaPrecio
type EDescripcionListaPrecio struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EGrupoPListaPrecio Estructura de campo de ListaPrecio
type EGrupoPListaPrecio struct {
	GrupoP   []bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusListaPrecio Estructura de campo de ListaPrecio
type EEstatusListaPrecio struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraListaPrecio Estructura de campo de ListaPrecio
type EFechaHoraListaPrecio struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//ListaPrecio estructura de ListaPrecios mongo
type ListaPrecio struct {
	ID bson.ObjectId
	ENombreListaPrecio
	EDescripcionListaPrecio
	EGrupoPListaPrecio
	EEstatusListaPrecio
	EFechaHoraListaPrecio
}

//SListaPrecio estructura de ListaPrecios para la vista
type SListaPrecio struct {
	SEstado bool
	SMsj    string
	ListaPrecio
}
