import { useEffect, useState } from 'react'

interface Task {
  id: string
  title: string
  done: boolean
  created_at: string
}

export default function App() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [title, setTitle] = useState('')

  useEffect(() => {
    fetch('/api/tasks')
      .then(r => r.json())
      .then(setTasks)
  }, [])

  async function addTask(e: React.FormEvent) {
    e.preventDefault()
    if (!title.trim()) return
    const res = await fetch('/api/tasks', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title }),
    })
    const task: Task = await res.json()
    setTasks(prev => [...prev, task])
    setTitle('')
  }

  async function toggleTask(task: Task) {
    const res = await fetch(`/api/tasks/${task.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ done: !task.done }),
    })
    const updated: Task = await res.json()
    setTasks(prev => prev.map(t => (t.id === updated.id ? updated : t)))
  }

  async function deleteTask(id: string) {
    await fetch(`/api/tasks/${id}`, { method: 'DELETE' })
    setTasks(prev => prev.filter(t => t.id !== id))
  }

  return (
    <div style={{ maxWidth: 600, margin: '2rem auto', fontFamily: 'sans-serif' }}>
      <h1>Tasks</h1>
      <form onSubmit={addTask} style={{ display: 'flex', gap: 8, marginBottom: '1rem' }}>
        <input
          value={title}
          onChange={e => setTitle(e.target.value)}
          placeholder="New task…"
          style={{ flex: 1, padding: '0.5rem' }}
        />
        <button type="submit">Add</button>
      </form>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {tasks.map(task => (
          <li key={task.id} style={{ display: 'flex', alignItems: 'center', gap: 8, padding: '0.4rem 0', borderBottom: '1px solid #eee' }}>
            <input
              type="checkbox"
              checked={task.done}
              onChange={() => toggleTask(task)}
            />
            <span style={{ flex: 1, textDecoration: task.done ? 'line-through' : 'none' }}>
              {task.title}
            </span>
            <button onClick={() => deleteTask(task.id)}>✕</button>
          </li>
        ))}
      </ul>
    </div>
  )
}