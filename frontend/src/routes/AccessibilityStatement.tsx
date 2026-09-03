import { Link } from 'react-router-dom'
import { Wordmark } from '@/components/Wordmark'

export function AccessibilityStatement() {
  return (
    <div className="min-h-screen bg-surface-200 p-4">
      <div className="mx-auto flex max-w-sm flex-col gap-6 pb-12">
        <h1 className="sr-only">Accessibility statement</h1>
        <div className="flex items-center justify-between">
          <Wordmark />
          <Link
            to="/login"
            className="text-label-caps uppercase font-display text-ink-600"
          >
            Back
          </Link>
        </div>

        <h2 className="text-h1 font-display text-brand-blue">
          Accessibility statement
        </h2>

        <p className="text-body font-sans text-ink-600">
          Ticket Tout targets RGAA (Référentiel Général d'Amélioration de
          l'Accessibilité) level AA. Every screen follows a heading structure
          a screen reader can navigate, form fields carry visible labels,
          interactive elements are real buttons and links rather than styled
          divs, and status information (pending, declined, expired) is
          always stated in words, never conveyed by color alone.
        </p>

      </div>
    </div>
  )
}
