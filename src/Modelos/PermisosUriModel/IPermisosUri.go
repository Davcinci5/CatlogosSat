
package PermisosUriModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IPermisosUri interface con los métodos de la clase
type IPermisosUri interface {
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
func (p PermisosUriMgo) InsertaMgo() bool {
	result := false
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}

	err = PermisosUris.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p PermisosUriMgo) InsertaElastic() bool {
	var PermisosUriE PermisosUriElastic

	PermisosUriE.Grupo = p.Grupo
	PermisosUriE.PermisoNegado = p.PermisoNegado
	PermisosUriE.PermisoAceptado = p.PermisoAceptado
	insert := MoConexion.InsertaElastic(MoVar.TipoPermisosUri, p.ID.Hex(), PermisosUriE)
	if !insert {
		fmt.Println("Error al insertar PermisosUri en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p PermisosUriMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = PermisosUris.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p PermisosUriMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPermisosUri, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar PermisosUri en Elastic")
		return false
	}

	if !p.InsertaElastic(){
		fmt.Println("Error al actualizar PermisosUri en Elastic, se perdió Referencia.")
		return false
	}
	
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p PermisosUriMgo) ReemplazaMgo() bool {
	result := false
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	err = PermisosUris.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un PermisosUri en elastic
func (p PermisosUriMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPermisosUri, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar PermisosUri en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoPermisosUri, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar PermisosUri en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p PermisosUriMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}
	n, e := PermisosUris.Find(bson.M{field: valor}).Count()
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
func (p PermisosUriMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}
	n, e := PermisosUris.Find(bson.M{"_id": p.ID}).Count()
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
func (p PermisosUriMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoPermisosUri, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p PermisosUriMgo) EliminaByIDMgo() bool {
	result := false
	s, PermisosUris, err := MoConexion.GetColectionMgo(MoVar.ColeccionPermisosUri)
	if err != nil {
		fmt.Println(err)
	}
	e := PermisosUris.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p PermisosUriMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPermisosUri, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar PermisosUri en Elastic")
		return false
	}
	return true
}
