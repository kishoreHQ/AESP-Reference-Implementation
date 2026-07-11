import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useArtifacts } from '@/shared/api/hooks'
import { Card } from '@/shared/ui/Card'
import { Mono } from '@/shared/ui/Mono'

export function ArtifactsPage() {
  const { data } = useArtifacts()
  const [sel, setSel] = useState(data?.[0]?.id)
  const item = data?.find((a) => a.id === sel) ?? data?.[0]

  return (
    <div className="mx-auto grid max-w-6xl gap-4 p-4 md:grid-cols-2 md:p-6">
      <div>
        <h1 className="mb-4 font-display text-[26px] font-bold">Artifacts</h1>
        <div className="space-y-2">
          {data?.map((a) => (
            <Card key={a.id} className={`cursor-pointer ${item?.id === a.id ? 'border-accent' : ''}`} onClick={() => setSel(a.id)}>
              <div className="font-medium">{a.name}</div>
              <Mono>
                {a.digest.slice(0, 18)}… · v{a.version} · {(a.sizeBytes / 1024).toFixed(1)}KB
              </Mono>
              <div className="mt-1 text-[12px] text-ink-1">
                mission{' '}
                <Link className="text-accent" to={`/missions/${a.missionId}`}>
                  {a.missionId}
                </Link>
              </div>
            </Card>
          ))}
        </div>
      </div>
      <div>
        <h2 className="mb-2 font-display text-[16px] font-semibold">Preview</h2>
        {item && (
          <Card>
            <div className="mb-2 text-[12px] text-ink-1">Provenance: {item.provenance.join(' → ')}</div>
            <pre className="max-h-[480px] overflow-auto rounded-[6px] bg-bg-0 p-3 font-mono text-[12.5px]">
              {item.contentPreview || 'No preview'}
            </pre>
            <Link className="mt-3 inline-block text-accent" to={`/missions/${item.missionId}`}>
              Open producing mission →
            </Link>
          </Card>
        )}
      </div>
    </div>
  )
}
