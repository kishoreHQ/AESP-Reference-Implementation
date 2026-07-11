import { NavLink, Outlet, useLocation } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { useApprovals, useHealth, useMissions } from '@/shared/api/hooks'
import { MissionSpine } from '@/shared/spine/MissionSpine'
import { detectProfile } from '@/profiles/detect'
import { useUi } from '@/shared/store/ui'
import { Toasts } from '@/shared/ui/Toasts'
import { EventBridge } from '@/shared/events/EventBridge'
import { useQueryClient } from '@tanstack/react-query'
import { cn } from '@/shared/lib/cn'
import { StatusDot } from '@/shared/ui/StatusDot'
import { AgentRail } from '@/features/agents/AgentRail'
import { BrainRail } from '@/features/brain/BrainRail'
import {
  IconApprovals,
  IconArtifacts,
  IconEval,
  IconFleet,
  IconLogo,
  IconMemory,
  IconMissions,
  IconSettings,
  IconSpine,
} from '@/shared/ui/Icons'

const sections: { label: string; items: { to: string; label: string; icon: typeof IconMissions; end?: boolean }[] }[] = [
  {
    label: 'Main',
    items: [
      { to: '/', label: 'Dashboard', icon: IconLogo, end: true },
      { to: '/missions', label: 'Missions', icon: IconMissions },
      { to: '/approvals', label: 'Approvals', icon: IconApprovals },
      { to: '/board', label: 'Board', icon: IconMissions },
      { to: '/routines', label: 'Routines', icon: IconEval },
    ],
  },
  {
    label: 'Orchestration',
    items: [
      { to: '/fleet', label: 'Fleet', icon: IconFleet },
      { to: '/connect', label: 'Connect', icon: IconFleet },
      { to: '/evaluations', label: 'Evaluations', icon: IconEval },
    ],
  },
  {
    label: 'Knowledge',
    items: [
      { to: '/memory', label: 'Memory', icon: IconMemory },
      { to: '/artifacts', label: 'Artifacts', icon: IconArtifacts },
    ],
  },
  {
    label: 'Config',
    items: [{ to: '/settings', label: 'Settings', icon: IconSettings }],
  },
]

