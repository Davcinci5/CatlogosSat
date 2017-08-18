package CargaCombos

import (
	"fmt"
	"strconv"

	"../../Modelos/AlmacenModel"
	"../../Modelos/CatalogoModel"
	"../../Modelos/ConexionModel"
	"../../Modelos/EquipoCajaModel"
	"../../Modelos/ImpuestoModel"
	"../../Modelos/UnidadModel"
	"gopkg.in/mgo.v2/bson"
)

//CargaComboSexo carga el combo de sexos
func CargaComboSexo(SexoSelect string) string {
	var Sexos = []string{"Masculino", "Femenino"}
	templ := ``
	for _, Sexo := range Sexos {
		if Sexo == SexoSelect {
			templ += `<option value="` + Sexo + `" selected>` + Sexo + `</option>`
		} else {
			templ += `<option value="` + Sexo + `">` + Sexo + `</option>`
		}
	}
	return templ
}

//CargaComboMostrarEnIndex carga las opciones de mostrar en el index
func CargaComboMostrarEnIndex(Muestra int) string {
	var Cantidades = []int{5, 10, 15, 20}
	templ := ``

	for _, v := range Cantidades {
		if Muestra == v {
			templ += `<option value="` + strconv.Itoa(v) + `" selected>` + strconv.Itoa(v) + `</option>`
		} else {
			templ += `<option value="` + strconv.Itoa(v) + `">` + strconv.Itoa(v) + `</option>`
		}
	}
	return templ
}

//CargaComboCatalogo recibe la clave del catálogo, el identificador opcional
//y regresa el template del combo del catálogo con el identificador seleccionado si así se desea
func CargaComboCatalogo(Clave int, ID string) string {

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	if Clave == 0 {
		Catalogos := CatalogoModel.GetAll()
		for _, v := range Catalogos {
			if ID == v.ID.Hex() {
				templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Nombre + `</option>`
			} else {
				templ += `<option value="` + v.ID.Hex() + `">` + v.Nombre + `</option>`
			}
		}
	} else {
		Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))
		for _, v := range Catalogo.Valores {
			if ID == v.ID.Hex() {
				templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Valor + `</option>`
			} else {
				templ += `<option value="` + v.ID.Hex() + `">` + v.Valor + `</option>`
			}
		}
	}
	return templ
}

//CargaComboCatalogo2 recibe la clave del catálogo, el identificador opcional
//y regresa el template del combo del catálogo con el identificador seleccionado si así se desea
func CargaComboCatalogo2(Clave int, nombre string) string {

	templ := ``

	if nombre != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	if Clave == 0 {
		Catalogos := CatalogoModel.GetAll()
		for _, v := range Catalogos {
			if nombre == v.Nombre {
				templ += `<option value="` + v.Nombre + `" selected>` + v.Nombre + `</option>`
			} else {
				templ += `<option value="` + v.Nombre + `">` + v.Nombre + `</option>`
			}
		}
	} else {
		Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))
		for _, v := range Catalogo.Valores {
			if nombre == v.Valor {
				templ += `<option value="` + v.Valor + `" selected>` + v.Valor + `</option>`
			} else {
				templ += `<option value="` + v.Valor + `">` + v.Valor + `</option>`
			}
		}
	}
	return templ
}

//CargaComboImpuesto estrae un select con los valores de los impuestos cargados
//Recibe un nombre de impuesto (tipo ObjectId), y un id del subimpuesto seleccionado
//Regresa un string tipo select
func CargaComboImpuesto(IDGrupoImpuesto bson.ObjectId, IDImpuestoDefault bson.ObjectId) string {

	templ := ``

	if IDImpuestoDefault != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	if IDGrupoImpuesto != "" {

		Impuesto := ImpuestoModel.GetManyEspecificByFields("Nombre", IDGrupoImpuesto)
		for _, v := range Impuesto {
			if IDImpuestoDefault == v.ID {
				templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Datos.Nombre + `</option>`
			} else {
				templ += `<option value="` + v.ID.Hex() + `">` + v.Datos.Nombre + `</option>`
			}
		}
	}
	return templ
}

