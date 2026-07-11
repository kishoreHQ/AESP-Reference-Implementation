import type { ReactNode } from 'react'
import { cn } from '@/shared/lib/cn'

export function Mono({ children, className }: { children: ReactNode; className?: string }) {
  return <span className={cn('font-mono text-[12px] tracking-tight text-ink-1', className)}>{children}</span>
}
