package CotizacionModel

import (
	"fmt"

	"../../Modelos/ClienteModel"
	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//ICotizacion interface con los métodos de la clase
type ICotizacion interface {
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
func (p CotizacionMgo) InsertaMgo() bool {
	result := false
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}

	err = Cotizacions.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p CotizacionMgo) InsertaElastic() bool {
	var CotizacionE CotizacionElastic

	//CotizacionE.Operacion = p.Operacion
	//CotizacionE.Usuario = p.Usuario
	//CotizacionE.Equipo = p.Equipo
	if bson.IsObjectIdHex(p.Cliente.Hex()) {
		esta, cli := ClienteModel.NombreCliente(p.Cliente)
		if !esta {
			fmt.Println("No se encontro el cliente ")
			return false
		}
		CotizacionE.Cliente = cli
	}
	//CotizacionE.Grupo = p.Grupo
	CotizacionE.Nombre = p.Nombre
	CotizacionE.Telefono = p.Telefono
	CotizacionE.Correo = p.Correo
	//CotizacionE.FormaDePago = p.FormaDePago
	//CotizacionE.FormaDeEnvío = p.FormaDeEnvío
	// CotizacionE.Buscar = p.Buscar
	// CotizacionE.Lista = p.Lista
	// CotizacionE.Carrito = p.Carrito
	// CotizacionE.Resumen = p.Resumen
	//CotizacionE.Estatus = p.Estatus
	CotizacionE.FechaInicio = p.FechaInicio
	CotizacionE.FechaFin = p.FechaFin
	insert := MoConexion.InsertaElastic(MoVar.TipoCotizacion, p.ID.Hex(), CotizacionE)
	if !insert {
		fmt.Println("Error al insertar Cotizacion en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p CotizacionMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Cotizacions.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p CotizacionMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCotizacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Cotizacion en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Cotizacion en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p CotizacionMgo) ReemplazaMgo() bool {
	result := false
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	err = Cotizacions.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Cotizacion en elastic
func (p CotizacionMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCotizacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Cotizacion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoCotizacion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Cotizacion en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p CotizacionMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Cotizacions.Find(bson.M{field: valor}).Count()
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
func (p CotizacionMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Cotizacions.Find(bson.M{"_id": p.ID}).Count()
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
func (p CotizacionMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoCotizacion, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p CotizacionMgo) EliminaByIDMgo() bool {
	result := false
	s, Cotizacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionCotizacion)
	if err != nil {
		fmt.Println(err)
	}
	e := Cotizacions.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p CotizacionMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCotizacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Cotizacion en Elastic")
		return false
	}
	return true
}
