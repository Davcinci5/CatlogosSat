

	//##############################< SCRIPTS JS >##########################################
	//################################< Facturacion.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {
		$("#Operacion").focus();

		
		var validator = valida();			
	});

    function valida(){
	var validator = $("#Form_Alta_Facturacion").validate({
		rules: {
			
						Calle : {
								
						rangelength : [4, 100],
					
						required : true
					
							},
						NumInterior : {
								
						rangelength : [1, 5],
					
						required : true
					
							},
						NumExterior : {
								
						required : true,
					
						rangelength : [1, 5]
					
							},
						Colonia : {
								
						required : true,
					
						rangelength : [1, 50]
					
							},
						Municipio : {
								
						required : true,
					
						rangelength : [4, 100]
					
							},
						Estado : {
								
						required : true,
					
						rangelength : [4, 50]
					
							},
						Pais : {
								
						required : true,
					
						rangelength : [4, 30]
					
							},
						CP : {
								
						required : true,
					
						rangelength : [5, 5]
					
							},
						Estatus : {
								
							}
		},
		messages: {
			
						Calle : {
								
						required : "El campo Calle es requerido.",
						rangelength : "La longitud del campo Calle debe estar entre  [4, 100]"
							},
						NumInterior : {
								
						required : "El campo NumInterior es requerido.",
						rangelength : "La longitud del campo NumInterior debe estar entre  [1, 5]"
							},
						NumExterior : {
								
						required : "El campo NumExterior es requerido.",
						rangelength : "La longitud del campo NumExterior debe estar entre  [1, 5]"
							},
						Colonia : {
								
						required : "El campo Colonia es requerido.",
						rangelength : "La longitud del campo Colonia debe estar entre  [1, 50]"
							},
						Municipio : {
								
						required : "El campo Municipio es requerido.",
						rangelength : "La longitud del campo Municipio debe estar entre  [4, 100]"
							},
						Estado : {
								
						rangelength : "La longitud del campo Estado debe estar entre  [4, 50]",
						required : "El campo Estado es requerido."
							},
						Pais : {
								
						required : "El campo Pais es requerido.",
						rangelength : "La longitud del campo Pais debe estar entre  [4, 30]"
							},
						CP : {
								
						required : "El campo CP es requerido.",
						rangelength : "La longitud del campo CP debe estar entre  [5, 5]"
							},
						Estatus : {
								
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


function EditaFacturacion(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Facturacions').val() != ""){
			window.location = '/Facturacions/edita/' + $('#Facturacions').val();
		}else{
			alertify.error("Debe Seleccionar un Facturacion para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Facturacions/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Facturacions';
		}
	}

}


function DetalleFacturacion(){
	if ($('#Facturacions').val() != ""){
		window.location = '/Facturacions/detalle/' + $('#Facturacions').val();
	}else{
	alertify.error("Debe Seleccionar un Facturacion para editar");
	}
}


