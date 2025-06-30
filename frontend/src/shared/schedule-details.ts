// Schedule with tasks data type
export type Task = {
  id: string
  schedule_id: string
  description: string
  status: 'completed' | 'in-progress' | 'pending' | 'not_completed'
  reason?: string // Optional, required if status is 'not_completed'
  created_at: string
  updated_at: string
}

export type ScheduleDetails = {
  id: string
  client_name: string
  shift_time: string
  shiftDate: string
  shiftTimeOnly: string
  shiftTimezone: string
  shiftDateWithDay: string
  location: string
  status: 'completed' | 'in-progress' | 'upcoming' | 'cancelled'
  start_time: string | null
  start_latitude: number | null
  start_longitude: number | null
  end_time: string | null
  end_latitude: number | null
  end_longitude: number | null
  created_at: string
  updated_at: string
  tasks: Task[]
}

type SuccessResponseWithData = {
  success: true
  message: string
  data: ScheduleDetails
}

type ErrorResponse = {
  success: false
  message: string
  error: {
    code: number
    message: string
    details: string
  }
}

export type ScheduleDetailsResponse = SuccessResponseWithData | ErrorResponse

export const formatShiftTime = (sched: ScheduleDetails) => {
  const date = new Date(sched.shift_time) // Convert ISO string to Date object

  // Format the date to "Mon, 15 Jan 2025" format
  const formattedDate = new Intl.DateTimeFormat('en-GB', {
    weekday: 'short', // Abbreviated weekday (e.g., "Mon")
    day: '2-digit', // Two-digit day (e.g., "15")
    month: 'short', // Abbreviated month (e.g., "Jan")
    year: 'numeric', // Full year (e.g., "2025")
  }).format(date)

  const parsedDate = new Date(sched.shift_time)

  // Extract time part (HH:mm:ss)
  const time = parsedDate.toISOString().split('T')[1].split('.')[0] // '09:00:49'

  // Extract timezone information
  const timezoneOffset = parsedDate.getTimezoneOffset()
  const timezone = `UTC${timezoneOffset > 0 ? '-' : '+'}${Math.abs(timezoneOffset) / 60}` // UTC+X or UTC-X

  sched.shiftDate = formattedDate
  sched.shiftTimeOnly = time
  sched.shiftTimezone = timezone
}

// Define types for the update task request and response
export type UpdateTaskRequest = {
  status:
    | 'completed'
    | 'in-progress'
    | 'pending'
    | 'not_completed'
    | 'cancelled' // Expanded status types
  reason?: string | null // Reason can be string or null
}

export type UpdateTaskResponse = {
  success: boolean
  message: string
  data?: Task // Optional: return the updated task data
  error?: {
    code: number
    message: string
    details: string
  }
}
