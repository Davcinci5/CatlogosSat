package BugModel

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

//BugMgo estructura de Bugs mongo
type BugMgo struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	Tipo            string        `bson:"Tipo"`
	Titulo          string        `bson:"Titulo"`
	Descripcion     string        `bson:"Descripcion"`
	Usuario         string        `bson:"Usuario,omitempty"`
	Metodo          string        `bson:"Metodo,omitempty"`
	EsAjax          bool          `bson:"EsAjax,omitempty"`
	EstatusPeticion string        `bson:"EstatusPeticion"`
	Estatus         bson.ObjectId `bson:"Estatus,omitempty"`
	Ruta            string        `bson:"Ruta"`
	FechaHora       time.Time     `bson:"FechaHora"`
}

//BugElastic estructura de Bugs para insertar en Elastic
type BugElastic struct {
	Tipo            string    `json:"Tipo"`
	Titulo          string    `json:"Titulo"`
	Descripcion     string    `json:"Descripcion"`
	Usuario         string    `json:"Usuario,omitempty"`
	Metodo          string    `json:"Metodo,omitempty"`
	EsAjax          bool      `json:"EsAjax,omitempty"`
	EstatusPeticion string    `json:"EstatusPeticion"`
	Estatus         string    `json:"Estatus,omitempty"`
	Ruta            string    `json:"Ruta"`
	FechaHora       time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []BugMgo {
	var result []BugMgo
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}
	err = Bugs.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Bugs.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) BugMgo {
	var result BugMgo
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}
	err = Bugs.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []BugMgo {
	var result []BugMgo
	var aux BugMgo
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = BugMgo{}
		Bugs.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) BugMgo {
	var result BugMgo
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)

	if err != nil {
		fmt.Println(err)
	}
	err = Bugs.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result BugMgo
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}
	err = Bugs.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboBugs regresa un combo de Bug de mongo
func CargaComboBugs(ID string) string {
	Bugs := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Bugs {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Titulo + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Titulo + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Bugs []BugMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
			<th>#</th>
			
				<th>Tipo</th>					
				
				<th>Titulo</th>					
				
				<th>Descripcion</th>					
				
				<th>Usuario</th>					
				
				<th>Metodo</th>					
				
				<th>EsAjax</th>					
				
				<th>EstatusPeticion</th>					
				
				<th>Estatus</th>					
				
				<th>Ruta</th>					
				
				<th>FechaHora</th>					
				</tr>`

	for k, v := range Bugs {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Bugs/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Tipo + `</td>`

		cuerpo += `<td>` + v.Titulo + `</td>`

		cuerpo += `<td>` + v.Descripcion + `</td>`

		cuerpo += `<td>` + v.Usuario + `</td>`

		cuerpo += `<td>` + v.Metodo + `</td>`

		cuerpo += `<td>` + strconv.FormatBool(v.EsAjax) + `</td>`

		cuerpo += `<td>` + v.EstatusPeticion + `</td>`

		cuerpo += `<td>` + v.Estatus.Hex() + `</td>`

		cuerpo += `<td>` + v.Ruta + `</td>`

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

	queryTilde = queryTilde.Field("Tipo")
	queryQuotes = queryQuotes.Field("Tipo")

	queryTilde = queryTilde.Field("Titulo")
	queryQuotes = queryQuotes.Field("Titulo")

	queryTilde = queryTilde.Field("Descripcion")
	queryQuotes = queryQuotes.Field("Descripcion")

	queryTilde = queryTilde.Field("Usuario")
	queryQuotes = queryQuotes.Field("Usuario")

	queryTilde = queryTilde.Field("EstatusPeticion")
	queryQuotes = queryQuotes.Field("EstatusPeticion")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoBug, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoBug, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
