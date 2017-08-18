

function ReportarProblema(){
    var titulo = $("#TituloProblema").val().trim();
    var contenido = $("#ContenidoProblema").val().trim();
    var ruta = window.location.pathname;
    
    if ( titulo != "" && contenido != "" ){
        $('#Loading').show();
        $.ajax({
            url: "/Bugs/alta",
            type: 'POST',
            dataType: 'json',
            data: {
				Titulo: titulo,
                Contenido: contenido,
                Ruta:ruta
            },
            success: function(data) {
                if (data != null) {
                    if ( data.SEstado ) {
                        //Cerrar Modal
                        $("#TituloProblema").val("");
                        $("#ContenidoProblema").val("");
                        $("#CerrarModalProblemas").click();
                        alertify.success(data.SMsj);                        
					} else {
                        alertify.error(data.SMsj);
                        eval(data.SFuncion);
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                }
                $('#Loading').hide();
            },
            error: function(a, b, c) {
                alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
                $('#Loading').hide();
            },
        });
    }else{
        alertify.error("Asigna un Titulo y Contenido a tu problema porfavor");
    }
        
}

