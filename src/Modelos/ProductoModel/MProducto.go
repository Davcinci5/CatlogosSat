package ProductoModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modelos/CatalogoModel"
	"../../Modelos/UnidadModel"

	"../../Modulos/Conexiones"
	"../../Modulos/ConsultasSql"
	"../../Modulos/General"
	"../../Modulos/Imagenes"
	"../../Modulos/Variables"

	"github.com/leekchan/accounting"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//ProductoMgo estructura de Productos mongo
type ProductoMgo struct {
	ID            bson.ObjectId   `bson:"_id,omitempty"`
	Nombre        string          `bson:"Nombre"`
	Codigos       CodigoMgo       `bson:"Codigos"`
	Tipo          bson.ObjectId   `bson:"Tipo"`
	Imagenes      []bson.ObjectId `bson:"Imagenes,omitempty"`
	Unidad        bson.ObjectId   `bson:"Unidad,omitempty"`
	Mmv           float64         `bson:"Mmv"`
	VentaFraccion bool            `bson:"VentaFraccion"`
	Etiquetas     []string        `bson:"Etiquetas"`
	Almacenes     []bson.ObjectId `bson:"Almacenes,omitempty"`
	Estatus       bson.ObjectId   `bson:"Estatus"`
	FechaHora     time.Time       `bson:"FechaHora"`
}

//ProductoElastic estructura de Productos para insertar en Elastic
type ProductoElastic struct {
	Nombre        string    `json:"Nombre"`
	Codigos       CodigoMgo `json:"Codigos"`
	Tipo          string    `json:"Tipo"`
	Unidad        string    `json:"Unidad"`
	VentaFraccion bool      `json:"VentaFraccion"`
	Etiquetas     []string  `json:"Etiquetas"`
	Estatus       string    `json:"Estatus"`
	FechaHora     time.Time `json:"FechaHora"`
}

//CodigoMgo subestructura de Producto
type CodigoMgo struct {
	Claves  []string `bson:"Claves"`
	Valores []string `bson:"Valores"`
}

