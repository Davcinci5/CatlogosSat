//##############################< SCRIPTS JS >##########################################
//################################< PuntoVenta.js >#####################################
//#########################< VALIDACIONES DE JEQUERY >##################################

$(document).click(function (e) {
    $('.popover-markup>.trigger').popover({
        container: "body",
        placement: "auto top",
        trigger: "click",
        html: true,
        title: function() {
            return $(this).parent().find('.head').html();
        },
        content: function() {
            return $(this).parent().find('.content').html();
        }
    });

    if (($('.popover').has(e.target).length == 0) || $(e.target).is('.close')) {
        $('.popover-markup>.trigger').popover('hide');
    }

    $('.popover-markup>.trigger').click(function (e) {
        e.stopPropagation();
    });

});



$(document).ready(function() {
    $(window).keydown(function(event) {
        if (event.keyCode == 13) {
            event.preventDefault();
        }
    });

    var validator = valida();

    $('.popover-markup>.trigger').popover({
        container: "body",
        html: true,
        placement: "auto top",
        trigger: "click",
        title: function() {
            return $(this).parent().find('.head').html();
        },
        content: function() {
            return $(this).parent().find('.content').html();
        }

    });


    $('.popover-markup>.trigger').click(function (e) {
        e.stopPropagation();
    });

    $('#Codigo').keydown(function(e) {
        if (e.which == 13 || e.keyCode == 13) {
            e.preventDefault();
            traeProducto();
        }
    });

    $('#identidicadorOp').keydown(function(e) {
        if (e.which == 13 || e.keyCode == 13) {
            e.preventDefault();
            traeOperacion();
        }
    });


});

function DestruyePop() {
    $('.popover-markup>.trigger').popover("destroy");
}


