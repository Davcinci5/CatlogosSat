/**
 * funcion que valida el formulario  antes de ser enviado a guardar 
 */
function ValidaInformacionAntesGuardar(){
	var msg="";
	$("#MensagesClientes").html("");
	if ($('.direccionCliente').length > 0){
		$("#SufijosDirPersona").html("")
		$( ".direccionCliente" ).each(function() {	
			$("#SufijosDirPersona").append("<input type='hidden' name='NumDirCliente' value='"+$( this ).attr( "data-num-direccion-cliente" )+"'>");
		});
	}else{
		msg+="Agrege Como Minimo una Direccion<br/>";
	}

	if (!$('input[name="CorreosPrincipal"]').length > 0){
		msg+="Agrege Como Minimo un Correo Electronico<br/>";
	}
	if (!$('input[name="TelefonoPrincipal"]').length > 0){
		msg+="Agrege Como Minimo un Telefono<br/>";
	}
	$( ".divPerContacto" ).each(function( index ) {
		var numContPer=$( this ).attr("data-persona-contacto")
		if($("#Nombre"+numContPer).val()==""){
			msg+="Capturar Nombre de Personas Contacto<br/>";
			addErrorEchoAmano("Nombre"+numContPer,"Capturar El Nombre De La Persona Contacto");	
		}else{
			removeErrorLlaValidoEchoAmano("Nombre"+numContPer);
		}
	});
	
	if(msg!=""){
		$("#MensagesClientes").append('<div class="alert alert-danger alert-dismissible" role="alert"><button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button><strong><h4 class="text-center">'+msg+'</h4></strong></div>');
		return false;
	}

	$( ".divPerContacto" ).each(function( index ) {
		var numContPer=$( this ).attr("data-persona-contacto")
		$( ".direccionTr" +numContPer).each(function( index ) {
			var numDirPer=$( this ).attr("data-num-dir-contper")
			$("#SufijosDirPersona"+numContPer).append("<input type='hidden' name='NumDirPerCont"+numContPer+"' value='"+numContPer+"-"+numDirPer+"'>");
		});
		$("#identificadoresPersonas").append("<input type='hidden' name='NumPer' value='"+numContPer+"'>");
	});
	return true;
}

/**
 * Funcion que agrega de Manera dinamica la vista de nuevas personas contacto
 * @param {*} Id el numero de la persona creada
 */
