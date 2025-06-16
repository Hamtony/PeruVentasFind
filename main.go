package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"
)

type OrdenCompra struct {
    FechaProceso     string  `csv:"FECHA_PROCESO"`
    AcuerdoMarco     string  `csv:"ACUERDO_MARCO"`
    Entidad          string  `csv:"ENTIDAD"`
    Proveedor        string  `csv:"PROVEEDOR"`
    TipoProcedimiento string `csv:"TIPO_PROCEDIMIENTO"`
    Subtotal         float64 `csv:"SUB_TOTAL"`
    IGV              float64 `csv:"IGV"`
    Total            float64 `csv:"TOTAL"`
}

func main() {
    file, err := os.Open("ordenes_perucompras.csv")
    if err != nil {
        log.Fatalf("No se pudo abrir el archivo CSV: %v", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = ',' // Asegura separador estándar

    records, err := reader.ReadAll()
    if err != nil {
        log.Fatalf("Error al leer el archivo CSV: %v", err)
    }

    fmt.Printf("Cargadas %d filas\n", len(records)-1) // -1 porque la primera es cabecera

    // Leer encabezado
    headers := records[0]
    headerIndex := make(map[string]int)
    for i, h := range headers {
        headerIndex[h] = i
    }

    // Parsear las filas
    var ordenes []OrdenCompra
    for _, row := range records[1:] {
        subtotal, _ := strconv.ParseFloat(row[headerIndex["SUB_TOTAL"]], 64)
        igv, _ := strconv.ParseFloat(row[headerIndex["IGV"]], 64)
        total, _ := strconv.ParseFloat(row[headerIndex["TOTAL"]], 64)

        orden := OrdenCompra{
            FechaProceso:     row[headerIndex["FECHA_PROCESO"]],
            AcuerdoMarco:     row[headerIndex["ACUERDO_MARCO"]],
            Entidad:          row[headerIndex["ENTIDAD"]],
            Proveedor:        row[headerIndex["PROVEEDOR"]],
            TipoProcedimiento: row[headerIndex["TIPO_PROCEDIMIENTO"]],
            Subtotal:         subtotal,
            IGV:              igv,
            Total:            total,
        }

        ordenes = append(ordenes, orden)
    }

    for i := 0; i < 3 && i < len(ordenes); i++ {
        fmt.Printf("%+v\n", ordenes[i])
    }

    // También podrías agregar filtros, agrupamientos, etc.
}
