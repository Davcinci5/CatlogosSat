package ClienteModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EIDPersonaCliente Estructura de campo de Cliente
type EIDPersonaCliente struct {
	ID bson.ObjectId
	ENombrePersona
	ESexo
	EFechaNacimiento
	ETipoPersona
	EGruposPersona
	EPredecesorPersona
	ENotificacionPersona
	EEstatusPersona
	EFechaHoraPersona
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaNacimiento Estructura que contiene la fecha de nacimiento de la persona
type EFechaNacimiento struct {
	FechaNacimiento string
	IEstatus        bool
	IMsj            string
	Ihtml           template.HTML
}

//ESexo Estructura del sexo persona
type ESexo struct {
	Sexo     string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoCliente Estructura de campo de Cliente
type ETipoCliente struct {
	TipoCliente bson.ObjectId
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//ERFCCliente Estructura de campo de Cliente
type ERFCCliente struct {
	RFC      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDireccionesCliente Estructura de campo de Cliente
type EDireccionesCliente struct {
	Direcciones    Direccion
	NumDirecciones int
	IEstatus       bool
	IMsj           string
	Ihtml          template.HTML
}

//EMediosDeContactoCliente Estructura de campo de Cliente
type EMediosDeContactoCliente struct {
	MediosDeContacto MediosContacto
	IEstatus         bool
	IMsj             string
	Ihtml            template.HTML
}

//EPersonasContactoCliente Estructura de campo de Cliente
type EPersonasContactoCliente struct {
	PersonasContacto PersonaContacto
	NumPerCont       int
	IEstatus         bool
	IMsj             string
	Ihtml            template.HTML
}

//EAlmacenesCliente Estructura de campo de Cliente
type EAlmacenesCliente struct {
	Almacenes bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//ENotificacionesCliente Estructura de campo de Cliente
type ENotificacionesCliente struct {
	Notificaciones Notificacion
	IEstatus       bool
	IMsj           string
	Ihtml          template.HTML
}

//EEstatusCliente Estructura de campo de Cliente
type EEstatusCliente struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraCliente Estructura de campo de Cliente
type EFechaHoraCliente struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Cliente estructura de Clientes mongo
type Cliente struct {
	ID bson.ObjectId
	ETipoCliente
	EIDPersonaCliente
	ERFCCliente
	EDireccionesCliente
	EMediosDeContactoCliente
	EPersonasContactoCliente
	EAlmacenesCliente
	ENotificacionesCliente
	EEstatusCliente
	EFechaHoraCliente
	Persona
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

//SCliente estructura de Clientes para la vista
type SCliente struct {
	SEstado bool
	SMsj    string
	Cliente
	SIndex
	SSesion
}

//ETipoDireccion estructura que contiene el tipo de direccion
type ETipoDireccion struct {
	TipoDireccion bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//ECalleDirecciones Estructura de campo de Direcciones
type ECalleDirecciones struct {
	Calle    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ENumInteriorDirecciones Estructura de campo de Direcciones
type ENumInteriorDirecciones struct {
	NumInterior string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//ENumExteriorDirecciones Estructura de campo de Direcciones
type ENumExteriorDirecciones struct {
	NumExterior string
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EColoniaDirecciones Estructura de campo de Direcciones
type EColoniaDirecciones struct {
	Colonia  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMunicipioDirecciones Estructura de campo de Direcciones
type EMunicipioDirecciones struct {
	Municipio bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EEstadoDirecciones Estructura de campo de Direcciones
type EEstadoDirecciones struct {
	Estado   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPaisDirecciones Estructura de campo de Direcciones
type EPaisDirecciones struct {
	Pais     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECPDirecciones Estructura de campo de Direcciones
type ECPDirecciones struct {
	CP       string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusDirecciones Estructura de campo de Direcciones
type EEstatusDirecciones struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Direccion subestructura de Cliente
type Direccion struct {
	ID bson.ObjectId
	ECalleDirecciones
	ENumInteriorDirecciones
	ENumExteriorDirecciones
	EColoniaDirecciones
	EMunicipioDirecciones
	EEstadoDirecciones
	EPaisDirecciones
	ECPDirecciones
	ETipoDireccion
	EEstatusDirecciones
}

//ECorreosMediosDeContacto Estructura de campo de MediosDeContacto
type ECorreosMediosDeContacto struct {
	Correos  Correo
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETelefonosMediosDeContacto Estructura de campo de MediosDeContacto
type ETelefonosMediosDeContacto struct {
	Telefonos Telefono
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EOtrosMediosDeContacto Estructura de campo de MediosDeContacto
type EOtrosMediosDeContacto struct {
	Otros    string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//MediosContacto subestructura de Cliente
type MediosContacto struct {
	ID bson.ObjectId
	ECorreosMediosDeContacto
	ETelefonosMediosDeContacto
	EOtrosMediosDeContacto
}

//EPrincipalCorreos Estructura de campo de Correos
type EPrincipalCorreos struct {
	Principal string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//ECorreosCorreos Estructura de campo de Correos
type ECorreosCorreos struct {
	Correos  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Correo subestructura de Cliente
type Correo struct {
	EPrincipalCorreos
	ECorreosCorreos
}

//EPrincipalTelefonos Estructura de campo de Telefonos
type EPrincipalTelefonos struct {
	Principal string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//ETelefonosTelefonos Estructura de campo de Telefonos
type ETelefonosTelefonos struct {
	Telefonos string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Telefono subestructura de Cliente
type Telefono struct {
	EPrincipalTelefonos
	ETelefonosTelefonos
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

//Notificacion subestructura de Cliente
type Notificacion struct {
	ID bson.ObjectId
	EMensajeNotificaciones
	ELeidoNotificaciones
	EFechaOcurrenciaNotificaciones
	EFechaVencimientoNotificaciones
}

//ENombreIDPersona Estructura de campo de IDPersona
type ENombreIDPersona struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoIDPersona Estructura de campo de IDPersona
type ETipoIDPersona struct {
	Tipo     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EGruposIDPersona Estructura de campo de IDPersona
type EGruposIDPersona struct {
	Grupos   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPredecesorIDPersona Estructura de campo de IDPersona
type EPredecesorIDPersona struct {
	Predecesor bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//ENotificacionIDPersona Estructura de campo de IDPersona
type ENotificacionIDPersona struct {
	Notificacion bson.ObjectId
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//EEstatusIDPersona Estructura de campo de IDPersona
type EEstatusIDPersona struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraIDPersona Estructura de campo de IDPersona
type EFechaHoraIDPersona struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Persona subestructura de Cliente
type Persona struct {
	ID bson.ObjectId
	ENombreIDPersona
	ETipoIDPersona
	EGruposIDPersona
	EPredecesorIDPersona
	ENotificacionIDPersona
	EEstatusIDPersona
	EFechaHoraIDPersona
}

//ENombrePersonasContacto Estructura de campo de PersonasContacto
type ENombrePersonasContacto struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EDireccionesPersonasContacto Estructura de campo de PersonasContacto
type EDireccionesPersonasContacto struct {
	Direcciones Direccion
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EMediosDeContactoPersonasContacto Estructura de campo de PersonasContacto
type EMediosDeContactoPersonasContacto struct {
	MediosDeContacto MediosContacto
	IEstatus         bool
	IMsj             string
	Ihtml            template.HTML
}

//EAlmacenesPersonasContacto Estructura de campo de PersonasContacto
type EAlmacenesPersonasContacto struct {
	Almacenes Almacen
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EEstatusPersonasContacto Estructura de campo de PersonasContacto
type EEstatusPersonasContacto struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//PersonaContacto subestructura de Cliente
type PersonaContacto struct {
	ID bson.ObjectId
	ENombrePersonasContacto
	EDireccionesPersonasContacto
	EMediosDeContactoPersonasContacto
	EAlmacenesPersonasContacto
	EEstatusPersonasContacto
}

//EIDContactoAlmacenes Estructura de campo de Almacenes
type EIDContactoAlmacenes struct {
	IDContacto bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//EIDDireccionAlmacenes Estructura de campo de Almacenes
type EIDDireccionAlmacenes struct {
	IDDireccion bson.ObjectId
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EIDAlmacenAlmacenes Estructura de campo de Almacenes
type EIDAlmacenAlmacenes struct {
	IDAlmacen bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Almacen subestructura de Cliente
type Almacen struct {
	EIDContactoAlmacenes
	EIDDireccionAlmacenes
	EIDAlmacenAlmacenes
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
