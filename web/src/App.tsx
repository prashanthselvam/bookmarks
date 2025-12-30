import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [message, setMessage] = useState<string>('Loading...')
  const [error, setError] = useState<string>('')

  useEffect(() => {
    fetch('https://backend-empty-sun-8345.fly.dev/')
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