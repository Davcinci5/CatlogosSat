package LoginModel

import (
	"fmt"
	"html/template"
)

func ConstruirTemplateUsuario(userOn, nivel, id string) string {
	fmt.Println("Hola Mundo desde el Template USUARIO")

	var template string

	return template
}

//SLogin estructura para la Vista de login
type SLogin struct {
	SEstado bool
	SMsj    string
	Login
	SSesion
}

//Login esstructura para variables del login
type Login struct {
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//SSesion estructura de variables de sesion de Usuarios del sistema
type SSesion struct {
	Name          string
	MenuPrincipal template.HTML
	MenuUsr       template.HTML
}
