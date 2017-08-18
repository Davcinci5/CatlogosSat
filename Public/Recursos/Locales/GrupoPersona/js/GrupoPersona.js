

	//##############################< SCRIPTS JS >##########################################
	//################################< GrupoPersona.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {			
		var validator = valida();			
	});
	function allowDrop(ev) {
		ev.preventDefault();
	}

	function drag(ev) {
		ev.dataTransfer.setData("text", ev.target.id);
	}

	function drop(ev) {
		ev.preventDefault();
		var data = ev.dataTransfer.getData("text");
		if(ev.target.id == "ingroup"){
			ev.target.appendChild(document.getElementById(data));
			$("#"+data+" :input").first().attr("name","Miembros");
		}else if(ev.target.id == "outgroup"){
			ev.target.appendChild(document.getElementById(data));
			$("#"+data+" :input").first().attr("name","PersonaGpo");			
		}		
	}

    function valida(){
	var validator = $("#Form_Alta_GrupoPersona").validate({
		rules: {
			
			Nombre : {
						
					required : true,
				
					rangelength : [5, 100]
				
					},
			Descripcion : {
						
					rangelength : [10, 250]
				
					},
		},
		messages: {
			
			Nombre : {
						
					rangelength : "La longitud del campo Nombre debe estar entre  [5, 100]",
					required : "El campo Nombre es requerido."
					},
			Descripcion : {
						
					rangelength : "La longitud del campo Descripcion debe estar entre  [10, 250]"
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


function EditaGrupoPersona(vista){
	if (vista == "Index" || vista ==""){
		if ($('#GrupoPersonas').val() != ""){
			window.location = '/GrupoPersonas/edita/' + $('#GrupoPersonas').val();
		}else{
			alertify.error("Debe Seleccionar un GrupoPersona para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/GrupoPersonas/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/GrupoPersonas';
		}
	}

}


function DetalleGrupoPersona(){
	if ($('#GrupoPersonas').val() != ""){
		window.location = '/GrupoPersonas/detalle/' + $('#GrupoPersonas').val();
	}else{
	alertify.error("Debe Seleccionar un GrupoPersona para editar");
	}
}


