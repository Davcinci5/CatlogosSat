

	//##############################< SCRIPTS JS >##########################################
	//################################< Caja.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {			
		var validator = valida();	
		//Funcion para el drag and drop de medios de pago.
		$(function () {		
			$(".source, .compra, .venta").sortable({
			connectWith: ".connected"
			});
		});		
		//Convierte en tabla las operaciones
		$('#tabla').bootstrapTable({});
	});

    function valida(){
	var validator = $("#Form_Alta_Caja").validate({
		rules: {
			
		},
		messages: {
			
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


function EditaCaja(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Cajas').val() != ""){
			window.location = '/Cajas/edita/' + $('#Cajas').val();
		}else{
			alertify.error("Debe Seleccionar un Caja para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Cajas/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Cajas';
		}
	}

}


function DetalleCaja(){
	if ($('#Cajas').val() != ""){
		window.location = '/Cajas/detalle/' + $('#Cajas').val();
	}else{
	alertify.error("Debe Seleccionar un Caja para editar");
	}
}
/*****************/
/*****************/
function buscaDocumento(id){
	//Leo ID del documento a buscar.
	if (id == null){
		var id = $("#documentoID").val();
	}else{
		var id = id;
	}
	//Valido si viene vacio
	if (id == ""){
		alertify.error("Debes indicar un ID de documento.");
		//Regreso el focus al campo de captura de ID de documento.
		$("#documentoID").focus();		
	}else{			
		$.ajax({
			url: '/Cajas/buscaDocumento/',
			type: 'POST',
			dataType: 'html',
			data:{ID:id},
			success : function(data){
				$('#detalleDocumento').html(data);
			}
		});
		
	}
}
function allowDrop(ev) {
    ev.preventDefault();
}

function drag(ev) {
    ev.dataTransfer.setData("text", ev.target.id);
}

function drop(ev) {
    ev.preventDefault();
    var data = ev.dataTransfer.getData("text");
    ev.target.appendChild(document.getElementById(data));
}
/*
 * @melchor
 * Funcion que permite el drag and drop de las series
 * @param:null
 * @return: null
 */
$(function () {
    $(".source, .compra, .venta").sortable({
      connectWith: ".connected"
    });
});
//Valida monto del pago
function validaMonto(ido,aceptacambio){
	var importeDocumento = $("#importeDocumento").val();	
	var suma = 0;
	$("ul.compra li input").each(function() {
		var id = $(this).attr("id");
		var valor = $("#"+id).val();
		//var suma = valor;
		if(valor == 'NaN' || valor == '') {valor = 0;}
		
		suma = parseFloat(suma) + parseFloat(valor);
		//alert(valor+'>>'+suma);
		if (parseFloat(valor) >= parseFloat(importeDocumento)){
			$(".fpagos:not(#"+ido+")").prop('disabled', true);
			$(".fpagos:not(#"+ido+")").prop('value', 0);
			if(aceptacambio == 'true'){
				$("#cambio").attr("value",parseFloat(suma)-parseFloat(importeDocumento))
			}
			
			return false
		}else{			
			if(parseFloat(suma) >= parseFloat(importeDocumento)){
				$(".fpagos:not(#"+ido+")").prop('disabled', true);
				if(aceptacambio == 'true'){
					$("#cambio").attr("value",parseFloat(suma)-parseFloat(importeDocumento))
				}
				return false	
			}else{
				$(".fpagos").prop('disabled', false);
				$("#cambio").attr("value",'Su cambio')
			}
		}
	});
}
//Inserta documento
function insertaDocumento(){
//Lectura de parametros
var importeDocumento = $("#importeDocumento").val();
var cajax = $("#cajax").val();
var operacion = $("#documentoIDx").val();
var fp = 0;
var suma = 0;

$("ul.compra li input").each(function() {
		var id = $(this).attr("id");
		var valor = $("#"+id).val();
		suma = parseFloat(suma) + parseFloat(valor);
		if (valor > 0 && suma>=importeDocumento){
			fp = 1;
		}
});
if (fp == 0){
		alertify.error("Debe seleccionar al menos una forma de pago igual o mayor al monto del documento.");
		//alert("Debe seleccionar al menos una forma de pago igual o mayor al monto del documento.");
	}else{
		if(confirm("¿Seguro?"))
                        {					
					$.ajax({
						url: '/Cajas/insertaDocumento/',
						type: 'POST',
						dataType: 'html',
						data:{cargo:suma,cajax:cajax,operacion:operacion},
						success : function(data){
							$('#detalleDocumento').html(data);
						}
					});
						}
	}
}
//Trae el saldo del dia para mostrar en el cierre de caja.
function getSaldos(cajaID){
	$("#myModal").modal();
	//Voy por el saldo
						$.ajax({
							url: '/Cajas/getSaldo/',
							type: 'POST',
							dataType: 'html',
							data:{caja:cajaID},
							success : function(data){
								$('#contenidoModal').html(data);						
							}
						});
}
//Inserta documento
function cierraCaja(cajaID){
	if(confirm("¿Seguro?"))
                {
					$.ajax({
						url: '/Cajas/cierraCaja/',
						type: 'POST',
						dataType: 'html',
						data:{caja:cajaID},
						success : function(data){
							alertify.success("Cierre de caja OK");
							//$("#myModal").modal('hide')
							$(".modal-body").html('<iframe src="../PDF/'+cajaID+'.pdf" width="100%" height="100%"></iframe>');
						}
					});
				}
}
//Abre Cajas
function showPDF(ID){
	$("#myModal").modal('show')
	$(".modal-body").html('<iframe src="../PDF/'+ID+'.pdf" width="100%" height="100%"></iframe>');
}


function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Cajas/search",
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
			url:"/Cajas/agrupa",
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



