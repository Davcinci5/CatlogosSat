package CatalogoModel

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Variables"
	"github.com/leekchan/accounting"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//################## < Estructuras de Direccion > ####################
//Paises subestructura de Catalogo
type Paises struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      string        `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	ClavePais   string        `bson:"ClavePais,omitempty"`
}

//Estados subestructura de Catalogo
type Estados struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      string        `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	ClaveEdo    string        `bson:"ClaveEdo,omitempty"`
}

//Municipio subestructura de Catalogo
type Municipio struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Nombre      string        `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	Ciudad      string        `bson:"Ciudad"`
	ClaveMpo    string        `bson:"ClaveMpo"`
	ClaveEdo    string        `bson:"ClaveEdo"`
}

//Colonia subestructura de Catalogo
type Colonia struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Nombre           string        `bson:"Nombre"`
	Descripcion      string        `bson:"Descripcion"`
	CP               string        `bson:"CP"`
	TipoAsentamiento string        `bson:"TipoAsentamiento"`
	Zona             string        `bson:"Zona"`
	ClaveMpo         string        `bson:"ClaveMpo"`
	ClaveEdo         string        `bson:"ClaveEdo"`
	Estatus          bson.ObjectId `bson:"Estatus,omitempty"`
}

//#########################< ESTRUCTURAS >##############################

//CatalogoMgo estructura de Catalogos mongo
type CatalogoMgo struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Clave       int64         `bson:"Clave"`
	Nombre      string        `bson:"Nombre"`
	Descripcion string        `bson:"Descripcion"`
	Valores     []ValoresMgo  `bson:"Valores"`
	Editable    bool          `bson:"Editable"`
	Estatus     bson.ObjectId `bson:"Estatus,omitempty"`
	FechaHora   time.Time     `bson:"FechaHora"`
}

//CatalogoElastic estructura de Catalogos para insertar en Elastic
type CatalogoElastic struct {
	Clave       int64        `json:"Clave"`
	Nombre      string       `json:"Nombre"`
	Descripcion string       `json:"Descripcion"`
	Valores     []ValoresMgo `json:"Valores"`
	Estatus     string       `json:"Estatus"`
	FechaHora   time.Time    `json:"FechaHora"`
}

//ValoresMgo subestructura de Catalogo
type ValoresMgo struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Valor         string        `bson:"Valor"`
	Editable      bool          `bson:"Editable"`
	Estatus       bson.ObjectId `bson:"Estatus,omitempty"`
	Clave         string        `bson:"Clave,omitempty"`
	LocalOFederal string        `bson:"LocalOfederal,omitempty"`
	Retencion     bool          `bson:"Retencion,omitempty"`
	Traslado      bool          `bson:"Traslado,omitempty"`
	TipoImpuesto  bson.ObjectId `bson:"TipoImpuesto,omitempty"`
	Clasificacion bson.ObjectId `bson:"Clasificacion,omitempty"`
	Decimales     int           `bson:"Decimales,omitempty"`
	Variacion     float64       `bson:"Variacion,omitempty"`
	ValorMaximo   float64       `bson:"ValorMaximo,omitempty"`
	AplicaMoral   bool          `bson:"AplicaMoral,omitempty"`
	AplicaFisica  bool          `bson:"AplicaFisica,omitempty"`
}

//#########################< FUNCIONES DE DIRECCION >###############################
//GeyAllPaises
func GetAllPaises() []Paises {
	var result []Paises
	s, Paises, err := MoConexion.GetColectionMgo("Paises")

	if err != nil {
		fmt.Println(err)
	}
	err = Paises.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetAllEstados regeresa todos los estados de Mongo
func GetAllEstados() []Estados {
	var result []Estados
	s, Estados, err := MoConexion.GetColectionMgo("Estados")

	if err != nil {
		fmt.Println(err)
	}
	err = Estados.Find(nil).Sort("Nombre").All(&result)

	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetNameEstado metodo que obtiene el nombre  del estado por el id
func GetNameEstado(ID bson.ObjectId) string {
	var result Estados
	s, Estados, err := MoConexion.GetColectionMgo("Estados")

	if err != nil {
		fmt.Println(err)
	}
	err = Estados.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.Nombre
}

//GetClaveEstado funcion que obtiene la clave de un estado
func GetClaveEstado(ID bson.ObjectId) string {
	var result Estados
	s, Estados, err := MoConexion.GetColectionMgo("Estados")

	if err != nil {
		fmt.Println(err)
	}
	err = Estados.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ClaveEdo
}

//GetAllMunicipiosForClaveEstado funcion que obtiene los municipos por la clave del estado
func GetAllMunicipiosForClaveEstado(ClaveEstado string) []Municipio {

	var results []Municipio
	s, Municipios, err := MoConexion.GetColectionMgo("Municipios")
	if err != nil {
		fmt.Println(err)
	}

	err = Municipios.Find(bson.M{"ClaveEdo": ClaveEstado}).Sort("Nombre").All(&results)

	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return results
}

//GetClaveMunicipio funcion que obtiene la clave de un municipio
func GetClaveMunicipio(ID bson.ObjectId) string {

	var result Municipio
	s, Municipios, err := MoConexion.GetColectionMgo("Municipios")

	if err != nil {
		fmt.Println(err)
	}
	err = Municipios.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ClaveMpo

}

//GetNameMunicipio funcion que obtiene el nombre del minicipio de acuerdo al id
func GetNameMunicipio(ID bson.ObjectId) string {
	var result Municipio
	s, Municipios, err := MoConexion.GetColectionMgo("Municipios")

	if err != nil {
		fmt.Println(err)
	}
	err = Municipios.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.Nombre
}

//GetAllColoniasForClaveMunicipio metodo que obtiene las colonias por clave de municipio
func GetAllColoniasForClaveMunicipio(ClaveMpo string) []Colonia {
	var results []Colonia
	s, Colonias, err := MoConexion.GetColectionMgo("Colonias")
	if err != nil {
		fmt.Println(err)
	}

	err = Colonias.Find(bson.M{"ClaveMpo": ClaveMpo}).Sort("Nombre").All(&results)

	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return results

}

//GetNameColonia funcion que obtiene el nombre de la colonia de acuerdo al id
func GetNameColonia(ID bson.ObjectId) string {
	var result Colonia
	s, Colonias, err := MoConexion.GetColectionMgo("Colonias")

	if err != nil {
		fmt.Println(err)
	}
	err = Colonias.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.Nombre
}

//GetCPForColonia Metodo que obtiene el codogo postal de una colonia
func GetCPForColonia(ID bson.ObjectId) string {
	var result Colonia
	s, Colonias, err := MoConexion.GetColectionMgo("Colonias")

	if err != nil {
		fmt.Println(err)
	}
	err = Colonias.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.CP
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []CatalogoMgo {
	var result []CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
	}
	err = Catalogos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Catalogos.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) CatalogoMgo {
	var result CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
	}
	err = Catalogos.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []CatalogoMgo {
	var result []CatalogoMgo
	var aux CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, value := range Ides {
		aux = CatalogoMgo{}
		Catalogos.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecificByFields(field string, valor interface{}) CatalogoMgo {
	var result CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
	}
	err = Catalogos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetSubEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetSubEspecificByFields(field string, valor interface{}) ValoresMgo {
	var result CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	if err != nil {
		fmt.Println(err)
	}
	err = Catalogos.Find(bson.M{field: valor}).Select(bson.M{"Valores.$": 1, "_id": 0}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	if len(result.Valores) > 0 {
		return result.Valores[0]
	}
	return ValoresMgo{}
}

//ObtenerValoresCatalogo obtiene solamente los valores de un catalogo
//excluye el id y la clave
func ObtenerValoresCatalogo(idCatalogo bson.ObjectId) string {
	result := GetSubEspecificByFields("_id", idCatalogo)
	claves := result.Clave
	return claves
}

//ObtenerValoresCatalogoPorValor obtiene solamente los valores de un catalogo
//excluye el id y la clave
func ObtenerValoresCatalogoPorValor(idCatalogo bson.ObjectId) string {
	result := GetSubEspecificByFields("Valores._id", idCatalogo)
	claves := result.Valor
	return claves
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
	}
	err = Catalogos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//SiguienteClaveDisponible asigna al producto el siguiente secuencial de clave disponible
func SiguienteClaveDisponible() int64 {
	var result CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	if err != nil {
		fmt.Println(err)
	}
	e := Catalogos.Find(nil).Sort("-Clave").Limit(1).One(&result)
	if e != nil {
		fmt.Println(e)
	}

	s.Close()
	return result.Clave + int64(1)
}

//CargaComboCatalogos regresa un combo de unidades de mongo
func CargaComboCatalogos(ID string) string {
	Catalogos := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, v := range Catalogos {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Nombre + `</option>`
		} else {
			templ += `<option value="` + v.ID.Hex() + `">` + v.Nombre + `</option>`
		}

	}
	return templ
}

//RegresaIDEstatusActivo regresa el Id del catálogo de estatus de catálogos donde haya un activo.
//debe especificarse la clave del catálogo.
func RegresaIDEstatusActivo(Clave int) bson.ObjectId {
	var result bson.ObjectId
	Catalogo := GetEspecificByFields("Clave", int64(Clave))
	for _, v := range Catalogo.Valores {
		if strings.ToUpper(v.Valor) == "ACTIVO" {
			result = v.ID
		}
	}
	return result
}

//RegresaValoresCatalogosClave regresa valores de un catalogo con la clave especificada
func RegresaValoresCatalogosClave(clave int) CatalogoMgo {
	var result CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
	}
	err = Catalogos.Find(bson.M{"Clave": clave}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//RegresaNombreSubCatalogo regresa el nombre de un subcatalogo, pasandole como parametro el objectId
func RegresaNombreSubCatalogo(IDValor bson.ObjectId) string {
	var result CatalogoMgo
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)

	if err != nil {
		fmt.Println(err)
	}
	err = Catalogos.Find(bson.M{"Valores._id": IDValor}).One(&result)
	valores := result.Valores
	for _, value := range valores {
		if IDValor == value.ID {
			return value.Valor
		}
	}
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return ""
}

//RegresaNombrePermisos Regresa el valor de un permiso seleccionado, pasandole como identificador el objectId específico del permiso
func RegresaNombrePermisos(objectIds []bson.ObjectId) []string {
	var a []string
	for _, value := range objectIds {
		valor := RegresaNombreSubCatalogo(value)
		a = append(a, valor)
	}
	return a
}

//GetValorMagnitud regresa el valor del elemento de catálogo por clave y id del objeto.
func GetValorMagnitud(ID bson.ObjectId, Clave int) string {
	Magnitudes := GetEspecificByFields("Clave", int64(Clave))
	for _, v := range Magnitudes.Valores {
		if ID == v.ID {
			return v.Valor
		}
	}
	return ""
}

//GetValorDefaultCatalogoBugs regresa el (estatus)valor por defecto al dar de alta un bug
func GetValorDefaultCatalogoBugs(Clave int) string {
	Magnitudes := GetEspecificByFields("Clave", int64(Clave))
	for _, v := range Magnitudes.Valores {
		if v.Valor == "PENDIENTE" {
			return v.ID.Hex()
		}
	}
	return ""
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Catalogos []CatalogoMgo) (string, string) {
	floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
			<th>#</th>
			<th>Clave</th>
			<th>Nombre</th>
			<th>Descripción</th>
			<th>Fecha</th>
			</tr>`

	for k, v := range Catalogos {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Catalogos/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + floats.FormatMoney(v.Clave) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + v.Descripcion + `</td>`
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

	queryTilde = queryTilde.Field("Valores.Valor")
	queryQuotes = queryQuotes.Field("Valores.Valor")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoCatalogo, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoCatalogo, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}