function AgregaTemplatePersonaContacto(Id){
	var HmlPersonaContacto = 
				'<div id="divPerCon'+Id+'" class="divPerContacto " data-persona-contacto="'+Id+'">' +
					'<div class="CabContactoPersona" data-Num-Contacto-Persona="'+Id+'">' +
						'<div class="titlePerCon">' +
							'<span class="glyphicon glyphicon-user" ></span> ' +
							'Persona Contacto' +
							'<div class="botonDeleteCP"> <span class="glyphicon glyphicon-trash" ></span> Eliminar</div>' +
						'</div>' +
					'</div>' +
					'<div id="BodyContactoPersona'+Id+'" ><br/>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="Nombre'+Id+'">Nombre:</label> ' +
							'<div class="col-sm-5">' +
								'<input type="text" name="Nombre'+Id+'" id="Nombre'+Id+'" class="form-control" value="" >' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="Estado' + Id + '">Estado:</label> ' +
							'<div class="col-sm-5">' +
								'<select id="Estado' + Id + '" name="Estado' + Id + '" class="form-control selectpicker SelectEstado" data-estado-select="'+Id+'">' +
								'</select>' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="Municipio' + Id + '">Municipio:</label>' +
							'<div class="col-sm-5">' +
								'<select id="Municipio' + Id + '" name="Municipio' + Id + '" class="form-control selectpicker SelectMunicipio" data-municipo-select="'+Id+'">' +
								'</select>' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="Colonia' + Id + '">Colonia:</label>' +
							'<div class="col-sm-5">' +
								'<select id="Colonia' + Id + '" name="Colonia' + Id + '" class="form-control selectpicker SelectColonia" data-colonia-select="'+Id+'">' +
								'</select>' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="CP' + Id + '">CP:</label>' +
							'<div class="col-sm-5">' +
								'<input type="text" name="CP' + Id + '" id="CP' + Id + '" class="form-control" value="">' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="Calle' + Id + '">Calle:</label>' +
							'<div class="col-sm-5">' +
								'<input type="text" name="Calle' + Id + '" id="Calle' + Id + '" class="form-control" value="">' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="NumExterior' + Id + '">NumExterior:</label> ' +
							'<div class="col-sm-5">' +
								'<input type="text" name="NumExterior' + Id + '" id="NumExterior' + Id + '" class="form-control" value="">' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="NumInterior' + Id + '">NumInterior:</label>' +
							'<div class="col-sm-5">' +
								'<input type="text" name="NumInterior' + Id + '" id="NumInterior' + Id + '" class="form-control" value="">' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<label class="col-sm-4 control-label" for="TipoDireccion' + Id + '">Tipo Direccion:</label> ' +
							'<div class="col-sm-5">' +
								'<select id="TipoDireccion' + Id + '" name="TipoDireccion' + Id + '" class="form-control selectpicker">' +
								'</select>' +
							'</div>' +
						'</div>' +
						'<div class="form-group">' +
							'<div class="col-md-9">' +
								'<div id="SufijosDirPersona' + Id + '"></div>'+
								'<button id="AddNewDireccion'+Id+'" type="button" class="btn btn-primary pull-right btn-lg addNewDireccionPersona" data-NewDirect-button="'+Id+'"><span class="glyphicon glyphicon-plus " ></span>Agregar Direccion</button>' +
								'<input type="hidden" name="NumDirContPer'+Id+'"  id="NumDirContPer'+Id+'" value="0">'+
							'</div>' +
						'</div>' +
						'<h3 class="text-center"><strong>Direcciones</strong></h3>'+
						'<table class="table" id="tableDir'+Id+'">' +
							'<thead>' +
								'<tr>' +
									'<th>Cp</th>' +
									'<th>Estado</th>' +
									'<th>Municipio</th>' +
									'<th>Colonia</th>' +
									'<th>Calle</th>' +
									'<th>Num. Ext</th>' +
									'<th>Num. Int</th>' +
									'<th>Tipo</th>' +
								'</tr>' +
							'</thead>' +
							'<tbody id="cuerpoDirecciones'+Id+'">' +
							'</tbody>' +
						'</table>' +
						'<div class="row">' +
							'<div class="col-sm-12">' +
								'<div class="col-sm-4">' +
									'<div class="row text-center" style="padding-bottom:10px;">' +
										'<button id="AgregaEmail'+Id+'" name="AgregaEmail'+Id+'" value="AgregaEmail'+Id+'" data-new-mail="'+Id+'" type="button" class="btn btn-success btn-lg col-md-8 addNewMailPer" style="float:right;margin-right:10%;"><span class="glyphicon glyphicon-plus"></span>Agregar Email</button>' +
									'</div>' +
									'<div class="form-group">' +
										'<label class="col-sm-4 control-label" for="Email'+Id+'">Email:</label>' +
									'<div class="col-sm-8">' +
										'<input type="hidden" name="CorreoPrincipal'+Id+'" id="CorreoPrincipal'+Id+'" value="" >'+
										'<input type="text" name="Email'+Id+'" id="Email'+Id+'" class="form-control" >' +
									'</div>' +
								'</div>' +
								'<h3 class="text-center"><strong>Correos</strong></h3>'+
								'<div class="col-sm-12 table-responsive container" id="div_tabla_correos">' +
									'<table class="table table-condensed table-hover">' +
										'<thead class="thead-inverse">' +
											'<tr>' +
												'<th>Principal</th>' +
												'<th>Correos</th>' +
												'<th>Eliminar</th>' +
											'</tr>' +
										'</thead>' +
										'<tbody id="tbody_etiquetas_correos'+Id+'">' +
										'</tbody>' +
									'</table>' +
								'</div>' +
							'</div>' +
							'<div class="col-sm-4">' +
								'<div class="row text-center" style="padding-bottom:10px;">' +
									'<button id="AgregaTelefono'+Id+'" name="AgregaTelefono'+Id+'" type="button" data-new-Tel="'+Id+'" class="btn btn-success btn-lg col-md-8 addNewTelefonoPer" style="float:right;margin-right:10%;"><span class="glyphicon glyphicon-plus"></span>Agregar Telefono</button>' +
								'</div>' +
								'<div class="form-group">' +
									'<label class="col-sm-4 control-label" for="Telefono'+Id+'">Telefono:</label>' +
									'<div class="col-sm-8">' +
										'<input type="hidden" name="TelefonoPrincipal'+Id+'" id="TelefonoPrincipal'+Id+'" value="" >'+
										'<input type="text" name="Telefono'+Id+'" id="Telefono'+Id+'" class="form-control">' +
									'</div>' +
								'</div>' +
								'<h3 class="text-center"><strong>Telefonos</strong></h3>'+
								'<div class="col-sm-12 table-responsive container" id="div_tabla_telefonos">' +
									'<table class="table table-condensed table-hover">' +
										'<thead class="thead-inverse">' +
											'<tr>' +
												'<th>Principal</th>' +
												'<th>Telefonos</th>' +
												'<th>Eliminar</th>' +
											'</tr>' +
										'</thead>' +
										'<tbody id="tbody_etiquetas_telefonos'+Id+'">' +
										'</tbody>' +
									'</table>' +
								'</div>' +
							'</div>' +
							'<div class="col-sm-4">' +
								'<div class="row text-center" style="padding-bottom:10px;">' +
									'<button id="AgregaOtro'+Id+'" name="AgregaOtro'+Id+'" type="button" data-new-otro="'+Id+'" class="btn btn-success btn-lg col-md-8 addNewOtroContPer" style="float:right;margin-right:10%;"><span class="glyphicon glyphicon-plus"></span>Agregar Otro Medio</button>' +
								'</div>' +
								'<div class="form-group">' +
									'<label class="col-sm-4 control-label" for="Otro'+Id+'">Otros:</label>' +
									'<div class="col-sm-8">' +
										'<input type="text" name="Otro'+Id+'" id="Otro'+Id+'" class="form-control">' +
									'</div>' +
								'</div>' +
								'<h3 class="text-center"><strong>Otros</strong></h3>'+
								'<div class="col-sm-12 table-responsive container" id="div_tabla_otros">' +
									'<table class="table table-condensed table-hover">' +
										'<thead class="thead-inverse">' +
											'<tr>' +
												'<th>Otros</th>' +
												'<th>Eliminar</th>' +
											'</tr>' +
										'</thead>' +
										'<tbody id="tbody_etiquetas_otros'+Id+'">' +
										'</tbody>' +
									'</table>' +
								'</div>' +
							'</div>' +
						'</div>' +
					'</div>' +
				'</div>' ;
				return HmlPersonaContacto;
}

