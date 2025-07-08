# PeruVentasFind 🚀  

**Solución tecnológica para optimizar las ventas al Estado Peruano mediante programación concurrente y distribuida**  

## 📌 Descripción  

PeruVentasFind es una aplicación web diseñada para ayudar a proveedores y emprendedores a identificar oportunidades estratégicas de venta al Estado Peruano. Utilizando técnicas de **machine learning** y un sistema distribuido basado en **Go (Golang)**, analiza patrones en los datos abiertos de compras públicas para:  

- 🔍 **Detectar entidades estatales con mayor demanda** de productos específicos.  
- 📊 **Identificar irregularidades** (sobreprecios, concentración de proveedores).  
- 🎯 **Recomendar estrategias** basadas en datos históricos.  

## 🛠️ Tecnologías  

| **Categoría**       | **Tecnologías**                                                                 |
|----------------------|---------------------------------------------------------------------------------|
| **Backend**          | Go (Golang), Goroutines, TCP/IP, Redis (cola de mensajes)                      |
| **Base de Datos**    | MongoDB (almacenamiento de modelos y resultados)                               |
| **Frontend**         | React.js (interfaz interactiva)                                                |
| **Infraestructura**  | Docker, Docker Compose (contenedores para MongoDB, Redis y nodos de trabajo)   |
| **Metodología**      | Scrum (Kanban, Sprints)                                                        |

## 🚀 Funcionalidades  

✅ **Motor de recomendaciones**  
   - Clasifica entidades públicas por afinidad con productos específicos (ej: "computadoras portátiles").  

✅ **Procesamiento concurrente**  
   - Cluster de 3 nodos worker que distribuyen tareas usando **patrones Fan-Out/Fan-In**.  

✅ **API RESTful**  
   - Endpoints para consultas en tiempo real (`/api/recomendar`) y carga de datasets (`/api/dataset`).  

✅ **Persistencia escalable**  
   - MongoDB guarda modelos entrenados con **transacciones ACID**. Redis gestiona colas de tareas.  

## 📊 Arquitectura  

```mermaid
graph TD
    A[Usuario] -->|React| B[Frontend]
    B -->|HTTP| C[Main Node (Go)]
    C -->|TCP| D[Worker Node 1]
    C -->|TCP| E[Worker Node 2]
    C -->|TCP| F[Worker Node 3]
    C -->|Consultas| G[MongoDB]
    C -->|Caché/Colas| H[Redis]
