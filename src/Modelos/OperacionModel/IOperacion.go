package OperacionModel

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modulos/Conexiones"
	"../../Modulos/ConsultasSql"
	"../../Modulos/Variables"
	_ "github.com/bmizerany/pq"
	"gopkg.in/mgo.v2/bson"
)

//IOperacion interface con los métodos de la clase
type IOperacion interface {
	InsertaMgo() bool
	InsertaElastic() bool

	InsertaKardexAlmacen(conex *sql.DB, nombreTabla string) bool
	InsertaKardexYActualizaInventario(idAlmacen bson.ObjectId) (string, bool, error)
	InsertaProductoInventario(idAlmacen bson.ObjectId) bool
	InsertaKardexInsertaInventario(idAlmacen bson.ObjectId) (string, bool, error)

	ActualizaMgo(campos []string, valores []interface{}) bool
	ActualizaElastic(campos []string, valores []interface{}) bool //Reemplaza No Actualiza

	ReemplazaMgo() bool
	ReemplazaElastic() bool

	ConsultaExistenciaByFieldMgo(field string, valor string)

	ConsultaExistenciaByIDMgo() bool
	ConsultaExistenciaByIDElastic() bool

	ConsultaInventarioPostgres(nombreTabla string) []InventarioPostgres
	ConsultaKardexPostgres() bool

	EliminaByIDMgo() bool
	EliminaByIDElastic() bool
}

//################################################<<METODOS DE GESTION >>################################################################

//##################################<< INSERTAR >>###################################

//InsertaMgo es un método que crea un registro en Mongo
func (p OperacionMgo) InsertaMgo() bool {
	result := false
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}

	err = Operacions.Insert(p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}

	s.Close()
	return result
}

//InsertaElastic es un método que crea un registro en Mongo
func (p OperacionMgo) InsertaElastic() bool {
	insert := MoConexion.InsertaElastic(MoVar.TipoOperacion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al insertar Operacion en Elastic")
		return false
	}
	return true
}

//##########################<< UPDATE >>############################################

//ActualizaMgo es un método que encuentra y Actualiza un registro en Mongo
//IMPORTANTE --> Debe coincidir el número y orden de campos con el de valores
func (p OperacionMgo) ActualizaMgo(campos []string, valores []interface{}) bool {
	result := false
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	var Abson bson.M
	Abson = make(map[string]interface{})
	for k, v := range campos {
		Abson[v] = valores[k]
	}
	change := bson.M{"$set": Abson}
	err = Operacions.Update(bson.M{"_id": p.ID}, change)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ActualizaElastic es un método que encuentra y Actualiza un registro en Mongo
func (p OperacionMgo) ActualizaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoOperacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Operacion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoOperacion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Operacion en Elastic")
		return false
	}
	return true
}

//##########################<< REEMPLAZA >>############################################

//ReemplazaMgo es un método que encuentra y Actualiza un registro en Mongo
func (p OperacionMgo) ReemplazaMgo() bool {
	result := false
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	err = Operacions.Update(bson.M{"_id": p.ID}, p)
	if err != nil {
		fmt.Println(err)
	} else {
		result = true
	}
	s.Close()
	return result
}

//ReemplazaElastic es un método que encuentra y reemplaza un Operacion en elastic
func (p OperacionMgo) ReemplazaElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoOperacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Operacion en Elastic")
		return false
	}
	insert := MoConexion.InsertaElastic(MoVar.TipoOperacion, p.ID.Hex(), p)
	if !insert {
		fmt.Println("Error al actualizar Operacion en Elastic")
		return false
	}
	return true
}

//###########################<< CONSULTA EXISTENCIAS >>###################################

//ConsultaExistenciaByFieldMgo es un método que verifica si un registro existe en Mongo indicando un campo y un valor string
func (p OperacionMgo) ConsultaExistenciaByFieldMgo(field string, valor string) bool {
	result := false
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Operacions.Find(bson.M{field: valor}).Count()
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
func (p OperacionMgo) ConsultaExistenciaByIDMgo() bool {
	result := false
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}
	n, e := Operacions.Find(bson.M{"_id": p.ID}).Count()
	if e != nil {
		fmt.Println(e)
	}
	if n > 0 {
		result = true
	}
	s.Close()
	return result
}

