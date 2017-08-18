package MoConexion

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"../../Modulos/Variables"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//DataM es una estructura que contiene los datos de configuración en el archivo cfg
var DataM = MoVar.CargaSeccionCFG(MoVar.SecMongo)

//check verifica y escribe el error que elige el usuario y grita el error general
func check(err error, mensaje string) {
	if err != nil {
		fmt.Println("##########")
		fmt.Println(mensaje)
		fmt.Println("##########")
		panic(err)
	}
}

//GetColectionMgo regresa una colección específica de Mongo
func GetColectionMgo(coleccion string) (*mgo.Session, *mgo.Collection, error) {
	s, err := mgo.Dial(DataM.Servidor)
	if err != nil {
		return nil, nil, err
	}
	c := s.DB(DataM.NombreBase).C(coleccion)
	return s, c, nil
}

//GetConexionMgo regresa una sesion de mgo y error
func GetConexionMgo() (*mgo.Session, error) {
	session, err := mgo.Dial(DataM.Servidor)
	if err != nil {
		return session, err
	}
	return session, nil
}

//CloseConexionMgo cierra la sesion que se especifica
func CloseConexionMgo(sesion *mgo.Session) {
	sesion.Close()
}

//GetBaseMgo regresa un objeto database de mgo específico de una sesion específica
func GetBaseMgo(base string, sesion *mgo.Session) *mgo.Database {
	sesion.SetMode(mgo.Monotonic, true)
	return sesion.DB(base)
}

//InsertarImagen inserta una imagen en mongo y en el directorio ./Recursos/Imagenes
func InsertarImagen(file multipart.File, header *multipart.FileHeader) (string, error) {
	fmt.Println("FILE: ", file)
	defer file.Close()

	nombrefile := header.Filename
	dirpath := "./Recursos/Imagenes"

	//Comprobar directorio y crearlo
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		fmt.Println("el directorio no	 existe")
		os.MkdirAll(dirpath, 0777)
	}

	//subir imagen al servidor local
	out, err := os.Create("./Recursos/Imagenes/" + nombrefile)
	if err != nil {
		return "No es posible crear el archivo en el directorio, compruebe los permisos", err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "Error al escribir la imagen al directorio", err
	}

	//Inserar la imagen en Mongo
	idsImg, err := UploadImageToMongodb(dirpath, nombrefile)

	return idsImg, err
}

//UploadImageToMongodb Inserta Imagen en Mongo y devuelve su Id
func UploadImageToMongodb(path string, namefile string) (string, error) {

	db, err := mgo.Dial(DataM.Servidor)
	if err != nil {
		check(err, "Error al conectar con mongo")
		return "Error al conectar con mongo", err
	}
	base := db.DB(DataM.NombreBase)

	file2, err := os.Open(path + "/" + namefile)
	if err != nil {
		check(err, "Error al abrir el archivo o el archivo no existe")
		return "Error al abrir el archivo o el archivo no existe", err
	}
	defer file2.Close()

	stat, err := file2.Stat()
	if err != nil {
		check(err, "Error al leer el archivo")
		return "Error al leer el archivo", err
	}

	bs := make([]byte, stat.Size()) // read the file
	_, err = file2.Read(bs)
	if err != nil {
		check(err, "Error al crear objeto que contendrá el archivo")
		return "Error al crear objeto que contendrá el archivo", err
	}

	img, err := base.GridFS("Imagenes").Create(namefile)
	if err != nil {
		check(err, "error al crear archivo en mongo")
		return "Error al crear archivo en mongo", err
	}

	idsImg := img.Id()
	_, err = img.Write(bs)
	if err != nil {
		check(err, "error al escribir archivo en mongo")
		return "Error al escribir archivo en mongo", err
	}

	fmt.Println("File uploaded successfully to mongo ")
	err = img.Close()
	if err != nil {
		check(err, "error al cerrar img de mongo")
		return "Error al cerrar img de mongo", err
	}
	db.Close()
	idimg := getObjectIDToInterface(idsImg)

	return idimg.Hex(), nil
}

func getObjectIDToInterface(i interface{}) bson.ObjectId {
	var v = i.(bson.ObjectId)
	return v
}

//RegresaTagImagen regresa el tag de la imagen con ID correspondiente
func RegresaTagImagen(ID string) (string, error) {
	objid := bson.ObjectIdHex(ID)
	db, err := mgo.Dial(DataM.Servidor)
	if err != nil {
		check(err, "Error al conectar con mongo")
		return "Error al conectar con mongo", err
	}
	base := db.DB(DataM.NombreBase)

	img, err := base.GridFS("Imagenes").OpenId(objid)
	if err != nil {
		check(err, "Error al obtener imagen")
		return "Error al obtener imagen", err
	}
	b := make([]byte, img.Size())
	_, errim := img.Read(b)
	if errim != nil {
		fmt.Println("Error al leer la imagen...", err)
	}

	var tmp = ``
	///////////////////////////////////////////////////////////////////////////
	switch extension := filepath.Ext(img.Name()); extension {
	case ".jpg", ".jpeg":
		imagen, _ := jpeg.Decode(bytes.NewReader(b))
		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, imagen, nil); err != nil {
			log.Println("unable to encode image.")
			return "No es posible decodificar imagen", err
		}
		str := base64.StdEncoding.EncodeToString(buffer.Bytes())
		tmp = `data:image/jpg;base64,` + str

	case ".png":

		imagen, _ := png.Decode(bytes.NewReader(b))
		buffer := new(bytes.Buffer)
		if err := png.Encode(buffer, imagen); err != nil {
			log.Println("unable to encode image.")
			return "No es posible decodificar imagen", err
		}
		str := base64.StdEncoding.EncodeToString(buffer.Bytes())
		tmp = `data:image/png;base64,` + str

	}
	return tmp, nil
}
