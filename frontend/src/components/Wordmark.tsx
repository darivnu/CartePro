import { Logo } from './Logo'

export function Wordmark() {
  return (
    <div className="flex items-center gap-2">
      <Logo className="h-5 w-5" />
      <p className="text-wordmark font-display text-brand-blue">CartePro</p>
    </div>
  )
}
