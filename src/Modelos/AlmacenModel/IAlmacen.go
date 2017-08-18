package AlmacenModel

import (
	"fmt"

	"../../Modelos/CatalogoModel"
	"../../Modelos/ConexionModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IAlmacen interface con los métodos de la clase
type IAlmacen interface {
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
func (p AlmacenMgo) InsertaMgo() bool {
	result := false
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}

	err = Almacens.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p AlmacenMgo) InsertaElastic() bool {
	AlmacenE := p.PreparaDatosELastic()
	insert := MoConexion.InsertaElastic(MoVar.TipoAlmacen, p.ID.Hex(), AlmacenE)
	if !insert {
		fmt.Println("Error al insertar Almacen en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p AlmacenMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Almacens.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p AlmacenMgo) ActualizaElastic() error {
	AlmacenE := p.PreparaDatosELastic()
	err := MoConexion.ActualizaElastic(MoVar.ColeccionAlmacen, p.ID.Hex(), AlmacenE)
	return err
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p AlmacenMgo) ReemplazaMgo() bool {
	result := false
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	fmt.Println(p.ID)
	err = Almacens.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println("Error al actualizar ID: ", p.ID.Hex(), ":", err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Almacen en elastic
func (p AlmacenMgo) ReemplazaElastic() bool {
	AlmacenE := p.PreparaDatosELastic()

	delete := MoConexion.DeleteElastic(MoVar.TipoAlmacen, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Almacen en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoAlmacen, p.ID.Hex(), AlmacenE)
	if !insert {
		fmt.Println("Error al actualizar Almacen en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p AlmacenMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Almacens.Find(bson.M{field: valor}).Count()
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
func (p AlmacenMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Almacens.Find(bson.M{"_id": p.ID}).Count()
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
func (p AlmacenMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoAlmacen, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p AlmacenMgo) EliminaByIDMgo() bool {
	var result bool
	s, Almacens, err := MoConexion.GetColectionMgo(MoVar.ColeccionAlmacen)
	if err != nil {
		fmt.Println("Error al consultar la coleccion de almacenes", err)
	}
	e := Almacens.RemoveId(p.ID)
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
func (p AlmacenMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoAlmacen, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Almacen en Elastic")
		return false
	}
	return true
}

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p AlmacenMgo) PreparaDatosELastic() AlmacenElastic {
	var AlmacenE AlmacenElastic
	AlmacenE.Nombre = p.Nombre
	AlmacenE.Tipo = CatalogoModel.RegresaNombreSubCatalogo(p.Tipo)
	AlmacenE.Clasificacion = CatalogoModel.RegresaNombreSubCatalogo(p.Clasificacion)

	//Pendiente incluí nombre de Grupos y Dirección de almacén

	if p.Predecesor != "" {
		AlmacenPredesesor := GetOne(p.Predecesor)
		AlmacenE.Predecesor = AlmacenPredesesor.Nombre
	}
	if len(p.Sucesor) != 0 {
		for _, NombreSucesor := range p.Sucesor {
			AlmacenE.Sucesor = append(AlmacenE.Sucesor, GetOne(NombreSucesor).Nombre)
		}
	}
	if p.Conexion != "" {
		parametrosAlmacenes := ConexionModel.GetOne(p.Conexion)
		AlmacenE.NombreConexion = parametrosAlmacenes.Nombre
	}
	AlmacenE.Estatus = CatalogoModel.RegresaNombreSubCatalogo(p.Estatus)
	AlmacenE.FechaHora = p.FechaHora
	return AlmacenE
}

//GetDescendantsIDs  obtiene el arreglo de ID de los nodos descendientes
func (p AlmacenMgo) GetDescendantsIDs() []bson.ObjectId {
	Descendants := p.Sucesor
	for _, descendant := range Descendants {
		for _, e := range GetOne(descendant).Sucesor {
			Descendants = append(Descendants, e)
		}
	}
	fmt.Println("Descendants: ", Descendants)

	return Descendants
}

//GetAncestorsIDs  obtiene el arreglo de ID de los nodos Ancestros
func (p AlmacenMgo) GetAncestorsIDs() []bson.ObjectId {
	var Ancestors []bson.ObjectId
	if !MoGeneral.EstaVacio(p.Predecesor.Hex()) {
		Ancestors = append(Ancestors, p.Predecesor)
	}
	for _, ancestor := range Ancestors {
		aux := GetOne(ancestor)
		if !MoGeneral.EstaVacio(aux.Predecesor.Hex()) {
			Ancestors = append(Ancestors, aux.Predecesor)
		}
	}
	fmt.Println("Ancestors: ", Ancestors)
	return Ancestors
}

//OthersThanMe  obtiene el arreglo de ID de los nodos diferentes al actual
func (p AlmacenMgo) OthersThanMe() []bson.ObjectId {
	var OthersThanMe []bson.ObjectId               // lista de ids no disponibles para ser descendiente
	for _, ancestor := range p.GetAncestorsIDs() { //excluye a sus ancestros
		OthersThanMe = append(OthersThanMe, ancestor)
	}

	for _, descendant := range p.GetDescendantsIDs() { //excluye a sus descendientes
		OthersThanMe = append(OthersThanMe, descendant)
	}
	return OthersThanMe
}

//RelativesToMe  obtiene el arreglo de ID de los nodos diferentes al actual
func (p AlmacenMgo) RelativesToMe() []bson.ObjectId {
	var Relatives []bson.ObjectId // lista de ids no disponibles para ser descendiente
	Relatives = append(Relatives, p.ID)
	// for _, relative := range p.GetAncestorsIDs() { //Lista ancestros
	// 	Relatives = append(Relatives, relative)
	// }
	for _, relative := range Relatives { //excluye a sus parientes
		for _, item := range GetOne(relative).Sucesor {
			if IndexOfAlmacen(Relatives, item) == -1 {
				Relatives = append(Relatives, relative)
			}
		}
		Relatives = append(Relatives, GetOne(relative).Sucesor[:]...)
	}
	return Relatives
}

//AvailableOthersIDs  obtiene el arreglo de ID de los nodos Disponibles como Descendientes
func (p AlmacenMgo) AvailableOthersIDs() []AlmacenMgo {
	var UnavailableDescendants []bson.ObjectId // lista de ids no disponibles para ser descendiente

	// UnavailableDescendants = p.RelativesToMe() //se excluye a si mismo y sus familiares
	UnavailableDescendants = append(UnavailableDescendants, p.ID) //se excluye a si mismo

	for _, other := range p.OthersThanMe() { //excluye a sus otros que no son el mismo
		UnavailableDescendants = append(UnavailableDescendants, other)
	}
	fmt.Println("Unavailable descendants: ", UnavailableDescendants)
	AvailableDescendants := GetEspecificsByTag("_id", bson.M{"$nin": UnavailableDescendants})

	return AvailableDescendants
}

//AvailableOthersIDsWithSameType  obtiene el arreglo de ID de los nodos Disponibles como Descendientes
func (p AlmacenMgo) AvailableOthersIDsWithSameType() []AlmacenMgo {
	var UnavailableDescendants []bson.ObjectId // lista de ids no disponibles para ser descendiente

	UnavailableDescendants = append(UnavailableDescendants, p.ID) //se excluye a si mismo

	for _, other := range p.OthersThanMe() { //excluye a sus otros que no son el mismo ancestros y descendientes
		UnavailableDescendants = append(UnavailableDescendants, other)
	}
	params := []bson.M{}
	params = append(params, bson.M{"Tipo": p.Tipo})
	params = append(params, bson.M{"$nin": UnavailableDescendants})

	query := bson.M{"$and": params}
	AvailableDescendants := GetAllInQuery(query)

	fmt.Println("others related than me: ", query)
	return AvailableDescendants
}

//ReplaceDescendants  reemplaza Descendientes con un nuevo arreglo
func (p AlmacenMgo) ReplaceDescendants(descendantsArray []bson.ObjectId) {
	//Los hijos olvidan al Predecesor
	var AlmacenAux AlmacenMgo
	for _, descendant := range p.Sucesor {
		Almacen := GetOne(descendant)
		AlmacenAux.ID = Almacen.ID
		AlmacenAux.Nombre = Almacen.Nombre
		AlmacenAux.Tipo = Almacen.Tipo
		AlmacenAux.Clasificacion = Almacen.Clasificacion
		// AlmacenAux.Predecesor = Almacen.Predecesor
		AlmacenAux.Conexion = Almacen.Conexion
		AlmacenAux.Estatus = Almacen.Estatus
		AlmacenAux.FechaHora = Almacen.FechaHora
		AlmacenAux.ReemplazaMgo()
		AlmacenAux.ReemplazaElastic()
	}
	p.Sucesor = descendantsArray

	for _, descendant := range p.Sucesor {
		Almacen := GetOne(descendant)
		AlmacenAux.ID = Almacen.ID
		AlmacenAux.Nombre = Almacen.Nombre
		AlmacenAux.Tipo = Almacen.Tipo
		AlmacenAux.Clasificacion = Almacen.Clasificacion
		AlmacenAux.Predecesor = p.ID
		AlmacenAux.Conexion = Almacen.Conexion
		AlmacenAux.Estatus = Almacen.Estatus
		AlmacenAux.FechaHora = Almacen.FechaHora
		AlmacenAux.ReemplazaMgo()
		AlmacenAux.ReemplazaElastic()
	}

	p.ReemplazaMgo()
	p.ReemplazaElastic()
}

//ForgotParent  Olvida al predescesor
func (p AlmacenMgo) ForgotParent() {
	//Los hijos olvidan al Predecesor
	var AlmacenAux AlmacenMgo
	if !MoGeneral.EstaVacio(p.Predecesor) {
		Parent := GetOne(p.Predecesor)
		myIndexInParent := IndexOfAlmacen(Parent.Sucesor, p.ID)

		if myIndexInParent != -1 {
			Parent.Sucesor = append(Parent.Sucesor[:myIndexInParent], Parent.Sucesor[myIndexInParent+1:]...)
		}

		AlmacenAux.ID = p.ID
		AlmacenAux.Nombre = p.Nombre
		AlmacenAux.Tipo = p.Tipo
		AlmacenAux.Clasificacion = p.Clasificacion
		AlmacenAux.Sucesor = p.Sucesor
		AlmacenAux.Conexion = p.Conexion
		AlmacenAux.Estatus = p.Estatus
		AlmacenAux.FechaHora = p.FechaHora
		AlmacenAux.ReemplazaMgo()
		AlmacenAux.ReemplazaElastic()

		Parent.ReemplazaMgo()
		Parent.ReemplazaElastic()
		AlmacenAux.ReemplazaMgo()
		AlmacenAux.ReemplazaElastic()
	}

}
