package ConfiguracionControler

import (
	"fmt"
	"html/template"

	"../../Modulos/Session"
	iris "gopkg.in/kataras/iris.v6"
)

type SSesion struct {
	Name          string
	IDS           string
	MenuPrincipal template.HTML
	MenuUsr       template.HTML
}

//Configuracion renderea al GET de la Configuracion del sistema
func Configuracion(ctx *iris.Context) {

	fmt.Println("Entrando al GET de Configuracion")
	Send := SSesion{}

	NameUsrLoged, MenuPrincipal, MenuUsr, _ := Session.GetDataSession(ctx) //Retorna los datos de la session
	Send.Name = NameUsrLoged
	Send.MenuPrincipal = template.HTML(MenuPrincipal)
	Send.MenuUsr = template.HTML(MenuUsr)

	ctx.Render("Configuracion.html", Send)

}
