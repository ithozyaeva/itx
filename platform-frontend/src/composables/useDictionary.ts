import type { DictionaryItem } from '@/services/dictionaryService'
import type { ComputedRef } from 'vue'

import { dictionaryService } from '@/services/dictionaryService'
import { useQuery } from '@tanstack/vue-query'

import { computed } from 'vue'

export type DictionaryKey =
  | 'placeTypes'
  | 'memberRoles'
  | 'reviewStatuses'
  | 'grades'
  | 'referalLinkStatuses'

export const dictionaryKeys = {
  all: ['dictionaries'] as const,
  lists: () => [...dictionaryKeys.all, 'list'] as const,
  list: (key: DictionaryKey) => [...dictionaryKeys.lists(), key] as const,
}

export function useDictionary<TValue extends string, K extends DictionaryKey = DictionaryKey>(keys?: K[]) {
  const { data: dictionaries, isLoading, error, refetch } = useQuery({
    queryKey: dictionaryKeys.lists(),
    queryFn: () => dictionaryService.getAllDictionaries(),
    staleTime: Infinity,
    gcTime: Infinity,
  })

  const getDictionary = (key: DictionaryKey): DictionaryItem<TValue>[] => {
    return (dictionaries.value?.[key] || []) as DictionaryItem<TValue>[]
  }

  const computedDictionaries = {} as Record<K, ComputedRef<DictionaryItem<TValue>[]>>
  const computedDictionariesObject = {} as Record<`${K}Object`, ComputedRef<Record<TValue, string>>>

  if (keys && keys.length > 0) {
    keys.forEach((key) => {
      computedDictionaries[key] = computed(() => getDictionary(key))
      computedDictionariesObject[`${key}Object`] = computed(() =>
        Object.fromEntries(getDictionary(key).map(d => [d.value, d.label])) as Record<TValue, string>,
      )
    })
  }

  return {
    dictionaries,
    isLoading,
    error,
    refetch,
    getDictionary,
    ...computedDictionariesObject,
    ...computedDictionaries,
  }
}
