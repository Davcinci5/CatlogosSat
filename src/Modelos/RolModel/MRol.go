package RolModel

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

//RolMgo estructura de Rols mongo
type RolMgo struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Nombre      string          `bson:"Nombre"`
	Descripcion string          `bson:"Descripcion"`
	Permisos    []bson.ObjectId `bson:"Permisos"`
	Estatus     bson.ObjectId   `bson:"Estatus"`
	FechaHora   time.Time       `bson:"FechaHora"`
}

//RolElastic estructura de Rols para insertar en Elastic
type RolElastic struct {
	Nombre      string    `json:"Nombre"`
	Descripcion string    `json:"Descripcion"`
	Permisos    []string  `json:"Permisos"`
	Estatus     string    `json:"Estatus"`
	FechaHora   time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []RolMgo {
	var result []RolMgo
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}
	err = Rols.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Rols.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) RolMgo {
	var result RolMgo
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}
	err = Rols.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []RolMgo {
	var result []RolMgo
	var aux RolMgo
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = RolMgo{}
		Rols.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) RolMgo {
	var result RolMgo
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)

	if err != nil {
		fmt.Println(err)
	}
	err = Rols.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result RolMgo
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}
	err = Rols.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboRols regresa un combo de Rol de mongo
func CargaComboRols(ID string) string {
	Rols := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Rols {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//CargaComboCatalogoArrayID Selecciona Los ObjectsIds de un catalogo
func CargaComboCatalogoArrayID(Clave int, rol RolMgo) string {

	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))

	templ := ``

	for _, v := range rol.Permisos {

		for _, vv := range Catalogo.Valores {
			if v == vv.ID {
				templ += `<option value="` + vv.ID.Hex() + `" selected>` + vv.Valor + `</option>`
			}
		}
	}
	return templ
}

//CargaComboRolsMulti regresa un combo de Rol multi de mongo
func CargaComboRolsMulti(ID string) string {
	Rols := GetAll()
	templ := ``
	for _, v := range Rols {
		templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
	}
	return templ
}

//CargaComboRolesMultiArrayObjID regresa un combo de Rol multi de mongo
func CargaComboRolesMultiArrayObjID(ArrayObID []bson.ObjectId) string {
	Rols := GetAll()
	var templ string
	for _, v := range Rols {

		for _, vv := range ArrayObID {
			if vv == v.ID {
				templ += `<option value="` + vv.Hex() + `" selected>  ` + v.Nombre + ` </option> `
			}
		}
	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Rols []RolMgo) (string, string) {
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>			
					<th>Nombre</th>									
					<th>Descripcion</th>									
					<th>Permisos</th>									
					<th>Estatus</th>									
					<th>FechaHora</th>					
				</tr>`

	for k, v := range Rols {
		//cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Rols/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<tr id = "` + v.ID.Hex() + `">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + v.Descripcion + `</td>`
		cuerpo += `<td>` + "v.Permisos" + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 164) + `</td>`
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

	queryTilde = queryTilde.Field("Permisos")
	queryQuotes = queryQuotes.Field("Permisos")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoRol, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoRol, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
