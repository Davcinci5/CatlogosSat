package ConsultasSql

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"../../Modelos/AlmacenModel"
	"../../Modulos/Conexiones"
	"../../Modulos/General"
	"gopkg.in/mgo.v2/bson"
)

//InsertarDatosConexion llena los parametros necesarios para conectarse a la base de datos
func InsertarDatosConexion(almacen bson.ObjectId) MoConexion.ParametrosConexionPostgres {
	datosConexion := AlmacenModel.ObtenerParametrosConexion(almacen)
	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = datosConexion.Servidor
	paramConex.Usuario = datosConexion.UsuarioBD
	paramConex.Pass = datosConexion.PassBD
	paramConex.NombreBase = datosConexion.NombreBD
	return paramConex
}

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
	paramConex := InsertarDatosConexion(bson.ObjectIdHex(almacen))
	BasePsql, err := MoConexion.ConexionEspecificaPsql(paramConex)
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
func ConsultaExistenciaProductoActivo(almacen, idProducto string) (bool, error) {
	paramConex := InsertarDatosConexion(bson.ObjectIdHex(almacen))
	BasePsql, err := MoConexion.ConexionEspecificaPsql(paramConex)
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

//ConsultaImpuestoEnAlmacen consulta si existen impuestos en un almacen con un movimiento y producto asociado
func ConsultaImpuestoEnAlmacen(almacen, idProducto, Movimiento string) (bool, error) {
	paramConex := InsertarDatosConexion(bson.ObjectIdHex(almacen))
	BasePsql, err := MoConexion.ConexionEspecificaPsql(paramConex)
	if err != nil {
		return false, err
	}
	defer BasePsql.Close()
	Query := `SELECT COUNT(*) FROM public."ImpuestoV_` + almacen + `" WHERE "IdProducto"='` + idProducto + `' and "IdMovimiento" = '` + Movimiento + `'`
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

//ConsultaValorImpuestoEnAlmacen consulta si existen impuestos en un almacen con un movimiento y producto asociado
func ConsultaValorImpuestoEnAlmacen(almacen, idProducto string) ([]AlmacenModel.ImpuestoAlmacen, error) {
	var Impuestos []AlmacenModel.ImpuestoAlmacen
	var Impuesto AlmacenModel.ImpuestoAlmacen

	var Tratamiento string
	var Valor float64
	var TipoImpuesto string
	var Factor string

	paramConex := InsertarDatosConexion(bson.ObjectIdHex(almacen))
	BasePsql, err := MoConexion.ConexionEspecificaPsql(paramConex)
	if err != nil {
		fmt.Println(err)
		return Impuestos, err
	}
	defer BasePsql.Close()

	Query := fmt.Sprintf(`SELECT "TipoImpuesto", "Factor", "Valor", "Tratamiento" FROM public."ImpuestoV_%v" 
						WHERE "IdMovimiento"= (
							SELECT "IdMovimiento"
							FROM public."Kardex_%v"
							WHERE "IdProducto"='%v'
							AND "TipoOperacion" = '58efbf8bd2b2131778e9c928'
							ORDER BY "FechaHora" DESC LIMIT 1);`, almacen, almacen, idProducto)

	// Query := fmt.Sprintf(`SELECT "TipoImpuesto", "Factor", "Valor" FROM public."ImpuestoV_%v" AS IMP
	// 					JOIN public."Kardex_%v" AS KAR ON IMP."IdMovimiento"= KAR."IdMovimiento"
	// 					WHERE IMP."IdProducto"='%v' AND KAR."FechaHora" = (SELECT MAX("FechaHora") FROM public."Kardex_%v")
	// 					GROUP BY  IMP."TipoImpuesto", IMP."Valor", IMP."Factor", KAR."FechaHora" ORDER BY KAR."FechaHora"   DESC;`, almacen, almacen, idProducto, almacen)

	Elemento, err := BasePsql.Query(Query)
	if err != nil {
		fmt.Println(err)
		return Impuestos, err
	}

	for Elemento.Next() {
		err := Elemento.Scan(&TipoImpuesto, &Factor, &Valor, &Tratamiento)
		if err != nil {
			return Impuestos, err
		}
		Impuesto = AlmacenModel.ImpuestoAlmacen{}
		Impuesto.Tipo = MoGeneral.EliminarEspaciosInicioFinal(TipoImpuesto)
		Impuesto.Factor = MoGeneral.EliminarEspaciosInicioFinal(Factor)
		Impuesto.Valor = Valor
		Impuesto.Tratamiento = MoGeneral.EliminarEspaciosInicioFinal(Tratamiento)
		Impuestos = append(Impuestos, Impuesto)
	}

	return Impuestos, nil
}

//ConsultaProductoActivo consulta si existe un producto y devuelve la cantidad de existencias y el precio
func ConsultaProductoActivo(idAlmacen, idProducto string) (float64, float64, float64, bool, error) {
	var encontrado = false
	var existencia float64
	var costo float64
	var precio float64
	paramConex := InsertarDatosConexion(bson.ObjectIdHex(idAlmacen))
	BasePsql, err := MoConexion.ConexionEspecificaPsql(paramConex)
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
func ConsultaPrecioExistenciaYActualizaProductoEnAlmacen(Operacion, Movimiento, almacen, idProducto string, ValorPrevio, ValorNuevo float64) (bool, float64, float64, float64, float64, error) {

	esta, err := ConsultaExistenciaProductoActivo(almacen, idProducto)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ocurrió un error al consultar EXISTENCIA EN el producto : ", idProducto, "En el almacén: ", almacen)
		return false, 0, 0, 0, 0, err
	}

	Impuestos, err := ConsultaValorImpuestoEnAlmacen(almacen, idProducto)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ocurrió un error al consultar el Impuesto del producto : ", idProducto, " en el Almacen: ", almacen)
	}

	if esta {
		paramConex := InsertarDatosConexion(bson.ObjectIdHex(almacen))
		BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
		if err != nil {
			return false, 0, 0, 0, 0, err
		}

		BasePsql.Exec("set transaction isolation level serializable")

		Query := fmt.Sprintf(`SELECT "Existencia", "Precio", "Costo" FROM public."Inventario_%v" WHERE "IdProducto" = '%v' and "Estatus" = 'ACTIVO' FOR UPDATE`, almacen, idProducto)
		stmt, err := SesionPsql.Prepare(Query)
		if err != nil {
			return false, 0, 0, 0, 0, err
		}

		resultSet, err := stmt.Query()
		if err != nil {
			return false, 0, 0, 0, 0, err
		}

		var cantidad float64
		var precio float64
		var costo float64

		for resultSet.Next() {
			resultSet.Scan(&cantidad, &precio, &costo)
		}

		impuesto := CalculaTotalDeImpuestoPorProducto(Impuestos, precio)
		Resto := cantidad + ValorPrevio - ValorNuevo

		if ValorNuevo > cantidad+ValorPrevio {
			SesionPsql.Rollback()
			resultSet.Close()
			stmt.Close()
			BasePsql.Close()
			return false, cantidad, precio, impuesto, 0, nil
		}

		if Resto >= 0 {
			Query = fmt.Sprintf(`UPDATE  public."Inventario_%v"  SET  "Existencia" = %v WHERE "IdProducto" ='%v'`, almacen, Resto, idProducto)
			_, err := SesionPsql.Exec(Query)
			if err != nil {
				return false, cantidad, precio, impuesto, 0, err
			}

			Query = fmt.Sprintf(`SELECT COUNT(*) FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, Operacion, Movimiento, almacen, idProducto)
			stmt, err := SesionPsql.Prepare(Query)
			if err != nil {
				return false, 0, 0, 0, 0, err
			}

			resultSet, err := stmt.Query()
			if err != nil {
				return false, 0, 0, 0, 0, err
			}

			var Numero float64
			for resultSet.Next() {
				resultSet.Scan(&Numero)
			}

			if Numero > 0 {
				Query = fmt.Sprintf(`UPDATE  public."VentaTemporal"  SET "Cantidad"= "Cantidad" + %v, "Costo"=%v, "Precio"=%v, "Existencia"=%v, "Impuesto"=%v, "Descuento"=0  WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, ValorNuevo-ValorPrevio, costo, precio, Resto, ValorNuevo*impuesto, Operacion, Movimiento, almacen, idProducto)
			} else {
				Query = fmt.Sprintf(`INSERT INTO  public."VentaTemporal"  VALUES('%v', '%v', '%v', '%v', %v, %v, %v, %v, 0, %v, '%v')`, Operacion, Movimiento, idProducto, almacen, ValorNuevo-ValorPrevio, costo, precio, ValorNuevo*impuesto, Resto, time.Now().Format(time.RFC3339Nano))
			}
			_, err = SesionPsql.Exec(Query)
			if err != nil {
				SesionPsql.Rollback()
				resultSet.Close()
				stmt.Close()
				BasePsql.Close()
				return false, cantidad, precio, impuesto, 0, err
			}

			SesionPsql.Commit()
			resultSet.Close()
			stmt.Close()
			BasePsql.Close()

			return true, Resto, precio, impuesto, 0, nil
		}

	}
	return false, 0, 0, 0, 0, nil
}

