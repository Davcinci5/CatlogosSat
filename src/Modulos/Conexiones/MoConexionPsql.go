package MoConexion

import (
	"database/sql"
	"fmt"

	"../../Modulos/Variables"
)

//ParametrosConexionPostgres estructura que contiene los datos de conexion para postgres
type ParametrosConexionPostgres struct {
	Servidor   string
	Usuario    string
	Pass       string
	NombreBase string
}

//DataP es una estructura que contiene los datos de configuración en el archivo cfg
var DataP = MoVar.CargaSeccionCFG(MoVar.SecPsql)

//ConexionPsql abre una conexión a PostgreSql
func ConexionPsql() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DataP.Servidor, DataP.Usuario, DataP.Pass, DataP.NombreBase)
	db, err := sql.Open("postgres", dbinfo)
	return db, err
}

//ConexioServidorAlmacen establece una conecion al servidor del almacen
func ConexioServidorAlmacen(conex ParametrosConexionPostgres) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", conex.Servidor, conex.Usuario, conex.Pass, conex.NombreBase)
	db, err := sql.Open("postgres", dbinfo)
	return db, err
}

//ConexioServidorAlmacenPing  sirve para testear una conexion
func ConexioServidorAlmacenPing(paramConex ParametrosConexionPostgres) (bool, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", paramConex.Servidor, paramConex.Usuario, paramConex.Pass, paramConex.NombreBase)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()

	if err != nil {
		return false, err
	}
	return true, err
}

//ConexionEspecificaPsql abre una conexión especifica a PostgreSql
func ConexionEspecificaPsql(conex ParametrosConexionPostgres) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", conex.Servidor, conex.Usuario, conex.Pass, conex.NombreBase)
	db, err := sql.Open("postgres", dbinfo)
	return db, err
}

//IniciaSesionEspecificaPsql regresa una base de datos con su sesion especificada en el parametro para commit y rollback
func IniciaSesionEspecificaPsql(conex ParametrosConexionPostgres) (*sql.DB, *sql.Tx, error) {
	Psql, err := ConexionEspecificaPsql(conex)
	if err != nil {
		return nil, nil, err
	}

	tx, err := Psql.Begin()
	if err != nil {
		return nil, nil, err
	}
	return Psql, tx, nil
}

//IniciaSesionPsql regresa una base de datos con su sesion para commit y rollback
func IniciaSesionPsql() (*sql.DB, *sql.Tx, error) {
	Psql, err := ConexionPsql()
	if err != nil {
		return nil, nil, err
	}

	tx, err := Psql.Begin()
	if err != nil {
		return nil, nil, err
	}
	return Psql, tx, nil
}

