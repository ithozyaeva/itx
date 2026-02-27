import axios from 'axios'

export interface DictionaryItem<TValue extends string> {
  value: TValue
  label: string
}

export interface DictionaryMap {
  [key: string]: DictionaryItem<string>[]
}

export class DictionaryService {
  async getAllDictionaries(): Promise<DictionaryMap> {
    const response = await axios.get('/api/dictionaries')
    return response.data
  }
}

export const dictionaryService = new DictionaryService()
