package AlmacenModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreAlmacen Estructura de campo de Almacen
type ENombreAlmacen struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoAlmacen Estructura de campo de Almacen
type ETipoAlmacen struct {
	Tipo     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EClasificacionAlmacen Estructura de campo de Almacen
type EClasificacionAlmacen struct {
	Clasificacion bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EPredecesorAlmacen Estructura de campo de Almacen
type EPredecesorAlmacen struct {
	Predecesor bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//ESucesorAlmacen Estructura de campo de Almacen
type ESucesorAlmacen struct {
	Sucesor  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EListaCostoAlmacen Estructura de campo de Almacen
type EListaCostoAlmacen struct {
	ListaCosto bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//EListaPrecioAlmacen Estructura de campo de Almacen
type EListaPrecioAlmacen struct {
	ListaPrecio bson.ObjectId
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//EDireccionAlmacen Estructura de campo de Almacen
type EDireccionAlmacen struct {
	Direccion bson.ObjectId
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EGruposAlmacen Estructura de campo de Almacen
type EGruposAlmacen struct {
	Grupos   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EConexionAlmacen Estructura de campo de Almacen
type EConexionAlmacen struct {
	Conexion DatosConexion
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EEstatusAlmacen Estructura de campo de Almacen
type EEstatusAlmacen struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraAlmacen Estructura de campo de Almacen
type EFechaHoraAlmacen struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Almacen estructura de Almacens mongo
type Almacen struct {
	ID bson.ObjectId
	ENombreAlmacen
	ETipoAlmacen
	EClasificacionAlmacen
	EPredecesorAlmacen
	ESucesorAlmacen
	EListaCostoAlmacen
	EListaPrecioAlmacen
	EDireccionAlmacen
	EGruposAlmacen
	EConexionAlmacen
	EEstatusAlmacen
	EFechaHoraAlmacen
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

//SAlmacen estructura de Almacenes para la vista
type SAlmacen struct {
	SEstado              bool
	SMsj                 string
	AlmacenesDisponibles int
	Almacen
	SIndex
	SSesion
}

//SMovimientoAlmacen dfvdv
type SMovimientoAlmacen struct {
	SEstado   bool
	SMsj      string
	Almacenes []Almacen
	Ihtml     template.HTML
	SSesion
}

//EServidorConexion Estructura de campo de Conexion
type EServidorConexion struct {
	Servidor string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ENombreBDConexion Estructura de campo de Conexion
type ENombreBDConexion struct {
	NombreBD string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EUsuarioBDConexion Estructura de campo de Conexion
type EUsuarioBDConexion struct {
	UsuarioBD string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EPassBDConexion Estructura de campo de Conexion
type EPassBDConexion struct {
	PassBD   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//DatosConexion subestructura de Almacen
type DatosConexion struct {
	ID bson.ObjectId
	EServidorConexion
	ENombreBDConexion
	EUsuarioBDConexion
	EPassBDConexion
}
