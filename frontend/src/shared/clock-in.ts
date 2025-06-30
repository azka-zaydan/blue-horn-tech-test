export type StartVisitRequest = {
  latitude: number
  longitude: number
}

export type StartVisitResponse = {
  success: boolean
  message: string
  error?: {
    code: number
    message: string
    details: string
  }
}

export type StartVisitError = {
  success: false
  message: string
  error: {
    code: number
    message: string
    details: string
  }
}
