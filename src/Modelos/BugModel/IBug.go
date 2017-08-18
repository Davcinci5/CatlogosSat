package BugModel

import (
	"fmt"

	"../../Modelos/CatalogoModel"

	"gopkg.in/mgo.v2/bson"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
)

// CatalogoEstatusBugs estatus existentes para los bugs
var CatalogoEstatusBugs = 181

//IBug interface con los métodos de la clase
type IBug interface {
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
func (p BugMgo) InsertaMgo() bool {
	result := false
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}

	err = Bugs.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p BugMgo) InsertaElastic() bool {
	var BugE BugElastic

	BugE.Tipo = p.Tipo
	BugE.Titulo = p.Titulo
	BugE.Descripcion = p.Descripcion
	BugE.Usuario = p.Usuario
	BugE.Metodo = p.Metodo
	BugE.EsAjax = p.EsAjax
	BugE.EstatusPeticion = p.EstatusPeticion
	//BugE.Estatus = p.Estatus
	BugE.Estatus = CatalogoModel.GetValorMagnitud(p.Estatus, CatalogoEstatusBugs)

	BugE.Ruta = p.Ruta
	BugE.FechaHora = p.FechaHora
	insert := MoConexion.InsertaElastic(MoVar.TipoBug, p.ID.Hex(), BugE)
	if !insert {
		fmt.Println("Error al insertar Bug en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p BugMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Bugs.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p BugMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoBug, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Bug en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Bug en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p BugMgo) ReemplazaMgo() bool {
	result := false
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	err = Bugs.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Bug en elastic
func (p BugMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoBug, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Bug en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoBug, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Bug en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p BugMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Bugs.Find(bson.M{field: valor}).Count()
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
func (p BugMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Bugs.Find(bson.M{"_id": p.ID}).Count()
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
func (p BugMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoBug, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p BugMgo) EliminaByIDMgo() bool {
	result := false
	s, Bugs, err := MoConexion.GetColectionMgo(MoVar.ColeccionBug)
	if err != nil {
		fmt.Println(err)
	}
	e := Bugs.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p BugMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoBug, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Bug en Elastic")
		return false
	}
	return true
}
