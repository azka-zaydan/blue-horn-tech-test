import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/schedule/$scheduleId')({
  component: ScheduleDetails,
})

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
  ArrowLeft,
  Calendar,
  Clock,
  Mail,
  Phone,
  MapPin,
  CheckCircle,
  XCircle,
} from 'lucide-react' // Import CheckCircle and XCircle
import { useRouter, useCanGoBack } from '@tanstack/react-router'
import { useSchedulesById } from '@/hooks/useSchedules'
import { formatShiftTime } from '@/shared/schedule-details'
import { useEffect, useState } from 'react'
import { useStartVisit } from '@/hooks/useClockIn' // Assuming useClockIn also contains useStartVisit
import { useEndVisit } from '@/hooks/useClockOut' // Import useEndVisit from the same file or its dedicated file
import { getContext } from '@/integrations/tanstack-query/root-provider'

export default function ScheduleDetails() {
  const router = useRouter()
  const canGoBack = useCanGoBack()
  const { scheduleId } = Route.useParams()

  const { queryClient } = getContext()

  const [gettingLocation, setGettingLocation] = useState(false)
  const [mutationErrorMessage, setMutationErrorMessage] = useState<
    string | null
  >(null)

  const startVisitMutation = useStartVisit() // Renamed to avoid conflict
  const clockOutMutation = useEndVisit() // Initialize the clock-out mutation

  const { data, error, isLoading } = useSchedulesById(scheduleId)

  // Effect for handling Start Visit mutation success/error
  useEffect(() => {
    if (startVisitMutation.isSuccess) {
      setMutationErrorMessage(null)
      queryClient.invalidateQueries({ queryKey: ['schedules', scheduleId] })
      router.navigate({
        to: `/clock-out/$scheduleId`,
        params: { scheduleId },
      })
    }

    if (startVisitMutation.isError) {
      setMutationErrorMessage(
        startVisitMutation.error?.message || 'An unexpected error occurred.',
      )
    }
  }, [
    startVisitMutation.isSuccess,
    startVisitMutation.isError,
    startVisitMutation.error,
    router,
    scheduleId,
    queryClient,
  ])

  // Effect for handling Clock Out mutation success/error
  useEffect(() => {
    if (clockOutMutation.isSuccess) {
      setMutationErrorMessage(null)
      queryClient.invalidateQueries({ queryKey: ['schedules', scheduleId] })
      // Navigate to a report or dashboard page after successful clock-out
      router.navigate({
        to: `/`, // Example: navigate to a report page
      })
    }

    if (clockOutMutation.isError) {
      setMutationErrorMessage(
        clockOutMutation.error?.message || 'An unexpected error occurred.',
      )
    }
  }, [
    clockOutMutation.isSuccess,
    clockOutMutation.isError,
    clockOutMutation.error,
    router,
    scheduleId,
    queryClient,
  ])

  if (isLoading && !data) {
    return (
      <div className="flex items-center justify-center h-screen">
        Loading...
      </div>
    )
  }
  if (error && !data) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-red-500">Error loading data: {error.message}</div>
      </div>
    )
  }

  // After this block, data, data.success, and data.data are guaranteed to exist.
  if (!data || !data.success || !data.data) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-red-500">No schedule data found</div>
      </div>
    )
  }

  // formatShiftTime can safely be called here now that data.data is guaranteed
  formatShiftTime(data.data)

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
  }

  const handleClockIn = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    setMutationErrorMessage(null)
    setGettingLocation(true)

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords
          setGettingLocation(false)

          startVisitMutation.mutate({
            scheduleId: scheduleId,
            latitude: latitude,
            longitude: longitude,
          })
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

  // New handleClockOut function
  const handleClockOut = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault()
    setMutationErrorMessage(null)
    setGettingLocation(true) // Re-use gettingLocation state for clock-out

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords
          setGettingLocation(false)

          clockOutMutation.mutate({
            scheduleId: scheduleId,
            latitude: latitude,
            longitude: longitude,
          })
        },
        (error) => {
          console.error('Error getting user location for clock-out:', error)
          setGettingLocation(false)
          setMutationErrorMessage(
            'Could not get your current location for clock-out. Please enable location services.',
          )
        },
        { enableHighAccuracy: true, timeout: 5000, maximumAge: 0 },
      )
    } else {
      setGettingLocation(false)
      console.error(
        'Geolocation is not supported by this browser for clock-out.',
      )
      setMutationErrorMessage(
        'Geolocation is not supported by this browser for clock-out.',
      )
    }
  }

  const isClockInButtonDisabled =
    startVisitMutation.isPending || gettingLocation
  const isClockOutButtonDisabled = clockOutMutation.isPending || gettingLocation // Disable clock-out button when fetching location or mutation is pending

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col max-w-md mx-auto">
      <div className="bg-white px-6 py-4">
        <div className="flex items-center gap-4">
          {canGoBack ? (
            <button onClick={() => router.history.back()} className="p-1">
              <ArrowLeft className="h-6 w-6 text-gray-700" />
            </button>
          ) : null}
          <h1 className="text-lg font-semibold text-gray-900">
            Schedule Details
          </h1>
        </div>
      </div>

      <div className="flex-1 px-6 py-6 space-y-6">
        <div className="text-center space-y-4">
          <h2 className="text-lg font-semibold text-[#0D5D59]">Service A</h2>

          <div className="flex items-center justify-center gap-3">
            <Avatar className="h-12 w-12">
              <AvatarImage src={''} />
              <AvatarFallback className="bg-gray-200 text-gray-700">
                {getInitials(data.data.client_name)}
              </AvatarFallback>
            </Avatar>
            <div className="font-semibold text-gray-900 text-lg">
              {data.data.client_name}
            </div>
          </div>

          <div className="bg-blue-50 rounded-lg p-3">
            <div className="flex items-center justify-center gap-4 text-sm text-gray-600">
              <div className="flex items-center gap-1">
                <Calendar className="h-4 w-4 text-cyan-600" />
                <span>{data.data.shiftDate}</span>
              </div>
              <span>|</span>
              <div className="flex items-center gap-1">
                <Clock className="h-4 w-4 text-cyan-600" />
                <span>{data.data.shiftTimeOnly}</span>
              </div>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900">Client Contact:</h3>
          <div className="space-y-3">
            <div className="flex items-center gap-3">
              <Mail className="h-5 w-5 text-gray-600" />
              <span className="text-gray-700">placeholder@mail.com</span>
            </div>
            <div className="flex items-center gap-3">
              <Phone className="h-5 w-5 text-gray-600" />
              <span className="text-gray-700">088888888</span>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900">Address:</h3>
          <div className="flex items-start gap-3">
            <MapPin className="h-5 w-5 text-gray-600 mt-0.5" />
            <span className="text-gray-700">{data.data.location}</span>
          </div>
        </div>

        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900">Tasks:</h3>
          <div className="space-y-4">
            {data.data.tasks && data.data.tasks.length > 0 ? (
              data.data.tasks.map((task) => (
                <div key={task.id} className="space-y-2">
                  <div className="flex items-center gap-2">
                    {' '}
                    {/* Added flex container for icon and text */}
                    {task.status === 'completed' && (
                      <CheckCircle className="h-5 w-5 text-green-500" />
                    )}
                    {task.status === 'not_completed' && (
                      <XCircle className="h-5 w-5 text-red-500" />
                    )}
                    <h4 className="font-semibold text-[#0D5D59]">
                      {task.description}
                    </h4>
                  </div>
                  <p className="text-sm text-gray-600 leading-relaxed">
                    {task.description}
                  </p>
                </div>
              ))
            ) : (
              <p className="text-sm text-gray-600">
                No tasks defined for this schedule.
              </p>
            )}
          </div>
        </div>

        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900">Service Notes</h3>
          <p className="text-sm text-gray-600 leading-relaxed">service notes</p>
        </div>
      </div>

      <div className="p-6">
        {data.data.status === 'upcoming' && (
          <>
            <Button
              className="w-full bg-[#0D5D59] hover:bg-[#0D5D59]/90 py-4 text-lg font-semibold"
              onClick={handleClockIn}
              disabled={isClockInButtonDisabled}
            >
              {gettingLocation
                ? 'Getting Location...'
                : startVisitMutation.isPending
                  ? 'Clocking In...'
                  : 'Clock-In Now'}
            </Button>
            {mutationErrorMessage && (
              <p className="text-red-500 text-sm mt-2 text-center">
                Error: {mutationErrorMessage}
              </p>
            )}
          </>
        )}
        {data.data.status === 'in-progress' && (
          <>
            <Button
              className="w-full bg-[#0D5D59] hover:bg-[#0D5D59]/90 py-4 text-lg font-semibold"
              onClick={handleClockOut} // Use the new handleClockOut function
              disabled={isClockOutButtonDisabled}
            >
              {gettingLocation
                ? 'Getting Location...'
                : clockOutMutation.isPending
                  ? 'Clocking Out...'
                  : 'Clock-Out Now'}
            </Button>
            {mutationErrorMessage && (
              <p className="text-red-500 text-sm mt-2 text-center">
                Error: {mutationErrorMessage}
              </p>
            )}
          </>
        )}
      </div>
    </div>
  )
}
