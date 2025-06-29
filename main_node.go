// Archivo: main_node.go
package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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
	rdb                *redis.Client
	ctx               = context.Background()
)

func main() {
	ctx := context.TODO()
	clientOpts := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	mongoCollection = client.Database("perucompras").Collection("consultas")
	fmt.Println("Conectado a MongoDB")

	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("No se pudo conectar a Redis:", err)
	}
	fmt.Println("Conectado a Redis")

	// Entrenamiento automático cada 30 min
	go entrenarModeloPeriodicamente()
	go iniciarTCP()

	http.HandleFunc("/api/dataset", manejarCSV)
	http.HandleFunc("/api/recomendar", manejarRecomendacion)
	fmt.Println("API REST escuchando en :8080")
	http.ListenAndServe(":8080", nil)
}

func entrenarModeloPeriodicamente() {
	for {
		fmt.Println("Entrenando modelo...")
		entrenarModeloDesdeCSV("data/ReportePCBienes_cleaned.csv")
		time.Sleep(30 * time.Minute)
	}
}

func entrenarModeloDesdeCSV(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Println("No se pudo abrir el archivo CSV para entrenamiento:", err)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	_, _ = r.Read() // Saltar encabezado

	contador := make(map[string]map[string]int) // producto -> entidad -> frecuencia

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(rec) < 3 {
			continue
		}
		producto := strings.ToLower(rec[1])
		entidad := rec[2]

		if contador[producto] == nil {
			contador[producto] = make(map[string]int)
		}
		contador[producto][entidad]++
	}

	for prod, entidades := range contador {
		total := 0
		for _, freq := range entidades {
			total += freq
		}
		var ranking []Resultado
		for entidad, freq := range entidades {
			ranking = append(ranking, Resultado{
				Entidad: entidad,
				Score:   float64(freq) / float64(total),
			})
		}
		sort.Slice(ranking, func(i, j int) bool {
			return ranking[i].Score > ranking[j].Score
		})
		if len(ranking) > 10 {
			ranking = ranking[:10]
		}
		bytes, _ := json.Marshal(ranking)
		rdb.Set(ctx, prod, bytes, time.Hour)
	}
	fmt.Println("Modelo actualizado y almacenado en Redis")
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
	productoKey := strings.ToLower(producto)

	// Buscar en cache Redis
	cached, err := rdb.Get(ctx, productoKey).Result()
	if err == nil {
		var res []Resultado
		err := json.Unmarshal([]byte(cached), &res)
		if err == nil && len(res) > 0 {
			fmt.Println("Resultado recuperado del modelo en Redis")
			return res
		}
	}

	// Respuesta por defecto
	res := []Resultado{
		{"UNIVERSIDAD NACIONAL DE PIURA", 0.91},
		{"MINISTERIO DE EDUCACIÓN", 0.88},
		{"GOBIERNO REGIONAL DE LIMA", 0.85},
	}
	fmt.Println("Se usa respuesta por defecto")
	return res
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
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error al escribir el archivo CSV", http.StatusInternalServerError)
	}
}

func manejarRecomendacion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido, usa POST", http.StatusMethodNotAllowed)
		return
	}

	var consulta Consulta
	err := json.NewDecoder(r.Body).Decode(&consulta)
	if err != nil || strings.TrimSpace(consulta.Producto) == "" {
		http.Error(w, "Error al leer el cuerpo o campo 'producto' inválido", http.StatusBadRequest)
		return
	}

	res := procesarConsulta(consulta.Producto)
	err = registrarEnMongo(consulta.Producto, res)
	if err != nil {
		http.Error(w, "Error al registrar la consulta en MongoDB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func registrarEnMongo(producto string, resultados []Resultado) error {
	doc := Registro{
		Producto:  producto,
		Resultados: resultados,
		Timestamp: time.Now(),
	}
	_, err := mongoCollection.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Println("Error al registrar en MongoDB:", err)
		return err
	}
	return nil
}
