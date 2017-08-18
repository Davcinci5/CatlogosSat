package OperacionModel

import (
	"html/template"
	"time"

	"gopkg.in/mgo.v2/bson"
)

//#########################< ESTRUCTURAS >##############################

//EInventarioPostgres estructura que será usada para insertar y leer datos en postgres
type EInventarioPostgres struct {
	IDProducto string
	Existencia float64
	Estatus    string
	Costo      float64
	Precio     float64
	Encontrado bool
}

//SKardex estructura del kardex para mostrar en la vista
type SKardex struct {
	SEstado bool
	SMsj    string
	Kardex  []Kardex
	SSesion
}

//Kardex estructura de kardex en postgres
type Kardex struct {
	NombreAlmacen  string
	IDMovimiento   string
	IDProducto     string
	Cantidad       float64
	Costo          float64
	Precio         float64
	ImpuestoTotal  float64
	DescuentoTotal float64
	TipoOperacion  string
	Existencia     float64
	FechaHora      time.Time
}

//EUsuarioOrigenOperacion Estructura de campo de Operacion
type EUsuarioOrigenOperacion struct {
	UsuarioOrigen bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EUsuarioDestinoOperacion Estructura de campo de Operacion
type EUsuarioDestinoOperacion struct {
	UsuarioDestino bson.ObjectId
	IEstatus       bool
	IMsj           string
	Ihtml          template.HTML
}

//EFechaHoraRegistroOperacion Estructura de campo de Operacion
type EFechaHoraRegistroOperacion struct {
	FechaHoraRegistro time.Time
	IEstatus          bool
	IMsj              string
	Ihtml             template.HTML
}

//ETipoOperacionOperacion Estructura de campo de Operacion
type ETipoOperacionOperacion struct {
	TipoOperacion bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EEstatusOperacion Estructura de campo de Operacion
type EEstatusOperacion struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPredecesorOperacion Estructura de campo de Operacion
type EPredecesorOperacion struct {
	Predecesor bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//EMovimientosOperacion Estructura de campo de Operacion
type EMovimientosOperacion struct {
	Movimientos []Movimiento
	IEstatus    bool
	IMsj        string
	Ihtml       template.HTML
}

//Operacion estructura de Operacions mongo
type Operacion struct {
	ID bson.ObjectId
	EUsuarioOrigenOperacion
	EUsuarioDestinoOperacion
	EFechaHoraRegistroOperacion
	ETipoOperacionOperacion
	EEstatusOperacion
	EPredecesorOperacion
	EMovimientosOperacion
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

//SOperacion estructura de Operaciones para la vista
type SOperacion struct {
	SEstado bool
	SMsj    string
	Operacion
	SSesion
	SIndex
}

//SSesion estructura de variables de sesion de Usuarios del sistema
type SSesion struct {
	Name          string
	MenuPrincipal template.HTML
	MenuUsr       template.HTML
}

//EIDMovimientoMovimientos Estructura de campo de Movimientos
type EIDMovimientoMovimientos struct {
	IDMovimiento bson.ObjectId
	IEstatus     bool
	IMsj         string
	Ihtml        template.HTML
}

//EAlmacenOrigenMovimientos Estructura de campo de Movimientos
type EAlmacenOrigenMovimientos struct {
	AlmacenOrigen bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EAlmacenDestinoMovimientos Estructura de campo de Movimientos
type EAlmacenDestinoMovimientos struct {
	AlmacenDestino bson.ObjectId
	IEstatus       bool
	IMsj           string
	Ihtml          template.HTML
}

//ERutaMovimientos Estructura de campo de Movimientos
type ERutaMovimientos struct {
	Ruta     []bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EPredecesorMovimientos Estructura de campo de Movimientos
type EPredecesorMovimientos struct {
	Predecesor bson.ObjectId
	IEstatus   bool
	IMsj       string
	Ihtml      template.HTML
}

//EEstatusMovimientos Estructura de campo de Movimientos
type EEstatusMovimientos struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//ETransaccionesMovimientos Estructura de campo de Movimientos
type ETransaccionesMovimientos struct {
	Transacciones []Transaccion
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//Movimiento subestructura de Operacion
type Movimiento struct {
	EIDMovimientoMovimientos
	EAlmacenOrigenMovimientos
	EAlmacenDestinoMovimientos
	ERutaMovimientos
	EPredecesorMovimientos
	EEstatusMovimientos
	ETransaccionesMovimientos
}

//EIDTransaccionTransacciones Estructura de campo de Transacciones
type EIDTransaccionTransacciones struct {
	IDTransaccion bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EAlmacenOrigenTransacciones Estructura de campo de Transacciones
type EAlmacenOrigenTransacciones struct {
	AlmacenOrigen bson.ObjectId
	IEstatus      bool
	IMsj          string
	Ihtml         template.HTML
}

//EAlmacenDestinoTransacciones Estructura de campo de Transacciones
type EAlmacenDestinoTransacciones struct {
	AlmacenDestino bson.ObjectId
	IEstatus       bool
	IMsj           string
	Ihtml          template.HTML
}

//EEstatusTransacciones Estructura de campo de Transacciones
type EEstatusTransacciones struct {
	Estatus  bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EMotivoTransacciones Estructura de campo de Transacciones
type EMotivoTransacciones struct {
	Motivo   bson.ObjectId
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//EFechaHoraAplicacionTransacciones Estructura de campo de Transacciones
type EFechaHoraAplicacionTransacciones struct {
	FechaHoraAplicacion time.Time
	IEstatus            bool
	IMsj                string
	Ihtml               template.HTML
}

//Transaccion subestructura de Operacion
type Transaccion struct {
	EIDTransaccionTransacciones
	EAlmacenOrigenTransacciones
	EAlmacenDestinoTransacciones
	EEstatusTransacciones
	EMotivoTransacciones
	EFechaHoraAplicacionTransacciones
}
