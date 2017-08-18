package ListaCostoModel
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


//ListaCostoMgo estructura de ListaCostos mongo
type ListaCostoMgo struct {
 ID    bson.ObjectId   `bson:"_id,omitempty"`
 Nombre   string  `bson:"Nombre"`
 Descripcion   string  `bson:"Descripcion"`
 GrupoP   []bson.ObjectId  `bson:"GrupoP"`
 Estatus   bson.ObjectId  `bson:"Estatus"`
 FechaHora   time.Time  `bson:"FechaHora"`
}



	//#########################< FUNCIONES GENERALES MGO >###############################

	//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
		func GetAll() []ListaCostoMgo {
			var result []ListaCostoMgo
			s, ListaCostos, err := MoConexion.GetColectionMgo(MoVar.ColeccionListaCosto)
			if err != nil {
				fmt.Println(err)
			}
			err = ListaCostos.Find(nil).All(&result)
			if err != nil {
				fmt.Println(err)
			}
			s.Close()
			return result
		}

		//GetOne Regresa un documento específico de Mongo (Por Coleccion)
		func GetOne(ID bson.ObjectId) ListaCostoMgo {
			var result ListaCostoMgo
			s, ListaCostos, err := MoConexion.GetColectionMgo(MoVar.ColeccionListaCosto)
			if err != nil {
				fmt.Println(err)
			}
			err = ListaCostos.Find(bson.M{"_id": ID}).One(&result)
			if err != nil {
				fmt.Println(err)
			}
			s.Close()
			return result
		}

		//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
		func GetEspecifics(Ides []bson.ObjectId) []ListaCostoMgo {
			var result []ListaCostoMgo
			var aux ListaCostoMgo
			s, ListaCostos, err := MoConexion.GetColectionMgo(MoVar.ColeccionListaCosto)
			if err != nil {
				fmt.Println(err)
			}
			for _, value := range Ides {
				aux = ListaCostoMgo{}
				ListaCostos.Find(bson.M{"_id": value}).One(&aux)
				result = append(result, aux)
			}
			s.Close()
			return result
		}

		//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
		func GetEspecificByFields(field string, valor interface{}) ListaCostoMgo {
			var result ListaCostoMgo
			s, ListaCostos, err := MoConexion.GetColectionMgo(MoVar.ColeccionListaCosto)

			if err != nil {
				fmt.Println(err)
			}
			err = ListaCostos.Find(bson.M{field: valor}).One(&result)
			if err != nil {
				fmt.Println(err)
			}
			s.Close()
			return result
		}

		//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
		func GetIDByField(field string, valor interface{}) bson.ObjectId {
			var result ListaCostoMgo
			s, ListaCostos, err := MoConexion.GetColectionMgo(MoVar.ColeccionListaCosto)
			if err != nil {
				fmt.Println(err)
			}
			err = ListaCostos.Find(bson.M{field: valor}).One(&result)
			if err != nil {
				fmt.Println(err)
			}
			s.Close()
			return result.ID
		}

		//CargaComboListaCostos regresa un combo de ListaCosto de mongo
		func CargaComboListaCostos(ID string) string {
			ListaCostos := GetAll()

			templ := ``

			if ID != "" {
				templ =  `<option value="">--SELECCIONE--</option> `
			} else {
				templ =  `<option value="" selected>--SELECCIONE--</option> `
			}

			for _, v := range ListaCostos {
				if ID == v.ID.Hex() {
					templ += `<option value=" ` + v.ID.Hex() +  `" selected>  ` + v.Nombre +  ` </option> `
				} else {
					templ += `<option value=" ` + v.ID.Hex() +  `">  ` + v.Nombre +  ` </option> `
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

	docs, err = MoConexion.BuscaElastic(MoVar.TipoListaCosto, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoListaCosto, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
				