function valida() {

    var validator = $("#Form_Alta_PuntoVenta").validate({
        rules: {

        },
        messages: {

        },
        errorElement: "em",
        errorPlacement: function(error, element) {
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
        success: function(label, element) {
            if (!$(element).next("span")[0]) {
                $("<span class='glyphicon glyphicon-ok form-control-feedback'></span>").insertAfter($(element));
            }
        },
        highlight: function(element, errorClass, validClass) {
            $(element).parents(".col-sm-5").addClass("has-error").removeClass("has-success");
            $(element).next("span").addClass("glyphicon-remove").removeClass("glyphicon-ok");
        },
        unhighlight: function(element, errorClass, validClass) {
            $(element).parents(".col-sm-5").addClass("has-success").removeClass("has-error");
            $(element).next("span").addClass("glyphicon-ok").removeClass("glyphicon-remove");
        }
    });
    return validator;
}


function EditaPuntoVenta(vista) {
    if (vista == "Index" || vista == "") {
        if ($('#PuntoVentas').val() != "") {
            window.location.href = '/PuntoVentas/edita/' + $('#PuntoVentas').val();
        } else {
            alertify.error("Debe Seleccionar un PuntoVenta para editar");
        }
    } else if (vista == "Detalle") {
        if ($('#ID').val() != "") {
            window.location.href = '/PuntoVentas/edita/' + $('#Operacion').val();
        } else {
            alertify.error("No se puede editar debido a un error de referencias, favor de intentar en el index");
            window.location.href = '/PuntoVentas';
        }
    }

}


function DetallePuntoVenta() {
    if ($('#PuntoVentas').val() != "") {
        window.location.href = '/PuntoVentas/detalle/' + $('#PuntoVentas').val();
    } else {
        alertify.error("Debe Seleccionar un PuntoVenta para editar");
    }
}


//Funcionoes
function quitarProducto(e) {
    var opera = $("#Operacion").val();
    var idproducto = e.parentElement.id;
    $('#Loading').show();
    $.ajax({
        url: "/PuntoVentas/quitarProducto",
        type: 'POST',
        dataType: 'json',
        data: {
            Producto: idproducto,
            Operacion: opera
        },
        success: function(data) {
            if (data != null) {
                if (data.SEstado) {
                    $("#idbody").empty();
                    $("#idbody").append(data.SIhtml);
                    $("#calculadora").empty();
                    $("#calculadora").append(data.SCalculadora);
                    // $("[name='Cantidad']").focusout(validaCantidad)
                } else {
                    //Hubo Pedos
                    alertify.error(data.SMsj);
                }
            } else {
                alertify.error("No se recibio información de lado del servidor, intente de nuevo.");
            }
            console.log(data);
            $('#Loading').hide();
        },
        error: function(a, b, c) {
            console.log(a);
            console.log(b);
            console.log(c);
            alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
            $('#Loading').hide();
        },
    });

}

function traeProducto() {
    var codigo = $("#Codigo").val();
    var opera = $("#Operacion").val();
    var almacen = $("#Almacen").val();

    if (codigo != "") {
        $('#Loading').show();
        $.ajax({
            url: "/PuntoVentas/traeProducto",
            type: 'POST',
            dataType: 'json',
            data: {
                Codigo: codigo,
                Operacion: opera,
                Almacen: almacen
            },
            success: function(data) {
                if (data != null) {
                    if (data.SEstado) {
                        if (data.SElastic) {
                            if (data.SBusqueda != "") {
                                $("#BProductos").empty();
                                $("#BProductos").append(data.SBusqueda);
                                $('#BuscaProductos').trigger("click");
                            } else {
                                alertify.error(data.SMsj);
                                alertify.error("Artículo no disponible en ningún almacen asociado al producto.");
                            }
                        } else {
                            $("#idbody").empty();
                            $("#idbody").append(data.SIhtml);
                            $("#calculadora").empty();
                            $("#calculadora").append(data.SCalculadora);
                        }
                    } else {
                        alertify.error(data.SMsj);
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                }
                console.log(data);
                $('#Loading').hide();
            },
            error: function(a, b, c) {
                console.log(a);
                console.log(b);
                console.log(c);
                alertify.error("Hubo un problema al recibir información del servidor, verifique su conexión.");
                $('#Loading').hide();
            },
        });
        $("#Codigo").val("");
        $("#Codigo").focus();
    } else {
        var validator = valida();
        validator.showErrors({
            "Codigo": "¿Es en serio? ¿Quieres buscar con una cadena vacía? =("
        });

        $("#Codigo").focus();
    }
}


function ValidaNumero(numero) {
    var ValorPrevio = Number(numero.getAttribute("previo"));
    var VentaFraccion = numero.getAttribute("fraccion");
    var Valor = numero.value;

    if (VentaFraccion === "true") {

        if (/^([0-9])*$/.test(Valor) || /^\d*\.?\d*$/.test(Valor)) {
            if (Number(Valor) <= 0) {
                alertify.error("No debe introducir valores menores o iguales a cero.");
                numero.value = ValorPrevio.toFixed(3);
            }
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
            numero.value = ValorPrevio.toFixed(3);
        }

    } else {
        if (/^([0-9])*$/.test(Valor)) {
            if (Number(Valor) <= 0) {
                alertify.error("No debe introducir valores menores o iguales a cero.");
                numero.value = ValorPrevio.toFixed(1);
            }
        } else if (/^\d*\.?\d*$/.test(Valor)) {
            alertify.error("No se pueden vender decimales de este producto.");
            numero.value = ValorPrevio.toFixed(1);
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
            numero.value = ValorPrevio.toFixed(1);
        }

    }
}

function ValidaNumeroModal(numero) {
    var ValorPrevio = Number(numero.getAttribute("previo"));
    var VentaFraccion = numero.getAttribute("fraccion");
    var Valor = numero.value;

    if (VentaFraccion === "true") {

        if (/^([0-9])*$/.test(Valor) || /^\d*\.?\d*$/.test(Valor)) {
            if (Number(Valor) < 0) {
                alertify.error("No debe introducir valores menores o iguales a cero.");
                numero.value = Number(0).toFixed(3);
            } else if (ValorPrevio < Number(Valor)) {
                alertify.error("Deberías ir con un loquero, No deberías estar vendiendo.");
                numero.value = Number(0).toFixed(3);
            }
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
            numero.value = Number(0).toFixed(3);
        }

    } else {
        if (/^([0-9])*$/.test(Valor)) {
            if (Number(Valor) < 0) {
                alertify.error("No debe introducir valores menores o iguales a cero.");
                numero.value = Number(0).toFixed(0);
            } else if (ValorPrevio < Number(Valor)) {
                alertify.error("Deberías ir con un loquero, No deberías estar vendiendo.");
                numero.value = Number(0).toFixed(0);
            }
        } else if (/^\d*\.?\d*$/.test(Valor)) {
            alertify.error("No se pueden vender decimales de este producto.");
            numero.value = Number(0).toFixed(0);
        } else {
            alertify.error("Debe introducir un valor adecuado para la venta.");
            numero.value = Number(0).toFixed(0);
        }

    }
}

function AplicaPeticion(peticion) {
    event.preventDefault();

    var Movimiento = peticion.parentNode.parentNode.getAttribute("id");
    var ValorPrevio = Number(peticion.getAttribute("previo"));
    var ValorNuevo = Number(peticion.value);
    var Almacen = peticion.getAttribute("almacen");

    if (ValorPrevio != ValorNuevo) {
        $('#Loading').show();
        var producto = peticion.id;
        var operacion = $("#Operacion").val();

        $.ajax({
            url: "/PuntoVentas/modificaCantidad",
            type: 'POST',
            dataType: 'json',
            data: {
                Producto: producto,
                Movimiento: Movimiento,
                Operacion: operacion,
                Cantidad: ValorNuevo,
                Previo: ValorPrevio,
                Almacen: Almacen
            },
            success: function(data) {
                if (data != null) {
                    if (data.SEstado) {
                        $("#idbody").empty();
                        $("#idbody").append(data.SIhtml);
                        $("#calculadora").empty();
                        $("#calculadora").append(data.SCalculadora);
                    } else {
                        //Hubo Pedos
                        peticion.value = ValorPrevio.toFixed(1);
                        alertify.error(data.SMsj);
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                }
                console.log(data);
                $('#Loading').hide();
            },
            error: function(data) {
                $('#Loading').hide();
            },
        });
    }

}

function AplicaPeticionModal() {
    event.preventDefault();

    $('#Loading').show();

    var productos = [];
    var almacenes = [];
    var nuevos = [];
    $.each($('#BProductos > tr'), function(i, val) {

        if (val.children[7].id == "thisone") {
            var valornuevo = Number(val.children[7].children[0].value);
            var previo = Number(val.children[7].children[0].attributes[4].value);

            if (valornuevo <= previo && valornuevo > 0) {
                almacenes.push(val.children[7].children[0].attributes[2].value);
                productos.push(val.children[7].children[0].attributes[1].value);
                nuevos.push(val.children[7].children[0].value);
            }
        }
    });
    console.log(productos, almacenes, nuevos);

    var operacion = $("#Operacion").val();

    $.ajax({
        url: "/PuntoVentas/modificaCantidadModal",
        type: 'POST',
        dataType: 'json',
        data: {
            Productos: productos,
            Almacenes: almacenes,
            Operacion: operacion,
            Cantidades: nuevos
        },
        success: function(data) {
            $("#BProductos").empty();
            $('#BuscaProductos').trigger("click");

            if (data != null) {
                if (data.SEstado) {
                    $("#idbody").empty();
                    $("#idbody").append(data.SIhtml);
                    $("#calculadora").empty();
                    $("#calculadora").append(data.SCalculadora);

                    alertify.message(data.SMsj);
                } else {
                    //Hubo Pedos
                    alertify.error(data.SMsj);
                }
            } else {
                alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
            }
            console.log(data);
            $('#Loading').hide();
        },
        error: function(data) {
            $("#BProductos").empty();
            $('#BuscaProductos').trigger("click");
            $('#Loading').hide();
            alertify.error("Ocurrió un problema interno, disculpe, por favor, verifique su conexión y vuelva a intntarlo.");
        },
    });

}

function AplicaPeticionModal2(peticion) {
    event.preventDefault();

    var ValorNuevo = Number(peticion.children[7].children[0].value);
    var ValorPrevio = Number(peticion.children[7].children[0].attributes[4].value);
    var Movimiento = peticion.children[7].children[0].attributes[2].value;
    var Producto = peticion.children[7].children[0].attributes[1].value;

    if (ValorNuevo > 0) {
        if (ValorPrevio != ValorNuevo) {
            $('#Loading').show();

            var operacion = $("#Operacion").val();

            $.ajax({
                url: "/PuntoVentas/modificaCantidadModal2",
                type: 'POST',
                dataType: 'json',
                data: {
                    Producto: Producto,
                    Almacen: Movimiento,
                    Operacion: operacion,
                    Cantidad: ValorNuevo
                },
                success: function(data) {
                    if (data != null) {
                        if (data.SEstado) {
                            $("#idbody").empty();
                            $("#idbody").append(data.SIhtml);
                            $("#calculadora").empty();
                            $("#calculadora").append(data.SCalculadora);
                        } else {
                            //Hubo Pedos
                            peticion.value = ValorPrevio.toFixed(1);
                            alertify.error(data.SMsj);
                        }
                    } else {
                        alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                    }
                    console.log(data);
                    $('#Loading').hide();
                },
                error: function(data) {
                    $('#Loading').hide();
                },
            });
        }
    } else {
        alertify.error("En serio? No puedes agregar un producto al carrito si no especificas una cantidad.");
    }
}



function traeOperacion() {
    $('#Loading').show();
    var validator = valida();

    var operacion = $("#identidicadorOp").val();
    if (operacion != "") {
        $.ajax({
            url: "/PuntoVentas/traeOperacion",
            type: 'POST',
            dataType: 'json',
            data: {
                Operacion: operacion,
            },
            success: function(data) {
                if (data != null) {
                    if (data.SEstado) {
                        $("#Operacion").val(data.ID);
                        $("#idbody").empty();
                        $("#idbody").append(data.SIhtml);
                        $("#calculadora").empty();
                        $("#calculadora").append(data.SCalculadora);
                        $("#identidicadorOp").val("");
                        alertify.success("Operacion Encontrada");
                    } else {
                        //Hubo Pedos
                        alertify.error(data.SMsj);
                        $("#identidicadorOp").focus();
                        $("#idbody").empty();
                        $("#calculadora").empty();
                    }
                } else {
                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
                    $("#identidicadorOp").focus();
                    $("#idbody").empty();
                    $("#calculadora").empty();
                }
                $('#Loading').hide();
            },
            error: function(data) {
                $('#Loading').hide();
            },
        });

    } else {
        validator.showErrors({
            "identidicadorOp": "No puede buscar valores vacíos"
        });
        $("#identidicadorOp").focus();
        $("#idbody").empty();
        $("#calculadora").empty();
    }
}

function generaTicket(operacion) {

    $.ajax({
        url: '/PuntoVentas/imprime',
        type: 'POST',
        dataType: 'json',
        data: {
            operacion: operacion
        },
        success: function(data) {
            if (data != null) {
                if (data.SEstado) {
                    var doc = new jsPDF();

                    var specialElementHandlers = {
                        '#editor': function(element, renderer) {
                            return true;
                        }
                    };

                    doc.fromHTML(data.SIhtml, 15, 15, {
                        'width': 170,
                        'elementHandlers': specialElementHandlers
                    });
                    doc.save('Ticket_' + data.ID + '.pdf');
                } else {
                    //Hubo Pedos
                    alertify.error(data.SMsj);
                }
            } else {
                alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");
            }
        },
    });
}


function ValidaVenta() {
    console.log($('#idbody').children().length);
    if ($('#idbody').children().length > 0) {

        $("#Form_Alta_PuntoVenta").submit();
    } else {
        alertify.error("No puede procesar el pago sin artículos.");
    }
}

function BuscaPagina(num) {
    $('#Loading').show();

    $.ajax({
        url: "/PuntoVentas/search",
        type: 'POST',
        dataType: 'json',
        data: {
            Pag: num,
        },
        success: function(data) {
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
        error: function(data) {
            $('#Loading').hide();
        },
    });
}


function SubmitGroup() {
    $('#Loading').show();
    $.ajax({
        url: "/PuntoVentas/agrupa",
        type: 'POST',
        dataType: 'json',
        data: {
            Grupox: $('#Grupos').val(),
            searchbox: $('#searchbox').val()
        },
        success: function(data) {
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
        error: function(data) {
            $('#Loading').hide();
        },
    });
}

function ActualizaImpuestoVenta(este) {
    var $this = $(este);
    var valorprevio = $this.attr("Valor");
    var valornuevo = este.parentNode.children[0].value;

    if (/^([0-9])*$/.test(valornuevo) || /^\d*\.?\d*$/.test(valornuevo)) {
        if (Number(valornuevo) < 0) {
            alertify.error("No debe introducir valores menores o iguales a cero.");
            este.parentNode.children[0].value = Number(valorprevio).toFixed(2);
        }else if (valornuevo == "" || /^-?[0-9]*e\^-?[0-9]*$/.test(valornuevo) || /^-?[0-9]*e\-?[0-9]*$/.test(valornuevo)){
            alertify.error("No debe introducir valores raros o con notación científica.");
            este.parentNode.children[0].value = Number(valorprevio).toFixed(2);
        }else{
            DestruyePop();
            if (Number(valornuevo) != Number(valorprevio)){ 

                var operacion = $("#Operacion").val();
                var producto = $this.attr("Producto");
                var factor = $this.attr("Factor");
                var tipo = $this.attr("Tipo");
                var almacen = $this.attr("Almacen");
                var movimiento = $this.attr("Movimiento");
                var precio = $this.attr("Precio");
            

                    $('#Loading').show();
                    $.ajax({
                        url: '/PuntoVentas/modificaImpuesto',
                        type: 'POST',
                        dataType: 'json',
                        data : {
                            Operacion: operacion,
                            Producto: producto,
                            Factor: factor,
                            Tipo: tipo,
                            Almacen: almacen,
                            ValorPrevio: valorprevio,
                            ValorNuevo: valornuevo,
                            Movimiento: movimiento,
                            Precio: precio
                        },
                        success : function(data) {
                                if (data != null){
                                    if (data.SEstado){	
                                        $("#idbody").empty();
                                        $("#idbody").append(data.SIhtml);
                                        $("#calculadora").empty();
                                        $("#calculadora").append(data.SCalculadora);
                                        alertify.message(data.SMsj);
                                    }else{
                                        alertify.error(data.SMsj);
                                    }
                                }else{
                                    alertify.error("Hubo un problema al recibir información del servidor, favor de volver a intentar.");	
                                }
                            $('#Loading').hide();
                        },
                        error: function(data) { 
                            $('#Loading').hide();
                        },
                    });
                }
        }
    } else {
        alertify.error("No puedes modificar Impuestos con esos números.");
        este.parentNode.children[0].value = Number(valorprevio).toFixed(2);
    }



}