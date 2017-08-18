package ConexionModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IConexion interface con los métodos de la clase
type IConexion interface {
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
func (p ConexionMgo) InsertaMgo() bool {
	result := false
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}

	err = Conexions.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p ConexionMgo) InsertaElastic() bool {
	ConexionE := p.PreparaDatosELastic()
	insert := MoConexion.InsertaElastic(MoVar.TipoConexion, p.ID.Hex(), ConexionE)
	if !insert {
		fmt.Println("Error al insertar Conexion en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p ConexionMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Conexions.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p ConexionMgo) ActualizaElastic() error {
	ConexionE := p.PreparaDatosELastic()
	err := MoConexion.ActualizaElastic(MoVar.TipoConexion, p.ID.Hex(), ConexionE)
	return err
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p ConexionMgo) ReemplazaMgo() bool {
	result := false
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	err = Conexions.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Conexion en elastic
func (p ConexionMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoConexion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Conexion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoConexion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Conexion en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p ConexionMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Conexions.Find(bson.M{field: valor}).Count()
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
func (p ConexionMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Conexions.Find(bson.M{"_id": p.ID}).Count()
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
func (p ConexionMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoConexion, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p ConexionMgo) EliminaByIDMgo() bool {
	result := false
	s, Conexions, err := MoConexion.GetColectionMgo(MoVar.ColeccionConexion)
	if err != nil {
		fmt.Println(err)
	}
	e := Conexions.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		fmt.Println("Error al eliminar el almacen en Mongdb: ", p.ID, e)
		result = false
	} else {
		result = true
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p ConexionMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoConexion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Conexion en Elastic")
		return false
	}
	return true
}

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p ConexionMgo) PreparaDatosELastic() ConexionElastic {
	var ConexionE ConexionElastic
	ConexionE.Nombre = p.Nombre
	ConexionE.Servidor = p.Servidor
	ConexionE.NombreBD = p.NombreBD
	ConexionE.UsuarioBD = p.UsuarioBD
	ConexionE.FechaHora = p.FechaHora
	return ConexionE
}
