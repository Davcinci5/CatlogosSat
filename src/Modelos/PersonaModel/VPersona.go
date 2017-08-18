package PersonaModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombrePersona Estructura de campo de Persona
type ENombrePersona struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoPersona Estructura de campo de Persona
type ETipoPersona struct {
	Tipo     []bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EGruposPersona Estructura de campo de Persona
type EGruposPersona struct {
	Grupos   []bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPredecesorPersona Estructura de campo de Persona
type EPredecesorPersona struct {
	Predecesor bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//ENotificacionPersona Estructura de campo de Persona
type ENotificacionPersona struct {
	Notificacion []bson.ObjectId
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//EEstatusPersona Estructura de campo de Persona
type EEstatusPersona struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraPersona Estructura de campo de Persona
type EFechaHoraPersona struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Persona estructura de Personas mongo
type Persona struct {
	ID bson.ObjectId
	ENombrePersona
	ETipoPersona
	EGruposPersona
	EPredecesorPersona
	ENotificacionPersona
	EEstatusPersona
	EFechaHoraPersona
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

//SPersona estructura de Personas para la vista
type SPersona struct {
	SEstado bool
	SMsj    string
	Persona
	SIndex
	SSesion
}

//EMensajeNotificacion Estructura de campo de Notificacion
type EMensajeNotificacion struct {
	Mensaje  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ELeidoNotificacion Estructura de campo de Notificacion
type ELeidoNotificacion struct {
	Leido    bool
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaOcurrenciaNotificacion Estructura de campo de Notificacion
type EFechaOcurrenciaNotificacion struct {
	FechaOcurrencia time.Time
	IEstatus        bool
	IMsj            string
	Ihtml           template.HTML
}

//EFechaVencimientoNotificacion Estructura de campo de Notificacion
type EFechaVencimientoNotificacion struct {
	FechaVencimiento time.Time
	IEstatus         bool
	IMsj             string
	Ihtml            template.HTML
}

//Notificacion subestructura de Persona
type Notificacion struct {
	ID bson.ObjectId
	EMensajeNotificacion
	ELeidoNotificacion
	EFechaOcurrenciaNotificacion
	EFechaVencimientoNotificacion
}
