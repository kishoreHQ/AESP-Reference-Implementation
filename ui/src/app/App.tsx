import { Navigate, Route, Routes } from 'react-router-dom'
import { AppShell } from './AppShell'
import { DashboardPage } from '@/features/dashboard/DashboardPage'
import { MissionsPage } from '@/features/missions/MissionsPage'
import { MissionDetailPage } from '@/features/missions/MissionDetailPage'
import { ApprovalsPage } from '@/features/approvals/ApprovalsPage'
import { FleetPage } from '@/features/fleet/FleetPage'
import { MemoryPage } from '@/features/memory/MemoryPage'
import { ArtifactsPage } from '@/features/artifacts/ArtifactsPage'
import { EvaluationsPage } from '@/features/evaluations/EvaluationsPage'
import { ReplayPage } from '@/features/replay/ReplayPage'
import { SettingsPage } from '@/features/settings/SettingsPage'
import { ConnectionsPage } from '@/features/connections/ConnectionsPage'
import { LiveSessionPage } from '@/features/agents/LiveSessionPage'
import { ControlRoomPage } from '@/features/agents/ControlRoomPage'
import { BoardPage } from '@/features/board/BoardPage'
import { RoutinesPage } from '@/features/routines/RoutinesPage'

export default function App() {
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route index element={<DashboardPage />} />
        <Route path="missions" element={<MissionsPage />} />
        <Route path="missions/:id" element={<MissionDetailPage />} />
        <Route path="approvals" element={<ApprovalsPage />} />
        <Route path="fleet" element={<FleetPage />} />
        <Route path="connect" element={<ConnectionsPage />} />
        <Route path="agents/:id" element={<ControlRoomPage />} />
        <Route path="sessions/new" element={<LiveSessionPage />} />
        <Route path="sessions/:id" element={<LiveSessionPage />} />
        <Route path="board" element={<BoardPage />} />
        <Route path="routines" element={<RoutinesPage />} />
        <Route path="memory" element={<MemoryPage />} />
        <Route path="artifacts" element={<ArtifactsPage />} />
        <Route path="evaluations" element={<EvaluationsPage />} />
        <Route path="replay/:runId" element={<ReplayPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="spine" element={<Navigate to="/missions" replace />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  )
}
