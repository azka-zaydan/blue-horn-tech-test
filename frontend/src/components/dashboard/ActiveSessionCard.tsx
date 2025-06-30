import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { MapPin, Clock } from 'lucide-react'
import { useNavigate } from '@tanstack/react-router'

interface ActiveSessionProps {
  customerName: string
  location: string
  time: string
  duration: string
  avatarUrl?: string
  scheduleId: string
}

export function ActiveSessionCard({
  customerName,
  location,
  time,
  duration,
  avatarUrl,
  scheduleId,
}: ActiveSessionProps) {
  const navigate = useNavigate()

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
  }

  return (
    <div className="bg-[#0D5D59] rounded-2xl p-6 text-white space-y-4">
      <div className="text-center">
        <div className="text-3xl font-bold tracking-wider">{duration}</div>
      </div>

      <div className="flex items-center gap-3">
        <Avatar className="h-12 w-12">
          <AvatarImage src={avatarUrl} />
          <AvatarFallback className="bg-white/20 text-white">
            {getInitials(customerName)}
          </AvatarFallback>
        </Avatar>
        <div>
          <div className="font-semibold text-lg">{customerName}</div>
        </div>
      </div>

      <div className="space-y-2">
        <div className="flex items-center gap-2 text-white/90">
          <MapPin className="h-4 w-4" />
          <span className="text-sm">{location}</span>
        </div>
        <div className="flex items-center gap-2 text-white/90">
          <Clock className="h-4 w-4" />
          <span className="text-sm">{time}</span>
        </div>
      </div>

      <Button
        className="w-full bg-white text-[#0D5D59] hover:bg-white/90 font-semibold"
        onClick={() =>
          navigate({ to: '/clock-out/$scheduleId', params: { scheduleId } })
        }
      >
        Clock-Out
      </Button>
    </div>
  )
}
