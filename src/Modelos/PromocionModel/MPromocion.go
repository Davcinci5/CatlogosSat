package PromocionModel

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

//PromocionMgo estructura de Promocions mongo
type PromocionMgo struct {
	ID                    bson.ObjectId `bson:"_id,omitempty"`
	Nombre                string        `bson:"Nombre"`
	Descripcion           string        `bson:"Descripcion"`
	PorcentajeDesc        OfertaMgo     `bson:"PorcentajeDesc"`
	PrecioOferta          OfertaMgo     `bson:"PrecioOferta"`
	OfertaMonto           OfertaMgo     `bson:"OfertaMonto"`
	OfertaPiezaPieza      OfertaMgo     `bson:"OfertaPiezaPieza"`
	OfertaPiezaPorcentaje OfertaMgo     `bson:"OfertaPiezaPorcentaje"`
	Estatus               bson.ObjectId `bson:"Estatus,omitempty"`
	FechaInicio           time.Time     `bson:"FechaInicio"`
	FechaFin              time.Time     `bson:"FechaFin"`
}

//PromocionElastic estructura de Promocions para insertar en Elastic
type PromocionElastic struct {
	Nombre                string        `json:"Nombre"`
	Descripcion           string        `json:"Descripcion"`
	PorcentajeDesc        OfertaElastic `json:"PorcentajeDesc"`
	PrecioOferta          OfertaElastic `json:"PrecioOferta"`
	OfertaMonto           OfertaElastic `json:"OfertaMonto"`
	OfertaPiezaPieza      OfertaElastic `json:"OfertaPiezaPieza"`
	OfertaPiezaPorcentaje OfertaElastic `json:"OfertaPiezaPorcentaje"`
	Estatus               string        `json:"Estatus,omitempty"`
	FechaInicio           time.Time     `json:"FechaInicio"`
	FechaFin              time.Time     `json:"FechaFin"`
}

//OfertaElastic subestructura de Promocion
type OfertaElastic struct {
	Cantidad string `json:"Cantidad"`
	Valor    string `json:"Valor"`
}

//OfertaMgo subestructura de Promocion
type OfertaMgo struct {
	Cantidad string `bson:"Cantidad"`
	Valor    string `bson:"Valor"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []PromocionMgo {
	var result []PromocionMgo
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	err = Promocions.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Promocions.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) PromocionMgo {
	var result PromocionMgo
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	err = Promocions.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []PromocionMgo {
	var result []PromocionMgo
	var aux PromocionMgo
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = PromocionMgo{}
		Promocions.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) PromocionMgo {
	var result PromocionMgo
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)

	if err != nil {
		fmt.Println(err)
	}
	err = Promocions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result PromocionMgo
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	err = Promocions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboPromocions regresa un combo de Promocion de mongo
func CargaComboPromocions(ID string) string {
	Promocions := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Promocions {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Promocions []PromocionMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
			<th>#</th>
			
				<th>Nombre</th>					
				
				<th>Descripcion</th>					
				
				<th>Estatus</th>					
				
				<th>FechaInicio</th>					
				
				<th>FechaFin</th>					
				</tr>`

	for k, v := range Promocions {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Promocions/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`

		cuerpo += `<td>` + v.Descripcion + `</td>`

		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 134) + `</td>`

		cuerpo += `<td>` + v.FechaInicio.Format(time.RFC1123) + `</td>`

		cuerpo += `<td>` + v.FechaFin.Format(time.RFC1123) + `</td>`

		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//########## GET NAMES ####################################

//GetNamePromocion regresa el nombre del Promocion con el ID especificado
func GetNamePromocion(ID bson.ObjectId) string {
	var result PromocionMgo
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	Promocions.Find(bson.M{"_id": ID}).One(&result)

	s.Close()
	return result.Nombre
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

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	queryTilde = queryTilde.Field("FechaInicio")
	queryQuotes = queryQuotes.Field("FechaInicio")

	queryTilde = queryTilde.Field("FechaFin")
	queryQuotes = queryQuotes.Field("FechaFin")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoPromocion, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoPromocion, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
