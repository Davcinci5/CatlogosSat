

	//##############################< SCRIPTS JS >##########################################
	//################################< Kit.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {			
		var validator = valida();			
	});

    function valida(){
	var validator = $("#Form_Alta_Kit").validate({
		rules: {
			
			Nombre : {
						
					required : true,
				
					rangelength : [5, 150]
				
					},
			Tipo : {
						
					required : true
				
					},
			Aplicación : {
						
					required : true
				
					},
						Almacen : {
									
						required : true
					
								},
						IDProducto : {
								
						required : true
					
							},
						Cantidad : {
								
						required : true
					
							},
						Precio : {
								
						required : true
					
							}
		},
		messages: {
			
			Nombre : {
						
					required : "El campo Nombre es requerido.",
					rangelength : "La longitud del campo Nombre debe estar entre  [5, 150]"
					},
			Tipo : {
						
					required : "El campo Tipo es requerido."
					},
			Aplicación : {
						
					required : "El campo Aplicación es requerido."
					},
						Almacen : {
									
						required : "El campo Almacen es requerido."
								},
						IDProducto : {
								
						required : "El campo IDProducto es requerido."
							},
						Cantidad : {
								
						required : "El campo Cantidad es requerido."
							},
						Precio : {
								
						required : "El campo Precio es requerido."
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


function EditaKit(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Kits').val() != ""){
			window.location = '/Kits/edita/' + $('#Kits').val();
		}else{
			alertify.error("Debe Seleccionar un Kit para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Kits/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Kits';
		}
	}

}


function DetalleKit(){
	if ($('#Kits').val() != ""){
		window.location = '/Kits/detalle/' + $('#Kits').val();
	}else{
	alertify.error("Debe Seleccionar un Kit para editar");
	}
}