//ProductoPgrs estructrura del Producto para postgres
type ProductoPgrs struct {
	IDProducto string
	Existencia float64
	Estatus    string
	Costo      float64
	Precio     float64
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []ProductoMgo {
	var result []ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	err = Productos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Productos.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) ProductoMgo {
	var result ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	err = Productos.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []ProductoMgo {
	var result []ProductoMgo
	var aux ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = ProductoMgo{}
		Productos.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) ProductoMgo {
	var result ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)

	if err != nil {
		fmt.Println(err)
	}
	err = Productos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	err = Productos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboProductos regresa un combo de Producto de mongo
func CargaComboProductos(ID string) string {
	Productos := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Productos {
		if ID == v.ID.Hex() {
			templ += `<option value=" ` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value=" ` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//########## GET NAMES ####################################

//GetNameProducto regresa el nombre del Producto con el ID especificado
func GetNameProducto(ID bson.ObjectId) string {
	var result ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	Productos.Find(bson.M{"_id": ID}).One(&result)

	s.Close()
	return result.Nombre
}

//GetProductoMongo Regresa un Producto from valor de  codigo
func GetProductoMongo(codigoBarra string) ProductoMgo {
	var result ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	defer s.Close()
	if err != nil {
		fmt.Println(err)
	}
	// Productos.Find(bson.M{"Codigos.Valores": codigoBarra}).One(&result)
	Productos.Find(bson.M{"Codigos.Valores": codigoBarra}).One(&result)

	return result
}

//GetProductoPostgres Regresa la existencia del producto en los almacenes condicionados por estado
func GetProductoPostgres(id, origen, destino string) (ProductoPgrs, ProductoPgrs) {
	datosConexion := AlmacenModel.ObtenerParametrosConexion(bson.ObjectIdHex(origen))

	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	conexOrigen, err1 := MoConexion.ConexionEspecificaPsql(paramConex)
	if err1 != nil {
		fmt.Println("Error en la conexion origen:", err1)
	}

	datosConexion = AlmacenModel.ObtenerParametrosConexion(bson.ObjectIdHex(destino))

	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	conexDestino, err2 := MoConexion.ConexionEspecificaPsql(paramConex)
	if err2 != nil {
		fmt.Println("Error en la conexion destino:", err2)
	}
	defer conexOrigen.Close()
	defer conexDestino.Close()

	var ProductoOrigen ProductoPgrs
	var ProductoDestino ProductoPgrs

	queryOrigen := `SELECT "IdProducto","Existencia","Estatus","Costo","Precio" FROM public."Inventario_` + origen + `" WHERE "IdProducto"='` + id + `'`
	queryDestino := `SELECT "IdProducto","Existencia","Estatus","Costo","Precio" FROM public."Inventario_` + destino + `" WHERE "IdProducto"='` + id + `'`
	resultado, _ := conexOrigen.Query(queryOrigen)
	for resultado.Next() {
		err := resultado.Scan(&ProductoOrigen.IDProducto, &ProductoOrigen.Existencia, &ProductoOrigen.Estatus, &ProductoOrigen.Costo, &ProductoOrigen.Precio)
		if err != nil {
			fmt.Println(err)
		}
	}

	resultado2, err2 := conexDestino.Query(queryDestino)
	if err2 != nil {
		fmt.Println("Errores en la consulta del almacen destino: ", err2)
	}
	for resultado2.Next() {
		err := resultado2.Scan(&ProductoDestino.IDProducto, &ProductoDestino.Existencia, &ProductoDestino.Estatus, &ProductoDestino.Costo, &ProductoDestino.Precio)
		if err != nil {
			fmt.Println(err)
		}
	}
	return ProductoOrigen, ProductoDestino
}

//RealizarTraslado modifica inventario de los almacenes
func RealizarTraslado(AlmacenOrigen, AlmacenDestino, IDProducto string, CantidadOrigen, CantidadDestino float64) (bool, string) {
	datosConexion := AlmacenModel.ObtenerParametrosConexion(bson.ObjectIdHex(AlmacenOrigen))

	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	datosConexion = AlmacenModel.ObtenerParametrosConexion(bson.ObjectIdHex(AlmacenDestino))

	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	conexOrigen, err1 := MoConexion.ConexionEspecificaPsql(paramConex)
	conexDestino, err2 := MoConexion.ConexionEspecificaPsql(paramConex)
	if err1 != nil && err2 != nil {
		fmt.Println(err1, err2)
	}
	defer conexOrigen.Close()
	defer conexDestino.Close()

	var ProductoDestino ProductoPgrs
	var OpExitosa = true
	var Msj string

	CantidadOrigenS := strconv.FormatFloat(CantidadOrigen, 'f', -1, 64)
	CantidadDestinoS := strconv.FormatFloat(CantidadDestino, 'f', -1, 64)

	//verificar si el producto id_prod existe en el almacen destino, si no existe crearlo
	// como verificamos eso
	QueryExistenciaDelProducto := `SELECT "IdProducto" FROM public."Inventario_` + AlmacenDestino + `" WHERE "IdProducto"='` + IDProducto + `'`
	res, err := conexDestino.Query(QueryExistenciaDelProducto)
	if err != nil {
		fmt.Println("Error al realizar la consulta: ", err)
	}

	for res.Next() {
		err := res.Scan(&ProductoDestino.IDProducto)
		if err != nil {
			fmt.Println(err)
		}
	}

	var QueryModificaAlmacenOrigen string
	var QueryModificaAlmacenDestino string

	if ProductoDestino.IDProducto != "" {
		QueryModificaAlmacenOrigen = `UPDATE public."Inventario_` + AlmacenOrigen + `" SET "Existencia"='` + CantidadOrigenS + `' WHERE "IdProducto"='` + IDProducto + `'`
		QueryModificaAlmacenDestino = `UPDATE public."Inventario_` + AlmacenDestino + `" SET "Existencia"='` + CantidadDestinoS + `' WHERE "IdProducto"='` + IDProducto + `'`
		_, err := conexOrigen.Exec(QueryModificaAlmacenOrigen)
		if err != nil {
			OpExitosa = false
			Msj += "Error al actualizar el almacen origen (1)"
			fmt.Println(err)
		}
		_, err = conexDestino.Exec(QueryModificaAlmacenDestino)
		if err != nil {
			OpExitosa = false
			Msj += "Error al actualizar el almacen destino"
			fmt.Println(err)
		}

	} else {
		QueryModificaAlmacenOrigen = `UPDATE public."Inventario_` + AlmacenOrigen + `" SET "Existencia"='` + CantidadOrigenS + `' WHERE "IdProducto"='` + IDProducto + `'`
		//QueryModificaAlmacenDestino = `INSERT INTO public."Inventario_` + AlmacenDestino + `" (IdProducto,Existencia,Estatus,Costo,Precio) VALUES ('` + IDProducto + `','` + CantidadDestinoS + `','ACTIVO',`+ProductoDestino.Costo+`,4)`
		QueryModificaAlmacenDestino := fmt.Sprintf(`INSERT INTO  public."Inventario_%v" ("IdProducto","Existencia","Estatus","Costo","Precio") values('%v','%v','%v','%v','%v')`, AlmacenDestino, IDProducto, CantidadDestinoS, "ACTIVO", ProductoDestino.Costo, ProductoDestino.Precio)
		_, err := conexOrigen.Exec(QueryModificaAlmacenOrigen)
		if err != nil {
			OpExitosa = false
			Msj += "Error al actualizar el almacen origen (2)"
			fmt.Println(err)
		}
		_, err = conexDestino.Exec(QueryModificaAlmacenDestino)
		if err != nil {
			OpExitosa = false
			Msj += "Error al insertar en el almacen destino"
			fmt.Println(err)
		}
	}
	return OpExitosa, Msj
}

//RealizarAjuste afecta los inventarios de los almacenes
func RealizarAjuste(AlmacenOrigen, IDProducto string, CantidadOrigen float64) (bool, string) {
	datosConexion := AlmacenModel.ObtenerParametrosConexion(bson.ObjectIdHex(AlmacenOrigen))

	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	conex, err := MoConexion.ConexionEspecificaPsql(paramConex)
	if err != nil {
		fmt.Println(err)
	}
	var OpExitosa = true
	var Msj string
	defer conex.Close()

	CantidadOrigenS := strconv.FormatFloat(CantidadOrigen, 'f', -1, 64)
	query := `UPDATE public."Inventario_` + AlmacenOrigen + `" SET "Existencia"='` + CantidadOrigenS + `' WHERE "IdProducto"='` + IDProducto + `'`
	_, err = conex.Exec(query)
	if err != nil {
		OpExitosa = false
		Msj += "Error Alm Origen"
		fmt.Println(err)
	}
	return OpExitosa, Msj
}

//BusquedasProductosNombre obtiene un producto de acuerdo a un nombre especificado
func BusquedasProductosNombre(nombre string) ProductoMgo {
	var result ProductoMgo
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)

	if err != nil {
		fmt.Println(err)
	}
	err = Productos.Find(bson.M{"Nombre": nombre}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Productos []ProductoMgo) (string, string) {
	etiquetas := ``
	codigos := ``
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>			
					<th>Nombre</th>									
					<th>Códigos</th>									
					<th>Tipo</th>									
					<th>Unidad</th>									
					<th>Venta por Fracciones</th>									
					<th>Etiquetas</th>									
					<th>Estatus</th>									
					<th>FechaHora</th>					
				</tr>`

	for k, v := range Productos {
		etiquetas = ``
		codigos = ``
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Productos/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`

		for i, cla := range v.Codigos.Claves {
			codigos += cla + `:` + v.Codigos.Valores[i] + `,`
		}

		codigos = codigos[:len(codigos)-1]
		cuerpo += `<td>` + codigos + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo, 162) + `</td>`
		cuerpo += `<td>` + UnidadModel.GetNombreUnidadByField("Datos._id", v.Unidad) + `</td>`
		if v.VentaFraccion {
			cuerpo += `<td>SI</td>`
		} else {
			cuerpo += `<td>NO</td>`
		}

		for _, val := range v.Etiquetas {
			etiquetas += val + `,`
		}
		etiquetas = etiquetas[:len(etiquetas)-1]
		cuerpo += `<td>` + etiquetas + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 161) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//GeneraTablaConImpuestos genera una tabla utilizada en el modulo de ventas
func GeneraTablaConImpuestos(Productos []ProductoMgo, Almacen string) string {
	floats := accounting.Accounting{Symbol: "", Precision: 2}
	tabla := ``
	for _, producto := range Productos {
		if bson.IsObjectIdHex(producto.ID.Hex()) {

			htmlProducto := `<tr id="filaProducto` + producto.ID.Hex() + `">`
			htmlProducto += `<td><label>` + producto.Nombre + `</td>`
			htmlProducto += `<td>`
			for _, value := range producto.Codigos.Valores {
				htmlProducto += `<label>` + value + `</label>`
			}
			htmlProducto += `</td>`
			tipo := CatalogoModel.RegresaNombreSubCatalogo(producto.Tipo)
			if tipo == "" {
				tipo = "Tipo no encontrado"
			}
			htmlProducto += `<td><label>` + tipo + `</label></td>`
			unidad := UnidadModel.RegresaNombreUnidad(producto.Unidad)
			if unidad == "" {
				unidad = "Unidad no encontrado"
			}

			htmlProducto += `<td><label>` + unidad + `</label></td>`
			estatus := CatalogoModel.RegresaNombreSubCatalogo(producto.Estatus)
			if estatus == "" {
				estatus = "Estatus desconocido"
			}
			htmlProducto += `<td><label>` + estatus + `</label></td>`
			htmlProducto += `<td>`
			for _, value := range producto.Etiquetas {
				htmlProducto += `<label>` + value + `</label>`
			}
			htmlProducto += `</td>`
			htmlImagenes, _ := Imagenes.CargaTemplateImagenesCompras(producto.Imagenes)

			htmlProducto += `<td>` + htmlImagenes + `</td>`

			var tipos, factores, tratamientos, valores string
			Impuestos, err := ConsultasSql.ConsultaValorImpuestoEnAlmacen(Almacen, producto.ID.Hex())
			if err != nil {
				fmt.Println(err)
			} else {
				// init := `[`
				// fin := `]`
				for _, v := range Impuestos {
					tipos += v.Tipo + `,`
					factores += v.Factor + `,`
					tratamientos += v.Tratamiento + `,`
					valores += floats.FormatMoney(v.Valor) + `,`

					// html += `<tr>
					// 			<td><input type="hidden" class="form-control" name="TipoImpuestoLista" value="` + v.Tipo + `"><input type="text" class="form-control" value="` + v.Tipo + `" readonly></td>
					// 			<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="ValorImpuestoLista" value="` + floats.FormatMoney(v.Valor) + `" readonly></td>
					// 			<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="FactorImpuestoLista" value="` + v.Factor + `" readonly></td>
					// 			<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="TratamientoImpuestoLista" value="` + v.Tratamiento + `" readonly></td>
					// 			<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>
					// 		</tr>`
				}
			}

			if len(tipos) > 0 {
				tipos = tipos[0 : len(tipos)-1]
			}
			if len(factores) > 0 {
				factores = factores[0 : len(factores)-1]
			}
			if len(tratamientos) > 0 {
				tratamientos = tratamientos[0 : len(tratamientos)-1]
			}
			if len(valores) > 0 {
				valores = valores[0 : len(valores)-1]
			}

			htmlProducto += `<td>
								<button id="` + producto.ID.Hex() + `" tipos="` + tipos + `" factores="` + factores + `" tratamientos="` + tratamientos + `" valores="` + valores + `" type="button" class="btn btn-info" onClick='leerFilaProducto(this)'>`

			htmlProducto += `		<span class="glyphicon glyphicon-ok-sign">Seleccionar</span>`

			htmlProducto += `	</button>
							 </td>`

			htmlProducto += `</tr>`

			tabla += htmlProducto
		}
	}
	return tabla
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

	queryTilde = queryTilde.Field("Codigos.Valores")
	queryQuotes = queryQuotes.Field("Codigos.Valores")

	queryTilde = queryTilde.Field("Etiquetas")
	queryQuotes = queryQuotes.Field("Etiquetas")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoProducto, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoProducto, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}

//BusquedaElastic realiza una busqueda en elastic, devuelve un error en caso de no conectarse o no encontrar informacion
//Funcion realizada por Ramon, para devolver informacion y un error
func BusquedaElastic(texto string) (*elastic.SearchResult, error) {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("Nombre")
	queryQuotes = queryQuotes.Field("Nombre")

	queryTilde = queryTilde.Field("Codigos.Valores")
	queryQuotes = queryQuotes.Field("Codigos.Valores")

	queryTilde = queryTilde.Field("Etiquetas")
	queryQuotes = queryQuotes.Field("Etiquetas")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	//var docs *elastic.SearchResult
	//var err error

	docs, err := MoConexion.BusquedaElastic(MoVar.TipoProducto, queryTilde)
	if err != nil {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento", err)
		return docs, err
	} else {
		if docs.Hits.TotalHits == 0 {
			docs, err := MoConexion.BusquedaElastic(MoVar.TipoProducto, queryQuotes)
			if err != nil {
				fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
				return docs, err
			}
		}
	}
	return docs, err
}
