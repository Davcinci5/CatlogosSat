package PuntoVentaModel

import (
	"fmt"
	"strconv"
	"time"

	"github.com/leekchan/accounting"

	"../../Modelos/AlmacenModel"
	"../../Modelos/ProductoModel"
	"../../Modelos/UnidadModel"
	"../../Modelos/UsuarioModel"
	"../../Modulos/ConsultasSql"

	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Imagenes"
	"../../Modulos/Variables"

	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//PuntoVentaMgo estructura de PuntoVentas mongo
type PuntoVentaMgo struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Operacion bson.ObjectId `bson:"Operacion,omitempty"`
	Usuario   bson.ObjectId `bson:"Usuario,omitempty"`
	Equipo    bson.ObjectId `bson:"Equipo,omitempty"`
	Codigo    string        `bson:"Codigo"`
	Carrito   string        `bson:"Carrito"`
	Resumen   string        `bson:"Resumen"`
	Estatus   bson.ObjectId `bson:"Estatus,omitempty"`
	FechaHora time.Time     `bson:"FechaHora"`
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
func GetAll() []PuntoVentaMgo {
	var result []PuntoVentaMgo
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}
	err = PuntoVentas.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) PuntoVentaMgo {
	var result PuntoVentaMgo
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}
	err = PuntoVentas.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []PuntoVentaMgo {
	var result []PuntoVentaMgo
	var aux PuntoVentaMgo
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = PuntoVentaMgo{}
		PuntoVentas.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) PuntoVentaMgo {
	var result PuntoVentaMgo
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)

	if err != nil {
		fmt.Println(err)
	}
	err = PuntoVentas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result PuntoVentaMgo
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}
	err = PuntoVentas.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

// GetIDUsuario Regresa el ID del Usuario especificado
func GetIDUsuario(NombreUsuario string) string {
	return UsuarioModel.GetIDByField("Usuario", NombreUsuario).Hex()
}

//########################< FUNCIONES GENERALES PSQL >#############################

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

	Query := fmt.Sprintf(`SELECT "Movimiento", "Producto", "Almacen", "Cantidad", "Costo", "Precio", "Impuesto" FROM public."VentaTemporal" WHERE "Operacion" = '%v' ORDER BY "FechaHora" ASC`, idOperacion)
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
		resultSet.Scan(&producto.Movimiento, &producto.Producto, &producto.Almacen, &producto.Cantidad, &producto.Costo, &producto.Precio, &producto.Impuesto)
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
			html += `<td><input name="Cantidad" id = "` + producto.ID.Hex() + `" almacen="` + productoPostgres.Almacen + `" fraccion="` + strconv.FormatBool(producto.VentaFraccion) + `" previo="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 3, 64) + `" onchange="ValidaNumero(this)" onfocusout="AplicaPeticion(this)" type="number" class="form-control" value="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 3, 64) + `"></td>`
		} else {
			html += `<td><input name="Cantidad" id = "` + producto.ID.Hex() + `" almacen="` + productoPostgres.Almacen + `" fraccion="` + strconv.FormatBool(producto.VentaFraccion) + `" previo="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 1, 64) + `" onchange="ValidaNumero(this)" onfocusout="AplicaPeticion(this)" type="number" class="form-control" value="` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 1, 64) + `"></td>`
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
		html += `<td id = "` + producto.ID.Hex() + `"><button type="button" class="btn btn-danger deleteButton" onclick="quitarProducto(this)">`
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

	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()

	return html, calculadora
}

