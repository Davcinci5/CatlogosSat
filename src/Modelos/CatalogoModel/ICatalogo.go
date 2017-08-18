//#########################< MODELOS >#############################
//#########################< ICatalogo.go >###############################
//################< ESTRUCTURA Y FUNCIONES >#######################

package CatalogoModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//ICatalogo interface con los métodos de la clase
type ICatalogo interface {
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

	SiguienteClaveDisponible()
}

//################################################<<METODOS DE GESTION >>################################################################

//##################################<< INSERTAR >>###################################

//InsertaMgo es un método que crea un registro en Mongo
func (p CatalogoMgo) InsertaMgo() bool {
	result := false
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	if err != nil {
		fmt.Println(err)
	}

	err = Catalogos.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p CatalogoMgo) InsertaElastic() bool {
	CatalogoE := p.PreparaDatosELastic()
	insert := MoConexion.InsertaElastic(MoVar.TipoCatalogo, p.ID.Hex(), CatalogoE)
	if !insert {
		fmt.Println("Error al insertar Catalogo en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p CatalogoMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Catalogos.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p CatalogoMgo) ActualizaElastic() error {
	CatalogoE := p.PreparaDatosELastic()
	err := MoConexion.ActualizaElastic(MoVar.TipoCatalogo, p.ID.Hex(), CatalogoE)
	return err
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p CatalogoMgo) ReemplazaMgo() bool {
	result := false
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	err = Catalogos.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Catalogo en elastic
func (p CatalogoMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCatalogo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Catalogo en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoCatalogo, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Catalogo en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p CatalogoMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	if err != nil {
		fmt.Println("Error al obtener la conexion en mongo.", err)
	}
	n, e := Catalogos.Find(bson.M{field: valor}).Count()
	if e != nil {
		fmt.Println("Error al consultar la existencia en Mongo", e)
	}
	if n > 0 {
		result = true
	}
	s.Close()
	return result
}

//ConsultaExistenciaByIDMgo es un método que encuentra un registro en Mongo buscándolo por ID
func (p CatalogoMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Catalogos.Find(bson.M{"_id": p.ID}).Count()
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
func (p CatalogoMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoCatalogo, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p CatalogoMgo) EliminaByIDMgo() bool {
	result := false
	s, Catalogos, err := MoConexion.GetColectionMgo(MoVar.ColeccionCatalogo)
	if err != nil {
		fmt.Println(err)
	}
	e := Catalogos.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p CatalogoMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCatalogo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Catalogo en Elastic")
		return false
	}
	return true
}

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p CatalogoMgo) PreparaDatosELastic() CatalogoElastic {
	var CatalogoE CatalogoElastic
	CatalogoE.Clave = p.Clave
	CatalogoE.Nombre = p.Nombre
	CatalogoE.Descripcion = p.Descripcion
	CatalogoE.Valores = p.Valores
	CatalogoE.Estatus = RegresaNombreSubCatalogo(p.Estatus)
	CatalogoE.FechaHora = p.FechaHora
	return CatalogoE
}
