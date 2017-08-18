
package KitModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IKit interface con los métodos de la clase
type IKit interface {
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
func (p KitMgo) InsertaMgo() bool {
	result := false
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}

	err = Kits.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p KitMgo) InsertaElastic() bool {
	insert := MoConexion.InsertaElastic(MoVar.TipoKit, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al insertar Kit en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p KitMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Kits.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p KitMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoKit, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Kit en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoKit, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Kit en Elastic")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p KitMgo) ReemplazaMgo() bool {
	result := false
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	err = Kits.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Kit en elastic
func (p KitMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoKit, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Kit en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoKit, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Kit en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p KitMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Kits.Find(bson.M{field: valor}).Count()
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
func (p KitMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Kits.Find(bson.M{"_id": p.ID}).Count()
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
func (p KitMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoKit, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p KitMgo) EliminaByIDMgo() bool {
	result := false
	s, Kits, err := MoConexion.GetColectionMgo(MoVar.ColeccionKit)
	if err != nil {
		fmt.Println(err)
	}
	e := Kits.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p KitMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoKit, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Kit en Elastic")
		return false
	}
	return true
}
