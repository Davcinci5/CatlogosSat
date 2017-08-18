
//#########################< MODELOS >#############################
//#########################< IUnidad.go >###############################
//################< ESTRUCTURA Y FUNCIONES >#######################

package UnidadModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IUnidad interface con los métodos de la clase
type IUnidad interface {
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
func (p UnidadMgo) InsertaMgo() bool {
	result := false
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}

	err = Unidads.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p UnidadMgo) InsertaElastic() bool {
	insert := MoConexion.InsertaElastic(MoVar.TipoUnidad, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al insertar Unidad en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p UnidadMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Unidads.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p UnidadMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoUnidad, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Unidad en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoUnidad, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Unidad en Elastic")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p UnidadMgo) ReemplazaMgo() bool {
	result := false
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	err = Unidads.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Unidad en elastic
func (p UnidadMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoUnidad, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Unidad en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoUnidad, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Unidad en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p UnidadMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Unidads.Find(bson.M{field: valor}).Count()
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
func (p UnidadMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Unidads.Find(bson.M{"_id": p.ID}).Count()
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
func (p UnidadMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoUnidad, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p UnidadMgo) EliminaByIDMgo() bool {
	result := false
	s, Unidads, err := MoConexion.GetColectionMgo(MoVar.ColeccionUnidad)
	if err != nil {
		fmt.Println(err)
	}
	e := Unidads.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p UnidadMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoUnidad, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Unidad en Elastic")
		return false
	}
	return true
}
