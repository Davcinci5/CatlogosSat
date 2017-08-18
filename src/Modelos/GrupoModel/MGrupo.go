package GrupoModel

import (
	"fmt"
	"strconv"
	"time"

	"github.com/leekchan/accounting"

	"../../Modelos/AlmacenModel"
	"../../Modelos/CajaModel"
	"../../Modelos/CatalogoModel"
	"../../Modelos/ClienteModel"
	"../../Modelos/DispositivoModel"
	"../../Modelos/PersonaModel"
	"../../Modelos/ProductoModel"
	"../../Modelos/UnidadModel"
	"../../Modelos/UsuarioModel"

	"../../Modulos/Conexiones"
	"../../Modulos/General"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/olivere/elastic.v5"
)

//#########################< ESTRUCTURAS >##############################

//GrupoMgo estructura de Grupos mongo
type GrupoMgo struct {
	ID            bson.ObjectId   `bson:"_id,omitempty"`
	Nombre        string          `bson:"Nombre"`
	Descripcion   string          `bson:"Descripcion"`
	PermiteVender bool            `bson:"PermiteVender,omitempty"`
	Tipo          bson.ObjectId   `bson:"Tipo,omitempty"`
	Miembros      []bson.ObjectId `bson:"Miembros,omitempty"`
	Estatus       bson.ObjectId   `bson:"Estatus"`
	FechaHora     time.Time       `bson:"FechaHora"`
	FechaEdicion  []time.Time     `bson:"FechaEdicion"`
}

//GrupoElastic estructura de Grupos para insertar en Elastic
type GrupoElastic struct {
	Nombre        string    `json:"Nombre"`
	Descripcion   string    `json:"Descripcion"`
	PermiteVender string    `json:"PermiteVender,omitempty"`
	Tipo          string    `json:"Tipo,omitempty"`
	Miembros      []string  `json:"Miembros,omitempty"`
	Estatus       string    `json:"Estatus"`
	FechaHora     time.Time `json:"FechaHora"`
}

//#########################< FUNCIONES GENERALES MGO >###############################

//GetAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func GetAll() []GrupoMgo {
	var result []GrupoMgo
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}
	err = Grupos.Find(nil).All(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//CountAll Regresa todos los documentos existentes de Mongo (Por Coleccion)
func CountAll() int {
	var result int
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)

	if err != nil {
		fmt.Println(err)
	}
	result, err = Grupos.Find(nil).Count()
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetOne Regresa un documento específico de Mongo (Por Coleccion)
func GetOne(ID bson.ObjectId) GrupoMgo {
	var result GrupoMgo
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}
	err = Grupos.Find(bson.M{"_id": ID}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetEspecifics rsegresa un conjunto de documentos específicos de Mongo (Por Coleccion)
func GetEspecifics(Ides []bson.ObjectId) []GrupoMgo {
	var result []GrupoMgo
	var aux GrupoMgo
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range Ides {
		aux = GrupoMgo{}
		Grupos.Find(bson.M{"_id": value}).One(&aux)
		result = append(result, aux)
	}
	s.Close()
	return result
}

//GetEspecificByFields regresa un documento de Mongo especificando un campo y un determinado valor
func GetEspecificByFields(field string, valor interface{}) GrupoMgo {
	var result GrupoMgo
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)

	if err != nil {
		fmt.Println(err)
	}
	err = Grupos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result
}

//GetIDByField regresa un documento específico de Mongo (Por Coleccion)
func GetIDByField(field string, valor interface{}) bson.ObjectId {
	var result GrupoMgo
	s, Grupos, err := MoConexion.GetColectionMgo(MoVar.ColeccionGrupo)
	if err != nil {
		fmt.Println(err)
	}
	err = Grupos.Find(bson.M{field: valor}).One(&result)
	if err != nil {
		fmt.Println(err)
	}
	s.Close()
	return result.ID
}

//CargaComboGrupos regresa un combo de Grupo de mongo
func CargaComboGrupos(ID string) string {
	Grupos := GetAll()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option> `
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option> `
	}

	for _, v := range Grupos {
		if ID == v.ID.Hex() {
			templ += `<option value="` + v.ID.Hex() + `" selected>  ` + v.Nombre + ` </option> `
		} else {
			templ += `<option value="` + v.ID.Hex() + `">  ` + v.Nombre + ` </option> `
		}

	}
	return templ
}

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Grupos []GrupoMgo) (string, string) {
	//floats := accounting.Accounting{Symbol: "", Precision: 2}
	cuerpo := ``

	cabecera := `<tr>
			    <th>#</th>
				<th>Nombre</th>					
				<th>Descripcion</th>					
				<th>Tipo de Objeto</th>					
				<th>Estatus</th>					
				<th>Fecha De Creación</th>					
				</tr>`

	for k, v := range Grupos {
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Grupos/detalle/` + v.ID.Hex() + `';">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + v.Nombre + `</td>`
		cuerpo += `<td>` + v.Descripcion + `</td>`
		cuerpo += `<td>` + CatalogoModel.ObtenerValoresCatalogoPorValor(v.Tipo) + `</td>`
		cuerpo += `<td>` + CatalogoModel.ObtenerValoresCatalogoPorValor(v.Estatus) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//########################< FUNCIONES GENERALES PSQL >#############################

//######################< FUNCIONES GENERALES ELASTIC >############################

//BuscarEnElastic busca el texto solicitado en los campos solicitados
func BuscarEnElastic(texto string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	queryTilde = queryTilde.Field("Nombre")
	queryQuotes = queryQuotes.Field("Nombre")

	queryTilde = queryTilde.Field("Descripcion")
	queryQuotes = queryQuotes.Field("Descripcion")

	queryTilde = queryTilde.Field("PermiteVender")
	queryQuotes = queryQuotes.Field("PermiteVender")

	queryTilde = queryTilde.Field("Tipo")
	queryQuotes = queryQuotes.Field("Tipo")

	queryTilde = queryTilde.Field("Miembros")
	queryQuotes = queryQuotes.Field("Miembros")

	queryTilde = queryTilde.Field("Estatus")
	queryQuotes = queryQuotes.Field("Estatus")

	var docs *elastic.SearchResult
	var err bool

	docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupo, queryTilde)
	if err {
		fmt.Println("No Match 1st Try")
	}

	if docs.Hits.TotalHits == 0 {
		docs, err = MoConexion.BuscaElastic(MoVar.TipoGrupo, queryQuotes)
		if err {
			fmt.Println("No Match 2nd Try")
		}
	}

	return docs
}

//BuscaObjetosEnElastic busca el texto solicitado en los campos solicitados
func BuscaObjetosEnElastic(texto, tipo string) *elastic.SearchResult {
	textoTilde, textoQuotes := MoGeneral.ConstruirCadenas(texto)

	queryTilde := elastic.NewQueryStringQuery(textoTilde)
	queryQuotes := elastic.NewQueryStringQuery(textoQuotes)

	var docs *elastic.SearchResult
	var err error

	switch tipo {

	case "CLIENTES":
		queryTilde = queryTilde.Field("RFC")
		queryQuotes = queryQuotes.Field("RFC")

		queryTilde = queryTilde.Field("Estatus")
		queryQuotes = queryQuotes.Field("Estatus")

		queryTilde = queryTilde.Field("Otros")
		queryQuotes = queryQuotes.Field("Otros")

		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoCliente, queryTilde)
		if err != nil {
			fmt.Println(err)
			fmt.Println("No Match 1st Try")
		}

		if docs.Hits.TotalHits == 0 {
			docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoCliente, queryQuotes)
			if err != nil {
				fmt.Println(err)
				fmt.Println("No Match 2nd Try")
			}
		}
		break
	case "PERSONAS":

		queryTilde = queryTilde.Field("Nombre")
		queryQuotes = queryQuotes.Field("Nombre")

		queryTilde = queryTilde.Field("Tipo")
		queryQuotes = queryQuotes.Field("Tipo")

		queryTilde = queryTilde.Field("Grupos")
		queryQuotes = queryQuotes.Field("Grupos")

		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoPersona, queryTilde)
		if err != nil {
			fmt.Println(err)
			fmt.Println("No Match 1st Try")
		}

		if docs.Hits.TotalHits == 0 {
			docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoPersona, queryQuotes)
			if err != nil {
				fmt.Println(err)
				fmt.Println("No Match 2nd Try")
			}
		}
		break
	case "USUARIOS":

		queryTilde = queryTilde.Field("IDPersona")
		queryQuotes = queryQuotes.Field("IDPersona")

		queryTilde = queryTilde.Field("Usuario")
		queryQuotes = queryQuotes.Field("Usuario")

		queryTilde = queryTilde.Field("Correos")
		queryQuotes = queryQuotes.Field("Correos")

		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoUsuario, queryTilde)
		if err != nil {
			fmt.Println(err)
			fmt.Println("No Match 1st Try")
		}

		if docs.Hits.TotalHits == 0 {
			docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoUsuario, queryQuotes)
			if err != nil {
				fmt.Println(err)
				fmt.Println("No Match 2nd Try")
			}
		}
		break
	case "PRODUCTOS":

		queryTilde = queryTilde.Field("Nombre")
		queryQuotes = queryQuotes.Field("Nombre")

		queryTilde = queryTilde.Field("Codigos.Valores")
		queryQuotes = queryQuotes.Field("Codigos.Valores")

		queryTilde = queryTilde.Field("Etiquetas")
		queryQuotes = queryQuotes.Field("Etiquetas")

		queryTilde = queryTilde.Field("Estatus")
		queryQuotes = queryQuotes.Field("Estatus")

		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoProducto, queryTilde)
		if err != nil {
			fmt.Println(err)
			fmt.Println("No Match 1st Try")
		}
		if docs.Hits.TotalHits == 0 {
			docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoProducto, queryQuotes)
			if err != nil {
				fmt.Println(err)
				fmt.Println("No Match 2nd Try")
			}
		}
		break
	case "ALMACENES":

		queryTilde = queryTilde.Field("Nombre")
		queryQuotes = queryQuotes.Field("Nombre")

		queryTilde = queryTilde.Field("Tipo")
		queryQuotes = queryQuotes.Field("Tipo")

		queryTilde = queryTilde.Field("Clasificacion")
		queryQuotes = queryQuotes.Field("Clasificacion")

		queryTilde = queryTilde.Field("Estatus")
		queryQuotes = queryQuotes.Field("Estatus")

		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoAlmacen, queryTilde)
		if err != nil {
			fmt.Println(err)
			fmt.Println("No Match 1st Try")
		}

		if docs.Hits.TotalHits == 0 {
			docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoAlmacen, queryQuotes)
			if err != nil {
				fmt.Println(err)
				fmt.Println("No Match 2nd Try")
			}
		}
		break
	// case "PROVEEDORES":
	// 	docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoProveedor, queryTilde)
	// 	if err != nil {
	// 		fmt.Println("No Match 1st Try")
	// 	}

	// 	if docs.Hits.TotalHits == 0 {
	// 		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoProveedor, queryQuotes)
	// 		if err != nil {
	// 			fmt.Println("No Match 2nd Try")
	// 		}
	// 	}
	// break
	case "CAJAS":

		queryTilde = queryTilde.Field("Usuario")
		queryQuotes = queryQuotes.Field("Usuario")

		queryTilde = queryTilde.Field("Caja")
		queryQuotes = queryQuotes.Field("Caja")

		queryTilde = queryTilde.Field("Operacion")
		queryQuotes = queryQuotes.Field("Operacion")

		queryTilde = queryTilde.Field("Estatus")
		queryQuotes = queryQuotes.Field("Estatus")

		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoCaja, queryTilde)
		if err != nil {
			fmt.Println(err)
			fmt.Println("No Match 1st Try")
		}

		if docs.Hits.TotalHits == 0 {
			docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoCaja, queryQuotes)
			if err != nil {
				fmt.Println(err)
				fmt.Println("No Match 2nd Try")
			}
		}
		break
	case "DISPOSITIVOS":

		queryTilde = queryTilde.Field("Nombre")
		queryQuotes = queryQuotes.Field("Nombre")

		queryTilde = queryTilde.Field("Mac")
		queryQuotes = queryQuotes.Field("Mac")
		docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoDispositivo, queryTilde)
		if err != nil {
			fmt.Println(err)
			fmt.Println("No Match 1st Try")
		}

		if docs.Hits.TotalHits == 0 {
			docs, err = MoConexion.BuscaElasticAvanzada(MoVar.TipoDispositivo, queryQuotes)
			if err != nil {
				fmt.Println(err)
				fmt.Println("No Match 2nd Try")
			}
		}
		break
	default:
		fmt.Println("No se encontró el tipo de objeto.")
		break
	}
	return docs
}

