import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { cn } from '@/shared/lib/cn'

export function Button({
  children,
  variant = 'primary',
  size = 'md',
  className,
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost' | 'warn'
  size?: 'sm' | 'md'
  children: ReactNode
}) {
  const v = {
    primary:
      'bg-gradient-to-b from-[#33ddff] to-[#00a1e6] text-[#041018] font-semibold hover:brightness-110 shadow-[0_0_20px_rgba(0,191,255,0.25)]',
    secondary:
      'bg-bg-2 text-ink-0 border border-[var(--line-strong)] hover:border-[rgba(0,191,255,0.35)] hover:bg-bg-3',
    danger: 'bg-fail/90 text-white hover:brightness-110 shadow-[0_0_16px_rgba(248,113,113,0.2)]',
    ghost: 'bg-transparent text-ink-1 hover:bg-bg-2 hover:text-ink-0',
    warn: 'bg-gradient-to-b from-[#ffb347] to-[#ff9f1c] text-[#1a1006] font-semibold shadow-[0_0_18px_rgba(255,159,28,0.25)]',
  }
  const s = {
    sm: 'min-h-9 px-2.5 text-[13px] gap-1.5',
    md: 'min-h-10 px-3.5 text-[13.5px] gap-2',
  }
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center rounded-[8px] font-medium transition-all duration-[var(--motion)] disabled:pointer-events-none disabled:opacity-40',
        v[variant],
        s[size],
        className,
      )}
      {...props}
    >
      {children}
    </button>
  )
}
