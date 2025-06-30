import React, { useState } from 'react'

const categorias = [
  "LLANTAS, NEUMÁTICOS Y ACCESORIOS",
  "COMPUTADORAS PORTÁTILES",
  "EQUIPO MÉDICO",
  "MATERIAL DE OFICINA"
]

export default function ProductRecommender() {
  const [categoria, setCategoria] = useState('')
  const [resultados, setResultados] = useState<any[]>([])

  const consultar = async () => {
    const res = await fetch('http://localhost:8080/api/recomendar', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ producto: categoria })
    })
    const data = await res.json()
    setResultados(data)
  }

  return (
    <div>
      <h2>Selecciona una categoría para obtener recomendaciones</h2>
      <select onChange={(e) => setCategoria(e.target.value)} value={categoria}>
        <option value="">-- Selecciona --</option>
        {categorias.map((cat) => (
          <option key={cat} value={cat}>{cat}</option>
        ))}
      </select>
      <button onClick={consultar} disabled={!categoria}>Consultar</button>

      <ul>
        {resultados.map((r, i) => (
          <li key={i}>
            {r.entidad} - Score: {(r.score * 100).toFixed(1)}%
          </li>
        ))}
      </ul>
    </div>
  )
}