export function AppShell() {
  const profile = detectProfile()
  const health = useHealth()
  const missions = useMissions()
  const approvals = useApprovals()
  const { railExpanded, toggleRail, spineOpen, setSpineOpen, theme, setTheme } = useUi()
  const loc = useLocation()
  const qc = useQueryClient()
  const pending = approvals.data?.length ?? 0
  const healthy = health.data?.status === 'ok'
  const [clock, setClock] = useState(() => new Date().toLocaleTimeString([], { hour12: false }))

  useEffect(() => {
    const t = setInterval(() => setClock(new Date().toLocaleTimeString([], { hour12: false })), 1000)
    return () => clearInterval(t)
  }, [])

  useEffect(() => {
    EventBridge.connect()
    const unsub = EventBridge.subscribe((ev) => {
      if (ev.type.startsWith('mission') || ev.type.startsWith('approval') || ev.type.startsWith('node')) {
        qc.invalidateQueries({ queryKey: ['missions'] })
        qc.invalidateQueries({ queryKey: ['approvals'] })
        qc.invalidateQueries({ queryKey: ['tree'] })
      }
    })
    return () => {
      unsub()
    }
  }, [qc])

  const showSpine = loc.pathname.startsWith('/missions/')

  return (
    <div className="flex h-full min-h-0 flex-col bg-bg-0 text-ink-0">
      <div className="hidden min-h-0 flex-1 md:flex">
        {/* Sidebar — Control Hub style */}
        <nav
          className={cn(
            'flex shrink-0 flex-col border-r border-[var(--line)] bg-bg-1/95 backdrop-blur-md transition-[width] duration-[var(--motion)]',
            railExpanded ? 'w-[240px]' : 'w-[72px]',
          )}
          aria-label="Primary"
        >
          <button
            className="flex h-16 items-center gap-2.5 border-b border-[var(--line)] px-3.5 text-left"
            onClick={toggleRail}
            aria-label="Toggle navigation"
          >
            <span className="relative flex h-9 w-9 items-center justify-center rounded-[10px] border border-[rgba(0,191,255,0.3)] bg-[rgba(0,191,255,0.1)] text-accent glow-cyan">
              <IconLogo size={18} />
            </span>
            {railExpanded && (
              <div className="min-w-0">
                <div className="truncate font-display text-[12px] font-bold tracking-[0.12em] text-ink-0">
                  CONTROL<span className="text-accent">HUB</span>
                </div>
                <div className="truncate font-mono text-[9px] text-ink-2">AESP Agent OS · host UI</div>
              </div>
            )}
          </button>

          <div className="flex-1 overflow-y-auto px-2 py-3">
            {sections.map((sec) => (
              <div key={sec.label} className="mb-4">
                {railExpanded && (
                  <div className="mb-1.5 px-2 text-[9px] font-bold uppercase tracking-[0.16em] text-ink-2">
                    {sec.label}
                  </div>
                )}
                <div className="flex flex-col gap-0.5">
                  {sec.items.map((n) => {
                    const Icon = n.icon
                    return (
                      <NavLink
                        key={n.to}
                        to={n.to}
                        end={n.end}
                        title={n.label}
                        className={({ isActive }) =>
                          cn(
                            'relative flex min-h-10 items-center gap-3 rounded-[10px] px-2.5 text-[13px] font-medium transition-all',
                            'text-ink-1 hover:bg-bg-2 hover:text-ink-0',
                            isActive &&
                              'bg-gradient-to-r from-[rgba(0,191,255,0.16)] to-transparent text-cyan-100 shadow-[inset_0_0_0_1px_rgba(0,191,255,0.15)]',
                            !railExpanded && 'justify-center px-0',
                          )
                        }
                      >
                        {({ isActive }) => (
                          <>
                            {isActive && (
                              <span className="absolute left-0 top-1/2 h-5 w-[3px] -translate-y-1/2 rounded-r bg-accent shadow-[0_0_8px_rgba(0,191,255,0.8)]" />
                            )}
                            <Icon size={18} className={isActive ? 'text-accent' : undefined} />
                            {railExpanded && <span className="truncate">{n.label}</span>}
                            {n.to === '/approvals' && pending > 0 && (
                              <span
                                className={cn(
                                  'flex h-5 min-w-5 items-center justify-center rounded-full bg-warn px-1 font-mono text-[10px] font-bold text-[#1a1006] shadow-[0_0_12px_rgba(255,159,28,0.4)]',
                                  railExpanded ? 'ml-auto' : 'absolute right-1 top-1',
                                )}
                              >
                                {pending}
                              </span>
                            )}
                          </>
                        )}
                      </NavLink>
                    )
                  })}
                </div>
              </div>
            ))}
          </div>

          <div className="space-y-2 border-t border-[var(--line)] p-3">
            <div
              className={cn(
                'flex items-center gap-2 rounded-[10px] border border-[var(--line)] bg-bg-0/70 px-2.5 py-2',
                !railExpanded && 'justify-center',
              )}
            >
              <StatusDot status={healthy ? 'ok' : 'degraded'} />
              {railExpanded && (
                <div className="min-w-0">
                  <div className="truncate text-[11px] font-semibold text-ink-0">{profile.label}</div>
                  <div className="truncate font-mono text-[10px] text-ink-2">
                    {clock} · {profile.useMocks ? 'mocks' : 'live'}
                  </div>
                </div>
              )}
            </div>
            {railExpanded && (
              <div className="flex gap-1">
                <button
                  className="flex-1 rounded-[8px] border border-[var(--line)] bg-bg-2 px-2 py-1.5 text-[10px] font-semibold uppercase tracking-wide text-ink-1 hover:border-[rgba(0,191,255,0.35)] hover:text-accent"
                  onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
                >
                  {theme}
                </button>
                <button
                  className="flex-1 rounded-[8px] border border-[var(--line)] bg-bg-2 px-2 py-1.5 text-[10px] font-semibold uppercase tracking-wide text-ink-1 hover:border-[rgba(0,191,255,0.35)] hover:text-accent"
                  onClick={() => setSpineOpen(!spineOpen)}
                >
                  spine
                </button>
              </div>
            )}
          </div>
        </nav>

        <div className="hidden w-[200px] shrink-0 xl:block">
          <AgentRail />
        </div>

        <main className="min-w-0 flex-1 overflow-auto">
          <Outlet />
        </main>

        {loc.pathname === '/' && (
          <div className="hidden w-[280px] shrink-0 2xl:block">
            <BrainRail />
          </div>
        )}

        {(showSpine || spineOpen) && (
          <div
            className={cn(
              'hidden shrink-0 border-l border-[var(--line)] bg-bg-1/90 lg:block',
              'w-[300px]',
            )}
          >
            <MissionSpine missions={missions.data} selectedMissionId={loc.pathname.split('/')[2]} />
          </div>
        )}
      </div>

      {/* Mobile */}
      <div className="flex min-h-0 flex-1 flex-col md:hidden">
        <header className="flex h-12 items-center justify-between border-b border-[var(--line)] bg-bg-1/95 px-3 backdrop-blur-md">
          <div className="flex items-center gap-2">
            <span className="flex h-7 w-7 items-center justify-center rounded-[8px] border border-[rgba(0,191,255,0.3)] bg-[rgba(0,191,255,0.1)] text-accent">
              <IconLogo size={14} />
            </span>
            <span className="font-display text-[12px] font-bold tracking-[0.1em]">
              CONTROL<span className="text-accent">HUB</span>
            </span>
          </div>
          <span className="online-pill text-[9px]">
            <span className="h-1.5 w-1.5 rounded-full bg-neon-green spine-pulse" />
            online
          </span>
        </header>
        <main className="min-h-0 flex-1 overflow-auto pb-[calc(3.75rem+env(safe-area-inset-bottom))]">
          <Outlet />
        </main>
        <nav
          className="fixed bottom-0 left-0 right-0 z-40 grid h-[3.75rem] grid-cols-5 border-t border-[var(--line)] bg-bg-1/95 pb-[env(safe-area-inset-bottom)] backdrop-blur-md"
          aria-label="Mobile"
        >
          {[
            { to: '/', label: 'Home', icon: IconLogo },
            { to: '/approvals', label: 'Approve', icon: IconApprovals, badge: pending },
            { to: '/missions', label: 'Missions', icon: IconMissions },
            { to: '/fleet', label: 'Fleet', icon: IconFleet },
            { to: '/settings', label: 'More', icon: IconSettings },
          ].map((t) => {
            const Icon = t.icon
            return (
              <NavLink
                key={t.to}
                to={t.to}
                end={t.to === '/'}
                className={({ isActive }) =>
                  cn(
                    'relative flex flex-col items-center justify-center gap-0.5 text-[10px] font-medium text-ink-2',
                    isActive && 'text-accent',
                  )
                }
              >
                <Icon size={18} />
                {t.label}
                {t.badge ? (
                  <span className="absolute right-[16%] top-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-warn px-1 font-mono text-[9px] font-bold text-[#1a1006]">
                    {t.badge}
                  </span>
                ) : null}
              </NavLink>
            )
          })}
        </nav>
      </div>

      <button
        className="fixed bottom-4 right-4 z-30 hidden items-center gap-1.5 rounded-[10px] border border-[var(--line)] bg-bg-1 px-3 py-2 text-[12px] font-medium text-ink-1 hover:border-[rgba(0,191,255,0.4)] hover:text-accent hover:glow-cyan md:flex lg:hidden"
        onClick={() => setSpineOpen(!spineOpen)}
      >
        <IconSpine size={14} /> Spine
      </button>
      {spineOpen && !showSpine && (
        <div className="fixed inset-y-0 right-0 z-40 hidden w-[300px] border-l border-[var(--line)] bg-bg-1 md:block lg:hidden">
          <MissionSpine missions={missions.data} />
        </div>
      )}

      <Toasts />
    </div>
  )
}
