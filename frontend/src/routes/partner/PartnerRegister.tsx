import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useRegisterPartner } from '../../partners/usePartnersData'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Wordmark } from '@/components/Wordmark'

export function PartnerRegisterPage() {
  const [businessName, setBusinessName] = useState('')
  const [siret, setSiret] = useState('')
  const [category, setCategory] = useState('')
  const [address, setAddress] = useState('')
  const [region, setRegion] = useState('')
  const [contactEmail, setContactEmail] = useState('')
  const [password, setPassword] = useState('')
  const navigate = useNavigate()
  const register = useRegisterPartner()

  function handleSubmit(event: React.SyntheticEvent<HTMLFormElement>) {
    event.preventDefault()
    register.mutate(
      {
        business_name: businessName,
        siret: Number(siret),
        category,
        address,
        region,
        contact_email: contactEmail,
        password,
      },
      {
        onSuccess: () => {
          navigate('/partner', { replace: true })
        },
      },
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-200 p-4">
      <div className="flex w-full max-w-sm flex-col items-center gap-6">
        <h1 className="sr-only">Create a CartePro partner account</h1>
        <Wordmark />
        <Card className="w-full rounded-[12px] border border-surface-400 shadow-card-highlight ring-0">
          <CardHeader>
            <CardTitle className="text-[1.5rem] leading-[1.15] tracking-[-0.02em] font-bold text-brand-blue">
              Register your business
            </CardTitle>
            <CardDescription className="text-[0.6875rem] leading-[1.5] text-ink-600">
              Create a partner account to start accepting CartePro payments.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="businessName"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Business name
                </Label>
                <Input
                  id="businessName"
                  type="text"
                  autoComplete="organization"
                  value={businessName}
                  onChange={(event) => setBusinessName(event.target.value)}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="siret"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Siret
                </Label>
                <Input
                  id="siret"
                  type="text"
                  inputMode="numeric"
                  pattern="[0-9]*"
                  value={siret}
                  onChange={(event) => setSiret(event.target.value.replace(/\D/g, ''))}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="category"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Category
                </Label>
                <Input
                  id="category"
                  type="text"
                  value={category}
                  onChange={(event) => setCategory(event.target.value)}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="address"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Address
                </Label>
                <Input
                  id="address"
                  type="text"
                  autoComplete="street-address"
                  value={address}
                  onChange={(event) => setAddress(event.target.value)}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="region"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Region
                </Label>
                <Input
                  id="region"
                  type="text"
                  autoComplete="address-level1"
                  value={region}
                  onChange={(event) => setRegion(event.target.value)}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label
                  htmlFor="contactEmail"
                  className="text-[0.6875rem] leading-none tracking-[0.16em] font-semibold uppercase font-display text-ink-600"
                >
                  Contact email
                </Label>
                <Input
                  id="contactEmail"
                  type="email"
                  autoComplete="email"
                  value={contactEmail}
                  onChange={(event) => setContactEmail(event.target.value)}
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
                  autoComplete="new-password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  className="h-11 rounded-[8px] px-3 text-[0.8125rem] leading-[1.5] md:text-[0.8125rem]"
                />
              </div>
              {register.isError && (
                <p className="text-body text-danger-600">{register.error.message}</p>
              )}
              <Button
                type="submit"
                disabled={register.isPending}
                className="h-11 w-full rounded-[8px] text-[0.8125rem] leading-none tracking-[0.12em] font-semibold uppercase font-display"
              >
                {register.isPending ? 'Creating account...' : 'Create account'}
              </Button>
            </form>
          </CardContent>
        </Card>
        <div className="flex items-center gap-1.5 text-body-strong font-sans text-ink-900">
          Already have an account?
          <Link to="/login" className="underline">
            Log in
          </Link>
        </div>
        <div className="flex items-center gap-1.5 text-body-strong font-sans text-ink-900">
          Looking to use your card?
          <Link to="/register" className="underline">
            Register as a client
          </Link>
        </div>
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
