import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useLogin } from '../auth/useAuth'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Wordmark } from '@/components/Wordmark'

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
    <div className="flex min-h-screen items-center justify-center bg-surface-200 p-4">
      <div className="flex w-full max-w-sm flex-col items-center gap-6">
        <h1 className="sr-only">Log in to Ticket Tout</h1>
        <Wordmark />
        <Card className="w-full rounded-[12px] border border-surface-400 shadow-card-highlight ring-0">
          <CardHeader>
            <CardTitle className="text-[1.5rem] leading-[1.15] tracking-[-0.02em] font-bold text-brand-blue">
              Log in
            </CardTitle>
            <CardDescription className="text-[0.6875rem] leading-[1.5] text-ink-600">
              Enter your credentials to continue.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="email"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Email
                </Label>
                <Input
                  id="email"
                  type="email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="password"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Password
                </Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              {login.isError && (
                <p className="text-body text-danger-600">{login.error.message}</p>
              )}
              <Button
                type="submit"
                disabled={login.isPending}
                className="h-11 w-full rounded-[8px] text-[0.8125rem] leading-none tracking-[0.12em] font-semibold uppercase font-display"
              >
                {login.isPending ? 'Logging in...' : 'Log in'}
              </Button>
            </form>
          </CardContent>
        </Card>
        <div className="flex items-center gap-4">
          <Link to="/terms" className="text-caption font-sans text-ink-600">
            Terms of Use
          </Link>
          <Link to="/accessibility" className="text-caption font-sans text-ink-600">
            Accessibility
          </Link>
        </div>
      </div>
    </div>
  )
}
