window.addEventListener("beforeunload", function (event) {
    $('#Loading').show();
});

$( document ).ready(function() {

     $('#Loading').hide();
    //asignacion de click al boton de ocultar
    $( "#iconMenuLateral" ).on("click", function() {
        ocultarMenuLateral();
    });
    console.log("hola");
     reajustarZisemenu();
    $( window ).resize(function() {
       reajustarZisemenu();
    });

        ocultarMenuLateral();

    /*reajustar el menu*/  
    // var resizeElement = document.getElementById('SeccionCuerpo');
    // addResizeListener(resizeElement, reajustarMenuSizeSecccionCuerpo);
});

function reajustarMenuSizeSecccionCuerpo(){
    var anchoWindow=$(window).width();
    if(anchoWindow > 992){
        var altoContMenu = $("#menuPanel").height();
        var altoContenido=$("#SeccionCuerpo").height();
        if(altoContenido>altoContMenu){
             $("#menuLateral").css('min-height', altoContenido+20+'px'); 
        }
    }
}

//funcion que reajustar  el alto del menu lateral
function reajustarZisemenu(){

    var alto=$(window).height();
    var altoTotal=alto-100;
    var altoContenido=$("#SeccionCuerpo").height();
    var anchoWindow=$(window).width();
    var altoWindow=($(window).height())-50;
    var altoContMenu = $("#menuPanel").height();
    if(anchoWindow <= 992){
            $("#menuLateral").css('min-height',altoContMenu+20+'px'); 
            hideMenu();                
        }else{
            if(altoWindow>altoContenido && altoWindow>altoContMenu  ){
                $("#menuLateral").css('min-height', altoWindow+'px'); 
            }else{
                if(altoContenido>altoTotal){
                    $("#menuLateral").css('min-height', altoContenido+20+'px'); 
                }else{
                    $("#menuLateral").css('min-height', altoTotal+'px');
                }
            }
            showMenu();
        }
}

//funcion para ocultar o mostrar el menu lateral
function ocultarMenuLateral(){
        if( $('#menuLateral').is(":visible") ){
            hideMenu2();
        }else{
            showMenu();
        }
}

function showMenu(){
     $('#SeccionCuerpo').removeClass("col-md-12");
     $('#SeccionCuerpo').addClass("col-md-10");
     $('#menuLateral').show("fast");
}

function hideMenu(){
     $('#menuLateral').hide("fast",function() {
        $('#SeccionCuerpo').removeClass("col-md-10");
        $('#SeccionCuerpo').addClass("col-md-12");
	});
}

function hideMenu2(){
    $('#menuLateral').hide("fast");
    $('#SeccionCuerpo').removeClass("col-md-10");
    $('#SeccionCuerpo').addClass("col-md-12");
}
