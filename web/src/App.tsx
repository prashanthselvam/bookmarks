import { useEffect, useState } from 'react'
import './App.css'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

function App() {
  const [message, setMessage] = useState<string>('Loading...')
  const [error, setError] = useState<string>('')

  useEffect(() => {
    fetch(API_URL)
      .then(res => res.text())
      .then(data => setMessage(data))
      .catch(err => setError(err.message))
  }, [])

  return (
    <div className="App">
      <h1>Bookmarks App</h1>
      {error ? (
        <p style={{ color: 'red' }}>Error: {error}</p>
      ) : (
        <p>API says: {message}</p>
      )}
    </div>
  )
}

export default App