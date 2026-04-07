import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { useEffect, useState } from 'react';
export default function App() {
    const [tasks, setTasks] = useState([]);
    const [title, setTitle] = useState('');
    useEffect(() => {
        fetch('/api/tasks')
            .then(r => r.json())
            .then(setTasks);
    }, []);
    async function addTask(e) {
        e.preventDefault();
        if (!title.trim())
            return;
        const res = await fetch('/api/tasks', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ title }),
        });
        const task = await res.json();
        setTasks(prev => [...prev, task]);
        setTitle('');
    }
    async function toggleTask(task) {
        const res = await fetch(`/api/tasks/${task.id}`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ done: !task.done }),
        });
        const updated = await res.json();
        setTasks(prev => prev.map(t => (t.id === updated.id ? updated : t)));
    }
    async function deleteTask(id) {
        await fetch(`/api/tasks/${id}`, { method: 'DELETE' });
        setTasks(prev => prev.filter(t => t.id !== id));
    }
    return (_jsxs("div", { style: { maxWidth: 600, margin: '2rem auto', fontFamily: 'sans-serif' }, children: [_jsx("h1", { children: "Tasks" }), _jsxs("form", { onSubmit: addTask, style: { display: 'flex', gap: 8, marginBottom: '1rem' }, children: [_jsx("input", { value: title, onChange: e => setTitle(e.target.value), placeholder: "New task\u2026", style: { flex: 1, padding: '0.5rem' } }), _jsx("button", { type: "submit", children: "Add" })] }), _jsx("ul", { style: { listStyle: 'none', padding: 0 }, children: tasks.map(task => (_jsxs("li", { style: { display: 'flex', alignItems: 'center', gap: 8, padding: '0.4rem 0', borderBottom: '1px solid #eee' }, children: [_jsx("input", { type: "checkbox", checked: task.done, onChange: () => toggleTask(task) }), _jsx("span", { style: { flex: 1, textDecoration: task.done ? 'line-through' : 'none' }, children: task.title }), _jsx("button", { onClick: () => deleteTask(task.id), children: "\u2715" })] }, task.id))) })] }));
}
