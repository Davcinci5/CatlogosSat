

	//##############################< SCRIPTS JS >##########################################
	//################################< Cotizacion.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################
	
	$( document ).ready( function () {
		var validator = valida();
		$("#FormaDePagoDetalle").val($("#FormaDePago").val());
		$("#FormaDeEnvíoDetalle").val($("#FormaDeEnvío").val());

		if( $("#BClientes").children().length == 0) {
			$("#tabcliente").hide();				
		}
		$("#tabBusquedaCliente").hide();	
		$("#tabBusquedaProductos").hide();
		if ( $("#BProductos").children().length == 0 ){
			$("#tabProductos").hide();
		}	
		if( $("#BClientes").children().length > 0){
			$("#DatosMinimos").hide();
		}
		
		$( "#Nombre,#Telefono,#Correo,#FormaDePago,#FormaDeEnvío" ).change(function() {
			ActualizaDetalle();
		});

		$('#ClienteNom').keydown(function(e) {
			if (e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
				traeCliente();
			}
		});
		$('#Buscar').keydown(function(e) {
			if (e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
				traeProductos();
			}
		});

		$('#Nombre').keydown(function(e) {
			if (e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
				$("#Telefono").focus();
			}
		});
		$('#Telefono').keydown(function(e) {
			if (e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
				$("#Correo").focus();
			}
		});

		$("#Correo").keydown(function(e) {
			if (e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
			}
		});

		$("input[name='Cantidad']").keydown(function(e) {
			if (e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
			}
		});

			

		$('#rootwizard').bootstrapWizard({
			  'tabClass': 'nav nav-pills',
			// 'onTabClick':function(tab, navigation, index) {
			// 				return false;
			// 			},
	  		'onNext': function(tab, navigation, index) {
	  			var $valid = $("#Form_Alta_Cotizacion").valid();
	  			if(!$valid) {
	  				validator.focusInvalid();
	  				return false;
	  			}
	  		}
		  });
		  
		  $(document).on('click', '.deleteButton', function () {
			$(this).parent().parent().remove();
			$('#tabcliente').hide();
			$("#DatosMinimos").show();
		});

		$("#Form_Alta_Cotizacion").submit(function(){
			
			var somestatus = true;
			if ($("#BClientes").children().length === 0){
				if( $("#Nombre").val() === "" && $("#Correo").val() === "" || $("#Telefono").val()==="" ){
					alertify.error("Debes agregar almenos los datos minimos del cliente");
					$("a[title='Cliente']").click();
					$("#Nombre").valid();
					$("#Correo").valid() ;
					$("#Telefono").valid();
					somestatus = false;
				}else{
					if( !$("#Nombre").valid() && !$("#Correo").valid() || !$("#Telefono").valid()){
						$("a[title='Cliente']").click();
						somestatus = false;
					}
				}				
			}

			if (somestatus) {
				if ( $("#BProductos").children().length === 0 ){
					somestatus = false;
					$("a[title='Productos']").click()
					alertify.error("Agrega al menos un producto al carrito")
				}
			}
			if (somestatus) {
				if( $("#FormaDePago").val() === "" || $("#FormaDeEnvío").val() === "" ){
					somestatus = false;
					$("a[title='Formas De Pago']").click();
					$("#FormaDePago").valid();
					$("#FormaDeEnvío").valid();
					alertify.error("Los campos de Forma de Pago y Entrega son obligatorios")
				}
			}

			if(!somestatus){
				event.preventDefault();
			} 

			

			
		});
		 
	});

    function valida(){
		var validator = $("#Form_Alta_Cotizacion").validate({
			rules: {
				Nombre:{
					minlength:8
				},
				Telefono:{
					rangelength:[10,13]
				},
				Correo:{
					email:true
				},
				FormaDePago:{
					required: true
				},
				FormaDeEnvío:{
					required:true
				}			
		},
		messages: {
				Nombre:{
					minlength: "Debe contener al menos 8 caracteres"
				},
				Telefono:{
					rangelength:"La longitud debe estar entre [10,13]"
				},
				Correo:{
					email: "Correo invalido"
				},
				FormaDePago:{
					required: "Campo Requerido"
				},
				FormaDeEnvío:{
					required: "Campo requerido"
				}			
		},
		errorElement: "em",
		errorPlacement: function ( error, element ) {
			error.addClass( "help-block" );
			element.parents( ".col-sm-5" ).addClass( "has-feedback" );

			if ( element.prop( "type" ) === "checkbox" ) {
				error.insertAfter( element.parent( "label" ) );
			} else {
				error.insertAfter( element );
			}

			if ( !element.next( "span" )[ 0 ] ) {
				$( "<span class='glyphicon glyphicon-remove form-control-feedback'></span>" ).insertAfter( element );
			}
		},
		success: function ( label, element ) {
			if ( !$( element ).next( "span" )[ 0 ] ) {
				$( "<span class='glyphicon glyphicon-ok form-control-feedback'></span>" ).insertAfter( $( element ) );
			}
		},
		highlight: function ( element, errorClass, validClass ) {
			$( element ).parents( ".col-sm-5" ).addClass( "has-error" ).removeClass( "has-success" );
			$( element ).next( "span" ).addClass( "glyphicon-remove" ).removeClass( "glyphicon-ok" );
		},
		unhighlight: function ( element, errorClass, validClass ) {
			$( element ).parents( ".col-sm-5" ).addClass( "has-success" ).removeClass( "has-error" );
			$( element ).next( "span" ).addClass( "glyphicon-ok" ).removeClass( "glyphicon-remove" );
		}
		});	

	return validator;

	}

