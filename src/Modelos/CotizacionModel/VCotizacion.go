package CotizacionModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EOperacionCotizacion Estructura de campo de Cotizacion
type EOperacionCotizacion struct {
	Operacion bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EUsuarioCotizacion Estructura de campo de Cotizacion
type EUsuarioCotizacion struct {
	Usuario  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEquipoCotizacion Estructura de campo de Cotizacion
type EEquipoCotizacion struct {
	Equipo   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EClienteCotizacion Estructura de campo de Cotizacion
type EClienteCotizacion struct {
	Cliente       bson.ObjectId
	NombreCliente string
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EGrupoCotizacion Estructura de campo de Cotizacion
type EGrupoCotizacion struct {
	Grupo    bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ENombreCotizacion Estructura de campo de Cotizacion
type ENombreCotizacion struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETelefonoCotizacion Estructura de campo de Cotizacion
type ETelefonoCotizacion struct {
	Telefono string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECorreoCotizacion Estructura de campo de Cotizacion
type ECorreoCotizacion struct {
	Correo   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFormaDePagoCotizacion Estructura de campo de Cotizacion
type EFormaDePagoCotizacion struct {
	FormaDePago bson.ObjectId
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EFormaDeEnvíoCotizacion Estructura de campo de Cotizacion
type EFormaDeEnvíoCotizacion struct {
	FormaDeEnvío bson.ObjectId
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//EBuscarCotizacion Estructura de campo de Cotizacion
type EBuscarCotizacion struct {
	Buscar   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EListaCotizacion Estructura de campo de Cotizacion
type EListaCotizacion struct {
	Lista    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECarritoCotizacion Estructura de campo de Cotizacion
type ECarritoCotizacion struct {
	Carrito  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EResumenCotizacion Estructura de campo de Cotizacion
type EResumenCotizacion struct {
	Resumen  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusCotizacion Estructura de campo de Cotizacion
type EEstatusCotizacion struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaInicioCotizacion Estructura de campo de Cotizacion
type EFechaInicioCotizacion struct {
	FechaInicio time.Time
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EFechaFinCotizacion Estructura de campo de Cotizacion
type EFechaFinCotizacion struct {
	FechaFin time.Time
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Cotizacion estructura de Cotizacions mongo
type Cotizacion struct {
	ID bson.ObjectId
	EOperacionCotizacion
	EUsuarioCotizacion
	EEquipoCotizacion
	EClienteCotizacion
	EGrupoCotizacion
	ENombreCotizacion
	ETelefonoCotizacion
	ECorreoCotizacion
	EFormaDePagoCotizacion
	EFormaDeEnvíoCotizacion
	EBuscarCotizacion
	EListaCotizacion
	ECarritoCotizacion
	EResumenCotizacion
	EEstatusCotizacion
	EFechaInicioCotizacion
	EFechaFinCotizacion
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

//SCotizacion estructura de Cotizaciones para la vista
type SCotizacion struct {
	SEstado bool
	SMsj    string
	Cotizacion
	SIndex
	SSesion
}

//SDataCliente estructura de Productos para la vista
type SDataCliente struct {
	SEstado bool
	SMsj    string
	SIhtml  template.HTML
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
