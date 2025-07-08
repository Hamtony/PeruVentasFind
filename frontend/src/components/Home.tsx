// src/components/Home.tsx
import { useEffect, useState } from 'react'
import axios from 'axios'
import styles from './Home.module.css'

interface EntidadRecomendada {
  entidad: string
  score: number
}

function Home() {
  const [entidades, setEntidades] = useState<EntidadRecomendada[]>([])
  const [producto, setProducto] = useState('')
  const [categorias, setCategorias] = useState<string[]>([])

  useEffect(() => {
    axios.get('http://localhost:8080/api/categorias')
      .then(res => setCategorias(res.data))
      .catch(err => console.error('Error al obtener categorías:', err))
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const res = await axios.post(
        'http://localhost:8080/api/recomendar',
        { producto },
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      )
      setEntidades(res.data)
    } catch (error) {
      console.error('Error al obtener recomendaciones:', error)
    }
  }

  return (
    <div className={styles.page}>
      <div className={styles.container}>
        <h1 className={styles.title}>Recomendador de Entidades Públicas</h1>
        <form onSubmit={handleSubmit} className={styles.form}>
          <select
            value={producto}
            onChange={(e) => setProducto(e.target.value)}
            className={styles.input}
          >
            <option value="">Seleccione una categoría</option>
            {categorias.map((cat, idx) => (
              <option key={idx} value={cat}>{cat}</option>
            ))}
          </select>
          <button type="submit" className={styles.button}>Recomendar</button>
        </form>
        <div className={styles.results}>
          {entidades.map((entidad, idx) => (
            <div key={idx} className={styles.entity}>
              <strong>{entidad.entidad}</strong> — Score: {entidad.score.toFixed(4)}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default Home