/**
 * Funcion que Agrega los Estados al Select
 * @param {*} Id parametro para saber a que numero de persona  agregar los estados puede ser Blank
 */
function GetEstadosForSelect(Id) {
	$('#Loading').show();
	$.ajax({
		url: "/Clientes/GetEstadosForSelect",
		type: 'POST',
		dataType: 'json',
		success: function (data) {
			if (data != null) {
				if (data.SEstado) {
					$("#Estado"+Id).html(data.SMsj);
				} else {
					alertify.error("Error al obtener los estados");
				}
			} else {
				alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
			}
			$('#Loading').hide();
		},
		error: function (data) {
			$('#Loading').hide();
		}
	});
}

/**
 * Funcion que obtiene los municipios cunado selecciona  un Estado
 * @param {*} id Parametro que para saber que numero de persona  agregar los municipios 
 * @param {*} Estado Paramero Para Saber el ID del estado Selecionado 
 */
function GetMunicipiosForEStado(id,Estado){

	$('#Loading').show();
	$.ajax({
			url: "/Clientes/GetMunicipiosForClaveEstado",
			type: 'POST',
			dataType: 'json',
			data: { "ID": Estado },
			success: function (data) {
				if (data != null) {
					if (data.SEstado = true) {
						$("#Municipio"+id).html(data.SMsj);
					}
				} else {
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}
				alertify.success("Municipios Agregados Con Exito");
				$('#Loading').hide();
			},
			error: function (data) {
				$('#Loading').hide();
			},
		});
}

