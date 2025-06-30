// Define the types for the clock-out request and response
export type ClockOutRequest = {
  latitude: number
  longitude: number
}

export type ClockOutResponse = {
  success: boolean
  message: string
  error?: {
    code: number
    message: string
    details: string
  }
}

export type ClockOutError = {
  success: false
  message: string
  error: {
    code: number
    message: string
    details: string
  }
}
