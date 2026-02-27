import { useToast } from '@/components/ui/toast'

export enum ErrorType {
  NETWORK = 'network',
  VALIDATION = 'validation',
  AUTHENTICATION = 'authentication',
  AUTHORIZATION = 'authorization',
  SERVER = 'server',
  UNKNOWN = 'unknown',
}

export interface AppError {
  type: ErrorType
  message: string
  originalError?: any
}

/** Базовый обработчик ошибок */
export function handleError(error: any): AppError {
  const { toast } = useToast()
  let appError: AppError

  if (error.name === 'HTTPError') {
    const status = error.response?.status

    if (status === 401 || status === 403) {
      appError = {
        type: ErrorType.AUTHENTICATION,
        message: 'Ошибка аутентификации. Пожалуйста, войдите снова.',
        originalError: error,
      }
    }
    else if (status === 400) {
      appError = {
        type: ErrorType.VALIDATION,
        message: 'Проверьте правильность введенных данных',
        originalError: error,
      }
    }
    else if (status >= 500) {
      appError = {
        type: ErrorType.SERVER,
        message: 'Ошибка сервера. Пожалуйста, попробуйте позже.',
        originalError: error,
      }
    }
    else {
      appError = {
        type: ErrorType.UNKNOWN,
        message: 'Произошла неизвестная ошибка',
        originalError: error,
      }
    }
  }
  else if (error.name === 'NetworkError') {
    appError = {
      type: ErrorType.NETWORK,
      message: 'Проблемы с сетевым подключением',
      originalError: error,
    }
  }
  else {
    appError = {
      type: ErrorType.UNKNOWN,
      message: error.message || 'Произошла неизвестная ошибка',
      originalError: error,
    }
  }

  toast({
    title: 'Ошибка',
    description: appError.message,
    variant: 'destructive',
  })

  console.error('[App Error]', appError)

  return appError
}
