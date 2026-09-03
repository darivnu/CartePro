import { cn } from '@/lib/utils'

interface SimulationNoticeProps {
  tone?: 'light' | 'dark'
  className?: string
}

export function SimulationNotice({ tone = 'light', className }: SimulationNoticeProps) {
  return (
    <p
      className={cn(
        'text-micro-caps uppercase font-sans',
        tone === 'dark' ? 'text-graphite-400' : 'text-ink-600',
        className,
      )}
    >
      Simulation — no real monetary value
    </p>
  )
}
