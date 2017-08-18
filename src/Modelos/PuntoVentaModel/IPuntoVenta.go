package PuntoVentaModel

import (
	"fmt"
	"time"

	"../../Modulos/Conexiones"
	"../../Modulos/ConsultasSql"
	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
)

//IPuntoVenta interface con los métodos de la clase
type IPuntoVenta interface {
	InsertaMgo() bool
	InsertaElastic() bool

	ActualizaMgo(campos []string, valores []interface{}) bool
	ActualizaElastic(campos []string, valores []interface{}) bool //Reemplaza No Actualiza

	ReemplazaMgo() bool
	ReemplazaElastic() bool

	ConsultaExistenciaByFieldMgo(field string, valor string) bool

	ConsultaExistenciaByIDMgo() bool
	ConsultaExistenciaByIDElastic() bool

	EliminaByIDMgo() bool
	EliminaByIDElastic() bool
}

//################################################<<METODOS DE GESTION >>################################################################

//##################################<< INSERTAR >>###################################

//InsertaMgo es un método que crea un registro en Mongo
func (p PuntoVentaMgo) InsertaMgo() bool {
	result := false
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}

	err = PuntoVentas.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p PuntoVentaMgo) InsertaElastic() bool {
	insert := MoConexion.InsertaElastic(MoVar.TipoPuntoVenta, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al insertar PuntoVenta en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p PuntoVentaMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = PuntoVentas.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p PuntoVentaMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPuntoVenta, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar PuntoVenta en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoPuntoVenta, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar PuntoVenta en Elastic")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p PuntoVentaMgo) ReemplazaMgo() bool {
	result := false
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	err = PuntoVentas.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un PuntoVenta en elastic
func (p PuntoVentaMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPuntoVenta, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar PuntoVenta en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoPuntoVenta, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar PuntoVenta en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p PuntoVentaMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}
	n, e := PuntoVentas.Find(bson.M{field: valor}).Count()
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
func (p PuntoVentaMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}
	n, e := PuntoVentas.Find(bson.M{"_id": p.ID}).Count()
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
func (p PuntoVentaMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoPuntoVenta, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p PuntoVentaMgo) EliminaByIDMgo() bool {
	result := false
	s, PuntoVentas, err := MoConexion.GetColectionMgo(MoVar.ColeccionPuntoVenta)
	if err != nil {
		fmt.Println(err)
	}
	e := PuntoVentas.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p PuntoVentaMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoPuntoVenta, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar PuntoVenta en Elastic")
		return false
	}
	return true
}

//InsertaKardexAlmacen inserta las entradas o salidas en el kardex correspondiente al almacen
//Recibe como parametro el nombre de la tabla, una conexion, y los datos
func (d DatosVentaTemporal) InsertaKardexAlmacen() bool {
	paramConex := ConsultasSql.InsertarDatosConexion(bson.ObjectIdHex(d.Almacen))
	BasePsql, err := MoConexion.ConexionEspecificaPsql(paramConex)
	if err != nil {
		fmt.Println("Error al iniciar la sesion", err)
		return false
	}
	defer BasePsql.Close()

	Query := fmt.Sprintf(`INSERT INTO public."Kardex_%v" VALUES('%v','%v','%v',%v,%v,%v,%v,%v,'%v',%v,'%v')`, d.Almacen, d.Operacion, d.Movimiento, d.Producto, d.Cantidad, d.Costo, d.Precio, d.Impuesto, d.Descuento, "58efbf8bd2b2131778e9c929", d.Existencia, time.Now().Format("2006-01-02 15:04:05"))
	row, err := BasePsql.Exec(Query)
	if err != nil {
		fmt.Println("Error al ejecutar la consulta", err)
		fmt.Println(Query)
		return false
	}

	_, err = row.RowsAffected()
	if err != nil {
		fmt.Println("Error al consultar los afectados", err)
		return false
	}
	return true
}

//InsertaImpuestoAlmacen inserta las entradas o salidas en el kardex correspondiente al almacen
//Recibe como parametro el nombre de la tabla, una conexion, y los datos
func (d DatosVentaTemporal) InsertaImpuestoAlmacen(Impuestos []string) bool {
	BasePsql, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer BasePsql.Close()

	Query := fmt.Sprintf(`INSERT INTO public."Impuesto_%v" VALUES('%v','%v','%v','%v')`, d.Almacen, d.Movimiento, d.Producto, d.IDimpuesto, d.Impuesto)
	row, err := BasePsql.Exec(Query)
	if err != nil {
		fmt.Println(err)
		return false
	}

	afectadas, err := row.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(afectadas, " Registros afectados.")

	return true
}

//InsertaDescuentoAlmacen inserta las entradas o salidas en el kardex correspondiente al almacen
//Recibe como parametro el nombre de la tabla, una conexion, y los datos
func (d DatosVentaTemporal) InsertaDescuentoAlmacen(Impuestos []string) bool {
	BasePsql, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer BasePsql.Close()

	Query := fmt.Sprintf(`INSERT INTO public."Impuesto_%v" VALUES('%v','%v','%v','%v')`, d.Almacen, d.Movimiento, d.Producto, d.IDdescuento, d.Descuento)
	row, err := BasePsql.Exec(Query)
	if err != nil {
		fmt.Println(err)
		return false
	}

	afectadas, err := row.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(afectadas, " Registros afectados.")

	return true
}

//ActualizaKardexAlmacen inserta las entradas o salidas en el kardex correspondiente al almacen
//Recibe como parametro el nombre de la tabla, una conexion, y los datos
func (d DatosVentaTemporal) ActualizaKardexAlmacen() bool {
	BasePsql, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer BasePsql.Close()

	Query := fmt.Sprintf(`DELETE FROM public."Kardex_%v" WHERE "IdMovimiento" = '%v'`, d.Almacen, d.Movimiento)
	row, err := BasePsql.Exec(Query)
	if err != nil {
		fmt.Println(err)
		return false
	}

	Query = fmt.Sprintf(`INSERT INTO public."Kardex_%v" VALUES('%v','%v','%v','%v','%v','%v','%v','%v')`, d.Almacen, d.Movimiento, d.Producto, d.Cantidad, d.Costo, d.Precio, "58efbf8bd2b2131778e9c929", d.Existencia, time.Now().Format("2006-01-02 15:04:05"))
	row, err = BasePsql.Exec(Query)
	if err != nil {
		fmt.Println(err)
		return false
	}

	afectadas, err := row.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(afectadas, " Registros afectados.")

	return true
}
