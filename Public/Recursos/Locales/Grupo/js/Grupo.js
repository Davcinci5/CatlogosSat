	//##############################< SCRIPTS JS >##########################################
	//################################< Grupo.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################
	
	$( document ).ready( function () {

		$(window).keydown(function(event){
			if(event.keyCode == 13) {
				event.preventDefault();
			}
		});

		var validator = valida();	
		
		// $('select[name=Tipo]').kendoMultiSelect().data("kendoMultiSelect");

		$('#inputBuscaBase').keydown(function(e) {
			if(e.which == 13 || e.keyCode == 13) {
				e.preventDefault();							
				ConsultaBase();
			}
        });		

	});

//Funciones para el manejo de los checkbox en las vistas
//var SkuSeleccionados = [];
var SkuSeleccionados = new Array();
var ClaveSatSeleccionado;

function GuardaSeleccionados(seleccionado, valor){
    if (seleccionado==true){
        SkuSeleccionados.push(valor);
    }else{
        var indiceEncontrado = SkuSeleccionados.indexOf(valor);
        removeA(SkuSeleccionados, valor);
    }	
}

/**
 * Funcion que seleciona o deselecciona  los checbox de productos
 * @param {*} seleccionado  El checbox de todos
 **/
function SelecionarTodos(seleccionado){

	if (seleccionado==true){
        $( ".ProdExtraido" ).each(function() {
			var id=$(this).attr( "id" );
  			$(this).prop('checked', true);
			  SkuSeleccionados.push(id);
		});
    }else{
		$( ".ProdExtraido" ).each(function() {
			var id=$(this).attr( "id");
  			$(this).prop('checked', false);
			  removeA(SkuSeleccionados, id);
		});
    }	
}


function ClickSelecionarTodos(seleccionado){
	if (seleccionado==true){
        SkuSeleccionados.push(valor);
    }else{
        var indiceEncontrado = SkuSeleccionados.indexOf(valor);
        removeA(SkuSeleccionados, valor);
    }	
}

function Verificarseleccionados(){
    for (var x=0; x<SkuSeleccionados.length; x++){
		if(document.getElementById(SkuSeleccionados[x])) {
			// document.getElementById(SkuSeleccionados[x]).setAttribute("checked", "checked");
		}
    }
	var contador=1;
	$( ".ProdExtraido" ).each(function() {
		if( !$(this).is(':checked') ) {
    		contador=0;
			return false;
		}
	});

	if(contador==1){
		// document.getElementById("selectTodos").setAttribute("checked", "checked");
	}	
}

function removeA(arr) {
    var what, a = arguments, L = a.length, ax;
    while (L > 1 && arr.length) {
        what = a[--L];
        while ((ax= arr.indexOf(what)) !== -1) {
            arr.splice(ax, 1);
        }
    }
    return arr;
}

function AsignarClaveSat(valor){
	ClaveSatSeleccionado = valor;
}

