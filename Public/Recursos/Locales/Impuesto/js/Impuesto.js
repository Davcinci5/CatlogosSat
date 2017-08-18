

//##############################< SCRIPTS JS >##########################################
//################################< Impuesto.js >#####################################
//#########################< VALIDACIONES DE JEQUERY >##################################

$(document).ready(function () {
	var validator = valida();

$('#Abreviatura').change(function(e){
	alertify.message("Cambió el Tipo de Impuesto");
});

$('#Clasificacion').change(function(e){
	alertify.message("Cambió clasificación");
});

$('#SubClasificacion').change(function(e){
	alertify.message("Cambió sub clasificación");
});


});

function valida() {
	var validator = $("#Form_Alta_Impuesto").validate({
		rules: {

			Abreviatura: {

				required: true

			},
			Clasificacion: {

				required: true

			},
			Nombre: {

				required: true,

				rangelength: [3, 150]

			},
			ValorMin: {

				required: true,
				min : 0

			},
			ValorMax: {

				required: true,
				min : 0

			},
			Tipo: {

				required: true

			},
			Unidad: {

				required: true

			},
		},
		messages: {

			Abreviatura: {
				required: "El campo Abreviatura es requerido."
			},
			Clasificacion: {

				required: "El campo Clasificacion es requerido."
			},
			Nombre: {

				required: "El campo Nombre es requerido.",
				rangelength: "La longitud del campo Nombre debe estar entre  [3, 150]"
			},
			ValorMin: {

				required: "El campo Valor Mínimo es requerido."
			},
			ValorMax: {

				required: "El campo Valor Máximo es requerido."
			},
			Tipo: {

				required: "El campo Tipo es requerido."
			},
			Unidad: {

				required: "El campo Unidad es requerido."
			},
		},
		errorElement: "em",
		errorPlacement: function (error, element) {
			error.addClass("help-block");
			element.parents(".col-sm-5").addClass("has-feedback");

			if (element.prop("type") === "checkbox") {
				error.insertAfter(element.parent("label"));
			} else {
				error.insertAfter(element);
			}

			if (!element.next("span")[0]) {
				$("<span class='glyphicon glyphicon-remove form-control-feedback'></span>").insertAfter(element);
			}
		},
		success: function (label, element) {
			if (!$(element).next("span")[0]) {
				$("<span class='glyphicon glyphicon-ok form-control-feedback'></span>").insertAfter($(element));
			}
		},
		highlight: function (element, errorClass, validClass) {
			$(element).parents(".col-sm-5").addClass("has-error").removeClass("has-success");
			$(element).next("span").addClass("glyphicon-remove").removeClass("glyphicon-ok");
		},
		unhighlight: function (element, errorClass, validClass) {
			$(element).parents(".col-sm-5").addClass("has-success").removeClass("has-error");
			$(element).next("span").addClass("glyphicon-ok").removeClass("glyphicon-remove");
		}
	});
	return validator;
}


function EditaImpuesto(vista) {
	if (vista == "Index" || vista == "") {
		if ($('#Impuestos').val() != "") {
			window.location = '/Impuestos/edita/' + $('#Impuestos').val();
		} else {
			alertify.error("Debe Seleccionar un Impuesto para editar");
		}
	} else if (vista == "Detalle") {
		if ($('#ID').val() != "") {
			window.location = '/Impuestos/edita/' + $('#ID').val();
		} else {
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Impuestos';
		}
	}

}


function DetalleImpuesto() {
	if ($('#Impuestos').val() != "") {
		window.location = '/Impuestos/detalle/' + $('#Impuestos').val();
	} else {
		alertify.error("Debe Seleccionar un Impuesto para editar");
	}
}


function BuscaPagina(num){
			$('#Loading').show();
			$.ajax({
			url:"/Impuestos/search",
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
			url:"/Impuestos/agrupa",
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


