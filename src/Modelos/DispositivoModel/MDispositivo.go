package DispositivoModel

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

//DispositivoMgo estructura de Dispositivos mongo
type DispositivoMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      string        `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	Predecesor  bson.ObjectId `bson:"Predecesor"`
	Mac         string        `bson:"Mac"`
	Campos      []string      `bson:"Campos"`
	Estatus     bson.ObjectId `bson:"Estatus"`
	FechaHora   time.Time     `bson:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []DispositivoMgo {
	var result []DispositivoMgo
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}
	err = Dispositivos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) DispositivoMgo {
	var result DispositivoMgo
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}
	err = Dispositivos.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []DispositivoMgo {
	var result []DispositivoMgo
	var aux DispositivoMgo
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = DispositivoMgo{}
		Dispositivos.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) DispositivoMgo {
	var result DispositivoMgo
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)

	if err != nil {
		fmt.Println(err)
	}
	err = Dispositivos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result DispositivoMgo
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}
	err = Dispositivos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboDispositivos regresa un combo de Dispositivo de mongo
func CargaComboDispositivos(ID string) string {
	Dispositivos := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Dispositivos {
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

	docs, err = MoConexion.BuscaElastic(MoVar.TipoDispositivo, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoDispositivo, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
