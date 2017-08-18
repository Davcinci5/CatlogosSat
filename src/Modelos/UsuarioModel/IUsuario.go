package UsuarioModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IUsuario interface con los métodos de la clase
type IUsuario interface {
	InsertaMgo() bool
	InsertaElastic() bool

	ActualizaMgo(campos []string, valores []interface{}) bool
	ActualizaElastic(campos []string, valores []interface{}) bool //Reemplaza No Actualiza

	ReemplazaMgo() bool
	ReemplazaElastic() bool

	ConsultaExistenciaByFieldMgo(field string, valor string)

	ConsultaExistenciaByIDMgo() bool
	ConsultaExistenciaByIDElastic() bool

	EliminaByIDMgo() bool
	EliminaByIDElastic() bool
}

//################################################<<METODOS DE GESTION >>################################################################

//##################################<< INSERTAR >>###################################

//InsertaMgo es un método que crea un registro en Mongo
func (p UsuarioMgo) InsertaMgo() bool {
	result := false
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}

	err = Usuarios.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p UsuarioMgo) InsertaElastic() bool {
	UsuarioE := p.PreparaDatosELastic()
	insert := MoConexion.InsertaElastic(MoVar.TipoUsuario, p.ID.Hex(), UsuarioE)
	if !insert {
		fmt.Println("Error al insertar Usuario en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p UsuarioMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Usuarios.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p UsuarioMgo) ActualizaElastic() error {
	UsuarioE := p.PreparaDatosELastic()
	err := MoConexion.ActualizaElastic(MoVar.TipoUsuario, p.ID.Hex(), UsuarioE)
	return err
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p UsuarioMgo) ReemplazaMgo() bool {
	result := false
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	err = Usuarios.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Usuario en elastic
func (p UsuarioMgo) ReemplazaElastic() bool {
	var UsuarioE UsuarioElastic
	UsuarioE.Persona = p.ID.Hex()
	UsuarioE.Usuario = p.Usuario
	UsuarioE.Correos = p.MediosDeContacto.Correos.Correos
	UsuarioE.Telefonos = p.MediosDeContacto.Telefonos.Telefonos
	UsuarioE.Otros = p.MediosDeContacto.Otros
	UsuarioE.Cajas = CargaNombreCajas(p.Cajas)
	UsuarioE.Estatus = CargaNombreEstatus(p.Estatus)
	UsuarioE.FechaHora = p.FechaHora
	fmt.Println("Usuario Elastic: ", UsuarioE)

	delete := MoConexion.DeleteElastic(MoVar.TipoUsuario, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Usuario en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoUsuario, p.ID.Hex(), UsuarioE)
	if !insert {
		fmt.Println("Error al actualizar Usuario en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p UsuarioMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Usuarios.Find(bson.M{field: valor}).Count()
	if e != nil {
		fmt.Println(e)
	}
	if n > 0 {
		result = true
	}
	s.Close()
	return result
}

//ConsultaExistenciaByIDMgo es un método que encuentra un registro en Mongo buscándolo por ID
func (p UsuarioMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Usuarios.Find(bson.M{"_id": p.ID}).Count()
	if e != nil {
		fmt.Println(e)
	}
	if n > 0 {
		result = true
	}
	s.Close()
	return result
}

//ConsultaExistenciaByIDElastic es un método que encuentra un registro en Mongo buscándolo por ID
func (p UsuarioMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoUsuario, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p UsuarioMgo) EliminaByIDMgo() bool {
	result := false
	s, Usuarios, err := MoConexion.GetColectionMgo(MoVar.ColeccionUsuario)
	if err != nil {
		fmt.Println(err)
	}
	e := Usuarios.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p UsuarioMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoUsuario, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Usuario en Elastic")
		return false
	}
	return true
}

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p UsuarioMgo) PreparaDatosELastic() UsuarioElastic {
	var UsuarioE UsuarioElastic
	UsuarioE.Persona = p.IDPersona.Hex()
	UsuarioE.Usuario = p.Usuario
	UsuarioE.Correos = p.MediosDeContacto.Correos.Correos
	UsuarioE.Telefonos = p.MediosDeContacto.Telefonos.Telefonos
	UsuarioE.Otros = p.MediosDeContacto.Otros
	UsuarioE.Cajas = CargaNombreCajas(p.Cajas)
	UsuarioE.Estatus = CargaNombreEstatus(p.Estatus)
	UsuarioE.FechaHora = p.FechaHora
	return UsuarioE
}
