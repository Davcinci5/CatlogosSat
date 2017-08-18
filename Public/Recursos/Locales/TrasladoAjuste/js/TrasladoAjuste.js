$(document).ready(function () {


	$('#cabeceraajuste').hide();
	$('#cabeceratraslado').hide();

	$('#contenedortabla').hide();
	$('#contenedorbotonguardar').hide();
	$('#contenedornoencontrado').hide();

	$("body").on("click", ".eliminar", function (e) {
		e.preventDefault();
		$(this).closest('tr').remove();
		var articulos_en_cuestion = $('#contenedordeproductos tr').length;
		if (articulos_en_cuestion == 0) {
			$('#contenedortabla').hide();
			$('#contenedorbotonguardar').hide();
			$('#contenedornoencontrado').hide();

		}
		// return false;
	});
});

function AjusteTraslado() {

	var origen = $("#s_origen option:selected").val();
	var destino = $("#s_destino option:selected").val();

	$('#contenedortabla').hide();
	$('#contenedornoencontrado').empty();
	$('#contenedordeproductos').empty();
	$('#contenedorbotonguardar').hide();

	if (origen == "Seleccione Origen.." || destino == "Seleccione Destino..") {
		$('#contenedorajustes').empty();
	}

	if (origen != "Seleccione Origen.." && destino != "Seleccione Destino..") {
		var tipomov = "";
		if (origen == destino) {
					$("#textoTipoMovimiento").html("Ajuste:");
					$("#tipomovimiento").val("ajuste");
					$('#cabeceraajuste').show();
					$('#cabeceratraslado').hide();			
		} else {

					$('#cabeceraajuste').hide();
					$('#cabeceratraslado').show();
					$("#textoTipoMovimiento").html("Traslado:");
					$("#tipomovimiento").val("traslado");	
		}
				document.getElementById("elarticulo").disabled=false;
				$('#contenedorajustes').focus();
	}
}

function AgregarArticulo(ID) {

	var inputCodigo = $('#elarticulo');
	var noEncontrado = $('#contenedornoencontrado');
	var AgregarCantidad = $("input:checkbox[name='agregarcantidad']");
	var IsAgregarOn = AgregarCantidad.is(":checked")

	var Productos = $('#contenedordeproductos');
	var Tabla = $('#contenedortabla');
	var Guardar = $('#contenedorbotonguardar')

	var origen = $("#s_origen option:selected").val();
	var destino = $("#s_destino option:selected").val();
	var tipo_movimiento = $('#tipomovimiento').val();

	var agregado = false;
	var IdArtArticulos = $(".renglon");

	 IdArtArticulos.each(function () {
		if ($(this).attr("id") == ("row-"+ID)) {
	 		agregado = true;
	 	}
	 });

	 var rows_articulos = $('#contenedordeproductos tr').length;
	// noEncontrado.empty();
	 if (agregado) {
		 alertify.alert("Ya as agregado el articulo cambiar la cantidad");
	 } else {
		$.ajax({
			url: '/TrasladoAjuste/GetArticuloAjustar',
			type: 'POST',
			dataType: 'html',
			data: { ID: ID, origen: origen, destino: destino, articulos_agregados: rows_articulos, tipomovimiento: tipo_movimiento, agregarcantidad: IsAgregarOn },
			success: function (data) {

				switch (data) {

					case '<h3 class="text-center">Origen: Desactivado     Destino: Activo</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					case '<h3 class="text-center">Origen: Bloqueado      Destino: Activo</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					case '<h3 class="text-center">Origen: No Existe     Destino: Activo</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					case '<h3 class="text-center">Origen: Activo      Destino: Desactivado</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					case '<h3 class="text-center">Origen: Activo      Destino: Bloqueado</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					case '<h3 class="text-center">Origen: Activo      Destino: No Existe</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					case '<h3 class="text-center">Articulo no encontrado</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					case '<h3 class="text-center">Origen: No Existe      Destino: No Existe</h3>':
						noEncontrado.show();
						noEncontrado.empty();
						noEncontrado.append(data);
						inputCodigo.val("");
						inputCodigo.focus();
						break;
					default:
						if (AgregarCantidad.is(":checked")) {
							inputCodigo.val("");
						} else {
							inputCodigo.val("");
							inputCodigo.focus();
						}
						noEncontrado.empty();
						Guardar.show();
						Tabla.show();
						Productos.append(data);
				}
			},
			complete:function(){
				$("#CerrarModal").trigger("click");
			}
		});
	 }
}


