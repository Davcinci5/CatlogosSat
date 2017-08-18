package KitModel

import (
	"fmt"
	"time"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//KitMgo estructura de Kits mongo
type KitMgo struct {
	ID           bson.ObjectId   `bson:"_id,omitempty"`
	Nombre       string          `bson:"Nombre"`
	Codigo       string          `bson:"Codigo"`
	Tipo         bson.ObjectId   `bson:"Tipo"`
	Aplicación   bson.ObjectId   `bson:"Aplicación"`
	Imagenes     []bson.ObjectId `bson:"Imagenes,omitempty"`
	Conformacion []ConformaMgo   `bson:"Conformacion"`
	Estatus      bson.ObjectId   `bson:"Estatus"`
	FechaHora    time.Time       `bson:"FechaHora"`
}

//ConformaMgo subestructura de Kit
type ConformaMgo struct {
	Almacen   bson.ObjectId     `bson:"Almacen"`
	Productos []DataProductoMgo `bson:"Productos"`
}

//DataProductoMgo subestructura de Kit
type DataProductoMgo struct {
	IDProducto bson.ObjectId `bson:"IDProducto"`
	Cantidad   float64       `bson:"Cantidad"`
	Precio     float64       `bson:"Precio,omitempty"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []KitMgo {
	var result []KitMgo
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}
	err = Kits.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) KitMgo {
	var result KitMgo
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}
	err = Kits.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []KitMgo {
	var result []KitMgo
	var aux KitMgo
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = KitMgo{}
		Kits.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) KitMgo {
	var result KitMgo
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)

	if err != nil {
		fmt.Println(err)
	}
	err = Kits.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result KitMgo
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}
	err = Kits.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboKits regresa un combo de Kit de mongo
func CargaComboKits(ID string) string {
	Kits := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Kits {
		if ID == v.ID.Hex() {
			templ += `<option value=" ` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value=" ` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
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

	docs, err = MoConexion.BuscaElastic(MoVar.TipoKit, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoKit, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
