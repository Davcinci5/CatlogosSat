package GrupoAlmacenModel

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

//GrupoAlmacenMgo estructura de GrupoAlmacens mongo
type GrupoAlmacenMgo struct {
	ID            bson.ObjectId   `bson:"_id,omitempty"`
	Nombre        string          `bson:"Nombre"`
	Descripcion   string          `bson:"Descripcion"`
	PermiteVender bool            `bson:"PermiteVender"`
	Miembros      []bson.ObjectId `bson:"Miembros,omitempty"`
	Estatus       bson.ObjectId   `bson:"Estatus"`
	FechaHora     time.Time       `bson:"FechaHora"`
}

//GrupoAlmacenElastic estructura de GrupoAlmacens para insertar en Elastic
type GrupoAlmacenElastic struct {
	Nombre        string    `json:"Nombre"`
	Descripcion   string    `json:"Descripcion"`
	PermiteVender string    `json:"PermiteVender"`
	Miembros      []string  `json:"Miembros,omitempty"`
	Estatus       string    `json:"Estatus"`
	FechaHora     time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []GrupoAlmacenMgo {
	var result []GrupoAlmacenMgo
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = GrupoAlmacens.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)

	if err != nil {
		fmt.Println(err)
	}
	result, err = GrupoAlmacens.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) GrupoAlmacenMgo {
	var result GrupoAlmacenMgo
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = GrupoAlmacens.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []GrupoAlmacenMgo {
	var result []GrupoAlmacenMgo
	var aux GrupoAlmacenMgo
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = GrupoAlmacenMgo{}
		GrupoAlmacens.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) GrupoAlmacenMgo {
	var result GrupoAlmacenMgo
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)

	if err != nil {
		fmt.Println(err)
	}
	err = GrupoAlmacens.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result GrupoAlmacenMgo
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = GrupoAlmacens.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboGrupoAlmacens regresa un combo de GrupoAlmacen de mongo
func CargaComboGrupoAlmacens(ID string) string {
	GrupoAlmacens := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range GrupoAlmacens {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(GrupoAlmacens []GrupoAlmacenMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
			<th>#</th>
			
				<th>Nombre</th>					
				
				<th>Descripcion</th>					
				
				<th>PermiteVender</th>					
				
				<th>Estatus</th>					
				
				<th>FechaHora</th>					
				</tr>`

	for k, v := range GrupoAlmacens {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/GrupoAlmacens/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`

		cuerpo += `<td>` + v.Descripcion + `</td>`

		cuerpo += `<td>` + strconv.FormatBool(v.PermiteVender) + `</td>`

		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 132) + `</td>`

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

	queryTilde = queryTilde.Field("PermiteVender")
	queryQuotes = queryQuotes.Field("PermiteVender")

	queryTilde = queryTilde.Field("Miembros")
	queryQuotes = queryQuotes.Field("Miembros")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupoAlmacen, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupoAlmacen, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}

//CargaComboGrupoAlmacenArray Recibe un arreglo IDs y regresa los option con los ids del arreglo seleccionados
func CargaComboGrupoAlmacenArray(ArrayObID []string) string {
	GrupoAlmacen := GetAll()
	var templ string
	for _, v := range GrupoAlmacen {
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
