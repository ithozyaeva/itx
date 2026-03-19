import type { DictionaryKey } from '@/composables/useDictionary'
import { describe, expect, it } from 'vitest'
import { dictionaryKeys } from '@/composables/useDictionary'

describe('useDictionary', () => {
  describe('dictionaryKeys', () => {
    it('generates correct "all" key', () => {
      expect(dictionaryKeys.all).toEqual(['dictionaries'])
    })

    it('generates correct "lists" key', () => {
      expect(dictionaryKeys.lists()).toEqual(['dictionaries', 'list'])
    })

    it('generates correct "list" key for a specific dictionary', () => {
      expect(dictionaryKeys.list('placeTypes')).toEqual(['dictionaries', 'list', 'placeTypes'])
      expect(dictionaryKeys.list('memberRoles')).toEqual(['dictionaries', 'list', 'memberRoles'])
      expect(dictionaryKeys.list('reviewStatuses')).toEqual(['dictionaries', 'list', 'reviewStatuses'])
      expect(dictionaryKeys.list('grades')).toEqual(['dictionaries', 'list', 'grades'])
      expect(dictionaryKeys.list('referalLinkStatuses')).toEqual(['dictionaries', 'list', 'referalLinkStatuses'])
    })

    it('returns new array instances on each call', () => {
      const lists1 = dictionaryKeys.lists()
      const lists2 = dictionaryKeys.lists()
      expect(lists1).toEqual(lists2)
      expect(lists1).not.toBe(lists2)
    })
  })

  describe('DictionaryKey type', () => {
    it('accepts valid dictionary keys', () => {
      const validKeys: DictionaryKey[] = [
        'placeTypes',
        'memberRoles',
        'reviewStatuses',
        'grades',
        'referalLinkStatuses',
      ]
      expect(validKeys).toHaveLength(5)
    })
  })
})
