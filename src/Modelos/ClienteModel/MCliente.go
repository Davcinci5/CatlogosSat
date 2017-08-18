package ClienteModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/PersonaModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//ClienteMgo estructura de Clientes mongo
type ClienteMgo struct {
	ID               bson.ObjectId        `bson:"_id,omitempty"`
	TipoCliente      bson.ObjectId        `bson:"TipoCliente,omitempty"`
	IDPersona        bson.ObjectId        `bson:"IDPersona,omitempty"`
	RFC              string               `bson:"RFC"`
	Direcciones      []DireccionMgo       `bson:"Direcciones"`
	MediosDeContacto MediosContactoMgo    `bson:"MediosDeContacto"`
	PersonasContacto []PersonaContactoMgo `bson:"PersonasContacto"`
	Almacenes        []bson.ObjectId      `bson:"Almacenes,omitempty"`
	Notificaciones   []NotificacionMgo    `bson:"Notificaciones"`
	Estatus          bson.ObjectId        `bson:"Estatus,omitempty"`
	FechaHora        time.Time            `bson:"FechaHora"`
}

//ClienteElastic estructura de Clientes para insertar en Elastic
type ClienteElastic struct {
	TipoCliente      string                   `bson:"TipoCliente,omitempty"`
	IDPersona        string                   `json:"IDPersona,omitempty"`
	RFC              string                   `json:"RFC"`
	Direcciones      []DireccionElastic       `json:"Direcciones"`
	MediosDeContacto MediosContactoElastic    `json:"MediosDeContacto"`
	PersonasContacto []PersonaContactoElastic `json:"PersonasContacto"`
	Almacenes        []string                 `json:"Almacenes,omitempty"`
	Notificaciones   []NotificacionElastic    `json:"Notificaciones"`
	Estatus          string                   `json:"Estatus,omitempty"`
	FechaHora        time.Time                `json:"FechaHora"`
}

//DireccionElastic subestructura de Cliente
type DireccionElastic struct {
	Calle       string `json:"Calle"`
	NumInterior string `json:"NumInterior"`
	NumExterior string `json:"NumExterior"`
	Colonia     string `json:"Colonia,omitempty"`
	Municipio   string `json:"Municipio,omitempty"`
	Estado      string `json:"Estado,omitempty"`
	Pais        string `json:"Pais,omitempty"`
	CP          string `json:"CP"`
	Estatus     string `json:"Estatus,omitempty"`
}

//DireccionMgo subestructura de Cliente
type DireccionMgo struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	TipoDireccion bson.ObjectId `TipoDireccion:"_id"`
	Calle         string        `bson:"Calle"`
	NumInterior   string        `bson:"NumInterior"`
	NumExterior   string        `bson:"NumExterior"`
	Colonia       bson.ObjectId `bson:"Colonia,omitempty"`
	Municipio     bson.ObjectId `bson:"Municipio,omitempty"`
	Estado        bson.ObjectId `bson:"Estado,omitempty"`
	Pais          string        `bson:"Pais,omitempty"`
	CP            string        `bson:"CP"`
	Estatus       string        `bson:"Estatus,omitempty"`
}

//MediosContactoElastic subestructura de Cliente
type MediosContactoElastic struct {
	Correos   CorreoElastic   `json:"Correos"`
	Telefonos TelefonoElastic `json:"Telefonos"`
	Otros     []string        `json:"Otros"`
}

//MediosContactoMgo subestructura de Cliente
type MediosContactoMgo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Correos   CorreoMgo     `bson:"Correos"`
	Telefonos TelefonoMgo   `bson:"Telefonos"`
	Otros     []string      `bson:"Otros"`
}

//CorreoMgo subestructura de Cliente
type CorreoMgo struct {
	Principal string   `bson:"Principal"`
	Correos   []string `bson:"Correos"`
}

//TelefonoMgo subestructura de Cliente
type TelefonoMgo struct {
	Principal string   `bson:"Principal"`
	Telefonos []string `bson:"Telefonos"`
}

//CorreoElastic subestructura de Cliente
type CorreoElastic struct {
	Principal string   `json:"Principal"`
	Correos   []string `json:"Correos"`
}

//TelefonoElastic subestructura de Cliente
type TelefonoElastic struct {
	Principal string   `json:"Principal"`
	Telefonos []string `json:"Telefonos"`
}