function EditaCotizacion(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Cotizacions').val() != ""){
			window.location = '/Cotizacions/edita/' + $('#Cotizacions').val();
		}else{
			alertify.error("Debe Seleccionar un Cotizacion para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Cotizacions/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Cotizacions';
		}
	}

}


function DetalleCotizacion(){
	if ($('#Cotizacions').val() != ""){
		window.location = '/Cotizacions/detalle/' + $('#Cotizacions').val();
	}else{
	alertify.error("Debe Seleccionar un Cotizacion para editar");
	}
}



function BuscaPagina(num){
			$('#Loading').show();
			$.ajax({
			url:"/Cotizacions/search",
			type: 'POST',
			dataType:'json',
			data:{
				Pag : num,
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#Cabecera").empty();						
						$("#Cabecera").append(data.SCabecera);
						$("#Cuerpo").empty();						
						$("#Cuerpo").append(data.SBody);
						$("#Paginacion").empty();		
						$("#Paginacion").append(data.SPaginacion);						
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}				
				$('#Loading').hide();	 
			},
		  error: function(data){
				$('#Loading').hide();
		  },
		});
}


 function SubmitGroup(){
	 $('#Loading').show();
			$.ajax({
			url:"/Cotizacions/agrupa",
			type: 'POST',
			dataType:'json',
			data:{
				Grupox : $('#Grupos').val(),
				searchbox: $('#searchbox').val()
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#Cabecera").empty();						
						$("#Cabecera").append(data.SCabecera);
						$("#Cuerpo").empty();						
						$("#Cuerpo").append(data.SBody);
						$("#Paginacion").empty();		
						$("#Paginacion").append(data.SPaginacion);						
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}
				$('#Loading').hide(); 
			},
		  error: function(data){
			  $('#Loading').hide();
		  },
		});
}