function AgregarSkuSeleccionadosAEnviar(){
	if (SkuSeleccionados.length>0){
		$("#SkuSeleccionados").val(JSON.stringify(SkuSeleccionados));
		return true;
	}else{
		return false;
	}
}

    function valida(){
	var validator = $("#Form_Alta_Grupo").validate({
		rules: {
			
			Nombre : {					
					required : true,				
					rangelength : [5, 100]				
					},
			Descripcion : {
					required : true,
					rangelength : [10, 250]				
					},
			PermiteVender : {						
					rangelength : [10, 250]				
					},
			Tipo : {						
					required : true				
					},
		},
		messages: {			
			Nombre : {						
					required : "El campo Nombre es requerido.",
					rangelength : "La longitud del campo Nombre debe estar entre  [5, 100]"
					},
			Descripcion : {
					required : "El campo Descripción es requerido.",
					rangelength : "La longitud del campo Descripcion debe estar entre  [10, 250]"
					},
			PermiteVender : {						
					rangelength : "La longitud del campo PermiteVender debe estar entre  [10, 250]"
					},
			Tipo : {						
					required : "El campo Tipo es requerido."
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

function EditaGrupo(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Grupos').val() != ""){
			window.location = '/Grupos/edita/' + $('#Grupos').val();
		}else{
			alertify.error("Debe Seleccionar un Grupo para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Grupos/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Grupos';
		}
	}

}

function DetalleGrupo(){
	if ($('#Grupos').val() != ""){
		window.location = '/Grupos/detalle/' + $('#Grupos').val();
	}else{
	alertify.error("Debe Seleccionar un Grupo para editar");
	}
}

function ConsultaBase(){
    var cadena = $("#inputBuscaBase").val();
	SkuSeleccionados=[];
    if (cadena != ""){
        $('#Loading').show();
        $.ajax({
			url:"/Grupos/ConsultaBase",
			type: 'POST',
			dataType:'json',
			data:{
				GrupoBase : $('#GruposBase').val(),
				Cadena: $('#inputBuscaBase').val(),
				Filtro: $('input:checked[type=radio][name=filtroBase]').val(),
				Busqueda: $('#AvanzadaBase').val(),
				Tipo: $('#Tipo').val()
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraBase").empty();						
						$("#CabeceraBase").append(data.SCabecera);
						$("#BodyBase").empty();						
						$("#BodyBase").append(data.SBody);
						$("#PaginacionBaseUp").empty();		
						$("#PaginacionBaseUp").append(data.SPaginacion);
						$("#PaginacionBaseDown").empty();		
						$("#PaginacionBaseDown").append(data.SPaginacion);										
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}				
	 				$('#Loading').hide(); 
			},
		  error: function(data){
            alertify.error("Error inesperado, favor de intentar más tarde.");
            				$('#Loading').hide(); 
		  },
		});
    }else{
        alertify.error("Introduce una cadena de texto a consultar.");
        $("#inputBuscaBase").focus();
    }
}

function BuscaPagina(num){
			$.ajax({
			url:"/Grupos/search",
			type: 'POST',
			dataType:'json',
			data:{
				Pag : num,
				Filtro: $('input:checked[type=radio][name=filtroBase]').val()
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraBase").empty();						
						$("#CabeceraBase").append(data.SCabecera);
						$("#BodyBase").empty();						
						$("#BodyBase").append(data.SBody);
						$("#PaginacionBaseUp").empty();		
						$("#PaginacionBaseUp").append(data.SPaginacion);
						$("#PaginacionBaseDown").empty();		
						$("#PaginacionBaseDown").append(data.SPaginacion);						
						Verificarseleccionados();				
					}else{						
						alertify.error(data.SMsj);
					}
				}else{
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}				
			},
		  error: function(data){
            alertify.error("Error inesperado, favor de intentar más tarde.");
		  },
		});
}

 function SubmitGroup(){
	 	$('#Loading').show();
			$.ajax({
			url:"/Grupos/agrupa",
			type: 'POST',
			dataType:'json',
			data:{
				Grupox : $('#Grupos').val(),
				searchbox: $('#searchbox').val()
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraBase").empty();						
						$("#CabeceraBase").append(data.SCabecera);
						$("#BodyBase").empty();						
						$("#BodyBase").append(data.SBody);
						$("#PaginacionBaseUp").empty();		
						$("#PaginacionBaseUp").append(data.SPaginacion);
						$("#PaginacionBaseDown").empty();		
						$("#PaginacionBaseDown").append(data.SPaginacion);							
						Verificarseleccionados();							
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

 function SubmitGroupBase(){
	 $('#Loading').show();
			$.ajax({
			url:"/Grupos/agrupaB",
			type: 'POST',
			dataType:'json',
			data:{
				GrupoBase : $('#GruposBase').val(),
				Cadena: $('#inputBuscaBase').val(),
				Filtro: $('input:checked[type=radio][name=filtroBase]').val(),
				Busqueda: $('#AvanzadaBase').val(),
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){			
						$("#CabeceraBase").empty();						
						$("#CabeceraBase").append(data.SCabecera);
						$("#BodyBase").empty();						
						$("#BodyBase").append(data.SBody);
						$("#PaginacionBaseUp").empty();		
						$("#PaginacionBaseUp").append(data.SPaginacion);						
						$("#PaginacionBaseDown").empty();		
						$("#PaginacionBaseDown").append(data.SPaginacion);					
					}else{
						$("#CabeceraBase").empty();						
						$("#BodyBase").empty();						
						$("#PaginacionBaseUp").empty();								
						$("#PaginacionBaseDown").empty();			
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

