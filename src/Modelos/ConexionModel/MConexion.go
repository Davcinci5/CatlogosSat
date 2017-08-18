package ConexionModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//ConexionMgo estructura de Conexions mongo
type ConexionMgo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Nombre    string        `bson:"Nombre"`
	Servidor  string        `bson:"Servidor"`
	NombreBD  string        `bson:"NombreBD"`
	UsuarioBD string        `bson:"UsuarioBD,omitempty"`
	PassBD    string        `bson:"PassBD,omitempty"`
	FechaHora time.Time     `bson:"FechaHora"`
}

//ConexionElastic estructura de Conexions para insertar en Elastic
type ConexionElastic struct {
	Nombre    string    `json:"Nombre"`
	Servidor  string    `json:"Servidor"`
	NombreBD  string    `json:"NombreBD"`
	UsuarioBD string    `json:"UsuarioBD,omitempty"`
	PassBD    string    `json:"PassBD,omitempty"`
	FechaHora time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []ConexionMgo {
	var result []ConexionMgo
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}
	err = Conexions.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Conexions.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) ConexionMgo {
	var result ConexionMgo
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}
	err = Conexions.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []ConexionMgo {
	var result []ConexionMgo
	var aux ConexionMgo
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = ConexionMgo{}
		Conexions.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) ConexionMgo {
	var result ConexionMgo
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)

	if err != nil {
		fmt.Println(err)
	}
	err = Conexions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result ConexionMgo
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}
	err = Conexions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboConexions regresa un combo de Conexion de mongo
func CargaComboConexions(ID string) string {
	Conexions := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Conexions {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Conexions []ConexionMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
				<th>#</th>			
				<th>Nombre</th>									
				<th>Servidor</th>									
				<th>NombreBD</th>									
				<th>UsuarioBD</th>									
				<th>PassBD</th>									
				<th>FechaHora</th>					
				</tr>`

	for k, v := range Conexions {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Conexions/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + v.Servidor + `</td>`
		cuerpo += `<td>` + v.NombreBD + `</td>`
		cuerpo += `<td>` + v.UsuarioBD + `</td>`
		cuerpo += `<td>` + v.PassBD + `</td>`
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

	queryTilde = queryTilde.Field("Servidor")
	queryQuotes = queryQuotes.Field("Servidor")

	queryTilde = queryTilde.Field("NombreBD")
	queryQuotes = queryQuotes.Field("NombreBD")

	queryTilde = queryTilde.Field("UsuarioBD")
	queryQuotes = queryQuotes.Field("UsuarioBD")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoConexion, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoConexion, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
