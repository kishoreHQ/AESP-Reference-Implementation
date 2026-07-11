import { useMemo, useRef, useState } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import { Link, useNavigate } from 'react-router-dom'
import type { ExecutionTree, Mission, TreeNode } from '@/shared/types/host'
import { StatusDot } from '@/shared/ui/StatusDot'
import { Mono } from '@/shared/ui/Mono'
import { cn } from '@/shared/lib/cn'
import { Skeleton } from '@/shared/ui/Skeleton'
import { IconChevron } from '@/shared/ui/Icons'

function flatten(node: TreeNode, depth = 0, acc: { node: TreeNode; depth: number }[] = []) {
  acc.push({ node, depth })
  node.children?.forEach((c) => flatten(c, depth + 1, acc))
  return acc
}

export function MissionSpine({
  missions,
  tree,
  selectedMissionId,
  onSelectNode,
  compact,
}: {
  missions?: Mission[]
  tree?: ExecutionTree
  selectedMissionId?: string
  onSelectNode?: (id: string) => void
  compact?: boolean
}) {
  const nav = useNavigate()
  const parentRef = useRef<HTMLDivElement>(null)
  const [collapsed, setCollapsed] = useState<Record<string, boolean>>({})

  const rows = useMemo(() => {
    if (tree) {
      const flat = flatten(tree.root)
      return flat.filter((r) => {
        let p = r.node.parentId
        while (p) {
          if (collapsed[p]) return false
          const parent = flat.find((x) => x.node.id === p)
          p = parent?.node.parentId
        }
        return true
      })
    }
    return (missions || []).map((m) => ({
      node: {
        id: m.id,
        label: m.name,
        kind: 'agent' as const,
        status: (m.state === 'awaiting_approval'
          ? 'blocked'
          : m.state === 'running'
            ? 'running'
            : m.state === 'succeeded'
              ? 'succeeded'
              : m.state === 'failed'
                ? 'failed'
                : 'pending') as TreeNode['status'],
        costUsd: m.costUsd,
        durationMs: m.elapsedMs,
      },
      depth: 0,
      mission: m,
    }))
  }, [tree, missions, collapsed])

  const virtualizer = useVirtualizer({
    count: rows.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 40,
    overscan: 14,
  })

  if (!missions && !tree) {
    return (
      <div className="space-y-2 p-4">
        <Skeleton className="h-3 w-20" />
        <Skeleton className="h-9 w-full" />
        <Skeleton className="h-9 w-full" />
      </div>
    )
  }

  return (
    <aside className={cn('flex h-full min-h-0 flex-col bg-bg-1', compact && 'bg-transparent')} aria-label="Mission Spine">
      <div className="flex items-center justify-between border-b border-line px-3.5 py-3">
        <div>
          <div className="font-display text-[12px] font-bold uppercase tracking-[0.08em] text-ink-1">
            Mission Spine
          </div>
          <div className="mt-0.5 font-mono text-[11px] text-ink-2">
            {tree ? `${tree.nodeCount} nodes` : `${missions?.length ?? 0} missions`}
          </div>
        </div>
        <div className="h-1.5 w-1.5 rounded-full bg-live spine-pulse" title="Live" />
      </div>

      <div ref={parentRef} className="min-h-0 flex-1 overflow-auto">
        <div style={{ height: virtualizer.getTotalSize(), position: 'relative' }}>
          {virtualizer.getVirtualItems().map((vi) => {
            const row = rows[vi.index] as { node: TreeNode; depth: number; mission?: Mission }
            const hasChildren = (row.node.children?.length ?? 0) > 0
            const isSel = selectedMissionId === row.node.id
            return (
              <div
                key={row.node.id + vi.index}
                className={cn(
                  'absolute left-0 top-0 flex w-full cursor-pointer items-center gap-1.5 border-b border-line/40 pr-2 text-[12.5px] transition-colors',
                  'hover:bg-bg-2/80',
                  isSel && 'bg-[var(--accent-dim)]',
                )}
                style={{
                  height: vi.size,
                  transform: `translateY(${vi.start}px)`,
                  paddingLeft: 10 + row.depth * 14,
                }}
                onClick={() => {
                  if (row.mission) nav(`/missions/${row.mission.id}`)
                  else onSelectNode?.(row.node.id)
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    if (row.mission) nav(`/missions/${row.mission.id}`)
                    else onSelectNode?.(row.node.id)
                  }
                }}
                role="treeitem"
                tabIndex={0}
              >
                {hasChildren ? (
                  <button
                    className="flex h-6 w-6 shrink-0 items-center justify-center rounded text-ink-2 hover:bg-bg-3 hover:text-ink-0"
                    aria-label={collapsed[row.node.id] ? 'Expand' : 'Collapse'}
                    onClick={(e) => {
                      e.stopPropagation()
                      setCollapsed((c) => ({ ...c, [row.node.id]: !c[row.node.id] }))
                    }}
                  >
                    <IconChevron dir={collapsed[row.node.id] ? 'right' : 'down'} />
                  </button>
                ) : (
                  <span className="w-6 shrink-0" />
                )}
                <StatusDot status={row.node.status} />
                <span className="min-w-0 flex-1 truncate font-medium text-ink-0">{row.node.label}</span>
                {!compact && row.node.costUsd != null && (
                  <Mono className="shrink-0 text-[11px]">${row.node.costUsd.toFixed(2)}</Mono>
                )}
              </div>
            )
          })}
        </div>
      </div>

      {missions && missions.length === 0 && (
        <div className="p-4 text-[12.5px] leading-relaxed text-ink-1">
          No missions in the spine.{' '}
          <Link className="text-accent hover:underline" to="/">
            Launch one
          </Link>
        </div>
      )}
    </aside>
  )
}
