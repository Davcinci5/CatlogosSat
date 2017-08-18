package GrupoModel

import (
	"fmt"
	"strconv"

	"../../Modelos/AlmacenModel"
	"../../Modelos/CajaModel"
	"../../Modelos/CatalogoModel"
	"../../Modelos/ClienteModel"
	"../../Modelos/DispositivoModel"
	"../../Modelos/PersonaModel"
	"../../Modelos/ProductoModel"
	"../../Modelos/UsuarioModel"

	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"../../Modulos/Variables"

	"gopkg.in/mgo.v2/bson"
)

//IGrupo interface con los métodos de la clase
type IGrupo interface {
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
func (p GrupoMgo) InsertaMgo() bool {
	result := false
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}

	err = Grupos.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p GrupoMgo) InsertaElastic() bool {
	var GrupoE GrupoElastic

	GrupoE.Nombre = p.Nombre
	GrupoE.Descripcion = p.Descripcion
	GrupoE.PermiteVender = strconv.FormatBool(p.PermiteVender)
	GrupoE.Tipo = CatalogoModel.ObtenerValoresCatalogoPorValor(p.Tipo)

	switch GrupoE.Tipo {
	case "CLIENTES":
		Clientes := ClienteModel.GetEspecifics(p.Miembros)
		for _, v := range Clientes {
			if !MoGeneral.EstaVacio(v) {
				Persona := PersonaModel.GetOne(v.IDPersona)
				GrupoE.Miembros = append(GrupoE.Miembros, Persona.Nombre)
			}
		}
		break
	case "PERSONAS":
		Personas := PersonaModel.GetEspecifics(p.Miembros)
		for _, v := range Personas {
			if !MoGeneral.EstaVacio(v) {
				GrupoE.Miembros = append(GrupoE.Miembros, v.Nombre)
			}
		}
		break
	case "USUARIOS":
		Usuarios := UsuarioModel.GetEspecifics(p.Miembros)

		for _, v := range Usuarios {
			if !MoGeneral.EstaVacio(v) {
				GrupoE.Miembros = append(GrupoE.Miembros, v.Usuario)
			}
		}
		break
	case "PRODUCTOS":
		Productos := ProductoModel.GetEspecifics(p.Miembros)
		for _, v := range Productos {
			if !MoGeneral.EstaVacio(v) {
				GrupoE.Miembros = append(GrupoE.Miembros, v.Nombre)
			}
		}
		break
	case "ALMACENES":
		Almacens := AlmacenModel.GetEspecifics(p.Miembros)

		for _, v := range Almacens {
			if !MoGeneral.EstaVacio(v) {
				GrupoE.Miembros = append(GrupoE.Miembros, v.Nombre)
			}
		}
		// case "PROVEEDORES":
		break
	case "CAJAS":
		Cajas := CajaModel.GetEspecifics(p.Miembros)
		for _, v := range Cajas {
			if !MoGeneral.EstaVacio(v) {
				GrupoE.Miembros = append(GrupoE.Miembros, v.ID.Hex())
			}
		}
		break
	case "DISPOSITIVOS":
		Dispositivos := DispositivoModel.GetEspecifics(p.Miembros)

		for _, v := range Dispositivos {
			if !MoGeneral.EstaVacio(v) {
				GrupoE.Miembros = append(GrupoE.Miembros, v.Nombre)
			}
		}
		break
	default:
		fmt.Println("No se encontró el tipo de objeto.")
		return false
		break
	}

	GrupoE.Estatus = CatalogoModel.ObtenerValoresCatalogoPorValor(p.Estatus)
	GrupoE.FechaHora = p.FechaHora
	insert := MoConexion.InsertaElastic(MoVar.TipoGrupo, p.ID.Hex(), GrupoE)
	if !insert {
		fmt.Println("Error al insertar Grupo en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p GrupoMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Grupos.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p GrupoMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Grupo en Elastic")
		return false
	}

	if !p.InsertaElastic() {
		fmt.Println("Error al actualizar Grupo en Elastic, se perdió Referencia.")
		return false
	}

	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p GrupoMgo) ReemplazaMgo() bool {
	result := false
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	err = Grupos.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Grupo en elastic
func (p GrupoMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Grupo en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoGrupo, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Grupo en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p GrupoMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Grupos.Find(bson.M{field: valor}).Count()
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
func (p GrupoMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Grupos.Find(bson.M{"_id": p.ID}).Count()
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
func (p GrupoMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoGrupo, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p GrupoMgo) EliminaByIDMgo() bool {
	result := false
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}
	e := Grupos.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p GrupoMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoGrupo, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Grupo en Elastic")
		return false
	}
	return true
}
