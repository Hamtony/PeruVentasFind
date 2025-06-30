// worker_node.go
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Consulta struct {
	Producto string `json:"producto"`
}

type Resultado struct {
	Entidad string  `json:"entidad"`
	Score   float64 `json:"score"`
}

var ctx = context.Background()
var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println("Worker conectado a Redis")

	port := os.Getenv("WORKER_PORT")
	if port == "" {
		port = "9001"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Worker escuchando en puerto", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error al aceptar conexión:", err)
			continue
		}
		go manejarConsulta(conn)
	}
}

func manejarConsulta(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		var consulta Consulta
		if err := json.Unmarshal([]byte(line), &consulta); err != nil {
			log.Println("Error al decodificar consulta:", err)
			continue
		}

		key := strings.ToLower(consulta.Producto)
		cached, err := rdb.Get(ctx, key).Result()
		if err != nil {
			log.Println("Producto no encontrado en Redis")
			conn.Write([]byte("[]\n"))
			return
		}
		conn.Write([]byte(cached + "\n"))
	}
}
