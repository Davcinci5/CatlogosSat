package UsuarioModel

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/EquipoCajaModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//UsuarioMgo estructura de Usuarios mongo
type UsuarioMgo struct {
	ID               bson.ObjectId          `bson:"_id,omitempty"`
	IDPersona        bson.ObjectId          `bson:"IDPersona,omitempty"`
	Usuario          string                 `bson:"Usuario"`
	Credenciales     CredencialesUsuarioMgo `bson:"Credenciales"`
	Roles            []bson.ObjectId        `bson:"Roles"`
	MediosDeContacto MediosContactoMgo      `bson:"MediosDeContacto"`
	Cajas            []bson.ObjectId        `bson:"Cajas"`
	Notificaciones   []NotificacionMgo      `bson:"Notificaciones"`
	Estatus          bson.ObjectId          `bson:"Estatus,omitempty"`
	FechaHora        time.Time              `bson:"FechaHora"`
	IsAdmin          bool                   `bson:"Administrador"`
}

//CredencialesUsuarioMgo Estructura para definir las credenciales del usuario Mongo
type CredencialesUsuarioMgo struct {
	Pin         string `bson:"Pin"`
	Contraseña  string `bson:"Contraseña"`
	CodigoBarra string `bson:"CodigoBarra"`
	Huella      string `bson:"Huella"`
}

//UsuarioElastic estructura de Usuarios para insertar en Elastic
type UsuarioElastic struct {
	Persona string `json:"Persona,omitempty"`
	Usuario string `json:"Usuario"`
	//Credenciales []string  `json:"Credenciales"`
	Roles     []string  `json:"Roles"`
	Correos   []string  `json:"Correos"`
	Telefonos []string  `json:"Telefonos"`
	Otros     []string  `json:"Otros"`
	Cajas     []string  `json:"Cajas"`
	Estatus   string    `json:"Estatus"`
	FechaHora time.Time `json:"FechaHora"`
}

// MediosContactoElastic subestructura de Usuario
// type MediosContactoElastic struct {
// 	Correos   []string `json:"Correos,omitempty"`
// 	Telefonos []string `json:"Telefonos,omitempty"`
// 	Otros     []string `json:"Otros,omitempty"`
// }

//MediosContactoMgo subestructura de Usuario
type MediosContactoMgo struct {
	Correos   CorreoMgo   `bson:"Correos"`
	Telefonos TelefonoMgo `bson:"Telefonos"`
	Otros     []string    `bson:"Otros"`
}

//CorreoMgo subestructura de Usuario
type CorreoMgo struct {
	Principal string   `bson:"Principal"`
	Correos   []string `bson:"Correos"`
}

//TelefonoMgo subestructura de Usuario
type TelefonoMgo struct {
	Principal string   `bson:"Principal"`
	Telefonos []string `bson:"Telefonos"`
}

//NotificacionMgo subestructura de Usuario
type NotificacionMgo struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	Mensaje          string        `bson:"Mensaje"`
	Leido            bool          `bson:"Leido"`
	FechaOcurrencia  time.Time     `bson:"FechaOcurrencia"`
	FechaVencimiento time.Time     `bson:"FechaVencimiento"`
}

//PersonaMgo estructura de Personas mongo
type PersonaMgo struct {
	ID           bson.ObjectId   `bson:"_id,omitempty"`
	Nombre       string          `bson:"Nombre"`
	Tipo         []bson.ObjectId `bson:"Tipo,omitempty"`
	Grupos       []bson.ObjectId `bson:"Grupos,omitempty"`
	Predecesor   bson.ObjectId   `bson:"Predecesor,omitempty"`
	Notificacion []bson.ObjectId `bson:"Notificacion,omitempty"`
	Estatus      bson.ObjectId   `bson:"Estatus,omitempty"`
	FechaHora    time.Time       `bson:"FechaHora,omitempty"`
}