//AgregaACarritoCotizacion consulta si existe un producto y devuelve la cantidad de existencias y el precio
func AgregaACarritoCotizacion(Operacion, almacen, idProducto string, ValorNuevo float64) (bool, error) {
	fmt.Println("Operacion", Operacion)
	fmt.Println("Almacen", almacen)
	fmt.Println("IdProducto", idProducto)
	fmt.Println("ValorNuevo", ValorNuevo)
	esta, err := ConsultaExistenciaProductoActivo(almacen, idProducto)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ocurrió un error al consultar EXISTENCIA EN el producto : ", idProducto, " en el Almacen: ", almacen)
		return false, err
	}

	Impuestos, err := ConsultaValorImpuestoEnAlmacen(almacen, idProducto)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ocurrió un error al consultar el Impuesto del producto : ", idProducto, " en el Almacen: ", almacen)
		return false, err
	}

	if esta {
		paramConex := InsertarDatosConexion(bson.ObjectIdHex(almacen))
		BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
		if err != nil {
			fmt.Println("Ocurrio un error al IniciaSesionEspecificaPsql", err)
			return false, err
		}

		Query := fmt.Sprintf(`SELECT "Existencia", "Precio", "Costo" FROM public."Inventario_%v" WHERE "IdProducto" = '%v' and "Estatus" = 'ACTIVO'`, almacen, idProducto)
		stmt, err := SesionPsql.Prepare(Query)
		if err != nil {
			return false, err
		}

		resultSet, err := stmt.Query()
		if err != nil {
			return false, err
		}

		var cantidad float64
		var precio float64
		var costo float64

		for resultSet.Next() {
			resultSet.Scan(&cantidad, &precio, &costo)
		}
		fmt.Println(cantidad, precio, costo)

		Query = fmt.Sprintf(`SELECT COUNT(*) FROM public."Cotizacion" WHERE "Operacion" = '%v' AND "Almacen"='%v' AND "Producto"='%v' `, Operacion, almacen, idProducto)
		stmt, err = SesionPsql.Prepare(Query)
		if err != nil {
			return false, err
		}

		resultSet, err = stmt.Query()
		if err != nil {
			return false, err
		}

		var Numero float64
		for resultSet.Next() {
			resultSet.Scan(&Numero)
		}

		impuesto := CalculaTotalDeImpuestoPorProducto(Impuestos, precio)
		Resto := cantidad
		fmt.Println("Resto:", Resto)

		if Numero > 0 {

			Query = fmt.Sprintf(`SELECT "Cantidad" FROM public."Cotizacion" WHERE "Operacion" = '%v' AND "Almacen" = '%v' AND "Producto" = '%v'`, Operacion, almacen, idProducto)
			con, err := BasePsql.Query(Query)
			if err != nil {
				fmt.Println(err)
				return false, err
			}
			var Prev float64
			for con.Next() {
				con.Scan(&Prev)
			}
			fmt.Println("Cantidad de Cotizacion", Prev)

			Query = fmt.Sprintf(`UPDATE  public."Cotizacion"  SET "Cantidad"= "Cantidad" + %v, "Costo"=%v, "Precio"=%v, "Existencia"=%v, "Impuesto"=%v, "Descuento"=0  WHERE "Operacion" = '%v'  AND "Almacen"='%v' AND "Producto"='%v' `, ValorNuevo, costo, precio, Resto, ValorNuevo*impuesto, Operacion, almacen, idProducto)
			_, err = SesionPsql.Exec(Query)
			if err != nil {
				fmt.Println("Ocurrio un error al insertar un nuevo producto:", err)
				SesionPsql.Rollback()
				resultSet.Close()
				stmt.Close()
				BasePsql.Close()
				return false, err
			}

		} else {

			Query = fmt.Sprintf(`INSERT INTO  public."Cotizacion"  VALUES('%v', '%v','%v', '%v', %v, %v, %v, %v, 0, %v, '%v')`, Operacion, Operacion, idProducto, almacen, ValorNuevo, costo, precio, ValorNuevo*impuesto, Resto, time.Now().Format(time.RFC3339Nano))

			_, err = SesionPsql.Exec(Query)
			if err != nil {
				fmt.Println("Ocurrio un error al insertar un nuevo producto:", err)
				SesionPsql.Rollback()
				resultSet.Close()
				stmt.Close()
				BasePsql.Close()
				return false, err
			}

		}
		SesionPsql.Commit()
		resultSet.Close()
		stmt.Close()
		BasePsql.Close()

		return true, nil

	}
	return false, err
}

