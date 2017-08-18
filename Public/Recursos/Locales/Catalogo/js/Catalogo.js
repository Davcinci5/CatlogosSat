

	//##############################< SCRIPTS JS >##########################################
	//################################< Catalogo.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################

$( document ).ready( function () {
	if (document.getElementById("tbody_etiquetas_catalogo").children.length == 0){
		$('#div_tabla_catalogo').hide();
	}	
	
	var validator = valida();		
	
	$('#AgregaCampo').click(function () {
		if ($('#Valor').val() != ""){
			$('#div_tabla_catalogo').show();
			$("#tbody_etiquetas_catalogo").append(
			'<tr>\n\
				<td><input type="hidden" class="form-control" name="ValoresIds" value=""><input type="text" class="form-control" onblur="ValidaCampo(this)"  id="inputTX" name="Valores" value="' + $("#Valor").val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-primary editButton"><span class="glyphicon glyphicon-pencil btn-xs"></span></button><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
			</tr>');

             $("#Valor").val("");
		}else{
			validator.showErrors({
			"Valor": "No puede agregar valores vacíos"
			});
			 $("#Valor").focus();
		}
	});	

	$('#Form_Alta_Catalogo').keydown(function(e) {
		if(e.which == 13 || e.keyCode == 13) {
			if ($('#Valor').val() != ""){
				e.preventDefault();
				$('#AgregaCampo').trigger("click");
				validator.element("#Valor");
			}else{
				if (document.getElementById("tbody_etiquetas_catalogo").children.length == 0){
					e.preventDefault();
					validator.showErrors({
						"Valor": "No puede dar de alta Catálogos sin valores"
					});
					
				}else if($("#Nombre")==''){
					validator.showErrors({
						"Nombre": "El Nombre es requerido"
					});
					$("#Nombre").focus();
				}
			}  
		}
	});


});

$(document).on('click', '.deleteButton', function () {
	$(this).parent().parent().remove();
	if (document.getElementById("tbody_etiquetas_catalogo").children.length == 0){
		$('#div_tabla_catalogo').hide();
	}
});

$(document).on('click', '.editButton', function () {
	$(this).parent().parent().children().children()[1].readOnly = false;
	$(this).parent().parent().children().children()[1].focus();
});

	function valida(){
		var validator =	$( "#Form_Alta_Catalogo" ).validate({
		rules: {			
			Nombre : {						
					required : true,				
					rangelength : [5, 100]				
					},
			Valor : {
					rangelength : [0, 100]					
					},
			Descripcion : {
					rangelength : [20, 250]
					}
				},
		messages: {
			
			Nombre : {						
					required : "El campo Nombre es requerido.",
					rangelength : "La longitud del campo Nombre debe estar entre 5 y 100"
					},
			Valor : {							
					rangelength : "La longitud del campo Valor debe estar entre 1 y 100"
				},
			Descripcion : {
						
					rangelength : "La longitud del campo Descripcion debe estar entre  [20, 250]"
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

function EditaCatalogo(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Catalogos').val() != ""){
			window.location = '/Catalogos/edita/' + $('#Catalogos').val();
		}else{
			alertify.error("Debe Seleccionar un Catálogo para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Catalogos/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Catalogos';
		}
	}

}


function DetalleCatalogo(){
	if ($('#Catalogos').val() != ""){
		window.location = '/Catalogos/detalle/' + $('#Catalogos').val();
	}else{
	alertify.error("Debe Seleccionar un Catálogo para editar");
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
		input.value = "--SIN VALOR--> (ESTATUS INACTIVO O CANCELADO) Valor Provicional mientras se adaptan los estatus de los valores de cada catálogo)."
		alertify.error("El Campo No debe ir vacío.");
		input.readOnly = true;		
	}
}

function BuscaPagina(num){
			$('#Loading').show();

			$.ajax({
			url:"/Catalogos/search",
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
			url:"/Catalogos/agrupa",
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


