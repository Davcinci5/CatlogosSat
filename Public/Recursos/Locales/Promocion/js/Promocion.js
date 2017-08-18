

	//##############################< SCRIPTS JS >##########################################
	//################################< Promocion.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################
	
	$( document ).ready( function () {	
		var validator = valida();			
	});

    function valida(){
	var validator = $("#Form_Alta_Promocion").validate({
		rules: {
			
			Nombre : {
						
					required : true
				
					},
			Descripcion : {
						
					required : true
				
					},
			FechaInicio : {
						
					required : true
				
					},
			FechaFin : {
						
					required : true
				
					},
		},
		messages: {
			
			Nombre : {
						
					required : "El campo Nombre es requerido."
					},
			Descripcion : {
						
					required : "El campo Descripcion es requerido."
					},
			FechaInicio : {
						
					required : "El campo FechaInicio es requerido."
					},
			FechaFin : {
						
					required : "El campo FechaFin es requerido."
					},
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

function EditaPromocion(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Promocions').val() != ""){
			window.location = '/Promocions/edita/' + $('#Promocions').val();
		}else{
			alertify.error("Debe Seleccionar un Promocion para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Promocions/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Promocions';
		}
	}

}


function DetallePromocion(){
	if ($('#Promocions').val() != ""){
		window.location = '/Promocions/detalle/' + $('#Promocions').val();
	}else{
	alertify.error("Debe Seleccionar un Promocion para editar");
	}
}



function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Promocions/search",
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
			url:"/Promocions/agrupa",
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


