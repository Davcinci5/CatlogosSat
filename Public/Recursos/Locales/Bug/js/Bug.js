

	//##############################< SCRIPTS JS >##########################################
	//################################< Bug.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################
	
	$( document ).ready( function () {	
		var validator = valida();			
	});

    function valida(){
	var validator = $("#Form_Alta_Bug").validate({
		rules: {
			
			Tipo : {
						
					rangelength : [5, 50],
				
					required : true
				
					},
			Titulo : {
						
					required : true,
				
					rangelength : [5, 100]
				
					},
			Descripcion : {
						
					required : true,
				
					rangelength : [10, 50]
				
					},
		},
		messages: {
			
			Tipo : {
						
					required : "El campo Tipo es requerido.",
					rangelength : "La longitud del campo Tipo debe estar entre  [5, 50]"
					},
			Titulo : {
						
					required : "El campo Titulo es requerido.",
					rangelength : "La longitud del campo Titulo debe estar entre  [5, 100]"
					},
			Descripcion : {
						
					required : "El campo Descripcion es requerido.",
					rangelength : "La longitud del campo Descripcion debe estar entre  [10, 50]"
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

function EditaBug(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Bugs').val() != ""){
			window.location = '/Bugs/edita/' + $('#Bugs').val();
		}else{
			alertify.error("Debe Seleccionar un Bug para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Bugs/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Bugs';
		}
	}

}


function DetalleBug(){
	if ($('#Bugs').val() != ""){
		window.location = '/Bugs/detalle/' + $('#Bugs').val();
	}else{
	alertify.error("Debe Seleccionar un Bug para editar");
	}
}



function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Bugs/search",
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
			url:"/Bugs/agrupa",
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


