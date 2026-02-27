import type { RemovableRef } from '@vueuse/core'
import { useLocalStorage } from '@vueuse/core'

export function useToken(): RemovableRef<null | string> {
  return useLocalStorage<string>('tg_token', null)
}
