import type { StartVisitRequest, StartVisitResponse } from '@/shared/clock-in'
import { useMutation } from '@tanstack/react-query'
// Hapus import axios, karena kita tidak lagi menggunakannya
// import axios from 'axios'

export const startVisit = async (
  scheduleId: string,
  data: StartVisitRequest,
): Promise<StartVisitResponse> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/schedules/${scheduleId}/start`,
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

  const responseData: StartVisitResponse = await response.json()
  return responseData
}

export const useStartVisit = () => {
  return useMutation({
    mutationFn: ({
      scheduleId,
      latitude,
      longitude,
    }: {
      scheduleId: string
      latitude: number
      longitude: number
    }) => startVisit(scheduleId, { latitude, longitude }),
    onError: (error) => {
      console.error('Error starting visit:', error.message)
    },
  })
}
