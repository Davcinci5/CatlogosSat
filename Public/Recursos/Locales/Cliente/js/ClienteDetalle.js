	$(document).ready(function () {
		
			$(".divPerContacto").each(function(){
				var numPer=$(this).attr("data-persona-contacto")
				$("#Nombre"+numPer).attr("disabled", true);
				$("#CP"+numPer).attr("disabled", true);
				$("#Estado"+numPer).attr("disabled", true);
				$("#Municipio"+numPer).attr("disabled", true);
				$("#Colonia"+numPer).attr("disabled", true);
				$("#Calle"+numPer).attr("disabled", true);
				$("#NumExterior"+numPer).attr("disabled", true);
				$("#NumInterior"+numPer).attr("disabled", true);
				$("#TipoDireccion"+numPer).attr("disabled", true);
				$("#Email"+numPer).attr("disabled", true);
				$("#Telefono"+numPer).attr("disabled", true);
				$("#Otro"+numPer).attr("disabled", true);
				
			});

			
			$(".CollapseDetalle").hide();			

			$(".deleteDirCP").each(function(){
				$(this).addClass("disabled");
				$( this ).removeClass("deleteDirCP");
			});
			
			$(".deleteButton").each(function(){
				$(this).addClass("disabled");
				$( this ).removeClass("deleteButton");
			});

			$(".addNewDireccionPersona").each(function(){
				$(this).addClass("disabled");
				$( this ).removeClass("addNewDireccionPersona");
			});
			
			$(".addNewMailPer").each(function(){
				$(this).addClass("disabled");
				
				$( this ).removeClass("addNewMailPer");
			});
			
			$(".addNewTelefonoPer").each(function(){
				$(this).addClass("disabled");
				$( this ).removeClass("addNewTelefonoPer");
			});

			$(".addNewOtroContPer").each(function(){
				$(this).addClass("disabled");
				$( this ).removeClass("addNewOtroContPer");
			});

			$(".botonDeleteCP").each(function(){
				$( this ).remove();
			});
		}); 