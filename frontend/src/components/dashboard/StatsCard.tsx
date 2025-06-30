import { cn } from '@/lib/utils'

interface StatsCardProps {
  value: number
  label: string
  variant: 'error' | 'warning' | 'success'
  className?: string
}

const variantStyles = {
  error: 'text-[#D32F2F]',
  warning: 'text-[#ED6C02]',
  success: 'text-[#2E7D32]',
} as const

export function StatsCard({
  value,
  label,
  variant,
  className,
}: StatsCardProps) {
  return (
    <div
      className={cn(
        'flex flex-col items-center justify-center gap-1 flex-1 p-4 bg-white rounded-2xl shadow-sm border border-gray-100',
        className,
      )}
    >
      <div className={cn('text-3xl font-bold', variantStyles[variant])}>
        {value}
      </div>
      <div className="text-xs text-gray-600 text-center leading-tight">
        {label}
      </div>
    </div>
  )
}