//CargaComboImpuestos construye el combo de impuestos por Grupo de Impuestos
func CargaComboImpuestos(Tipo, impuesto string) string {

	templ := ``

	if impuesto != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	if Tipo != "" {

		Impuesto := ImpuestoModel.GetManyEspecificByFields("TipoImpuesto", bson.ObjectIdHex(Tipo))
		for _, v := range Impuesto {
			if impuesto == v.ID.Hex() {
				templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Datos.Nombre + `</option>`
			} else {
				templ += `<option value="` + v.ID.Hex() + `">` + v.Datos.Nombre + `</option>`
			}
		}
	}
	return templ
}

//CargaComboTipoDeImpuestos carga todos los tipos de impuestos
func CargaComboTipoDeImpuestos(IDImpuestoDefault string) string {
	TiposDeImpuesto := ImpuestoModel.GetAllTiposDeImpuestos()

	templ := ``

	if IDImpuestoDefault != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, v := range TiposDeImpuesto {
		if IDImpuestoDefault == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Descripcion + `</option>`
		} else {
			templ += `<option value="` + v.ID.Hex() + `">` + v.Descripcion + `</option>`
		}
	}

	return templ
}

//CargaComboTipoDeFactor carga todos los tipos de impuestos
func CargaComboTipoDeFactor(IDImpuestoDefault string) string {

	TiposDeFactores := ImpuestoModel.GetAllTiposDeFactores()

	templ := ``

	if IDImpuestoDefault != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, v := range TiposDeFactores {
		if IDImpuestoDefault == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Descripcion + `</option>`
		} else {
			templ += `<option value="` + v.ID.Hex() + `">` + v.Descripcion + `</option>`
		}
	}

	return templ
}

//CargaComboAlmacenes Obtiene los nombres de los almacenes por clasificacion
//recibe el identificador de la clasificación (cliente, propio, proveedor, transporte, etc)
func CargaComboAlmacenes(idClasificacion bson.ObjectId) string {
	templ := `<option value="">--SELECCIONE--</option>`

	almacen, _ := AlmacenModel.GetEspecificsByTagAndTestConexion("Clasificacion", idClasificacion)
	numAlm := len(almacen)
	if numAlm > 0 {
		for _, v := range almacen {
			if v.Conexion != "" {
				templ += `<option value="` + v.ID.Hex() + `">` + v.Nombre + `</option>`
			}
		}
	} else {
		fmt.Println("No se encontraron almacenes")
	}

	return templ
}

//CargaComboUnidades recibe un Id en caso de solicitar combo seleccionado
func CargaComboUnidades(ID string) string {
	unidades := UnidadModel.GetAll()
	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	opc := ``

	for _, v := range unidades {
		templ += `<optgroup label="` + v.Magnitud + `">`
		opc = ``
		for _, val := range v.Datos {
			if ID == val.ID.Hex() {
				opc += `<option value="` + val.ID.Hex() + `" selected>` + val.Abreviatura + `</option>`
			} else {
				opc += `<option value="` + val.ID.Hex() + `">` + val.Abreviatura + `</option>`
			}
		}
		templ = templ + opc + `</optgroup>`
	}
	return templ
}

