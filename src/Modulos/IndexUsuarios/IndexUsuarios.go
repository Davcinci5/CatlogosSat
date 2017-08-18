package IndexUsuarios

import (
	"strconv"
	"time"

	"../../Modelos/CatalogoModel"
	"../../Modelos/GrupoPersonaModel"
	"../../Modelos/PersonaModel"
	"../../Modelos/UsuarioModel"
)

//GeneraTemplatesBusqueda crea templates de tabla de búsqueda
func GeneraTemplatesBusqueda(Usuarios []UsuarioModel.UsuarioMgo) (string, string) {
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>			
					<th>Nombre</th>									
					<th>Tipo</th>
					<th>Grupo</th>
					<th>Usuario</th>				
					<th>Estatus</th>									
					<th>FechaHora</th>					
				</tr>`

	for k, v := range Usuarios {
		persona := PersonaModel.GetOne(v.IDPersona)
		cuerpo += `<tr id = "` + v.ID.Hex() + `" onclick="window.location.href = '/Usuarios/detalle/` + v.ID.Hex() + `';">`
		//cuerpo += `<tr id = "` + v.ID.Hex() + `">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + persona.Nombre + `</td>`
		var tipoPersona string
		for _, value := range persona.Tipo {
			tipoPersona += CatalogoModel.RegresaNombreSubCatalogo(value)
		}
		cuerpo += `<td>` + tipoPersona + `</td>`
		var grupoPersona string
		for _, value := range persona.Grupos {
			grupoPersona += GrupoPersonaModel.CargaNombreGrupo(value)
		}
		cuerpo += `<td>` + grupoPersona + `</td>`
		cuerpo += `<td>` + v.Usuario + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 167) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}

//GeneraTemplatesBusquedaSesiones crea templates de tabla de búsqueda
func GeneraTemplatesBusquedaSesiones(Usuarios []UsuarioModel.UsuarioMgo) (string, string) {
	cuerpo := ``

	cabecera := `<tr>
					<th>#</th>			
					<th>Nombre</th>									
					<th>Tipo</th>
					<th>Grupo</th>
					<th>Usuario</th>				
					<th>Estatus</th>									
					<th>FechaHora</th>					
				</tr>`

	for k, v := range Usuarios {
		persona := PersonaModel.GetOne(v.IDPersona)
		cuerpo += `<tr id = "` + v.Usuario + `" onclick="window.location.href = '/Sesiones/detalle/` + v.Usuario + `';">`
		//cuerpo += `<tr id = "` + v.ID.Hex() + `">`
		cuerpo += `<td>` + strconv.Itoa(k+1) + `</td>`
		cuerpo += `<td>` + persona.Nombre + `</td>`
		var tipoPersona string
		for _, value := range persona.Tipo {
			tipoPersona += CatalogoModel.RegresaNombreSubCatalogo(value)
		}
		cuerpo += `<td>` + tipoPersona + `</td>`
		var grupoPersona string
		for _, value := range persona.Grupos {
			grupoPersona += GrupoPersonaModel.CargaNombreGrupo(value)
		}
		cuerpo += `<td>` + grupoPersona + `</td>`
		cuerpo += `<td>` + v.Usuario + `</td>`
		cuerpo += `<td>` + CatalogoModel.GetValorMagnitud(v.Estatus, 167) + `</td>`
		cuerpo += `<td>` + v.FechaHora.Format(time.RFC1123) + `</td>`
		cuerpo += `</tr>`
	}

	return cabecera, cuerpo
}