//GeneraTemplateBusquedaDeProductosEnElasticYAlmacenesParaPuntoVenta Como su nombre indica regresa un template
//para mostrar en el Punto De Venta una búsqueda en Elastic por Almacenes
func GeneraTemplateBusquedaDeProductosEnElasticYAlmacenesParaPuntoVenta(Productos []ProductoModel.ProductoMgo) string {
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
					busqueda += `<tr>`
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
						busqueda += `<td id="thisone"><input name="Cantidad" id = "` + v.ID.Hex() + `" almacen="` + b.Hex() + `" fraccion="` + strconv.FormatBool(v.VentaFraccion) + `" previo="` + strconv.FormatFloat(Cantidad, 'f', 3, 64) + `" onchange="ValidaNumeroModal(this)" onblur="ValidaNumeroModal(this)" type="number" class="form-control" value=""></td>`
					} else {
						busqueda += `<td id="thisone"><input name="Cantidad" id = "` + v.ID.Hex() + `" almacen="` + b.Hex() + `" fraccion="` + strconv.FormatBool(v.VentaFraccion) + `" previo="` + strconv.FormatFloat(Cantidad, 'f', 1, 64) + `" onchange="ValidaNumeroModal(this)" onblur="ValidaNumeroModal(this)" type="number" class="form-control" value=""></td>`
					}
					// busqueda += `<td><input type="checkbox" name="ModalSeleccione" value="` + v.ID.Hex() + `"></td>`
					busqueda += `</tr>`
				}
			}

		}

	}

	return busqueda
}

