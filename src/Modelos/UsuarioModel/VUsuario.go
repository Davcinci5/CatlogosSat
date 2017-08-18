package UsuarioModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EIDPersonaUsuario Estructura de campo de Usuario
type EIDPersonaUsuario struct {
	IDPersona bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EUsuarioUsuario Estructura de campo de Usuario
type EUsuarioUsuario struct {
	Usuario  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECredencialesUsuario Estructura de campo de Usuario
type ECredencialesUsuario struct {
	Pin         Pin
	Contraseña  Contraseña
	CodigoBarra CodigoBarra
	Huella      Huella
}

//Pin Estructura para definir la credencial Pin
type Pin struct {
	Pin      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Contraseña Estructura para definir la credencial Contraseña
type Contraseña struct {
	Contraseña string
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//CodigoBarra Estructura para definir la credencial CodigoBarra
type CodigoBarra struct {
	CodigoBarra string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//Huella Estructura para definir la credencial Huella
type Huella struct {
	Huella   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ERolesUsuario Estructura de campo de Usuario
type ERolesUsuario struct {
	Roles    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMediosDeContactoUsuario Estructura de campo de Usuario
type EMediosDeContactoUsuario struct {
	ECorreos     Correos
	ETelefonos   Telefonos
	EOtrosMedios OtrosMedios
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//Correos Estructura para los campos correos de Usuario
type Correos struct {
	Principal string
	Correos   []string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Telefonos Estructura para los campos correos de Usuario
type Telefonos struct {
	Principal string
	Telefonos []string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//OtrosMedios Estructura para los campos correos de Usuario
type OtrosMedios struct {
	OtrosMedios []string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//ECajasUsuario Estructura de campo de Usuario
type ECajasUsuario struct {
	Cajas    []string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ENotificacionesUsuario Estructura de campo de Usuario
type ENotificacionesUsuario struct {
	Notificaciones Notificacion
	IEstatus       bool
	IMsj           string
	Ihtml          template.HTML
}

//EEstatusUsuario Estructura de campo de Usuario
type EEstatusUsuario struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraUsuario Estructura de campo de Usuario
type EFechaHoraUsuario struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Usuario estructura de Usuarios mongo
type Usuario struct {
	ID bson.ObjectId
	EPersonaUsuario
	EUsuarioUsuario
	ECredencialesUsuario
	ERolesUsuario
	EMediosDeContactoUsuario
	ECajasUsuario
	ENotificacionesUsuario
	EEstatusUsuario
	EFechaHoraUsuario
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

//SUsuario estructura de Usuarios para la vista
type SUsuario struct {
	SEstado bool
	SMsj    string
	Usuario
	SIndex
	SSesion
}

//EMensajeNotificaciones Estructura de campo de Notificaciones
type EMensajeNotificaciones struct {
	Mensaje  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ELeidoNotificaciones Estructura de campo de Notificaciones
type ELeidoNotificaciones struct {
	Leido    bool
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaOcurrenciaNotificaciones Estructura de campo de Notificaciones
type EFechaOcurrenciaNotificaciones struct {
	FechaOcurrencia time.Time
	IEstatus        bool
	IMsj            string
	Ihtml           template.HTML
}

//EFechaVencimientoNotificaciones Estructura de campo de Notificaciones
type EFechaVencimientoNotificaciones struct {
	FechaVencimiento time.Time
	IEstatus         bool
	IMsj             string
	Ihtml            template.HTML
}

//Notificacion subestructura de Usuario
type Notificacion struct {
	ID bson.ObjectId
	EMensajeNotificaciones
	ELeidoNotificaciones
	EFechaOcurrenciaNotificaciones
	EFechaVencimientoNotificaciones
}

//ENombrePersona Estructura de campo de IDPersona
type ENombrePersona struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoPersona Estructura de campo de IDPersona
type ETipoPersona struct {
	Tipo     []string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EGruposPersona Estructura de campo de IDPersona
type EGruposPersona struct {
	Grupos   []string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPredecesorPersona Estructura de campo de IDPersona
type EPredecesorPersona struct {
	Predecesor bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//ENotificacionPersona Estructura de campo de IDPersona
type ENotificacionPersona struct {
	Notificacion string
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//EEstatusPersona Estructura de campo de IDPersona
type EEstatusPersona struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraPersona Estructura de campo de IDPersona
type EFechaHoraPersona struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EPersonaUsuario subestructura de Usuario
type EPersonaUsuario struct {
	ID bson.ObjectId
	ENombrePersona
	ETipoPersona
	EGruposPersona
	EPredecesorPersona
	ENotificacionPersona
	EEstatusPersona
	EFechaHoraPersona
}
