$(document).ready(function(){

	
	$('#cabeceraajuste').hide();
	$('#cabeceratraslado').hide();

	$('#contenedortabla').hide();
	$('#contenedorbotonguardar').hide();
	$('#contenedornoencontrado').hide();

		$("body").on("click",".eliminar", function(e){
   			e.preventDefault();
    		$(this).closest('tr').remove(); 
    		var articulos_en_cuestion = $('#contenedordeproductos tr').length;
    		if (articulos_en_cuestion == 0){
    				$('#contenedortabla').hide();
					$('#contenedorbotonguardar').hide();
					$('#contenedornoencontrado').hide();

    		}      	      
	       // return false;
	   });
});

function AjusteTraslado(){

    var origen = $("#s_origen option:selected").val();
	var destino = $("#s_destino option:selected").val();

	
	
	$('#contenedortabla').hide();
	$('#contenedornoencontrado').empty();
	$('#contenedordeproductos').empty();		
	$('#contenedorbotonguardar').hide();

	if (origen == "Seleccione Origen.." || destino == "Seleccione Destino.."){
		$('#contenedorajustes').empty();
	}

	if (origen != "Seleccione Origen.." && destino != "Seleccione Destino.."){
		
		if (origen == destino){
			
			$.ajax({
			url: '/Almacens/Operacion/Movimientos',
			type: 'POST',
			dataType: 'html',
			data:{origen:origen,destino:destino,tipo:"ajuste"},
			success : function(data){
				$('#cabeceraajuste').show();
				$('#cabeceratraslado').hide();				
				$('#contenedorajustes').html(data);
				$('#contenedorajustes').focus();
				}
			});

		}else{

			$.ajax({
			url: '/Almacens/Operacion/Movimientos',
			type: 'POST',
			dataType: 'html',
			data:{origen:origen,destino:destino,tipo:"traslado"},
			success : function(data){
				$('#cabeceraajuste').hide();
				$('#cabeceratraslado').show();
				$('#contenedorajustes').html(data);
				$('#contenedorajustes').focus();

				}
			});		

		}
	}

}

function GetArticulo(){
	var inputCodigo = $('#elarticulo'); 
    var noEncontrado = $('#contenedornoencontrado'); 
    var AgregarCantidad = $("input:checkbox[name='agregarcantidad']");
    var IsAgregarOn = AgregarCantidad.is(":checked")

    console.log("Objeto -> ",AgregarCantidad);
    
    var Productos = $('#contenedordeproductos'); 
    var Tabla = $('#contenedortabla'); 
    var Guardar =$('#contenedorbotonguardar')    
    
	var cod_b_art =	$('#elarticulo').val();
	var origen = $("#s_origen option:selected").val();
	var destino = $("#s_destino option:selected").val();
    var tipo_movimiento = $('#tipomovimiento').val();

	if (cod_b_art != ""){
        
        var agregado = false;
	    var codigosdearticulos = $(".codigosdearticulos");

        codigosdearticulos.each(function(){
           if ($(this).html() == cod_b_art){
                ids = $(this).attr('id');
                idsplit = ids.split(":");
                inputTarget = $("#"+idsplit[0]);
                agregado = true;
            }
        });
        
		var articulos_en_cuestion = $('#contenedordeproductos tr').length;
		noEncontrado.empty();
    
    if(agregado){
            
            inputCodigo.val("");                
            noEncontrado.empty();
            if (confirm("Ya has agregado "+ inputTarget.val()+", ¿cambiar la cantidad?")){    
                inputTarget.select();
                inputTarget.focus();
            }else{
                inputCodigo.val("");
                inputCodigo.focus();
            }
	}else{		
			$.ajax({
				url: '/Almacens/Operacion/Movimientos/Ajustar',
				type: 'POST',
				dataType: 'html',
				data:{cod_b_art:cod_b_art,origen:origen,destino:destino,articulos_agregados:articulos_en_cuestion,tipomovimiento:tipo_movimiento,agregarcantidad:IsAgregarOn},
				success : function(data){
			
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
                            if (AgregarCantidad.is(":checked") ){
                                inputCodigo.val("");                                
                            }else{
                                inputCodigo.val("");
                                inputCodigo.focus();
                            }
					    	noEncontrado.empty();
					       	Guardar.show();
							Tabla.show();
							Productos.append(data);
					}
				}
			});						
	}
	
	}else{
		noEncontrado.empty();
        inputCodigo.focus();
		alert("Introduce un código de barra");
	}
}


