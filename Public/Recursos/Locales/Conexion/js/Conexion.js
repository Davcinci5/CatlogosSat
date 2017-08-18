

	//##############################< SCRIPTS JS >##########################################
	//################################< Conexion.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################
	
	$( document ).ready( function () {	
		var validator = valida();			
	});

    function valida(){
	var validator = $("#Form_Alta_Conexion").validate({
		rules: {
			
			Nombre : {
						
					rangelength : [5, 50],
				
					required : true
				
					},
			Servidor : {
						
					required : true,
				
					rangelength : [5, 100]
				
					},
			NombreBD : {
						
					required : true,
				
					rangelength : [1, 50]
				
					},
		},
		messages: {
			
			Nombre : {
						
					required : "El campo Nombre es requerido.",
					rangelength : "La longitud del campo Nombre debe estar entre  [5, 50]"
					},
			Servidor : {
						
					required : "El campo Servidor es requerido.",
					rangelength : "La longitud del campo Servidor debe estar entre  [5, 100]"
					},
			NombreBD : {
						
					rangelength : "La longitud del campo NombreBD debe estar entre  [1, 50]",
					required : "El campo NombreBD es requerido."
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

function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Conexions/search",
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

function testConexion(){
	Servidor =$("#Servidor").val();
	NombreBD = $("#NombreBD").val();
	UsuarioBD = $("#UsuarioBD").val();
	PassBD = $("#PassBD").val();
	 $.ajax({
  		type: "POST",
  		url: '/Conexions/testConexion',
		data : { "Servidor":Servidor,"NombreBD":NombreBD,"UsuarioBD":UsuarioBD,"PassBD":PassBD },
  		dataType: "json",
  		async: true,
		beforeSend :function(){
			 $('#Loading').show();
		  },
  		success: function (data) {
			  if(data.Estatus){
				 $( "#alertConexion" ).removeClass("alert alert-danger" )
				 $("#alertConexion").addClass('alert alert-success');
				 $("#iconConexion").html("<span class='glyphicon glyphicon-ok' ></span>");
				 $("#mensajeConexion").html("<h4> "+data.Mensage+"</h4>");
			  }else{  
				 $("#alertConexion" ).removeClass("alert alert-success" )
				 $("#alertConexion").addClass('alert alert-danger');
				 $("#iconConexion").html("<span class='glyphicon glyphicon-remove' ></span>");
				 $("#mensajeConexion").html("<h4> "+data.Mensage+"</h4><h4>"+data.MensageError+"</h4>");
			  }
  		},
  		error: function () {
   			alert('Error occured');
  		},
		complete:function(){
			$('#Loading').hide();
		}
 	});
}

 function SubmitGroup(){
	 $('#Loading').show();
			$.ajax({
			url:"/Conexions/agrupa",
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