/**
 * Funcion que obtiene las colonias deacuerdo al municipio Selecionado
 * @param {*} id Parametro que para saber que numero de persona  agregar las Colonias
 * @param {*} Municipio Paramero Para Saber el ID del Municipio Selecionado
 */
function GetColoniasForMunicipio(id,Municipio){
	$('#Loading').show();
	$.ajax({
			url: "/Clientes/GetColoniasForClaveMunicipio",
			type: 'POST',
			dataType: 'json',
			data: { "ID": Municipio },
			success: function (data) {
				if (data != null) {
					if (data.SEstado = true) {
						$("#Colonia"+id).html(data.SMsj);
					}
				} else {
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}
				alertify.success("Colonias Agregadas con exito");
				$('#Loading').hide();
			},
			error: function (data) {
				$('#Loading').hide();
			},
		});
}

/**
 * Funcion que obtiene  el codigo postal deacuerdo a una colonia selecionada
 * @param {*} id Parametro que para saber que numero de persona  agregar el cp
 * @param {*} Colonia Paramero Para Saber el ID del la colonia Selecionada
 */
function GetCPForColonia(id,Colonia){
		$('#Loading').show();
		$.ajax({
			url: "/Clientes/GetCPForColonia",
			type: 'POST',
			dataType: 'json',
			data: { "ID": Colonia },
			success: function (data) {
				if (data != null) {
					if (data.SEstado = true) {
						$("#CP"+id).val(data.SMsj);
					}
				} else {
					alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
				}
				alertify.success("CP Agregado Con Exito");
				$('#Loading').hide();
			},
			error: function (data) {
				$('#Loading').hide();
			},
		});
}

/**
 * Funcion que obtiene los tipos de direciones
 * @param {*} id Parametro que para saber que numero de persona  agregar los Tipos de direcionnes
 */
function getTipoDireciones(id){	
		$('#Loading').show();
		$.ajax({
			url: "/Clientes/GetTipoDireciones",
			type: 'POST',
			dataType: 'json',
			success: function (data) {
				if (data != null) {
					if (data.SEstado = true) {
						$("#TipoDireccion"+id).html(data.SMsj);
					}
				} else {
					alertify.error("Hubo un problema al obtener tipos de direcciones");
				}
				$('#Loading').hide();
			},
			error: function (data) {
				$('#Loading').hide();
			},
		});
}

/**
 * Funcion que agrega una nueva direccion ala persona
 * @param {*} Id Parametro que para saber que numero de persona  agregar la direccion
 */
