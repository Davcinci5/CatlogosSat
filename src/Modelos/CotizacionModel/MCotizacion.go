package CotizacionModel

import (
	"fmt"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modelos/ProductoModel"
	"../../Modelos/UnidadModel"
	"../../Modulos/ConsultasSql"
	"../../Modulos/Imagenes"
	"github.com/leekchan/accounting"

	"../../Modelos/CatalogoModel"
	"../../Modelos/ClienteModel"
	"../../Modelos/PersonaModel"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//CotizacionMgo estructura de Cotizacions mongo
type CotizacionMgo struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Operacion    bson.ObjectId `bson:"Operacion,omitempty"`
	Usuario      bson.ObjectId `bson:"Usuario,omitempty"`
	Equipo       bson.ObjectId `bson:"Equipo,omitempty"`
	Cliente      bson.ObjectId `bson:"Cliente,omitempty"`
	Grupo        bson.ObjectId `bson:"Grupo,omitempty"`
	Nombre       string        `bson:"Nombre"`
	Telefono     string        `bson:"Telefono"`
	Correo       string        `bson:"Correo"`
	FormaDePago  bson.ObjectId `bson:"FormaDePago,omitempty"`
	FormaDeEnvío bson.ObjectId `bson:"FormaDeEnvío,omitempty"`
	Buscar       string        `bson:"Buscar"`
	Lista        string        `bson:"Lista"`
	Carrito      string        `bson:"Carrito"`
	Resumen      string        `bson:"Resumen"`
	Estatus      bson.ObjectId `bson:"Estatus,omitempty"`
	FechaInicio  time.Time     `bson:"FechaInicio"`
	FechaFin     time.Time     `bson:"FechaFin"`
}

//CotizacionElastic estructura de Cotizacions para insertar en Elastic
type CotizacionElastic struct {
	Operacion    string    `json:"Operacion,omitempty"`
	Usuario      string    `json:"Usuario,omitempty"`
	Equipo       string    `json:"Equipo,omitempty"`
	Cliente      string    `json:"Cliente,omitempty"`
	Grupo        string    `json:"Grupo,omitempty"`
	Nombre       string    `json:"Nombre"`
	Telefono     string    `json:"Telefono"`
	Correo       string    `json:"Correo"`
	FormaDePago  string    `json:"FormaDePago,omitempty"`
	FormaDeEnvío string    `json:"FormaDeEnvío,omitempty"`
	Buscar       string    `json:"Buscar"`
	Lista        string    `json:"Lista"`
	Carrito      string    `json:"Carrito"`
	Resumen      string    `json:"Resumen"`
	Estatus      string    `json:"Estatus,omitempty"`
	FechaInicio  time.Time `json:"FechaInicio"`
	FechaFin     time.Time `json:"FechaFin"`
}

