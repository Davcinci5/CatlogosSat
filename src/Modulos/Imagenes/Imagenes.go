package Imagenes

import (
	"html/template"

	"../../Modulos/Conexiones"
	"gopkg.in/mgo.v2/bson"
)

//CargaTemplateImagenes Regresa template de imagenes que hacen referencia al arreglo de ObjectIds
func CargaTemplateImagenes(IDImg []bson.ObjectId) (template.HTML, error) {
	htmlImagenes := ``
	for _, id := range IDImg {
		ruta, err := MoConexion.RegresaTagImagen(id.Hex())
		if err != nil {
			return template.HTML(htmlImagenes), err
		}
		htmlImagenes += `<img id="Imagen" alt="Responsive image" name="ImagenProducto" width="200px" height="200px";" src="` + ruta + `">`
	}
	return template.HTML(htmlImagenes), nil
}

//CargaTemplateImagenesCompras retorna imagenes para el modulo de compras
func CargaTemplateImagenesCompras(IDImg []bson.ObjectId) (string, error) {
	htmlImagenes := ``
	for _, id := range IDImg {
		ruta, err := MoConexion.RegresaTagImagen(id.Hex())
		if err != nil {
			return htmlImagenes, err
		}
		htmlImagenes += `<img id="Imagen" alt="Responsive image" name="ImagenProducto" width="80px" height="80px";" src="` + ruta + `">`
	}
	return htmlImagenes, nil
}

//CargaTemplateImagenesVentas retorna imagenes para el modulo de ventas
func CargaTemplateImagenesVentas(IDImg []bson.ObjectId) (string, error) {
	htmlImagenes := ``
	for _, id := range IDImg {
		ruta, err := MoConexion.RegresaTagImagen(id.Hex())
		if err != nil {
			return htmlImagenes, err
		}
		htmlImagenes += `<img id="Imagen" alt="Responsive image" name="ImagenProducto" width="45px" height="45px";" src="` + ruta + `">`
	}
	return htmlImagenes, nil
}