//PersonaElastic subestructura de Usuario
type PersonaElastic struct {
	Nombre       string    `json:"Nombre"`
	Tipo         string    `json:"Tipo,omitempty"`
	Grupos       string    `json:"Grupos,omitempty"`
	Predecesor   string    `json:"Predecesor,omitempty"`
	Notificacion string    `json:"Notificacion,omitempty"`
	Estatus      string    `json:"Estatus,omitempty"`
	FechaHora    time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []UsuarioMgo {
	var result []UsuarioMgo
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	err = Usuarios.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Usuarios.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) UsuarioMgo {
	var result UsuarioMgo
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	err = Usuarios.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []UsuarioMgo {
	var result []UsuarioMgo
	var aux UsuarioMgo
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = UsuarioMgo{}
		Usuarios.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) UsuarioMgo {
	var result UsuarioMgo
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)

	if err != nil {
		fmt.Println(err)
	}
	err = Usuarios.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//ObtenerUsuariosMismoNombre regresa todos los usuarios que tienen el mismo nombre
func ObtenerUsuariosMismoNombre(usuario string) []UsuarioMgo {
	var result []UsuarioMgo
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	err = Usuarios.Find(bson.M{"Usuario": usuario}).All(&result)
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result UsuarioMgo
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	err = Usuarios.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboUsuarios regresa un combo de Usuario de mongo
func CargaComboUsuarios(ID string) string {
	Usuarios := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Usuarios {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Usuario + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Usuario + ` </option> `
		}

	}
	return templ
}

// CargaNombreCajas de un arreglo de IDS regresa el nombre de cada caja que corresponde cada ID
func CargaNombreCajas(Cajas []bson.ObjectId) []string {
	return EquipoCajaModel.CargaNombresCajas(Cajas)
}

// CargaNombreEstatus recibe un ID y regresa el nombre del tipo de estatus de usuario
func CargaNombreEstatus(ID bson.ObjectId) string {
	return CatalogoModel.GetValorMagnitud(ID, 167)
}

type InitConf struct {
	ID       bson.ObjectId   `bson:"_id,omitempty"`
	Nombre   string          `bson:"Nombre"`
	Terceros []string        `bson:"Terceros"`
	Intentos int             `bson:"Intentos"`
	Minutos  int             `bson:"Minutos"`
	OldPass  bool            `bson:"OldPass"`
	Opciones []bson.ObjectId `bson:"Opciones"`
}

type Opciones struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Nombre       string        `bson:"Nombre"`
	Longitud     int           `bson:"Longitud"`
	Comprobacion bool          `bson:"Comprobacion"`
}

//CargarOpcionesUsuario Carga las Opciones de la Seguridad de Acceso al sistema
func CargarOpcionesUsuario() []Opciones {

	var result InitConf
	var option Opciones
	var options []Opciones
	s, Configuracion, err := MoConexion.GetColectionMgo("InitConf")
	defer s.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = Configuracion.Find(bson.M{"Nombre": "Configuracion01"}).One(&result)

	for _, v := range result.Opciones {
		err = Configuracion.Find(bson.M{"_id": v}).One(&option)
		options = append(options, option)
	}

	return options
}

//CargarTemplateCredenciales Cargara el template para el alta de usuarios, en el area de opciones de seguridad de logeo
func CargarTemplateCredenciales(mapavalores map[string]string, mapaboleano map[string]bool) string {
	var result InitConf
	var option Opciones
	var options []Opciones
	s, Configuracion, err := MoConexion.GetColectionMgo("InitConf")
	defer s.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = Configuracion.Find(bson.M{"Nombre": "Configuracion01"}).One(&result)

	for _, v := range result.Opciones {
		err = Configuracion.Find(bson.M{"_id": v}).One(&option)
		options = append(options, option)
	}

	var tmpl = ``

	for _, x := range options {

		if x.Nombre == "PIN" {

			if mapaboleano == nil && mapavalores == nil {

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
				} else {
					tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
				}

				tmpl += fmt.Sprintf(`<div class="input-group">`)
				tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
				tmpl += fmt.Sprintf(`</span>`)
				tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
				tmpl += fmt.Sprintf(`</div>`)
				tmpl += fmt.Sprintf(`</div>`)

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)
				}
			} else if mapaboleano != nil && mapavalores != nil {

				if mapaboleano[x.Nombre] == false {
					//Dividir la string que viene el el mapa de valores para pintar 1 y 1
					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}

				} else if mapaboleano[x.Nombre] == true {
					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadena = mapavalores[x.Nombre]
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {

						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}
				}
			}
		}

		if x.Nombre == "CONTRASEÑA" {
			if mapaboleano == nil && mapavalores == nil {

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
				} else {
					tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
				}

				tmpl += fmt.Sprintf(`<div class="input-group">`)
				tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
				tmpl += fmt.Sprintf(`</span>`)
				tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
				tmpl += fmt.Sprintf(`</div>`)
				tmpl += fmt.Sprintf(`</div>`)

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)
				}
			} else if mapaboleano != nil && mapavalores != nil {

				if mapaboleano[x.Nombre] == false {
					//Dividir la string que viene el el mapa de valores para pintar 1 y 1
					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}

				} else if mapaboleano[x.Nombre] == true {

					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadena = mapavalores[x.Nombre]
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {

						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}
				}
			}
		}

		if x.Nombre == "HUELLA" {
			if mapaboleano == nil && mapavalores == nil {

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
				} else {
					tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
				}

				tmpl += fmt.Sprintf(`<div class="input-group">`)
				tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
				tmpl += fmt.Sprintf(`</span>`)
				tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
				tmpl += fmt.Sprintf(`</div>`)
				tmpl += fmt.Sprintf(`</div>`)

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)
				}
			} else if mapaboleano != nil && mapavalores != nil {

				if mapaboleano[x.Nombre] == false {
					//Dividir la string que viene el el mapa de valores para pintar 1 y 1
					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}

				} else if mapaboleano[x.Nombre] == true {

					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadena = mapavalores[x.Nombre]
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {

						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}
				}
			}
		}

		if x.Nombre == "CODIGOBARRA" {
			if mapaboleano == nil && mapavalores == nil {

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
				} else {
					tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
				}

				tmpl += fmt.Sprintf(`<div class="input-group">`)
				tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
				tmpl += fmt.Sprintf(`</span>`)
				tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
				tmpl += fmt.Sprintf(`</div>`)
				tmpl += fmt.Sprintf(`</div>`)

				if x.Comprobacion == true {
					tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required >`, x.Nombre, x.Nombre, x.Longitud)
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)
				}
			} else if mapaboleano != nil && mapavalores != nil {

				if mapaboleano[x.Nombre] == false {
					//Dividir la string que viene el el mapa de valores para pintar 1 y 1
					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}

				} else if mapaboleano[x.Nombre] == true {

					var cadena string
					var cadenas []string

					if x.Comprobacion == true && mapaboleano[x.Nombre] == false {
						cadena = mapavalores[x.Nombre]
						cadenas = strings.Split(cadena, "|")

					} else {
						cadena = mapavalores[x.Nombre]
						cadenas = append(cadenas, cadena)
						cadenas = append(cadenas, cadena)
					}

					if x.Comprobacion == true {
						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
					} else {
						tmpl += fmt.Sprintf(`<div class="col-sm-10">`)
					}
					tmpl += fmt.Sprintf(`<div class="input-group">`)
					tmpl += fmt.Sprintf(`<span class="input-group-addon">%v`, x.Nombre)
					tmpl += fmt.Sprintf(`</span>`)
					tmpl += fmt.Sprintf(`<input type="text" name="%v1" id="%v1" class="form-control" placeholder="%v digitos" required value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[0])
					tmpl += fmt.Sprintf(`</div>`)
					tmpl += fmt.Sprintf(`</div>`)

					if x.Comprobacion == true {

						tmpl += fmt.Sprintf(`<div class="col-sm-6">`)
						tmpl += fmt.Sprintf(`<div class="input-group">`)
						tmpl += fmt.Sprintf(`<span class="input-group-addon">Comprobar`)
						tmpl += fmt.Sprintf(`</span>`)
						tmpl += fmt.Sprintf(`<input type="text" name="%v2" id="%v2" class="form-control" placeholder="%v digitos" required  value="%v">`, x.Nombre, x.Nombre, x.Longitud, cadenas[1])
						tmpl += fmt.Sprintf(`</div>`)
						tmpl += fmt.Sprintf(`</div>`)
					}
				}
			}
		}

	}
	return tmpl
}

