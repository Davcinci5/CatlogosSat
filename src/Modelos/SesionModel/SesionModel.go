package SesionModel

import "html/template"

//SSesiones estructura para la vista de Sesiones
type SSesiones struct {
	SEstado bool
	SMsj    string
	Sesion
	SIndex
	SSesion
}

//SIndex estructura de variables de index
type SIndex struct {
	SResultados bool
	SRMsj       string
	SCabecera   template.HTML
	SBody       template.HTML
	SPaginacion template.HTML
	SGrupo      template.HTML
}

//SSesion estructura de variables de sesion de Usuarios del sistema
type SSesion struct {
	Name          string
	MenuPrincipal template.HTML
	MenuUsr       template.HTML
}

//Sesion estructura para representar los datos de la sesion de un usuario
type Sesion struct {
	Nombre   string
	IEstatus bool
	IMsj     string
	Ihtml    template.HTML
}

//SAjaxSesiones estructura para la Respuesta AJAX de Sesiones
type SAjaxSesiones struct {
	SEstado  bool
	SMsj     string
	SFuncion string
	Ihtml    template.HTML
}