function AddNewDireccionPersona(id){

 		var errores=0;
  		if($('#Estado'+id).val() == ""){
			addErrorEchoAmano("Estado"+id,"Seleccionar El estado");	
			errores++;
		}else{
			removeErrorLlaValidoEchoAmano("Estado"+id);
		}
		if($('#Municipio'+id).val() == ""||$('#Municipio'+id).val() == null){
			addErrorEchoAmano("Municipio"+id,"Seleccionar El estado");	
			errores++;
		}else{
			removeErrorLlaValidoEchoAmano("Municipio"+id);
		}

		if($('#Colonia'+id).val() == ""||$('#Colonia'+id).val() == null){
			addErrorEchoAmano("Colonia"+id,"Seleccionar la Colonia");	
			errores++;
		}else{
			removeErrorLlaValidoEchoAmano("Colonia"+id);
		}

		 if($('#CP'+id).val() == ""){
			addErrorEchoAmano("CP"+id,"Agregar El Codigo Postal");
		 	errores++;
		 }else{
			 removeErrorLlaValidoEchoAmano("CP"+id);
		 }

		 if($('#Calle'+id).val() == ""){
			 addErrorEchoAmano("Calle"+id,"Capturar la Calle.");
		 	errores++;
		 }else{
			 removeErrorLlaValidoEchoAmano("Calle"+id);
		 }
		 if($('#NumExterior'+id).val() == ""){
			addErrorEchoAmano("NumExterior"+id,"Capturar el Numero Exterior");
		 	errores++;
		 }else{
			 removeErrorLlaValidoEchoAmano("NumExterior"+id);
		 }	
		 if($('#TipoDireccion'+id).val() == ""||$('#TipoDireccion'+id).val() == null){
			addErrorEchoAmano("TipoDireccion"+id,"selecionar el tipo de direccion");
		 	errores++;
		 }else{
			 removeErrorLlaValidoEchoAmano("TipoDireccion"+id);
		 }

		 if(errores==0){
			var Estado=$('#Estado'+id).val();
			var EstadoText=$("#Estado"+id+" option:selected").text();
			var Municipio= $('#Municipio'+id).val();
			var MunicipioText=$("#Municipio"+id+" option:selected").text();
			var Colonia=$('#Colonia'+id).val();
			var ColoniaText=$("#Colonia"+id+" option:selected").text();
			var CP=$('#CP'+id).val();
			var Calle=$('#Calle'+id).val();
			var NumExterior=$('#NumExterior'+id).val();
			var NumInterior=$('#NumInterior'+id).val();
			var TipoDireccion=$('#TipoDireccion'+id).val();
			var TipoDireccionText=$("#TipoDireccion"+id+" option:selected").text();
			var NumDirPer=parseInt( $("#NumDirContPer"+id).val())+1;

			var tr="<tr class='direccionTr"+id+"' data-Num-Dir-ContPer='"+NumDirPer+"' >"+
				"<td><input type='hidden' name='EstadoPC"+id+"-"+NumDirPer+"' value='"+Estado+"'>"+EstadoText+"</td>"+
				"<td><input type='hidden' name='MunicipioPC"+id+"-"+NumDirPer+"' value='"+Municipio+"'>"+MunicipioText+"</td>"+
				"<td><input type='hidden' name='ColoniaPC"+id+"-"+NumDirPer+"' value='"+Colonia+"'>"+ColoniaText+"</td>"+
				"<td><input type='hidden' name='cpPC"+id+"-"+NumDirPer+"' value='"+CP+"'>"+CP+"</td>"+
				"<td><input type='hidden' name='CallePC"+id+"-"+NumDirPer+"' value='"+Calle+"'>"+Calle+"</td>"+
				"<td><input type='hidden' name='NumExteriorPC"+id+"-"+NumDirPer+"' value='"+NumExterior+"'>"+NumExterior+"</td>"+
				"<td><input type='hidden' name='NumInteriorPC"+id+"-"+NumDirPer+"' value='"+NumInterior+"'>"+NumInterior+"</td>"+
				"<td><input type='hidden' name='TipoDireccionPC"+id+"-"+NumDirPer+"' value='"+TipoDireccion+"'>"+TipoDireccionText+"</td>"+
				"<td><button type='button' class='btn btn-danger deleteDirCP'><span class='glyphicon glyphicon-trash btn-xs'></span></button></td>"+
				"</tr>";

			$( "#cuerpoDirecciones"+id ).append( tr );
            $("#NumDirContPer"+id).val(NumDirPer)
			$('#Estado'+id).val("");
			$('#Municipio'+id).html("");
			$('#Colonia'+id).html("");
			$('#CP'+id).val("");
			$('#Calle'+id).val("");
			$('#NumExterior'+id).val("");
			$('#NumInterior'+id).val("");
			$('#TipoDireccion'+id).val("");

		 }
}

/**
 * Funcion que agrega una nuevo mail ala persona
 * @param {*} id Parametro que para saber que numero de persona  agregar el mail
 */
