import type { ReactNode } from 'react'

/** Minimal stroke icons — no emoji chrome. */

type P = { className?: string; size?: number }

function Svg({ size = 18, className, children }: P & { children: ReactNode }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.75"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden
    >
      {children}
    </svg>
  )
}

export function IconMissions(p: P) {
  return (
    <Svg {...p}>
      <rect x="3" y="4" width="18" height="16" rx="2" />
      <path d="M3 10h18M9 4v16" />
    </Svg>
  )
}
export function IconApprovals(p: P) {
  return (
    <Svg {...p}>
      <path d="M12 3l8 4v6c0 5-3.5 8-8 9-4.5-1-8-4-8-9V7l8-4z" />
      <path d="M9 12l2 2 4-4" />
    </Svg>
  )
}
export function IconFleet(p: P) {
  return (
    <Svg {...p}>
      <circle cx="12" cy="12" r="3" />
      <path d="M12 2v3M12 19v3M2 12h3M19 12h3M4.9 4.9l2.1 2.1M17 17l2.1 2.1M4.9 19.1L7 17M17 7l2.1-2.1" />
    </Svg>
  )
}
export function IconMemory(p: P) {
  return (
    <Svg {...p}>
      <path d="M4 7h16v10H4z" />
      <path d="M8 7V5a4 4 0 018 0v2M8 17v2a4 4 0 008 0v-2" />
    </Svg>
  )
}
export function IconArtifacts(p: P) {
  return (
    <Svg {...p}>
      <path d="M14 3H6a2 2 0 00-2 2v14a2 2 0 002 2h12a2 2 0 002-2V9l-6-6z" />
      <path d="M14 3v6h6" />
    </Svg>
  )
}
export function IconEval(p: P) {
  return (
    <Svg {...p}>
      <path d="M4 19V5M4 19h16" />
      <path d="M8 16l3-5 3 3 4-7" />
    </Svg>
  )
}
export function IconSettings(p: P) {
  return (
    <Svg {...p}>
      <circle cx="12" cy="12" r="3" />
      <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
    </Svg>
  )
}
export function IconSearch(p: P) {
  return (
    <Svg {...p} size={16}>
      <circle cx="11" cy="11" r="7" />
      <path d="M20 20l-3-3" />
    </Svg>
  )
}
export function IconPlus(p: P) {
  return (
    <Svg {...p} size={16}>
      <path d="M12 5v14M5 12h14" />
    </Svg>
  )
}
export function IconChevron(p: P & { dir?: 'down' | 'right' }) {
  return (
    <Svg {...p} size={14}>
      {p.dir === 'right' ? <path d="M9 6l6 6-6 6" /> : <path d="M6 9l6 6 6-6" />}
    </Svg>
  )
}
export function IconLogo(p: P) {
  return (
    <Svg {...p} size={20}>
      <path d="M4 18V6l8-3 8 3v12l-8 3-8-3z" />
      <path d="M12 3v18M4 10h16" />
    </Svg>
  )
}
export function IconSpine(p: P) {
  return (
    <Svg {...p}>
      <path d="M8 4v16M8 8h8M8 12h10M8 16h6" />
    </Svg>
  )
}
