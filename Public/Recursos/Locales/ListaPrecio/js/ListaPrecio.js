

	//##############################< SCRIPTS JS >##########################################
	//################################< ListaPrecio.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {			
		var validator = valida();			
	});

    function valida(){
	var validator = $("#Form_Alta_ListaPrecio").validate({
		rules: {
			
			Nombre : {
						
					required : true,
				
					rangelength : [5, 100]
				
					},
			Descripcion : {
						
					rangelength : [20, 250]
				
					},
			GrupoP : {
						
					required : true
				
					},
		},
		messages: {
			
			Nombre : {
						
					required : "El campo Nombre es requerido.",
					rangelength : "La longitud del campo Nombre debe estar entre  [5, 100]"
					},
			Descripcion : {
						
					rangelength : "La longitud del campo Descripcion debe estar entre  [20, 250]"
					},
			GrupoP : {
						
					required : "El campo GrupoP es requerido."
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


function EditaListaPrecio(vista){
	if (vista == "Index" || vista ==""){
		if ($('#ListaPrecios').val() != ""){
			window.location = '/ListaPrecios/edita/' + $('#ListaPrecios').val();
		}else{
			alertify.error("Debe Seleccionar un ListaPrecio para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/ListaPrecios/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/ListaPrecios';
		}
	}

}


function DetalleListaPrecio(){
	if ($('#ListaPrecios').val() != ""){
		window.location = '/ListaPrecios/detalle/' + $('#ListaPrecios').val();
	}else{
	alertify.error("Debe Seleccionar un ListaPrecio para editar");
	}
}


