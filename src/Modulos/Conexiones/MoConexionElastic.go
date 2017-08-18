package MoConexion

import (
	"context"
	"fmt"

	"../../Modulos/Variables"
	elastic "gopkg.in/olivere/elastic.v5"
)

//DataE es una estructura que contiene los datos de configuración en el archivo cfg
var DataE = MoVar.CargaSeccionCFG(MoVar.SecElastic)

//ctx Contexto
var ctx = context.Background()

//GetClienteElastic crea un nuevo cliente a elastic
func GetClienteElastic() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(DataE.BaseURL))
	if err != nil {
		fmt.Println("Error al crear un nuevo cliente", err)
		return nil, err
	}
	return client, nil
}

//VerificaIndex verifica que el indice exista en elastic
func VerificaIndex(client *elastic.Client) bool {
	exists, err := client.IndexExists(MoVar.IndexElastic).Do(ctx)
	if err != nil {
		fmt.Println("Error al verificar index de elastic", err)
		return false
	}
	if !exists {
		fmt.Printf("\n El Indice: %s no existe ", MoVar.IndexElastic)
		return false
	}
	return true
}

//VerificaIndexName verifica que el indice exista en elastic
func VerificaIndexName(client *elastic.Client, name string) bool {
	exists, err := client.IndexExists(name).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if !exists {
		fmt.Printf("\n El Indice: %s no existe ", name)
		return false
	}
	return true
}

//VerificaIndexName2 verifica que el indice exista en elastic
func VerificaIndexName2(client *elastic.Client, name string) error {
	exists, err := client.IndexExists(name).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !exists {
		fmt.Printf("\n El Indice: %s no existe ", name)
		return err
	}
	return nil
}

//InsertaElastic inserta un articulo de minisuper en elastic
func InsertaElastic(Type string, ID string, Data interface{}) bool {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println("Error al obtener el cliente: ", err)
		return false
	}
	defer client.Stop()
	if VerificaIndex(client) {
		Put, err := client.Index().Index(MoVar.IndexElastic).Type(Type).Id(ID).BodyJson(Data).Do(ctx)
		if err != nil {
			fmt.Println("Error al obtener el cliente de elastic ", err)
			return false
		}
		fmt.Printf("\nIndexado en el index %s, con type %s\n", Put.Index, Put.Type)
		return true
	}
	return false
}

//ActualizaElastic  actualiza correctamente un documento en elasticsearch
func ActualizaElastic(Type string, ID string, Data interface{}) error {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println("Error al obtener el cliente elasticSearch: ", err)
		return err
	}
	_, err = client.Update().Index(MoVar.IndexElastic).Type(Type).Id(ID).Doc(Data).DetectNoop(true).Do(context.TODO())
	if err != nil {
		fmt.Println("Error al Actualizar en elasticSearch", err)
		return err
	}
	defer client.Stop()
	return err
}

//DeleteElastic elimina un docuemnto de elastic por ID
func DeleteElastic(Type string, ID string) bool {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer client.Stop()

	_, err = client.Get().Index(MoVar.IndexElastic).Type(Type).Id(ID).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, err = client.Delete().Index(MoVar.IndexElastic).Type(Type).Id(ID).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//ConsultaElastic elimina un docuemnto de elastic por ID
func ConsultaElastic(Type string, ID string) bool {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer client.Stop()

	_, err = client.Get().Index(MoVar.IndexElastic).Type(Type).Id(ID).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//BuscaElastic busca documentos por un texto dado
func BuscaElastic(Type string, consulta *elastic.QueryStringQuery) (*elastic.SearchResult, bool) {
	client, err := GetClienteElastic()
	if err != nil {
		return nil, false
	}
	defer client.Stop()

	docs, err := client.Search().Index(MoVar.IndexElastic).Type(Type).Query(consulta).Do(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}

	return docs, true
}

//BusquedaElastic realiza una consulta al servidor de elastic, regresa resultados o un error
//Funcion realizada por ramon, para traer un error y así poder mostrar el error en la pantalla del cliente
func BusquedaElastic(Type string, consulta *elastic.QueryStringQuery) (*elastic.SearchResult, error) {
	client, err := GetClienteElastic()
	if err != nil {
		return nil, err
	}
	defer client.Stop()

	docs, err := client.Search().Index(MoVar.IndexElastic).Type(Type).Query(consulta).Do(ctx)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

//BuscaElasticAvanzada busca documentos por un texto dado
func BuscaElasticAvanzada(Type string, consulta *elastic.QueryStringQuery) (*elastic.SearchResult, error) {
	client, err := GetClienteElastic()
	if err != nil {
		return nil, err
	}
	defer client.Stop()
	var docs *elastic.SearchResult

	err1 := VerificaIndexName2(client, MoVar.IndexElastic)
	if err1 != nil {
		return nil, err1
	}

	docs, err2 := client.Search().Index(MoVar.IndexElastic).Type(Type).Query(consulta).From(0).Size(10000).Do(ctx)
	if err2 != nil {
		return nil, err2
	}

	return docs, nil
}

//FlushElastic hace flush a determinado index de elastic
func FlushElastic() {
	client, err := GetClienteElastic()
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.Flush().Index(MoVar.IndexElastic).Do(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
