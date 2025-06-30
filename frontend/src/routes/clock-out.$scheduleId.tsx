// src/pages/ClockOut.tsx
import { createFileRoute } from '@tanstack/react-router'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { ArrowLeft, Check, X, MapPin, Save } from 'lucide-react'
import { useRouter, useCanGoBack, useNavigate } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { type Task } from '@/shared/schedule-details'
import { formatDuration } from '@/shared/dashboard'
import { useSchedulesById } from '@/hooks/useSchedules'
import { useUpdateTask } from '@/hooks/useTasks'
import { getContext } from '@/integrations/tanstack-query/root-provider'
import { useEndVisit } from '@/hooks/useClockOut' // Import useEndVisit

export const Route = createFileRoute('/clock-out/$scheduleId')({
  component: ClockOut,
})

function TaskItem({
  task,
  onToggle,
  onReasonChange,
  onSaveReason,
  isUpdating,
}: {
  task: Task
  onToggle: (id: string) => void
  onReasonChange: (id: string, reason: string) => void
  onSaveReason: (id: string, reason: string) => void
  isUpdating: boolean
}) {
  return (
    <div className="space-y-3">
      <h4 className="font-semibold text-[#0D5D59]">{task.description}</h4>
      <p className="text-sm text-gray-600 leading-relaxed">
        {task.description}
      </p>
      <div className="flex items-center gap-4">
        <button
          onClick={() => onToggle(task.id)}
          className={`flex items-center gap-1 ${
            task.status === 'completed' ? 'text-green-600' : 'text-gray-500'
          }`}
          disabled={isUpdating}
        >
          <Check className="h-4 w-4" />
          <span className="text-sm font-medium">Yes</span>
        </button>

        <span className="text-gray-400">|</span>

        <button
          onClick={() => onToggle(task.id)}
          className={`flex items-center gap-1 ${
            task.status !== 'completed' ? 'text-red-600' : 'text-gray-500'
          }`}
          disabled={isUpdating}
        >
          <X className="h-4 w-4" />
          <span className="text-sm font-medium">No</span>
        </button>
      </div>
      {task.status !== 'completed' && (
        <div className="flex items-center gap-2">
          <input
            type="text"
            placeholder="Add reason..."
            value={task.reason || ''}
            onChange={(e) => onReasonChange(task.id, e.target.value)}
            className="flex-1 px-4 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#0D5D59] focus:border-transparent"
            disabled={isUpdating}
          />
          <Button
            size="sm"
            className="bg-[#0D5D59] hover:bg-[#0D5D59]/90"
            onClick={() => onSaveReason(task.id, task.reason || '')}
            disabled={isUpdating || !task.reason}
          >
            <Save className="h-4 w-4" />
          </Button>
        </div>
      )}
      {isUpdating && <p className="text-xs text-blue-500 mt-1">Updating...</p>}
    </div>
  )
}

