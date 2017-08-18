

	//##############################< SCRIPTS JS >##########################################
	//################################< Usuario.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {			
		var validator = valida();		
		$('select[name=Cajas]').kendoMultiSelect().data("kendoMultiSelect");
		$('select[name=Grupos]').kendoMultiSelect().data("kendoMultiSelect");
		$('select[name=Tipo]').kendoMultiSelect().data("kendoMultiSelect");
		$("body").on("click",".eliminar", function(e){
			$(this).parent('div').parent('div').parent('div').remove();
		});
		$('input[type=radio][name=CorreosPrincipal]').change(function() {
			$("#CorreoPrincipal").val($('input[name=CorreosPrincipal]:checked', '#Form_Alta_Usuario').val());
		});
		$('input[type=radio][name=TelefonosPrincipal]').change(function() {
			$("#TelefonoPrincipal").val($('input[name=TelefonosPrincipal]:checked', '#Form_Alta_Usuario').val());
		});
		// $("[name='PrincipalCorreo']").change(function(){ $("[name='PrincipalCorreo']").val( $(this).parent().parent().find("#Correos").val())});
		// $("[name='PrincipalTelefono']").change(function(){ $("[name='PrincipalTelefono']").val( $(this).parent().parent().find("#Telefonos").val())});

		$('#AgregaEmail').click(function () {				
		if($('#Email').val() == ""){
			$("#Email").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Email": "No puede agregar campos vacíos"
			});
			
		}else{
			$("#Email").parent().removeClass("has-feedback has-error");	
			if($('#Email').valid()){
			$("#Email").parent().removeClass("has-feedback has-error");
			$('#div_tabla_correos').show();		
			$("#tbody_etiquetas_correos").append(
				'<tr>\n\
				<td><input type="radio" name="CorreosPrincipal" value="' + $("#Email").val() + '" checked></td>\n\
				<td><input type="text" class="form-control" name="Correos" value="' + $("#Email").val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');
				$('input[type=radio][name=CorreosPrincipal]').change(function() {
					$("#CorreoPrincipal").val($('input[name=CorreosPrincipal]:checked', '#Form_Alta_Usuario').val());
				});
				$("#CorreoPrincipal").val($("#Email").val());	
				$("#Email").val("");
				$("#Email").focus();			
			}else{
				$("#Email").parent().addClass("has-feedback has-error");
				$("#Email").focus();				
			}			
		}
		});	
	$('#AgregaTelefono').click(function () {				
		if($('#Telefono').val() == ""){
			$("#Telefono").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Telefono": "No puede agregar campos vacíos"
			});
			
		}else{
			$("#Telefono").parent().removeClass("has-feedback has-error");	
			if($('#Telefono').valid()){
			$("#Telefono").parent().removeClass("has-feedback has-error");
			$('#div_tabla_telefonos').show();		
			$("#tbody_etiquetas_telefonos").append(
				'<tr>\n\
				<td><input type="radio" name="TelefonosPrincipal" value="' + $("#Telefono").val() + '" checked></td>\n\
				<td><input type="text" class="form-control" name="Telefonos" value="' + $("#Telefono").val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');
				$('input[type=radio][name=TelefonosPrincipal]').change(function() {
					$("#TelefonoPrincipal").val($('input[name=TelefonosPrincipal]:checked', '#Form_Alta_Usuario').val());
				});
				$("#TelefonoPrincipal").val($("#Telefono").val());
				$("#Telefono").val("");
				$("#Telefono").focus();	

							
			}else{
				$("#Telefono").parent().addClass("has-feedback has-error");
				$("#Telefono").focus();				
			}			
		}
	});	
	
	$('#AgregaOtro').click(function () {				
		if($('#Otro').val() == ""){
			$("#Otro").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Otro": "No puede agregar campos vacíos"
			});
			
		}else{
			validator.showErrors({
				"Otro": ""
			});
			$("#Otro").parent().removeClass("has-feedback has-error");	
			$('#div_tabla_otros').show();		
			$("#tbody_etiquetas_otros").append(
				'<tr>\n\
				<td><input type="text" class="form-control" name="Otros" value="' + $("#Otro").val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');
				$("#Otro").val("");
				$("#Otro").focus();	
		}
	});	
	
	$('#Email').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
			e.preventDefault();
			$('#AgregaEmail').trigger("click");
		}
	});
	$('#Telefono').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
			e.preventDefault();
			$('#AgregaTelefono').trigger("click");
		}
	});
	$('#Otro').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
			e.preventDefault();
			$('#AgregaOtro').trigger("click");
		}
	});






		$(document).on('click', '.deleteButton', function () {
			$(this).parent().parent().remove();

			if (document.getElementById("tbody_etiquetas_correos").children.length == 0){
				$('#div_tabla_correos').hide();
			}
			if (document.getElementById("tbody_etiquetas_telefonos").children.length == 0){
				$('#div_tabla_telefonos').hide();
			}
			if (document.getElementById("tbody_etiquetas_otros").children.length == 0){
				$('#div_tabla_otros').hide();
			}

			// if (document.getElementById("tbody_etiquetas_etiquetas").children.length == 0){
			// 	$('#div_tabla_etiquetas').hide();
			// }
		});

	});



    function valida(){
	var validator = $("#Form_Alta_Usuario").validate({
		rules: {
			Nombre : {
						
					required : true,
				
					rangelength : [5, 100]
				
					},
			Tipo : {
						
					required : true
				
					},
			Grupos : {
						
					required : false
				
					},
			Predecesor : {
						
					required : false
				
					},
			Usuario : {
					required:true,
						
					rangelength : [5,20]
				
					},
			Contraseña : {

					required:true,
					minlength:8

					},
			confirmaContraseña : {
					required:true,
					minlength:8,
					equalTo:"#Contraseña"

					},
			Pin : {
					rangelength:[4,4]

					},
			confirmaPin : {
					rangelength:[4,4],
					equalTo:"#Pin"

					},
			Roles : {
						
					required : true
				
					},
			Principal : {
						
					required : true
		
					},
			 Email:{
			    	email:true

			       },
			 Telefono:{
			    	rangelength:[10,13]

			       },
		},
		messages: {
			
			Nombre : {
						
					required : "El campo Nombre es requerido.",
					rangelength : "La longitud del campo Nombre debe estar entre  [5, 100]"
					},
			Tipo : {
						
					required : "El campo Tipo es requerido."
					},
			Grupos : {
						
					required : "El campo Grupos es requerido."
					},
			Predecesor : {
						
					required : "Introduce un Predecesor si se ocupa"
					},
			Usuario : {
					required:"El nombre de usuario es requerido",						
					rangelength : "La longitud del campo Usuario debe estar entre  [5,20]"
					},
	   Contraseña : {

					required:"La contraseña es requerida",
					minlength:"Debe contener al menos 8 caracteres"

    				},
confirmaContraseña : {
					required:"La contraseña es requerida",
					minlength:"Debe contener al menos 8 caracteres",
					equalTo:"La contraseña no es igual"

					},
			  Pin : {
					rangelength:"Debe contener una longitud de 4"
					},
	  confirmaPin : {
					rangelength:"Debe contener una longitud de 4",
					equalTo:"Pin incorrecto"
					},
			Roles : {
					required : "El campo Roles es requerido."
					},
			  Email:{
		     		email:"Debe ser un email ejemplo@gmail.com"
			        },
			  Telefono:{
		     		rangelength:"Deben ser 14 digitos"
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

function AgregarContacto(){
	var TmpMedioExtra = '<div class="form-group">';
	TmpMedioExtra += '<div  class="col-sm-offset-4 col-sm-5">';
	TmpMedioExtra += '<div class="col-sm-5"> <div class="input-group"> <span class="input-group-addon">	<input type="radio" id="PrincipalCorreo" name="PrincipalCorreo" class="CheckPrincipalCorreo"> </span> <input type="text" name="Correos" id="Correos" class="form-control inputCorreos" placeholder="Correo Electronico" > </div> </div>';
	TmpMedioExtra += '<div class="col-sm-5"> <div class="input-group"> <span class="input-group-addon">	<input type="radio" id="PrincipalTelefono" name="PrincipalTelefono" class="CheckPrincipalTelefono"> </span> <input type="text" name="Telefonos" id="Telefonos" class="form-control inputTelefono" placeholder="Telefono"> </div> </div>';
	TmpMedioExtra += '<div class="col-sm-2"> <input type="button" class="btn btn-danger eliminar" value="- Eliminar" </div>';
	TmpMedioExtra += '</div>';
	TmpMedioExtra += '</div>';
	$("#MediosExtras").append(TmpMedioExtra);
	$("[name='PrincipalCorreo']").change(function(){ $("[name='PrincipalCorreo']").val( $(this).parent().parent().find("#Correos").val())});
	$("[name='PrincipalTelefono']").change(function(){ $("[name='PrincipalTelefono']").val( $(this).parent().parent().find("#Telefonos").val())});
}

function EditaUsuario(vista){
 if (vista == "Index" || vista ==""){
  if ($('#Usuarios').val() != ""){
   window.location = '/Usuarios/edita/' + $('#Usuarios').val();
  }else{
   alertify.error("Debe Seleccionar un Usuario para editar");
  }
 }else if(vista == "Detalle"){
  if ($('#ID').val() != ""){
   window.location = '/Usuarios/edita/' + $('#ID').val();
  }else{
   alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
   window.location = '/Usuarios';
  }
 }

}


function DetalleUsuario(){
 if ($('#Usuarios').val() != ""){
  window.location = '/Usuarios/detalle/' + $('#Usuarios').val();
 }else{
 alertify.error("Debe Seleccionar un Usuario para editar");
 }
}

function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Usuarios/search",
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
			url:"/Usuarios/agrupa",
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





