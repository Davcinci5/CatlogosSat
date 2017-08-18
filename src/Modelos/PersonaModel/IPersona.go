package PersonaModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IPersona interface con los métodos de la clase
type IPersona interface {
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
func (p PersonaMgo) InsertaMgo() bool {
	result := false
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}

	err = Personas.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p PersonaMgo) InsertaElastic() bool {
	PersonaE := p.PreparaDatosELastic()
	insert := MoConexion.InsertaElastic(MoVar.TipoPersona, p.ID.Hex(), PersonaE)
	if !insert {
		fmt.Println("Error al insertar Persona en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p PersonaMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Personas.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p PersonaMgo) ActualizaElastic() error {
	PersonaE := p.PreparaDatosELastic()
	err := MoConexion.ActualizaElastic(MoVar.ColeccionPersona, p.ID.Hex(), PersonaE)
	return err
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p PersonaMgo) ReemplazaMgo() bool {
	result := false
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	err = Personas.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Persona en elastic
func (p PersonaMgo) ReemplazaElastic() bool {
	var PersonaE PersonaElastic

	PersonaE.Nombre = p.Nombre
	PersonaE.Tipo = CargaNombresTiposPersonas(p.Tipo)
	PersonaE.Grupos = CargaNombresGruposPersonas(p.Grupos)
	PersonaE.Predecesor = NombrePredecesor(p.Predecesor)
	PersonaE.Estatus = CargaNombreEstatusGrupo(p.Estatus)
	PersonaE.FechaHora = p.FechaHora

	fmt.Println("Persona Elastic", PersonaE)

	delete := MoConexion.DeleteElastic(MoVar.TipoPersona, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Persona en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoPersona, p.ID.Hex(), PersonaE)
	if !insert {
		fmt.Println("Error al actualizar Persona en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p PersonaMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Personas.Find(bson.M{field: valor}).Count()
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
func (p PersonaMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Personas.Find(bson.M{"_id": p.ID}).Count()
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
func (p PersonaMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoPersona, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p PersonaMgo) EliminaByIDMgo() bool {
	result := false
	s, Personas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPersona)
	if err != nil {
		fmt.Println(err)
	}
	e := Personas.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p PersonaMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPersona, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Persona en Elastic")
		return false
	}
	return true
}

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p PersonaMgo) PreparaDatosELastic() PersonaElastic {
	var PersonaE PersonaElastic
	PersonaE.Nombre = p.Nombre
	PersonaE.Tipo = CargaNombresTiposPersonas(p.Tipo)
	PersonaE.Grupos = CargaNombresGruposPersonas(p.Grupos)
	PersonaE.Predecesor = NombrePredecesor(p.Predecesor)
	PersonaE.Estatus = CargaNombreEstatusGrupo(p.Estatus)
	PersonaE.FechaHora = p.FechaHora
	return PersonaE
}
