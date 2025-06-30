import { cn } from '@/lib/utils'
import type { ScheduleStatus } from '@/shared/dashboard'
import { Badge } from '@/components/ui/badge'

interface StatusChipProps {
  status: ScheduleStatus
}

const statusConfig = {
  upcoming: {
    label: 'upcoming',
    className: 'bg-[#616161] text-white hover:bg-[#616161]/80',
  },
  'in-progress': {
    label: 'In progress',
    className: 'bg-[#ED6C02] text-white hover:bg-[#ED6C02]/80',
  },
  completed: {
    label: 'Completed',
    className: 'bg-[#2E7D32] text-white hover:bg-[#2E7D32]/80',
  },
  cancelled: {
    label: 'Cancelled',
    className: 'bg-[#D32F2F] text-white hover:bg-[#D32F2F]/80',
  },
} as const

export function StatusChip({ status }: StatusChipProps) {
  const config = statusConfig[status]

  return (
    <Badge
      variant="secondary"
      className={cn('text-xs font-normal px-2.5 py-1', config.className)}
    >
      {config.label}
    </Badge>
  )
}
