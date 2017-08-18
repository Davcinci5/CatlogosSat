package Redis

import "github.com/go-redis/redis"

func getCliente() *redis.Client {
	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}

//CrearHash crea un hash referenciado por un ID de un Conjunto
func CrearHash(ID string, fields map[string]interface{}) error {
	client := getCliente()
	status := client.HMSet(ID, fields)
	err := status.Err()
	client.Close()
	return err
}

//RegresaDatoHash regresa el valor almacenado en un registro
func RegresaDatoHash(ID string, Field string) (string, error) {
	client := getCliente()
	status := client.HGet(ID, Field)
	client.Close()
	return status.Val(), status.Err()
}

//RegresaValoresHash regresa el valor almacenado en un registro
func RegresaValoresHash(ID string) (map[string]string, error) {
	client := getCliente()
	status := client.HGetAll(ID)
	client.Close()
	return status.Val(), status.Err()
}

// ActualizarHash actualiza un registro hash
func ActualizarHash(ID string, fields map[string]interface{}) error {
	client := getCliente()
	status := client.HMSet(ID, fields)
	err := status.Err()
	client.Close()
	return err
}

// EliminarFieldHash elimina un registro hash
func EliminarFieldHash(ID string, field string) error {
	client := getCliente()
	status := client.HDel(ID, field)
	client.Close()
	return status.Err()
}

//EliminarHash eleimina todos los registros de un Hash
func EliminarHash(ID string, fields []string) error {
	var err error
	for _, val := range fields {
		err = EliminarFieldHash(ID, val)
	}
	return err
}
