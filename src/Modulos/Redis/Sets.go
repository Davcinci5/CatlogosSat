package Redis

import "fmt"

// InsertarAConjunto inserta a un conjunto un valor dado
func InsertarAConjunto(nombreConjunto string, valor string) (bool, error) {
	client := getCliente()
	ds := client.SAdd(nombreConjunto, valor)
	_, err := ds.Result()
	if err != nil {
		fmt.Println(err)
		client.Close()
		return false, err
	}
	//client.BgSave()
	client.Close()
	return true, err
}

//ObtenerMiembrosdeGrupo Regresa todos los IDS de los miembros de un conjunto
func ObtenerMiembrosdeGrupo(nombreConjuntoOiD string) ([]string, error) {
	client := getCliente()
	ds := client.SMembers(nombreConjuntoOiD)
	resul, err := ds.Result()
	if err != nil {
		fmt.Println(err)
	}
	client.Close()
	return resul, err
}

//ObtenerMiemmbroDeGrupo verifica si el nombreConjunto pasado cuenta con el elemento en su interior
func ObtenerMiemmbroDeGrupo(nombreConjuntoOID string, valor string) (bool, error) {
	client := getCliente()
	ds := client.SIsMember(nombreConjuntoOID, valor)
	res, err := ds.Result()
	if err != nil {
		fmt.Println(err)
	}
	client.Close()
	return res, err
}

//InsertaRedis crea un conjunto en redis
func InsertaRedis(IDUsuarioOGrupo string, conjuntoDeIDS []string) error {
	var err error
	for _, val := range conjuntoDeIDS {
		_, err = InsertarAConjunto(IDUsuarioOGrupo, val)
	}
	return err
}

// EliminarConjunto elimina el conjunto dado
func EliminarConjunto(nombreconjunto string) error {
	client := getCliente()
	ds := client.Del(nombreconjunto)
	_, err := ds.Result()
	//client.BgSave()
	client.Close()
	return err
}

// EliminaMiembroDeGrupo Elimina el  miembro (valor) del conjunto especimicado
func EliminaMiembroDeGrupo(nombreConjunto string, valor string) error {
	client := getCliente()
	ds := client.SRem(nombreConjunto, valor)
	client.Close()
	return ds.Err()
}
