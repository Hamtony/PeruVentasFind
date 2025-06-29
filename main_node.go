// Archivo: main_node.go
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Consulta struct {
	Producto string `json:"producto"`
}

type Resultado struct {
	Entidad string  `json:"entidad"`
	Score   float64 `json:"score"`
}

type Registro struct {
	Producto  string     `bson:"producto"`
	Resultados []Resultado `bson:"resultados"`
	Timestamp time.Time  `bson:"timestamp"`
}

var (
	mu                sync.Mutex
	resultadosGlobales []Resultado
	mongoCollection    *mongo.Collection
)

func main() {
	// Conexión a MongoDB
	ctx := context.TODO()
	clientOpts := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	mongoCollection = client.Database("perucompras").Collection("consultas")
	fmt.Println("Conectado a MongoDB")

	// Servidor TCP para comunicación con workers
	go iniciarTCP()

	// API REST para el frontend
	http.HandleFunc("/api/dataset", manejarCSV)
	http.HandleFunc("/api/recomendar", manejarRecomendacion)
	fmt.Println("API REST escuchando en :8080")
	http.ListenAndServe(":8080", nil)
}

func iniciarTCP() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Nodo principal (TCP) escuchando en puerto 8000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error al aceptar conexión:", err)
			continue
		}
		go manejarConexion(conn)
	}
}

func manejarConexion(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		var consulta Consulta
		if err := json.Unmarshal([]byte(line), &consulta); err != nil {
			log.Println("Error al decodificar consulta:", err)
			continue
		}

		fmt.Printf("Consulta recibida por TCP: %s\n", consulta.Producto)
		resultados := procesarConsulta(consulta.Producto)
		resJSON, _ := json.Marshal(resultados)
		conn.Write(resJSON)
		conn.Write([]byte("\n"))

		mu.Lock()
		resultadosGlobales = resultados
		mu.Unlock()

		registrarEnMongo(consulta.Producto, resultados)
	}
}

func procesarConsulta(producto string) []Resultado {
	producto = strings.ToLower(producto)
	return []Resultado{
		{"UNIVERSIDAD NACIONAL DE PIURA", 0.91},
		{"MINISTERIO DE EDUCACIÓN", 0.88},
		{"GOBIERNO REGIONAL DE LIMA", 0.85},
	}
}

func manejarCSV(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data/ReportePCBienes_cleaned.csv")
	if err != nil {
		http.Error(w, "No se pudo abrir el archivo CSV", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=ReportePCBienes_cleaned.csv")
	_, err = bufio.NewReader(file).WriteTo(w)
	if err != nil {
		http.Error(w, "Error al escribir el archivo CSV", http.StatusInternalServerError)
	}
}

func manejarRecomendacion(w http.ResponseWriter, r *http.Request) {
	producto := r.URL.Query().Get("producto")
	if producto == "" {
		http.Error(w, "Falta parámetro 'producto'", http.StatusBadRequest)
		return
	}

	res := procesarConsulta(producto)
	registrarEnMongo(producto, res)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func registrarEnMongo(producto string, resultados []Resultado) {
	doc := Registro{
		Producto:  producto,
		Resultados: resultados,
		Timestamp: time.Now(),
	}
	_, err := mongoCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Println("Error al registrar en MongoDB:", err)
	}
}
