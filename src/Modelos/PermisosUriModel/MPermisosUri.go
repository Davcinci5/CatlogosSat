package PermisosUriModel

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/GrupoPersonaModel"
	"../../Modelos/UsuarioModel"

	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//CatalogoUris Coleccion de los catalogos de URIS
var CatalogoUris = 177

//#########################< ESTRUCTURAS >##############################

//PermisosUriMgo estructura de PermisosUris mongo
type PermisosUriMgo struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	Grupo           string        `bson:"Grupo,omitempty"`
	PermisoNegado   string        `bson:"PermisoNegado"`
	PermisoAceptado string        `bson:"PermisoAceptado"`
}

//PermisosUriElastic estructura de PermisosUris para insertar en Elastic
type PermisosUriElastic struct {
	Grupo           string `json:"Grupo,omitempty"`
	PermisoNegado   string `json:"PermisoNegado"`
	PermisoAceptado string `json:"PermisoAceptado"`
}

type PermisoAux struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Nombre    string        `bson:"Nombre"`
	EsUsuario bool          `bson:"Usuario"`
	Estatus   bson.ObjectId `bson:"Estatus"`
	FechaHora time.Time     `bson:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []PermisosUriMgo {
	var result []PermisosUriMgo
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}
	err = PermisosUris.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)

	if err != nil {
		fmt.Println(err)
	}
	result, err = PermisosUris.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) PermisosUriMgo {
	var result PermisosUriMgo
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}
	err = PermisosUris.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []PermisosUriMgo {
	var result []PermisosUriMgo
	var aux PermisosUriMgo
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = PermisosUriMgo{}
		PermisosUris.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) PermisosUriMgo {
	var result PermisosUriMgo
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)

	if err != nil {
		fmt.Println(err)
	}
	err = PermisosUris.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result PermisosUriMgo
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}
	err = PermisosUris.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboPermisosUris regresa un combo de PermisosUri de mongo
func CargaComboPermisosUris(ID string) string {
	PermisosUris := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range PermisosUris {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Grupo + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Grupo + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(PermisosUris []PermisosUriMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
				<th>#</th>
				<th>Grupo</th>					
				</tr>`

	for k, v := range PermisosUris {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/PermisosUris/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Grupo + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda2(PermisosUris []PermisoAux) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	estatus := ``
	cuerpo := ``

	cabecera := `<tr>
				<th>#</th>
				<th>Nombre</th>
				<th>Tipo</th>				
				<th>Estatus</th>					
				<th>Fecha</th>					
				</tr>`

	for k, v := range PermisosUris {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/PermisosUris/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		if v.EsUsuario {
			cuerpo += `<td>USUARIO</td>`
		} else {
			cuerpo += `<td>GRUPO</td>`
		}
		estatus = CatalogoModel.GetValorMagnitud(v.Estatus, 167)
		if estatus == "" {
			estatus = CatalogoModel.GetValorMagnitud(v.Estatus, 146)
		}
		cuerpo += `<td>` + estatus + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//GeneraEtiquetasURIS genera la etiquetas para los permisos de cada URI
func GeneraEtiquetasURIS() string {
	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(CatalogoUris))
	templ := ``
	for _, values := range Catalogo.Valores {
		templ += `<a id="` + values.ID.Hex() + `" class="list-group-item" draggable="true" ondragstart="drag(event)"> <input value="` + values.ID.Hex() + `" name="PermisosNegados" readonly hidden>` + strings.ToLower(values.Valor) + `</a>`

	}
	return templ
}

//GeneraEtiquetasEspecificPermisosURIS genera la etiquetas para de cada URI con permisos
func GeneraEtiquetasEspecificPermisosURIS(ides []string) string {
	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(CatalogoUris))
	templ := ``
	for _, values := range Catalogo.Valores {
		for _, valor := range ides {
			if valor == values.ID.Hex() {
				templ += `<a id="` + values.ID.Hex() + `" class="list-group-item" draggable="true" ondragstart="drag(event)"> <input value="` + values.ID.Hex() + `" name="PermisosAceptados" readonly hidden>` + strings.ToLower(values.Valor) + `</a>`
			}
		}
	}
	return templ
}

//GeneraEtiquetasEspecificNoPermisosURIS genera las etiquetas de cada URI que no estan asignadas
func GeneraEtiquetasEspecificNoPermisosURIS(ides []string) string {
	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(CatalogoUris))
	URISAll := Catalogo.Valores
	templ := ``
	for _, values := range URISAll { //Recorre todos los ids de las URIS
		existe := false
		for _, existente := range ides { //Recorre los IDS de Uris permitidas
			if values.ID.Hex() == existente {
				existe = true //si existe lo marca
			}
		}
		if !existe {
			templ += `<a id="` + values.ID.Hex() + `" class="list-group-item" draggable="true" ondragstart="drag(event)"> <input value="` + values.ID.Hex() + `" name="PermisosNegados" readonly hidden>` + strings.ToLower(values.Valor) + `</a>`
		}
	}
	return templ
}

//JoinGroupUser join user and groups for show
func JoinGroupUser(Usuarios *[]UsuarioModel.UsuarioMgo, Grupos *[]GrupoPersonaModel.GrupoPersonaMgo) *[]PermisoAux {
	var Auxiliares []PermisoAux
	var aux PermisoAux

	for _, v := range *Usuarios {
		aux = PermisoAux{}
		aux.EsUsuario = true
		aux.Nombre = v.Usuario
		aux.ID = v.ID
		aux.Estatus = v.Estatus
		aux.FechaHora = v.FechaHora
		Auxiliares = append(Auxiliares, aux)
	}

	for _, v := range *Grupos {
		aux = PermisoAux{}
		aux.EsUsuario = false
		aux.Nombre = v.Nombre
		aux.ID = v.ID
		aux.Estatus = v.Estatus
		aux.FechaHora = v.FechaHora
		Auxiliares = append(Auxiliares, aux)
	}

	return &Auxiliares
}

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("Usuario")
	queryQuotes = queryQuotes.Field("Usuario")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoUsuario, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoUsuario, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}

//BuscarEnElasticG busca el texto solicitado en los campos solicitados
func BuscarEnElasticG(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("Nombre")
	queryQuotes = queryQuotes.Field("Nombre")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupoPersona, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupoPersona, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
