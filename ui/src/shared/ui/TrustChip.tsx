import type { TrustLabel } from '@/shared/types/host'
import { Badge } from './Badge'

const tone: Record<TrustLabel, 'ok' | 'accent' | 'default' | 'warn' | 'fail'> = {
  system: 'accent',
  verified: 'ok',
  agent: 'default',
  retrieved: 'warn',
  untrusted: 'fail',
  'poison-suspect': 'fail',
}

export function TrustChip({ trust }: { trust: TrustLabel }) {
  return <Badge tone={tone[trust]}>{trust}</Badge>
}
