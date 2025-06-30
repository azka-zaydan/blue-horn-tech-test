// src/hooks/useSchedules.ts
import type { SchedulesResponse } from '@/shared/dashboard' // Import types
import type { ScheduleDetailsResponse } from '@/shared/schedule-details'
import {
  type InfiniteData,
  useInfiniteQuery,
  useQuery,
} from '@tanstack/react-query' // Import InfiniteData

// Utility function to fetch paginated schedules
const fetchSchedules = async (
  page: number,
  pageSize: number,
): Promise<SchedulesResponse> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/schedules?page=${page}&pageSize=${pageSize}`,
  )

  if (!response.ok) {
    const errorData = await response
      .json()
      .catch(() => ({ message: 'Failed to parse error response.' }))
    throw new Error(
      errorData.message ||
        `Network response was not ok, status: ${response.status}`,
    )
  }
  return response.json() as Promise<SchedulesResponse>
}

// Custom hook for infinite scrolling schedules
export const useInfiniteSchedules = (pageSize: number = 5) => {
  // Corrected generic type parameters for useInfiniteQuery
  return useInfiniteQuery<
    SchedulesResponse, // TQueryFnData: The type of data returned by queryFn for a single page
    Error, // TError: The type of error
    InfiniteData<SchedulesResponse, number>, // TData: The type of the aggregated data returned by the hook
    string[], // TQueryKey: The type of the queryKey
    number // TPageParam: The type of the pageParam
  >({
    queryKey: ['schedules'],
    queryFn: async ({ pageParam = 1 }) => fetchSchedules(pageParam, pageSize),
    getNextPageParam: (lastPage) => {
      if (lastPage.success && lastPage.pagination) {
        const nextPage = lastPage.pagination.page + 1
        return nextPage <= lastPage.pagination.total_pages
          ? nextPage
          : undefined
      }
      return undefined
    },
    initialPageParam: 1,
    refetchOnWindowFocus: true,
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: 2, // Retry failed requests up to 2 times
  })
}

// Utility function to fetch a single schedule by ID
const fetchSchedulesById = async (
  scheduleId: string,
): Promise<ScheduleDetailsResponse> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/schedules/${scheduleId}`,
  )

  if (!response.ok) {
    const errorData = await response
      .json()
      .catch(() => ({ message: 'Failed to parse error response.' }))
    throw new Error(
      errorData.message ||
        `Network response was not ok for schedule ${scheduleId}, status: ${response.status}`,
    )
  }
  return response.json() as Promise<ScheduleDetailsResponse>
}

// Custom hook to fetch a single schedule by ID
export const useSchedulesById = (scheduleId: string) => {
  return useQuery<ScheduleDetailsResponse, Error>({
    queryKey: ['schedules', scheduleId],
    queryFn: () => fetchSchedulesById(scheduleId),
    refetchOnWindowFocus: true,
    staleTime: 1000 * 60 * 5, // 5 minutes
    retry: 2,
    enabled: !!scheduleId, // Only run query if scheduleId is available
  })
}