function Ejecutar(){

	var i = 0;
	var codigos = []
	var nombres = []
	var precios = []
	var existencias = []
	var operaciones = []
	var cantidades = []
	var existenciasd = []
	var origen = $("#s_origen option:selected").val();
	var destino = $("#s_destino option:selected").val();

	$("#contenedordeproductos tr.renglon").each(function(){
	
		var cds = $('#codigo_b'+i).html()
		var nms = $('#desc_b'+i).html()
		var prc = $('#precio_b'+i).html()
		var exs = $('#origen_b'+i).html()


		if (origen == destino){
			var ope = $('input:radio[name=operacion'+cds+']:checked').val(); 
		}else{
			var exsd = $('#destino_b'+i).html()
			existenciasd.push(exsd)
		}


		var cant = $('#'+cds).val()

		codigos.push(cds)
		nombres.push(nms)
		precios.push(prc)
		existencias.push(exs)
		operaciones.push(ope)
		cantidades.push(cant)
		i++;        	 
    });

    $.ajax({
		url: '/Almacens/Operacion/Movimientos/Realizar',
		type: 'POST',
		dataType: 'html',
		data:{ codigos:codigos, nombres:nombres,precios:precios, existencias:existencias, operaciones:operaciones, cantidades:cantidades,origen:origen,destino:destino,existenciasd:existenciasd},
		success : function(data){
			if (data == "no"){
				alert("No se puede surtir desde el almacen de origen.")
				location.reload("/ajustetraslado");
			}else{
				alert("Operacion Existosa")
				location.reload("/ajustetraslado");
			}
			}
		});

}


function BuscarArticulo(){
	alert("Ir a buscar articulo");	
}


function FocusInputCodigos(){
    var inputCodigo = $('#elarticulo'); 
    inputCodigo.val("");
    inputCodigo.focus();
    return false
}

//obtenerAlmacenDefecto lee el select que se muestra en la vista y regresa el nombre del almacen seleccionado
function obtenerAlmacenDefecto(){
	var almacenDefecto = $('#s_origen').val();
	return almacenDefecto;
}

//obtenerProductoIngresado lee el input de la vista en donde se ingresa el nombre o codigo del producto y regresa el valor ingresado
function obtenerProductoIngresado(){
	return $('#scanner_search3').val();
}

//leerFilaProducto recorre la ventana modal en donde se muestan las caracteristicas del producto
//Al seleccionar el boton "selecciona" cierra la modal y guarda los atributos en una coleccion
function leerFilaProducto(identificador) {
	var id = $("#filaProducto" + identificador);
	var numeros_hijos = id.children().length;
	var datos = id.children();
	var trs = [];
	var filas = [];

	for (var i = 0; i < numeros_hijos; i++) {
		if (i < numeros_hijos - 1) {
			trs.push(datos[i]);
		}
	}
	alert(id)
	// filas.push(trs);
	// $("#CerrarModal").trigger("click");
	// valida();
	// var elementos = crearEtiquetasProductos();
	// MostrarDetallesProducto(identificador, filas, elementos);
}

//buscarProductos realiza una busqueda general en la coleccion de mongo, sin tomar en cuenta el almacen por defecto
//posteriormente se tendrá que buscar en el almacen por defecto
//Obtiene los atributos del producto y lo muestra en una ventana modal
function buscarProductos() {
 	var almacenDefecto = obtenerAlmacenDefecto();
 	var producto = obtenerProductoIngresado();
 	//var validator = valida();
 	$("#Carrito").empty();
 	if (almacenDefecto != "") {
  		if (producto != "") {
  			$.ajax({
    			url: "/ConsultarProductos",
    			type: 'POST',
    			dataType: 'json',
    			data: { nombreProducto: producto },
    			success: function (data) {
					if(data.SEstado == false){
						alertify.error(data.SMsj); 
					}else{
						$("#Carrito").empty();
						$("#Carrito").append(data.SIhtml);
						// $('#btnModal').trigger("click");
					}
   				}
   			});
   
  		}else {
//    validator.showErrors({
//    "Codigo":"El campo no debe estar vacio"
//   });
  	    }
 	}else {
//   validator.showErrors({
//    "AlmacenDefecto":"Debe seleccionar un almacen"
//   });
 	}
}