//##############################< SCRIPTS JS >##########################################
//################################< Compra.js >#####################################
//#########################< VALIDACIONES DE JEQUERY >##################################

$(document).ready(function() {
    if (document.getElementById("tbody_impuestos").children.length == 0) {
        $('#div_tabla_impuestos').hide();
    }

    var validator = valida();

    $('#Codigo').keydown(function(e) {
        if (e.which == 13 || e.keyCode == 13) {
            e.preventDefault();
            $('#btnBuscar').trigger("click");
        }
    });

    $('#CantidadComprada,#CostoGeneral,#PrecioGeneral').keydown(function(e) {
        if (e.which == 13 || e.keyCode == 13) {
            e.preventDefault();
            var validator = valida();
            if ($("#CantidadComprada").val() == "") {
                validator.showErrors({
                    CantidadComprada: "El campo cantidad es obligatorio"
                });
                $("#CantidadComprada").focus();
            } else if ($("#CostoGeneral").val() == "") {
                validator.showErrors({
                    CostoGeneral: "El campo costo es obligatorio"
                });
                $("#CostoGeneral").focus();

            } else if ($("#PrecioGeneral").val() == "") {
                validator.showErrors({
                    CostoGeneral: "El campo precio es obligatorio"
                });
                $("#PrecioGeneral").focus();
            } else if ($("#AlmacenesProductos").children().length == 0) {
                alertify.error("Ningun Producto selecionado");
            }


        }
    });
    // $('#CostoGeneral').keydown(function(e) {
    // 	if(e.which == 13 || e.keyCode == 13) {
    // 			e.preventDefault();			
    // 	}
    // });	
    // $('#PrecioGeneral').keydown(function(e) {
    // 	if(e.which == 13 || e.keyCode == 13) {
    // 			e.preventDefault();			
    // 	}
    // });		

    $('#AgregaImpuesto').click(function() {

        if ($('#Impuesto').val() == "") {
            validator.showErrors({
                "Impuesto": "No puede agregar Impuestos sin valor."
            });

            $("#Impuesto").focus();
        } else if (Number($('#Impuesto').val()) < 0) {
            validator.showErrors({
                "Impuesto": "No puede agregar Impuestos con valores negativos."
            });
        } else if (!(/^\d*\.?\d*$/.test($('#Impuesto').val()))) {
            validator.showErrors({
                "Impuesto": "No puede agregar Impuestos con valores no numéricos."
            });
        } else if ($('#TipoImpuesto').val() == "") {
            validator.showErrors({
                "TipoImpuesto": "Especifique el tipo."
            });
        } else {
            $('#div_tabla_impuestos').show();
            $("#tbody_impuestos").append(
                '<tr>\n\
				<td><input type="hidden" class="form-control" name="TipoImpuestoLista" value="' + $('#TipoImpuesto').val() + '"><input type="text" class="form-control" value="' + $("#TipoImpuesto option:selected").text() + '" readonly></td>\n\
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="ValorImpuestoLista" value="' + $("#Impuesto").val() + '" readonly></td>\n\
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="FactorImpuestoLista" value="' + $('input:checked[type=radio][name=Factor]').val() + '" readonly></td>\n\
				<td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="TratamientoImpuestoLista" value="' + $('input:checked[type=radio][name=Retenido]').val() + '" readonly></td>\n\
				<td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
				</tr>');

            $("#Impuesto").val("");
            $("#Impuesto").focus();
        }
    });

    $('#Impuesto').keydown(function(e) {
        if (e.which == 13 || e.keyCode == 13) {
            e.preventDefault();
            $('#AgregaImpuesto').trigger("click");
        }
    });



});

$(document).on('click', '.deleteButton', function() {
    $(this).parent().parent().remove();
    if (document.getElementById("tbody_impuestos").children.length == 0) {
        $('#div_tabla_impuestos').hide();
    }
});
//Valida realiza la validacion de todo el formulario
function valida() {
    var validator = $("#Form_Alta_Compra").validate({
        rules: {
            AlmacenDefecto: {
                required: true
            },
            Codigo: {
                required: true
            },
            CantidadComprada: {
                required: true,
            },
            CostoGeneral: {
                required: true,
            },
            PrecioGeneral: {
                required: true
            },
        },
        messages: {
            AlmacenDefecto: {
                required: "Seleccione un almacen"
            },
            Codigo: {
                required: "El campo no debe estar vacio"
            },
            CantidadComprada: {
                required: "El campo Cantidad es requerido.",
                type: "Debe ser numero"
            },
            CostoGeneral: {
                required: "El campo Costo es requerido.",
            },
            PrecioGeneral: {
                required: "El campo Precio es requerido."
            },
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

//obtenerAlmacenDefecto lee el select que se muestra en la vista y regresa el nombre del almacen seleccionado
function obtenerAlmacenDefecto() {
    var almacenDefecto = $('#AlmacenDefecto').val();
    return almacenDefecto;
}

//obtenerProductoIngresado lee el input de la vista en donde se ingresa el nombre o codigo del producto y regresa el valor ingresado
function obtenerProductoIngresado() {
    return $('#Codigo').val();
}

//buscarProductos realiza una busqueda general en la coleccion de mongo, sin tomar en cuenta el almacen por defecto
//posteriormente se tendrá que buscar en el almacen por defecto
//Obtiene los atributos del producto y lo muestra en una ventana modal
function buscarProductos() {
    $('#Loading').show();
    var almacenDefecto = obtenerAlmacenDefecto();
    var producto = obtenerProductoIngresado();
    var validator = valida();
    $("#Carrito").empty();
    if (almacenDefecto != "") {
        if (producto != "") {
            $.ajax({
                url: "/ConsultarProductos",
                type: 'POST',
                dataType: 'json',
                async: false,
                data: { nombreProducto: producto,
                        Almacen: almacenDefecto
                     },
                success: function(data) {
                    if (data.SEstado == false) {
                        $('#Loading').hide();
                        alertify.error(data.SMsj);
                    } else {
                        $('#Loading').hide();
                        $("#Carrito").empty();
                        $("#Carrito").append(data.SIhtml);
                        $('#btnModal').trigger("click");
                    }
                },
                error: function() {
                    $('#Loading').hide();
                    alertify.error("El servidor de búsqueda no responde...");
                }
            });

        } else {
            $('#Loading').hide();
            validator.showErrors({
                "Codigo": "El campo no debe estar vacio"
            });
        }
    } else {
        $('#Loading').hide();
        validator.showErrors({
            "AlmacenDefecto": "Debe seleccionar un almacen"
        });
    }
}

//leerFilaProducto recorre la ventana modal en donde se muestan las caracteristicas del producto
//Al seleccionar el boton "selecciona" cierra la modal y guarda los atributos en una coleccion
function leerFilaProducto(identificador) {

    var id = $("#filaProducto" + identificador.id);
    var numeros_hijos = id.children().length;
    var datos = id.children();
    var trs = [];
    var filas = [];
    for (var i = 0; i < numeros_hijos; i++) {
        if (i < numeros_hijos - 1) {
            trs.push(datos[i]);
        }
    }
    filas.push(trs);
    $("#CerrarModal").trigger("click");
    valida();
    var elementos = crearEtiquetasProductos();
    MostrarDetallesProducto(identificador.id, filas, elementos);

    var $this = $(identificador);
    MostrarImpuestosDeProducto($this.attr("tipos"),$this.attr("factores"),$this.attr("tratamientos"),$this.attr("valores"));

}

//MostrarImpuestosDeProducto recibe strings de datos de impuestos para mostrarlos en la tabla correspondiente
function MostrarImpuestosDeProducto(tipos, factores, tratamientos, valores) {

    $('#tbody_impuestos').empty();
    var Tipos = tipos.split(",")
    var Factores = factores.split(",")
    var Tratamientos = tratamientos.split(",")
    var Valores = valores.split(",")

    if (Tipos.length >0) {
        $('#div_tabla_impuestos').show();
        
        for (var i = 0; i < Tipos.length; i++) {
                $("#tbody_impuestos").append(
                    '<tr>\n\
                    <td><input type="hidden" class="form-control" name="TipoImpuestoLista" value="' + Tipos[i] + '"><input type="text" class="form-control" value="' + Tipos[i] + '" readonly></td>\n\
                    <td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="ValorImpuestoLista" value="' + Valores[i] + '" readonly></td>\n\
                    <td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="FactorImpuestoLista" value="' + Factores[i] + '" readonly></td>\n\
                    <td><input type="hidden" class="form-control" value=""><input type="text" class="form-control" name="TratamientoImpuestoLista" value="' + Tratamientos[i] + '" readonly></td>\n\
                    <td><button type="button" class="btn btn-danger deleteButton"><span class="glyphicon glyphicon-trash btn-xs"></span></button></td>\n\
                    </tr>');
        }
    }

}

//MostrarDetallesProducto recibe una coleccion de elementos a mostrar (atributos del producto)		
//los formatea (coloca en un input) y los muestra en pantalla
function MostrarDetallesProducto(idProducto, elementos, etiquetas) {
    $("#listaProductos").html("");
    var idOculto = "<input type='hidden' name='idProducto' value='" + idProducto + "'>";

    for (var i = 0; i < elementos.length; i++) {
        for (var j = 0; j < elementos[i].length; j++) {
            var html;
            var numeroElementos = elementos[i][j].childElementCount;
            if (numeroElementos > 1) {
                html = "";
                var subElementos = elementos[i][j].childNodes;
                subElementos.forEach(
                    function(value, key, listObj, argument) {
                        var infor = value.innerText;
                        html += "<input type='text' value='" + infor + "' readonly></input>";
                    },
                    "arg"
                );
            } else {
                html = "";
                var infor = elementos[i][j].innerText;
                if (j != elementos[i].length - 1) {
                    html += "<input type='text' value='" + infor + "' readonly></input>";
                } else {
                    if (elementos[i][j].hasChildNodes()) {
                        html += `<img id="Imagen" alt="Responsive image" name="ImagenProducto" width="80px" height="80px" ;"="" src="` + elementos[i][j].firstElementChild.src + `">`;
                    }
                }
            }
            var etiqueta = "<label>" + etiquetas[j] + "</label>";
            $("#listaProductos").append("<tr><td>" + etiqueta + "</td><td>" + html + "</tr>");
        }
    }
    $("#listaProductos").append(idOculto);
    obtenerIdsAlmacenesMongo(idProducto);
}

//crearEtiquetasProductos crea un arreglo de etiquetas que describirán los atributos de un producto en específico
//Regresa un array con dichas etiquetas
function crearEtiquetasProductos() {
    var etiquetasProductos = new Array("Descripcion:", "Codigos:", "Tipo:", "Unidad:", "Estatus:", "Etiquetas:", "Imagen:");
    return etiquetasProductos;
}

//obtenerIdsAlmacenesMongo consulta los almacenes de tipo "propio" de la coleccion de mongo
//la funcion regresa un arreglo de identificadores (a excepcion del almacen por defecto)
//por cada identificador de almacen encontrado, se dirige a la funcion de consulta de detalles de dicho almacen
function obtenerIdsAlmacenesMongo(idProducto) {
    var almacenDefecto = obtenerAlmacenDefecto();
    $.ajax({
        type: "POST",
        url: '/ConsultarAlmacenesMongo',
        dataType: "json",
        async: true,
        success: function(data) {
            $("#AlmacenesProductos").html("");
            for (x = 0; x < data.length; x++) {
                //Unicamente buscara los detalles de los almacenes que no sean por defecto
                if (almacenDefecto != data[x]) {
                    //var idAlmacen = "Inventario_" + data[x];
                    var idAlmacen = data[x];
                    if (esPar(x)) {
                        var nombreAlmacen = data[x + 1];
                        obtenerInformacionAlmacenesPostgres(idProducto, idAlmacen, nombreAlmacen);
                    }
                }
            }
        },
        error: function() {
            alert('Error occured');
        }
    });
}

//esPar determina si el numero que recibe es par o no (true, false)
function esPar(numero) {
    if ((numero % 2) == 0) {
        return true;
    }
    return false;
}

//checkAlmacen Asigna el valor 0 o 1 al almacen seleccionado para dar de AltasAlmacenes
//Requiere el identificador de un input de tipo hidden para enviar el valor por post
function checkAlmacen(seleccionado, idInput, idProducto) {
    if (seleccionado) {
        //verificara la existencia en el Inventario corrrespondiente
        $.ajax({
            type: "POST",
            url: '/ConsultarExistenciaInventario',
            dataType: "text",
            async: true,
            data: { identificadorAlmacen: idInput, idproducto: idProducto },
            success: function(data) {
                if (data == "EXISTE") {
                    alert("El producto ya esta dado de alta en el almacen\nEste producto ya no se podra dar de alta");
                    $('#Chk' + idInput).prop('checked', false);
                    $('#Chk' + idInput).prop('disabled', 'disabled');
                } else {
                    $('#' + idInput).val(1);
                    //inserta el valor 0 a la cantidad del almacen
                    $('#Alm' + idInput).val(0);
                }
            },
            error: function() {
                alert('Error occured');
            }
        });
    } else {
        $('#' + idInput).val(0);
        //Elimina el valor de la cantidad en el almacen
        $('#Alm' + idInput).innerText = "";
    }
}

//obtenerInformacionAlmacenesPostgres Consulta el inventario de postgres y devuelve la existencia y los precios del almacen requerido
function obtenerInformacionAlmacenesPostgres(idProducto, idAlmacen, nombreAlmacen) {
    $('#Loading').show();
    $.ajax({
        type: "POST",
        url: '/ConsultarAlmacenesPostgres',
        dataType: "json",
        async: true,
        data: { idproducto: idProducto, identificadorAlmacen: idAlmacen },
        success: function(data) {
            if (data.Encontrado) {
                var checkAlta = "<td><input type='checkbox' id='Chk" + idAlmacen + "' onClick='checkAlmacen(this.checked,`" + idAlmacen + "`,`" + idProducto + "`)' class='custom-control-input'></td>";
                //var inputCheckedAlmacenes = "<input type='hidden' id='"+idAlmacen+"' name='AltasAlmacenes' value='0'>";
                //var producto = "<td>"+data.IDProducto+"</td>";
                var existe = "<td>" + data.Existencia + "</td>";
                //var estatus = "<td>"+data.Estatus+"</td>";
                var costo = "<td>" + data.Costo + "</td>";
                var precio = "<td>" + data.Precio + "</td>";
                var inputCantidad = "<td><input type='input' id='Alm" + idAlmacen + "' name='Cantidades' onfocusout='actualizaCantidad(this)' ></td>";
                var inputAlmacenes = "<input type='hidden' name='Almacenes' value='" + idAlmacen + "'>";

                //var html = "<tr>"+checkAlta+inputCheckedAlmacenes+"<td>"+nombreAlmacen+"</td>"+existe+costo+precio+inputCantidad+inputAlmacenes+"</tr>";
                var html = "<tr>" + "<td>" + nombreAlmacen + "</td>" + existe + costo + precio + inputCantidad + inputAlmacenes + "</tr>";
                $("#AlmacenesProductos").append(html);
            }
        },
        error: function() {
            alert('Error occured');
        }
    });
    $('#Loading').hide();
}

function actualizaCantidad(e) {
    var name = e.name;
    var valor = Number(e.value);
    if (!isNaN(valor)) {
        if (valor < $("#CantidadComprada").val() || valor == 0) {
            var reciduo = $("#CantidadComprada").val() - valor;
            $("#CantidadComprada").val(reciduo);
        } else {
            var validator = valida();
            alertify.error("No se puede surtir esa candidad, Debe ser una cantidad menor a la compra");
        }

    }
}

function ValidaNumero(numero) {
    var Valor = numero.value;
    if (/^\d*\.?\d*$/.test(Valor)) {
        if (Number(Valor) < 0) {
            alertify.error("No debe introducir valores menores a cero.");
            numero.value = "";
        }
    } else {
        alertify.error("No es un número adecuado para dar de alta un impuesto.");
        numero.value = "";
    }
}