//NotificacionElastic subestructura de cliente
type NotificacionElastic struct {
	Mensaje          string `json:"Mensaje"`
	Leido            string `json:"Leido"`
	FechaOcurrencia  string `json:"FechaOcurrencia"`
	FechaVencimiento string `json:"FechaVencimiento"`
}

//NotificacionMgo subestructura de Cliente
type NotificacionMgo struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Mensaje          string        `bson:"Mensaje"`
	Leido            string        `bson:"Leido"`
	FechaOcurrencia  string        `bson:"FechaOcurrencia"`
	FechaVencimiento string        `bson:"FechaVencimiento"`
}

//PersonaElastic subestructura de Cliente
type PersonaElastic struct {
	Nombre     string `json:"Nombre"`
	Tipo       string `json:"Tipo,omitempty"`
	Grupos     string `json:"Grupos,omitempty"`
	Predecesor string `json:"Predecesor,omitempty"`
	Estatus    string `json:"Estatus,omitempty"`
	FechaHora  string `json:"FechaHora"`
}

//PersonaMgo subestructura de Cliente
type PersonaMgo struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Nombre       string        `bson:"Nombre"`
	Tipo         string        `bson:"Tipo,omitempty"`
	Grupos       string        `bson:"Grupos,omitempty"`
	Predecesor   string        `bson:"Predecesor,omitempty"`
	Notificacion string        `bson:"Notificacion,omitempty"`
	Estatus      string        `bson:"Estatus,omitempty"`
	FechaHora    string        `bson:"FechaHora"`
}

//PersonaContactoElastic subestructura de Cliente
type PersonaContactoElastic struct {
	Nombre           string                `json:"Nombre"`
	Direcciones      []DireccionElastic    `json:"Direcciones"`
	MediosDeContacto MediosContactoElastic `json:"MediosDeContacto"`
	Almacenes        []AlmacenElastic      `json:"Almacenes"`
	Estatus          string                `json:"Estatus,omitempty"`
}

//PersonaContactoMgo subestructura de Cliente
type PersonaContactoMgo struct {
	ID               bson.ObjectId     `bson:"_id,omitempty"`
	Nombre           string            `bson:"Nombre"`
	Direcciones      []DireccionMgo    `bson:"Direcciones"`
	MediosDeContacto MediosContactoMgo `bson:"MediosDeContacto"`
	Almacenes        []AlmacenMgo      `bson:"Almacenes"`
	Estatus          string            `bson:"Estatus,omitempty"`
}

//AlmacenElastic subestructura de Cliente
type AlmacenElastic struct {
	IDContacto  string `json:"IDContacto,omitempty"`
	IDDireccion string `json:"IDDireccion,omitempty"`
	IDAlmacen   string `json:"IDAlmacen,omitempty"`
}

//AlmacenMgo subestructura de Cliente
type AlmacenMgo struct {
	IDContacto  string `bson:"IDContacto,omitempty"`
	IDDireccion string `bson:"IDDireccion,omitempty"`
	IDAlmacen   string `bson:"IDAlmacen,omitempty"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []ClienteMgo {
	var result []ClienteMgo
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}
	err = Clientes.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Clientes.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) ClienteMgo {
	var result ClienteMgo
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}
	err = Clientes.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []ClienteMgo {
	var result []ClienteMgo
	var aux ClienteMgo
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = ClienteMgo{}
		Clientes.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) ClienteMgo {
	var result ClienteMgo
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)

	if err != nil {
		fmt.Println(err)
	}
	err = Clientes.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result ClienteMgo
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}
	err = Clientes.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Clientes []ClienteMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
				<th>#</th>
				<th>Nombre</th>
				<th>RFC</th>								
				<th>Estatus</th>					
				</tr>`

	for k, v := range Clientes {
		Persona := PersonaModel.GetOne(v.IDPersona)
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Clientes/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + Persona.Nombre + `</td>`
		cuerpo += `<td>` + v.RFC + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 137) + `</td>`

		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

// NombreCliente Regresa el nombre del cliente, devuelve false si no existe
func NombreCliente(ID bson.ObjectId) (bool, string) {
	Cli := GetOne(ID)
	if !MoGeneral.EstaVacio(Cli) {
		Perso := PersonaModel.GetOne(Cli.IDPersona)
		if !MoGeneral.EstaVacio(Perso) {
			return true, Perso.Nombre
		}
	}
	return false, ""
}

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("RFC")
	queryQuotes = queryQuotes.Field("RFC")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoCliente, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoCliente, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