//CargaComboCatalogoMulti Carga el combo para un multi select
func CargaComboCatalogoMulti(Clave int, ID string) string {

	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))

	templ := ``

	for _, v := range Catalogo.Valores {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>` + v.Valor + `</option>`
		} else {
			templ += `<option value="` + v.ID.Hex() + `">` + v.Valor + `</option>`
		}
	}
	return templ
}

//CargaEstatusActivoEnAlta Extrae el Object Id del status 'ACTIVO' de un catalogo
func CargaEstatusActivoEnAlta(Clave int) bson.ObjectId {

	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))

	var objectid bson.ObjectId

	for _, v := range Catalogo.Valores {

		if v.Valor == "ACTIVO" {
			objectid = v.ID
		}
	}

	return objectid

}

//CargaCatalogoByID Extrae el Valor deseado de un catalgo desde el numero de catalogo y el id del valor deseado
func CargaCatalogoByID(Clave int, id bson.ObjectId) string {
	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))
	var cadena string
	for _, v := range Catalogo.Valores {
		if v.ID == id {
			cadena = v.Valor
		}
	}
	return cadena
}

//ArrayStringToObjectID Convierte un Arreglo de String a uno de Objects Ids
func ArrayStringToObjectID(ArrayStr []string) []bson.ObjectId {
	var ArrayID []bson.ObjectId
	for _, d := range ArrayStr {
		if bson.IsObjectIdHex(d) {
			ArrayID = append(ArrayID, bson.ObjectIdHex(d))
		}
	}
	return ArrayID
}

//CargaComboCajas carga el listado de cajas en un select para seleccionar en la apertura de cajas. By @melchormendoza
func CargaComboCajas() string {
	templ := `<option value="">--SELECCIONE--</option>`

	almacen := EquipoCajaModel.GetAll()
	for _, v := range almacen {
		templ += `<option value="` + v.ID.Hex() + `">` + v.Nombre + `</option>`
	}
	return templ
}

//CargaComboConexiones carga un listado de conexiones acia postgres
func CargaComboConexiones(idConexion string) string {
	conexiones := ConexionModel.GetAll()
	templ := ``

	if idConexion != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, val := range conexiones {
		if idConexion == val.ID.Hex() {
			templ += `<option value="` + val.ID.Hex() + `" selected>` + val.Nombre + `[` + val.Servidor + `]` + `</option>`
		} else {
			templ += `<option value="` + val.ID.Hex() + `">` + val.Nombre + `[` + val.Servidor + `]` + `</option>`
		}
	}

	return templ
}

// CargaComboTiposPersonas Carga el combo de un catalogo, seleccionando las opciones mediante un arreglo de ID
func CargaComboTiposPersonas(Clave int, ID []string) string {
	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))
	templ := ``
	for _, vv := range Catalogo.Valores {
		existe := false
		for _, v := range ID {
			if v == vv.ID.Hex() {
				existe = true
			}
		}
		if existe {
			templ += `<option value="` + vv.ID.Hex() + `" selected>` + vv.Valor + `</option>`
		} else {
			templ += `<option value="` + vv.ID.Hex() + `">` + vv.Valor + `</option>`
		}
	}
	return templ
}

//CargaIDTipoPersonaPorDefecto regresa el ID de la persona por defecto que en este caso es Cliente
func CargaIDTipoPersonaPorDefecto(Clave int) bson.ObjectId {
	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))
	var objectid bson.ObjectId
	for _, vv := range Catalogo.Valores {
		if vv.Valor == "USUARIO" {
			objectid = vv.ID
		}
	}
	return objectid
}

//ExisteEnCatalogo Extrae el ID del valor de un catalogo, primer parametro nuemro de catalogo, y segundo la URI
func ExisteEnCatalogo(Clave int, valor string) string {
	Catalogo := CatalogoModel.GetEspecificByFields("Clave", int64(Clave))
	var idusr bson.ObjectId
	var existe = false
	for _, v := range Catalogo.Valores {
		if v.Valor == valor {
			existe = true
			idusr = v.ID
		}
	}
	if !existe {
		fmt.Println("No existe esta Uri")
		URIS := CatalogoModel.RegresaValoresCatalogosClave(Clave)
		nuevaURI := CatalogoModel.ValoresMgo{}
		nuevaURI.ID = bson.NewObjectId()
		nuevaURI.Valor = valor

		URIS.Valores = append(URIS.Valores, nuevaURI)
		fmt.Println(URIS)
		if URIS.ReemplazaMgo() {
			existe = true
			idusr = nuevaURI.ID
		}
		fmt.Println(existe)
	}
	return idusr.Hex()
}
