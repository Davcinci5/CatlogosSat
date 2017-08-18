package ProductoModel

import (
	"fmt"

	"../../Modelos/CatalogoModel"
	"../../Modelos/UnidadModel"
	"../../Modulos/Conexiones"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IProducto interface con los métodos de la clase
type IProducto interface {
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
func (p ProductoMgo) InsertaMgo() bool {
	result := false
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}

	err = Productos.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p ProductoMgo) InsertaElastic() bool {

	ProductoE := p.PreparaDatosELastic()

	insert := MoConexion.InsertaElastic(MoVar.TipoProducto, p.ID.Hex(), ProductoE)
	if !insert {
		fmt.Println("Error al insertar Producto en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p ProductoMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Productos.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaAlmacenes agrega referencia de almacenes al producto siempre y cuando no exista en el arreglo de almacenes
func (p ProductoMgo) ActualizaAlmacenes(almacen bson.ObjectId) bool {
	result := false
	var existe bool
	producto := GetOne(p.ID)
	for _, value := range producto.Almacenes {
		if value == almacen {
			existe = true
		}
	}
	if !existe {
		fmt.Println("No existe")
		PushToArray := bson.M{"$push": bson.M{"Almacenes": almacen}}
		s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
		err = Productos.Update(bson.M{"_id": p.ID}, PushToArray)
		if err != nil {
			fmt.Println(err)
		} else {
			result = true
		}
		s.Close()
	}
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p ProductoMgo) ActualizaElastic() error {
	ProductoE := p.PreparaDatosELastic()
	err := MoConexion.ActualizaElastic(MoVar.TipoProducto, p.ID.Hex(), ProductoE)
	return err
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p ProductoMgo) ReemplazaMgo() bool {
	result := false
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	err = Productos.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Producto en elastic
func (p ProductoMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoProducto, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Producto en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoProducto, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Producto en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p ProductoMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Productos.Find(bson.M{field: valor}).Count()
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
func (p ProductoMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Productos.Find(bson.M{"_id": p.ID}).Count()
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
func (p ProductoMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoProducto, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p ProductoMgo) EliminaByIDMgo() bool {
	result := false
	s, Productos, err := MoConexion.GetColectionMgo(MoVar.ColeccionProducto)
	if err != nil {
		fmt.Println(err)
	}
	e := Productos.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p ProductoMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoProducto, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Producto en Elastic")
		return false
	}
	return true
}

//###################################<< FUNCIONES EXTRAS >> ############################################

//PreparaDatosELastic  obtiene los datos por defecto de mongo y los convierte en string de tal forma que
//se inserteadecuadamente en elastic
func (p ProductoMgo) PreparaDatosELastic() ProductoElastic {
	var ProductoE ProductoElastic
	ProductoE.Nombre = p.Nombre
	ProductoE.Codigos = p.Codigos
	ProductoE.Tipo = CatalogoModel.RegresaNombreSubCatalogo(p.Tipo)
	ProductoE.Unidad = UnidadModel.RegresaNombreUnidad(p.Unidad)
	ProductoE.VentaFraccion = p.VentaFraccion
	ProductoE.Etiquetas = p.Etiquetas
	ProductoE.Estatus = CatalogoModel.RegresaNombreSubCatalogo(p.Estatus)
	ProductoE.FechaHora = p.FechaHora
	return ProductoE
}