//CreaComboNumeros crea el combo de números que se le indique.
func CreaComboNumeros(Num int) string {

	templ := `<select class="form-control">`
	templ += `<option value="">--SELECCIONE--</option>`

	for i := 1; i <= Num; i++ {
		templ += `<option value="` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</option>`
	}
	templ += `</select>`
	return templ
}

//AplicaOperacionParaPago alica a los almacenes de postgresql los moviientos generados en la preventa del carrito
func AplicaOperacionParaPago(Operacion string) bool {

	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		fmt.Println("Error al iniciar la sesion psql", err)
		return false
	}

	Query := fmt.Sprintf(`SELECT "Operacion","Movimiento", "Almacen", "Producto", "Cantidad", "Costo", "Precio", "Impuesto", "Descuento", "Existencia" FROM public."VentaTemporal" WHERE "Operacion" = '%v'`, Operacion)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		fmt.Println("Error al preparar la consulta", err)
		return false
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return false
	}

	Datas := []DatosVentaTemporal{}
	for resultSet.Next() {
		Data := DatosVentaTemporal{}
		resultSet.Scan(&Data.Operacion, &Data.Movimiento, &Data.Almacen, &Data.Producto, &Data.Cantidad, &Data.Costo, &Data.Precio, &Data.Impuesto, &Data.Descuento, &Data.Existencia)
		Datas = append(Datas, Data)
	}
	for _, v := range Datas {
		v.InsertaKardexAlmacen()
		//v.InsertaImpuestoAlmacen()
		//v.InsertaDescuentoAlmacen()
	}

	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()

	return true
}

//ActualizaOperacionParaPago alica a los almacenes de postgresql los moviientos generados en la preventa del carrito
func ActualizaOperacionParaPago(Operacion string) bool {
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		fmt.Println(err)
		return false
	}

	Query := fmt.Sprintf(`SELECT "Movimiento", "Almacen", "Producto", "Cantidad", "Costo", "Precio", "Impuesto", "Descuento", "Existencia" FROM public."VentaTemporal" WHERE "Operacion" = '%v'`, Operacion)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return false
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return false
	}

	Datas := []DatosVentaTemporal{}
	for resultSet.Next() {
		Data := DatosVentaTemporal{}
		resultSet.Scan(&Data.Movimiento, &Data.Almacen, &Data.Producto, &Data.Cantidad, &Data.Costo, &Data.Precio, &Data.Impuesto, &Data.Descuento, &Data.Existencia)
		Datas = append(Datas, Data)
	}

	for _, v := range Datas {
		v.ActualizaKardexAlmacen()
		//v.InsertaImpuestoAlmacen()
		//v.InsertaDescuentoAlmacen()
	}

	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()
	return true
}

//GetHTMLMontoOperacion regresa un template con un ticket para Cobro
func GetHTMLMontoOperacion(Operacion string) string {
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		fmt.Println("Error al iniciar sesion en Psql", err)
		return ""
	}

	Query := fmt.Sprintf(`SELECT "Movimiento", "Almacen", "Producto", "Cantidad", "Costo", "Precio", "Impuesto", "Descuento", "Existencia" FROM public."VentaTemporal" WHERE "Operacion" = '%v'`, Operacion)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		fmt.Println("Error al preparar la consulta", err)
		return ""
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	Datas := []DatosVentaTemporal{}
	for resultSet.Next() {
		Data := DatosVentaTemporal{}
		resultSet.Scan(&Data.Movimiento, &Data.Almacen, &Data.Producto, &Data.Cantidad, &Data.Costo, &Data.Precio, &Data.Impuesto, &Data.Descuento, &Data.Existencia)
		Datas = append(Datas, Data)
	}

	var Monto float64
	var Impuestos float64

	for _, v := range Datas {
		Monto += v.Cantidad * v.Precio
		Impuestos += v.Impuesto
	}
	money := accounting.Accounting{Symbol: "$", Precision: 2}

	templ := `<h1 style="text-align: center;"><span style="text-decoration: underline;">Referencia De Compra Para Pago</span></h1>
			<p>&nbsp;</p>
			<hr>
			<h2 style="text-align: center;"><em>#Operaci&oacute;n : ` + Operacion + `</em></h2>
			<p>&nbsp;</p>
			<h2 style="text-align: center;"><em>RESUMEN DE COMPRA</em></h2>			
			<h2 style="text-align: center;">			
			 		<div style="text-align: center;">
                      <table style="text-align: center;">
                        <thead>
                          <tr>
                            <th>Concepto</th>
                            <th>Monto</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr>
                            <td>SubTotal:</td>
                            <td>` + money.FormatMoney(Monto) + `</td>
                          </tr>
                          <tr>
                            <td>Impuestos:</td>
                            <td>` + money.FormatMoney(Impuestos) + `</td>
                          </tr>
                          <tr>
                            <td>Total a Pagar:</td>
                            <td>` + money.FormatMoney(Monto+Impuestos) + `</td>
                          </tr>
                        </tbody>
                      </table>                      
                  </div>
			</h2>
			<h2 style="text-align: center;">&nbsp;</h2>
			<hr>
			<h1 style="text-align: center;"><em>Gracias por su Compra :)</em></h1>`

	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()
	return templ

}

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoPuntoVenta, queryTilde)
	if err {
		fmt.Println("Ocurrió un error al consultar en Elastic en el primer intento")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoPuntoVenta, queryQuotes)
		if err {
			fmt.Println("Ocurrió un error al consultar en Elastic en el segundo intento")
		}
	}

	return docs
}

//######################### RUTINAS ADICIONALES ###################

//GeneraBotonesPago regresa el template de botones para procesar pago con una operacion precargada
func GeneraBotonesPago(operacion string) string {
	templ := `<div class="col-md-4">
					<button onclick="window.location.href = '/PuntoVentas/edita/` + operacion + `';" type="button" class="btn-lg btn-secondary btn-block" style="height: 200px;">
						<h1 style="color: black;">Edita Compra</h1>
					</button>
				</div>
				<div class="col-md-4">
					<button onclick="generaTicket('` + operacion + `')" type="button" class="btn-lg btn-primary btn-block" style="height: 200px;">
					<h1 style="color: white;">Imprime Referencia</h1>
					</button>
				</div>
				<div class="col-md-4">
					<button onclick="window.location.href = '/PuntoVentas/cobra/` + operacion + `'; return false;" class="btn-lg btn-success btn-block" style="height: 200px;">
					<h1 style="color: white;">Cobrar</h1>
					</button>
				</div>`

	return templ
}

//ConsultaDatosOperacionAll devuelve el listado total de operaciones pendientes de pago by @melchormendoza
func ConsultaDatosOperacionAll() (string, []DatosVentaTemporal) {

	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	Query := fmt.Sprintf(`SELECT sum("Cantidad"),sum("Cantidad"*"Precio"),sum("Impuesto"),"Operacion","Almacen" FROM public."VentaTemporal" GROUP BY "Operacion","Almacen"`)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	productos := []DatosVentaTemporal{}
	for resultSet.Next() {
		producto := DatosVentaTemporal{}
		resultSet.Scan(&producto.Cantidad, &producto.Precio, &producto.Impuesto, &producto.Operacion, &producto.Almacen)
		productos = append(productos, producto)
	}

	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()
	return "_", productos
}

//ConsultaDatosOperacionByID devuelve el detalle de una operacion para el cobro by @melchormendoza
func ConsultaDatosOperacionByID(idOperacion string) (string, string, float64, string) {
	html := ``
	calculadora := ``
	var totales []float64   //para almacenar los totales por poducto
	var impuestos []float64 //para almacenar los impuestos por producto
	money := accounting.Accounting{Symbol: "$", Precision: 2}
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		fmt.Println(err)
		return "", "", 0.00, "Error al iniciar Sesion"
	}
	Query := fmt.Sprintf(`SELECT "Movimiento", "Producto", "Almacen", "Cantidad", "Costo", "Precio","Impuesto" FROM public."VentaTemporal" WHERE "Operacion" = '%v'`, idOperacion)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		fmt.Println(err)
		return "", "", 0.00, ""
	}
	resultSet, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return "", "", 0.00, ""
	}

	productos := []DatosVentaTemporal{}
	for resultSet.Next() {
		producto := DatosVentaTemporal{}
		resultSet.Scan(&producto.Movimiento, &producto.Producto, &producto.Almacen, &producto.Cantidad, &producto.Costo, &producto.Precio, &producto.Impuesto)
		productos = append(productos, producto)
	}

	for _, productoPostgres := range productos {
		producto := ProductoModel.GetOne(bson.ObjectIdHex(productoPostgres.Producto))

		html += `<tr id = "` + productoPostgres.Movimiento + `">`
		html += `<td>`
		aux := ``
		for _, valores := range producto.Codigos.Valores {
			aux += valores + `,`
		}
		html += aux[:len(aux)-1]
		html += `</td>`
		html += `<td>` + producto.Nombre + `</td>`
		html += `<td>` + strconv.FormatFloat(productoPostgres.Cantidad, 'f', 1, 64) + `</td>`
		html += `<td>` + money.FormatMoney(productoPostgres.Precio) + `</td>`
		unidad := UnidadModel.GetSubEspecificByFields("Datos._id", producto.Unidad)
		html += `<td>` + unidad.Abreviatura + `</td>`
		precio := productoPostgres.Cantidad * productoPostgres.Precio
		//html += `<td>` + money.FormatMoney(precio) + `</td>`
		html += `<td>` + money.FormatMoney(precio) + `</td>`
		//html += `<td><select class="selectpicker" aria-invalid="false" disabled><option>Mostrador</option></select></td>`
		//html += `<td><select class="selectpicker" aria-invalid="false" disabled><option>Sin Ruta</option><option>Transporte 1</option></select></td>`
		//html += `<td id = "` + producto.ID.Hex() + `"><button type="button" class="btn btn-danger deleteButton" onclick="quitarProducto(this)">`
		//html += `<span class="glyphicon glyphicon-trash btn-xs"></span>`
		//html += `</button></td>`
		html += `</tr>`

		impuestos = append(impuestos, productoPostgres.Impuesto)
		totales = append(totales, precio)
	}
	fmt.Println("Total Productos:", totales)
	var total, impuesto float64
	total = 0
	impuesto = 0
	for i := 0; i < len(totales); i++ {
		total += totales[i]
		impuesto += impuestos[i]
	}
	/// template para calculadora de totales
	calculadora += `<tr>`
	calculadora += `<td>SubTotal:</td>`
	calculadora += `<td>` + money.FormatMoney(total) + `</td>`
	calculadora += `</tr>`
	calculadora += `<tr>`
	calculadora += `<td>Impuestos:</td>`
	calculadora += `<td>` + money.FormatMoney(impuesto) + `</td>`
	calculadora += `</tr>`
	calculadora += `<tr>`
	calculadora += `<td>Total:</td>`
	calculadora += `<td>` + money.FormatMoney(total) + `</td>`
	calculadora += `</tr>`

	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()

	return html, calculadora, total, idOperacion
}
