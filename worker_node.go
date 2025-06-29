// Archivo: worker_node.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Consulta struct {
	Producto string `json:"producto"`
}

type Resultado struct {
	Entidad string  `json:"entidad"`
	Score   float64 `json:"score"`
}

func main() {
	mainAddr := os.Getenv("MAIN_ADDR")
	if mainAddr == "" {
		mainAddr = "main_node:8000"
	}
	fmt.Printf("Conectando al nodo principal en %s...\n", mainAddr)

	conn, err := net.Dial("tcp", mainAddr)
	if err != nil {
		log.Fatalf("No se pudo conectar al nodo principal: %v", err)
	}
	defer conn.Close()

	consulta := Consulta{Producto: "COMPUTADORAS PORTÁTILES"}
	payload, _ := json.Marshal(consulta)
	conn.Write(payload)
	conn.Write([]byte("\n"))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		var resultados []Resultado
		if err := json.Unmarshal([]byte(line), &resultados); err != nil {
			log.Println("Error al decodificar resultados:", err)
			continue
		}
		fmt.Println("Resultados recibidos:")
		for _, r := range resultados {
			fmt.Printf("Entidad: %s, Score: %.2f\n", r.Entidad, r.Score)
		}
		break // Solo una consulta para esta demo
	}
}
