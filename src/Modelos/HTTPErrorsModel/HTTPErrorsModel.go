package HTTPErrorsModel

import "html/template"

//HTTPError es una estructura que contiene las partes del error a enviar a la vista.
type HTTPError struct {
	Tipo    template.HTML
	Mensaje template.HTML
	Soporte template.HTML
}