//GeneraTemplatesBusquedaObjetos regresa cabecera y cuerpo de consulta de objetos específica
func GeneraTemplatesBusquedaObjetos(Objetos []bson.ObjectId, Tipo string, All bool) (string, string, int) {
	var Cabecera, Cuerpo string
	var Contador int

	switch Tipo {
	case "CLIENTES":
		Clientes := ClienteModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>
				<th class="text-center">`

		if All {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" checked>`
		} else {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">`
		}

		Cabecera += `</th>
				 	<th>Nombre</th>
					<th>RFC</th>								
					<th>Estatus</th>					
				</tr>`

		for k, v := range Clientes {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Persona := PersonaModel.GetOne(v.IDPersona)
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`

				if All {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)" checked></td>`
				} else {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
				}

				Cuerpo += `<td>` + Persona.Nombre + `</td>`
				Cuerpo += `<td>` + v.RFC + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 137) + `</td>`
				Cuerpo += `</tr>`
			}

		}
		break
	case "PERSONAS":
		Personas := PersonaModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
						<th>#</th>
						<th class="text-center">`

		if All {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" checked>`
		} else {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">`
		}
		Cabecera += `</th>
						<th>Nombre</th>									
						<th>Tipo</th>									
						<th>Grupos</th>									
						<th>Predecesor</th>									
						<th>Estatus</th>									
						<th>FechaHora</th>					
					</tr>`

		for k, v := range Personas {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				if All {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)" checked></td>`
				} else {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
				}
				Cuerpo += `<td>` + v.Nombre + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo[0], 159) + `</td>`
				Cuerpo += `<td>` + "v.Grupos" + `</td>`
				Cuerpo += `<td>` + v.Predecesor.Hex() + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 160) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "USUARIOS":
		Usuarios := UsuarioModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>
						<th class="text-center">`

		if All {
			Cabecera += `(Todos)<input name="Objeto" class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" checked>`
		} else {
			Cabecera += `(Todos)<input name="Objeto" class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">`
		}
		Cabecera += `</th>
						<th>Usuario</th>
				<th>Nombre</th>
				<th>Sexo</th>								
				<th>Estatus</th>					
				</tr>`

		for k, v := range Usuarios {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Persona := PersonaModel.GetOne(v.IDPersona)
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				if All {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)" checked></td>`
				} else {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
				}
				Cuerpo += `<td>` + v.Usuario + `</td>`
				Cuerpo += `<td>` + Persona.Nombre + `</td>`
				Cuerpo += `<td>` + Persona.Sexo + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 137) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "PRODUCTOS":
		etiquetas := ``
		codigos := ``
		Productos := ProductoModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
					<th>#</th>
						<th class="text-center">`

		if All {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" checked>`
		} else {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">`
		}
		Cabecera += `</th>
						<th>Nombre</th>									
					<th>Códigos</th>									
					<th>Tipo</th>									
					<th>Unidad</th>									
					<th>Venta por Fracciones</th>									
					<th>Etiquetas</th>									
					<th>Estatus</th>									
					<th>FechaHora</th>					
				</tr>`

		for k, v := range Productos {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				etiquetas = ``
				codigos = ``
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				if All {
					Cuerpo += `<td class="text-center"><input  name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)" checked></td>`
				} else {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
				}
				Cuerpo += `<td>` + v.Nombre + `</td>`

				for i, cla := range v.Codigos.Claves {
					codigos += cla + `:` + v.Codigos.Valores[i] + `,`
				}

				codigos = codigos[:len(codigos)-1]
				Cuerpo += `<td>` + codigos + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo, 162) + `</td>`
				Cuerpo += `<td>` + UnidadModel.GetNombreUnidadByField("Datos._id", v.Unidad) + `</td>`
				if v.VentaFraccion {
					Cuerpo += `<td>SI</td>`
				} else {
					Cuerpo += `<td>NO</td>`
				}

				for _, val := range v.Etiquetas {
					etiquetas += val + `,`
				}
				etiquetas = etiquetas[:len(etiquetas)-1]
				Cuerpo += `<td>` + etiquetas + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 161) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "ALMACENES":
		Almacens := AlmacenModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>
						<th class="text-center">`

		if All {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" checked>`
		} else {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">`
		}
		Cabecera += `</th>
						<th>Nombre</th>									
				<th>Tipo</th>									
				<th>Clasificacion</th>									
				<th>Predecesor</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

		for k, v := range Almacens {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				if All {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)" checked></td>`
				} else {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
				}
				Cuerpo += `<td>` + v.Nombre + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo, 132) + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Clasificacion, 133) + `</td>`
				Cuerpo += `<td>` + v.Predecesor.Hex() + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 134) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		// case "PROVEEDORES":
		break
	case "CAJAS":
		floats := accounting.Accounting{Symbol: "$", Precision: 2}
		Cajas := CajaModel.GetEspecifics(Objetos)

		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>
						<th class="text-center">`

		if All {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" checked>`
		} else {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">`
		}
		Cabecera += `</th>
						<th>Usuario</th>									
				<th>Caja</th>									
				<th>Cargo</th>									
				<th>Abono</th>									
				<th>Saldo</th>									
				<th>Operacion</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

		for k, v := range Cajas {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `" >`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				if All {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)" checked></td>`
				} else {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
				}
				Cuerpo += `<td>` + v.Usuario.Hex() + `</td>`
				Cuerpo += `<td>` + v.Caja.Hex() + `</td>`
				Cuerpo += `<td>` + floats.FormatMoney(v.Cargo) + `</td>`
				Cuerpo += `<td>` + floats.FormatMoney(v.Abono) + `</td>`
				Cuerpo += `<td>` + floats.FormatMoney(v.Saldo) + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Operacion, 169) + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 135) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "DISPOSITIVOS":
		Dispositivos := DispositivoModel.GetEspecifics(Objetos)

		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>
						<th class="text-center">`

		if All {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)" checked>`
		} else {
			Cabecera += `(Todos)<input class="ProdAsignar"  type="checkbox" value="todos" onClick="SelecionarTodos(this.checked)">`
		}
		Cabecera += `</th>
						<th>Nombre</th>									
				<th>Descripcion</th>									
				<th>Mac</th>																													
				</tr>`

		for k, v := range Dispositivos {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `" >`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				if All {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)" checked></td>`
				} else {
					Cuerpo += `<td class="text-center"><input name="Objeto" class="ProdExtraido" type="checkbox" id="` + v.ID.Hex() + `" value="` + v.ID.Hex() + `" onClick="GuardaSeleccionados(this.checked, this.id)"></td>`
				}
				Cuerpo += `<td>` + v.Nombre + `</td>`
				Cuerpo += `<td>` + v.Descripcion + `</td>`
				Cuerpo += `<td>` + v.Mac + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	default:
		fmt.Println("No se encontró el tipo de Objeto.")
		break
	}

	return Cabecera, Cuerpo, Contador
}

