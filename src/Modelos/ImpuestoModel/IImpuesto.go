package ImpuestoModel

import (
	"fmt"

	"../../Modelos/CatalogoModel"
	"../../Modelos/UnidadModel"
	"../../Modulos/Conexiones"
	"../../Modulos/Variables"

	"gopkg.in/mgo.v2/bson"
)

//IImpuesto interface con los métodos de la clase
type IImpuesto interface {
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
func (p ImpuestoMgo) InsertaMgo() bool {
	result := false
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}

	err = Impuestos.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p ImpuestoMgo) InsertaElastic() bool {
	var ImpuestoE ImpuestoElastic
	var dataImpuestosElastic DataImpuestoElastic

	ImpuestoE.Descripcion = p.Descripcion
	ImpuestoE.Clasificacion = CatalogoModel.RegresaNombreSubCatalogo(p.Clasificacion)

	value := p.Datos

	dataImpuestosElastic.Nombre = value.Nombre
	dataImpuestosElastic.Min = value.Min
	dataImpuestosElastic.Max = value.Max
	dataImpuestosElastic.TipoFactor = CatalogoModel.RegresaNombreSubCatalogo(value.TipoFactor)
	dataImpuestosElastic.Unidad = UnidadModel.RegresaNombreUnidad(value.Unidad)
	dataImpuestosElastic.FechaHora = value.FechaHora

	ImpuestoE.Datos = dataImpuestosElastic
	ImpuestoE.Estatus = CatalogoModel.RegresaNombreSubCatalogo(p.Estatus)
	ImpuestoE.FechaHora = p.FechaHora

	insert := MoConexion.InsertaElastic(MoVar.TipoImpuesto, p.ID.Hex(), ImpuestoE)
	if !insert {
		fmt.Println("Error al insertar Impuesto en Elastic")
		return false
	}

	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p ImpuestoMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Impuestos.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p ImpuestoMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoImpuesto, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Impuesto en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Impuesto en Elastic, se perdió Referencia.")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p ImpuestoMgo) ReemplazaMgo() bool {
	result := false
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	err = Impuestos.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Impuesto en elastic
func (p ImpuestoMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoImpuesto, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Impuesto en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoImpuesto, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Impuesto en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p ImpuestoMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Impuestos.Find(bson.M{field: valor}).Count()
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
func (p ImpuestoMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Impuestos.Find(bson.M{"_id": p.ID}).Count()
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
func (p ImpuestoMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoImpuesto, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p ImpuestoMgo) EliminaByIDMgo() bool {
	result := false
	s, Impuestos, err := MoConexion.GetColectionMgo(MoVar.ColeccionImpuesto)
	if err != nil {
		fmt.Println(err)
	}
	e := Impuestos.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p ImpuestoMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoImpuesto, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Impuesto en Elastic")
		return false
	}
	return true
}
