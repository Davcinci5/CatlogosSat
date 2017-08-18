package AlmacenModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/ConexionModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//AlmacenMgo estructura de Almacens mongo
type AlmacenMgo struct {
	ID            bson.ObjectId   `bson:"_id,omitempty"`
	Nombre        string          `bson:"Nombre"`
	Tipo          bson.ObjectId   `bson:"Tipo"`
	Clasificacion bson.ObjectId   `bson:"Clasificacion,omitempty"`
	Predecesor    bson.ObjectId   `bson:"Predecesor,omitempty"`
	Sucesor       []bson.ObjectId `bson:"Sucesor,omitempty"`
	ListaCosto    []bson.ObjectId `bson:"ListaCosto,omitempty"`
	ListaPrecio   []bson.ObjectId `bson:"ListaPrecio,omitempty"`
	Direccion     bson.ObjectId   `bson:"Direccion,omitempty"`
	Grupos        []bson.ObjectId `bson:"Grupos,omitempty"`
	Conexion      bson.ObjectId   `bson:"Conexion"`
	Estatus       bson.ObjectId   `bson:"Estatus"`
	FechaHora     time.Time       `bson:"FechaHora,omitempty"`
}

//AlmacenElastic estructura de Almacens para insertar en Elastic
type AlmacenElastic struct {
	Nombre         string    `json:"Nombre"`
	Tipo           string    `json:"Tipo"`
	Clasificacion  string    `json:"Clasificacion"`
	Predecesor     string    `json:"Predecesor,omitempty"`
	Sucesor        []string  `json:"Sucesor,omitempty"`
	Direccion      string    `json:"Direccion,omitempty"`
	Grupos         []string  `json:"Grupos,omitempty"`
	NombreConexion string    `json:"NombreConexion"`
	Estatus        string    `json:"Estatus"`
	FechaHora      time.Time `json:"FechaHora"`
}

//DatosConexionMgo subestructura de Almacen
type DatosConexionMgo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Servidor  string        `bson:"Servidor"`
	NombreBD  string        `bson:"NombreBD"`
	UsuarioBD string        `bson:"UsuarioBD"`
	PassBD    string        `bson:"PassBD"`
}

//MovimientoAlmacenes estructura para arreglos de almacenes
type MovimientoAlmacenes struct {
	Almacenes []AlmacenMgo
}

//ImpuestoAlmacen estrcutura de impuestos de almacen en Psql
type ImpuestoAlmacen struct {
	Tratamiento string
	Tipo        string
	Factor      string
	Valor       float64
}

//#########################< FUNCIONES GENERALES MGO >###############################

