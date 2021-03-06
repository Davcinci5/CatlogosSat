package PromocionModel

import (
	"fmt"

	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IPromocion interface con los métodos de la clase
type IPromocion interface {
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
func (p PromocionMgo) InsertaMgo() bool {
	result := false
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}

	err = Promocions.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p PromocionMgo) InsertaElastic() bool {
	var PromocionE PromocionElastic

	PromocionE.Nombre = p.Nombre
	PromocionE.Descripcion = p.Descripcion
	// PromocionE.PorcentajeDesc = p.PorcentajeDesc
	// PromocionE.PrecioOferta = p.PrecioOferta
	// PromocionE.OfertaMonto = p.OfertaMonto
	// PromocionE.OfertaPiezaPieza = p.OfertaPiezaPieza
	// PromocionE.OfertaPiezaPorcentaje = p.OfertaPiezaPorcentaje
	// PromocionE.Estatus = p.Estatus
	PromocionE.FechaInicio = p.FechaInicio
	PromocionE.FechaFin = p.FechaFin
	insert := MoConexion.InsertaElastic(MoVar.TipoPromocion, p.ID.Hex(), PromocionE)
	if !insert {
		fmt.Println("Error al insertar Promocion en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p PromocionMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Promocions.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p PromocionMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPromocion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Promocion en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Promocion en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p PromocionMgo) ReemplazaMgo() bool {
	result := false
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	err = Promocions.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Promocion en elastic
func (p PromocionMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPromocion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Promocion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoPromocion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Promocion en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p PromocionMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Promocions.Find(bson.M{field: valor}).Count()
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
func (p PromocionMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Promocions.Find(bson.M{"_id": p.ID}).Count()
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
func (p PromocionMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoPromocion, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p PromocionMgo) EliminaByIDMgo() bool {
	result := false
	s, Promocions, err := MoConexion.GetColectionMgo(MoVar.ColeccionPromocion)
	if err != nil {
		fmt.Println(err)
	}
	e := Promocions.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p PromocionMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPromocion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Promocion en Elastic")
		return false
	}
	return true
}
