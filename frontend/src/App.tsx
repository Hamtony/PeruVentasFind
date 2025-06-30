import { useState } from 'react';
import axios from 'axios';

interface Entidad {
  entidad: string;
  score: number;
}

function App() {
  const [producto, setProducto] = useState('');
  const [entidades, setEntidades] = useState<Entidad[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const baseURL = 'http://localhost:8080/api';

  const recomendar = async () => {
    setLoading(true);
    setError('');
    setEntidades([]);
    try {
      const res = await axios.post(`${baseURL}/recomendar`, { producto });
      setEntidades(res.data);
    } catch (err) {
      console.error('Error al obtener recomendaciones:', err);
      setError('Ocurrió un error al obtener recomendaciones.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6 flex flex-col items-center">
      <h1 className="text-3xl font-bold text-blue-800 mb-6">
        Recomendador de Entidades Estatales
      </h1>

      <div className="w-full max-w-md bg-white rounded-2xl shadow p-6">
        <input
          type="text"
          value={producto}
          onChange={(e) => setProducto(e.target.value)}
          placeholder="Ejemplo: Computadora portátil"
          className="w-full p-3 mb-4 rounded-xl border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-400"
        />
        <button
          onClick={recomendar}
          disabled={!producto || loading}
          className="w-full bg-blue-600 text-white p-3 rounded-xl hover:bg-blue-700 transition disabled:opacity-50"
        >
          {loading ? 'Consultando...' : 'Consultar'}
        </button>

        {error && <p className="text-red-600 mt-4">{error}</p>}

        {entidades.length > 0 && (
          <div className="mt-6">
            <h2 className="text-xl font-semibold mb-2">Entidades sugeridas:</h2>
            <ul className="space-y-3">
              {entidades.map((item, index) => (
                <li
                  key={index}
                  className="bg-blue-100 border-l-4 border-blue-600 p-4 rounded-lg"
                >
                  <p className="font-medium">{item.entidad}</p>
                  <p className="text-sm text-gray-700">
                    Score: {(item.score * 100).toFixed(2)}%
                  </p>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
