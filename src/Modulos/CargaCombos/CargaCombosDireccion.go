package CargaCombos

import (
	"../../Modelos/CatalogoModel"
)

//CargaComboxPaises  funcion para cargar los combox de paises
func CargaComboxPaises(ID string) string {

	Paises := CatalogoModel.GetAllPaises()
	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, Pais := range Paises {
		if ID == Pais.ID.Hex() {
			templ += `<option value="` + Pais.ID.Hex() + `" selected>` + Pais.Nombre + `</option>`
		} else {
			templ += `<option value="` + Pais.ID.Hex() + `">` + Pais.Nombre + `</option>`
		}
	}

	return templ
}

//CargaComboEstados funccion que carga los combos de Estado
func CargaComboEstados(ID string) string {
	Estados := CatalogoModel.GetAllEstados()

	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, Estado := range Estados {
		if ID == Estado.ID.Hex() {
			templ += `<option value="` + Estado.ID.Hex() + `" selected>` + Estado.Nombre + `</option>`
		} else {
			templ += `<option value="` + Estado.ID.Hex() + `">` + Estado.Nombre + `</option>`
		}
	}

	return templ
}

//CargaComboMunicipiosForClaveEstado  funcion que carga  los municipios por clave del estado
func CargaComboMunicipiosForClaveEstado(ClaveEdo, ID string) string {

	Municipios := CatalogoModel.GetAllMunicipiosForClaveEstado(ClaveEdo)
	templ := ``

	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}

	for _, Municipio := range Municipios {
		if ID == Municipio.ID.Hex() {
			templ += `<option value="` + Municipio.ID.Hex() + `" selected>` + Municipio.Nombre + `</option>`
		} else {
			templ += `<option value="` + Municipio.ID.Hex() + `">` + Municipio.Nombre + `</option>`
		}
	}

	return templ
}

//CargaComboColoniasForClaveMunicipio funcion que carga las colonias desde la clave del municipio
func CargaComboColoniasForClaveMunicipio(ClaveMpo string, ID string) string {

	Coloniass := CatalogoModel.GetAllColoniasForClaveMunicipio(ClaveMpo)
	templ := ``
	if ID != "" {
		templ = `<option value="">--SELECCIONE--</option>`
	} else {
		templ = `<option value="" selected>--SELECCIONE--</option>`
	}
	for _, Colonias := range Coloniass {
		if ID == Colonias.ID.Hex() {
			templ += `<option value="` + Colonias.ID.Hex() + `" selected>` + Colonias.Nombre + `</option>`
		} else {
			templ += `<option value="` + Colonias.ID.Hex() + `">` + Colonias.Nombre + `</option>`
		}

	}
	return templ
}
