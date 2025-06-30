import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/schedule-completed')({
  component: ScheduleCompleted,
})

import { Button } from '@/components/ui/button'
import { X, Calendar, Clock } from 'lucide-react'
import { useNavigate } from '@tanstack/react-router'

interface ScheduleCompletedData {
  date: string
  time: string
  duration: string
}

// Mock data - in a real app this would come from the previous clock-out session
const mockCompletedData: ScheduleCompletedData = {
  date: 'Mon, 15 January 2025',
  time: '10:30 - 11:30 SGT',
  duration: '(1 hour)',
}

function SuccessAnimation() {
  return (
    <div className="relative flex items-center justify-center w-36 h-36">
      {/* Enhanced illustration matching Figma design */}
      <svg
        width="140"
        height="140"
        viewBox="0 0 141 140"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        {/* Main orange circle with check */}
        <path
          d="M71.001 34.978H71C52.4516 34.978 37.415 50.0145 37.415 68.563V68.564C37.415 87.1125 52.4516 102.149 71 102.149H71.001C89.5495 102.149 104.586 87.1125 104.586 68.564V68.563C104.586 50.0145 89.5495 34.978 71.001 34.978Z"
          fill="#FF8046"
        />
        <path
          d="M56.256 68.851L65.951 78.546L85.744 58.751"
          stroke="white"
          strokeWidth="5"
          strokeLinecap="round"
          strokeLinejoin="round"
        />

        {/* Decorative curved lines */}
        <path
          d="M105.097 107.24C107.972 107.906 110.535 109.527 112.368 111.84C113.778 114.33 114.645 117.091 114.91 119.94"
          stroke="#BAE0DB"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <path
          d="M31.979 84.343C29.862 86.78 24.321 91.154 19.087 89.153C12.545 86.653 12.737 82.803 14.854 81.264C16.971 79.725 19.854 82.226 20.434 84.343C21.014 86.46 20.049 93.579 14.854 96.273C9.659 98.967 2.732 97.62 1 96.273"
          stroke="#EFE0FF"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />
        <path
          d="M115.391 16.228C112.891 15.907 107.656 17.382 106.732 25.849C105.808 34.316 98.395 36.047 94.803 35.854"
          stroke="#FF8046"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />

        {/* Decorative circles */}
        <circle cx="21.108" cy="51.109" r="4.618" fill="#F99A3D" />
        <circle cx="49.032" cy="14.881" r="3.848" fill="#FFE7EE" />
        <circle cx="60.769" cy="22.18" r="2.296" fill="#FFE7EE" />
        <circle cx="122.126" cy="65.759" r="4.137" fill="#EFE0FF" />
        <circle cx="36.116" cy="112.708" r="4.137" fill="#FFE7EE" />
        <circle cx="87.683" cy="124.83" r="4.137" fill="#F9D2DD" />
        <circle cx="74.887" cy="120.693" r="2.296" fill="#FF8046" />
      </svg>
    </div>
  )
}

export default function ScheduleCompleted() {
  const navigate = useNavigate()
  const data = mockCompletedData

  const handleGoHome = () => {
    navigate({ to: '/' })
  }

  const handleClose = () => {
    navigate({ to: '/' })
  }

  return (
    <div className="min-h-screen bg-[#0D5D59] flex flex-col max-w-md mx-auto relative">
      <button
        onClick={handleClose}
        className="absolute top-4 right-6 p-2 text-white hover:bg-white/10 rounded-full transition-colors z-10"
      >
        <X className="h-6 w-6" />
      </button>

      <div className="flex-1 flex flex-col items-center justify-center px-6 py-12 space-y-12">
        <SuccessAnimation />

        <div className="text-center space-y-8">
          <h1 className="text-3xl font-bold text-white">Schedule Completed</h1>

          <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-6 space-y-4 border border-white/20">
            <div className="flex items-center justify-center gap-3 text-white">
              <Calendar className="h-5 w-5 text-cyan-300" />
              <span className="text-lg">{data.date}</span>
            </div>

            <div className="flex items-center justify-center gap-3 text-white">
              <Clock className="h-5 w-5 text-cyan-300" />
              <div className="text-center">
                <div className="text-lg">{data.time}</div>
                <div className="text-sm text-white/80">{data.duration}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="p-6">
        <Button
          onClick={handleGoHome}
          variant="outline"
          className="w-full bg-transparent border-2 border-white text-white hover:bg-white hover:text-[#0D5D59] py-4 text-lg font-semibold transition-all duration-300"
        >
          Go to Home
        </Button>
      </div>
    </div>
  )
}
