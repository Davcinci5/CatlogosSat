package PuntoVentaModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EOperacionPuntoVenta Estructura de campo de PuntoVenta
type EOperacionPuntoVenta struct {
	Operacion bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EUsuarioPuntoVenta Estructura de campo de PuntoVenta
type EUsuarioPuntoVenta struct {
	Usuario  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEquipoPuntoVenta Estructura de campo de PuntoVenta
type EEquipoPuntoVenta struct {
	Equipo   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECodigoPuntoVenta Estructura de campo de PuntoVenta
type ECodigoPuntoVenta struct {
	Codigo   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECarritoPuntoVenta Estructura de campo de PuntoVenta
type ECarritoPuntoVenta struct {
	Carrito  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EResumenPuntoVenta Estructura de campo de PuntoVenta
type EResumenPuntoVenta struct {
	Resumen  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusPuntoVenta Estructura de campo de PuntoVenta
type EEstatusPuntoVenta struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraPuntoVenta Estructura de campo de PuntoVenta
type EFechaHoraPuntoVenta struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//PuntoVenta estructura de PuntoVentas mongo
type PuntoVenta struct {
	ID bson.ObjectId
	EOperacionPuntoVenta
	EUsuarioPuntoVenta
	EEquipoPuntoVenta
	ECodigoPuntoVenta
	ECarritoPuntoVenta
	EResumenPuntoVenta
	EEstatusPuntoVenta
	EFechaHoraPuntoVenta
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

//SPuntoVenta estructura de PuntoVentas para la vista
type SPuntoVenta struct {
	SEstado bool
	SMsj    string
	PuntoVenta
	SIndex
	SSesion
}

//SDataProducto estructura de Productos para la vista
type SDataProducto struct {
	ID           string
	SEstado      bool
	SMsj         string
	SIhtml       template.HTML
	SCalculadora template.HTML
	SElastic     bool
	SBusqueda    template.HTML
	SIndex
	SSesion
}
