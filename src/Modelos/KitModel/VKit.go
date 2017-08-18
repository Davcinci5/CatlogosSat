package KitModel

import (

		"html/template"
		"time"

		"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreKit Estructura de campo de Kit
type ENombreKit struct {
 Nombre   string 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//ECodigoKit Estructura de campo de Kit
type ECodigoKit struct {
 Codigo   string 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//ETipoKit Estructura de campo de Kit
type ETipoKit struct {
 Tipo   bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EAplicaciónKit Estructura de campo de Kit
type EAplicaciónKit struct {
 Aplicación   bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EImagenesKit Estructura de campo de Kit
type EImagenesKit struct {
 Imagenes   []bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EConformacionKit Estructura de campo de Kit
type EConformacionKit struct {
 Conformacion   []Conforma 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EEstatusKit Estructura de campo de Kit
type EEstatusKit struct {
 Estatus   bson.ObjectId 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}


//EFechaHoraKit Estructura de campo de Kit
type EFechaHoraKit struct {
 FechaHora   time.Time 
 	IEstatus bool
 	IMsj     string
 	Ihtml  template.HTML
}
//Kit estructura de Kits mongo
type Kit struct {
 ID    bson.ObjectId   
  ENombreKit 
  ECodigoKit 
  ETipoKit 
  EAplicaciónKit 
  EImagenesKit 
  EConformacionKit 
  EEstatusKit 
  EFechaHoraKit 
}


//SKit estructura de Kits para la vista
type SKit struct {
 SEstado bool
 SMsj	   string
 Kit
}
	

//EAlmacenConformacion Estructura de campo de Conformacion
type EAlmacenConformacion struct {
 Almacen   bson.ObjectId 
 IEstatus bool
 IMsj     string
 Ihtml  template.HTML
}

//EProductosConformacion Estructura de campo de Conformacion
type EProductosConformacion struct {
 Productos   []DataProducto 
 IEstatus bool
 IMsj     string
 Ihtml  template.HTML
}

//Conforma subestructura de Kit
type Conforma struct {
  EAlmacenConformacion 
  EProductosConformacion 
}

//EIDProductoProductos Estructura de campo de Productos
type EIDProductoProductos struct {
 IDProducto   bson.ObjectId 
 IEstatus bool
 IMsj     string
 Ihtml  template.HTML
}

//ECantidadProductos Estructura de campo de Productos
type ECantidadProductos struct {
 Cantidad   float64 
 IEstatus bool
 IMsj     string
 Ihtml  template.HTML
}

//EPrecioProductos Estructura de campo de Productos
type EPrecioProductos struct {
 Precio   float64 
 IEstatus bool
 IMsj     string
 Ihtml  template.HTML
}

//DataProducto subestructura de Kit
type DataProducto struct {
  EIDProductoProductos 
  ECantidadProductos 
  EPrecioProductos 
}

