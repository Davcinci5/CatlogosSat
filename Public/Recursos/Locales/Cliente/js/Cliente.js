

//##############################< SCRIPTS JS >##########################################
//################################< Cliente.js >#####################################
//#########################< VALIDACIONES DE JEQUERY >##################################

$(document).ready(function () {

	var validator = valida();

	//Asignacion de evento cuando se seleciona un estado se asignan
	//de manera automatica  los municipios
	$(document).on( "change" , ".SelectEstado", function(){
		var id=$(this).attr('data-estado-select');
		var atrid=$(this).attr('id');
		var Estado = $("#"+atrid+" option:selected").val();
		
		if(Estado==""){
			alertify.error("selecionar el Estado");
			return false;
		}
		GetMunicipiosForEStado(id,Estado);
	} );

	//Asignacion de evento cuando se seleciona un municipio se asignan
	//de manera automatica las Colonias
	$(document).on( "change" , ".SelectMunicipio", function(){
		var id=$(this).attr('data-municipo-select');
		var atrid=$(this).attr('id');
		var Municipio = $("#"+atrid+" option:selected").val();
		if(Municipio==""){
			alertify.error("selecionar el Municipio");
			return false;
		}
		GetColoniasForMunicipio(id,Municipio);
	} );
	
	//Asignacion de evento cuando se seleciona una colonia se asigna
	//de manera automatica el CP
	$(document).on( "change" , ".SelectColonia", function(){
		var id=$(this).attr('data-colonia-select');
		var atrid=$(this).attr('id');
		var Colonia = $("#"+atrid+" option:selected").val();
		if(Colonia==""){
			alertify.error("Selecionar la Colonia");
			return false;
		}
		GetCPForColonia(id,Colonia);
	} );

	//Asignacion de evento para agregar una nueva direccion al Contacto Persona
	$(document).on( "click",".addNewDireccionPersona", function(){
		var id=$(this).attr('data-NewDirect-button');
		AddNewDireccionPersona(id,validator);
	});


	//Asignacion de evento para agregar un nuevo correo  al Contacto Persona
	$(document).on( "click",".addNewMailPer", function(){
		var id=$(this).attr('data-new-mail');
		addNewMailPersona(id);
	});

	//Asiganacion de evento para agregar un nuevo telefono al Contacto Persona 
	$(document).on( "click",".addNewTelefonoPer", function(){
		var id=$(this).attr('data-new-Tel');
		addNewTelefonoPersona(id);
	});

	//Asignacion de evento para agregar otro contacto al Contacto persona
	$(document).on( "click",".addNewOtroContPer", function(){
		var id=$(this).attr('data-new-otro');
		addNewOtroPersona(id);
	});

	//Agrega el Html de Una Persona Contacto
	$("#AddPersonaContacto").click(function(){
		var numPerCont=parseInt($("#ContadorPerCont").val())+1;
		$("#PersonasContactos").append(AgregaTemplatePersonaContacto(numPerCont));
		GetEstadosForSelect(numPerCont);
		getTipoDireciones(numPerCont);
		$("#ContadorPerCont").val(numPerCont)
		alertify.success("Persona Contacto Agregada con Exito.");
	});

	$(document).on( "click",".CabContactoPersona", function(){
		var id=$(this).attr('data-Num-Contacto-Persona');
		$("#BodyContactoPersona"+id).toggle("slow" );
	});

	//Se asigna el boton de eliminar  el contacto Persona
	$(document).on( "click",".botonDeleteCP", function(event){
		//var id=$(this).attr('data-Num-Contacto-Persona');
		var divCP=$(this).parent().parent().parent();
		divCP.remove();
		event.stopPropagation();

		//$("#BodyContactoPersona"+id).toggle("slow" );
	});

    //asignacion de el  evento clic para eliminar  una direccion de persona Contacto
	$(document).on( "click",".deleteDirCP", function(){
		  var tr=$(this).parent().parent();
		  alertify.confirm('Eliminar Cliente', 'Desea eliminar la direccion de el Cliente', function(){ 
			  tr.remove();
			  alertify.success('Direccion Del Cliente Eliminada Con Exito'); 
			 }, function(){});
		
		});

//asignacio b para saber si esuna persona fisica o moral 
	$( "#Tipo" ).change(function() {
		if(this.value=="5936efac8c649f1b8839e48d"){
			$("#contentPersonaFisica").show();
		}else if(this.value=="5936efac8c649f1b8839e48e"){
			$("#contentPersonaFisica").hide();
		}		
	});
	
	//captacion de el boton de eliminar direccion de persona Contacto
	$("#Form_Alta_Cliente").submit(function(e) {	
		if(ValidaInformacionAntesGuardar()){
			return;
		}else{
			e.preventDefault();
		}
	});

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	
	$( "#AddNewDireccion" ).click(function() {
		alert
        var errores=0;
  		if($('#Estado').val() == ""){
			$("#Estado").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Estado": "Seleccionar  el estado"
			});	
			errores++;
		}
		
		if($('#Municipio').val() == ""||$('#Municipio').val() == null){

			$("#Municipio").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Municipio": "Seleccionar  el Municipio"
			});	
			errores++;
		}

		if($('#Colonia').val() == ""||$('#Colonia').val() == null){
			$("#Colonia").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Colonia": "Seleccionar la Colonia"
			});	
			errores++;
		}

		if($('#CP').val() == ""){
			$("#CP").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"CP": "Seleccionar o Capturar el Codigo Postal"
			});	
			errores++;
		}

		if($('#Calle').val() == ""){
			$("#Calle").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Calle": "Capturar la calle"
			});
			errores++;
		}

		if($('#NumExterior').val() == ""){
			$("#NumExterior").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"NumExterior": "Capturar el Numero Exterior."
			});
			errores++;
		}
		
		if($('#TipoDireccion').val() == ""||$('#TipoDireccion').val() == null){
			$("#TipoDireccion").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"TipoDireccion": "Seleccionar el tipo de direccion"
			});
			errores++;
		}

		var numDireccionesCliente=0;

		if(errores==0){
			var Estado=$('#Estado').val();
			var EstadoText=$("#Estado option:selected").text();
			var Municipio= $('#Municipio').val();
			var MunicipioText=$("#Municipio option:selected").text();
			var Colonia=$('#Colonia').val();
			var ColoniaText=$("#Colonia option:selected").text();
			var CP=$('#CP').val();
			var Calle=$('#Calle').val();
			var NumExterior=$('#NumExterior').val();
			var NumInterior=$('#NumInterior').val();
			var TipoDireccion=$('#TipoDireccion').val();
			var TipoDireccionText=$("#TipoDireccion option:selected").text();

			var numDireccionesCliente=parseInt( $("#NumDirClient").val())+1;
			
			var tr="<tr class='direccionCliente' data-num-direccion-cliente="+numDireccionesCliente+">"+
				"<td><input type='hidden' name='Estador"+numDireccionesCliente+"' value='"+Estado+"'>"+EstadoText+"</td>"+
				"<td><input type='hidden' name='Municipior"+numDireccionesCliente+"' value='"+Municipio+"'>"+MunicipioText+"</td>"+
				"<td><input type='hidden' name='Coloniar"+numDireccionesCliente+"' value='"+Colonia+"'>"+ColoniaText+"</td>"+
				"<td><input type='hidden' name='cpr"+numDireccionesCliente+"' value='"+CP+"'>"+CP+"</td>"+
				"<td><input type='hidden' name='Caller"+numDireccionesCliente+"' value='"+Calle+"'>"+Calle+"</td>"+
				"<td><input type='hidden' name='NumExteriorr"+numDireccionesCliente+"' value='"+NumExterior+"'>"+NumExterior+"</td>"+
				"<td><input type='hidden' name='NumInteriorr"+numDireccionesCliente+"' value='"+NumInterior+"'>"+NumInterior+"</td>"+
				"<td><input type='hidden' name='TipoDireccionr"+numDireccionesCliente+"' value='"+TipoDireccion+"'>"+TipoDireccionText+"</td>"+
				"<td><button type='button' class='btn btn-danger deleteDirCP'><span class='glyphicon glyphicon-trash btn-xs'></span></button></td>"+
			"</tr>";

			$( "#cuerpoDirecciones" ).append( tr );
			$("#NumDirClient").val(numDireccionesCliente);

			$('#Estado').val("");
			$('#Municipio').html("");
			$('#Colonia').html("");
			$('#CP').val("");
			$('#Calle').val("");
			$('#NumExterior').val("");
			$('#NumInterior').val("");
			$('#TipoDireccion').val("");


		}
	});

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
						$("#CorreoPrincipal").val($('input[name=CorreosPrincipal]:checked', '#Form_Alta_Cliente').val());
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