//CargaMediosDeContacto Regresa el template que carga los medios de contacto
// func CargaMediosDeContacto(contactos []MediosContacto) string {

// 	tmpl := ``

// 	tmpl += fmt.Sprintf(`
// 		<div class="form-group">
// 			<label class="col-sm-4 control-label" for="MediosDeContacto">Medios de Contacto:</label>
// 			<div class="col-sm-5">
// 	`)

// 	for _, value := range contactos {

// 		fmt.Println(value)

// 		tmpl += fmt.Sprintf(`
// 				<div class="col-sm-6">
// 					<div class="input-group">
// 						<span class="input-group-addon">
// 				`)

// 		if value.Correos.Principal {
// 			tmpl += fmt.Sprintf(`
// 							<input type="radio" id="PrincipalCorreo" name="PrincipalCorreo" class="CheckPrincipalCorreo" disabled checked>
// 				`)

// 		} else {

// 			tmpl += fmt.Sprintf(`
// 							<input type="radio" id="PrincipalCorreo" name="PrincipalCorreo" class="CheckPrincipalCorreo" disabled>
// 				`)
// 		}

// 		tmpl += fmt.Sprintf(`
// 						</span>
// 						<input type="text" name="Correos" id="Correos" class="form-control inputCorreos" placeholder="Correo Electronico"  value="%v" readonly>
// 					</div>
// 				</div>
// 				`, value.Correos.Email)

