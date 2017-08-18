package OperacionModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//OperacionElastic estructura de PuntoVentas para insertar en Elastic
type OperacionElastic struct {
	Operacion      string    `json:"Operacion,omitempty"`
	UsuarioOrigen  string    `json:"UsuarioOrigen,omitempty"`
	UsuarioDestino string    `json:"UsuarioDestino,omitempty"`
	Monto          string    `json:"Monto"`
	TipoOperacion  string    `json:"TipoOperacion,omitempty"`
	Predecesor     string    `json:"Predecesor,omitempty"`
	Estatus        string    `json:"Estatus,omitempty"`
	FechaHora      time.Time `json:"FechaHora"`
}

//OperacionMgo estructura de Operacions mongo
type OperacionMgo struct {
	ID                bson.ObjectId   `bson:"_id,omitempty"`
	UsuarioOrigen     bson.ObjectId   `bson:"UsuarioOrigen,omitempty"`
	UsuarioDestino    bson.ObjectId   `bson:"UsuarioDestino,omitempty"`
	FechaHoraRegistro time.Time       `bson:"FechaHoraRegistro"`
	TipoOperacion     bson.ObjectId   `bson:"TipoOperacion,omitempty"`
	Estatus           bson.ObjectId   `bson:"Estatus,omitempty"`
	Predecesor        bson.ObjectId   `bson:"Predecesor,omitempty"`
	Movimientos       []MovimientoMgo `bson:"Movimientos,omitempty"`
}

//MovimientoMgo subestructura de Operacion
type MovimientoMgo struct {
	IDMovimiento   bson.ObjectId    `bson:"IDMovimiento,omitempty"`
	AlmacenOrigen  bson.ObjectId    `bson:"AlmacenOrigen,omitempty"`
	AlmacenDestino bson.ObjectId    `bson:"AlmacenDestino,omitempty"`
	Ruta           []bson.ObjectId  `bson:"Ruta,omitempty"`
	Predecesor     bson.ObjectId    `bson:"Predecesor,omitempty"`
	Estatus        bson.ObjectId    `bson:"Estatus,omitempty"`
	Transacciones  []TransaccionMgo `bson:"Transacciones,omitempty"`
}

//TransaccionMgo subestructura de Operacion
type TransaccionMgo struct {
	IDTransaccion       bson.ObjectId `bson:"IDTransaccion,omitempty"`
	AlmacenOrigen       bson.ObjectId `bson:"AlmacenOrigen,omitempty"`
	AlmacenDestino      bson.ObjectId `bson:"AlmacenDestino,omitempty"`
	Estatus             bson.ObjectId `bson:"Estatus,omitempty"`
	Motivo              bson.ObjectId `bson:"Motivo,omitempty"`
	FechaHoraAplicacion time.Time     `bson:"FechaHoraAplicacion"`
}

//ImpuestoPostgres estructura que será usada para insertar y leer datos en postgres
type ImpuestoPostgres struct {
	IDMovimiento   string
	IDProducto     string
	TipoDeImpuesto string
	Valor          string
	Factor         string
	Tratamiento    string
	FechaHora      time.Time
}

//InventarioPostgres estructura que será usada para insertar y leer datos en postgres
type InventarioPostgres struct {
	IDProducto string
	Existencia float64
	Estatus    string
	Costo      float64
	Precio     float64
}

//KardexPostgres estructura que sera usada para almacenenar el registro de las entradas y salidas en un almacen de postgres
type KardexPostgres struct {
	IDOperacion    string
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

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []OperacionMgo {
	var result []OperacionMgo
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Operacions.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Operacions.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) OperacionMgo {
	var result OperacionMgo
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Operacions.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []OperacionMgo {
	var result []OperacionMgo
	var aux OperacionMgo
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = OperacionMgo{}
		Operacions.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) OperacionMgo {
	var result OperacionMgo
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)

	if err != nil {
		fmt.Println(err)
	}
	err = Operacions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecificsByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificsByFields(field string, valor interface{}) []OperacionMgo {
	var result []OperacionMgo
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)

	if err != nil {
		fmt.Println(err)
	}
	err = Operacions.Find(bson.M{field: valor}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecificsByFields2 regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificsByFields2(fields []string, valores []interface{}) []OperacionMgo {
	var result []OperacionMgo
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}

	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range fields {
		Abson[v] = valores[k]
	}
	find := Abson
	err = Operacions.Find(find).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result OperacionMgo
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Operacions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//GeneraTemplatesBusquedaParaPuntoDeVenta crea templates de tabla de búsqueda
func GeneraTemplatesBusquedaParaPuntoDeVenta(Operaciones []OperacionMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>			
					<th>Operacion</th>									
					<th>Usuario Origen</th>									
					<th>Usuario Destino</th>									
					<th>Monto</th>																								
					<th>Estatus</th>									
					<th>Fecha</th>					
				</tr>`

	for k, v := range Operaciones {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/PuntoVentas/edita/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.ID.Hex() + `</td>`
		cuerpo += `<td>` + v.UsuarioOrigen.Hex() + `</td>`
		cuerpo += `<td>` + v.UsuarioDestino.Hex() + `</td>`
		cuerpo += `<td>` + "Monto" + `</td>`
		//cuerpo += `<td>` + floats.FormatMoney(v.Codigo) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 171) + `</td>`
		cuerpo += `<td>` + v.FechaHoraRegistro.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

/*
//ObtenerParametrosConexion Obtiene los datos de una conexion pasandole el identificador del almacén
func ObtenerParametrosConexion(idAlmacen bson.ObjectId) ConexionModel.ConexionMgo {
	datosAlmacen := AlmacenModel.GetOne(idAlmacen)
	var datosConexion ConexionModel.ConexionMgo
	if datosAlmacen.Conexion != "" {
		datosConexion = ConexionModel.GetOne(datosAlmacen.Conexion)
	} else {
		fmt.Println("El inventario no tiene datos de conexion, favor de verificar.", idAlmacen, datosAlmacen.Nombre)
	}
	return datosConexion
}
*/

//########## GET NAMES ####################################

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoOperacion, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoOperacion, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
