package UnidadModel

import (
	"fmt"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//UnidadMgo estructura de Unidads mongo
type UnidadMgo struct {
	ID          bson.ObjectId   `bson:"_id,omitempty"`
	Magnitud    string          `bson:"Magnitud"`
	Descripcion string          `bson:"Descripcion"`
	Datos       []DataUnidadMgo `bson:"Datos"`
	Estatus     bson.ObjectId   `bson:"Estatus,omitempty"`
	FechaHora   time.Time       `bson:"FechaHora"`
}

//DataUnidadMgo subestructura de Unidad
type DataUnidadMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      string        `bson:"Nombre"`
	Abreviatura string        `bson:"Abreviatura"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []UnidadMgo {
	var result []UnidadMgo
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	err = Unidads.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) UnidadMgo {
	var result UnidadMgo
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	err = Unidads.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []UnidadMgo {
	var result []UnidadMgo
	var aux UnidadMgo
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = UnidadMgo{}
		Unidads.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetSubByField regresa un documento de Mongo especificando un subdocumento y un campo con un determinado valor
func GetSubByField(sub string, field string, valor interface{}) UnidadMgo {
	var result UnidadMgo
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)

	if err != nil {
		fmt.Println(err)
	}

	err = Unidads.Find(nil).Select(bson.M{sub: bson.M{"$elemMatch": bson.M{field: valor}}}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result UnidadMgo
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	err = Unidads.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) DataUnidadMgo {
	var result DataUnidadMgo
	s, Unidades, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)

	if err != nil {
		fmt.Println(err)
	}
	err = Unidades.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetSubEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetSubEspecificByFields(field string, valor interface{}) DataUnidadMgo {
	var result UnidadMgo
	s, Unidades, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	err = Unidades.Find(bson.M{field: valor}).Select(bson.M{"Datos.$": 1, "_id": 0}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	if len(result.Datos) > 0 {
		return result.Datos[0]
	}
	return DataUnidadMgo{}
}

//GetNombreUnidadByField regresa la abreviatura de una unidad específica
func GetNombreUnidadByField(field string, valor interface{}) string {
	Data := GetSubEspecificByFields(field, valor)
	return Data.Abreviatura
}

//CargaComboMagnitudes crea el combo de magnitudes de la coleccion de catalogos
func CargaComboMagnitudes(Magnitudes CatalogoModel.CatalogoMgo, ID string) string {
	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, v := range Magnitudes.Valores {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Valor + `</option>`
		} else {
			templ += `<option value="` + v.ID.Hex() + `">` + v.Valor + `</option>`
		}
	}
	return templ
}

//CargaComboUnidades regresa un combo de unidades de mongo
func CargaComboUnidades(ID string) string {
	Unidades := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, v := range Unidades {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Magnitud + `</option>`
		} else {
			templ += `<option value="` + v.ID.Hex() + `">` + v.Magnitud + `</option>`
		}

	}
	return templ
}

//GetNameMagnitud regresa el nombre del Unidad con el ID especificado
func GetNameMagnitud(ID string) string {
	Magnitudes := CatalogoModel.GetEspecificByFields("Clave", int64(166))
	for _, v := range Magnitudes.Valores {
		if ID == v.ID.Hex() {
			return v.Valor
		}
	}
	return ""
}

//RegresaNombreUnidad regresa el nombre de la unidad, pasandole como parametro el identificador unico de la unidad
func RegresaNombreUnidad(IDUnidad bson.ObjectId) string {
	var result UnidadMgo
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)

	if err != nil {
		fmt.Println(err)
	}
	err = Unidads.Find(bson.M{"Datos._id": IDUnidad}).One(&result)
	valores := result.Datos
	for _, value := range valores {
		if IDUnidad == value.ID {
			if value.Abreviatura != "" {
				return value.Abreviatura
			}
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return ""
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

	docs, err = MoConexion.BuscaElastic(MoVar.TipoUnidad, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoUnidad, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
