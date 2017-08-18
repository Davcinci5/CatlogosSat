package ClienteModel

import (
	"fmt"

	"../../Modelos/CatalogoModel"
	"../../Modelos/PersonaModel"
	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//ICliente interface con los métodos de la clase
type ICliente interface {
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
func (p ClienteMgo) InsertaMgo() bool {
	result := false
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}

	err = Clientes.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p ClienteMgo) InsertaElastic() bool {
	ClienteE := p.PreparaDatosELastic()
	insert := MoConexion.InsertaElastic(MoVar.TipoCliente, p.ID.Hex(), ClienteE)
	if !insert {
		fmt.Println("Error al insertar Cliente en Elastic")
		return false
	}
	return true
}

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p ClienteMgo) PreparaDatosELastic() ClienteElastic {
	var ClienteE ClienteElastic
	ClienteE.IDPersona = PersonaModel.GetOne(p.IDPersona).Nombre
	ClienteE.TipoCliente = CatalogoModel.RegresaNombreSubCatalogo(p.TipoCliente)
	ClienteE.RFC = p.RFC
	var direccionesElastic []DireccionElastic
	for _, v := range p.Direcciones {
		var direccionElastic DireccionElastic
		direccionElastic.Calle = v.Calle
		direccionElastic.NumInterior = v.NumInterior
		direccionElastic.NumExterior = v.NumExterior
		direccionElastic.Colonia = CatalogoModel.GetNameColonia(v.Colonia)
		direccionElastic.Municipio = CatalogoModel.GetNameMunicipio(v.Municipio)
		direccionElastic.Estado = CatalogoModel.GetNameEstado(v.Estado)
		direccionElastic.Pais = v.Pais
		direccionElastic.CP = v.CP
		direccionElastic.Estatus = CatalogoModel.RegresaNombreSubCatalogo(p.Estatus)
		direccionesElastic = append(direccionesElastic, direccionElastic)
	}
	ClienteE.Direcciones = direccionesElastic
	ClienteE.MediosDeContacto.Correos.Principal = p.MediosDeContacto.Correos.Principal
	ClienteE.MediosDeContacto.Correos.Correos = p.MediosDeContacto.Correos.Correos
	ClienteE.MediosDeContacto.Telefonos.Principal = p.MediosDeContacto.Telefonos.Principal
	ClienteE.MediosDeContacto.Telefonos.Telefonos = p.MediosDeContacto.Telefonos.Telefonos
	ClienteE.MediosDeContacto.Otros = p.MediosDeContacto.Otros

	var personasContactoElastic []PersonaContactoElastic
	for _, v := range p.PersonasContacto {
		var personaContactoElastic PersonaContactoElastic
		var almacenesElastic []AlmacenElastic
		for _, v1 := range v.Almacenes {
			var almacenElastic AlmacenElastic
			almacenElastic.IDContacto = v1.IDContacto
			almacenElastic.IDDireccion = v1.IDDireccion
			almacenElastic.IDAlmacen = v1.IDAlmacen
			almacenesElastic = append(almacenesElastic, almacenElastic)
		}
		personaContactoElastic.Almacenes = almacenesElastic

		var direccionesElastic []DireccionElastic
		for _, v2 := range p.Direcciones {
			var direccionElastic DireccionElastic
			direccionElastic.Calle = v2.Calle
			direccionElastic.NumInterior = v2.NumInterior
			direccionElastic.NumExterior = v2.NumExterior
			direccionElastic.Colonia = CatalogoModel.GetNameColonia(v2.Colonia)
			direccionElastic.Municipio = CatalogoModel.GetNameMunicipio(v2.Municipio)
			direccionElastic.Estado = CatalogoModel.GetNameEstado(v2.Estado)
			direccionElastic.Pais = v2.Pais
			direccionElastic.CP = v2.CP
			direccionElastic.Estatus = CatalogoModel.RegresaNombreSubCatalogo(bson.ObjectIdHex(v2.Estatus))
			direccionesElastic = append(direccionesElastic, direccionElastic)
		}

		personaContactoElastic.Direcciones = direccionesElastic
		personaContactoElastic.Estatus = CatalogoModel.RegresaNombreSubCatalogo(bson.ObjectIdHex(v.Estatus))

		personaContactoElastic.Direcciones = direccionesElastic
		personaContactoElastic.MediosDeContacto.Correos.Principal = v.MediosDeContacto.Correos.Principal
		personaContactoElastic.MediosDeContacto.Correos.Correos = v.MediosDeContacto.Correos.Correos
		personaContactoElastic.MediosDeContacto.Telefonos.Principal = v.MediosDeContacto.Telefonos.Principal
		personaContactoElastic.MediosDeContacto.Telefonos.Telefonos = v.MediosDeContacto.Telefonos.Telefonos
		personaContactoElastic.MediosDeContacto.Otros = v.MediosDeContacto.Otros

		personaContactoElastic.Nombre = v.Nombre
		personasContactoElastic = append(personasContactoElastic, personaContactoElastic)
	}
	ClienteE.PersonasContacto = personasContactoElastic
	for _, v3 := range p.Almacenes {
		ClienteE.Almacenes = append(ClienteE.Almacenes, v3.Hex())
	}

	//ClienteE.Notificaciones = p.Notificaciones
	ClienteE.Estatus = p.Estatus.Hex()
	ClienteE.FechaHora = p.FechaHora
	return ClienteE
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p ClienteMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Clientes.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p ClienteMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCliente, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Cliente en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Cliente en Elastic, se perdió Referencia.")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p ClienteMgo) ReemplazaMgo() bool {
	result := false
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	err = Clientes.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Cliente en elastic
func (p ClienteMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCliente, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Cliente en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoCliente, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Cliente en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p ClienteMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Clientes.Find(bson.M{field: valor}).Count()
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
func (p ClienteMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Clientes.Find(bson.M{"_id": p.ID}).Count()
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
func (p ClienteMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoCliente, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p ClienteMgo) EliminaByIDMgo() bool {
	result := false
	s, Clientes, err := MoConexion.GetColectionMgo(MoVar.ColeccionCliente)
	if err != nil {
		fmt.Println(err)
	}
	e := Clientes.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p ClienteMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoCliente, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Cliente en Elastic")
		return false
	}
	return true
}
