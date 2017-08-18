
package CajaModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//ICaja interface con los métodos de la clase
type ICaja interface {
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
func (p CajaMgo) InsertaMgo() bool {
	result := false
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}

	err = Cajas.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p CajaMgo) InsertaElastic() bool {
	var CajaE CajaElastic

	// CajaE.Usuario = p.Usuario
	// CajaE.Caja = p.Caja
	// CajaE.Cargo = p.Cargo
	// CajaE.Abono = p.Abono
	// CajaE.Saldo = p.Saldo
	// CajaE.Operacion = p.Operacion
	// CajaE.Estatus = p.Estatus
	// CajaE.FechaHora = p.FechaHora

	insert := MoConexion.InsertaElastic(MoVar.TipoCaja, p.ID.Hex(), CajaE)
	if !insert {
		fmt.Println("Error al insertar Caja en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p CajaMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Cajas.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p CajaMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCaja, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Caja en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Caja en Elastic, se perdió Referencia.")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p CajaMgo) ReemplazaMgo() bool {
	result := false
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	err = Cajas.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Caja en elastic
func (p CajaMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCaja, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Caja en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoCaja, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Caja en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p CajaMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Cajas.Find(bson.M{field: valor}).Count()
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
func (p CajaMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Cajas.Find(bson.M{"_id": p.ID}).Count()
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
func (p CajaMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoCaja, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p CajaMgo) EliminaByIDMgo() bool {
	result := false
	s, Cajas, err := MoConexion.GetColectionMgo(MoVar.ColeccionCaja)
	if err != nil {
		fmt.Println(err)
	}
	e := Cajas.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p CajaMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCaja, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Caja en Elastic")
		return false
	}
	return true
}
