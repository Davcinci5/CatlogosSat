package main

import (
	"fmt"

	"./src/Controladores/AlmacenControler"
	"./src/Controladores/BugControler"
	"./src/Controladores/CajaControler"
	"./src/Controladores/CatalogoControler"
	"./src/Controladores/ClienteControler"
	"./src/Controladores/CompraControler"
	"./src/Controladores/ConexionControler"
	"./src/Controladores/ConfiguracionControler"
	"./src/Controladores/CotizacionControler"
	"./src/Controladores/DispositivoControler"
	"./src/Controladores/EmpresaControler"
	"./src/Controladores/EquipoCajaControler"
	"./src/Controladores/FacturacionControler"
	"./src/Controladores/GrupoAlmacenControler"
	"./src/Controladores/GrupoControler"
	"./src/Controladores/GrupoPersonaControler"
	"./src/Controladores/HTTPErrorsControler"
	"./src/Controladores/ImpuestoControler"
	"./src/Controladores/LoginControler"
	"./src/Controladores/MediosPagoControler"
	"./src/Controladores/PermisosUriControler"
	"./src/Controladores/PersonaControler"
	"./src/Controladores/ProductoControler"
	"./src/Controladores/PromocionControler"
	"./src/Controladores/PuntoVentaControler"
	"./src/Controladores/RolControler"
	"./src/Controladores/SesionesController"
	"./src/Controladores/TrasladoAjusteControler"
	"./src/Controladores/UnidadControler"
	"./src/Controladores/UsuarioControler"

	"./src/Modulos/Variables"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

func main() {

	//###################### Start ####################################
	app := iris.New()
	app.Adapt(httprouter.New())
	sesiones := sessions.New(sessions.Config{Cookie: "cookiekore"})
	app.Adapt(sesiones)

	app.Adapt(view.HTML("./Public/Vistas", ".html").Reload(true))

	app.Set(iris.OptionCharset("UTF-8"))

	app.StaticWeb("/icono", "./Public/Recursos/Generales/img")

	app.StaticWeb("/css", "./Public/Recursos/Generales/css")
	app.StaticWeb("/js", "./Public/Recursos/Generales/js")
	app.StaticWeb("/Plugins", "./Public/Recursos/Generales/Plugins")
	app.StaticWeb("/scripts", "./Public/Recursos/Generales/scripts")
	app.StaticWeb("/img", "./Public/Recursos/Generales/img")
	app.StaticWeb("/Comprobantes", "./Public/Recursos/Locales/comprobantes")
	app.StaticWeb("/Locales", "./Public/Recursos/Locales")

	//###################### CFG ######################################

	var DataCfg = MoVar.CargaSeccionCFG(MoVar.SecDefault)

	//###################### Ruteo ####################################

	//###################### HTTP ERRORS ################################

	//Error 403
	app.OnError(iris.StatusForbidden, HTTPErrors.Error403)
	//Error 404
	app.OnError(iris.StatusNotFound, HTTPErrors.Error404)

	//Error 500
	app.OnError(iris.StatusInternalServerError, HTTPErrors.Error500)

	// PRUEBAS DE ERRORES
	app.Get("/500", func(ctx *iris.Context) {
		ctx.EmitError(iris.StatusInternalServerError)
	})
	app.Get("/404", func(ctx *iris.Context) {
		ctx.EmitError(iris.StatusNotFound)
	})

	//###################### Login ################################
	//Login

	app.Get("/Login/saveMac", LoginControler.GetMac)

	// get
	app.Get("/", LoginControler.LoginGet)
	//Post
	app.Post("/", LoginControler.LoginPost)

	//Password Recovery
	//Get
	app.Get("/Recuperar", LoginControler.RecuperarGet)

	//Post
	app.Post("/Recuperar", LoginControler.RecuperarPost)

	//NuevoLogin
	//Get
	app.Get("/NuevoLogin", LoginControler.NuevoGet)

	app.Get("/Logout", LoginControler.Logout)
	//Post

	//Index
	app.Get("/Index", LoginControler.IndexGet)

	//###################### Catalogo ################################
	//Index (Búsqueda)
	app.Get("/Catalogos", CatalogoControler.IndexGet)
	app.Post("/Catalogos", CatalogoControler.IndexPost)
	app.Post("/Catalogos/search", CatalogoControler.BuscaPagina)
	app.Post("/Catalogos/agrupa", CatalogoControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Catalogos/alta", CatalogoControler.AltaGet)
	app.Post("/Catalogos/alta", CatalogoControler.AltaPost)

	//Edicion
	app.Get("/Catalogos/edita", CatalogoControler.EditaGet)
	app.Post("/Catalogos/edita", CatalogoControler.EditaPost)
	app.Get("/Catalogos/edita/:ID", CatalogoControler.EditaGet)
	app.Post("/Catalogos/edita/:ID", CatalogoControler.EditaPost)

	//Detalle
	app.Get("/Catalogos/detalle", CatalogoControler.DetalleGet)
	app.Post("/Catalogos/detalle", CatalogoControler.DetallePost)
	app.Get("/Catalogos/detalle/:ID", CatalogoControler.DetalleGet)
	app.Post("/Catalogos/detalle/:ID", CatalogoControler.DetallePost)

	//Rutinas adicionales

	//###################### Unidad ################################
	//Index (Búsqueda)
	app.Get("/Unidades", UnidadControler.Index)
	app.Post("/Unidades", UnidadControler.Index)

	//Alta
	app.Get("/Unidades/alta", UnidadControler.AltaGet)
	app.Post("/Unidades/alta", UnidadControler.AltaPost)

	//Edicion
	app.Get("/Unidades/edita", UnidadControler.EditaGet)
	app.Post("/Unidades/edita", UnidadControler.EditaPost)
	app.Get("/Unidades/edita/:ID", UnidadControler.EditaGet)
	app.Post("/Unidades/edita/:ID", UnidadControler.EditaPost)

	//Detalle
	app.Get("/Unidades/detalle", UnidadControler.DetalleGet)
	app.Post("/Unidades/detalle", UnidadControler.DetallePost)
	app.Get("/Unidades/detalle/:ID", UnidadControler.DetalleGet)
	app.Post("/Unidades/detalle/:ID", UnidadControler.DetallePost)

	//Rutinas adicionales
	app.Post("/ConsultaUnidadesDeMagnitud", UnidadControler.ConsultaUnidadesDeMagnitud)
	app.Post("/AgregaMagnitudDesdeUnidades", UnidadControler.AgregaMagnitudDesdeUnidades)

	//###################### Cliente ################################
	//Index (Búsqueda)
	app.Get("/Clientes", ClienteControler.IndexGet)
	app.Post("/Clientes", ClienteControler.IndexPost)
	app.Post("/Clientes/search", ClienteControler.BuscaPagina)
	//Alta
	app.Get("/Clientes/alta", ClienteControler.AltaGet)
	app.Post("/Clientes/alta", ClienteControler.AltaPost)

	//Edicion
	app.Get("/Clientes/edita", ClienteControler.EditaGet)
	app.Post("/Clientes/edita", ClienteControler.EditaPost)
	app.Get("/Clientes/edita/:ID", ClienteControler.EditaGet)
	app.Post("/Clientes/edita/:ID", ClienteControler.EditaPost)

	//Detalle
	app.Get("/Clientes/detalle", ClienteControler.DetalleGet)
	app.Post("/Clientes/detalle", ClienteControler.DetallePost)
	app.Get("/Clientes/detalle/:ID", ClienteControler.DetalleGet)
	app.Post("/Clientes/detalle/:ID", ClienteControler.DetallePost)

	//Rutinas adicionales
	app.Post("/Clientes/GetMunicipiosForClaveEstado", ClienteControler.GetMunicipiosForClaveEstado)
	app.Post("/Clientes/GetColoniasForClaveMunicipio", ClienteControler.GetColoniasForClaveMunicipio)
	app.Post("/Clientes/GetCPForColonia", ClienteControler.GetCPForColonia)
	app.Post("/Clientes/GetEstadosForSelect", ClienteControler.GetEstadosForSelect)
	app.Post("/Clientes/GetTipoDireciones", ClienteControler.GetTipoDireciones)

	//###################### Producto ################################
	//Index (Búsqueda)
	app.Get("/Productos", ProductoControler.IndexGet)
	app.Post("/Productos", ProductoControler.IndexPost)
	app.Post("/Productos/search", ProductoControler.BuscaPagina)
	app.Post("/Productos/agrupa", ProductoControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Productos/alta", ProductoControler.AltaGet)
	app.Post("/Productos/alta", ProductoControler.AltaPost)

	//Edicion
	app.Get("/Productos/edita", ProductoControler.EditaGet)
	app.Post("/Productos/edita", ProductoControler.EditaPost)
	app.Get("/Productos/edita/:ID", ProductoControler.EditaGet)
	app.Post("/Productos/edita/:ID", ProductoControler.EditaPost)

	//Detalle
	app.Get("/Productos/detalle", ProductoControler.DetalleGet)
	app.Post("/Productos/detalle", ProductoControler.DetallePost)
	app.Get("/Productos/detalle/:ID", ProductoControler.DetalleGet)
	app.Post("/Productos/detalle/:ID", ProductoControler.DetallePost)

	//###################### Almacen ################################
	//Index (Búsqueda)
	app.Get("/Almacens", AlmacenControler.IndexGet)
	app.Post("/Almacens", AlmacenControler.IndexPost)
	app.Post("/Almacens/search", AlmacenControler.BuscaPagina)
	app.Post("/Almacens/agrupa", AlmacenControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Almacens/alta", AlmacenControler.AltaGet)
	app.Post("/Almacens/alta", AlmacenControler.AltaPost)

	//Edicion
	app.Get("/Almacens/edita", AlmacenControler.EditaGet)
	app.Post("/Almacens/edita", AlmacenControler.EditaPost)
	app.Get("/Almacens/edita/:ID", AlmacenControler.EditaGet)
	app.Post("/Almacens/edita/:ID", AlmacenControler.EditaPost)

	//Detalle
	app.Get("/Almacens/detalle", AlmacenControler.DetalleGet)
	app.Post("/Almacens/detalle", AlmacenControler.DetallePost)
	app.Get("/Almacens/detalle/:ID", AlmacenControler.DetalleGet)
	app.Post("/Almacens/detalle/:ID", AlmacenControler.DetallePost)

	//Rutinas adicionales
	app.Get("/Almacens/Operacion/Ajuste", AlmacenControler.MovimientoAjusteGet)
	app.Get("/Almacens/Operacion/Traslado", AlmacenControler.MovimientoTrasladoGet)

	//################ Operaciones -.- Movimientos #######################

	app.Get("/Almacens/Operacion/Movimientos", AlmacenControler.MovimientosGet)
	app.Post("/Almacens/Operacion/Movimientos", AlmacenControler.MovimientosPost)
	app.Post("/Almacens/Operacion/Movimientos/Ajustar", AlmacenControler.MovimientosAjustePost)

	app.Post("/Almacens/Operacion/Movimientos/Realizar", AlmacenControler.RealizarMovimiento)

	//################ Operaciones---------traslados Ajustes #####################
	app.Get("/TrasladoAjuste/Ajuste", TrasladoAjusteControler.MovimientoAjusteGet)
	app.Get("/TrasladoAjuste/Traslado", TrasladoAjusteControler.MovimientoTrasladoGet)
	app.Get("/TrasladoAjuste/Movimientos", TrasladoAjusteControler.MovimientosGet)
	//app.Post("/TrasladoAjuste/Movimientos", TrasladoAjusteControler.MovimientosPost)
	app.Post("/TrasladoAjuste/Ajustar", TrasladoAjusteControler.MovimientosAjustePost)
	app.Post("/TrasladoAjuste/GetArticuloAjustar", TrasladoAjusteControler.MovimientosAjustePost)

	app.Post("/TrasladoAjuste/Realizar", TrasladoAjusteControler.RealizarMovimiento)

	app.Get("/TrasladoAjuste/Productos", TrasladoAjusteControler.ConsultarProductos)
	app.Post("/TrasladoAjuste/Productos", TrasladoAjusteControler.ConsultarProductos)

	//###################### GrupoPersona ################################
	//Index (Búsqueda)
	app.Get("/GrupoPersonas", GrupoPersonaControler.IndexGet)
	app.Post("/GrupoPersonas", GrupoPersonaControler.IndexPost)

	//Alta
	app.Get("/GrupoPersonas/alta", GrupoPersonaControler.AltaGet)
	app.Post("/GrupoPersonas/alta", GrupoPersonaControler.AltaPost)

	//Edicion
	app.Get("/GrupoPersonas/edita", GrupoPersonaControler.EditaGet)
	app.Post("/GrupoPersonas/edita", GrupoPersonaControler.EditaPost)
	app.Get("/GrupoPersonas/edita/:ID", GrupoPersonaControler.EditaGet)
	app.Post("/GrupoPersonas/edita/:ID", GrupoPersonaControler.EditaPost)

	//Detalle
	app.Get("/GrupoPersonas/detalle", GrupoPersonaControler.DetalleGet)
	app.Post("/GrupoPersonas/detalle", GrupoPersonaControler.DetallePost)
	app.Get("/GrupoPersonas/detalle/:ID", GrupoPersonaControler.DetalleGet)
	app.Post("/GrupoPersonas/detalle/:ID", GrupoPersonaControler.DetallePost)

	//Rutinas adicionales

	/*
		//###################### ListaCosto ################################
		//Index (Búsqueda)
		app.Get("/ListaCostos", ListaCostoControler.IndexGet)
		app.Post("/ListaCostos", ListaCostoControler.IndexPost)

		//Alta
		app.Get("/ListaCostos/alta", ListaCostoControler.AltaGet)
		app.Post("/ListaCostos/alta", ListaCostoControler.AltaPost)

		//Edicion
		app.Get("/ListaCostos/edita", ListaCostoControler.EditaGet)
		app.Post("/ListaCostos/edita", ListaCostoControler.EditaPost)
		app.Get("/ListaCostos/edita/:ID", ListaCostoControler.EditaGet)
		app.Post("/ListaCostos/edita/:ID", ListaCostoControler.EditaPost)

		//Detalle
		app.Get("/ListaCostos/detalle", ListaCostoControler.DetalleGet)
		app.Post("/ListaCostos/detalle", ListaCostoControler.DetallePost)
		app.Get("/ListaCostos/detalle/:ID", ListaCostoControler.DetalleGet)
		app.Post("/ListaCostos/detalle/:ID", ListaCostoControler.DetallePost)

		//Rutinas adicionales

		//###################### ListaPrecio ################################
		//Index (Búsqueda)
		app.Get("/ListaPrecios", ListaPrecioControler.IndexGet)
		app.Post("/ListaPrecios", ListaPrecioControler.IndexPost)

		//Alta
		app.Get("/ListaPrecios/alta", ListaPrecioControler.AltaGet)
		app.Post("/ListaPrecios/alta", ListaPrecioControler.AltaPost)

		//Edicion
		app.Get("/ListaPrecios/edita", ListaPrecioControler.EditaGet)
		app.Post("/ListaPrecios/edita", ListaPrecioControler.EditaPost)
		app.Get("/ListaPrecios/edita/:ID", ListaPrecioControler.EditaGet)
		app.Post("/ListaPrecios/edita/:ID", ListaPrecioControler.EditaPost)

		//Detalle
		app.Get("/ListaPrecios/detalle", ListaPrecioControler.DetalleGet)
		app.Post("/ListaPrecios/detalle", ListaPrecioControler.DetallePost)
		app.Get("/ListaPrecios/detalle/:ID", ListaPrecioControler.DetalleGet)
		app.Post("/ListaPrecios/detalle/:ID", ListaPrecioControler.DetallePost)

		//Rutinas adicionales
	*/
	//###################### Empresa ################################
	//Index (Búsqueda)

	app.Get("/Empresas", EmpresaControler.EditaGet)
	app.Post("/Empresas", EmpresaControler.EditaPost)

	//Alta
	app.Get("/Empresas/alta", EmpresaControler.AltaGet)
	app.Post("/Empresas/alta", EmpresaControler.AltaPost)

	//Edicion
	app.Get("/Empresas/edita", EmpresaControler.EditaGet)
	app.Post("/Empresas/edita", EmpresaControler.EditaPost)
	app.Get("/Empresas/edita/:ID", EmpresaControler.EditaGet)
	app.Post("/Empresas/edita/:ID", EmpresaControler.EditaPost)

	//Detalle
	app.Get("/Empresas/detalle", EmpresaControler.DetalleGet)
	app.Post("/Empresas/detalle", EmpresaControler.DetallePost)
	app.Get("/Empresas/detalle/:ID", EmpresaControler.DetalleGet)
	app.Post("/Empresas/detalle/:ID", EmpresaControler.DetallePost)

	//Rutinas adicionales

	//###################### Caja ################################
	//Index (Búsqueda)
	app.Get("/Cajas", CajaControler.IndexGet)
	app.Post("/Cajas", CajaControler.IndexPost)
	app.Post("/Cajas/search", CajaControler.BuscaPagina)
	app.Post("/Cajas/agrupa", CajaControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Cajas/alta", CajaControler.AltaGet)
	app.Post("/Cajas/alta", CajaControler.AltaPost)

	//Edicion
	app.Get("/Cajas/edita", CajaControler.EditaGet)
	app.Post("/Cajas/edita", CajaControler.EditaPost)
	app.Get("/Cajas/edita/:ID", CajaControler.EditaGet)
	app.Post("/Cajas/edita/:ID", CajaControler.EditaPost)

	//Detalle
	app.Get("/Cajas/detalle", CajaControler.DetalleGet)
	app.Post("/Cajas/detalle", CajaControler.DetallePost)
	app.Get("/Cajas/detalle/:ID", CajaControler.DetalleGet)
	app.Post("/Cajas/detalle/:ID", CajaControler.DetallePost)

	//Rutinas adicionales
	//Busca documento para COBRO/PAGO
	app.Post("/Cajas/buscaDocumento", CajaControler.BuscaDocumentoPost)
	//Inserta movimiento de PAGO
	app.Post("/Cajas/insertaDocumento", CajaControler.InsertaDocumento)
	//Cierra CAJA
	app.Post("/Cajas/cierraCaja", CajaControler.CierraCaja)
	app.Post("/Cajas/getSaldo", CajaControler.GetSaldoCaja)
	//###################### Dispositivo ################################
	//Index (Búsqueda)
	app.Get("/Dispositivos", DispositivoControler.IndexGet)
	app.Post("/Dispositivos", DispositivoControler.IndexPost)

	//Alta
	app.Get("/Dispositivos/alta", DispositivoControler.AltaGet)
	app.Post("/Dispositivos/alta", DispositivoControler.AltaPost)

	//Edicion
	app.Get("/Dispositivos/edita", DispositivoControler.EditaGet)
	app.Post("/Dispositivos/edita", DispositivoControler.EditaPost)
	app.Get("/Dispositivos/edita/:ID", DispositivoControler.EditaGet)
	app.Post("/Dispositivos/edita/:ID", DispositivoControler.EditaPost)

	//Detalle
	app.Get("/Dispositivos/detalle", DispositivoControler.DetalleGet)
	app.Post("/Dispositivos/detalle", DispositivoControler.DetallePost)
	app.Get("/Dispositivos/detalle/:ID", DispositivoControler.DetalleGet)
	app.Post("/Dispositivos/detalle/:ID", DispositivoControler.DetallePost)

	//Rutinas adicionales

	//###################### Impuesto ################################
	//Index (Búsqueda)
	app.Get("/Impuestos", ImpuestoControler.IndexGet)
	app.Post("/Impuestos", ImpuestoControler.IndexPost)
	app.Post("/Impuestos/search", ImpuestoControler.BuscaPagina)
	app.Post("/Impuestos/agrupa", ImpuestoControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Impuestos/alta", ImpuestoControler.AltaGet)
	app.Post("/Impuestos/alta", ImpuestoControler.AltaPost)

	//Edicion
	app.Get("/Impuestos/edita", ImpuestoControler.EditaGet)
	app.Post("/Impuestos/edita", ImpuestoControler.EditaPost)
	app.Get("/Impuestos/edita/:ID", ImpuestoControler.EditaGet)
	app.Post("/Impuestos/edita/:ID", ImpuestoControler.EditaPost)

	//Detalle
	app.Get("/Impuestos/detalle", ImpuestoControler.DetalleGet)
	app.Post("/Impuestos/detalle", ImpuestoControler.DetallePost)
	app.Get("/Impuestos/detalle/:ID", ImpuestoControler.DetalleGet)
	app.Post("/Impuestos/detalle/:ID", ImpuestoControler.DetallePost)

	/*
		//###################### Kit ################################
		//Index (Búsqueda)
		app.Get("/Kits", KitControler.IndexGet)
		app.Post("/Kits", KitControler.IndexPost)

		//Alta
		app.Get("/Kits/alta", KitControler.AltaGet)
		app.Post("/Kits/alta", KitControler.AltaPost)

		//Edicion
		app.Get("/Kits/edita", KitControler.EditaGet)
		app.Post("/Kits/edita", KitControler.EditaPost)
		app.Get("/Kits/edita/:ID", KitControler.EditaGet)
		app.Post("/Kits/edita/:ID", KitControler.EditaPost)

		//Detalle
		app.Get("/Kits/detalle", KitControler.DetalleGet)
		app.Post("/Kits/detalle", KitControler.DetallePost)
		app.Get("/Kits/detalle/:ID", KitControler.DetalleGet)
		app.Post("/Kits/detalle/:ID", KitControler.DetallePost)

		//Rutinas adicionales
	*/
	//###################### MediosPago ################################
	//Index (Búsqueda)
	app.Get("/MediosPagos", MediosPagoControler.IndexGet)
	app.Post("/MediosPagos", MediosPagoControler.IndexPost)
	app.Post("/MediosPagos/search", MediosPagoControler.BuscaPagina)
	app.Post("/MediosPagos/agrupa", MediosPagoControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/MediosPagos/alta", MediosPagoControler.AltaGet)
	app.Post("/MediosPagos/alta", MediosPagoControler.AltaPost)

	//Edicion
	app.Get("/MediosPagos/edita", MediosPagoControler.EditaGet)
	app.Post("/MediosPagos/edita", MediosPagoControler.EditaPost)
	app.Get("/MediosPagos/edita/:ID", MediosPagoControler.EditaGet)
	app.Post("/MediosPagos/edita/:ID", MediosPagoControler.EditaPost)

	//Detalle
	app.Get("/MediosPagos/detalle", MediosPagoControler.DetalleGet)
	app.Post("/MediosPagos/detalle", MediosPagoControler.DetallePost)
	app.Get("/MediosPagos/detalle/:ID", MediosPagoControler.DetalleGet)
	app.Post("/MediosPagos/detalle/:ID", MediosPagoControler.DetallePost)

	//Rutinas adicionales

	//###################### Persona ################################
	//Index (Búsqueda)
	app.Get("/Personas", PersonaControler.IndexGet)
	app.Post("/Personas", PersonaControler.IndexPost)
	app.Post("/Personas/search", PersonaControler.BuscaPagina)
	app.Post("/Personas/agrupa", PersonaControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Personas/alta", PersonaControler.AltaGet)
	app.Post("/Personas/alta", PersonaControler.AltaPost)

	//Edicion
	app.Get("/Personas/edita", PersonaControler.EditaGet)
	app.Post("/Personas/edita", PersonaControler.EditaPost)
	app.Get("/Personas/edita/:ID", PersonaControler.EditaGet)
	app.Post("/Personas/edita/:ID", PersonaControler.EditaPost)

	//Detalle
	app.Get("/Personas/detalle", PersonaControler.DetalleGet)
	app.Post("/Personas/detalle", PersonaControler.DetallePost)
	app.Get("/Personas/detalle/:ID", PersonaControler.DetalleGet)
	app.Post("/Personas/detalle/:ID", PersonaControler.DetallePost)

	//Rutinas adicionales

	//###################### Rol ################################
	//Index (Búsqueda)
	app.Get("/Sesiones", SesionesController.IndexGet)
	app.Post("/Sesiones", SesionesController.IndexPost)

	app.Get("/Sesiones/detalle", SesionesController.DetalleGet)
	app.Post("/Sesiones/detalle", SesionesController.DetallePost)
	app.Get("/Sesiones/detalle/:ID", SesionesController.DetalleGet)
	app.Post("/Sesiones/detalle/:ID", SesionesController.DetallePost)

	app.Get("/Sesiones/edita", SesionesController.EditaGet)
	app.Get("/Sesiones/edita/:ID", SesionesController.EditaGet)

	app.Get("/Sesiones/Activas", SesionesController.SesionesTotal)

	//Peticiones AJAX
	app.Get("/Sesiones/EliminaByID", SesionesController.EliminaByID)
	app.Get("/Sesiones/EliminaByName", SesionesController.EliminaByName)

	//Rutinas adicionales

	//###################### Rol ################################
	//Index (Búsqueda)
	app.Get("/Rols", RolControler.IndexGet)
	app.Post("/Rols", RolControler.IndexPost)
	app.Post("/Rols/search", RolControler.BuscaPagina)
	app.Post("/Rols/agrupa", RolControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Rols/alta", RolControler.AltaGet)
	app.Post("/Rols/alta", RolControler.AltaPost)

	//Edicion
	app.Get("/Rols/edita", RolControler.EditaGet)
	app.Post("/Rols/edita", RolControler.EditaPost)
	app.Get("/Rols/edita/:ID", RolControler.EditaGet)
	app.Post("/Rols/edita/:ID", RolControler.EditaPost)

	//Detalle
	app.Get("/Rols/detalle", RolControler.DetalleGet)
	app.Post("/Rols/detalle", RolControler.DetallePost)
	app.Get("/Rols/detalle/:ID", RolControler.DetalleGet)
	app.Post("/Rols/detalle/:ID", RolControler.DetallePost)

	//Rutinas adicionales

	//###################### Usuario ################################
	//Index (Búsqueda)
	app.Get("/Usuarios", UsuarioControler.IndexGet)
	app.Post("/Usuarios", UsuarioControler.IndexPost)
	app.Post("/Usuarios/search", UsuarioControler.BuscaPagina)
	app.Post("/Usuarios/agrupa", UsuarioControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Usuarios/alta", UsuarioControler.AltaGet)
	app.Post("/Usuarios/alta", UsuarioControler.AltaPost)

	//Edicion
	app.Get("/Usuarios/edita", UsuarioControler.EditaGet)
	app.Post("/Usuarios/edita", UsuarioControler.EditaPost)
	app.Get("/Usuarios/edita/:ID", UsuarioControler.EditaGet)
	app.Post("/Usuarios/edita/:ID", UsuarioControler.EditaPost)

	//Detalle
	app.Get("/Usuarios/detalle", UsuarioControler.DetalleGet)
	app.Post("/Usuarios/detalle", UsuarioControler.DetallePost)
	app.Get("/Usuarios/detalle/:ID", UsuarioControler.DetalleGet)
	app.Post("/Usuarios/detalle/:ID", UsuarioControler.DetallePost)

	//ModificarCuentaUsuario
	app.Get("/ModificarCuenta", UsuarioControler.EditaPropioGet)

	//Rutinas adicionales
	//VER PERFIL
	app.Get("/Cuenta/:ID", UsuarioControler.Perfil)
	//ADMINISTRAR USUARIO
	app.Get("/AdministrarUsuarios/", UsuarioControler.AdminUsers)
	//NOTIFICACIONES
	app.Get("/Notificaciones/:ID", UsuarioControler.NotificacionesDeUsuario)

	//###################### EquipoCaja ################################
	//Index (Búsqueda)
	app.Get("/EquipoCajas", EquipoCajaControler.IndexGet)
	app.Post("/EquipoCajas", EquipoCajaControler.IndexPost)

	//Alta
	app.Get("/EquipoCajas/alta", EquipoCajaControler.AltaGet)
	app.Post("/EquipoCajas/alta", EquipoCajaControler.AltaPost)

	//Edicion
	app.Get("/EquipoCajas/edita", EquipoCajaControler.EditaGet)
	app.Post("/EquipoCajas/edita", EquipoCajaControler.EditaPost)
	app.Get("/EquipoCajas/edita/:ID", EquipoCajaControler.EditaGet)
	app.Post("/EquipoCajas/edita/:ID", EquipoCajaControler.EditaPost)

	//Detalle
	app.Get("/EquipoCajas/detalle", EquipoCajaControler.DetalleGet)
	app.Post("/EquipoCajas/detalle", EquipoCajaControler.DetallePost)
	app.Get("/EquipoCajas/detalle/:ID", EquipoCajaControler.DetalleGet)
	app.Post("/EquipoCajas/detalle/:ID", EquipoCajaControler.DetallePost)

	//Rutinas adicionales

	//###################### PuntoVenta ################################
	//Index (Búsqueda)
	app.Get("/PuntoVentas", PuntoVentaControler.IndexGet)
	app.Post("/PuntoVentas", PuntoVentaControler.IndexPost)
	app.Post("/PuntoVentas/search", PuntoVentaControler.BuscaPagina)
	app.Post("/PuntoVentas/agrupa", PuntoVentaControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/PuntoVentas/alta", PuntoVentaControler.AltaGet)
	app.Post("/PuntoVentas/alta", PuntoVentaControler.AltaPost)

	//Edicion
	app.Get("/PuntoVentas/edita", PuntoVentaControler.EditaGet)
	app.Post("/PuntoVentas/edita", PuntoVentaControler.EditaPost)
	app.Get("/PuntoVentas/edita/:ID", PuntoVentaControler.EditaGet)
	app.Post("/PuntoVentas/edita/:ID", PuntoVentaControler.EditaPost)

	//Rutinas adicionales
	app.Post("/PuntoVentas/traeProducto", PuntoVentaControler.TraePrimerDatoDeProducto)
	app.Post("/PuntoVentas/quitarProducto", PuntoVentaControler.QuitarProducto)
	app.Post("/PuntoVentas/modificaCantidad", PuntoVentaControler.ModificarCantidad)
	app.Post("/PuntoVentas/modificaImpuesto", PuntoVentaControler.ModificarImpuesto)
	app.Post("/PuntoVentas/modificaCantidadModal", PuntoVentaControler.ModificarCantidadModal)
	app.Post("/PuntoVentas/modificaCantidadModal2", PuntoVentaControler.ModificarCantidadModal2)
	app.Post("/PuntoVentas/traeOperacion", PuntoVentaControler.TraeOperacion)
	app.Post("/PuntoVentas/imprime", PuntoVentaControler.DetallePost)
	app.Get("/PuntoVentas/cobra/:ID", CajaControler.CobrarDesdeVentaGet)

	//###################### Facturacion ################################
	//Index (Búsqueda)
	app.Get("/Facturacions", FacturacionControler.IndexGet)
	app.Post("/Facturacions", FacturacionControler.IndexPost)
	//Alta
	app.Get("/Facturacions/alta", FacturacionControler.AltaGet)
	app.Post("/Facturacions/alta", FacturacionControler.AltaPost)

	//Edicion
	app.Get("/Facturacions/edita", FacturacionControler.EditaGet)
	app.Post("/Facturacions/edita", FacturacionControler.EditaPost)
	app.Get("/Facturacions/edita/:ID", FacturacionControler.EditaGet)
	app.Post("/Facturacions/edita/:ID", FacturacionControler.EditaPost)

	//Detalle
	app.Get("/Facturacions/detalle", FacturacionControler.DetalleGet)
	app.Post("/Facturacions/detalle", FacturacionControler.DetallePost)
	app.Get("/Facturacions/detalle/:ID", FacturacionControler.DetalleGet)
	app.Post("/Facturacions/detalle/:ID", FacturacionControler.DetallePost)

	//Rutinas adicionales

	//###################### PuntoVenta ################################
	//Index (Búsqueda)
	app.Get("/Compras", CompraControler.IndexGet)

	//Alta
	app.Get("/Compras/alta", CompraControler.AltaGet)
	app.Post("/Compras/alta", CompraControler.AltaPost)

	//###################### Conexion ################################
	//Index (Búsqueda)
	app.Get("/Conexions", ConexionControler.IndexGet)
	app.Post("/Conexions", ConexionControler.IndexPost)
	app.Post("/Conexions/search", ConexionControler.BuscaPagina)
	app.Post("/Conexions/agrupa", ConexionControler.MuestraIndexPorGrupo)
	app.Post("/Conexions/testConexion", ConexionControler.TestConexion)

	//Alta
	app.Get("/Conexions/alta", ConexionControler.AltaGet)
	app.Post("/Conexions/alta", ConexionControler.AltaPost)

	//Edicion
	app.Get("/Conexions/edita", ConexionControler.EditaGet)
	app.Post("/Conexions/edita", ConexionControler.EditaPost)
	app.Get("/Conexions/edita/:ID", ConexionControler.EditaGet)
	app.Post("/Conexions/edita/:ID", ConexionControler.EditaPost)

	//Detalle
	app.Get("/Conexions/detalle", ConexionControler.DetalleGet)
	app.Post("/Conexions/detalle", ConexionControler.DetallePost)
	app.Get("/Conexions/detalle/:ID", ConexionControler.DetalleGet)
	app.Post("/Conexions/detalle/:ID", ConexionControler.DetallePost)

	//Rutinas adicionales

	//###################### PermisosUri ################################
	//Index (Búsqueda)
	app.Get("/PermisosUris", PermisosUriControler.IndexGet)
	app.Post("/PermisosUris", PermisosUriControler.IndexPost)
	app.Post("/PermisosUris/search", PermisosUriControler.BuscaPagina)
	app.Post("/PermisosUris/agrupa", PermisosUriControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/PermisosUris/alta/:ID", PermisosUriControler.AltaGet)
	app.Post("/PermisosUris/alta/:ID", PermisosUriControler.AltaPost)

	//Edicion
	app.Get("/PermisosUris/edita", PermisosUriControler.EditaGet)
	app.Post("/PermisosUris/edita", PermisosUriControler.EditaPost)
	app.Get("/PermisosUris/edita/:ID", PermisosUriControler.EditaGet)
	app.Post("/PermisosUris/edita/:ID", PermisosUriControler.EditaPost)

	//Detalle
	app.Get("/PermisosUris/detalle", PermisosUriControler.DetalleGet)
	app.Post("/PermisosUris/detalle", PermisosUriControler.DetallePost)
	app.Get("/PermisosUris/detalle/:ID", PermisosUriControler.DetalleGet)
	app.Post("/PermisosUris/detalle/:ID", PermisosUriControler.DetallePost)

	//Rutinas adicionales

	//###################### Promocion ################################
	//Index (Búsqueda)
	app.Get("/Promocions", PromocionControler.IndexGet)
	app.Post("/Promocions", PromocionControler.IndexPost)
	app.Post("/Promocions/search", PromocionControler.BuscaPagina)
	app.Post("/Promocions/agrupa", PromocionControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Promocions/alta", PromocionControler.AltaGet)
	app.Post("/Promocions/alta", PromocionControler.AltaPost)

	//Edicion
	app.Get("/Promocions/edita", PromocionControler.EditaGet)
	app.Post("/Promocions/edita", PromocionControler.EditaPost)
	app.Get("/Promocions/edita/:ID", PromocionControler.EditaGet)
	app.Post("/Promocions/edita/:ID", PromocionControler.EditaPost)

	//Detalle
	app.Get("/Promocions/detalle", PromocionControler.DetalleGet)
	app.Post("/Promocions/detalle", PromocionControler.DetallePost)
	app.Get("/Promocions/detalle/:ID", PromocionControler.DetalleGet)
	app.Post("/Promocions/detalle/:ID", PromocionControler.DetallePost)

	//Rutinas adicionales

	//###################### Cotizacion ################################
	//Index (Búsqueda)
	app.Get("/Cotizacions", CotizacionControler.IndexGet)
	app.Post("/Cotizacions", CotizacionControler.IndexPost)
	app.Post("/Cotizacions/search", CotizacionControler.BuscaPagina)
	app.Post("/Cotizacions/agrupa", CotizacionControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Cotizacions/alta", CotizacionControler.AltaGet)
	app.Post("/Cotizacions/alta", CotizacionControler.AltaPost)

	//Edicion
	app.Get("/Cotizacions/edita", CotizacionControler.EditaGet)
	app.Post("/Cotizacions/edita", CotizacionControler.EditaPost)
	app.Get("/Cotizacions/edita/:ID", CotizacionControler.EditaGet)
	app.Post("/Cotizacions/edita/:ID", CotizacionControler.EditaPost)

	//Detalle
	app.Get("/Cotizacions/detalle", CotizacionControler.DetalleGet)
	app.Post("/Cotizacions/detalle", CotizacionControler.DetallePost)
	app.Get("/Cotizacions/detalle/:ID", CotizacionControler.DetalleGet)
	app.Post("/Cotizacions/detalle/:ID", CotizacionControler.DetallePost)

	//Rutinas adicionales
	app.Post("/Cotizacions/traerClientes", CotizacionControler.TraerClientes)
	app.Post("/Cotizacions/traerCliente", CotizacionControler.TraerCliente) //Cliente Especifico
	app.Post("/Cotizacions/TraerProducto", CotizacionControler.TraerProducto)
	app.Post("/Cotizacions/traeProductos", CotizacionControler.TraerProductos)
	app.Post("/Cotizacions/ActualizaProductoCarrito", CotizacionControler.ActualizaProductoCarrito)
	app.Post("/Cotizacions/QuitarProducto", CotizacionControler.QuitarProducto)
	//Rutinas adicionales

	//###################### GrupoAlmacen ################################
	//Index (Búsqueda)
	app.Get("/GrupoAlmacens", GrupoAlmacenControler.IndexGet)
	app.Post("/GrupoAlmacens", GrupoAlmacenControler.IndexPost)
	app.Post("/GrupoAlmacens/search", GrupoAlmacenControler.BuscaPagina)
	app.Post("/GrupoAlmacens/agrupa", GrupoAlmacenControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/GrupoAlmacens/alta", GrupoAlmacenControler.AltaGet)
	app.Post("/GrupoAlmacens/alta", GrupoAlmacenControler.AltaPost)

	//Edicion
	app.Get("/GrupoAlmacens/edita", GrupoAlmacenControler.EditaGet)
	app.Post("/GrupoAlmacens/edita", GrupoAlmacenControler.EditaPost)
	app.Get("/GrupoAlmacens/edita/:ID", GrupoAlmacenControler.EditaGet)
	app.Post("/GrupoAlmacens/edita/:ID", GrupoAlmacenControler.EditaPost)

	//Detalle
	app.Get("/GrupoAlmacens/detalle", GrupoAlmacenControler.DetalleGet)
	app.Post("/GrupoAlmacens/detalle", GrupoAlmacenControler.DetallePost)
	app.Get("/GrupoAlmacens/detalle/:ID", GrupoAlmacenControler.DetalleGet)
	app.Post("/GrupoAlmacens/detalle/:ID", GrupoAlmacenControler.DetallePost)

	//Rutinas adicionales

	//###################### Grupo ################################
	//Index (Búsqueda)
	app.Get("/Grupos", GrupoControler.IndexGet)
	app.Post("/Grupos", GrupoControler.IndexPost)
	app.Post("/Grupos/search", GrupoControler.BuscaPagina)
	app.Post("/Grupos/agrupaB", GrupoControler.MuestraIndexPorGrupoB)
	app.Post("/Grupos/agrupa", GrupoControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Grupos/alta", GrupoControler.AltaGet)
	app.Post("/Grupos/alta", GrupoControler.AltaPost)

	//Edicion
	app.Get("/Grupos/edita", GrupoControler.EditaGet)
	app.Post("/Grupos/edita", GrupoControler.EditaPost)
	app.Get("/Grupos/edita/:ID", GrupoControler.EditaGet)
	app.Post("/Grupos/edita/:ID", GrupoControler.EditaPost)

	//Detalle
	app.Get("/Grupos/detalle", GrupoControler.DetalleGet)
	app.Post("/Grupos/detalle", GrupoControler.DetallePost)
	app.Get("/Grupos/detalle/:ID", GrupoControler.DetalleGet)
	app.Post("/Grupos/detalle/:ID", GrupoControler.DetallePost)

	//Rutinas adicionales
	app.Post("/Grupos/ConsultaBase", GrupoControler.ConsultaElasticBase)

	//###################### Bug ################################
	//Index (Búsqueda)
	app.Get("/Bugs", BugControler.IndexGet)
	app.Post("/Bugs", BugControler.IndexPost)
	app.Post("/Bugs/search", BugControler.BuscaPagina)
	app.Post("/Bugs/agrupa", BugControler.MuestraIndexPorGrupo)

	//Alta
	app.Get("/Bugs/alta", BugControler.AltaGet)
	app.Post("/Bugs/alta", BugControler.AltaPost)

	//Edicion
	app.Get("/Bugs/edita", BugControler.EditaGet)
	app.Post("/Bugs/edita", BugControler.EditaPost)
	app.Get("/Bugs/edita/:ID", BugControler.EditaGet)
	app.Post("/Bugs/edita/:ID", BugControler.EditaPost)

	//Detalle
	app.Get("/Bugs/detalle", BugControler.DetalleGet)
	app.Post("/Bugs/detalle", BugControler.DetallePost)
	app.Get("/Bugs/detalle/:ID", BugControler.DetalleGet)
	app.Post("/Bugs/detalle/:ID", BugControler.DetallePost)

	//###################### Otros #####################################

	//Funciones usadas por ajax (principalmente en el módulo de compras)
	app.Post("/ConsultarProductos", ProductoControler.ConsultarProductosConImpuestosDeAlmacen)
	app.Post("/ConsultaImpuesto", ImpuestoControler.ConsultarImpuestoPorGrupo)
	app.Post("/ConsultarAlmacenesMongo", AlmacenControler.ObtenerIDAlmacenes)
	app.Post("/ConsultarAlmacenesPostgres", AlmacenControler.ConsultarAlmacenesPostgres)
	app.Post("/ConsultarExistenciaInventario", AlmacenControler.ConsultarExistenciaInventario)

	//###################### Configuracion ################################
	//Configurar
	//Get
	app.Get("/Configurar/", ConfiguracionControler.Configuracion)

	//###################### Otros #####################################

	//                        ...

	//###################### Listen Server #############################

	//###################### Problemas de Sistema #############################

	if DataCfg.Puerto != "" {
		fmt.Println("Ejecutandose en el puerto: ", DataCfg.Puerto)
		fmt.Println("Acceder a la siguiente url: ", DataCfg.BaseURL)
		app.Listen(":" + DataCfg.Puerto)
	} else {
		fmt.Println("Ejecutandose en el puerto: 8080")
		fmt.Println("Acceder a la siguiente url: localhost")
		app.Listen(":8080")
	}

}
