package MediosPagoModel

import (
	"fmt"
	"strconv"
	"time"

	"github.com/leekchan/accounting"

	"../../Modelos/CatalogoModel"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//MediosPagoMgo estructura de MediosPagos mongo
type MediosPagoMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      bson.ObjectId `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	CodigoSat   string        `bson:"CodigoSat"`
	Tipo        bson.ObjectId `bson:"Tipo"`
	Comision    float64       `bson:"Comision"`
	Cambio      bool          `bson:"Cambio"`
	Estatus     bson.ObjectId `bson:"Estatus"`
	FechaHora   time.Time     `bson:"FechaHora"`
}

//MediosPagoElastic estructura de MediosPagos para insertar en Elastic
type MediosPagoElastic struct {
	Nombre      string    `json:"Nombre"`
	Descripcion string    `json:"Descripcion"`
	CodigoSat   string    `json:"CodigoSat"`
	Tipo        string    `json:"Tipo"`
	Comision    float64   `json:"Comision"`
	Cambio      bool      `json:"Cambio"`
	Estatus     string    `json:"Estatus"`
	FechaHora   time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []MediosPagoMgo {
	var result []MediosPagoMgo
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}
	err = MediosPagos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) MediosPagoMgo {
	var result MediosPagoMgo
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}
	err = MediosPagos.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)

	if err != nil {
		fmt.Println(err)
	}
	result, err = MediosPagos.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []MediosPagoMgo {
	var result []MediosPagoMgo
	var aux MediosPagoMgo
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = MediosPagoMgo{}
		MediosPagos.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) MediosPagoMgo {
	var result MediosPagoMgo
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)

	if err != nil {
		fmt.Println(err)
	}
	err = MediosPagos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result MediosPagoMgo
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}
	err = MediosPagos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboMediosPagos regresa un combo de MediosPago de mongo
func CargaComboMediosPagos(ID string) string {
	MediosPagos := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range MediosPagos {
		if ID == v.ID.Hex() {
			templ += `<option value=" ` + v.ID.Hex() + `" selected>  ` + v.Nombre.Hex() + ` </option> `
		} else {
			templ += `<option value=" ` + v.ID.Hex() + `">  ` + v.Nombre.Hex() + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(MediosPagos []MediosPagoMgo) (string, string) {
	floats := accounting.Accounting{Symbol: "$", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
				<th>#</th>			
				<th>Nombre</th>									
				<th>Descripcion</th>									
				<th>CodigoSat</th>					
				<th>Tipo</th>									
				<th>Comisión</th>									
				<th>¿Da Cambio?</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

	for k, v := range MediosPagos {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/MediosPagos/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Nombre, 156) + `</td>`
		cuerpo += `<td>` + v.Descripcion + `</td>`
		cuerpo += `<td>` + v.CodigoSat + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo, 157) + `</td>`
		cuerpo += `<td>` + floats.FormatMoney(v.Comision) + `</td>`
		if v.Cambio {
			cuerpo += `<td>SI</td>`
		} else {
			cuerpo += `<td>NO</td>`
		}
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 158) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("Nombre")
	queryQuotes = queryQuotes.Field("Nombre")

	queryTilde = queryTilde.Field("Descripcion")
	queryQuotes = queryQuotes.Field("Descripcion")

	queryTilde = queryTilde.Field("CodigoSat")
	queryQuotes = queryQuotes.Field("CodigoSat")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoMediosPago, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoMediosPago, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