//GeneraTemplatesBusquedaObjetosDetalle regresa cabecera y cuerpo de consulta de objetos específica
func GeneraTemplatesBusquedaObjetosDetalle(Objetos []bson.ObjectId, Tipo string) (string, string, int) {
	var Cabecera, Cuerpo string
	var Contador int

	switch Tipo {
	case "CLIENTES":
		Clientes := ClienteModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>
				<th>Nombre</th>
				<th>RFC</th>								
				<th>Estatus</th>					
				</tr>`

		for k, v := range Clientes {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Persona := PersonaModel.GetOne(v.IDPersona)
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				Cuerpo += `<td>` + Persona.Nombre + `</td>`
				Cuerpo += `<td>` + v.RFC + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 137) + `</td>`
				Cuerpo += `</tr>`
			}

		}
		break
	case "PERSONAS":
		Personas := PersonaModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
						<th>#</th>		
						<th>Nombre</th>									
						<th>Tipo</th>									
						<th>Grupos</th>									
						<th>Predecesor</th>									
						<th>Estatus</th>									
						<th>FechaHora</th>					
					</tr>`

		for k, v := range Personas {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				Cuerpo += `<td>` + v.Nombre + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo[0], 159) + `</td>`
				Cuerpo += `<td>` + "v.Grupos" + `</td>`
				Cuerpo += `<td>` + v.Predecesor.Hex() + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 160) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "USUARIOS":
		Usuarios := UsuarioModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>
				<th>Usuario</th>
				<th>Nombre</th>
				<th>Sexo</th>								
				<th>Estatus</th>					
				</tr>`

		for k, v := range Usuarios {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Persona := PersonaModel.GetOne(v.IDPersona)
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				Cuerpo += `<td>` + v.Usuario + `</td>`
				Cuerpo += `<td>` + Persona.Nombre + `</td>`
				Cuerpo += `<td>` + Persona.Sexo + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 137) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "PRODUCTOS":
		etiquetas := ``
		codigos := ``
		Productos := ProductoModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
					<th>#</th>		
					<th>Nombre</th>									
					<th>Códigos</th>									
					<th>Tipo</th>									
					<th>Unidad</th>									
					<th>Venta por Fracciones</th>									
					<th>Etiquetas</th>									
					<th>Estatus</th>									
					<th>FechaHora</th>					
				</tr>`

		for k, v := range Productos {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				etiquetas = ``
				codigos = ``
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				Cuerpo += `<td>` + v.Nombre + `</td>`

				for i, cla := range v.Codigos.Claves {
					codigos += cla + `:` + v.Codigos.Valores[i] + `,`
				}

				codigos = codigos[:len(codigos)-1]
				Cuerpo += `<td>` + codigos + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo, 162) + `</td>`
				Cuerpo += `<td>` + UnidadModel.GetNombreUnidadByField("Datos._id", v.Unidad) + `</td>`
				if v.VentaFraccion {
					Cuerpo += `<td>SI</td>`
				} else {
					Cuerpo += `<td>NO</td>`
				}

				for _, val := range v.Etiquetas {
					etiquetas += val + `,`
				}
				etiquetas = etiquetas[:len(etiquetas)-1]
				Cuerpo += `<td>` + etiquetas + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 161) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "ALMACENES":
		Almacens := AlmacenModel.GetEspecifics(Objetos)
		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>		
				<th>Nombre</th>									
				<th>Tipo</th>									
				<th>Clasificacion</th>									
				<th>Predecesor</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

		for k, v := range Almacens {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `">`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				Cuerpo += `<td>` + v.Nombre + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Tipo, 132) + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Clasificacion, 133) + `</td>`
				Cuerpo += `<td>` + v.Predecesor.Hex() + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 134) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		// case "PROVEEDORES":
		break
	case "CAJAS":
		floats := accounting.Accounting{Symbol: "$", Precision: 2}
		Cajas := CajaModel.GetEspecifics(Objetos)

		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>							
				<th>Usuario</th>									
				<th>Caja</th>									
				<th>Cargo</th>									
				<th>Abono</th>									
				<th>Saldo</th>									
				<th>Operacion</th>									
				<th>Estatus</th>									
				<th>FechaHora</th>					
				</tr>`

		for k, v := range Cajas {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `" >`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				Cuerpo += `<td>` + v.Usuario.Hex() + `</td>`
				Cuerpo += `<td>` + v.Caja.Hex() + `</td>`
				Cuerpo += `<td>` + floats.FormatMoney(v.Cargo) + `</td>`
				Cuerpo += `<td>` + floats.FormatMoney(v.Abono) + `</td>`
				Cuerpo += `<td>` + floats.FormatMoney(v.Saldo) + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Operacion, 169) + `</td>`
				Cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 135) + `</td>`
				Cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	case "DISPOSITIVOS":
		Dispositivos := DispositivoModel.GetEspecifics(Objetos)

		Cuerpo = ``
		Cabecera = `<tr>
				<th>#</th>						
				<th>Nombre</th>									
				<th>Descripcion</th>									
				<th>Mac</th>																													
				</tr>`

		for k, v := range Dispositivos {
			if !MoGeneral.EstaVacio(v) {
				Contador++
				Cuerpo += `<tr id = "` + v.ID.Hex() + `" >`
				Cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
				Cuerpo += `<td>` + v.Nombre + `</td>`
				Cuerpo += `<td>` + v.Descripcion + `</td>`
				Cuerpo += `<td>` + v.Mac + `</td>`
				Cuerpo += `</tr>`
			}
		}
		break
	default:
		fmt.Println("No se encontró el tipo de objeto.")
		break
	}

	return Cabecera, Cuerpo, Contador
}
