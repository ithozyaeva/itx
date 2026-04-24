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
  originalError?: unknown
}

const fallbackMessages: Record<string, string> = {
  401: 'Ошибка аутентификации. Пожалуйста, войдите снова.',
  403: 'Ошибка аутентификации. Пожалуйста, войдите снова.',
  500: 'Ошибка сервера. Пожалуйста, попробуйте позже.',
}

/** Базовый обработчик ошибок */
export async function handleError(error: unknown): Promise<AppError> {
  const { toast } = useToast()
  let appError: AppError

  const err = error as Record<string, unknown>

  if (err?.name === 'HTTPError') {
    const response = err.response as { status?: number, json?: () => Promise<unknown> } | undefined
    const status = response?.status
    let serverMessage: string | undefined

    try {
      const body = await response?.json?.() as Record<string, unknown> | undefined
      serverMessage = body?.error as string | undefined
    }
    catch {}

    if (status === 401 || status === 403) {
      appError = {
        type: ErrorType.AUTHENTICATION,
        message: serverMessage || fallbackMessages[401],
        originalError: error,
      }
    }
    else if (status === 400 || status === 422 || status === 429) {
      appError = {
        type: ErrorType.VALIDATION,
        message: serverMessage || 'Проверьте правильность введенных данных',
        originalError: error,
      }
    }
    else if (status >= 500) {
      appError = {
        type: ErrorType.SERVER,
        message: serverMessage || fallbackMessages[500],
        originalError: error,
      }
    }
    else {
      appError = {
        type: ErrorType.UNKNOWN,
        message: serverMessage || 'Произошла неизвестная ошибка',
        originalError: error,
      }
    }
  }
  else if (err?.name === 'NetworkError') {
    appError = {
      type: ErrorType.NETWORK,
      message: 'Проблемы с сетевым подключением',
      originalError: error,
    }
  }
  else {
    appError = {
      type: ErrorType.UNKNOWN,
      message: (err?.message as string) || 'Произошла неизвестная ошибка',
      originalError: error,
    }
  }

  toast({
    title: 'Ошибка',
    description: appError.message,
    variant: 'destructive',
  })

  if (import.meta.env.DEV) {
    console.error('[App Error]', appError)
  }

  return appError
}
