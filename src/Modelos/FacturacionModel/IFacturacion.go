package FacturacionModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IFacturacion interface con los métodos de la clase
type IFacturacion interface {
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
func (p FacturacionMgo) InsertaMgo() bool {
	result := false
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}

	err = Facturacions.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p FacturacionMgo) InsertaElastic() bool {
	insert := MoConexion.InsertaElastic(MoVar.TipoFacturacion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al insertar Facturacion en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p FacturacionMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Facturacions.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p FacturacionMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoFacturacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Facturacion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoFacturacion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Facturacion en Elastic")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p FacturacionMgo) ReemplazaMgo() bool {
	result := false
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	err = Facturacions.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Facturacion en elastic
func (p FacturacionMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoFacturacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Facturacion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoFacturacion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Facturacion en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p FacturacionMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Facturacions.Find(bson.M{field: valor}).Count()
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
func (p FacturacionMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Facturacions.Find(bson.M{"_id": p.ID}).Count()
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
func (p FacturacionMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoFacturacion, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p FacturacionMgo) EliminaByIDMgo() bool {
	result := false
	s, Facturacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionFacturacion)
	if err != nil {
		fmt.Println(err)
	}
	e := Facturacions.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p FacturacionMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoFacturacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Facturacion en Elastic")
		return false
	}
	return true
}
