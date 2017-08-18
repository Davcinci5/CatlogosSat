package ListaCostoModel

import (

		"html/template"
		"time"

		"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreListaCosto Estructura de campo de ListaCosto
type ENombreListaCosto struct {
 Nombre   string 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EDescripcionListaCosto Estructura de campo de ListaCosto
type EDescripcionListaCosto struct {
 Descripcion   string 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EGrupoPListaCosto Estructura de campo de ListaCosto
type EGrupoPListaCosto struct {
 GrupoP   []bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EEstatusListaCosto Estructura de campo de ListaCosto
type EEstatusListaCosto struct {
 Estatus   bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EFechaHoraListaCosto Estructura de campo de ListaCosto
type EFechaHoraListaCosto struct {
 FechaHora   time.Time 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}
//ListaCosto estructura de ListaCostos mongo
type ListaCosto struct {
 ID    bson.ObjectId   
  ENombreListaCosto 
  EDescripcionListaCosto 
  EGrupoPListaCosto 
  EEstatusListaCosto 
  EFechaHoraListaCosto 
}


//SListaCosto estructura de ListaCosto para la vista
type SListaCosto struct {
 SEstado bool
 SMsj	   string
 ListaCosto
}
	

