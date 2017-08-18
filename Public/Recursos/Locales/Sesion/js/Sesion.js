

	//##############################< SCRIPTS JS >##########################################
	//################################< Sesion.js >#####################################
	//#########################< VALIDACIONES DE JEQUERY >##################################
	
	$( document ).ready( function () {	
        $(document).on('click', '.deleteButton', function () {		
            DelSession($(this))
		});
	});


  
function EditaSesion(vista){
	if (vista == "Index" || vista ==""){
		if ($('#Sesiones').val() != ""){
			window.location = '/Sesiones/edita/' + $('#Sesiones').val();
		}else{
			alertify.error("Debe Seleccionar un PermisosUri para editar");
		}
	}else if(vista == "Detalle"){
		if ($('#ID').val() != ""){
			window.location = '/Sesiones/edita/' + $('#ID').val();
		}else{
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Sesiones';
		}
	}

}


function DetallePermisosUri(){
	if ($('#ID').val() != ""){
		window.location = '/Sesiones/detalle/' + $('#ID').val();
	}else{
	alertify.error("Debe Seleccionar un PermisosUri para editar");
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


//Funciones Extras
function DelSession(context){
			$('#Loading').show();
            ID = context.attr("id");
			$.ajax({
			url:"/Sesiones/EliminaByID",
			type: 'GET',
			dataType:'json',
			data:{
				ID : ID,
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){	
                        context.parent().parent().remove();		
                        alertify.success(data.SMsj);				
					}else{		
                        eval(data.SFuncion);				
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

function closeAll(){
			$('#Loading').show();
            ID = $("#ID").val()
			$.ajax({
			url:"/Sesiones/EliminaByName",
			type: 'GET',
			dataType:'json',
			data:{
				ID : ID,
			},
			success: function(data){
				if (data != null){
					if (data.SEstado){		
                        alertify.success(data.SMsj)
						setTimeout(function(){eval(data.SFuncion);},2000);
						
					}else{		
                        eval(data.SFuncion);				
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