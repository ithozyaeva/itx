import ky from 'ky'
import router from '@/router'
import { logout } from '@/services/authService'
import { handleError } from '@/services/errorService'

const api = ky.create({
  prefixUrl: '/api/admin/',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
  hooks: {
    beforeRequest: [
      (request) => {
        const tgToken = localStorage.getItem('tg_token')

        if (tgToken) {
          request.headers.set('X-Telegram-User-Token', tgToken)
        }
      },
    ],
    afterResponse: [
      async (_request, _options, response) => {
        if (response.ok)
          return response

        if (response.status === 401) {
          try {
            handleError(new Error('Сессия истекла'))
            logout()
            router.push('/login')
          }
          catch (error) {
            handleError(error)
          }
        }

        return response
      },
    ],
  },
})

export default api
