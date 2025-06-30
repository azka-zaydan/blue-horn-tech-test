import type {
  UpdateTaskRequest,
  UpdateTaskResponse,
} from '@/shared/schedule-details'
import { useMutation } from '@tanstack/react-query'

/**
 * Performs the API call to update a single task.
 * @param taskId The ID of the task to update.
 * @param data The update task request data (status and optional reason).
 * @returns A promise that resolves to the UpdateTaskResponse.
 * @throws An error if the API call fails or returns a non-OK status.
 */
export const updateTask = async (
  taskId: string,
  data: UpdateTaskRequest,
): Promise<UpdateTaskResponse> => {
  const response = await fetch(
    `${import.meta.env.VITE_API_URL}/tasks/${taskId}/update`,
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

  const responseData: UpdateTaskResponse = await response.json()
  return responseData
}

/**
 * A React Query mutation hook for updating a single task.
 * It provides a convenient way to manage the state of the update operation
 * (loading, error, success).
 */
export const useUpdateTask = () => {
  return useMutation<
    UpdateTaskResponse,
    Error,
    {
      taskId: string
      status: UpdateTaskRequest['status']
      reason?: UpdateTaskRequest['reason']
    }
  >({
    mutationFn: ({ taskId, status, reason }) =>
      updateTask(taskId, { status, reason }),
    onError: (error) => {
      console.error('Error updating task:', error.message)
    },
  })
}