export default function ClockOut() {
  const router = useRouter()
  const canGoBack = useCanGoBack()
  const navigate = useNavigate()
  const [duration, setDuration] = useState<string>('')
  const [tasks, setTasks] = useState<Task[]>([])
  const { scheduleId } = Route.useParams()

  const { queryClient } = getContext()

  const { data, error, isLoading } = useSchedulesById(scheduleId)
  const updateTaskMutation = useUpdateTask()
  const endVisitMutation = useEndVisit() // Initialize the end visit mutation

  // State to track which task is currently being updated
  const [updatingTaskId, setUpdatingTaskId] = useState<string | null>(null)
  // State for overall clock-out mutation error message
  const [clockOutErrorMessage, setClockOutErrorMessage] = useState<
    string | null
  >(null)
  // State to track if location is being fetched for clock-out
  const [gettingLocationForClockOut, setGettingLocationForClockOut] =
    useState(false)

  useEffect(() => {
    if (isLoading || !data || !data.success || !data.data) {
      return
    }

    if (data && data.success && data.data && data.data.start_time) {
      const interval = setInterval(() => {
        setDuration(formatDuration(String(data.data.start_time)))
      }, 1000)

      return () => clearInterval(interval)
    }

    const initialTasks = data.data.tasks.map((task) => ({
      ...task,
      status: task.status || 'not_completed',
      reason: task.reason || '',
    }))
    setTasks(initialTasks)
  }, [data, isLoading])

  // Effect to handle updateTaskMutation success/error
  useEffect(() => {
    if (updateTaskMutation.isSuccess) {
      setUpdatingTaskId(null)
      queryClient.invalidateQueries({ queryKey: ['schedules', scheduleId] })
    }
    if (updateTaskMutation.isError) {
      setUpdatingTaskId(null)
      console.error('Failed to update task:', updateTaskMutation.error)
      // Optionally display a user-friendly error for task update
    }
  }, [
    updateTaskMutation.isSuccess,
    updateTaskMutation.isError,
    queryClient,
    scheduleId,
  ])

  // Effect to handle End Visit mutation success/error
  useEffect(() => {
    if (endVisitMutation.isSuccess) {
      setClockOutErrorMessage(null)
      queryClient.invalidateQueries({ queryKey: ['schedules', scheduleId] })
      router.navigate({ to: '/' }) // Navigate to home or dashboard after successful clock-out
    }
    if (endVisitMutation.isError) {
      setClockOutErrorMessage(
        endVisitMutation.error?.message ||
          'An unexpected error occurred during clock-out.',
      )
      console.error('Failed to clock out:', endVisitMutation.error)
    }
  }, [
    endVisitMutation.isSuccess,
    endVisitMutation.isError,
    endVisitMutation.error,
    router,
    queryClient,
    scheduleId,
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

  if (!data || !data.success || !data.data) {
    return (
      <div className="flex items-center justify-center h-screen">
        <div className="text-red-500">No schedule data found</div>
      </div>
    )
  }

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
  }

  const handleTaskToggle = (taskId: string) => {
    setTasks((prev) =>
      prev.map((task) => {
        if (task.id === taskId) {
          const newStatus =
            task.status === 'completed' ? 'not_completed' : 'completed'
          const newReason = newStatus === 'completed' ? null : task.reason

          setUpdatingTaskId(taskId)
          updateTaskMutation.mutate({
            taskId: taskId,
            status: newStatus,
            reason: newReason,
          })

          return {
            ...task,
            status: newStatus,
            reason: newReason || '',
          }
        }
        return task
      }),
    )
  }

  const handleReasonChange = (taskId: string, reason: string) => {
    setTasks((prev) =>
      prev.map((task) =>
        task.id === taskId ? { ...task, reason: reason } : task,
      ),
    )
  }

  const handleSaveReason = (taskId: string, reason: string) => {
    setUpdatingTaskId(taskId)
    updateTaskMutation.mutate({
      taskId: taskId,
      status: tasks.find((t) => t.id === taskId)?.status || 'pending',
      reason: reason,
    })
  }

  const handleClockOut = () => {
    setClockOutErrorMessage(null)
    setGettingLocationForClockOut(true) // Indicate location fetching for clock-out

    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const { latitude, longitude } = position.coords
          setGettingLocationForClockOut(false)

          endVisitMutation.mutate({
            scheduleId: scheduleId,
            latitude: latitude,
            longitude: longitude,
          })
        },
        (error) => {
          console.error('Error getting user location for clock-out:', error)
          setGettingLocationForClockOut(false)
          setClockOutErrorMessage(
            'Could not get your current location for clock-out. Please enable location services.',
          )
        },
        { enableHighAccuracy: true, timeout: 5000, maximumAge: 0 },
      )
    } else {
      setGettingLocationForClockOut(false)
      console.error(
        'Geolocation is not supported by this browser for clock-out.',
      )
      setClockOutErrorMessage(
        'Geolocation is not supported by this browser for clock-out.',
      )
    }
  }

  const isClockOutButtonDisabled =
    endVisitMutation.isPending ||
    gettingLocationForClockOut ||
    updateTaskMutation.isPending

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col max-w-md mx-auto">
      <div className="bg-white px-6 py-4">
        <div className="flex items-center gap-4">
          {canGoBack ? (
            <button onClick={() => router.history.back()} className="p-1">
              <ArrowLeft className="h-6 w-6 text-gray-700" />
            </button>
          ) : null}
          <h1 className="text-lg font-semibold text-gray-900">Clock-Out</h1>
        </div>
      </div>

      <div className="flex-1 px-6 py-6 space-y-6">
        <div className="text-center">
          <div className="text-3xl font-bold tracking-wider text-gray-900">
            {duration}
          </div>
        </div>

        <div className="text-center space-y-4">
          <h2 className="text-lg font-semibold text-[#0D5D59]">
            Service Name A
          </h2>

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
        </div>

        <div className="space-y-4">
          <div>
            <h3 className="font-semibold text-gray-900 mb-2">Tasks:</h3>
            <p className="text-sm text-gray-600 mb-4">
              Please tick the tasks that you have done
            </p>
          </div>

          <div className="space-y-6">
            {tasks &&
              tasks.map((task) => (
                <TaskItem
                  key={task.id}
                  task={task}
                  onToggle={handleTaskToggle}
                  onReasonChange={handleReasonChange}
                  onSaveReason={handleSaveReason}
                  isUpdating={
                    updatingTaskId === task.id || endVisitMutation.isPending
                  } // Disable if any task is updating or clock-out is pending
                />
              ))}
          </div>
        </div>
        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900">Clock-In Location</h3>
          <div className="flex items-start gap-3">
            <div className="w-16 h-16 bg-gray-200 rounded-lg flex items-center justify-center">
              <MapPin className="h-6 w-6 text-gray-500" />
            </div>
            <div className="flex-1">
              <p className="text-sm text-gray-700">{data.data.location}</p>
              <p className="text-xs text-gray-500">
                Latitude: {data.data.start_latitude}, Longitude:{' '}
                {data.data.start_longitude}
              </p>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <h3 className="font-semibold text-gray-900">Service Notes:</h3>
          <p className="text-sm text-gray-600 leading-relaxed">Service Notes</p>
        </div>
      </div>

      <div className="p-6 space-y-3">
        <div className="flex gap-3">
          <Button
            variant="outline"
            className="flex-1 border-red-300 text-red-600 hover:bg-red-50"
            onClick={() => navigate({ to: '/' })}
            disabled={endVisitMutation.isPending} // Disable cancel if clock-out is in progress
          >
            Cancel Clock-In
          </Button>
          <Button
            className="flex-1 bg-[#0D5D59] hover:bg-[#0D5D59]/90"
            onClick={handleClockOut}
            disabled={isClockOutButtonDisabled}
          >
            {gettingLocationForClockOut
              ? 'Getting Location...'
              : endVisitMutation.isPending
                ? 'Clocking Out...'
                : 'Clock-Out'}
          </Button>
        </div>
        {clockOutErrorMessage && (
          <p className="text-red-500 text-sm mt-2 text-center">
            Error: {clockOutErrorMessage}
          </p>
        )}
      </div>
    </div>
  )
}
