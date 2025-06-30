import type { ClockOutRequest, ClockOutResponse } from '@/shared/clock-out'
import { useMutation } from '@tanstack/react-query'

/**
 * Performs the API call to clock out of a schedule.
 * @param scheduleId The ID of the schedule to clock out from.
 * @param data The clock-out request data, including latitude and longitude.
 * @returns A promise that resolves to the ClockOutResponse.
 * @throws An error if the API call fails or returns a non-OK status.
 */
export const endVisit = async (
  scheduleId: string,
  data: ClockOutRequest,
): Promise<ClockOutResponse> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/schedules/${scheduleId}/end`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data),
    },
  )

  if (!response.ok) {
    const errorData = await response
      .json()
      .catch(() => ({ message: 'Failed to parse error response.' }))

    throw new Error(
      errorData.message || `HTTP error! status: ${response.status}`,
    )
  }

  const responseData: ClockOutResponse = await response.json()
  return responseData
}

/**
 * A React Query mutation hook for clocking out of a visit.
 * It provides a convenient way to manage the state of the clock-out operation
 * (loading, error, success).
 */
export const useEndVisit = () => {
  return useMutation({
    mutationFn: ({
      scheduleId,
      latitude,
      longitude,
    }: {
      scheduleId: string
      latitude: number
      longitude: number
    }) => endVisit(scheduleId, { latitude, longitude }),

    onError: (error) => {
      console.error('Error clocking out:', error.message)
    },
  })
}
