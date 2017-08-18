package MediosPagoModel

import (
	"fmt"

	"../../Modelos/CatalogoModel"
	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IMediosPago interface con los métodos de la clase
type IMediosPago interface {
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
func (p MediosPagoMgo) InsertaMgo() bool {
	result := false
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}

	err = MediosPagos.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p MediosPagoMgo) InsertaElastic() bool {
	MediosPagoE := p.PreparaDatosELastic()
	insert := MoConexion.InsertaElastic(MoVar.TipoMediosPago, p.ID.Hex(), MediosPagoE)
	if !insert {
		fmt.Println("Error al insertar MediosPago en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p MediosPagoMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = MediosPagos.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p MediosPagoMgo) ActualizaElastic() error {
	MediosPagoE := p.PreparaDatosELastic()
	err := MoConexion.ActualizaElastic(MoVar.TipoMediosPago, p.ID.Hex(), MediosPagoE)
	return err
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p MediosPagoMgo) ReemplazaMgo() bool {
	result := false
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	err = MediosPagos.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un MediosPago en elastic
func (p MediosPagoMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoMediosPago, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar MediosPago en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoMediosPago, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar MediosPago en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p MediosPagoMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}
	n, e := MediosPagos.Find(bson.M{field: valor}).Count()
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
func (p MediosPagoMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}
	n, e := MediosPagos.Find(bson.M{"_id": p.ID}).Count()
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
func (p MediosPagoMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoMediosPago, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p MediosPagoMgo) EliminaByIDMgo() bool {
	result := false
	s, MediosPagos, err := MoConexion.GetColectionMgo(MoVar.ColeccionMediosPago)
	if err != nil {
		fmt.Println(err)
	}
	e := MediosPagos.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p MediosPagoMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoMediosPago, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar MediosPago en Elastic")
		return false
	}
	return true
}

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p MediosPagoMgo) PreparaDatosELastic() MediosPagoElastic {
	var MediosPagoE MediosPagoElastic
	MediosPagoE.Nombre = CatalogoModel.RegresaNombreSubCatalogo(p.Nombre)
	MediosPagoE.Descripcion = p.Descripcion
	MediosPagoE.CodigoSat = p.CodigoSat
	MediosPagoE.Tipo = CatalogoModel.RegresaNombreSubCatalogo(p.Tipo)
	MediosPagoE.Comision = p.Comision
	MediosPagoE.Cambio = p.Cambio
	MediosPagoE.Estatus = CatalogoModel.RegresaNombreSubCatalogo(p.Estatus)
	MediosPagoE.FechaHora = p.FechaHora
	return MediosPagoE
}