function Ejecutar() {

	var i = 0, ii = 0, contador = 0;
	var codigos = [];
	var nombres = [];
	var precios = [];
	var existencias = [];
	var operaciones = [];
	var cantidades = [];
	var existenciasd = [];
	var origen = $("#s_origen option:selected").val();
	var destino = $("#s_destino option:selected").val();

	$("#contenedordeproductos tr.renglon").each(function () {
		var cds = $('#codigo_b' + ii).html();
		var cant = $('#' + cds).val();
		var cantMax = $('#origen_b' + ii).text();

		if (cant > 0) {
			ii++;
		} else {
			var nms = $('#desc_b' + ii).html();
			alertify.error("Capturar Campos correctos en: " + nms);
			ii++;
			contador++;
		}
	});

	if (contador > 0) {
		return false;
	}

	$("#contenedordeproductos tr.renglon").each(function () {

		var cds = $('#codigo_b' + i).html();
		var nms = $('#desc_b' + i).html();
		var prc = $('#precio_b' + i).html();
		var exs = $('#origen_b' + i).html();

		if (origen == destino) {
			var ope = $('input:radio[name=operacion' + cds + ']:checked').val();
		} else {
			var exsd = $('#destino_b' + i).html();
			existenciasd.push(exsd);
		}
		var cant = $('#' + cds).val();

		codigos.push(cds);
		nombres.push(nms);
		precios.push(prc);
		existencias.push(exs);
		operaciones.push(ope);
		cantidades.push(cant);
		i++;
	});

	$.ajax({
		url: '/TrasladoAjuste/Realizar',
		type: 'POST',
		dataType: 'html',
		data: { codigos: codigos, nombres: nombres, precios: precios, existencias: existencias, operaciones: operaciones, cantidades: cantidades, origen: origen, destino: destino, existenciasd: existenciasd },
		success: function (data) {
			if (data == "no") {
				alertify.error("No se puede surtir desde el almacen de origen.");
				// location.reload("/ajustetraslado");
			} else {
				location.href = "/TrasladoAjuste/Movimientos?exito=1";
				alertify.success("Operacion Existosa");
				//location.reload("/ajustetraslado");
			}
		}
	});
}


function BuscarArticulo() {
	alert("Ir a buscar articulo");
}

function FocusInputCodigos() {
	var inputCodigo = $('#elarticulo');
	inputCodigo.val("");
	inputCodigo.focus();
	return false
}

//obtenerAlmacenDefecto lee el select que se muestra en la vista y regresa el nombre del almacen seleccionado
function obtenerAlmacenDefecto() {
	var almacenDefecto = $('#s_origen').val();
	return almacenDefecto;
}

//obtenerProductoIngresado lee el input de la vista en donde se ingresa el nombre o codigo del producto y regresa el valor ingresado
function obtenerProductoIngresado() {
	return $('#scanner_search3').val();
}

//buscarProductos realiza una busqueda general en la coleccion de mongo, sin tomar en cuenta el almacen por defecto
//posteriormente se tendrá que buscar en el almacen por defecto
//Obtiene los atributos del producto y lo muestra en una ventana modal
function buscarProducto() {
	var codigoPalabra = $('#elarticulo').val();
	if (codigoPalabra == "") {
		alertify.error("Capturar  el Nombre o Codigo del Articulo a buscar.");
	} else {
		$.ajax({
			url: "/TrasladoAjuste/Productos",
			type: "POST",
			dataType: "json",
			data: { nombreProducto: codigoPalabra },
			success: function (data) {
				if (data.SEstado == false) {
					alertify.error(data.SMsj);
				} else {
					$("#Carrito").empty();
					$("#Carrito").append(data.SIhtml);
					$('#modal-search').modal('show');
					// $('#btnModal').trigger("click");
				}
			}
		});
	}
}