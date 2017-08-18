package LoginControler

import (
	"fmt"
	"html/template"
	//"html/template"
	"log"
	"net/smtp"

	"../../Modelos/LoginModel"
	"../../Modelos/UsuarioModel"
	"../../Modulos/General"
	"../../Modulos/Session"
	"gopkg.in/kataras/iris.v6"
)

//LoginGet Renderea el Get del Login
func LoginGet(ctx *iris.Context) {
	data := ctx.Session().Get("data")
	if data != nil {
		ctx.Redirect("/Index", iris.StatusFound)
	}
	ctx.Render("Login.html", nil)
}

//LoginPost Renderea el Post del Login
func LoginPost(ctx *iris.Context) {
	var Send LoginModel.SLogin

	username := ctx.Request.FormValue("userName")
	userpass := ctx.Request.FormValue("userPass")

	if username != "" && userpass != "" {
		usuario := UsuarioModel.GetEspecificByFields("Usuario", username)
		if MoGeneral.EstaVacio(usuario) {
			Send.Login.IEstatus = true
			Send.Login.Ihtml = "Credenciales <strong>Invalidas!</strong>."
		} else {
			if usuario.Usuario == username && usuario.Credenciales.Contraseña == userpass {
				err := Session.CreateSession(usuario.Usuario, ctx)
				if err != nil {
					Send.SEstado = true
					Send.SMsj = err.Error()
				} else {
					ctx.Redirect("/Index", iris.StatusFound)
				}

			} else {
				Send.Login.IEstatus = true
				Send.Login.Ihtml = "Credenciales <strong>¡No!</strong>Validas."
			}
		}

	} else {
		Send.Login.IEstatus = true
		Send.Login.Ihtml = "Los campos <strong>Usuario</strong> y <strong>Contraseña</strong> son obligatorios "
	}
	ctx.Render("Login.html", Send)

}

//RecuperarGet Renderea al Get de Recuperar Contraseña
func RecuperarGet(ctx *iris.Context) {
	ctx.Render("RecuperarContraseña.html", nil)
}

//RecuperarPost Renderea al Post de Recuperar Contraseña
func RecuperarPost(ctx *iris.Context) {
	var correo string
	var username string
	var msg template.HTML
	correo = ctx.Request.FormValue("userCorreo")
	username = ctx.Request.FormValue("userName")
	if correo != "" && username != "" {
		//Buscar el nombre de usuario en la base de datos mongoDb (deberá ser exacto)
		usuariosEncontrados := UsuarioModel.ObtenerUsuariosMismoNombre(username)
		numeroUsuarios := len(usuariosEncontrados)
		if numeroUsuarios == 0 {
			msg = template.HTML("<div class='alert alert-danger alert-dismissable'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a><strong>¡Ups!</strong> Hemos detectado incongruencias en tus credenciales.</div>")
		} else if numeroUsuarios > 1 {
			msg = template.HTML("<div class='alert alert-danger alert-dismissable'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a><strong>¡Ups!</strong> Hemos detectado incongruencias en tus credenciales.</div>")
		} else {
			//verificar si el correo que ingreso, existe
			usuarioExistente := usuariosEncontrados[0].MediosDeContacto.Correos.Correos
			var mailFound bool
			for _, value := range usuarioExistente {
				correoEncontrado := value
				if correoEncontrado == correo {
					mailFound = true
				}
			}
			//Si el correo existe en el usuario, se envia el correo
			if mailFound {
				secreto, correo, nombreusuario, err := UsuarioModel.ExistMailandUpdate(correo, username)
				if err == nil {
					enviarCorreo(nombreusuario, secreto, correo)
					msg = template.HTML("<div class='alert alert-success alert-dismissable'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a><strong>¡Ok!</strong> Se ha enviado un correo a su cuenta.</div>")
				} else {
					fmt.Println("El correo que ingreso no existe...")
					msg = template.HTML("<div class='alert alert-danger alert-dismissable'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a><strong>¡Ups!</strong> No fue posible enviarle el correo.</div>")
				}
			} else {
				fmt.Println("Correo no existente en el usuario")
				msg = template.HTML("<div class='alert alert-danger alert-dismissable'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a><strong>¡Ups!</strong> Hemos detectado incongruencias en tus credenciales.</div>")
			}
		}
	} else {
		fmt.Println("No trae informacion.")
		msg = template.HTML("<div class='alert alert-danger alert-dismissable'><a href='#' class='close' data-dismiss='alert' aria-label='close'>&times;</a><strong>¡Ups!</strong> Debes identificarte.</div>")
	}
	ctx.Render("RecuperarContraseña.html", msg)
}

func enviarCorreo(nombreusuario, secreto, correo string) {
	body := "<p>Estimado usuario: <strong>" + nombreusuario + "</strong></p> \n <p>Su contrase&ntilde;a ha sido reestablecida al valor <strong>" + secreto + "</strong> .</p> \n Verifique su ingreso correcto."
	msg := "To: " + correo + "\r\n" +
		"Content-type: text/html" + "\r\n" +
		"Subject: Cambio de credenciales" + "\r\n\r\n" +
		body + "\r\n"
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		"test.miderp@hotmail.com",
		"D3m0st3n3s",
		"smtp.live.com",
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		"smtp.live.com:587",
		auth,
		"test.miderp@hotmail.com",
		[]string{correo},
		[]byte(msg),
	)
	//err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, password, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		log.Println(err)
		fmt.Println("Ha sucedido un error al procesar su solicitud,\n intente de nuevo.")
	} else {
		fmt.Println("Su solicitud ha sido procesada.")
	}
}

//NuevoGet Renderea al Get de Nuevo Login
func NuevoGet(ctx *iris.Context) {
	ctx.Render("NuevoLogin.html", nil)
}

//IndexGet Renderea al Get de Nuevo Login
func IndexGet(ctx *iris.Context) {
	var Send LoginModel.SLogin
	NameUsrLoged, MenuPrincipal, MenuUsr, errSes := Session.GetDataSession(ctx) //Retorna los datos de la session
	Send.SSesion.Name = NameUsrLoged
	Send.SSesion.MenuPrincipal = template.HTML(MenuPrincipal)
	Send.SSesion.MenuUsr = template.HTML(MenuUsr)
	if errSes != nil {
		Send.SEstado = false
		Send.SMsj = errSes.Error()
		ctx.Render("ZError.html", Send)
		return
	}

	ctx.Render("IndexDashboard.html", Send)

}

//Logout Cierra Sesion
func Logout(ctx *iris.Context) {
	var Send LoginModel.SLogin
	errSes := Session.DestroySession(ctx)
	if errSes != nil {
		Send.SEstado = false
		Send.SMsj = errSes.Error()
		ctx.Render("ZError.html", Send)
		return
	}
}

//GetMac funcion
func GetMac(ctx *iris.Context) {
	fmt.Println(ctx.FormValues())
	ctx.HTML(iris.StatusOK, "Se a recivido la  las direcciones correctamente ")
}