//ConsultaInventarioPostgres extrae las cantidades existentes del inventario
func (p InventarioPostgres) ConsultaInventarioPostgres(idInventario bson.ObjectId) (EInventarioPostgres, error) {
	Conexion := AlmacenModel.ObtenerParametrosConexion(idInventario)
	var parametros MoConexion.ParametrosConexionPostgres
	parametros.Servidor = Conexion.Servidor
	parametros.Usuario = Conexion.UsuarioBD
	parametros.Pass = Conexion.PassBD
	parametros.NombreBase = Conexion.NombreBD
	var inv EInventarioPostgres

	conex, err := MoConexion.ConexioServidorAlmacen(parametros)
	if err != nil {
		fmt.Println("Errores con: ", idInventario, err)
		return inv, err
	}
	consulta := `SELECT * FROM  "Inventario_` + idInventario.Hex() + `" WHERE "IdProducto" = '` + p.IDProducto + `'`
	rows, errsql := conex.Query(consulta)
	if errsql != nil {
		return inv, err
	}
	var encontrado bool
	for rows.Next() {
		encontrado = true
		err := rows.Scan(&inv.IDProducto, &inv.Existencia, &inv.Costo, &inv.Precio, &inv.Estatus)
		inv.Encontrado = true
		if err != nil {
			fmt.Println(err)
			inv.Encontrado = false
			return inv, err
		}
	}
	if !encontrado {
		inv.IDProducto = "NINGUNO"
		inv.Existencia = 0
		inv.Estatus = "NINGUNO"
		inv.Encontrado = true
		inv.Costo = 0
		inv.Precio = 0
		return inv, err
	}
	rows.Close()
	conex.Close()
	return inv, err
}

//ConsultaKardexPostgres Extrae la informacion del kardex d postgres pasandole como parametro el identificador de movimiento
func (p KardexPostgres) ConsultaKardexPostgres() bool {
	return false
}

//ConsultaExistenciaByIDElastic es un método que encuentra un registro en Mongo buscándolo por ID
func (p OperacionMgo) ConsultaExistenciaByIDElastic() bool {
	result := MoConexion.ConsultaElastic(MoVar.TipoOperacion, p.ID.Hex())
	return result
}

//##################################<< ELIMINACIONES >>#################################################

//EliminaByIDMgo es un método que elimina un registro en Mongo
func (p OperacionMgo) EliminaByIDMgo() bool {
	result := false
	s, Operacions, err := MoConexion.GetColectionMgo(MoVar.ColeccionOperacion)
	if err != nil {
		fmt.Println(err)
	}
	e := Operacions.RemoveId(bson.M{"_id": p.ID})
	if e != nil {
		result = true
	} else {
		fmt.Println(e)
	}
	s.Close()
	return result
}

//EliminaByIDElastic es un método que elimina un registro en Mongo
func (p OperacionMgo) EliminaByIDElastic() bool {
	delete := MoConexion.DeleteElastic(MoVar.TipoOperacion, p.ID.Hex())
	if !delete {
		fmt.Println("Error al actualizar Operacion en Elastic")
		return false
	}
	return true
}

//##############################<< INSERTA EN POSTGRES >>##################################################

//InsertaKardexAlmacen inserta las entradas o salidas en el kardex correspondiente al almacen
//Recibe como parametro el nombre de la tabla, una conexion, y los datos
func (p KardexPostgres) InsertaKardexAlmacen(nombreTabla string) bool {
	conex, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println("Error en InsertarKardexAlmacen:", err)
		return false
	}
	defer conex.Close()
	//Convierte todas las variables en cadenas, para formar la consulta requerida
	cantidad := strconv.FormatFloat(p.Cantidad, 'f', -1, 64)
	costo := strconv.FormatFloat(p.Costo, 'f', -1, 64)
	precio := strconv.FormatFloat(p.Precio, 'f', -1, 64)
	existencia := strconv.FormatFloat(p.Existencia, 'f', -1, 64)
	fechaPostgres := p.FechaHora
	fechaPostgresFormat := fechaPostgres.Format("2006-01-02 15:04:05")
	//Realiza la consulta para inserta posteriormente en postgres
	query := `INSERT INTO "` + nombreTabla + `" VALUES('` + p.IDMovimiento + `', '` + p.IDProducto + `',` + cantidad + `,` + costo + `,` + precio + `,'` + p.TipoOperacion + `',` + existencia + `,'` + fechaPostgresFormat + `')`
	con, errsql := conex.Query(query)
	if errsql != nil {
		fmt.Println("Ocurrio un error en la base de datos postgres: ", errsql)
		return false
	}
	//Cierra las conexiones y en su caso devuelve verdadero (asumiendo que la consulta se realizó con éxito)
	con.Close()
	return true
}