//CargaComboAlamcenesArray Recibe un arreglo IDs y regresa los option con los ids del arreglo seleccionados
func CargaComboAlamcenesArray(ArrayObID []string) string {
	Almacenes := GetAll()
	templ := ""
	for _, v := range Almacenes {
		existe := false
		for _, vv := range ArrayObID {
			if vv == v.ID.Hex() {
				existe = true
			}
		}
		if existe {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}
	}
	return templ
}

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []AlmacenMgo {
	var result []AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Almacens.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) AlmacenMgo {
	var result AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []AlmacenMgo {
	var result []AlmacenMgo
	var aux AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = AlmacenMgo{}
		Almacens.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) AlmacenMgo {
	var result AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)

	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecificsByTag regresa un arreglo de almacenes pasandole la etiqueta y el valor a buscar
func GetEspecificsByTag(tag string, valor interface{}) []AlmacenMgo {
	var result []AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)

	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(bson.M{tag: valor}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetAllInQuery Regresa todos los documentos existentes de Mongo que coincidan con la consulta "query" (Por Coleccion)
func GetAllInQuery(query interface{}) []AlmacenMgo {
	var result []AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecificsByTagAndTestConexion regresa un arreglo de almacenes siempre y cuando responda el servidor correspondiente.
func GetEspecificsByTagAndTestConexion(tag string, valor interface{}) ([]AlmacenMgo, []string) {
	var result []AlmacenMgo
	var almVerdaderos []AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)

	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(bson.M{tag: valor}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	var errores []string
	if len(result) > 0 {
		for _, value := range result {
			if value.Conexion != "" {
				//datosConexion := ConexionModel.GetOne(value.Conexion)
				Conexion := ConexionModel.GetOne(value.Conexion)
				var parametros MoConexion.ParametrosConexionPostgres
				parametros.Servidor = Conexion.Servidor
				parametros.Usuario = Conexion.UsuarioBD
				parametros.Pass = Conexion.PassBD
				parametros.NombreBase = Conexion.NombreBD
				ping, err := MoConexion.ConexioServidorAlmacenPing(parametros)
				if err != nil {
					fmt.Println("Error al intentar conectar con el almacen: ", err)
				}
				if ping {
					almVerdaderos = append(almVerdaderos, value)
				} else {
					errores = append(errores, value.Nombre+": "+err.Error())
				}
			}
		}
		s.Close()
		return almVerdaderos, errores
	}
	s.Close()
	return almVerdaderos, errores

}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//AlmacenFields Regresa el ID, el Nombre y el Tipo de Almacen
func AlmacenFields() []AlmacenMgo {
	var result []AlmacenMgo
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	err = Almacens.Find(nil).Select(bson.M{"_id": 1, "Nombre": 1, "Tipo": 1}).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CargaComboGrupoAlmacenesMulti Regresa
func CargaComboGrupoAlmacenesMulti(ID string) string {
	Almacens := GetAll()
	templ := ``
	for _, v := range Almacens {

		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}
	}
	return templ
}

//CargaComboAlmacens regresa un combo de Almacen de mongo
func CargaComboAlmacens(ID string) string {
	Almacens := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Almacens {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Almacens []AlmacenMgo) (string, string) {
	cuerpo := ``
	cabecera := `<tr>
				<th>#</th>			
				<th>Nombre</th>									
				<th>Tipo</th>									
				<th>Clasificacion</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

	for k, v := range Almacens {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Almacens/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo, 132) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Clasificacion, 133) + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 134) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//ObtenerParametrosConexion Obtiene los datos de una conexion pasandole el identificador del almacén
func ObtenerParametrosConexion(idAlmacen bson.ObjectId) ConexionModel.ConexionMgo {
	datosAlmacen := GetOne(idAlmacen)
	var datosConexion ConexionModel.ConexionMgo
	if datosAlmacen.Conexion != "" {
		datosConexion = ConexionModel.GetOne(datosAlmacen.Conexion)
	} else {
		fmt.Println("El inventario no tiene datos de conexion, favor de verificar.", idAlmacen, datosAlmacen.Nombre)
	}
	return datosConexion
}

//########################< FUNCIONES GENERALES PSQL >#############################

//CreaTablasPostgres  funcion para crear las tablas en postgres de almacen
func CreaTablasPostgres(ID string, IDConexion bson.ObjectId) string {
	//Conexion a postgreSQL
	Conexion := ConexionModel.GetOne(IDConexion)
	var parametros MoConexion.ParametrosConexionPostgres
	parametros.Servidor = Conexion.Servidor
	parametros.Usuario = Conexion.UsuarioBD
	parametros.Pass = Conexion.PassBD
	parametros.NombreBase = Conexion.NombreBD
	BasePosGres, err := MoConexion.ConexioServidorAlmacen(parametros)

	query := "CREATE TABLE public.\"Kardex_" + ID + "\"(" +
		"\"IdOperacion\" character(25) NOT NULL," +
		"\"IdMovimiento\" character(25) NOT NULL," +
		"\"IdProducto\" character(25) NOT NULL," +
		"\"Cantidad\" numeric," +
		"\"Costo\" numeric," +
		"\"Precio\" numeric," +
		"\"ImpuestoTotal\" numeric," +
		"\"DescuentoTotal\" numeric," +
		"\"TipoOperacion\" character(25) NOT NULL," +
		"\"Existencia\" numeric," +
		"\"FechaHora\" timestamp without time zone," +
		"CONSTRAINT \"LLaveUunicaK_" + ID + "\" PRIMARY KEY (\"IdOperacion\", \"IdMovimiento\", \"IdProducto\"  ))"
	con, err := BasePosGres.Query(query)
	if err != nil {
		fmt.Println("tiene un error", err)
		check(err, "Error al crear kardex con Postgres")
	}
	query = "CREATE TABLE public.\"Inventario_" + ID + "\"(" +
		"\"IdProducto\" character(25) NOT NULL," +
		"\"Existencia\" numeric," +
		"\"Costo\" numeric," +
		"\"Precio\" numeric," +
		"\"Estatus\" character varying," +
		"CONSTRAINT \"LlaveUnicaI_" + ID + "\" PRIMARY KEY (\"IdProducto\"))"

	con, err = BasePosGres.Query(query)
	if err != nil {
		fmt.Println("tiene un error", err)
		check(err, "Error al crear inventario con Postgres")
	}

	query = "CREATE TABLE public.\"ImpuestoC_" + ID + "\"" +
		"(" +
		"\"IdMovimiento\" character(25) NOT NULL," +
		"\"IdProducto\" character(25) NOT NULL," +
		"\"TipoImpuesto\" character(25)," +
		"\"Factor\" character(25)," +
		"\"Tratamiento\" character(25)," +
		"\"Valor\" numeric," +
		"CONSTRAINT \"LLavePrimariaC_" + ID + "\" PRIMARY KEY (\"IdMovimiento\", \"IdProducto\", \"TipoImpuesto\", \"Factor\", \"Tratamiento\", \"Valor\" ))"

	con, err = BasePosGres.Query(query)
	if err != nil {
		fmt.Println("tiene un error", err)
		check(err, "Error al crear impuesto con Postgres")
	}

	query = "CREATE TABLE public.\"ImpuestoV_" + ID + "\"" +
		"(" +
		"\"IdMovimiento\" character(25) NOT NULL," +
		"\"IdProducto\" character(25) NOT NULL," +
		"\"TipoImpuesto\" character(25)," +
		"\"Factor\" character(25)," +
		"\"Tratamiento\" character(25)," +
		"\"Valor\" numeric," +
		"CONSTRAINT \"LLavePrimariaV_" + ID + "\" PRIMARY KEY (\"IdMovimiento\", \"IdProducto\", \"TipoImpuesto\", \"Factor\", \"Tratamiento\", \"Valor\" ))"

	con, err = BasePosGres.Query(query)
	if err != nil {
		fmt.Println("tiene un error", err)
		check(err, "Error al crear impuesto con Postgres")
	}

	query = "CREATE TABLE public.\"Descuento_" + ID + "\"" +
		"(" +
		"\"IdMovimiento\" character(25) NOT NULL," +
		"\"IdProducto\" character(25) NOT NULL," +
		"\"IdDescuento\" character(25)," +
		"\"Valor\" numeric," +
		"CONSTRAINT \"LLavePrimariaDescuento_" + ID + "\" PRIMARY KEY (\"IdMovimiento\", \"IdProducto\", \"IdDescuento\")" +
		")"
	con, err = BasePosGres.Query(query)
	if err != nil {
		fmt.Println("tiene un error", err)
		check(err, "Error al crear Descuento con Postgres")
	}
	con.Close()
	BasePosGres.Close()
	return "1"
}

func check(err error, cadena string) {
	if err != nil {
		panic(err)
	}
}

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("Nombre")
	queryQuotes = queryQuotes.Field("Nombre")

	queryTilde = queryTilde.Field("Tipo")
	queryQuotes = queryQuotes.Field("Tipo")

	queryTilde = queryTilde.Field("Clasificacion")
	queryQuotes = queryQuotes.Field("Clasificacion")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoAlmacen, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoAlmacen, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}

// IndexOfAlmacen regresa el indice del elemento TestingObject en los ids de almacen seleccionados SelectedAlmacens
func IndexOfAlmacen(SelectedAlmacens []bson.ObjectId, TestingObject bson.ObjectId) int {
	for i, selected := range SelectedAlmacens {
		if selected == TestingObject {
			return i
		}
	}
	return -1
}
