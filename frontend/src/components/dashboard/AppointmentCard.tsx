import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { StatusChip } from './StatusChip'
import type { Schedule } from '@/shared/dashboard'
import { MoreHorizontal, MapPin, Calendar, Clock } from 'lucide-react'
import { useNavigate } from '@tanstack/react-router'
import type React from 'react'
import { useStartVisit } from '@/hooks/useClockIn' // Import the useStartVisit hook
import { getContext } from '@/integrations/tanstack-query/root-provider' // Import getContext
import { useState } from 'react' // Import useState

interface AppointmentCardProps {
  schedule: Schedule
}

export function AppointmentCard({ schedule }: AppointmentCardProps) {
  const navigate = useNavigate()
  const mutation = useStartVisit() // Initialize the mutation hook
  const { queryClient } = getContext() // Get the queryClient
  const [gettingLocation, setGettingLocation] = useState(false) // State to manage location fetching status
  const [mutationErrorMessage, setMutationErrorMessage] = useState<
    string | null
  >(null) // State for mutation errors

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
  }

  const handleCardClick = () => {
    navigate({
      to: '/schedule/$scheduleId',
      params: { scheduleId: schedule.id },
    })
  }

  const handleStartVisit = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation() // Prevent card click when button is clicked
    setMutationErrorMessage(null) // Clear any previous error messages
    setGettingLocation(true) // Indicate that location is being fetched

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords
          setGettingLocation(false)

          mutation.mutate(
            {
              scheduleId: schedule.id,
              latitude: latitude,
              longitude: longitude,
            },
            {
              onSuccess: () => {
                setMutationErrorMessage(null)
                queryClient.invalidateQueries({
                  queryKey: ['schedules', schedule.id],
                })
                navigate({
                  to: `/clock-out/$scheduleId`,
                  params: { scheduleId: schedule.id },
                })
              },
              onError: (error) => {
                setMutationErrorMessage(
                  error?.message || 'An unexpected error occurred.',
                )
              },
            },
          )
        },
        (error) => {
          console.error('Error getting user location:', error)
          setGettingLocation(false)
          setMutationErrorMessage(
            'Could not get your current location. Please enable location services.',
          )
        },
        { enableHighAccuracy: true, timeout: 5000, maximumAge: 0 },
      )
    } else {
      setGettingLocation(false)
      console.error('Geolocation is not supported by this browser.')
      setMutationErrorMessage('Geolocation is not supported by this browser.')
    }
  }

  const renderActionButtons = () => {
    const handleButtonClick = (e: React.MouseEvent, action?: string) => {
      e.stopPropagation() // Prevent card click when button is clicked
      if (action === 'clock-out') {
        navigate({
          to: '/clock-out/$scheduleId',
          params: { scheduleId: schedule.id },
        })
      }
      if (action === 'view-progress') {
        navigate({
          to: '/schedule/$scheduleId',
          params: { scheduleId: schedule.id },
        })
      }
      if (action === 'view-report') {
        navigate({
          to: '/schedule/$scheduleId',
          params: { scheduleId: schedule.id },
        })
      }
    }

    const isClockInButtonDisabled = mutation.isPending || gettingLocation

    switch (schedule.status) {
      case 'upcoming':
        return (
          <>
            <Button
              size="sm"
              className="w-full bg-[#0D5D59] hover:bg-[#0D5D59]/90"
              onClick={handleStartVisit} // Use the new handleStartVisit function
              disabled={isClockInButtonDisabled}
            >
              {gettingLocation
                ? 'Getting Location...'
                : mutation.isPending
                  ? 'Clocking In...'
                  : 'Clock-In Now'}
            </Button>
            {mutationErrorMessage && (
              <p className="text-red-500 text-sm mt-2 text-center">
                Error: {mutationErrorMessage}
              </p>
            )}
          </>
        )

      case 'in-progress':
        return (
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              className="flex-1 border-blue-300 text-[#0D5D59] hover:bg-blue-50"
              onClick={(e: React.MouseEvent) =>
                handleButtonClick(e, 'view-progress')
              }
            >
              View Progress
            </Button>
            <Button
              size="sm"
              className="flex-1 bg-[#0D5D59] hover:bg-[#0D5D59]/90"
              onClick={(e: React.MouseEvent) =>
                handleButtonClick(e, 'clock-out')
              }
            >
              Clock-Out Now
            </Button>
          </div>
        )

      case 'completed':
        return (
          <Button
            variant="outline"
            size="sm"
            className="w-full border-blue-300 text-[#0D5D59] hover:bg-blue-50"
            onClick={(e: React.MouseEvent) =>
              handleButtonClick(e, 'view-report')
            }
          >
            View Report
          </Button>
        )

      case 'cancelled':
        return null

      default:
        return null
    }
  }

  return (
    <div
      className="bg-white rounded-2xl shadow-sm p-4 space-y-3 cursor-pointer hover:shadow-md transition-shadow"
      onClick={handleCardClick}
    >
      <div className="flex items-center justify-between">
        <StatusChip status={schedule.status} />
        <button
          className="p-1 hover:bg-gray-100 rounded"
          onClick={(e) => e.stopPropagation()}
        >
          <MoreHorizontal className="h-5 w-5 text-gray-800" />
        </button>
      </div>

      <div className="flex items-center gap-2">
        <Avatar className="h-10 w-10">
          <AvatarImage src={''} />
          <AvatarFallback className="bg-gray-200 text-gray-700">
            {getInitials(schedule.client_name)}
          </AvatarFallback>
        </Avatar>
        <div className="flex-1">
          <div className="font-semibold text-gray-900 text-sm">
            {schedule.client_name}
          </div>
          <div className="text-xs text-gray-600">{schedule.client_name}</div>
        </div>
      </div>

      <div className="flex items-center gap-1">
        <MapPin className="h-4 w-4 text-gray-600" />
        <span className="text-xs text-gray-600">{schedule.location}</span>
      </div>

      <div className="bg-blue-50 rounded-lg p-2">
        <div className="flex items-center justify-between text-xs text-gray-600">
          <div className="flex items-center gap-1">
            <Calendar className="h-4 w-4 text-cyan-600" />
            <span>{schedule.shiftDate}</span>
          </div>
          <span>|</span>
          <div className="flex items-center gap-1">
            <Clock className="h-4 w-4 text-cyan-600" />
            <span>{schedule.shiftTimeOnly}</span>
          </div>
        </div>
      </div>

      {renderActionButtons()}
    </div>
  )
}
