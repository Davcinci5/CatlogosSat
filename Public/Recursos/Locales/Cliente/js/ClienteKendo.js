$(document).ready(function () {

	$('select[name=Almacenes]').kendoMultiSelect().data("kendoMultiSelect");
	$('select[name=Grupos]').kendoMultiSelect().data("kendoMultiSelect");
	//$('select[name=Tipo]').kendoMultiSelect().data("kendoMultiSelect");

	$('#FechaNacimiento').datepicker({
    	isRTL: false,
    	format: 'dd-mm-yyyy',
   	 	autoclose:true,
    	language: 'es'
	});

});