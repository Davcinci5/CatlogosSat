package GrupoAlmacenModel

import (

		"html/template"
		"time"

		"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreGrupoAlmacen Estructura de campo de GrupoAlmacen
type ENombreGrupoAlmacen struct {
 Nombre   string 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EDescripcionGrupoAlmacen Estructura de campo de GrupoAlmacen
type EDescripcionGrupoAlmacen struct {
 Descripcion   string 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EPermiteVenderGrupoAlmacen Estructura de campo de GrupoAlmacen
type EPermiteVenderGrupoAlmacen struct {
 PermiteVender   bool 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EMiembrosGrupoAlmacen Estructura de campo de GrupoAlmacen
type EMiembrosGrupoAlmacen struct {
 Miembros   []bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EEstatusGrupoAlmacen Estructura de campo de GrupoAlmacen
type EEstatusGrupoAlmacen struct {
 Estatus   bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EFechaHoraGrupoAlmacen Estructura de campo de GrupoAlmacen
type EFechaHoraGrupoAlmacen struct {
 FechaHora   time.Time 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}
//GrupoAlmacen estructura de GrupoAlmacens mongo
type GrupoAlmacen struct {
 ID    bson.ObjectId   
  ENombreGrupoAlmacen 
  EDescripcionGrupoAlmacen 
  EPermiteVenderGrupoAlmacen 
  EMiembrosGrupoAlmacen 
  EEstatusGrupoAlmacen 
  EFechaHoraGrupoAlmacen 
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


//SGrupoAlmacen estructura de GrupoPersonas para la vista
type SGrupoAlmacen struct {
 SEstado bool
 SMsj	   string
 GrupoAlmacen
SIndex
SSesion
}
	

