package FacturacionModel

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

//FacturacionMgo estructura de Facturacions mongo
type FacturacionMgo struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	Operacion       bson.ObjectId `bson:"Operacion,omitempty"`
	Detalle         bson.ObjectId `bson:"Detalle,omitempty"`
	Cliente         bson.ObjectId `bson:"Cliente"`
	DomicilioFiscal DireccionMgo  `bson:"DomicilioFiscal"`
	Estatus         bson.ObjectId `bson:"Estatus,omitempty"`
	FechaHora       time.Time     `bson:"FechaHora"`
}

//DireccionMgo subestructura de Facturacion
type DireccionMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Calle       string        `bson:"Calle"`
	NumInterior string        `bson:"NumInterior"`
	NumExterior string        `bson:"NumExterior"`
	Colonia     bson.ObjectId `bson:"Colonia,omitempty"`
	Municipio   bson.ObjectId `bson:"Municipio,omitempty"`
	Estado      bson.ObjectId `bson:"Estado,omitempty"`
	Pais        bson.ObjectId `bson:"Pais,omitempty"`
	CP          string        `bson:"CP"`
	Estatus     bson.ObjectId `bson:"Estatus,omitempty"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []FacturacionMgo {
	var result []FacturacionMgo
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Facturacions.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) FacturacionMgo {
	var result FacturacionMgo
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Facturacions.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []FacturacionMgo {
	var result []FacturacionMgo
	var aux FacturacionMgo
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = FacturacionMgo{}
		Facturacions.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) FacturacionMgo {
	var result FacturacionMgo
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)

	if err != nil {
		fmt.Println(err)
	}
	err = Facturacions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result FacturacionMgo
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Facturacions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
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

	docs, err = MoConexion.BuscaElastic(MoVar.TipoFacturacion, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoFacturacion, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