//InsertaKardexYActualizaInventario inserta los elementos en el kardex correspondiente
//Asi mismo consulta la existencia del producto en el inventario
//Actualiza la existencia en el kardex
//Y por ultimo actualiza la existencia en el inventario
func (p KardexPostgres) InsertaKardexYActualizaInventario(idAlmacen bson.ObjectId) (string, bool, error) {
	datosConexion := AlmacenModel.ObtenerParametrosConexion(idAlmacen)

	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	//Se asegura de que el producto esté en el inventario correspondiente (para ello se ocupa la variable booleana)
	existencia, _, _, encontrado, err := ConsultasSql.ConsultaProductoActivo(idAlmacen.Hex(), p.IDProducto)
	if err != nil {
		return "No encontrado", false, err
	}
	var SesionPsql *sql.Tx
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	if err != nil {
		return "Error en la sesion", false, err
	}
	BasePsql.Exec("set transaction isolation level serializable")
	nuevaExistencia := existencia + p.Cantidad
	Query := fmt.Sprintf(`INSERT INTO public."Kardex_%v" VALUES('%v','%v','%v',%v,%v,%v,%v,%v,'%v',%v,'%v')`, idAlmacen.Hex(), p.IDOperacion, p.IDMovimiento, p.IDProducto, p.Cantidad, p.Costo, p.Precio, p.ImpuestoTotal, p.DescuentoTotal, p.TipoOperacion, nuevaExistencia, time.Now().Format(time.RFC3339))
	_, err = SesionPsql.Exec(Query)
	if err != nil {
		return "Error al insertar en el kardex", false, err
	}

	if encontrado {
		Query = fmt.Sprintf(`UPDATE  public."Inventario_%v"  SET  "Existencia" = %v, "Costo" = %v, "Precio" = %v  WHERE "IdProducto" ='%v'`, idAlmacen.Hex(), nuevaExistencia, p.Costo, p.Precio, p.IDProducto)
		_, err = SesionPsql.Exec(Query)
		if err != nil {
			fmt.Println("Ha ocurrido un error en la actualizacion", Query)
			SesionPsql.Rollback()
			BasePsql.Close()
		}
		SesionPsql.Commit()
		return "Realizado", true, err
	}
	query := fmt.Sprintf(`INSERT INTO public."Inventario_%v" VALUES('%v',%v,%v,%v,'%v')`, idAlmacen.Hex(), p.IDProducto, p.Cantidad, p.Precio, p.Costo, "ACTIVO")
	_, errsql := SesionPsql.Exec(query)
	if errsql != nil {
		SesionPsql.Rollback()
		BasePsql.Close()
		fmt.Println("Error al insertar el producto")
		fmt.Println(query)
		return "Error al relizar la operacion", false, err
	}
	SesionPsql.Commit()
	BasePsql.Close()
	return "Operacion realizada con éxito", true, err
}

//InsertaKardexInsertaInventario inserta un elemento en el kardex y en el almacen
//generalmente solo se ocupa cuando se da de alta un producto
func (p KardexPostgres) InsertaKardexInsertaInventario(idAlmacen bson.ObjectId, inv InventarioPostgres) (string, bool, error) {
	datosConexion := AlmacenModel.ObtenerParametrosConexion(idAlmacen)

	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	//BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	if err != nil {
		return "Error en la sesion", false, err
	}
	BasePsql.Exec("set transaction isolation level serializable")
	Query := fmt.Sprintf(`INSERT INTO public."Kardex_%v" VALUES('%v','%v','%v',%v,%v,%v,%v,%v,'%v',%v,'%v')`, idAlmacen.Hex(), p.IDOperacion, p.IDMovimiento, p.IDProducto, p.Cantidad, p.Costo, p.Precio, p.ImpuestoTotal, p.DescuentoTotal, p.TipoOperacion, inv.Existencia, time.Now().Format(time.RFC3339))
	fmt.Println("This One:", Query)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		return "Error en la consulta a postgres", false, err
	}
	_, err = stmt.Query()
	if err != nil {
		return "Error al buscar en inventario", false, err
	}
	/*
		if ValorNuevo > cantidad+ValorPrevio {
			SesionPsql.Rollback()
			resultSet.Close()
			stmt.Close()
			BasePsql.Close()
			return false, cantidad, precio, nil
		}
	*/
	query := fmt.Sprintf(`INSERT INTO public."Inventario_%v" VALUES('%v',%v,%v,%v,'%v')`, idAlmacen.Hex(), inv.IDProducto, inv.Existencia, inv.Costo, inv.Precio, inv.Estatus)
	fmt.Println(query)
	_, errsql := BasePsql.Query(query)
	if errsql != nil {
		fmt.Println("Ocurrio un error en la base de datos postgres: ", errsql)
		return "Error en la insercion de inventario", false, errsql
	}
	SesionPsql.Commit()
	stmt.Close()
	BasePsql.Close()
	return "Realizado", true, nil
}