function traeCliente() {
    var cliente = $("#ClienteNom").val();
    if (cliente != "") {
        $('#Loading').show();
        $.ajax({
            url: "/Cotizacions/traerClientes",
            type: 'POST',
            dataType: 'json',
            data: {
                Cliente: cliente
            },
            success: function(data) {
                if (data != null) {
                    if ( data.SEstado ) {
						$("#tabBusquedaCliente").show();
						$("#BBusquedaClientes").empty();
						$("#BBusquedaClientes").append(data.SIhtml);
					} else {
                        alertify.error(data.SMsj);
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                }
                $('#Loading').hide();
            },
            error: function(a, b, c) {
                console.log(a);
                console.log(b);
                console.log(c);
                alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
                $('#Loading').hide();
            },
        });
        $("#ClienteNom").val("");
        $("#ClienteNom").focus();
    } else {
        var validator = valida();
        validator.showErrors({
            "ClienteNom": "Campo Obligatorio"
        });

        $("#ClienteNom").focus();
    }
}


function SeleccionarCliente(clid){
	var id = $(clid).parent().parent().attr("id")

	if (id != "") {
        $('#Loading').show();
        $.ajax({
            url: "/Cotizacions/traerCliente",
            type: 'POST',
            dataType: 'json',
            data: {
                Cliente: id
            },
            success: function(data) {
                if (data != null) {
                    if ( data.SEstado ) {
						$("#tabBusquedaCliente").hide();
						$("#BBusquedaClientes").empty();
						$("#BClientes").empty();
						$("#BClientes").append(data.SIhtml);
						$("#BClientesDetalle").empty();
						$("#BClientesDetalle").append(data.SIhtml);
						$("#tabcliente").show();
						$("#DatosMinimos").hide();
						ActualizaDetalle();
					} else {
                        alertify.error(data.SMsj);
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                }
                $('#Loading').hide();
            },
            error: function(a, b, c) {
                alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
                $('#Loading').hide();
            },
        });
    } else {
        alertify.error("Existio un problema con la referencia del Cliente");
    }
}



function SeleccionarProducto(clid){
	var estatus = ValidaNumero($(clid).parent().parent().find("input"));
	var id = $(clid).parent().parent().attr("id");
	var almacen = $(clid).parent().parent().attr("almacen");
	var Cot = $("#ID").val();
	var Cantidad = $(clid).parent().parent().find("#thisone").children().val();
	if ( estatus ){
		if (id != "") {
			$('#Loading').show();
			$.ajax({
				url: "/Cotizacions/TraerProducto",
				type: 'POST',
				dataType: 'json',
				data: {
					Producto: id,
					Almacen: almacen,
					Operacion: Cot,
					Cantidad: Cantidad
				},
				success: function(data) {
					if (data != null) {
						if ( data.SEstado ) {
							$("#tabBusquedaProductos").hide();
							$("#BBusquedaProductos").empty();
							$("#BProductos").empty();
							$("#BProductos").append(data.SIhtml);
							$("#tabProductos").show();
							$("#calculadora").empty();
							$("#calculadora").append(data.SCalculadora);
							//Detalle
							$("#BProductosDetalle").empty();
							$("#BProductosDetalle").append(data.SIhtml);
							$("#calculadoraDetalle").empty();
							$("#calculadoraDetalle").append(data.SCalculadora);
							alertify.success(data.SMsj);
						} else {
							alertify.error(data.SMsj);
						}
					} else {
						alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
					}
					$('#Loading').hide();
				},
				error: function(a, b, c) {
					alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
					$('#Loading').hide();
				},
			});
		} else {
			alertify.error("Existio un problema con la referencia del Cliente");
		}
	}else{
		$(clid).parent().parent().find("input").focus();
	}
}


function ValidaNumero(numero) {
    var Existencia = Number($(numero).attr("existencia"));
    var VentaFraccion = $(numero).attr("fraccion");
	var Valor = $(numero).val();
	
	var status = true;

	

    if (VentaFraccion === "true") {

        if (/^([0-9])*$/.test(Valor) || /^\d*\.?\d*$/.test(Valor)) {
            if (Number(Valor) < 0) {
                alertify.error("No debe introducir valores menores a cero.");
				numero.value = Number(0).toFixed(3);
				status = false;
            }
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
			// numero.value = Number(0).toFixed(3);
			status = false;
        }

    } else {
        if (/^([0-9])*$/.test(Valor)) {
            if (Number(Valor) <= 0) {
                alertify.error("No debe introducir valores menores o iguales a cero.");
				// numero.value = Number(1);
				status = false;
            } 
        } else if (/^\d*\.?\d*$/.test(Valor)) {
            alertify.error("No se pueden vender decimales de este producto.");
			// numero.value = Number(1);
			status = false;
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
			// numero.value = Number(1);
			status = false;
        }

	}
	return status;
}

function ValidaNumeroCarrito(numero) {
    var prev = Number(numero.getAttribute("previo"));
    var Existencia = Number(numero.getAttribute("existencia"));
    var VentaFraccion = numero.getAttribute("fraccion");
	var Valor = numero.value;
	

    if (VentaFraccion === "true") {

        if (/^([0-9])*$/.test(Valor) || /^\d*\.?\d*$/.test(Valor)) {
            if (Number(Valor) < 0) {
                alertify.error("No debe introducir valores menores a cero.");
                numero.value = prev;
            } 
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
            numero.value = prev;
        }

    } else {
        if (/^([0-9])*$/.test(Valor)) {
            if (Number(Valor) <= 0) {
                alertify.error("No debe introducir valores menores o iguales a cero.");
                numero.value = prev;
            }
        } else if (/^\d*\.?\d*$/.test(Valor)) {
            alertify.error("No se pueden vender decimales de este producto.");
            numero.value = prev;
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
            numero.value = prev;
        }

    }
}



function traeProductos() {
	var Prod = $("#Buscar").val();
	var Cli = $("#ClienteId").val();
	var Cot = $("#ID").val();
	
    if (Prod != "") {
		
        $('#Loading').show();
        $.ajax({
            url: "/Cotizacions/traeProductos",
            type: 'POST',
            dataType: 'json',
            data: {
				Producto: Prod,
				Cliente:Cli,
				CotizacionID:Cot,
            },
            success: function(data) {
                if ( data != null ) {
                    if ( data.SEstado ) {
                        if ( data.SElastic ) {
                            if ( data.SBusqueda != "" ) {
								$("#tabBusquedaProductos").show();
                                $("#BBusquedaProductos").empty();
                                $("#BBusquedaProductos").append(data.SBusqueda);
                                // $('#BuscaProductos').trigger("click");
                            } else {
                                alertify.error(data.SMsj);
                            }
                        } else {
                            $("#idbody").empty();
                            $("#idbody").append(data.SIhtml);
                            $("#calculadora").empty();
                            $("#calculadora").append(data.SCalculadora);
                        }
                    } else {
                        alertify.error(data.SMsj);
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                }
                $('#Loading').hide();
            },
            error: function(a, b, c) {
                alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
                $('#Loading').hide();
            },
        });
        $("#Buscar").val("");
        $("#Buscar").focus();
    } else {
        var validator = valida();
        validator.showErrors({
            "Buscar": "EL campo no debe estar vacío"
        });

        $("#Buscar").focus();
    }
}



function AplicaPeticion(clid){
	var id = $(clid).attr("id");
	var almacen = $(clid).attr("almacen");
	var Cot = $("#ID").val();
	var Cantidad = $(clid).val();

	var prev = $(clid).attr("previo");
	if (Number(prev) != Number(Cantidad)){
		if (id != "") {
			$('#Loading').show();
			$.ajax({
				url: "/Cotizacions/ActualizaProductoCarrito",
				type: 'POST',
				dataType: 'json',
				data: {
					Producto: id,
					Almacen: almacen,
					Operacion: Cot,
					Cantidad: Cantidad
				},
				success: function(data) {
					if (data != null) {
						if ( data.SEstado ) {
							$("#tabBusquedaProductos").hide();
							$("#BBusquedaProductos").empty();
							$("#BProductos").empty();
							$("#BProductos").append(data.SIhtml);
							$("#tabProductos").show();
							$("#calculadora").empty();
							$("#calculadora").append(data.SCalculadora);
							//Detalle
							$("#BProductosDetalle").empty();
							$("#BProductosDetalle").append(data.SIhtml);
							$("#calculadoraDetalle").empty();
							$("#calculadoraDetalle").append(data.SCalculadora);
							alertify.success(data.SMsj);
						} else {
							alertify.error(data.SMsj);
						}
					} else {
						alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
					}
					$('#Loading').hide();
				},
				error: function(a, b, c) {
					alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
					$('#Loading').hide();
				},
			});
		} else {
			alertify.error("Existio un problema con la referencia del Cliente");
		}
	}
	
}


function quitarProducto(clid){
	var id = $(clid).parent().attr("id");
	var almacen = $(clid).parent().attr("almacen");
	var Cot = $("#ID").val();
	if (id != "") {
        $('#Loading').show();
        $.ajax({
            url: "/Cotizacions/QuitarProducto",
            type: 'POST',
            dataType: 'json',
            data: {
				Producto: id,
				Almacen: almacen,
				Operacion: Cot
            },
            success: function(data) {
                if (data != null) {
                    if ( data.SEstado ) {
						$("#BProductos").empty();
						$("#BProductos").append(data.SIhtml);
						$("#tabProductos").show();
						$("#calculadora").empty();
						$("#calculadora").append(data.SCalculadora);
						//Detalle
						$("#BProductosDetalle").empty();
						$("#BProductosDetalle").append(data.SIhtml);
						$("#calculadoraDetalle").empty();
						$("#calculadoraDetalle").append(data.SCalculadora);
						alertify.success(data.SMsj);
					} else {
                        alertify.error(data.SMsj);
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                }
                $('#Loading').hide();
            },
            error: function(a, b, c) {
                alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
                $('#Loading').hide();
            },
        });
    } else {
        alertify.error("Existio un problema con la referencia del Cliente");
    }
}

function ActualizaDetalle(){
	if( $("#BClientes").children().length > 0 ){
		$("#NombreDetalle").val($("#BClientes").children().attr("Nombre"));
		$("#TelefonoDetalle").val($("#BClientes").children().attr("Telefono"));
		$("#CorreoDetalle").val($("#BClientes").children().attr("Correo"));
	}else{
		$("#NombreDetalle").val($("#Nombre").val());
		$("#TelefonoDetalle").val($("#Telefono").val());
		$("#CorreoDetalle").val($("#Correo").val());
	}

	$("#FormaDePagoDetalle").val($("#FormaDePago").val());
	$("#FormaDeEnvíoDetalle").val($("#FormaDeEnvío").val());	
}