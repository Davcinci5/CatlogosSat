

	//##############################< SCRIPTS JS >##########################################
	//################################< Producto.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

	$( document ).ready( function () {
		if (document.getElementById("tbody_etiquetas_codigos").children.length == 0){
			$('#div_tabla_codigos').hide();
		}	

		if (document.getElementById("tbody_etiquetas_etiquetas").children.length == 0){
			$('#div_tabla_etiquetas').hide();
		}				
		var validator = valida();	

		$('#AgregaCampo').click(function () {
			if ($('#Codigo').val() == ""){
				validator.showErrors({
				"Codigo": "No puede agregar Codigos vacíos"
				});
				$("#Codigo").focus();
			}else if($('#Valcodigo').val() == ""){
				validator.showErrors({
				"Valcodigo": "No puede agregar valores vacías"
				});
				$("#Valcodigo").focus();
			}else{
				$('#div_tabla_codigos').show();
				$("#tbody_etiquetas_codigos").append(
				'<tr>\n\
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Codigos" value="' + $("#Codigo").val() + '" readonly></td>\n\
				<td><input type="text" class="form-control" name="Valcodigos" value="' + $("#Valcodigo").val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');

				$("#Codigo,#Valcodigo").val("");
				$("#Codigo").focus();
			}
		});	

		$('#AgregaEtiqueta').click(function () {
			if ($('#Etiqueta').val() == ""){
				validator.showErrors({
				"Etiqueta": "No puede agregar Etiquetas vacías"
				});
				$("#Etiqueta").focus();
			}else{
				$('#div_tabla_etiquetas').show();
				$("#tbody_etiquetas_etiquetas").append(
				'<tr>\n\
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="Etiquetas" value="' + $("#Etiqueta").val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');

				$("#Etiqueta").val("");
				$("#Etiqueta").focus();
			}
		});		
	

		$('#Codigo,#Valcodigo').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
				$('#AgregaCampo').trigger("click");
				// validator.element("#Codigo");
				// validator.element("#Valcodigo");

		}
	});	
	$('#Etiqueta').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
				e.preventDefault();
				$('#AgregaEtiqueta').trigger("click");
				// validator.element("#Etiqueta");

		}
	});	

	});

	$(document).on('click', '.deleteButton', function () {
		$(this).parent().parent().remove();

		if (document.getElementById("tbody_etiquetas_codigos").children.length == 0){
			$('#div_tabla_codigos').hide();
		}
		if (document.getElementById("tbody_etiquetas_etiquetas").children.length == 0){
			$('#div_tabla_etiquetas').hide();
		}
	});

	
    function valida(){
	var validator = $("#Form_Alta_Producto").validate({
		rules: {
			
			Nombre : {
						
					required : true,
				
					rangelength : [5, 150]
				
					},
			Tipo : {
						
					required : true
				
					},
			Unidad : {
						
					required : true
				
					},
						Claves : {
								
						required : true,
					
						rangelength : [5, 50]
					
							},
						Valores : {
								
						rangelength : [5, 25],
					
						required : true
					
							}
		},
		messages: {
			
			Nombre : {
						
					required : "El campo Descripción es requerido.",
					rangelength : "La longitud del campo Descripción debe estar entre  [5, 150]"
					},
			Tipo : {
						
					required : "El campo Tipo es requerido."
					},
			Unidad : {
						
					required : "El campo Unidad es requerido."
					},
						Claves : {
								
						required : "El campo Claves es requerido.",
						rangelength : "La longitud del campo Claves debe estar entre  [5, 50]"
							},
						Valores : {
								
						required : "El campo Valores es requerido.",
						rangelength : "La longitud del campo Valores debe estar entre  [5, 25]"
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


function EditaProducto(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Productos').val() != ""){
			window.location = '/Productos/edita/' + $('#Productos').val();
		}else{
			alertify.error("Debe Seleccionar un Producto para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Productos/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Productos';
		}
	}

}


function DetalleProducto(){
	if ($('#Productos').val() != ""){
		window.location = '/Productos/detalle/' + $('#Productos').val();
	}else{
	alertify.error("Debe Seleccionar un Producto para editar");
	}
}

function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Productos/search",
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
			url:"/Productos/agrupa",
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


