package RolModel

import (
	"fmt"

	"../../Modelos/CatalogoModel"
	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IRol interface con los métodos de la clase
type IRol interface {
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
func (p RolMgo) InsertaMgo() bool {
	result := false
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}

	err = Rols.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p RolMgo) InsertaElastic() bool {
	var RolE RolElastic

	RolE.Nombre = p.Nombre
	RolE.Descripcion = p.Descripcion
	RolE.Permisos = CatalogoModel.RegresaNombrePermisos(p.Permisos)
	RolE.Estatus = CatalogoModel.RegresaNombreSubCatalogo(p.Estatus)
	RolE.FechaHora = p.FechaHora

	insert := MoConexion.InsertaElastic(MoVar.TipoRol, p.ID.Hex(), RolE)
	if !insert {
		fmt.Println("Error al insertar Rol en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p RolMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Rols.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p RolMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoRol, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Rol en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Rol en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p RolMgo) ReemplazaMgo() bool {
	result := false
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	err = Rols.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Rol en elastic
func (p RolMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoRol, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Rol en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoRol, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Rol en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p RolMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Rols.Find(bson.M{field: valor}).Count()
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
func (p RolMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Rols.Find(bson.M{"_id": p.ID}).Count()
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
func (p RolMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoRol, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p RolMgo) EliminaByIDMgo() bool {
	result := false
	s, Rols, err := MoConexion.GetColectionMgo(MoVar.ColeccionRol)
	if err != nil {
		fmt.Println(err)
	}
	e := Rols.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p RolMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoRol, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Rol en Elastic")
		return false
	}
	return true
}