function valida() {
	var validator = $("#Form_Alta_Cliente").validate({
		rules: {
			Nombre: {
				required: true,
				rangelength: [6, 150]
			},
			RFC: {
				required: true,
				rangelength: [12, 13]
			}
		},
		messages: {
			Nombre: {
				required: "El campo Nombre es requerido.",
				rangelength: "La longitud del campo Nombre debe estar entre  [6, 150]"
			},
			RFC: {
				rangelength: "La longitud del campo RFC debe estar entre  [12, 13]",
				required: "El campo RFC es requerido."
			}
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

function EditaCliente(vista) {
	if (vista == "Index" || vista == "") {
		if ($('#Clientes').val() != "") {
			window.location = '/Clientes/edita/' + $('#Clientes').val();
		} else {
			alertify.error("Debe Seleccionar un Cliente para editar");
		}
	} else if (vista == "Detalle") {
		if ($('#ID').val() != "") {
			window.location = '/Clientes/edita/' + $('#ID').val();
		} else {
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Clientes';
		}
	}

}


function DetalleCliente() {
	if ($('#Clientes').val() != "") {
		window.location = '/Clientes/detalle/' + $('#Clientes').val();
	} else {
		alertify.error("Debe Seleccionar un Cliente para editar");
	}
}


function BuscaPagina(num) {
	$('#Loading').show();

	$.ajax({
		url: "/Clientes/search",
		type: 'POST',
		dataType: 'json',
		data: {
			Pag: num,
		},
		success: function (data) {
			if (data != null) {
				if (data.SEstado) {
					$("#Cabecera").empty();
					$("#Cabecera").append(data.SCabecera);
					$("#Cuerpo").empty();
					$("#Cuerpo").append(data.SBody);
					$("#Paginacion").empty();
					$("#Paginacion").append(data.SPaginacion);
				} else {
					alertify.error(data.SMsj);
				}
			} else {
				alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
			}
			$('#Loading').hide();
		},
		error: function (data) {
			$('#Loading').hide();
		},
	});
}


function SubmitGroup() {
	$('#Loading').show();
	$.ajax({
		url: "/Clientes/agrupa",
		type: 'POST',
		dataType: 'json',
		data: {
			Grupox: $('#Grupos').val(),
			searchbox: $('#searchbox').val()
		},
		success: function (data) {
			if (data != null) {
				if (data.SEstado) {
					$("#Cabecera").empty();
					$("#Cabecera").append(data.SCabecera);
					$("#Cuerpo").empty();
					$("#Cuerpo").append(data.SBody);
					$("#Paginacion").empty();
					$("#Paginacion").append(data.SPaginacion);
				} else {
					alertify.error(data.SMsj);
				}
			} else {
				alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
			}
			$('#Loading').hide();
		},
		error: function (data) {
			$('#Loading').hide();
		},
	});
}


function AgregaNuevaDireccion() {
	if($('#Estado').val() == ""){
			$("#Estado").parent().addClass("has-feedback has-error");
			validator.showErrors({
				"Estado": "Seleccionar  el estado"
			});	
		}
}