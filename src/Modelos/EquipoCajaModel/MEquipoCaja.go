package EquipoCajaModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modelos/CatalogoModel"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//EquipoCajaMgo estructura de EquipoCajas mongo
type EquipoCajaMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      string        `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	Usuario     bson.ObjectId `bson:"Usuario,omitempty"`
	Dispositivo bson.ObjectId `bson:"Dispositivo,omitempty"`
	Estatus     bson.ObjectId `bson:"Estatus,omitempty"`
	FechaHora   time.Time     `bson:"FechaHora"`
}

//EquipoCajaElastic estructura de EquipoCajas para insertar en Elastic
type EquipoCajaElastic struct {
	Nombre      string    `json:"Nombre"`
	Descripcion string    `json:"Descripcion"`
	Usuario     string    `json:"Usuario"`
	Dispositivo string    `json:"Dispositivo"`
	Estatus     string    `json:"Estatus"`
	FechaHora   time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []EquipoCajaMgo {
	var result []EquipoCajaMgo
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = EquipoCajas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)

	if err != nil {
		fmt.Println(err)
	}
	result, err = EquipoCajas.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) EquipoCajaMgo {
	var result EquipoCajaMgo
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = EquipoCajas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []EquipoCajaMgo {
	var result []EquipoCajaMgo
	var aux EquipoCajaMgo
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = EquipoCajaMgo{}
		EquipoCajas.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) EquipoCajaMgo {
	var result EquipoCajaMgo
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)

	if err != nil {
		fmt.Println(err)
	}
	err = EquipoCajas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result EquipoCajaMgo
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = EquipoCajas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboEquipoCajas regresa un combo de EquipoCaja de mongo
func CargaComboEquipoCajas(ID string) string {
	EquipoCajas := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range EquipoCajas {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//CargaComboCajasMulti regresa un combo de Caja de mongo
func CargaComboCajasMulti() string {
	Cajas := GetAll()

	templ := ``

	for _, v := range Cajas {
		templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
	}

	return templ
}

//CargaComboCajasMultiArrayObjID Carga el combo multi de cajas seleccionando los valores en el arreglo de ids
func CargaComboCajasMultiArrayObjID(ArrayObID []string) string {
	Cajas := GetAll()
	var templ string
	for _, v := range Cajas {
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

//CargaNombreCaja Regresa el nombre de una caja por Id
func CargaNombreCaja(ID bson.ObjectId) string {
	var result EquipoCajaMgo
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = EquipoCajas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.Nombre
}

//CargaNombresCajas Regresa el nombre de un arreglo de cajas por id
func CargaNombresCajas(IDS []bson.ObjectId) []string {
	var result []string
	for _, val := range IDS {
		result = append(result, CargaNombreCaja(val))
	}
	return result
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(EquipoCajas []EquipoCajaMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>			
					<th>Nombre</th>									
					<th>Descripcion</th>									
					<th>Usuario</th>									
					<th>Dispositivo</th>									
					<th>Estatus</th>									
					<th>FechaHora</th>					
				  </tr>`

	for k, v := range EquipoCajas {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/EquipoCajas/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + v.Descripcion + `</td>`
		cuerpo += `<td>` + v.Usuario.Hex() + `</td>`
		cuerpo += `<td>` + v.Dispositivo.Hex() + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 135) + `</td>`
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

	queryTilde = queryTilde.Field("Usuario")
	queryQuotes = queryQuotes.Field("Usuario")

	queryTilde = queryTilde.Field("Dispositivo")
	queryQuotes = queryQuotes.Field("Dispositivo")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoEquipoCaja, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoEquipoCaja, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}

//GetCajaAbiertaByUsuario devuelve caja abierta para un usuario
func GetCajaAbiertaByUsuario(IDUsuario bson.ObjectId, IDEstatus bson.ObjectId) EquipoCajaMgo {
	var result EquipoCajaMgo
	s, EquipoCajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEquipoCaja)
	if err != nil {
		fmt.Println(err)
	}
	err = EquipoCajas.Find(bson.M{"Usuario": IDUsuario, "Estatus": IDEstatus}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}
