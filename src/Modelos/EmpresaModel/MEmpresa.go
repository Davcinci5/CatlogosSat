package EmpresaModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	elastic "gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//EmpresaMgo estructura de Empresas mongo
type EmpresaMgo struct {
	ID               bson.ObjectId          `bson:"_id,omitempty"`
	DatosComerciales ComercialMgo           `bson:"DatosComerciales,omitempty"`
	DatosFiscales    FiscalMgo              `bson:"DatosFiscales,omitempty"`
	DatosFactura     FacturaMgo             `bson:"DatosFactura,omitempty"`
	CorreoYEnvio     ConfiguracionCorreoMgo `bson:"CorreoYEnvio,omitempty"`
}

//ComercialMgo subestructura de Empresa
type ComercialMgo struct {
	Nombre    string       `bson:"Nombre"`
	Domicilio DireccionMgo `bson:"Domicilio"`
	Contactos ContactoMgo  `bson:"Contactos"`
}

//FiscalMgo subestructura de Empresa
type FiscalMgo struct {
	RazonSocial string        `bson:"RazonSocial"`
	RFC         string        `bson:"RFC"`
	Domicilio   DireccionMgo  `bson:"Domicilio"`
	Contactos   []ContactoMgo `bson:"Contactos"`
}

//FacturaMgo subestructura de Empresa
type FacturaMgo struct {
	Key string `bson:"Key,omitempty"`
	Cer string `bson:"Cer,omitempty"`
	Pem string `bson:"Pem,omitempty"`
}

//ConfiguracionCorreoMgo subestructura de Empresa
type ConfiguracionCorreoMgo struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Correo  string        `bson:"Correo"`
	Pass    string        `bson:"Pass"`
	Tipo    string        `bson:"Tipo"`
	Puerto  string        `bson:"Puerto"`
	Cifrado string        `bson:"Cifrado,omitempty"`
}

//DireccionMgo subestructura de Empresa
type DireccionMgo struct {
	Calle       string `bson:"Calle"`
	NumInterior string `bson:"NumInterior,omitempty"`
	NumExterior string `bson:"NumExterior"`
	Colonia     string `bson:"Colonia"`
	Municipio   string `bson:"Municipio"`
	Estado      string `bson:"Estado"`
	Pais        string `bson:"Pais"`
	CP          string `bson:"CP"`
}

//ContactoMgo subestructura de Empresa
type ContactoMgo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Alias     string        `bson:"Alias,omitempty"`
	Email     string        `bson:"Email,omitempty"`
	Telefono  string        `bson:"Telefono,omitempty"`
	Movil     string        `bson:"Movil,omitempty"`
	Domicilio DireccionMgo  `bson:"Domicilio,omitempty"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []EmpresaMgo {
	var result []EmpresaMgo
	s, Empresas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEmpresa)
	if err != nil {
		fmt.Println(err)
	}
	err = Empresas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) EmpresaMgo {
	var result EmpresaMgo
	s, Empresas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEmpresa)
	if err != nil {
		fmt.Println(err)
	}
	err = Empresas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []EmpresaMgo {
	var result []EmpresaMgo
	var aux EmpresaMgo
	s, Empresas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEmpresa)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = EmpresaMgo{}
		Empresas.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) EmpresaMgo {
	var result EmpresaMgo
	s, Empresas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEmpresa)

	if err != nil {
		fmt.Println(err)
	}
	err = Empresas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result EmpresaMgo
	s, Empresas, err := MoConexion.GetColectionMgo(MoVar.ColeccionEmpresa)
	if err != nil {
		fmt.Println(err)
	}
	err = Empresas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()

	return result.ID
}

// //CargaComboEmpresas regresa un combo de Empresa de mongo
// func CargaComboEmpresas(ID string) string {
// 	Empresas := GetAll()

// 	templ := ``

// 	if ID != "" {
// 		templ =  `<option value="">--SELECCIONE--</option> `
// 	} else {
// 		templ =  `<option value="" selected>--SELECCIONE--</option> `
// 	}

// 	for _, v := range Empresas {
// 		if ID == v.ID.Hex() {
// 			templ += `<option value=" ` + v.ID.Hex() +  `" selected>  ` + v.Nombre +  ` </option> `
// 		} else {
// 			templ += `<option value=" ` + v.ID.Hex() +  `">  ` + v.Nombre +  ` </option> `
// 		}

// 	}
// 	return templ
// }

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoEmpresa, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoEmpresa, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
