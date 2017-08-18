package CajaModel

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"github.com/leekchan/accounting"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//CajaMgo estructura de Cajas mongo
type CajaMgo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Usuario   bson.ObjectId `bson:"Usuario,omitempty"`
	Caja      bson.ObjectId `bson:"Caja,omitempty"`
	Cargo     float64       `bson:"Cargo"`
	Abono     float64       `bson:"Abono"`
	Saldo     float64       `bson:"Saldo"`
	Operacion bson.ObjectId `bson:"Operacion,omitempty"`
	Estatus   bson.ObjectId `bson:"Estatus"`
	FechaHora time.Time     `bson:"FechaHora"`
}

//CajaElastic estructura de Cajas para insertar en Elastic
type CajaElastic struct {
	Usuario   string    `json:"Usuario,omitempty"`
	Caja      string    `json:"Caja,omitempty"`
	Cargo     float64   `json:"Cargo"`
	Abono     float64   `json:"Abono"`
	Saldo     float64   `json:"Saldo"`
	Operacion string    `json:"Operacion,omitempty"`
	Estatus   string    `json:"Estatus"`
	FechaHora time.Time `json:"FechaHora"`
}

//FormasDePago estructura provisional de formas de pago
type FormasDePago struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      string        `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	CodigoSat   string        `bson:"CosigoSat"`
	Comision    float64       `bson:"Comision"`
	Tipo        string        `bson:"Tipo"`
	Cabio       bool          `bson:"Cambio"`
}

//Operaciones estrucutura provisional de operaciones
type Operaciones struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Tipo     string        `bson:"Tipo"`
	Concepto string        `bson:"Concepto"`
	Monto    float64       `bson:"Monto"`
}

//Staff estructura provisional para generacion de XML
type Staff struct {
	XMLName   xml.Name `xml:"staff"`
	ID        int      `xml:"id"`
	FirstName string   `xml:"firstname"`
	LastName  string   `xml:"lastname"`
	UserName  string   `xml:"username"`
}

//Company estructura provisional para generacion de XML
type Company struct {
	XMLName xml.Name `xml:"company"`
	Staffs  []Staff  `xml:"staff"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []CajaMgo {
	var result []CajaMgo
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Cajas.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetAllFormas Provisional, devuelve las formas de Pago.
func GetAllFormas() []FormasDePago {
	var result []FormasDePago
	s, Cajas, err := MoConexion.GetColectionMgo("FormasDePago")
	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) CajaMgo {
	var result CajaMgo
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOperacionesByID Regresa un documento específico de Mongo (Por Coleccion)
func GetOperacionesByID(ID bson.ObjectId) []Operaciones {
	var result []Operaciones
	s, Cajas, err := MoConexion.GetColectionMgo("Operaciones")
	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(bson.M{"_id": ID}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOperacionesAll Regresa todos los documentos de una coleccion
func GetOperacionesAll() []Operaciones {
	var result []Operaciones
	s, Cajas, err := MoConexion.GetColectionMgo("Operaciones")
	if err != nil {
		fmt.Println(err)
	}
	//Trae todas las operaciones, donde estatus sea (0),pendientes de pago.
	err = Cajas.Find(bson.M{"Estatus": 0}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetCajaAbiertaByUsuario regresa un documento específico de Mongo (Por Coleccion)
func GetCajaAbiertaByUsuario(ID bson.ObjectId) CajaMgo {
	var result CajaMgo
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(bson.M{"Usuario": ID}).One(&result)
	/////
	/*
		o1 := bson.M{"$match": bson.M{"Usuario": ID}}
		o2 := bson.M{"$group": bson.M{"Usuario": ID, "Cargo": bson.M{"$sum": "$qty"}, "count": bson.M{"$sum": 1}}}
		operations := []bson.M{o1, o2}
		pipe := Cajas.Pipe(operations)
		err = pipe.All(&result)
	*/
	//////
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	fmt.Println(result)
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []CajaMgo {
	var result []CajaMgo
	var aux CajaMgo
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = CajaMgo{}
		Cajas.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) CajaMgo {
	var result CajaMgo
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)

	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result CajaMgo
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Cajas []CajaMgo) (string, string) {
	floats := accounting.Accounting{Symbol: "$", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
				<th>#</th>			
				<th>Usuario</th>									
				<th>Caja</th>									
				<th>Cargo</th>									
				<th>Abono</th>									
				<th>Saldo</th>									
				<th>Operacion</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

	for k, v := range Cajas {
		//cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Cajas/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<tr id = "` + v.ID.Hex() + `" >`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Usuario.Hex() + `</td>`
		cuerpo += `<td>` + v.Caja.Hex() + `</td>`
		cuerpo += `<td>` + floats.FormatMoney(v.Cargo) + `</td>`
		cuerpo += `<td>` + floats.FormatMoney(v.Abono) + `</td>`
		cuerpo += `<td>` + floats.FormatMoney(v.Saldo) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Operacion, 169) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 135) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

/*
Estos tres metodos se pasan al modelo de EquipoCaja model
*********************************************************
El motivo es por que se cambio la estructura
*********************************************************
//CargaComboCajas regresa un combo de Caja de mongo
func CargaComboCajas(ID string) string {
	Cajas := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Cajas {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}


//CargaComboCajasMulti regresa un combo de Caja de mongo
func CargaComboCajasMulti(ID string) string {
	Cajas := GetAll()

	templ := ``

	for _, v := range Cajas {
		templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
	}

	return templ
}

func CargaComboCajasMultiArrayObjID(ArrayObID []bson.ObjectId) string {
	Cajas := GetAll()
	var templ string
	for _, v := range Cajas {
		for _, vv := range ArrayObID {
			if vv == v.ID {
				templ += `<option value="` + vv.Hex() + `" selected>  ` + v.Nombre + ` </option> `
			}
		}
	}
	return templ
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

	queryTilde = queryTilde.Field("Usuario")
	queryQuotes = queryQuotes.Field("Usuario")

	queryTilde = queryTilde.Field("Operacion")
	queryQuotes = queryQuotes.Field("Operacion")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoCaja, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoCaja, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}

//GetSaldoCaja trae el saldo de caja abierta para este usuario para mostrarlo en la vista de operaciones
func GetSaldoCaja(idCaja string, idOperacion string) float64 {
	var result []CajaMgo
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = Cajas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	//
	fmt.Println(idCaja, idOperacion)
	var lineas int
	var sumacargo float64
	var sumaabono float64
	var saldo float64
	for _, row := range result {
		//fmt.Println(i)
		lineas = lineas + 1
		sumacargo = sumacargo + row.Cargo
		sumaabono = sumaabono + row.Abono

	}
	saldo = sumacargo - sumaabono
	//
	s.Close()
	return saldo
}
