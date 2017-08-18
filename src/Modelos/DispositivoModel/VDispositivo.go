package DispositivoModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreDispositivo Estructura de campo de Dispositivo
type ENombreDispositivo struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDescripcionDispositivo Estructura de campo de Dispositivo
type EDescripcionDispositivo struct {
	Descripcion string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EPredecesorDispositivo Estructura de campo de Dispositivo
type EPredecesorDispositivo struct {
	Predecesor bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//EMacDispositivo Estructura de campo de Dispositivo
type EMacDispositivo struct {
	Mac      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECamposDispositivo Estructura de campo de Dispositivo
type ECamposDispositivo struct {
	Campos   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusDispositivo Estructura de campo de Dispositivo
type EEstatusDispositivo struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraDispositivo Estructura de campo de Dispositivo
type EFechaHoraDispositivo struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Dispositivo estructura de Dispositivos mongo
type Dispositivo struct {
	ID bson.ObjectId
	ENombreDispositivo
	EDescripcionDispositivo
	EPredecesorDispositivo
	EMacDispositivo
	ECamposDispositivo
	EEstatusDispositivo
	EFechaHoraDispositivo
}

//SDispositivo estructura de Dispositivos para la vista
type SDispositivo struct {
	SEstado bool
	SMsj    string
	Dispositivo
}
