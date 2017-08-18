package PersonaModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/GrupoPersonaModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//PersonaMgo estructura de Personas mongo
type PersonaMgo struct {
	ID              bson.ObjectId   `bson:"_id,omitempty"`
	Nombre          string          `bson:"Nombre"`
	Sexo            string          `bson:"Sexo,omitempty"`
	FechaNacimiento time.Time       `bson:"FechaNacimiento,omitempty"`
	Tipo            []bson.ObjectId `bson:"Tipo,omitempty"`
	Grupos          []bson.ObjectId `bson:"Grupos,omitempty"`
	Predecesor      bson.ObjectId   `bson:"Predecesor,omitempty"`
	Notificacion    []bson.ObjectId `bson:"Notificacion,omitempty"`
	Estatus         bson.ObjectId   `bson:"Estatus,omitempty"`
	FechaHora       time.Time       `bson:"FechaHora,omitempty"`
}

//PersonaElastic estructura de Personas para insertar en Elastic
type PersonaElastic struct {
	Nombre     string    `json:"Nombre"`
	Tipo       []string  `json:"Tipo"`
	Grupos     []string  `json:"Grupos,omitempty"`
	Predecesor string    `json:"Predecesor,omitempty"`
	Estatus    string    `json:"Estatus"`
	FechaHora  time.Time `json:"FechaHora,omitempty"`
}

//NotificacionMgo subestructura de Persona
type NotificacionMgo struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Mensaje          string        `bson:"Mensaje"`
	Leido            bool          `bson:"Leido"`
	FechaOcurrencia  time.Time     `bson:"FechaOcurrencia"`
	FechaVencimiento time.Time     `bson:"FechaVencimiento"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []PersonaMgo {
	var result []PersonaMgo
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}
	err = Personas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Personas.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) PersonaMgo {
	var result PersonaMgo
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}
	err = Personas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []PersonaMgo {
	var result []PersonaMgo
	var aux PersonaMgo
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = PersonaMgo{}
		Personas.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) PersonaMgo {
	var result PersonaMgo
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)

	if err != nil {
		fmt.Println(err)
	}
	err = Personas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result PersonaMgo
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}
	err = Personas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboPersonas regresa un combo de Persona de mongo
func CargaComboPersonas(ID string) string {
	Personas := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Personas {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//ConstruirComboUsuarioPredecesores Combo de Predecesor para Personas
func ConstruirComboUsuarioPredecesores(ID string) string {
	var templ string
	var result []PersonaMgo

	templ += `<option value="">--SELECCIONE--</option>`
	result = GetAll()

	for _, v := range result {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Nombre + `</option>`
		} else {
			templ += `<option value="` + v.ID.Hex() + `">` + v.Nombre + `</option>`
		}

	}

	return templ

}

//CargaComboCatalogoArrayID Regresa un tmplat de option para un Select para Personsa.
func CargaComboCatalogoArrayID(Clave int, persona PersonaMgo) string {

	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))

	templ := ``

	for _, v := range persona.Tipo {

		for _, vv := range Catalogo.Valores {
			if v == vv.ID {
				templ += `<option value="` + vv.ID.Hex() + `" selected>` + vv.Valor + `</option>`
			}
		}
	}
	return templ
}

//CargaNombrePredecesor Funcion que Regresa el <option> del nombre de un Predecesor partiendo de un Object Id
func CargaNombrePredecesor(IDPredecesor bson.ObjectId) string {
	Predecesor := GetOne(IDPredecesor)
	var templ string
	templ += `<option value="` + IDPredecesor.Hex() + `">` + Predecesor.Nombre + `</option>`
	return templ
}

// CargaNombresTiposPersonas regresa nombre de los tipos de personas de acuerdo a los IDS
func CargaNombresTiposPersonas(IDS []bson.ObjectId) []string {
	var tipos []string
	for _, val := range IDS {
		tipos = append(tipos, CatalogoModel.RegresaNombreSubCatalogo(val))
	}
	return tipos
}

//CargaNombresGruposPersonas regresa un arreglo de nombres de grupos de acuerdo a los IDS
func CargaNombresGruposPersonas(IDS []bson.ObjectId) []string {
	var grupos []string
	for _, val := range IDS {
		grupos = append(grupos, GrupoPersonaModel.CargaNombreGrupo(val))
	}
	return grupos
}

// CargaNombreEstatusGrupo regresa el nombre un estatus de persona
func CargaNombreEstatusGrupo(ID bson.ObjectId) string {
	return CatalogoModel.GetValorMagnitud(ID, 160)
}

//NombrePredecesor Funcion que Regresa el nombre de un Predecesor partiendo de un Object Id
func NombrePredecesor(IDPredecesor bson.ObjectId) string {
	Predecesor := GetOne(IDPredecesor)
	return Predecesor.Nombre
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Personas []PersonaMgo) (string, string) {
	cuerpo := ``

	cabecera := `<tr>
				<th>#</th>			
				<th>Nombre</th>									
				<th>Tipo</th>									
				<th>Grupos</th>									
				<th>Predecesor</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

	for k, v := range Personas {
		//cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Personas/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<tr id = "` + v.ID.Hex() + `">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo[0], 159) + `</td>`
		cuerpo += `<td>` + "v.Grupos" + `</td>`
		cuerpo += `<td>` + v.Predecesor.Hex() + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 160) + `</td>`
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

	queryTilde = queryTilde.Field("Tipo")
	queryQuotes = queryQuotes.Field("Tipo")

	// queryTilde = queryTilde.Field("FechaHora")
	// queryQuotes = queryQuotes.Field("FechaHora")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoPersona, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoPersona, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
