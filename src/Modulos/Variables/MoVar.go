package MoVar

import (
	config "github.com/robfig/config"
)

const (

	//############# ARCHIVOS LOCALES ######################################

	//FileConfigName contiene el nombre del archivo CFG
	FileConfigName = "Acfg.cfg"

	//################ SECCIONES CFG  ######################################

	//SecDefault nombre de la seccion default del servidor en CFG
	SecDefault = "DEFAULT"
	//SecMongo nombre de la seccion de mongo en CFG
	SecMongo = "CONFIG_DB_MONGO"
	//SecPsql nombre de la seccion de postgresql en cfg
	SecPsql = "CONFIG_DB_POSTGRES"
	//SecElastic nombre de la seccion de postgresql en cfg
	SecElastic = "CONFIG_DB_ELASTIC"

	//############# COLUMNAS DE ALMACEN PSQL ######################################

	//############# COLECCIONES MONGO ######################################

	//Mongo---------------> Catalogo

	//ColeccionCatalogo nombre de la coleccion de Catalogo en mongo
	ColeccionCatalogo = "Catalogo"

	//Mongo---------------> Unidad

	//ColeccionUnidad nombre de la coleccion de Unidad en mongo
	ColeccionUnidad = "Unidad"

	//Mongo---------------> Cliente

	//ColeccionCliente nombre de la coleccion de Cliente en mongo
	ColeccionCliente = "Cliente"

	//Mongo---------------> Producto

	//ColeccionProducto nombre de la coleccion de Producto en mongo
	ColeccionProducto = "Producto"

	//Mongo---------------> Almacen

	//ColeccionAlmacen nombre de la coleccion de Almacen en mongo
	ColeccionAlmacen = "Almacen"

	//Mongo---------------> GrupoPersona

	//ColeccionGrupoPersona nombre de la coleccion de GrupoPersona en mongo
	ColeccionGrupoPersona = "GrupoPersona"

	//Mongo---------------> ListaCosto

	//ColeccionListaCosto nombre de la coleccion de ListaCosto en mongo
	ColeccionListaCosto = "ListaCosto"

	//Mongo---------------> ListaPrecio

	//ColeccionListaPrecio nombre de la coleccion de ListaPrecio en mongo
	ColeccionListaPrecio = "ListaPrecio"

	//Mongo---------------> Empresa

	//ColeccionEmpresa nombre de la coleccion de Empresa en mongo
	ColeccionEmpresa = "Empresa"

	//Mongo---------------> Caja

	//ColeccionCaja nombre de la coleccion de Caja en mongo
	ColeccionCaja = "Caja"

	//Mongo---------------> Dispositivo

	//ColeccionDispositivo nombre de la coleccion de Dispositivo en mongo
	ColeccionDispositivo = "Dispositivo"

	//Mongo---------------> Impuesto

	//ColeccionImpuesto nombre de la coleccion de Impuesto en mongo
	ColeccionImpuesto = "Impuesto"

	//Mongo---------------> Kit

	//ColeccionKit nombre de la coleccion de Kit en mongo
	ColeccionKit = "Kit"

	//Mongo---------------> MediosPago

	//ColeccionMediosPago nombre de la coleccion de MediosPago en mongo
	ColeccionMediosPago = "MediosPago"

	//Mongo---------------> Persona

	//ColeccionPersona nombre de la coleccion de Persona en mongo
	ColeccionPersona = "Persona"

	//Mongo---------------> Rol

	//ColeccionRol nombre de la coleccion de Rol en mongo
	ColeccionRol = "Rol"

	//Mongo---------------> Usuario

	//ColeccionUsuario nombre de la coleccion de Usuario en mongo
	ColeccionUsuario = "Usuario"

	//Mongo---------------> EquipoCaja

	//ColeccionEquipoCaja nombre de la coleccion de EquipoCaja en mongo
	ColeccionEquipoCaja = "EquipoCaja"

	//Mongo---------------> PuntoVenta

	//ColeccionPuntoVenta nombre de la coleccion de PuntoVenta en mongo
	ColeccionPuntoVenta = "PuntoVenta"

	//Mongo---------------> Facturacion

	//ColeccionFacturacion nombre de la coleccion de Facturacion en mongo
	ColeccionFacturacion = "Facturacion"

	//Mongo---------------> Conexion

	//ColeccionConexion nombre de la coleccion de Conexion en mongo
	ColeccionConexion = "Conexion"

	//Mongo---------------> PermisosUri

	//ColeccionPermisosUri nombre de la coleccion de PermisosUri en mongo
	ColeccionPermisosUri = "PermisosUri"

	//Mongo---------------> Operacion

	//ColeccionOperacion Nombre de la coleccion que almacena las operaciones Generales
	ColeccionOperacion = "Operacion"

	//Mongo---------------> Promocion

	//ColeccionPromocion nombre de la coleccion de Promocion en mongo
	ColeccionPromocion = "Promocion"

	//Mongo---------------> Cotizacion

	//ColeccionCotizacion nombre de la coleccion de Cotizacion en mongo
	ColeccionCotizacion = "Cotizacion"

	//Mongo---------------> GrupoAlmacen

	//ColeccionGrupoAlmacen nombre de la coleccion de GrupoAlmacen en mongo
	ColeccionGrupoAlmacen = "GrupoAlmacen"

	//Mongo---------------> Grupo

	//ColeccionGrupo nombre de la coleccion de Grupo en mongo
	ColeccionGrupo = "Grupo"

	//Mongo---------------> Bug

	//ColeccionBug nombre de la coleccion de Bug en mongo
	ColeccionBug = "Bug"

	//##########################<CATÁLOGOS DEL SISTEMA>######################

	//CatSysTipoImpuesto nombre de la colección de tipos de impuestos del sistema
	CatSysTipoImpuesto = "CatSysTiposDeImpuestos"

	//CatSysTipoFactor nombre de la colección de tipos de factores de impuestos del sistema
	CatSysTipoFactor = "CatSysTipoDeFactor"

	//CatSysClasificacionDeimpuestos nombre de la colección de tipos de factores de impuestos del sistema
	CatSysClasificacionDeimpuestos = "CatSysClasificacionDeimpuestos"

	//CatSysSubClasificacionDeimpuestos nombre de la colección de tipos de factores de impuestos del sistema
	CatSysSubClasificacionDeimpuestos = "CatSysSubClasificacionDeimpuestos"

	//################# DATOS ELASTIC ######################################

	//Elastic---------------> Catalogo

	//TipoCatalogo tipo a manejar en elastic
	TipoCatalogo = "Catalogo"

	//Elastic---------------> Unidad

	//TipoUnidad tipo a manejar en elastic
	TipoUnidad = "Unidad"

	//Elastic---------------> Cliente

	//TipoCliente tipo a manejar en elastic
	TipoCliente = "Cliente"

	//Elastic---------------> Producto

	//TipoProducto tipo a manejar en elastic
	TipoProducto = "Producto"

	//Elastic---------------> Almacen

	//TipoAlmacen tipo a manejar en elastic
	TipoAlmacen = "Almacen"

	//Elastic---------------> GrupoPersona

	//TipoGrupoPersona tipo a manejar en elastic
	TipoGrupoPersona = "GrupoPersona"

	//Elastic---------------> ListaCosto

	//TipoListaCosto tipo a manejar en elastic
	TipoListaCosto = "ListaCosto"

	//Elastic---------------> ListaPrecio

	//TipoListaPrecio tipo a manejar en elastic
	TipoListaPrecio = "ListaPrecio"

	//Elastic---------------> Empresa

	//TipoEmpresa tipo a manejar en elastic
	TipoEmpresa = "Empresa"

	//Elastic---------------> Caja

	//TipoCaja tipo a manejar en elastic
	TipoCaja = "Caja"

	//Elastic---------------> Dispositivo

	//TipoDispositivo tipo a manejar en elastic
	TipoDispositivo = "Dispositivo"

	//Elastic---------------> Impuesto

	//TipoImpuesto tipo a manejar en elastic
	TipoImpuesto = "Impuesto"

	//Elastic---------------> Kit

	//TipoKit tipo a manejar en elastic
	TipoKit = "Kit"

	//Elastic---------------> MediosPago

	//TipoMediosPago tipo a manejar en elastic
	TipoMediosPago = "MediosPago"

	//Elastic---------------> Persona

	//TipoPersona tipo a manejar en elastic
	TipoPersona = "Persona"

	//Elastic---------------> Rol

	//TipoRol tipo a manejar en elastic
	TipoRol = "Rol"

	//Elastic---------------> Usuario

	//TipoUsuario tipo a manejar en elastic
	TipoUsuario = "Usuario"

	//Elastic---------------> EquipoCaja

	//TipoEquipoCaja tipo a manejar en elastic
	TipoEquipoCaja = "EquipoCaja"

	//Elastic---------------> PuntoVenta

	//TipoPuntoVenta tipo a manejar en elastic
	TipoPuntoVenta = "PuntoVenta"

	//Elastic---------------> Operacion

	//TipoOperacion establecido temporalmente para minimizar los errores, posteriormente se obtendrá de algun catalogo
	TipoOperacion = "Venta"

	//Elastic---------------> Facturacion

	//TipoFacturacion tipo a manejar en elastic
	TipoFacturacion = "Facturacion"

	//Elastic---------------> Conexion

	//TipoConexion tipo a manejar en elastic
	TipoConexion = "Conexion"

	//Elastic---------------> PermisosUri

	//TipoPermisosUri tipo a manejar en elastic
	TipoPermisosUri = "PermisosUri"

	//Elastic---------------> Promocion

	//TipoPromocion tipo a manejar en elastic
	TipoPromocion = "Promocion"

	//Elastic---------------> Cotizacion

	//TipoCotizacion tipo a manejar en elastic
	TipoCotizacion = "Cotizacion"

	//Elastic---------------> GrupoAlmacen

	//TipoGrupoAlmacen tipo a manejar en elastic
	TipoGrupoAlmacen = "GrupoAlmacen"

	//Elastic---------------> Grupo

	//TipoGrupo tipo a manejar en elastic
	TipoGrupo = "Grupo"

//Elastic---------------> Bug

	//TipoBug tipo a manejar en elastic
	TipoBug = "Bug"


	//IndexElastic nombre del index a usar en elastic
	IndexElastic = "minisuperampliado"

	// * -----------------------CLAVES DE CATALOGO ------------------- * //

	// CatalogoDeEstatusDeCliente numero de catalogo en la coleccion Catalogos
	CatalogoDeEstatusDeCliente = 137
)

