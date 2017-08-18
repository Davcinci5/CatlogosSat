package ImpuestoModel

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

//TipoFactor estrcutura de catálogo del sistema
type TipoFactor struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Descripcion string        `bson:"DescripcionTipoFactor"`
}

//TipoImpuesto estrcutura de catálogo del sistema
type TipoImpuesto struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	LocalOFederal string        `bson:"LocalOfederal"`
	Retencion     bool          `bson:"Retencion"`
	Clave         string        `bson:"Clave"`
	Traslado      bool          `bson:"Traslado"`
	Descripcion   string        `bson:"Descripcion"`
}

//ClasificacionImpuesto estrcutura de catálogo del sistema
type ClasificacionImpuesto struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Descripcion  string        `bson:"Descripcion"`
	TipoImpuesto bson.ObjectId `bson:"TipoImpuesto,omitempty"`
}

//SubClasificacionImpuesto estrcutura de catálogo del sistema
type SubClasificacionImpuesto struct {
	ID            bson.ObjectId   `bson:"_id,omitempty"`
	Descripcion   string          `bson:"Descripcion"`
	Clasificacion bson.ObjectId   `bson:"Clasificacion,omitempty"`
	Datos         DataImpuestoMgo `bson:"Impuestos,omitempty"`
}

//ImpuestoMgo estructura de Impuestos mongo
type ImpuestoMgo struct {
	ID               bson.ObjectId   `bson:"_id,omitempty"`
	Descripcion      string          `bson:"Descripcion"`
	TipoImpuesto     bson.ObjectId   `bson:"TipoImpuesto,omitempty"`
	Clasificacion    bson.ObjectId   `bson:"Clasificacion,omitempty"`
	SubClasificacion bson.ObjectId   `bson:"SubClasificacion,omitempty"`
	Datos            DataImpuestoMgo `bson:"Impuestos,omitempty"`
	Estatus          bson.ObjectId   `bson:"Estatus,omitempty"`
	Editable         bool            `bson:"Editable"`
	FechaHora        time.Time       `bson:"FechaHora"`
}

//ImpuestoElastic estructura de Impuestos para insertar en Elastic
type ImpuestoElastic struct {
	Descripcion      string              `json:"Descripcion"`
	Clasificacion    string              `json:"Clasificacion"`
	SubClasificacion string              `json:"SubClasificacion"`
	Datos            DataImpuestoElastic `json:"Impuestos"`
	Estatus          string              `json:"Estatus"`
	Editable         bool                `json:"Editable"`
	FechaHora        time.Time           `json:"FechaHora"`
}

//DataImpuestoElastic subestructura de Impuesto
type DataImpuestoElastic struct {
	Nombre     string    `json:"Nombre"`
	Max        float64   `json:"ValorMaximo"`
	Min        float64   `json:"ValorMinimo"`
	TipoFactor string    `json:"Factor"`
	Unidad     string    `json:"Unidad"`
	FechaHora  time.Time `json:"FechaHora"`
}

//DataImpuestoMgo subestructura de Impuesto
type DataImpuestoMgo struct {
	Nombre     string        `bson:"Impuesto"`
	Max        float64       `bson:"ValorMaximo"`
	Min        float64       `bson:"ValorMinimo"`
	Retencion  bool          `bson:"Retencion"`
	Traslado   bool          `bson:"Traslado"`
	TipoFactor bson.ObjectId `bson:"Factor,omitempty"`
	Rango      bool          `bson:"RangoOFijo"`
	Unidad     bson.ObjectId `bson:"Unidad,omitempty"`
	FechaHora  time.Time     `bson:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []ImpuestoMgo {
	var result []ImpuestoMgo
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetAllTiposDeImpuestos Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAllTiposDeImpuestos() []TipoImpuesto {
	var result []TipoImpuesto
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.CatSysTipoImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetAllTiposDeFactores Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAllTiposDeFactores() []TipoFactor {
	var result []TipoFactor
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.CatSysTipoFactor)
	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetAllClasificacionDeImpuestos Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAllClasificacionDeImpuestos() []ClasificacionImpuesto {
	var result []ClasificacionImpuesto
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.CatSysClasificacionDeimpuestos)
	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetAllSubClasificacionDeImpuestos Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAllSubClasificacionDeImpuestos() []SubClasificacionImpuesto {
	var result []SubClasificacionImpuesto
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.CatSysSubClasificacionDeimpuestos)
	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Impuestos.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) ImpuestoMgo {
	var result ImpuestoMgo
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []ImpuestoMgo {
	var result []ImpuestoMgo
	var aux ImpuestoMgo
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = ImpuestoMgo{}
		Impuestos.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) ImpuestoMgo {
	var result ImpuestoMgo
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)

	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetManyEspecificByFields regresa muchos documentos de Mongo especificando un campo y un determinado valor
func GetManyEspecificByFields(field string, valor interface{}) []ImpuestoMgo {
	var result []ImpuestoMgo
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)

	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(bson.M{field: valor}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result ImpuestoMgo
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//GetImpuestosEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetImpuestosEspecificByFields(field string, valor interface{}) []ImpuestoMgo {
	var result []ImpuestoMgo
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)

	if err != nil {
		fmt.Println(err)
	}
	err = Impuestos.Find(bson.M{field: valor}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CargaComboImpuestos regresa un combo de Impuesto de mongo
func CargaComboImpuestos(ID string) string {
	Impuestos := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Impuestos {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Descripcion + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Descripcion + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Impuestos []ImpuestoMgo) (string, string) {
	floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>					
					<th>Grupo</th>					
					<th>Descripcion</th>
					<th>Clasificación</th>
					<th>Valor Mínimo</th>
					<th>Valor Máximo</th>
					<th>Factor</th>
					<th>Aplica A</th>					
					<th>Fecha</th>			
					<th>Estatus</th>		
				</tr>`

	for k, v := range Impuestos {

		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Impuestos/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.ID.Hex() + `</td>`
		cuerpo += `<td>` + v.Descripcion + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Clasificacion, 148) + `</td>`
		cuerpo += `<td>` + floats.FormatMoney(v.Datos.Min) + `</td>`
		cuerpo += `<td>` + floats.FormatMoney(v.Datos.Max) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Datos.TipoFactor, 150) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Datos.Unidad, 151) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 147) + `</td>`
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

	queryTilde = queryTilde.Field("Clasificacion")
	queryQuotes = queryQuotes.Field("Clasificacion")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoImpuesto, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoImpuesto, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
