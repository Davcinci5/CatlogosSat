package MoGeneral

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//########### GENERALES #######################################

//EstaVacio verifica si un objeto está vacío o no
func EstaVacio(object interface{}) bool {
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	if reflect.ValueOf(object).Kind() == reflect.Struct {
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}
	return false
}

//ConstruirCadenas recibe un texto y regresa dos que se utilizarán para buscar en elastic
func ConstruirCadenas(texto string) (string, string) {

	var palabras = []string{}
	var final = []string{}
	var final2 = []string{}
	var cadenafinal string
	var cadenafinal2 string

	nuevacadena := strings.Replace(texto, "/", "\\/", -1)
	nuevacadena = strings.Replace(nuevacadena, "~", "\\~", -1)
	nuevacadena = strings.Replace(nuevacadena, "^", "\\^", -1)
	nuevacadena = strings.Replace(nuevacadena, "+", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "[", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "]", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "{", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "}", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "(", "\\(", -1)
	nuevacadena = strings.Replace(nuevacadena, ")", "\\)", -1)
	nuevacadena = strings.Replace(nuevacadena, "|", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "=", "", -1)
	nuevacadena = strings.Replace(nuevacadena, ">", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "<", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "!", "", -1)
	nuevacadena = strings.Replace(nuevacadena, "&", "", -1)

	palabras = strings.Split(nuevacadena, " ")

	for _, valor := range palabras {
		if valor != "" {
			palabrita := valor + "~2"
			final = append(final, palabrita)
		}
	}

	for _, valor := range palabras {
		if valor != "" {
			palabrita := `+` + `"` + valor + `"`
			final2 = append(final2, palabrita)
		}
	}

	for _, value := range final {
		cadenafinal = cadenafinal + " " + value
	}

	for _, value := range final2 {
		cadenafinal2 = cadenafinal2 + " " + value
	}

	return cadenafinal, cadenafinal2
}

//Totalpaginas calcula el número de paginaciones de acuerdo al número
// de resultados encontrados y los que se quieren mostrar en la página.
func Totalpaginas(numeroRegistros int, limitePorPagina int) int {
	NumPagina := float32(numeroRegistros) / float32(limitePorPagina)
	NumPagina2 := int(NumPagina)
	if NumPagina > float32(NumPagina2) {
		NumPagina2++
	}
	return NumPagina2
}

//ConstruirPaginacion construtye la paginación en formato html para usarse en la página
func ConstruirPaginacion(paginasTotales int, pag int) string {
	var lt int
	var rt int

	lt = 1
	rt = paginasTotales

	if pag > 2 {
		lt = pag - 1
	}
	if paginasTotales > pag {
		rt = pag + 1
	}

	var templateP string
	templateP += `
	<nav aria-label="Page navigation">
		<ul class="pagination"  style="cursor: pointer;">
			<li>
				<a onclick="BuscaPagina(1)" aria-label="Inicio">
				<span aria-hidden="true">&laquo;</span>
				</a>
			</li>
			<li>
				<a onclick="BuscaPagina(` + strconv.Itoa(lt) + `)" aria-label="Inicio">
				<span aria-hidden="true">&lt;</span>
				</a>
			</li>			
			`

	templateP += ``
	for i := 0; i <= paginasTotales; i++ {
		if i == 1 {
			if i == pag {
				templateP += `<li class="active"><a onclick="BuscaPagina(` + strconv.Itoa(i) + `)">` + strconv.Itoa(i) + `</a></li>`
			} else {
				templateP += `<li><a onclick="BuscaPagina(` + strconv.Itoa(i) + `)">` + strconv.Itoa(i) + `</a></li>`
			}

		} else if i > 1 && i < 11 {
			if i == pag {
				templateP += `<li class="active"><a onclick="BuscaPagina(` + strconv.Itoa(i) + `)">` + strconv.Itoa(i) + `</a></li>`
			} else {
				templateP += `<li><a onclick="BuscaPagina(` + strconv.Itoa(i) + `)">` + strconv.Itoa(i) + `</a></li>`
			}
		} else if i > 11 && i == paginasTotales {

			if i == pag {
				templateP += `<li><span aria-hidden="true">...</span></li><li class="active"><a onclick="BuscaPagina(` + strconv.Itoa(i) + `)">` + strconv.Itoa(i) + `</a></li>`
			} else {
				templateP += `<li><span aria-hidden="true">...</span></li><li><a onclick="BuscaPagina(` + strconv.Itoa(i) + `)">` + strconv.Itoa(i) + `</a></li>`
			}
		}
	}
	templateP += `
		<li>
			<a onclick="BuscaPagina(` + strconv.Itoa(rt) + `)" aria-label="Inicio">
				<span aria-hidden="true">&gt;</span>
			</a>
		</li>			
		<li><a onclick="BuscaPagina(` + strconv.Itoa(paginasTotales) + `)" aria-label="Fin"><span aria-hidden="true">&raquo;</span></a></li></ul></nav>`
	return templateP
}

//MiURI retorna la uri a la cual se hace la peticion, sin parametros
func MiURI(URI, ID string) string {
	arr := strings.Split(URI, "/")
	nuevaURI := ""
	for _, val := range arr {
		if val != ID && val != "" {
			// fmt.Println("[", i, "]", "=", val)
			nuevaURI += "/" + val
		}
	}
	fmt.Println(nuevaURI)
	return nuevaURI
}

//EliminarEspaciosInicioFinal Elimina los espacios en blanco Al inicio y final de una cadena:
//recibe cadena, regresa cadena limpia de espacios al inicio o final o "" si solo contiene espacios
func EliminarEspaciosInicioFinal(cadena string) string {
	var cadenalimpia string
	cadenalimpia = cadena
	re := regexp.MustCompile("(^\\s+|\\s+$)")
	cadenalimpia = re.ReplaceAllString(cadenalimpia, "")
	return cadenalimpia
}
