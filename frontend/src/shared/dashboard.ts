export type Pagination = {
  page: number
  page_size: number
  total_items: number
  total_pages: number
}

export type ScheduleStatus =
  | 'completed'
  | 'in-progress'
  | 'upcoming'
  | 'cancelled'

export type Schedule = {
  id: string
  client_name: string
  shiftDate: string
  shiftTimeOnly: string
  shiftTimezone: string
  shiftDateWithDay: string
  shift_time: string
  location: string
  status: ScheduleStatus
  start_time: string | null
  start_latitude: number | null
  start_longitude: number | null
  end_time: string | null
  end_latitude: number | null
  end_longitude: number | null
  created_at: string
  updated_at: string
}

type SuccessResponseWithData = {
  success: true
  message: string
  data: Schedule[] | null // Can be an array of schedules or null
  pagination: Pagination
}

type SuccessResponseNoData = {
  success: true
  message: string
  data: null // Explicitly null when no data
  pagination: Pagination
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

export type SchedulesResponse =
  | SuccessResponseWithData
  | SuccessResponseNoData
  | ErrorResponse

// Function to parse shift_time into date, time, and timezone
export const parseShiftTime = (shiftTime: string) => {
  const parsedDate = new Date(shiftTime)

  // Extract date part (YYYY-MM-DD)
  const date = parsedDate.toISOString().split('T')[0] // '2025-06-28'

  // Extract time part (HH:mm:ss)
  const time = parsedDate.toISOString().split('T')[1].split('.')[0] // '09:00:49'

  // Extract timezone information
  const timezoneOffset = parsedDate.getTimezoneOffset()
  const timezone = `UTC${timezoneOffset > 0 ? '-' : '+'}${Math.abs(timezoneOffset) / 60}` // UTC+X or UTC-X

  // Format the date to "Mon, 15 Jan 2025" format
  const formattedDate = new Intl.DateTimeFormat('en-GB', {
    weekday: 'short', // Abbreviated weekday (e.g., "Mon")
    day: '2-digit', // Two-digit day (e.g., "15")
    month: 'short', // Abbreviated month (e.g., "Jan")
    year: 'numeric', // Full year (e.g., "2025")
  }).format(parsedDate)

  return { date, time, timezone, formattedDate }
}

export const getActiveSchedule = (schedules: Schedule[] | null) => {
  // Filter for 'in-progress' and 'upcoming' schedules\
  if (!schedules || schedules.length === 0) {
    return null // Return null if no schedules are provided
  }
  const activeSchedules = schedules.filter(
    (schedule) =>
      schedule.status === 'in-progress' || schedule.status === 'upcoming',
  )

  // Sort 'in-progress' schedules by the latest shift_time
  const inProgressSchedules = activeSchedules.filter(
    (schedule) => schedule.status === 'in-progress',
  )

  if (inProgressSchedules.length > 0) {
    // Sort 'in-progress' schedules by latest shift_time (descending order)
    inProgressSchedules.sort((a, b) => {
      const dateA = new Date(a.shift_time)
      const dateB = new Date(b.shift_time)
      return dateB.getTime() - dateA.getTime() // Latest first
    })

    // Return the most recent 'in-progress' schedule
    return inProgressSchedules[0]
  }

  // If there are no 'in-progress' schedules, return the first 'upcoming' schedule
  const upcomingSchedules = activeSchedules.filter(
    (schedule) => schedule.status === 'upcoming',
  )

  if (upcomingSchedules.length > 0) {
    // Return the first 'upcoming' schedule (no need to sort as we just want the first one)
    return upcomingSchedules[0]
  }

  // If no active schedule found, return null
  return null
}

// Updated calculateScheduleStats function to use the new properties
export const calculateScheduleStats = (schedules: Schedule[] | null) => {
  const stats = {
    missed: 0,
    upcoming: 0,
    completed: 0,
  }

  if (!schedules) {
    return stats // Return empty stats if no schedules are provided
  }

  schedules.forEach((schedule) => {
    // Parse shift_time for date, time, and timezone
    const { date, time, timezone, formattedDate } = parseShiftTime(
      schedule.shift_time,
    )
    schedule.shiftDate = date
    schedule.shiftTimeOnly = time
    schedule.shiftTimezone = timezone
    schedule.shiftDateWithDay = formattedDate

    if (schedule.status === 'completed') {
      stats.completed += 1
    } else if (schedule.status === 'upcoming') {
      stats.upcoming += 1
    } else if (schedule.status === 'in-progress') {
      // Assuming that "missed" means a shift that was scheduled but has passed
      // and is not marked as completed or upcoming
      if (new Date(schedule.shift_time) < new Date()) {
        stats.missed += 1
      }
    }
  })

  return stats
}

export const formatDuration = (startTime: string) => {
  if (!startTime) {
    return '00:00:00' // Return default duration if startTime is not provided
  }
  const start = new Date(startTime) // Convert startTime string to Date object
  const now = new Date() // Current date and time

  // Calculate the difference in milliseconds
  const diff = now.getTime() - start.getTime()

  // Convert milliseconds to hours, minutes, and seconds
  const hours = Math.floor(diff / (1000 * 60 * 60)) // Convert to hours
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60)) // Convert to minutes
  const seconds = Math.floor((diff % (1000 * 60)) / 1000) // Convert to seconds

  // Format the duration as "hh:mm:ss"
  const formattedDuration = `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`

  return formattedDuration
}
