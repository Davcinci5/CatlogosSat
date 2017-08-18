package ProductoModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//ENombreProducto Estructura de campo de Producto
type ENombreProducto struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ECodigosProducto Estructura de campo de Producto
type ECodigosProducto struct {
	Codigos  Codigo
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETipoProducto Estructura de campo de Producto
type ETipoProducto struct {
	Tipo     bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EImagenesProducto Estructura de campo de Producto
type EImagenesProducto struct {
	Imagenes []bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EUnidadProducto Estructura de campo de Producto
type EUnidadProducto struct {
	Unidad   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMmvProducto Estructura de campo de Producto
type EMmvProducto struct {
	Mmv      string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EVentaFraccionProducto Estructura de campo de Producto
type EVentaFraccionProducto struct {
	VentaFraccion bool
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EEtiquetasProducto Estructura de campo de Producto
type EEtiquetasProducto struct {
	Etiquetas string
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//EEstatusProducto Estructura de campo de Producto
type EEstatusProducto struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraProducto Estructura de campo de Producto
type EFechaHoraProducto struct {
	FechaHora time.Time
	IEstatus  bool
	IMsj      string
	Ihtml     template.HTML
}

//Producto estructura de Productos mongo
type Producto struct {
	ID bson.ObjectId
	ENombreProducto
	ECodigosProducto
	ETipoProducto
	EImagenesProducto
	EUnidadProducto
	EMmvProducto
	EVentaFraccionProducto
	EEtiquetasProducto
	EEstatusProducto
	EFechaHoraProducto
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

//SProducto estructura de Productos para la vista
type SProducto struct {
	SEstado bool
	SMsj    string
	Producto
	SIndex
	SSesion
}

//SDataProducto estructura de Productos para la vista
type SDataProducto struct {
	SEstado bool
	SMsj    string
	SIhtml  template.HTML
}

//EClavesCodigos Estructura de campo de Codigos
type EClavesCodigos struct {
	Claves   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EValoresCodigos Estructura de campo de Codigos
type EValoresCodigos struct {
	Valores  string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//Codigo subestructura de Producto
type Codigo struct {
	EClavesCodigos
	EValoresCodigos
}
