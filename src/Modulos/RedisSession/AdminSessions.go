package RedisSession

import (
	"../Redis"
)

//ConjuntoUsuarios Nombre del conjunto para los usuarios del sistema
var ConjuntoUsuarios = "UsuariosDeSistema"

// EliminarSesionesActivas elimina todas las sesiones abiertas de un suaurio
func EliminarSesionesActivas(Username string) error {
	sesiones, err := Redis.ObtenerMiembrosdeGrupo(Username)

	if err != nil {
		return err
	}

	for _, val := range sesiones {
		err = EliminarSesion(Username, val)
	}
	return err
}

// RegresaSesionesActivasDe devuelve las sesiones activas de un usuario
func RegresaSesionesActivasDe(Username string) ([]map[string]string, error) {
	sesiones, err := Redis.ObtenerMiembrosdeGrupo(Username)
	var arr []map[string]string
	if err != nil {
		return arr, err
	}
	for _, vals := range sesiones {
		var m map[string]string
		m, err = Redis.RegresaValoresHash(vals)
		if err != nil {
			return arr, err
		}
		m["IDSession"] = vals
		arr = append(arr, m)
	}
	return arr, err
}

// RegresaTodasSesionesActivas devuelve las sesiones activas de un usuario
func RegresaTodasSesionesActivas() ([]map[string]string, error) {
	var arr []map[string]string
	var sesiones []string
	var err error

	usuariosDelSistema, errs := Redis.ObtenerMiembrosdeGrupo(ConjuntoUsuarios) //Obtiene todos los usuarios de tu sistema
	if errs != nil {
		return arr, errs
	}

	for _, usr := range usuariosDelSistema {
		SesionesDeUsr, err := Redis.ObtenerMiembrosdeGrupo(usr)
		if err != nil {
			return arr, err
		}
		for _, x := range SesionesDeUsr {
			sesiones = append(sesiones, x)
		}
	}

	if err != nil {
		return arr, err
	}
	for _, vals := range sesiones {
		var m map[string]string
		m, err = Redis.RegresaValoresHash(vals)
		if err != nil {
			return arr, err
		}
		m["IDSession"] = vals
		arr = append(arr, m)
	}
	return arr, err
}

//EliminaSesionByID elimina y cierra la sesion con el Id indicado
func EliminaSesionByID(ID string) error {
	usuariosDelSistema, errs := Redis.ObtenerMiembrosdeGrupo(ConjuntoUsuarios) //Obtiene todos los usuarios de tu sistema
	if errs != nil {
		return errs
	}
	var err error

	for _, usr := range usuariosDelSistema {
		err := EliminarSesion(usr, ID)
		if err != nil {
			return err
		}
	}
	return err
}

//AgregarUsuariosDelSistema agrega una lista de usuarios a Redis, al conjunto llamado 'Usuarios',
//Este conjunto puede ser configurado
func AgregarUsuariosDelSistema(usuarios []string) error {
	return Redis.InsertaRedis(ConjuntoUsuarios, usuarios)
}

//AgregarUsuario agrega un usuario al conjunto de Usuarios en Redis
func AgregarUsuario(UserName string) (bool, error) {
	return Redis.InsertarAConjunto(ConjuntoUsuarios, UserName)
}

//EliminarUsuario elimina un usuario del conjunto de Usuarios en Redis
func EliminarUsuario(Username string) error {
	return Redis.EliminaMiembroDeGrupo(ConjuntoUsuarios, Username)
}
