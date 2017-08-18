$(document).ready(function () {
	var validator = valida();

	$("#Nombre").focus();

	$.validator.addMethod("regx", function (value, element, regexpr) {
		return regexpr.test(value);
	}, "Cadena no v&aacute;lida.");

	$("#copiarComercialEnFiscal").click(function () {
		$("#CalleFiscal").val($("#CalleComercial").val());
		$("#NumExtFiscal").val($("#NumExtComercial").val());
		$("#NumIntFiscal").val($("#NumIntComercial").val());
		$("#PaisFiscal").val($("#PaisComercial").val());
		$("#EstadoFiscal").val($("#EstadoComercial").val())
		$("#MunicipioFiscal").val($("#MunicipioComercial").val())
		$("#ColoniaLocalidadFiscal").val($("#ColoniaLocalidadComercial").val())
		$("#CPFiscal").val($("#CPComercial").val())

	});
});

function valida() {
	var validator = $("#Form_Alta_Empresa").validate({
		rules: {
			Nombre: {
				required: true,
				rangelength: [4, 1000]
			},
			RazonSocial: {
				required: true,
				rangelength: [4, 1000],
			},
			RFC: {
				required: true,
				rangelength: [12, 13],
				regx: /^(([a-zA-Z]{3})|(([a-zA-Z]){4}))\d{6}([a-zA-Z0-9]){3}$/
			},
			Key: {
				required: true,
				rangelength: [1, 1500]
			},
			Cer: {
				required: true,
				rangelength: [1, 1500]
			},
			Pem: {
				required: true,
				rangelength: [1, 1500]
			},
			Correo: {
				required: true,
				rangelength: [5, 100]
			},
			Pass: {
				required: true,

			},
			Tipo: {
				required: false,
			},
			Puerto: {
				range: [1, 9999],
				required: true
			},
			Cifrado: {
				required: true,
				rangelength: [1, 150]
			},
			CalleComercial: {
				required: true,
				rangelength: [2, 500]
			},
			NumIntComercial: {
				required: false,
				rangelength: [2, 300]
			},
			NumExtComercial: {
				required: true,
				rangelength: [2, 300]
			},
			ColoniaLocalidadComercial: {
				required: true,
				rangelength: [1, 500]
			},
			MunicipioComercial: {
				required: true,
				rangelength: [4, 500]
			},
			EstadoComercial: {
				required: true,
			},
			PaisComercial: {
				required: true,
				rangelength: [4, 200]
			},
			CPComercial: {
				required: true,
				regx: /^\d{5}$/
			},
			CalleFiscal: {
				required: true,
				rangelength: [2, 500]
			},
			NumIntFiscal: {
				required: false,
				rangelength: [2, 300]
			},
			NumExtFiscal: {
				required: true,
				rangelength: [2, 300]
			},
			ColoniaLocalidadFiscal: {
				required: true,
				rangelength: [1, 500]
			},
			MunicipioFiscal: {
				required: true,
				rangelength: [4, 500]
			},
			EstadoFiscal: {
				required: true,
			},
			PaisFiscal: {
				required: true,
				rangelength: [4, 200]
			},
			CPFiscal: {
				required: true,
				regx: /^\d{5}$/
			},
			Alias: {
				required: true,
				rangelength: [5, 50]
			},
			Email: {
				required: true,
				email: true,
				rangelength: [5, 100]
			},
			Telefono: {
				required: true,
				number: true,
				rangelength: [10, 13]
			},
			Movil: {
				required: true,
				number: true,
				rangelength: [10, 13]
			},

		},
		messages: {
			Nombre: {
				rangelength: "La longitud del campo Nombre debe estar entre  [4, 1000]",
				required: "El campo Nombre es requerido."
			},
			RazonSocial: {
				required: "El campo RazonSocial es requerido.",
				rangelength: "La longitud del campo RazonSocial debe estar entre  [4, 1000]",
			},
			RFC: {
				required: "El campo RFC es requerido.",
				rangelength: "La longitud del campo RFC debe estar entre  [12, 13]",
				regex: "La cadena no coincide con el patrón"
			},
			Key: {
				required: "El campo Key es requerido.",
				rangelength: "La longitud del campo Key debe estar entre  [[1, 1500]"
			},
			Cer: {
				required: "El campo Cer es requerido.",
				rangelength: "La longitud del campo Cer debe estar entre  [1, 1500]"
			},
			Pem: {
				required: "El campo Pem es requerido.",
				rangelength: "La longitud del campo Pem debe estar entre  [1, 1500]"
			},
			Correo: {
				required: "El campo Correo es requerido.",
				rangelength: "La longitud del campo Correo debe estar entre  [5, 100]"
			},
			Pass: {
				required: "El campo Pass es requerido."
			},
			Tipo: {
				required: "El campo Tipo es requerido.",
			},
			Puerto: {
				required: "El campo Puerto es requerido.",
				range: "el valor del campo Puerto debe estar entre  [1, 9999]"
			},
			Cifrado: {
				required: "El campo Cifrado es requerido.",
				rangelength: "La longitud del campo Cifrado debe estar entre  [1, 150]"
			},
			CalleComercial: {
				required: "El campo Calle es requerido.",
				rangelength: "La longitud del campo Calle debe estar entre  [4, 500]"
			},
			NumIntComercial: {
				required: "El campo NumInterior es requerido.",
				rangelength: "La longitud del campo NumInterior debe estar entre  [1, 300]"
			},
			NumExtComercial: {
				required: "El campo NumExterior es requerido.",
				rangelength: "La longitud del campo NumExterior debe estar entre  [1,300]"
			},
			ColoniaLocalidadComercial: {
				required: "El campo Colonia es requerido.",
				rangelength: "La longitud del campo Colonia debe estar entre  [1, 500]"
			},
			MunicipioComercial: {
				required: "El campo Municipio es requerido.",
				rangelength: "La longitud del campo Municipio debe estar entre  [4, 500]"
			},
			EstadoComercial: {
				required: "El campo Estado es requerido.",
			},
			PaisComercial: {
				rangelength: "La longitud del campo Pais debe estar entre  [4, 200]",
				required: "El campo Pais es requerido."
			},
			CPComercial: {
				required: "El campo CP es requerido.",
				regx: "El valor del campo CP debe contener cinco d&iacute;gitos"
			},
			CalleFiscal: {
				required: "El campo Calle es requerido.",
				rangelength: "La longitud del campo Calle debe estar entre  [4, 500]"
			},
			NumIntFiscal: {
				required: "El campo NumInterior es requerido.",
				rangelength: "La longitud del campo NumInterior debe estar entre  [1, 300]"
			},
			NumExtFiscal: {
				required: "El campo NumExterior es requerido.",
				rangelength: "La longitud del campo NumExterior debe estar entre  [1,300]"
			},
			ColoniaLocalidadFiscal: {
				required: "El campo Colonia es requerido.",
				rangelength: "La longitud del campo Colonia debe estar entre  [1, 500]"
			},
			MunicipioFiscal: {
				required: "El campo Municipio es requerido.",
				rangelength: "La longitud del campo Municipio debe estar entre  [4, 500]"
			},
			EstadoFiscal: {
				required: "El campo Estado es requerido.",
			},
			PaisFiscal: {
				rangelength: "La longitud del campo Pais debe estar entre  [4, 200]",
				required: "El campo Pais es requerido."
			},
			CPFiscal: {
				required: "El campo CP es requerido.",
				regx: "El valor del campo CP debe contener cinco d&iacute;gitos"
			},
			Alias: {

				required: "El campo Alias es requerido.",
				rangelength: "La longitud del campo Alias debe estar entre  [5, 50]"
			},
			Email: {
				email: "Debe ser un email a@b.c",
				rangelength: "La longitud del campo Email debe estar entre  [5, 100]",
				required: "El campo Email es requerido."
			},
			Telefono: {
				rangelength: "La longitud del campo Telefono debe estar entre  [10, 13]",
				required: "El campo Telefono es requerido."
			},
			Movil: {
				rangelength: "La longitud del campo Movil debe estar entre  [10, 13]",
				required: "El campo Movil es requerido."
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


function EditaEmpresa(vista) {
	if (vista == "Index" || vista == "") {
		if ($('#Empresas').val() != "") {
			window.location = '/Empresas/edita/' + $('#Empresas').val();
		} else {
			alertify.error("Debe Seleccionar un Empresa para editar");
		}
	} else if (vista == "Detalle") {
		if ($('#ID').val() != "") {
			window.location = '/Empresas/edita/' + $('#ID').val();
		} else {
			alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
			window.location = '/Empresas';
		}
	}

}


function DetalleEmpresa() {
	if ($('#Empresas').val() != "") {
		window.location = '/Empresas/detalle/' + $('#Empresas').val();
	} else {
		alertify.error("Debe Seleccionar un Empresa para editar");
	}
}
//LimpiaCadena Funcion que elimina espacios en blanco multiples entre palabras sustituyendo por " "
//y elimina espacios al inicio y al final de una cadena
// entrada: input 
//salida: input.value despues de  equivalente a trim
function LimpiaCadena(entrada) {
	text = entrada.value
	text = text.replace(/([\s]+)/g, ' ');
	text = text.replace(/^\s+|\s+$/, '');
	text = text.toUpperCase();
	entrada.value = text;
}