
package GrupoPersonaModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IGrupoPersona interface con los métodos de la clase
type IGrupoPersona interface {
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
func (p GrupoPersonaMgo) InsertaMgo() bool {
	result := false
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}

	err = GrupoPersonas.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p GrupoPersonaMgo) InsertaElastic() bool {
	insert := MoConexion.InsertaElastic(MoVar.TipoGrupoPersona, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al insertar GrupoPersona en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p GrupoPersonaMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = GrupoPersonas.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p GrupoPersonaMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupoPersona, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar GrupoPersona en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoGrupoPersona, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar GrupoPersona en Elastic")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p GrupoPersonaMgo) ReemplazaMgo() bool {
	result := false
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	err = GrupoPersonas.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un GrupoPersona en elastic
func (p GrupoPersonaMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupoPersona, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar GrupoPersona en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoGrupoPersona, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar GrupoPersona en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p GrupoPersonaMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	n, e := GrupoPersonas.Find(bson.M{field: valor}).Count()
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
func (p GrupoPersonaMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	n, e := GrupoPersonas.Find(bson.M{"_id": p.ID}).Count()
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
func (p GrupoPersonaMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoGrupoPersona, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p GrupoPersonaMgo) EliminaByIDMgo() bool {
	result := false
	s, GrupoPersonas, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoPersona)
	if err != nil {
		fmt.Println(err)
	}
	e := GrupoPersonas.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p GrupoPersonaMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupoPersona, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar GrupoPersona en Elastic")
		return false
	}
	return true
}