//DataCfg estructura de datos del entorno
type DataCfg struct {
	BaseURL    string
	Servidor   string
	Puerto     string
	Usuario    string
	Pass       string
	Protocolo  string
	NombreBase string
}

//#################<Funciones Generales>#######################################

//CargaSeccionCFG rellena los datos de la seccion a utilizar
func CargaSeccionCFG(seccion string) DataCfg {
	var d DataCfg
	var FileConfig, err = config.ReadDefault(FileConfigName)
	if err == nil {
		if FileConfig.HasOption(seccion, "baseurl") {
			d.BaseURL, _ = FileConfig.String(seccion, "baseurl")
		}
		if FileConfig.HasOption(seccion, "servidor") {
			d.Servidor, _ = FileConfig.String(seccion, "servidor")
		}
		if FileConfig.HasOption(seccion, "puerto") {
			d.Puerto, _ = FileConfig.String(seccion, "puerto")
		}
		if FileConfig.HasOption(seccion, "usuario") {
			d.Usuario, _ = FileConfig.String(seccion, "usuario")
		}
		if FileConfig.HasOption(seccion, "pass") {
			d.Pass, _ = FileConfig.String(seccion, "pass")
		}
		if FileConfig.HasOption(seccion, "protocolo") {
			d.Protocolo, _ = FileConfig.String(seccion, "protocolo")
		}
		if FileConfig.HasOption(seccion, "base") {
			d.NombreBase, _ = FileConfig.String(seccion, "base")
		}
	}
	return d
}
