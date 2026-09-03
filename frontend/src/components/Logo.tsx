interface LogoProps {
  className?: string
}

export function Logo({ className }: LogoProps) {
  return (
    <svg viewBox="0 0 64 64" className={className} aria-hidden="true">
      <path
        d="M14 0 H50 A14 14 0 0 1 64 14 V26.5 A5.5 5.5 0 0 0 64 37.5 V50 A14 14 0 0 1 50 64 H14 A14 14 0 0 1 0 50 V37.5 A5.5 5.5 0 0 0 0 26.5 V14 A14 14 0 0 1 14 0 Z"
        className="fill-brand-blue"
      />
      <rect x="27.5" y="21" width="9" height="28" rx="1.5" className="fill-paper-cream" />
      <rect x="15" y="15" width="34" height="9" rx="1.5" className="fill-accent-amber" />
    </svg>
  )
}
