

	//##############################< SCRIPTS JS >##########################################
	//################################< Unidad.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {
		if (document.getElementById("tbody_etiquetas_unidades").children.length == 0){
			$('#div_tabla_unidades').hide();
		}	


		var validatorModal = validaModal();
		var validator = valida();		
			

$('#Magnitud').change(function(){
	var magnitud = $('#Magnitud').val();

	if  (magnitud != ""){
		$.ajax({
			url: '/ConsultaUnidadesDeMagnitud',
			type: 'POST',
			dataType: 'html',
			data: { magnitud : magnitud},
			success: function (data) {
				$("#tbody_etiquetas_unidades").empty();
				$("#tbody_etiquetas_unidades").append(data);
			},
			error: function (data) {
				alertify.error("Ocurrió un error al momento de consultar las unidades de la magnitud, por favor intente de nuevo.")
			}
		});
	}
});



		$('#AgregaCampo').click(function () {
			if ($('#Nombre').val() == ""){
				validator.showErrors({
				"Nombre": "No puede agregar Nombres vacíos"
				});
				$("#Nombre").focus();
			}else if($('#Abreviatura').val() == ""){
				validator.showErrors({
				"Abreviatura": "No puede agregar Abreviaturas vacías"
				});
				$("#Abreviatura").focus();
			}else{
				$('#div_tabla_unidades').show();
				$("#tbody_etiquetas_unidades").append(
				'<tr>\n\
				<td><input type="hidden" class="form-control" name="DatosIds" value=""><input type="text" class="form-control" name="Nombres" value="' + $("#Nombre").val() + '" readonly></td>\n\
				<td><input type="text" class="form-control" name="Abreviaturas" value="' + $("#Abreviatura").val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button> <button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');

				$("#Nombre,#Abreviatura").val("");
				$("#Nombre").focus();
			}
		});	

	$('#Form_Alta_Unidad').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
				$('#AgregaCampo').trigger("click");
				// validator.element("#Nombre");
				// validator.element("#Abreviatura");
		}
	});


	});

	$(document).on('click', '.deleteButton', function () {
		$(this).parent().parent().remove();

		if (document.getElementById("tbody_etiquetas_unidades").children.length == 0){
			$('#div_tabla_unidades').hide();
		}
	});

	$(document).on('click', '.editButton', function () {
		$(this).parent().parent().children().children()[1].readOnly = false;
		$(this).parent().parent().children().children()[2].readOnly = false;
		$(this).parent().parent().children().children()[1].focus();
	});

    function valida(){
	var validator = $("#Form_Alta_Unidad").validate({
		rules: {
			
			Magnitud : {
						
					required : true
				
					},
			Descripcion : {
									
					rangelength : [20, 250]				
					},
			Nombre : {	

					rangelength : [4, 30]
					},
			Abreviatura : {
							
					rangelength : [1, 10]					
				}
				},
		messages: {
			
			Magnitud : {
						
					required : "El campo Magnitud es requerido."
					},
			Descripcion : {
						
					rangelength : "La longitud del campo Descripcion debe estar entre  [20, 250]"
					},
			Nombre : {

					rangelength : "La longitud del campo Nombre debe estar entre  [4, 30]"
					},
			Abreviatura : {
								
					rangelength : "La longitud del campo Abreviatura debe estar entre  [1, 10]"
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

function validaModal(){
	var validatorM = $("#Alta_Magnitud_Modal").validate({
		rules: {
			ModalMagnitud:{
				rangelength : [1, 25],
				required : true
			}
				},
		messages: {
			ModalMagnitud : {
								
					rangelength : "La longitud de la longitud debe estar entre  [1, 25]",
						required : "El campo Nombre es requerido."
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
	return validatorM;
}


function EditaUnidad(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Unidads').val() != ""){
			window.location = '/Unidades/edita/' + $('#Unidads').val();
		}else{
			alertify.error("Debe Seleccionar un Unidad para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Unidades/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Unidades';
		}
	}

}

function DetalleUnidad(){
	if ($('#Unidads').val() != ""){
		window.location = '/Unidades/detalle/' + $('#Unidads').val();
	}else{
	alertify.error("Debe Seleccionar un Unidad para ver Detalle");
	}
}

function ValidaCampo(input){
	if (input.value == ""){
		alertify.error("El Campo No debe ir vacío.");
		input.parentElement.parentElement.remove();
		if (document.getElementById("tbody_etiquetas_catalogo").children.length == 0){
			$('#div_tabla_catalogo').hide();
		}
	}else{
		input.readOnly = true;
	}
}

function ValidaCampo2(input){
	if (input.value == ""){
		input.value = "--SIN VALOR--"
		alertify.error("El Campo No debe ir vacío.");
		input.readOnly = true;		
	}
}
