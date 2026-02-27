import ky from 'ky'

export interface DictionaryItem<TValue extends string> {
  value: TValue
  label: string
}

export interface DictionaryMap {
  [key: string]: DictionaryItem<string>[]
}

// Create a separate API instance for public endpoints (without authentication)
const publicApi = ky.create({
  prefixUrl: '/api/',
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
})

export class DictionaryService {
  async getAllDictionaries(): Promise<DictionaryMap> {
    const response = await publicApi.get('dictionaries')
    return response.json()
  }
}

export const dictionaryService = new DictionaryService()