//InsertaProductoInventario Inserta un producto en el Inventario de postgres
func (p InventarioPostgres) InsertaProductoInventario(idAlmacen bson.ObjectId) bool {
	conex, err := MoConexion.ConexionPsql()
	if err != nil {
		fmt.Println("Error en InsertaProductoInventario:", err)
		return false
	}
	defer conex.Close()
	//Realiza la consulta para inserta posteriormente en postgres
	query := fmt.Sprintf(`INSERT INTO public."Inventario_%v" VALUES('%v',%v,%v,%v,'%v')`, idAlmacen.Hex(), p.IDProducto, p.Existencia, p.Costo, p.Precio, p.Estatus)
	fmt.Println(query)
	con, errsql := conex.Query(query)
	if errsql != nil {
		fmt.Println("Ocurrio un error en la base de datos postgres: ", errsql)
		return false
	}
	//Cierra las conexiones y en su caso devuelve verdadero (asumiendo que la consulta se realizó con éxito)
	con.Close()
	return true
}

//InsertaImpuestoEnAlmacenComprasPsql inserta un impuesto en un determinado almacen
func (p ImpuestoPostgres) InsertaImpuestoEnAlmacenComprasPsql(idAlmacen bson.ObjectId) (string, bool, error) {

	datosConexion := AlmacenModel.ObtenerParametrosConexion(idAlmacen)

	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	//BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	if err != nil {
		return "Error en la sesión del almacén", false, err
	}
	BasePsql.Exec("set transaction isolation level serializable")

	Query := ``
	if p.TipoDeImpuesto != "" {
		Query = fmt.Sprintf(`INSERT INTO  public."ImpuestoC_%v"  VALUES('%v', '%v', '%v', '%v', '%v', %v)`, idAlmacen.Hex(), p.IDMovimiento, p.IDProducto, p.TipoDeImpuesto, p.Factor, p.Tratamiento, p.Valor)
		_, err = SesionPsql.Exec(Query)
		if err != nil {
			SesionPsql.Rollback()
			BasePsql.Close()
			return "Error al Ejecutar la instrucción en Psql", false, err
		}
	}

	SesionPsql.Commit()
	BasePsql.Close()
	return "Realizado", true, nil
}

//InsertaImpuestoEnAlmacenVentasPsql inserta un impuesto en un determinado almacen
func (p ImpuestoPostgres) InsertaImpuestoEnAlmacenVentasPsql(idAlmacen bson.ObjectId) (string, bool, error) {

	datosConexion := AlmacenModel.ObtenerParametrosConexion(idAlmacen)

	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD

	//BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	if err != nil {
		return "Error en la sesión del almacén", false, err
	}
	BasePsql.Exec("set transaction isolation level serializable")

	if p.TipoDeImpuesto != "" {
		Query := ``
		Query = fmt.Sprintf(`INSERT INTO  public."ImpuestoV_%v"  VALUES('%v', '%v', '%v', '%v', '%v', %v)`, idAlmacen.Hex(), p.IDMovimiento, p.IDProducto, p.TipoDeImpuesto, p.Factor, p.Tratamiento, p.Valor)

		_, err = SesionPsql.Exec(Query)
		if err != nil {
			SesionPsql.Rollback()
			BasePsql.Close()
			return "Error al Ejecutar la instrucción en Psql", false, err
		}
	}

	SesionPsql.Commit()
	BasePsql.Close()
	return "Realizado", true, nil
}
