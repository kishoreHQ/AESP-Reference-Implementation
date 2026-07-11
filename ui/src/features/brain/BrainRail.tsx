import { useState } from 'react'
import { useAddJournal, useGoals, useJournal, useMemorySearch } from '@/shared/api/hooks'
import { TrustChip } from '@/shared/ui/TrustChip'
import { Mono } from '@/shared/ui/Mono'
import { Button } from '@/shared/ui/Button'

/** Right-rail brain: goals · journal · memory (unified memory, INV-04). */
export function BrainRail() {
  const goals = useGoals()
  const journal = useJournal()
  const addJ = useAddJournal()
  const [q, setQ] = useState('')
  const [entry, setEntry] = useState('')
  const mem = useMemorySearch(q)

  return (
    <aside className="flex h-full flex-col border-l border-[var(--line)] bg-bg-1/90">
      <div className="border-b border-[var(--line)] px-3 py-3">
        <div className="font-display text-[11px] font-bold uppercase tracking-[0.12em] text-ink-2">Brain</div>
        <div className="font-mono text-[10px] text-ink-2">unified memory · INV-04</div>
      </div>
      <div className="flex-1 space-y-4 overflow-auto p-3">
        <section>
          <div className="section-label">Goals</div>
          {(goals.data || []).map((g: any) => (
            <div key={g.id} className="mb-2 rounded-[8px] border border-[var(--line)] bg-bg-0/50 p-2">
              <div className="text-[12px] font-medium">{g.title}</div>
              <div className="mt-1 h-1.5 overflow-hidden rounded-full bg-bg-3">
                <div
                  className="h-full rounded-full bg-gradient-to-r from-cyan-400 to-neon-green"
                  style={{ width: `${Math.round((g.progress || 0) * 100)}%` }}
                />
              </div>
              <Mono className="text-[10px]">{Math.round((g.progress || 0) * 100)}%</Mono>
            </div>
          ))}
        </section>

        <section>
          <div className="section-label">Journal</div>
          <textarea
            className="input mb-1 min-h-[4rem] w-full text-[12px]"
            placeholder="Quick entry…"
            value={entry}
            onChange={(e) => setEntry(e.target.value)}
          />
          <Button
            size="sm"
            className="w-full"
            onClick={() => {
              if (entry.trim()) {
                addJ.mutate(entry)
                setEntry('')
              }
            }}
          >
            Save entry
          </Button>
          <ul className="mt-2 space-y-1">
            {(journal.data || []).slice(0, 5).map((j: any) => (
              <li key={j.id} className="rounded border border-[var(--line)] px-2 py-1 text-[11px] text-ink-1">
                {j.text}
              </li>
            ))}
          </ul>
        </section>

        <section>
          <div className="section-label">Memory</div>
          <input
            className="input mb-2 w-full text-[12px]"
            placeholder="Search…"
            value={q}
            onChange={(e) => setQ(e.target.value)}
          />
          <ul className="space-y-1">
            {(mem.data || []).slice(0, 6).map((m: any) => (
              <li key={m.id} className="rounded border border-[var(--line)] p-1.5 text-[11px]">
                <TrustChip trust={m.trust || 'agent'} />
                <div className="mt-0.5 text-ink-0">{m.text}</div>
              </li>
            ))}
          </ul>
        </section>
      </div>
    </aside>
  )
}
