package ListaCostoControler

import (
	"strconv"
	"html/template"

	"../../Modulos/Conexiones"
	"../../Modelos/ListaCostoModel"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/mgo.v2/bson"
)
		

//##########< Variables Generales > ############

var cadenaBusqueda string
var buscarEn string
var numeroRegistros int64
var paginasTotales int
//NumPagina especifica ***************
var NumPagina float32
//limitePorPagina especifica ***************
var limitePorPagina = 10
var result []ListaCostoModel.ListaCosto
var resultPage []ListaCostoModel.ListaCosto
var templatePaginacion = ``
		

//####################< INDEX (BUSQUEDA) >###########################

//IndexGet renderea al index de ListaCosto
func IndexGet(ctx *iris.Context) {
ctx.Render("ListaCostoIndex.html", nil)
}
//IndexPost regresa la peticon post que se hizo desde el index de ListaCosto
func IndexPost(ctx *iris.Context) {

		templatePaginacion = ``
		
		var resultados []ListaCostoModel.ListaCostoMgo
		var IDToObjID bson.ObjectId		
		var arrObjIds []bson.ObjectId
		var arrToMongo []bson.ObjectId

		cadenaBusqueda = ctx.FormValue("searchbox")
		buscarEn = ctx.FormValue("buscaren")

		if cadenaBusqueda != "" {

			docs := ListaCostoModel.BuscarEnElastic(cadenaBusqueda)

			if docs.Hits.TotalHits > 0 {
				numeroRegistros = docs.Hits.TotalHits

				paginasTotales = Totalpaginas()

				for _, item := range docs.Hits.Hits {
					IDToObjID = bson.ObjectIdHex(item.Id)
					arrObjIds = append(arrObjIds, IDToObjID)
				}

				if numeroRegistros <= int64(limitePorPagina) {
					for _, v := range arrObjIds[0:numeroRegistros] {
						arrToMongo = append(arrToMongo, v)
					}
				} else if numeroRegistros >= int64(limitePorPagina) {
					for _, v := range arrObjIds[0:limitePorPagina] {
						arrToMongo = append(arrToMongo, v)
					}
				}

				resultados = ListaCostoModel.GetEspecifics(arrToMongo)

				MoConexion.FlushElastic()

			}

		}

		templatePaginacion = ConstruirPaginacion()
		
		ctx.Render("ListaCostoIndex.html", map[string]interface{}{
			"result":          resultados,
			"cadena_busqueda": cadenaBusqueda,
			"PaginacionT":     template.HTML(templatePaginacion),
		})
				
		
}
//###########################< ALTA >################################

//AltaGet renderea al alta de ListaCosto
func AltaGet(ctx *iris.Context) {
ctx.Render("ListaCostoAlta.html", nil)
}
//AltaPost regresa la petición post que se hizo desde el alta de ListaCosto
func AltaPost(ctx *iris.Context) {

	//######### LEE TU OBJETO DEL FORMULARIO #########
	var ListaCosto ListaCostoModel.ListaCostoMgo
	ctx.ReadForm(&ListaCosto)
					
	//######### VALIDA TU OBJETO #########
	EstatusPeticion := true //True indica que hay un error
	//##### TERMINA TU VALIDACION ########
	
	//########## Asigna vairables a la estructura que enviarás a la vista
	ListaCosto.ID = bson.NewObjectId()
	
		
	//######### ENVIA TUS RESULTADOS #########
		var SListaCosto ListaCostoModel.SListaCosto

	//	SListaCosto.ListaCosto = ListaCosto //Asigamos el Objeto que hemos capturado para que pueda regresar los valores capturados a la vista.
				
	if EstatusPeticion{
		SListaCosto.SEstado = false //En la vista los errores se manejan al reves para hacer uso del rellenado por defecto de Go
		SListaCosto.SMsj = "La validación indica que el objeto capturado no puede procesarse" //La idea es después hacer un colector de errores y mensaje de éxito y enviarlo en esta variable.
		ctx.Render("ListaCostoAlta.html", SListaCosto)
	}else{
			
		//Si no hubo error se procede a realizar alguna acción con el objeto, en este caso, una inserción.
		if ListaCosto.InsertaMgo(){
			SListaCosto.SEstado = true 
			SListaCosto.SMsj = "Se ha realizado una inserción exitosa" 

			//SE PUEDE TOMA LA DECICIÓN QUE SE CREA MÁS PERTINENTE, EN ESTE CASO SE CONSIDERA EL DETALLE DEL OBJETO.
			ctx.Render("ListaCostoDetalle.html", SListaCosto)
						
		}else{
			SListaCosto.SEstado = false
			SListaCosto.SMsj = "Ocurrió un error al insertar el Objeto, intente más tarde" 
			ctx.Render("ListaCostoAlta.html", SListaCosto)
		}
	
}
		
}
//###########################< EDICION >###############################

//EditaGet renderea a la edición de ListaCosto
func EditaGet(ctx *iris.Context) {
ctx.Render("ListaCostoEdita.html", nil)
}
//EditaPost regresa el resultado de la petición post generada desde la edición de ListaCosto
func EditaPost(ctx *iris.Context) {
ctx.Render("ListaCostoEdita.html", nil)
}

//#################< DETALLE >####################################

//DetalleGet renderea al index.html
func DetalleGet(ctx *iris.Context) {
ctx.Render("ListaCostoDetalle.html", nil)
}

//DetallePost renderea al index.html
func DetallePost(ctx *iris.Context) {
ctx.Render("ListaCostoDetalle.html", nil)
}
//####################< RUTINAS ADICIONALES >##########################


//Totalpaginas calcula el número de paginaciones de acuerdo al número
// de resultados encontrados y los que se quieren mostrar en la página.
func Totalpaginas() int {

	NumPagina = float32(numeroRegistros) / float32(limitePorPagina)
	NumPagina2 := int(NumPagina)
	if NumPagina > float32(NumPagina2) {
		NumPagina2++
	}
	totalpaginas := NumPagina2
	return totalpaginas

}

//ConstruirPaginacion construtye la paginación en formato html para usarse en la página
func ConstruirPaginacion() string {
	var templateP string
	templateP += `
	<nav aria-label="Page navigation">
		<ul class="pagination">
			<li>
				<a href="/ListaCostos/1" aria-label="Primera">
				<span aria-hidden="true">&laquo;</span>
				</a>
			</li>`
	
	templateP += ``				
	for i := 0; i <= paginasTotales; i++ {
		if i == 1 {

			templateP += `<li class="active"><a href="/ListaCostos/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`
		}else if i > 1 && i < 11 {
			templateP += `<li><a href="/ListaCostos/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`

		}else if i > 11 && i == paginasTotales {
			templateP += `<li><span aria-hidden="true">...</span></li><li><a href="/ListaCostos/` + strconv.Itoa(i) + `">` + strconv.Itoa(i) + `</a></li>`				
		}
	}
	templateP += `<li><a href="/ListaCostos/` + strconv.Itoa(paginasTotales) + `" aria-label="Ultima"><span aria-hidden="true">&raquo;</span></a></li></ul></nav>`				
	return templateP
}


