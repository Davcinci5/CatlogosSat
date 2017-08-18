

	//##############################< SCRIPTS JS >##########################################
	//################################< Almacen.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {			
		var validator = valida();
		$( "#Tipo" ).change(function() {
			if ($( "#Tipo" ).val() == "58e56906e75770120c60bef2"){
				$('#Conexion').removeAttr("disabled");
			}else{
				$('#Conexion').attr("disabled", "disabled"); 
			}
		});
	});

    function valida(){
	var validator = $("#Form_Alta_Almacen").validate({
		rules: {
			
			Nombre : {
						
					required : true,
				
					rangelength : [5, 100]
				
					},
			Tipo : {
						
					required : true
				
					},
			Clasificacion : {
						
					required : true
				
					},
			//ListaCosto : {
						
			//		required : true
				
			//		},
			//ListaPrecio : {
						
			//		required : true
				
			//		},
						Estatus:{
							required : true
						},
						Servidor : {
								
						required : true
					
							},
						NombreBD : {
								
						required : true
					
							},
						UsuarioBD : {
								
						required : true
					
							},
						PassBD : {
								
						required : true
					
							}
		},
		messages: {
			
			Nombre : {
						
					required : "El campo Nombre es requerido.",
					rangelength : "La longitud del campo Nombre debe estar entre  [5, 100]"
					},
			Tipo : {
						
					required : "El campo Tipo es requerido."
					},
			Clasificacion : {
						
					required : "El campo Clasificacion es requerido."
					},
			ListaCosto : {
						
					required : "El campo ListaCosto es requerido."
					},
			ListaPrecio : {
						
					required : "El campo ListaPrecio es requerido."
					},
						Servidor : {
								
						required : "El campo Servidor es requerido."
							},
						NombreBD : {
								
						required : "El campo NombreBD es requerido."
							},
						UsuarioBD : {
								
						required : "El campo UsuarioBD es requerido."
							},
						PassBD : {
								
						required : "El campo PassBD es requerido."
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


function EditaAlmacen(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Almacens').val() != ""){
			window.location = '/Almacens/edita/' + $('#Almacens').val();
		}else{
			alertify.error("Debe Seleccionar un Almacen para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Almacens/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Almacens';
		}
	}

}


function DetalleAlmacen(){
	if ($('#Almacens').val() != ""){
		window.location = '/Almacens/detalle/' + $('#Almacens').val();
	}else{
	alertify.error("Debe Seleccionar un Almacen para editar");
	}
}

function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Almacens/search",
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
			url:"/Almacens/agrupa",
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


