interface LogoProps {
  className?: string
}

export function Logo({ className }: LogoProps) {
  return (
    <svg viewBox="0 0 64 64" className={className} aria-hidden="true">
      <rect width="64" height="64" rx="16" className="fill-brand-blue" />
      <path
        d="M40 18.1A16 16 0 1 0 40 45.9"
        fill="none"
        stroke="#FFFFFF"
        strokeWidth="9"
        strokeLinecap="round"
      />
      <circle cx="45" cy="32" r="5.5" className="fill-coral" />
    </svg>
  )
}