// 		tmpl += fmt.Sprintf(`

// 				<div class="col-sm-6">
// 					<div class="input-group">
// 						<span class="input-group-addon">
// 				`)
// 		if value.Telefonos.Principal {
// 			tmpl += fmt.Sprintf(`
// 							<input type="radio" id="PrincipalTelefono" name="PrincipalTelefono" class="CheckPrincipalTelefono" disabled checked>
// 				`)

// 		} else {

// 			tmpl += fmt.Sprintf(`
// 							<input type="radio" id="PrincipalTelefono" name="PrincipalTelefono" class="CheckPrincipalTelefono" disabled>
// 				`)
// 		}

// 		tmpl += fmt.Sprintf(`
// 						</span>
// 						<input type="text" name="Telefonos" id="Telefonos" class="form-control inputTelefono" placeholder="Telefono"  value="%v" readonly>
// 					</div>
// 				</div>

// 				`, value.Telefonos.Telefono)

// 	}

// 	tmpl += fmt.Sprintf(`
// 			</div>
// 		</div>
// 	`)

// 	tmpl += fmt.Sprintf(``)

// 	return tmpl

// }

func ComprobarUsuario(usuario string, password string) (*UsuarioMgo, bool) {
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	defer s.Close()
	var result UsuarioMgo
	bisonDocSearch := bson.M{"Usuario": usuario, "Credenciales.Contraseña": password}
	err = Usuarios.Find(bisonDocSearch).Select(bson.M{"_id": 1, "Usuario": 1, "Administrador": 1}).One(&result)
	if err != nil {
		return nil, false
	} else {
		return &result, true
	}

}

func ExistMailandUpdate(mail string, username string) (string, string, string, error) {
	s, Usuarios, err1 := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err1 != nil {
		fmt.Println(err1)
	}
	defer s.Close()

	var result UsuarioMgo
	var err error
	var correoprincipal string

	var bisonDocSearch bson.M

	if mail != "" {
		bisonDocSearch = bson.M{"MediosDeContacto.Correos.Email": mail}
	}

	if username != "" {
		bisonDocSearch = bson.M{"Usuario": username}
	}

	err = Usuarios.Find(bisonDocSearch).Select(bson.M{"_id": 1, "Usuario": 1, "Administrador": 1, "MediosDeContacto": 1, "Credenciales": 1}).One(&result)

	if err == nil {
		secreto := RandStringRunes(6)
		// hasher := md5.New()
		// hasher.Write([]byte(secreto))
		// secretomd5 := hex.EncodeToString(hasher.Sum(nil))

		// for i, v := range result.Credenciales {
		// 	var credsplit []string
		// 	credsplit = strings.Split(v, ":")
		// 	for _, vv := range credsplit {
		// 		if vv == "CONTRASEÑA" {
		// 			result.Credenciales[i] = "CONTRASEÑA:" + secreto
		// 		}
		// 	}
		// }
		result.Credenciales.Contraseña = secreto
		bisonDocFilter := bson.M{"$set": bson.M{"Credenciales": result.Credenciales}}
		bisonDocSearch = bson.M{"_id": result.ID}
		_, err = Usuarios.UpdateAll(bisonDocSearch, bisonDocFilter)

		// if mail == "" {
		// 	for _, valor := range result.MediosDeContacto {
		// 		if valor.Correos.Principal {
		// 			correoprincipal = valor.Correos.Email
		// 		}
		// 	}
		// } else {
		// 	correoprincipal = mail

		// }
		correoprincipal = result.MediosDeContacto.Correos.Principal
		return secreto, correoprincipal, result.Usuario, err
	}
	return "", "", "", err
}

func RandStringRunes(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("IDPersona")
	queryQuotes = queryQuotes.Field("IDPersona")

	queryTilde = queryTilde.Field("Usuario")
	queryQuotes = queryQuotes.Field("Usuario")

	queryTilde = queryTilde.Field("Credenciales")
	queryQuotes = queryQuotes.Field("Credenciales")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoUsuario, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoUsuario, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}