/*

//CreaTablaPsql crea una tabla en Postgres pasandole el nombre de la tabla a crear
func CreaTablaPsql(conex *sql.DB, nombreTabla string) bool {
	query := `CREATE TABLE IF NOT EXISTS ` + nombreTabla + ` (
  				IdProducto varchar(25) NOT NULL,
 				Precio numeric NOT NULL DEFAULT '0.0',
  				IdImpuesto varchar(25),
  				PRIMARY KEY (IdProducto, Precio, IdImpuesto)
			)`
	con, errsql := conex.Query(query)
	if errsql != nil {
		fmt.Println("Ocurrio un error en la base de datos postgres: ", errsql)
		return false
	}
	con.Close()
	return true
}

//ConsultaDatosDeProductoActivo consulta datos de un producto en un almacen dado
func ConsultaDatosDeProductoActivo(almacen, idProducto string) (bool, float64, float64, float64, error) {
	BasePsql, err := ConexionPsql()
	if err != nil {
		return false, 0, 0, 0, err
	}
	defer BasePsql.Close()

	Query := `SELECT "Existencia", "Costo", "Precio" FROM "Inventario_` + almacen + `" WHERE "IdProducto"='` + idProducto + `' and "Estatus" = 'ACTIVO'`
	Elemento, err := BasePsql.Query(Query)
	if err != nil {
		return false, 0, 0, 0, err
	}

	var Cantidad, Costo, Precio float64
	for Elemento.Next() {
		err := Elemento.Scan(&Cantidad, &Costo, &Precio)
		if err != nil {
			return false, 0, 0, 0, err
		}
	}

	return true, Cantidad, Costo, Precio, nil
}

//ConsultaExistenciaProductoActivo consulta si existe un producto y devuelve la cantidad de existencias y el precio
func ConsultaExistenciaProductoActivo(conex ParametrosConexionPostgres, almacen, idProducto string) (bool, error) {
	BasePsql, err := ConexionEspecificaPsql(conex)
	if err != nil {
		return false, err
	}
	defer BasePsql.Close()
	Query := `SELECT COUNT(*) FROM "Inventario_` + almacen + `" WHERE "IdProducto"='` + idProducto + `' and "Estatus" = 'ACTIVO'`
	Elemento, err := BasePsql.Query(Query)

	if err != nil {
		return false, err
	}

	var Numero int
	for Elemento.Next() {
		err := Elemento.Scan(&Numero)
		if err != nil {
			return false, err
		}
	}

	if Numero > 0 {
		return true, nil
	}

	return false, nil

}

//ConsultaProductoActivo consulta si existe un producto y devuelve la cantidad de existencias y el precio
func ConsultaProductoActivo(conex ParametrosConexionPostgres, idAlmacen, idProducto string) (float64, float64, float64, bool, error) {
	var encontrado = false
	var existencia float64
	var costo float64
	var precio float64

	BasePsql, err := ConexionEspecificaPsql(conex)
	if err != nil {
		return existencia, costo, precio, encontrado, err
	}
	defer BasePsql.Close()
	Query := `SELECT "Existencia", "Costo", "Precio" FROM "Inventario_` + idAlmacen + `" WHERE "IdProducto"='` + idProducto + `' and "Estatus" = 'ACTIVO'`
	Elemento, err := BasePsql.Query(Query)

	if err != nil {
		return existencia, costo, precio, encontrado, err
	}
	for Elemento.Next() {
		encontrado = true
		err := Elemento.Scan(&existencia, &costo, &precio)
		if err != nil {
			return existencia, costo, precio, encontrado, err
		}
	}
	return existencia, costo, precio, encontrado, err
}

//ConsultaPrecioExistenciaYActualizaProductoEnAlmacen consulta si existe un producto y devuelve la cantidad de existencias y el precio
func ConsultaPrecioExistenciaYActualizaProductoEnAlmacen(conex ParametrosConexionPostgres, Operacion, Movimiento, almacen, idProducto string, ValorPrevio, ValorNuevo float64) (bool, float64, float64, error) {

	esta, err := ConsultaExistenciaProductoActivo(conex, almacen, idProducto)
	if err != nil {
		return false, 0, 0, err
	}
	if esta {
		BasePsql, SesionPsql, err := IniciaSesionEspecificaPsql(conex)
		if err != nil {
			return false, 0, 0, err
		}

		BasePsql.Exec("set transaction isolation level serializable")

		Query := fmt.Sprintf(`SELECT "Existencia", "Precio", "Costo" FROM public."Inventario_%v" WHERE "IdProducto" = '%v' and "Estatus" = 'ACTIVO' FOR UPDATE`, almacen, idProducto)
		stmt, err := SesionPsql.Prepare(Query)
		if err != nil {
			return false, 0, 0, err
		}

		resultSet, err := stmt.Query()
		if err != nil {
			return false, 0, 0, err
		}

		var cantidad float64
		var precio float64
		var costo float64
		// var impuesto float64
		// var descuento float64

		for resultSet.Next() {
			resultSet.Scan(&cantidad, &precio, &costo)
		}

		Resto := cantidad + ValorPrevio - ValorNuevo

		if ValorNuevo > cantidad+ValorPrevio {
			SesionPsql.Rollback()
			resultSet.Close()
			stmt.Close()
			BasePsql.Close()
			return false, cantidad, precio, nil
		}

		if Resto >= 0 {
			Query = fmt.Sprintf(`UPDATE  public."Inventario_%v"  SET  "Existencia" = %v WHERE "IdProducto" ='%v'`, almacen, Resto, idProducto)
			_, err := SesionPsql.Exec(Query)
			if err != nil {
				return false, cantidad, precio, err
			}

			Query = fmt.Sprintf(`SELECT COUNT(*) FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, Operacion, Movimiento, almacen, idProducto)
			stmt, err := SesionPsql.Prepare(Query)
			if err != nil {
				return false, 0, 0, err
			}

			resultSet, err := stmt.Query()
			if err != nil {
				return false, 0, 0, err
			}

			var Numero float64
			for resultSet.Next() {
				resultSet.Scan(&Numero)
			}

			if Numero > 0 {
				Query = fmt.Sprintf(`UPDATE  public."VentaTemporal"  SET "Cantidad"= "Cantidad" + %v, "Costo"=%v, "Precio"=%v, "Existencia"=%v, "Impuesto"=0, "Descuento"=0  WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, ValorNuevo-ValorPrevio, costo, precio, Resto, Operacion, Movimiento, almacen, idProducto)
			} else {
				Query = fmt.Sprintf(`INSERT INTO  public."VentaTemporal"  VALUES('%v', '%v', '%v', '%v', %v, %v, %v, 0, 0, %v)`, Operacion, Movimiento, idProducto, almacen, ValorNuevo-ValorPrevio, costo, precio, Resto)
			}
			_, err = SesionPsql.Exec(Query)
			if err != nil {
				SesionPsql.Rollback()
				resultSet.Close()
				stmt.Close()
				BasePsql.Close()
				return false, cantidad, precio, err
			}

			SesionPsql.Commit()
			resultSet.Close()
			stmt.Close()
			BasePsql.Close()

			return true, Resto, precio, nil
		}

	}
	return false, 0, 0, nil
}

//ConsultaPrecioExistenciaYActualizaProductoEnAlmacenModal consulta si existe un producto y devuelve la cantidad de existencias y el precio
func ConsultaPrecioExistenciaYActualizaProductoEnAlmacenModal(conex ParametrosConexionPostgres, Operacion string, Movimiento, almacen, idProducto, ValorNuevo []string) (bool, error) {
	var bandera bool
	var ValoresPrevios []float64

	BasePsql, SesionPsql, err := IniciaSesionEspecificaPsql(conex)
	if err != nil {
		fmt.Println(err)
		bandera = true
	}

	for i, v := range idProducto {
		Query := fmt.Sprintf(`SELECT "Cantidad" FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Movimiento" = '%v' AND "Almacen" = '%v' AND "Producto" = '%v'`, Operacion, Movimiento[i], almacen[i], v)

		con, err := BasePsql.Query(Query)
		if err != nil {
			fmt.Println(err)
		}

		var Prev float64
		for con.Next() {
			con.Scan(&Prev)
		}

		ValoresPrevios = append(ValoresPrevios, Prev)
	}

	for k, v := range idProducto {

		esta, err := ConsultaExistenciaProductoActivo(conex, almacen[k], v)
		if err != nil {
			bandera = true
			fmt.Println(err)
		}
		if esta {
			BasePsql.Exec("set transaction isolation level serializable")

			Query := fmt.Sprintf(`SELECT "Existencia", "Precio", "Costo" FROM public."Inventario_%v" WHERE "IdProducto" = '%v' and "Estatus" = 'ACTIVO' FOR UPDATE`, almacen[k], v)
			stmt, err := SesionPsql.Prepare(Query)
			if err != nil {
				fmt.Println(err)
				bandera = true
			}

			resultSet, err := stmt.Query()
			if err != nil {
				fmt.Println(err)
				bandera = true
			}

			var cantidad float64
			var precio float64
			var costo float64

			for resultSet.Next() {
				resultSet.Scan(&cantidad, &precio, &costo)
			}
			VN, _ := strconv.ParseFloat(ValorNuevo[k], 64)

			Resto := cantidad + ValoresPrevios[k] - VN

			if Resto >= 0 {
				Query = fmt.Sprintf(`UPDATE  public."Inventario_%v"  SET  "Existencia" = %v WHERE "IdProducto" ='%v'`, almacen[k], Resto, v)
				_, err := SesionPsql.Exec(Query)
				if err != nil {
					fmt.Println(err)
					bandera = true
				}

				Query = fmt.Sprintf(`SELECT COUNT(*) FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, Operacion, Movimiento[k], almacen[k], v)
				stmt, err := SesionPsql.Prepare(Query)
				if err != nil {
					fmt.Println(err)
					bandera = true
				}

				resultSet, err := stmt.Query()
				if err != nil {
					fmt.Println(err)
					bandera = true
				}

				var Numero float64
				for resultSet.Next() {
					resultSet.Scan(&Numero)
				}

				if Numero > 0 {
					Query = fmt.Sprintf(`UPDATE  public."VentaTemporal"  SET "Cantidad"= "Cantidad" + %v, "Costo"=%v, "Precio"=%v, "Existencia"=%v, "Impuesto"=0, "Descuento"=0  WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, VN-ValoresPrevios[k], costo, precio, Resto, Operacion, Movimiento[k], almacen[k], v)
				} else {
					Query = fmt.Sprintf(`INSERT INTO  public."VentaTemporal"  VALUES('%v', '%v', '%v', '%v', %v, %v, %v, 0, 0, %v)`, Operacion, Movimiento[k], v, almacen[k], VN-ValoresPrevios[k], costo, precio, Resto)
				}
				_, err = SesionPsql.Exec(Query)
				if err != nil {
					fmt.Println(err)
					SesionPsql.Rollback()
					bandera = true
				}

			} else {
				fmt.Println(err)
				SesionPsql.Rollback()
				resultSet.Close()
				stmt.Close()
				bandera = false
			}
		}
	}

	SesionPsql.Commit()
	BasePsql.Close()
	if bandera {
		return false, err
	}

	return true, nil

}

//ActualizaVentaTemporal actualiza el almacen temporal de Ventas (EL CARRITO)
func ActualizaVentaTemporal(Operacion, Movimiento, Producto, Almacen string, Cantidad, Costo, Precio, Impuesto, Descuento float64) (bool, error) {
	BasePsql, err := ConexionPsql()
	if err != nil {
		return false, err
	}
	defer BasePsql.Close()

	Query := fmt.Sprintf(`SELECT COUNT(*) FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' FOR UPDATE`, Operacion, Movimiento, Almacen, Producto)
	Elemento, err := BasePsql.Query(Query)
	if err != nil {
		return false, err
	}

	var Numero int
	for Elemento.Next() {
		err := Elemento.Scan(&Numero)
		if err != nil {
			return false, err
		}
	}

	if Numero > 0 {
		Query = fmt.Sprintf(`UPDATE  public."VentaTemporal"  SET "Cantidad"= "Cantidad" + %v, "Costo"=%v, "Precio"=%v, "Impuesto"=0, "Descuento"=0 WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v'`, Cantidad, Costo, Precio, Operacion, Movimiento, Almacen, Producto)
	} else {
		Query = fmt.Sprintf(`INSERT INTO  public."VentaTemporal"  VALUES('%v', '%v', '%v', '%v', %v, %v, %v, 0, 0)`, Operacion, Movimiento, Producto, Almacen, Cantidad, Costo, Precio)
	}

	row, err := BasePsql.Exec(Query)
	if err != nil {
		return false, err
	}

	afectadas, err := row.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(afectadas, " Registros afectados.")

	return true, nil

}

//EliminaProductoCarritoYActualizaInventarioAlmacen elimina el producto del almacen temporal y afecta el inventario
func EliminaProductoCarritoYActualizaInventarioAlmacen(Operacion, idProducto string) (bool, error) {

	BasePsql, SesionPsql, err := IniciaSesionPsql()
	if err != nil {
		return false, err
	}

	BasePsql.Exec("set transaction isolation level serializable")
	Query := fmt.Sprintf(`SELECT "Operacion", "Movimiento","Producto","Almacen","Cantidad" FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Producto" = '%v'  FOR UPDATE`, Operacion, idProducto)
	fmt.Println(Query)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		return false, err
	}

	resultSet, err := stmt.Query()
	if err != nil {
		return false, err
	}

	var operacion string
	var movimiento string
	var idproducto string
	var almacen string
	var cantidad float64

	for resultSet.Next() {
		resultSet.Scan(&operacion, &movimiento, &idproducto, &almacen, &cantidad)
		fmt.Println(operacion, movimiento, idproducto, almacen, cantidad)
	}

	Query = fmt.Sprintf(`UPDATE  public."Inventario_%v"  SET  "Existencia" = "Existencia" + %v WHERE "IdProducto" ='%v'`, almacen, cantidad, idproducto)
	fmt.Println(Query)
	_, err = SesionPsql.Exec(Query)
	if err != nil {
		SesionPsql.Rollback()
		return false, err
	}

	Query = fmt.Sprintf(`DELETE FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, operacion, movimiento, almacen, idproducto)
	fmt.Println(Query)
	stmt, err = SesionPsql.Prepare(Query)
	if err != nil {
		SesionPsql.Rollback()
		return false, err
	}

	resultSet, err = stmt.Query()
	if err != nil {
		SesionPsql.Rollback()
		return false, err
	}

	SesionPsql.Commit()
	resultSet.Close()
	stmt.Close()
	BasePsql.Close()

	return true, err

}
*/
