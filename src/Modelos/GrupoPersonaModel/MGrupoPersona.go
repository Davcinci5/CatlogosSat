package GrupoPersonaModel

import (
	"fmt"
	"time"

	"../../Modelos/UsuarioModel"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"strings"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//GrupoPersonaMgo estructura de GrupoPersonas mongo
type GrupoPersonaMgo struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Nombre      string          `bson:"Nombre"`
	Descripcion string          `bson:"Descripcion"`
	Miembros    []bson.ObjectId `bson:"Miembros"`
	Estatus     bson.ObjectId   `bson:"Estatus"`
	FechaHora   time.Time       `bson:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []GrupoPersonaMgo {
	var result []GrupoPersonaMgo
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	err = GrupoPersonas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)

	if err != nil {
		fmt.Println(err)
	}
	result, err = GrupoPersonas.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) GrupoPersonaMgo {
	var result GrupoPersonaMgo
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	err = GrupoPersonas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []GrupoPersonaMgo {
	var result []GrupoPersonaMgo
	var aux GrupoPersonaMgo
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = GrupoPersonaMgo{}
		GrupoPersonas.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) GrupoPersonaMgo {
	var result GrupoPersonaMgo
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)

	if err != nil {
		fmt.Println(err)
	}
	err = GrupoPersonas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetArrayEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetArrayEspecificByFields(field string, valor interface{}) []GrupoPersonaMgo {
	var result []GrupoPersonaMgo
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)

	if err != nil {
		fmt.Println(err)
	}
	err = GrupoPersonas.Find(bson.M{field: valor}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result GrupoPersonaMgo
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	err = GrupoPersonas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboGrupoPersonas regresa un combo de GrupoPersona de mongo
func CargaComboGrupoPersonas(ID string) string {
	GrupoPersonas := GetAll()
	var templ string
	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range GrupoPersonas {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

func CargaComboGrupoPersonasMulti(ID string) string {
	GrupoPersonas := GetAll()

	var templ string

	for _, v := range GrupoPersonas {

		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}
	}
	return templ
}

//CargaComboGrupoPersonasArray Recibe un arreglo IDs y regresa los option con los ids del arreglo seleccionados
func CargaComboGrupoPersonasArray(ArrayObID []string) string {
	GrupoPersonas := GetAll()
	var templ string
	for _, v := range GrupoPersonas {
		existe := false
		for _, vv := range ArrayObID {
			if vv == v.ID.Hex() {
				existe = true
			}
		}
		if existe {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}
	}
	return templ
}

//CargaNombreGrupo regresa el nombre de un grupo
func CargaNombreGrupo(ID bson.ObjectId) string {
	var result GrupoPersonaMgo
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	err = GrupoPersonas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.Nombre

}

// GeneraEtiquetasPersonas regresa una lista de los Usuarios en el sistema
func GeneraEtiquetasPersonas() string {
	usuarios := UsuarioModel.GetAll()
	html := ``
	for _, values := range usuarios {
		html += `<a id="` + values.ID.Hex() + `" class="list-group-item" draggable="true" ondragstart="drag(event)"> <input value="` + values.ID.Hex() + `" name="PersonaGpo" readonly hidden>` + strings.ToLower(values.Usuario) + `</a>`
	}
	return html
}

// GeneraEtiquetasPersonasEnUnGpo regresa una lista de los Usuarios por IDS en el sistema
func GeneraEtiquetasPersonasEnUnGpo(MiembrosIDS []bson.ObjectId) string {
	usuarios := UsuarioModel.GetEspecifics(MiembrosIDS)
	html := ``
	for _, values := range usuarios {
		if values.ID.Hex() != "" {
			html += `<a id="` + values.ID.Hex() + `" class="list-group-item" draggable="true" ondragstart="drag(event)"> <input value="` + values.ID.Hex() + `" name="Miembros" readonly hidden>` + strings.ToLower(values.Usuario) + `</a>`
		}
	}
	return html
}

// GeneraEtiquetasPersonasFueraDeUnGpo regresa una lista de los Usuarios por IDS en el sistema
func GeneraEtiquetasPersonasFueraDeUnGpo(MiembrosIDS []bson.ObjectId) string {
	TodosUsuarios := UsuarioModel.GetAll()
	var ides []bson.ObjectId

	for _, values := range TodosUsuarios {
		ides = append(ides, values.ID)
	}

	var SinGrupo []bson.ObjectId
	for _, values := range ides { //Recorre todos los ids de grupos y usuarios
		existe := false
		for _, existente := range MiembrosIDS { //Por cada grupo o usuario, Recorre los existentes
			if values.Hex() == existente.Hex() {
				existe = true //si existe lo marca
			}
		}
		if !existe {
			SinGrupo = append(SinGrupo, values) //los que no existen se agregan al arreglo
		}
	}

	usuarios := UsuarioModel.GetEspecifics(SinGrupo)
	html := ``
	for _, values := range usuarios {
		if values.ID.Hex() != "" {
			html += `<a id="` + values.ID.Hex() + `" class="list-group-item" draggable="true" ondragstart="drag(event)"> <input value="` + values.ID.Hex() + `" name="PersonaGpo" readonly hidden>` + strings.ToLower(values.Usuario) + `</a>`
		}
	}
	return html
}

// ConvierteAObjectIDS recibe un arreglo de string de IDS y retorna un arreglo de bson.ObjectId
func ConvierteAObjectIDS(IDS []string) []bson.ObjectId {
	var NuevosIDS []bson.ObjectId
	for _, val := range IDS {
		if bson.IsObjectIdHex(val) {
			NuevosIDS = append(NuevosIDS, bson.ObjectIdHex(val))
		}
	}
	return NuevosIDS
}

//RegresaGruposDeUsuario Regresa un arreglo con los IDs de los grupos a los que pertenece el usuario
func RegresaGruposDeUsuario(ID bson.ObjectId) []string {
	gpo := GetArrayEspecificByFields("Miembros", ID)
	var aux []string
	for _, val := range gpo {
		aux = append(aux, val.ID.Hex())
	}
	return aux
}

// GetIDUsuario Regresa el ID del Usuario especificado
func GetIDUsuario(NombreUsuario string) string {
	return UsuarioModel.GetIDByField("Usuario", NombreUsuario).Hex()
}

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupoPersona, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupoPersona, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
