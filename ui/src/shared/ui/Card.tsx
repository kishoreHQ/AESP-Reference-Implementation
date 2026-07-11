import type { HTMLAttributes } from 'react'
import { cn } from '@/shared/lib/cn'

export function Card({ className, children, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        'rounded-[12px] border border-[var(--line)] bg-gradient-to-b from-[rgba(18,31,45,0.95)] to-[rgba(12,21,32,0.98)] p-4',
        'shadow-[inset_0_1px_0_rgba(255,255,255,0.03),0_8px_28px_rgba(0,0,0,0.25)]',
        'transition-[border-color,box-shadow,transform] duration-[var(--motion)]',
        className,
      )}
      {...props}
    >
      {children}
    </div>
  )
}