//DatosVentaTemporal para consultar carrito
type DatosVentaTemporal struct {
	Operacion   string
	Movimiento  string
	Almacen     string
	Producto    string
	Cantidad    float64
	Costo       float64
	Precio      float64
	IDimpuesto  string
	Impuesto    float64
	IDdescuento string
	Descuento   float64
	Existencia  float64
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []CotizacionMgo {
	var result []CotizacionMgo
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Cotizacions.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Cotizacions.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) CotizacionMgo {
	var result CotizacionMgo
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Cotizacions.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []CotizacionMgo {
	var result []CotizacionMgo
	var aux CotizacionMgo
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = CotizacionMgo{}
		Cotizacions.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) CotizacionMgo {
	var result CotizacionMgo
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)

	if err != nil {
		fmt.Println(err)
	}
	err = Cotizacions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result CotizacionMgo
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}
	err = Cotizacions.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboCotizacions regresa un combo de Cotizacion de mongo
func CargaComboCotizacions(ID string) string {
	Cotizacions := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Cotizacions {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Cotizacions []CotizacionMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>			
					<th>Operacion</th>									
					<th>Usuario</th>									
					<th>Cliente</th>									
					<th>Grupo</th>									
					<th>Nombre</th>									
					<th>Telefono</th>									
					<th>Correo</th>									
					<th>Estatus</th>									
					<th>FechaInicio</th>									
					<th>FechaFin</th>					
				</tr>`

	for k, v := range Cotizacions {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Cotizacions/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Operacion.Hex() + `</td>`
		cuerpo += `<td>` + v.Usuario.Hex() + `</td>`
		cuerpo += `<td>` + v.Cliente.Hex() + `</td>`
		cuerpo += `<td>` + v.Grupo.Hex() + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + v.Telefono + `</td>`
		cuerpo += `<td>` + v.Correo + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 171) + `</td>`
		cuerpo += `<td>` + v.FechaInicio.Format(time.RFC1123) + `</td>`
		cuerpo += `<td>` + v.FechaFin.Format(time.RFC1123) + `</td>`
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

	queryTilde = queryTilde.Field("Operacion")
	queryQuotes = queryQuotes.Field("Operacion")

	queryTilde = queryTilde.Field("Usuario")
	queryQuotes = queryQuotes.Field("Usuario")

	queryTilde = queryTilde.Field("Equipo")
	queryQuotes = queryQuotes.Field("Equipo")

	queryTilde = queryTilde.Field("Cliente")
	queryQuotes = queryQuotes.Field("Cliente")

	queryTilde = queryTilde.Field("Grupo")
	queryQuotes = queryQuotes.Field("Grupo")

	queryTilde = queryTilde.Field("Nombre")
	queryQuotes = queryQuotes.Field("Nombre")

	queryTilde = queryTilde.Field("Telefono")
	queryQuotes = queryQuotes.Field("Telefono")

	queryTilde = queryTilde.Field("Correo")
	queryQuotes = queryQuotes.Field("Correo")

	queryTilde = queryTilde.Field("FormaDePago")
	queryQuotes = queryQuotes.Field("FormaDePago")

	queryTilde = queryTilde.Field("FormaDeEnvío")
	queryQuotes = queryQuotes.Field("FormaDeEnvío")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	queryTilde = queryTilde.Field("FechaInicio")
	queryQuotes = queryQuotes.Field("FechaInicio")

	queryTilde = queryTilde.Field("FechaFin")
	queryQuotes = queryQuotes.Field("FechaFin")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoCotizacion, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoCotizacion, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}

// GeneraTemplateBusquedaClientes respuesta AJAX genera el template de la estructura de clientes para la tabla de busqueda
func GeneraTemplateBusquedaClientes(Personas []PersonaModel.PersonaMgo) string {

	html := ``
	for in, Per := range Personas {
		Cliente := ClienteModel.GetEspecificByFields("IDPersona", Per.ID)
		if !MoGeneral.EstaVacio(Cliente) {
			html += `
			<tr id="` + Cliente.ID.Hex() + `">
				<th scope="row">` + strconv.Itoa(in) + `</th>
				<td>` + Per.Nombre + `</td>
				<td>Av. Siempre Viva</td>
				<td>` + Cliente.MediosDeContacto.Correos.Principal + `</td>
				<td>` + Cliente.MediosDeContacto.Telefonos.Principal + `</td>
				`
			TiposEstatus := CatalogoModel.RegresaValoresCatalogosClave(MoVar.CatalogoDeEstatusDeCliente)
			for _, est := range TiposEstatus.Valores {
				if est.ID.Hex() == Cliente.Estatus.Hex() {
					html += `<td>` + est.Valor + `</td>`
				}
			}
			html += `				
				<td>
				<button type="button" class="btn btn-default" onclick="SeleccionarCliente(this);">
					<span class="glyphicon glyphicon-ok"></span>
				</button>
				</td>
			</tr>`
		}

	}

	return html
}

//GeneraTemplateBusquedaClienteEspecifico genera el templade para la tabla cliente
func GeneraTemplateBusquedaClienteEspecifico(Cliente ClienteModel.ClienteMgo) string {
	Persona := PersonaModel.GetOne(Cliente.IDPersona)
	html := `
	<tr id="` + Cliente.ID.Hex() + `" Nombre="` + Persona.Nombre + `" Telefono="` + Cliente.MediosDeContacto.Telefonos.Principal + `" Correo="` + Cliente.MediosDeContacto.Correos.Principal + `">
		<th scope="row">#</th>
		<td>` + Persona.Nombre + `</td>
		<td>Av. Siempre Viva</td>
		<td>` + Cliente.MediosDeContacto.Correos.Principal + `</td>
		<td>` + Cliente.MediosDeContacto.Telefonos.Principal + `</td>`

	TiposEstatus := CatalogoModel.RegresaValoresCatalogosClave(MoVar.CatalogoDeEstatusDeCliente)
	for _, est := range TiposEstatus.Valores {
		if est.ID.Hex() == Cliente.Estatus.Hex() {
			html += `<td>` + est.Valor + `</td>`
		}
	}

	html += `				
		<td>
			<button type="button" class="btn btn-default deleteButton">
				<span class="glyphicon glyphicon-trash"></span>
			</button>
		</td>
		<input name="ClienteId" id="ClienteId" type="hidden" value="` + Cliente.ID.Hex() + `" />
	</tr>`
	return html
}

// GeneraTemplateBusquedaDeProductos Genera el template de Busqueda de los productos encontrados por elastic
func GeneraTemplateBusquedaDeProductos(Productos []ProductoModel.ProductoMgo) string {
	money := accounting.Accounting{Symbol: "$", Precision: 2}
	busqueda := ``
	aux := ``
	disponibilidad := ``

	for _, v := range Productos {

		for _, b := range v.Almacenes {

			estado, Cantidad, _, Precio, err := ConsultasSql.ConsultaDatosDeProductoActivo(b.Hex(), v.ID.Hex())
			if err != nil {
				fmt.Println(err)
			}

			if estado {
				if Cantidad > 0 {

					almacen := AlmacenModel.GetOne(b)
					busqueda += `<tr id="` + v.ID.Hex() + `" almacen="` + b.Hex() + `">`
					htmlImagenes, _ := Imagenes.CargaTemplateImagenesVentas(v.Imagenes)
					busqueda += `<td>` + htmlImagenes + `</td>`
					aux = ``
					for _, valores := range v.Codigos.Valores {
						aux += valores + `,`
					}

					busqueda += `<td>` + aux[:len(aux)-1] + `</td>`
					busqueda += `<td>` + v.Nombre + `</td>`

					unidad := UnidadModel.GetSubEspecificByFields("Datos._id", v.Unidad)
					busqueda += `<td>` + unidad.Abreviatura + `</td>`
					busqueda += `<td>` + money.FormatMoney(Precio) + `</td>`

					if v.VentaFraccion {
						disponibilidad = strconv.FormatFloat(Cantidad, 'f', 3, 64)
					} else {
						disponibilidad = strconv.FormatFloat(Cantidad, 'f', 1, 64)
					}

					busqueda += `<td>` + almacen.Nombre + `</td>`
					busqueda += `<td>` + disponibilidad + `</td>`

					if v.VentaFraccion {
						busqueda += `<td id="thisone"><input name="Cantidad" id = "` + v.ID.Hex() + `" almacen="` + b.Hex() + `" fraccion="` + strconv.FormatBool(v.VentaFraccion) + `" existencia="` + strconv.FormatFloat(Cantidad, 'f', 3, 64) + `" type="number" class="form-control" value=""></td>`
					} else {
						busqueda += `<td id="thisone"><input name="Cantidad" id = "` + v.ID.Hex() + `" almacen="` + b.Hex() + `" fraccion="` + strconv.FormatBool(v.VentaFraccion) + `" existencia="` + strconv.FormatFloat(Cantidad, 'f', 1, 64) + `" type="number" class="form-control" value=""></td>`
					}

					busqueda += `				
						<td>
						<button type="button" class="btn btn-primary" onclick="SeleccionarProducto(this);">
							<span class="glyphicon glyphicon-ok"></span>
						</button>
						</td>`

					// busqueda += `<td><input type="checkbox" name="ModalSeleccione" value="` + v.ID.Hex() + `"></td>`
					busqueda += `</tr>`
					busqueda += `<script>`
					busqueda += `$("input[name='Cantidad']").keydown(function(e) {
									if (e.which == 13 || e.keyCode == 13) {
										e.preventDefault();
										SeleccionarProducto(this);
									}
								});`
					busqueda += `</script>`
				}
			}

		}

	}

	return busqueda
}

//ConsultaDatosOperacion genera un template html de los datos agregados al carrito
func ConsultaDatosOperacion(idOperacion string) (string, string) {
	html := ``
	calculadora := ``

	var subtotales []float64 //para almacenar los totales por poducto
	var totales []float64    //para almacenar los totales por poducto
	var impuestos []float64  //para almacenar los impuestos por producto

	money := accounting.Accounting{Symbol: "$", Precision: 2}
	floats := accounting.Accounting{Symbol: "", Precision: 2}
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		fmt.Println(err)
		return "", "Error al iniciar Sesion"
	}

	Query := fmt.Sprintf(`SELECT "Movimiento", "Producto", "Almacen", "Cantidad", "Costo", "Precio", "Impuesto","Existencia" FROM public."Cotizacion" WHERE "Operacion" = '%v' ORDER BY "FechaHora" ASC`, idOperacion)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	productos := []DatosVentaTemporal{}
	for resultSet.Next() {
		producto := DatosVentaTemporal{}
		resultSet.Scan(&producto.Movimiento, &producto.Producto, &producto.Almacen, &producto.Cantidad, &producto.Costo, &producto.Precio, &producto.Impuesto, &producto.Existencia)
		productos = append(productos, producto)
	}

	for _, productoPostgres := range productos {

		Impuestos, err := ConsultasSql.ConsultaValorImpuestoEnAlmacen(productoPostgres.Almacen, productoPostgres.Producto)
		if err != nil {
			fmt.Println(err)
			return "", ""
		}

		producto := ProductoModel.GetOne(bson.ObjectIdHex(productoPostgres.Producto))
		almacen := AlmacenModel.GetOne(bson.ObjectIdHex(productoPostgres.Almacen))

		html += `<tr id = "` + productoPostgres.Movimiento + `">`
		html += `<td>`
		aux := ``
		for _, valores := range producto.Codigos.Valores {
			aux += valores + `,`
		}
		html += aux[:len(aux)-1]
		html += `</td>`
		html += `<td>` + producto.Nombre + `</td>`
		html += `<td>` + almacen.Nombre + `</td>`
		html += `<td>` + money.FormatMoney(productoPostgres.Precio) + `</td>`

		if producto.VentaFraccion {
			html += `<td><input name="Producto" id = "` + producto.ID.Hex() + `" almacen="` + productoPostgres.Almacen + `" fraccion="` + strconv.FormatBool(producto.VentaFraccion) + `" existencia="` + strconv.FormatFloat(productoPostgres.Existencia, 'f', 3, 64) + `" previo="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 3, 64) + `" onchange="ValidaNumeroCarrito(this)" onfocusout="AplicaPeticion(this)" type="number" class="form-control" value="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 3, 64) + `"></td>`
		} else {
			html += `<td><input name="Producto" id = "` + producto.ID.Hex() + `" almacen="` + productoPostgres.Almacen + `" fraccion="` + strconv.FormatBool(producto.VentaFraccion) + `" existencia="` + strconv.FormatFloat(productoPostgres.Existencia, 'f', 1, 64) + `" previo="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 1, 64) + `" onchange="ValidaNumeroCarrito(this)" onfocusout="AplicaPeticion(this)" type="number" class="form-control" value="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 1, 64) + `"></td>`
		}

		unidad := UnidadModel.GetSubEspecificByFields("Datos._id", producto.Unidad)
		html += `<td>` + unidad.Abreviatura + `</td>`

		precio := productoPostgres.Cantidad * productoPostgres.Precio
		html += `<td>` + money.FormatMoney(precio) + `</td>`

		valimpuesto := ConsultasSql.CalculaTotalDeImpuestoPorProducto(Impuestos, productoPostgres.Precio)
		html += `<td>
					<div class="popover-markup">
					<a class="Impuestolink trigger">` + money.FormatMoney(productoPostgres.Cantidad*valimpuesto) + `</a>
						<div class="head hide">Edita Impuesto(s)</div>
						<div class="content hide">
							<div class="editableform-loading" style="display: none;">
								aqui va info al abrir
							</div>							
							<div class="form-inline">`

		for _, v := range Impuestos {

			html += `	
     							<div class="editable-label form-group">
									<label class="label-sm">` + v.Tipo + `(` + v.Factor + `)</label>
								</div>
								<div class="editable-input form-group" style="position: relative;">
									<input type="number" name="ValorImpuestoEdita" class="form-control input-sm" value="` + floats.FormatMoney(v.Valor) + `" style="padding-right: 24px;">											
									<button Precio= "` + floats.FormatMoney(productoPostgres.Precio) + `" Movimiento= "` + productoPostgres.Movimiento + `" Producto="` + productoPostgres.Producto + `" Tipo="` + v.Tipo + `"  Factor="` + v.Factor + `" Valor="` + floats.FormatMoney(v.Valor) + `" Almacen="` + productoPostgres.Almacen + `" type="button" onclick="ActualizaImpuestoVenta(this)" class="btn btn-primary btn-sm editable-submit" data-loading-text="<i class='fa fa-circle-o-notch fa-spin'></i>">
										<i class="glyphicon glyphicon-ok"></i>
									</button>
								</div>	
								<hr>`
		}

		html += `	</div>	
								<div class="editable-error-block help-block" style="display: none;">
									Aqui van errores
								</div>

						</div>
					</div>
				</td>`

		html += `<td>` + money.FormatMoney(productoPostgres.Cantidad*(productoPostgres.Precio+valimpuesto)) + `</td>`
		html += `<td id = "` + producto.ID.Hex() + `" almacen="` + productoPostgres.Almacen + `"><button type="button" class="btn btn-danger deleteButton" onclick="quitarProducto(this)">`
		html += `<span class="glyphicon glyphicon-trash btn-xs"></span>`
		html += `</button></td>`
		html += `</tr>`

		impuestos = append(impuestos, productoPostgres.Cantidad*valimpuesto)
		totales = append(totales, productoPostgres.Cantidad*(productoPostgres.Precio+valimpuesto))
		subtotales = append(subtotales, precio)
	}

	fmt.Println("Total Productos:", totales)
	var total, impuesto, subtotal float64
	total = 0
	impuesto = 0
	subtotal = 0
	for i := 0; i < len(totales); i++ {
		total += totales[i]
		impuesto += impuestos[i]
		subtotal += subtotales[i]
	}
	/// template para calculadora de totales
	calculadora += `<tr>`
	calculadora += `<td>SubTotal:</td>`
	calculadora += `<td>` + money.FormatMoney(subtotal) + `</td>`
	calculadora += `</tr>`
	calculadora += `<tr>`
	calculadora += `<td>Impuestos:</td>`
	calculadora += `<td>` + money.FormatMoney(impuesto) + `</td>`
	calculadora += `</tr>`
	calculadora += `<tr>`
	calculadora += `<td>Total:</td>`
	calculadora += `<td>` + money.FormatMoney(total) + `</td>`
	calculadora += `</tr>`

	calculadora += `<script>`
	calculadora += `$("input[name='Producto']").keydown(function(e) {
						if (e.which == 13 || e.keyCode == 13) {
							e.preventDefault();
						}
					});`
	calculadora += `</script>`
	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()

	return html, calculadora
}
