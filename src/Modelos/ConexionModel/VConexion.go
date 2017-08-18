package ConexionModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//TestConReturn estructura para la conexion
type TestConReturn struct {
	Estatus      bool
	Mensage      string
	MensageError string
}

//ENombreConexion Estructura de campo de Conexion
type ENombreConexion struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
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

//EFechaHoraConexion Estructura de campo de Conexion
type EFechaHoraConexion struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Conexion estructura de Conexions mongo
type Conexion struct {
	ID bson.ObjectId
	ENombreConexion
	EServidorConexion
	ENombreBDConexion
	EUsuarioBDConexion
	EPassBDConexion
	EFechaHoraConexion
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

//SConexion estructura de Conexiones para la vista
type SConexion struct {
	SEstado bool
	SMsj    string
	Conexion
	SIndex
	SSesion
}
