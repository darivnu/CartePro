import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useLogin } from '../auth/useAuth'

export function LoginPage() {
const [email, setEmail] = useState('')
const [password, setPassword] = useState('')
const navigate = useNavigate()
const login = useLogin()

function handleSubmit(event: React.SyntheticEvent<HTMLFormElement>) {
    event.preventDefault()
    login.mutate(
    { email, password },
    {
        onSuccess: () => {
        navigate('/', { replace: true })
        },
    },
    )
}

return (
    <form onSubmit={handleSubmit} className="p-4">
    <div>
        <label htmlFor="email">Email</label>
        <input
        id="email"
        type="email"
        value={email}
        onChange={(event) => setEmail(event.target.value)}
        />
    </div>
    <div>
        <label htmlFor="password">Password</label>
        <input
        id="password"
        type="password"
        value={password}
        onChange={(event) => setPassword(event.target.value)}
        />
    </div>
    {login.isError && <p>{login.error.message}</p>}
    <button type="submit" disabled={login.isPending}>
        {login.isPending ? 'Logging in...' : 'Log in'}
    </button>
    </form>
)
}
