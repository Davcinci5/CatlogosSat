package HTTPErrors

import (
	"html/template"

	"../../Modelos/HTTPErrorsModel"
	"gopkg.in/kataras/iris.v6"
)

//Error403 es una función que maneja el error 404
func Error403(ctx *iris.Context) {
	var Error HTTPErrorsModel.HTTPError
	Error.Mensaje = template.HTML(` No tienes acceso a esta pantalla y/o esta acción.
        							Te sugerimos que regreses a <a href="/">Minisuper Ampliado</a> e intentes de nuevo.`)

	Error.Tipo = template.HTML("ERROR 403 ACCESO DENEGADO")
	Error.Soporte = template.HTML(`Si crees que este mensaje no debería ser un error, contactanos <a href="mailto:help@freedcamp.com">ayuda@Minisuper.com</a>`)

	ctx.Render("AAHTTPErrors.html", Error)
}

//Error404 es una función que maneja el error 404
func Error404(ctx *iris.Context) {
	var Error HTTPErrorsModel.HTTPError
	Error.Mensaje = template.HTML(` No se encontró la URL que solicitaste. Por Favor, asegúrate de usar el link correcto.
        							Te sugerimos que regreses a <a href="/">Minisuper Ampliado</a> e intentes de nuevo.`)

	Error.Tipo = template.HTML("ERROR 404 NO ENCONTRADO")
	Error.Soporte = template.HTML(`Si crees que este mensaje no debería ser un error, contactanos <a href="mailto:help@freedcamp.com">ayuda@Minisuper.com</a>`)

	ctx.Render("AAHTTPErrors.html", Error)
}

//Error500 es una función que maneja el error 500
func Error500(ctx *iris.Context) {
	var Error HTTPErrorsModel.HTTPError
	Error.Mensaje = template.HTML(` Esto nos da mucha pena, no debió ocurrir, pero trabajaremos en ello.
        							Te sugerimos que regreses a <a href="/">Minisuper Ampliado</a>.`)
	Error.Tipo = template.HTML("ERROR 500 ERROR INTERNO")
	Error.Soporte = template.HTML(`Por favor comunícate con nosotros para reportar el problema a <a href="mailto:help@freedcamp.com">ayuda@Minisuper.com</a>`)

	ctx.Render("AAHTTPErrors.html", Error)
}
