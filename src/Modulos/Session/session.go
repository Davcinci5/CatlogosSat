package Session

import (
	"../RedisSession"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)

//LoginPage URI a donde se encuentra el login
var LoginPage = "/"

//IndexPage URI donde se encuentra el index del sistema
var IndexPage = "/Index"

//Datos estructura de variables de session
type Datos struct {
	ID       string
	Nombre   string
	MenuMain string
	MenuUsr  string
}

// CreateSession verifica que no exista una sesion, si no es asi crea una
func CreateSession(Name string, ctx *iris.Context) error {
	var err error
	ses := ctx.Session().Get("data")
	if ses == nil {
		data := Datos{}
		IDSesion := bson.NewObjectId()
		err = RedisSession.AgregarSesion(Name, IDSesion.Hex(), ctx)
		if err == nil {
			data.ID = IDSesion.Hex()
			data.Nombre = Name
			data.MenuMain = "<h1>Hola </h1> <h2> Menu Main </h2>"      // Recibe como cadena el template que vas a crear de tu menu (principal)
			data.MenuUsr = "<h1>Hola </h1> <h2> Menu De Usuario </h2>" // Recibe como cadena el template que vas a crear de tu menu (secundario)
			ctx.Session().Set("data", data)
		}
	} else {
		ctx.Redirect(IndexPage, iris.StatusFound)
	}
	return err
}

// GetDataSession obtiene los datos de la sesion
func GetDataSession(ctx *iris.Context) (Nombre, MenuMain, MenuUsuario string, Errores error) {
	data := ctx.Session().Get("data")
	var err error
	var permiso bool
	person := Datos{}

	if data != nil {
		person = data.(Datos)
		permiso, err = RedisSession.ComprobarSesion(person.Nombre, person.ID, ctx)
		if err != nil {
			return person.Nombre, person.MenuMain, person.MenuUsr, err
		}
		if !permiso {
			DestroySession(ctx)
		}
		return person.Nombre, person.MenuMain, person.MenuUsr, err
	}
	ctx.Redirect(LoginPage, iris.StatusFound)
	return
}

// DestroySession elimina la sesion existente
func DestroySession(ctx *iris.Context) error {
	data := ctx.Session().Get("data")
	if data != nil {
		person := data.(Datos)
		err := RedisSession.EliminarSesion(person.Nombre, person.ID)
		if err != nil {
			return err
		}
		ctx.Session().Delete("data")
		ctx.Redirect(LoginPage, iris.StatusFound)
	} else {
		ctx.Redirect(LoginPage, iris.StatusFound)
	}
	return nil
}

// GetInfoSession Regresa la informacion de todas las sesiones abiertas del usuario
func GetInfoSession(UserName string) ([]map[string]string, error) {
	return RedisSession.RegresaSesionesActivasDe(UserName)
}

//EliminarSesionesUsr Elimina todas las sesiones activas de un usuario
func EliminarSesionesUsr(UserName string) error {
	return RedisSession.EliminarSesionesActivas(UserName)
}

// RegresaSesionesActivas Retorna todas las sesiones activas en el sistema
func RegresaSesionesActivas() ([]map[string]string, error) {
	return RedisSession.RegresaTodasSesionesActivas()
}

//EliminaSesionByID Elimina una seesion indicada por un ID
func EliminaSesionByID(ID string) error {
	return RedisSession.EliminaSesionByID(ID)
}

//AgregarUsuariosDelSistema agrega todos los usuarios del sistema a redis, solo son los nombre de usuarios
func AgregarUsuariosDelSistema(usuariosDeSistema []string) error {
	return RedisSession.AgregarUsuariosDelSistema(usuariosDeSistema)
}

//AgregarUsuario agrega un usuario al conjunto de usuarios del sistema
func AgregarUsuario(UserName string) (bool, error) {
	return RedisSession.AgregarUsuario(UserName)
}

//EliminarUsuario elimina un usuario del conjunto de usuarios de sistema
func EliminarUsuario(UserName string) error {
	return RedisSession.EliminarUsuario(UserName)
}

//Funciones de Sesiones para respuestas AJAX

// GetDataSessionAJAX obtiene los datos de la sesion
func GetDataSessionAJAX(ctx *iris.Context) (SesionActiva bool, fn string, Nombre string, Errores error) {
	data := ctx.Session().Get("data")
	var err error
	var permiso bool
	var funcion string
	person := Datos{}

	if data != nil {
		person = data.(Datos)
		permiso, err = RedisSession.ComprobarSesion(person.Nombre, person.ID, ctx)
		if err != nil {
			return false, funcion, person.Nombre, err
		}
		if !permiso {
			DestroySessionAJAX(ctx)
			funcion = "location.reload();"
			return false, funcion, person.Nombre, err
		}
		return true, funcion, person.Nombre, err
	}
	funcion = "location.reload();"
	return false, funcion, person.Nombre, err

}

// DestroySessionAJAX elimina la sesion existente
func DestroySessionAJAX(ctx *iris.Context) error {
	data := ctx.Session().Get("data")
	if data != nil {
		person := data.(Datos)
		err := RedisSession.EliminarSesion(person.Nombre, person.ID)
		if err != nil {
			return err
		}
		ctx.Session().Delete("data")
	}
	return nil
}