function addNewMailPersona(id){
	var errores=0;
	var mail=$('#Email'+id).val();
  		if(mail == ""){
			addErrorEchoAmano("Email"+id,"Capturar El Correo Electronico." );	
			errores++;
		}else{
			removeErrorLlaValidoEchoAmano("Email"+id);
			var correoOno =  /^\w+([\.\+\-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,4})+$/;
			if(!correoOno.test(mail)){
				addErrorEchoAmano("Email"+id,"Formato de Correo Incorrecto.");
				errores++;
			}else{
				removeErrorLlaValidoEchoAmano("Email"+id);
			}
		}

		if(errores==0){
			$("#tbody_etiquetas_correos"+id).prepend(
				'<tr>\n\
				<td><input type="radio" name="CorreosPrincipal'+id+'" value="' + $("#Email"+id).val() + '" checked></td>\n\
				<td><input type="text" class="form-control" name="Correos'+id+'" value="' + $("#Email"+id).val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');
			
				$('input[type=radio][name=CorreosPrincipal'+id+']').change(function() {
						$("#CorreoPrincipal"+id).val($('input[name=CorreosPrincipal'+id+']:checked', '#Form_Alta_Cliente').val());
					});
				$("#CorreoPrincipal"+id).val($("#Email"+id).val());	
				$("#Email"+id).val("");
				$("#Email"+id).focus();
		}
	
}

/**
 * funcion que agrega un nuevo telefono ala persona contacto
 * @param {*} id Parametro que recive para saber que numero de persona  agregar el mail
 */
function addNewTelefonoPersona(id){
	var errores=0;
	 var telefono=$('#Telefono'+id).val();
  		if(telefono == ""){
			addErrorEchoAmano("Telefono"+id,"Capturar El Telefono");	
			errores++;
		}else{
			removeErrorLlaValidoEchoAmano("Telefono"+id);
		}

		if(errores==0){
				$("#tbody_etiquetas_telefonos"+id).prepend(
				'<tr>\n\
				<td><input type="radio" name="TelefonosPrincipal'+id+'" value="' + $("#Telefono"+id).val() + '" checked></td>\n\
				<td><input type="text" class="form-control" name="Telefonos'+id+'" value="' + $("#Telefono"+id).val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');
				$('input[type=radio][name=TelefonosPrincipal'+id+']').change(function() {
					$("#TelefonoPrincipal"+id).val($('input[name=TelefonosPrincipal'+id+']:checked', '#Form_Alta_Cliente').val());
				});
				$("#TelefonoPrincipal"+id).val($("#Telefono"+id).val());
				$("#Telefono"+id).val("");
				$("#Telefono"+id).focus();	

		}
}

/**
 * Funcion que agrega otro contacto ala Persona Contato
 * @param {*} id Parametro que recive un id para saber la persona contacto
 */
function addNewOtroPersona(id){
	var errores=0;
	 var otro=$('#Otro'+id).val();
  		if(otro == ""){
			addErrorEchoAmano("Otro"+id,"Capturar El Contacto");	
			errores++;
		}else{
			removeErrorLlaValidoEchoAmano("Otro"+id);
		}

		if(errores==0){
				$("#tbody_etiquetas_otros"+id).append(
				'<tr>\n\
				<td><input type="text" class="form-control" name="Otros'+id+'" value="' + $("#Otro"+id).val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');
				$("#Otro"+id).val("");
				$("#Otro"+id).focus();	

		}

}
///////////////////< Parte de la Validacion echo a mano >//////////////////////

function addErrorEchoAmano(id,menssage){
	var div=$("#"+id).parent();
	if($("#"+id).attr("aria-describedby") == id+"-error"){
		 removeErrorEchoAmano(div,id);
	}
	div.addClass("has-feedback has-error");
	$("#"+id).attr("aria-describedby",id+"-error");
	var span='<span class="glyphicon glyphicon-remove form-control-feedback"></span>';
	var em='<em id="Municipio-error" class="error help-block">'+menssage+'</em>';
	div.append(span);
	div.append(em);
}

function removeErrorEchoAmano(div,id){
	$("#"+id).removeAttr("aria-describedby");
	div.removeClass("has-feedback has-error");
	div.children( "span" ).remove();
	div.children( "em" ).remove();

}

function removeErrorLlaValidoEchoAmano(id){
	var div=$("#"+id).parent();
	$("#"+id).removeAttr("aria-describedby");
	div.removeClass("has-feedback has-error");
	div.children( "em" ).remove();
}