//ActualizaImpuestoDeProductoEnAlmacenPsql modifica el valor de un impuesto desde la vista
func ActualizaImpuestoDeProductoEnAlmacenPsql(Operacion, Movimientox, Factor, Tipo, Almacen, Producto string, ValorPrevio, ValorNuevo, Precio float64) error {
	paramConex := InsertarDatosConexion(bson.ObjectIdHex(Almacen))
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer BasePsql.Close()

	var Movimiento string

	BasePsql.Exec("set transaction isolation level serializable")

	Query := fmt.Sprintf(`SELECT "IdMovimiento"
							FROM public."Kardex_%v"
							WHERE "IdProducto"='%v'
							AND "TipoOperacion" = '58efbf8bd2b2131778e9c928'
							ORDER BY "FechaHora" DESC LIMIT 1`,
		Almacen, Producto)
	con, err := BasePsql.Query(Query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for con.Next() {
		con.Scan(&Movimiento)
	}

	Movimiento = MoGeneral.EliminarEspaciosInicioFinal(Movimiento)
	Query = fmt.Sprintf(`UPDATE  public."ImpuestoV_%v"  
							SET "Valor"=  %v   
							WHERE  "IdMovimiento"= '%v'
							AND "IdProducto"='%v' 
							AND "TipoImpuesto"='%v' 
							AND "Factor"= '%v'`,
		Almacen, ValorNuevo, Movimiento, Producto, Tipo, Factor)
	_, err = SesionPsql.Exec(Query)
	if err != nil {
		fmt.Println(err)
		SesionPsql.Rollback()
		return err
	}
	SesionPsql.Commit()
	return nil
}

//ActualizaVentaTemporalPorImpuestos actualiza VentaTemporal por afectación de Impuestos en Ventas
func ActualizaVentaTemporalPorImpuestos(Operacion, Movimientox, Almacen, Producto string, ValorNuevo, Precio float64) error {
	paramConex := InsertarDatosConexion(bson.ObjectIdHex(Almacen))
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var Impuestos []AlmacenModel.ImpuestoAlmacen
	var Impuesto AlmacenModel.ImpuestoAlmacen
	var Valor float64
	var TipoImpuesto string
	var Factorx string

	Query := fmt.Sprintf(`SELECT "TipoImpuesto", "Factor", "Valor" FROM public."ImpuestoV_%v" 
						WHERE "IdMovimiento"= (
							SELECT "IdMovimiento"
							FROM public."Kardex_%v"
							WHERE "IdProducto"='%v'
							AND "TipoOperacion" = '58efbf8bd2b2131778e9c928'
							ORDER BY "FechaHora" DESC LIMIT 1);`, Almacen, Almacen, Producto)

	Elemento, err := BasePsql.Query(Query)
	if err != nil {
		fmt.Println(err)
		SesionPsql.Rollback()
		return err
	}

	for Elemento.Next() {
		err := Elemento.Scan(&TipoImpuesto, &Factorx, &Valor)
		if err != nil {
			fmt.Println(err)
			SesionPsql.Rollback()
			return err
		}
		Impuesto = AlmacenModel.ImpuestoAlmacen{}
		Impuesto.Tipo = MoGeneral.EliminarEspaciosInicioFinal(TipoImpuesto)
		Impuesto.Factor = MoGeneral.EliminarEspaciosInicioFinal(Factorx)
		Impuesto.Valor = Valor
		Impuestos = append(Impuestos, Impuesto)
	}

	ImpuestoTotal := CalculaTotalDeImpuestoPorProducto(Impuestos, Precio)

	Query = fmt.Sprintf(`UPDATE  public."VentaTemporal"  SET "Impuesto"= "Cantidad" * %v  WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, ImpuestoTotal, Operacion, Movimientox, Almacen, Producto)
	_, err = SesionPsql.Exec(Query)
	if err != nil {
		fmt.Println(err)
		SesionPsql.Rollback()
		return err
	}
	SesionPsql.Commit()
	return nil
}

//CalculaTotalDeImpuestoPorProducto calcula el total de un grupo de impuestos con un precio dado
func CalculaTotalDeImpuestoPorProducto(Impuestos []AlmacenModel.ImpuestoAlmacen, Precio float64) float64 {
	var Total float64
	for _, v := range Impuestos {
		if v.Factor == "Tasa" {
			Total += Precio * (v.Valor / 100)
		} else if v.Factor == "Cuota" {
			Total += v.Valor
		} else if v.Factor == "Exento" {
			Total += 0
		} else {
			Total += 0
		}
	}
	return Total
}

//ConsultaPrecioExistenciaYActualizaProductoEnAlmacenModal consulta si existe un producto y devuelve la cantidad de existencias y el precio
func ConsultaPrecioExistenciaYActualizaProductoEnAlmacenModal(Operacion string, Movimiento, almacen, idProducto, ValorNuevo []string) error {
	var bandera bool
	var ValoresPrevios []float64
	var ValoresImpuestos [][]AlmacenModel.ImpuestoAlmacen

	var Impuestos []AlmacenModel.ImpuestoAlmacen
	var Impuesto AlmacenModel.ImpuestoAlmacen
	var Valor float64
	var TipoImpuesto string
	var Factor string

	fmt.Println("Almacenes:", almacen)
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		fmt.Println(err)
	}

	for i, v := range idProducto {
		Impuestos = []AlmacenModel.ImpuestoAlmacen{}
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

		Query = fmt.Sprintf(`SELECT "TipoImpuesto", "Factor", "Valor" FROM public."ImpuestoV_%v" 
						WHERE "IdMovimiento"= (
							SELECT "IdMovimiento"
							FROM public."Kardex_%v"
							WHERE "IdProducto"='%v'
							AND "TipoOperacion" = '58efbf8bd2b2131778e9c928'
							ORDER BY "FechaHora" DESC LIMIT 1);`, almacen[i], almacen[i], v)

		// Query = fmt.Sprintf(`SELECT "TipoImpuesto", "Factor", "Valor" FROM public."ImpuestoV_%v" AS IMP
		// 				JOIN public."Kardex_%v" AS KAR ON IMP."IdMovimiento"= KAR."IdMovimiento"
		// 				WHERE IMP."IdProducto"='%v' AND KAR."FechaHora" = (SELECT MAX("FechaHora") FROM public."Kardex_%v")
		// 				GROUP BY  IMP."TipoImpuesto", IMP."Valor", IMP."Factor", KAR."FechaHora" ORDER BY KAR."FechaHora"   DESC;`, almacen[i], almacen[i], v, almacen[i])

		con, err = BasePsql.Query(Query)
		if err != nil {
			fmt.Println(err)
		}

		for con.Next() {
			err := con.Scan(&TipoImpuesto, &Factor, &Valor)
			if err != nil {
				fmt.Println(err)
			}
			Impuesto = AlmacenModel.ImpuestoAlmacen{}
			Impuesto.Tipo = MoGeneral.EliminarEspaciosInicioFinal(TipoImpuesto)
			Impuesto.Factor = MoGeneral.EliminarEspaciosInicioFinal(Factor)
			Impuesto.Valor = Valor
			Impuestos = append(Impuestos, Impuesto)
		}

		ValoresImpuestos = append(ValoresImpuestos, Impuestos)
	}

	for k, v := range idProducto {
		esta, err := ConsultaExistenciaProductoActivo(almacen[k], v)
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
			VN = ValoresPrevios[k] + VN
			Resto := cantidad - VN

			if Resto >= 0 {

				ValorImpuesto := CalculaTotalDeImpuestoPorProducto(ValoresImpuestos[k], precio)

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
					Query = fmt.Sprintf(`UPDATE  public."VentaTemporal"  SET "Cantidad"= "Cantidad" + %v, "Costo"=%v, "Precio"=%v, "Existencia"=%v, "Impuesto"=%v, "Descuento"=0  WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, VN, costo, precio, Resto, ValorImpuesto*VN, Operacion, Movimiento[k], almacen[k], v)
				} else {
					Query = fmt.Sprintf(`INSERT INTO  public."VentaTemporal"  VALUES('%v', '%v', '%v', '%v', %v, %v, %v, %v, 0, %v, '%v')`, Operacion, Movimiento[k], v, almacen[k], VN, costo, precio, ValorImpuesto*VN, Resto, time.Now().Format(time.RFC3339Nano))
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
		return err
	}

	return nil

}

//EliminaProductoCarritoYActualizaInventarioAlmacen elimina el producto del almacen temporal y afecta el inventario
func EliminaProductoCarritoYActualizaInventarioAlmacen(Operacion, idProducto string) (bool, error) {
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		return false, err
	}

	BasePsql.Exec("set transaction isolation level serializable")
	Query := fmt.Sprintf(`SELECT "Operacion", "Movimiento","Producto","Almacen","Cantidad" FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Producto" = '%v'  FOR UPDATE`, Operacion, idProducto)
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
	}

	Query = fmt.Sprintf(`UPDATE  public."Inventario_%v"  SET  "Existencia" = "Existencia" + %v WHERE "IdProducto" ='%v'`, almacen, cantidad, idproducto)
	_, err = SesionPsql.Exec(Query)
	if err != nil {
		SesionPsql.Rollback()
		return false, err
	}

	Query = fmt.Sprintf(`DELETE FROM public."VentaTemporal" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, operacion, movimiento, almacen, idproducto)
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

//ConsultaExistenciaProductoEnAlmacenCotizacion consulta si existe un producto
func ConsultaExistenciaProductoEnAlmacenCotizacion(IDCotizacion, idProducto string) (bool, error) {
	var paramConex MoConexion.ParametrosConexionPostgres
	paramConex.Servidor = "192.168.1.110"
	paramConex.Usuario = "postgres"
	paramConex.Pass = "12345"
	paramConex.NombreBase = "IsmaelMinisuperAmpliado"

	BasePsql, err := MoConexion.ConexionEspecificaPsql(paramConex)
	if err != nil {
		return false, err
	}
	defer BasePsql.Close()
	Query := `SELECT COUNT(*) FROM "Cotizacion" WHERE "Cotizacion"='` + IDCotizacion + `' and "Producto"='` + idProducto + `'`
	Elemento, err := BasePsql.Query(Query)
	if err != nil {
		fmt.Println("Error al consultar las existencias de in producto en Cotización", err)
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

//ConsultaPrecioExistenciaYActualizaProductoEnCotizacion actualiza un producto de una cotizacion, de no existir lo inserta
func ConsultaPrecioExistenciaYActualizaProductoEnCotizacion(Operacion, Movimiento, almacen, idProducto string, ValorNuevo float64) (bool, error) {

	esta, err := ConsultaExistenciaProductoActivo(almacen, idProducto)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ocurrió un error al consultar EXISTENCIA EN el producto : ", idProducto, "En el almacén: ", almacen)
		return false, err
	}

	Impuestos, err := ConsultaValorImpuestoEnAlmacen(almacen, idProducto)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Ocurrió un error al consultar el Impuesto del producto : ", idProducto, " en el Almacen: ", almacen)
	}

	if esta {
		paramConex := InsertarDatosConexion(bson.ObjectIdHex(almacen))
		BasePsql, SesionPsql, err := MoConexion.IniciaSesionEspecificaPsql(paramConex)
		if err != nil {
			return false, err
		}

		Query := fmt.Sprintf(`SELECT "Existencia", "Precio", "Costo" FROM public."Inventario_%v" WHERE "IdProducto" = '%v' and "Estatus" = 'ACTIVO'`, almacen, idProducto)
		stmt, err := SesionPsql.Prepare(Query)
		if err != nil {
			return false, err
		}

		resultSet, err := stmt.Query()
		if err != nil {
			return false, err
		}

		var cantidad float64
		var precio float64
		var costo float64

		for resultSet.Next() {
			resultSet.Scan(&cantidad, &precio, &costo)
		}

		impuesto := CalculaTotalDeImpuestoPorProducto(Impuestos, precio)

		Query = fmt.Sprintf(`SELECT COUNT(*) FROM public."Cotizacion" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, Operacion, Movimiento, almacen, idProducto)
		stmt, err = SesionPsql.Prepare(Query)
		if err != nil {
			return false, err
		}

		resultSet, err = stmt.Query()
		if err != nil {
			return false, err
		}

		var Numero float64
		for resultSet.Next() {
			resultSet.Scan(&Numero)
		}

		if Numero > 0 {
			Query = fmt.Sprintf(`UPDATE  public."Cotizacion"  SET "Cantidad"= %v, "Costo"=%v, "Precio"=%v, "Existencia"=%v, "Impuesto"=%v, "Descuento"=0  WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, ValorNuevo, costo, precio, cantidad, ValorNuevo*impuesto, Operacion, Movimiento, almacen, idProducto)
		} else {
			Query = fmt.Sprintf(`INSERT INTO  public."VentaTemporal"  VALUES('%v', '%v', '%v', '%v', %v, %v, %v, %v, 0, %v, '%v')`, Operacion, Movimiento, idProducto, almacen, ValorNuevo, costo, precio, ValorNuevo*impuesto, cantidad, time.Now().Format(time.RFC3339Nano))
		}
		_, err = SesionPsql.Exec(Query)
		if err != nil {
			SesionPsql.Rollback()
			resultSet.Close()
			stmt.Close()
			BasePsql.Close()
			return false, err
		}

		SesionPsql.Commit()
		resultSet.Close()
		stmt.Close()
		BasePsql.Close()

		return true, nil

	}
	return false, nil
}

//EliminaProductoCarritoYActualizaInventarioCotizacion elimina el producto del almacen temporal y afecta el inventario
func EliminaProductoCarritoYActualizaInventarioCotizacion(Operacion, Movimiento, Almacen, idProducto string) (bool, error) {
	BasePsql, SesionPsql, err := MoConexion.IniciaSesionPsql()
	if err != nil {
		return false, err
	}

	Query := fmt.Sprintf(`DELETE FROM public."Cotizacion" WHERE "Operacion" = '%v' AND "Movimiento"= '%v' AND "Almacen"='%v' AND "Producto"='%v' `, Operacion, Movimiento, Almacen, idProducto)
	stmt, err := SesionPsql.Prepare(Query)
	if err != nil {
		SesionPsql.Rollback()
		return false, err
	}

	resultSet, err := stmt.Query()
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
