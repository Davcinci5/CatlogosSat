package DispositivoModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IDispositivo interface con los métodos de la clase
type IDispositivo interface {
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
func (p DispositivoMgo) InsertaMgo() bool {
	result := false
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}

	err = Dispositivos.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p DispositivoMgo) InsertaElastic() bool {
	insert := MoConexion.InsertaElastic(MoVar.TipoDispositivo, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al insertar Dispositivo en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p DispositivoMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Dispositivos.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p DispositivoMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoDispositivo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Dispositivo en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoDispositivo, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Dispositivo en Elastic")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p DispositivoMgo) ReemplazaMgo() bool {
	result := false
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	err = Dispositivos.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Dispositivo en elastic
func (p DispositivoMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoDispositivo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Dispositivo en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoDispositivo, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Dispositivo en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p DispositivoMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Dispositivos.Find(bson.M{field: valor}).Count()
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
func (p DispositivoMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Dispositivos.Find(bson.M{"_id": p.ID}).Count()
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
func (p DispositivoMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoDispositivo, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p DispositivoMgo) EliminaByIDMgo() bool {
	result := false
	s, Dispositivos, err := MoConexion.GetColectionMgo(MoVar.ColeccionDispositivo)
	if err != nil {
		fmt.Println(err)
	}
	e := Dispositivos.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p DispositivoMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoDispositivo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Dispositivo en Elastic")
		return false
	}
	return true
}
