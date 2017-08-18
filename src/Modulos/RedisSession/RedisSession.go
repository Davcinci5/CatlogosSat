package RedisSession

import (
	"time"

	"../Redis"
	"gopkg.in/kataras/iris.v6"
)

//AgregarSesion agrega una session a un usuario referenciada por un id por cada
// nueva sesion, si existen otras, agrega el id al grupo de sesiones del usuario
func AgregarSesion(Username string, IDSesion string, ctx *iris.Context) error {
	f := map[string]interface{}{
		"IP":              ctx.RemoteAddr(),
		"ServerHost":      ctx.ServerHost(),
		"UserAgent":       ctx.Request.UserAgent(),
		"VirtualHostName": ctx.VirtualHostname(),
		"Ubicacion":       ctx.Request.RequestURI,
		"Time":            time.Now().String(),
	}

	err := Redis.CrearHash(IDSesion, f)

	if err == nil {
		_, err = Redis.InsertarAConjunto(Username, IDSesion)
	}

	return err
}

// ComprobarSesion verifica que el ID de la sesion pertenezca al usuario, es decir, que este activa
func ComprobarSesion(Username string, ID string, ctx *iris.Context) (bool, error) {
	exist, err := Redis.ObtenerMiemmbroDeGrupo(Username, ID)
	if err == nil {
		if exist {
			f := map[string]interface{}{
				"Ubicacion": ctx.Request.RequestURI,
			}
			Redis.ActualizarHash(ID, f)
		}
	}
	return exist, err
}

//EliminarSesion Elimina una sesion de un usuario, sacandolo del sistema
func EliminarSesion(Username, ID string) error {
	arr := []string{"IP", "ServerHost", "UserAgent", "VirtualHostName", "Time", "Ubicacion"}
	err := Redis.EliminarHash(ID, arr)

	if err == nil {
		err = Redis.EliminaMiembroDeGrupo(Username, ID)
	}
	return err
}
