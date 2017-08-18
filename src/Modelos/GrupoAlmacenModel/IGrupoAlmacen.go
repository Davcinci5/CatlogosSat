package GrupoAlmacenModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IGrupoAlmacen interface con los métodos de la clase
type IGrupoAlmacen interface {
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
func (p GrupoAlmacenMgo) InsertaMgo() bool {
	result := false
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}

	err = GrupoAlmacens.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p GrupoAlmacenMgo) InsertaElastic() bool {
	var GrupoAlmacenE GrupoAlmacenElastic

	GrupoAlmacenE.Nombre = p.Nombre
	GrupoAlmacenE.Descripcion = p.Descripcion
	// GrupoAlmacenE.PermiteVender = p.PermiteVender
	// GrupoAlmacenE.Miembros = p.Miembros
	// GrupoAlmacenE.Estatus = p.Estatus
	GrupoAlmacenE.FechaHora = p.FechaHora
	insert := MoConexion.InsertaElastic(MoVar.TipoGrupoAlmacen, p.ID.Hex(), GrupoAlmacenE)
	if !insert {
		fmt.Println("Error al insertar GrupoAlmacen en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p GrupoAlmacenMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = GrupoAlmacens.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p GrupoAlmacenMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupoAlmacen, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar GrupoAlmacen en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar GrupoAlmacen en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p GrupoAlmacenMgo) ReemplazaMgo() bool {
	result := false
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	err = GrupoAlmacens.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un GrupoAlmacen en elastic
func (p GrupoAlmacenMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupoAlmacen, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar GrupoAlmacen en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoGrupoAlmacen, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar GrupoAlmacen en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p GrupoAlmacenMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	n, e := GrupoAlmacens.Find(bson.M{field: valor}).Count()
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
func (p GrupoAlmacenMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	n, e := GrupoAlmacens.Find(bson.M{"_id": p.ID}).Count()
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
func (p GrupoAlmacenMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoGrupoAlmacen, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p GrupoAlmacenMgo) EliminaByIDMgo() bool {
	result := false
	s, GrupoAlmacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupoAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	e := GrupoAlmacens.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p GrupoAlmacenMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupoAlmacen, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar GrupoAlmacen en Elastic")
		return false
	}
	return true
}
