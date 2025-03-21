import type { DeepPartial } from '@/lib/utils'
import type { Pagination, Registry } from '@/models/registry'
import { useToast } from '@/components/ui/toast'
import api from '@/lib/api'
import { handleError } from '@/services/errorService'
import { ref } from 'vue'

/**
 * Базовый сервис взаимодействяи сущностей с бэкендом.
 */
export class BaseService<T> {
  /** Префикс для API. */
  protected basePath: string
  /** Метка загрузки */
  public isLoading = ref(false)
  /** Массив сущностей */
  public items = ref<Registry<T>>({ items: [], total: 0 })
  /** Пагинация */
  public pagination = ref<Pagination>({
    limit: 10,
    offset: 0,
  })

  private toast = useToast()

  constructor(basePath: string) {
    this.basePath = basePath
  }

  /** Коллбек изменяющий пагинацию, с дальнейшим перезапросом getAll */
  changePagination = (page: number) => {
    this.pagination.value.offset = (page - 1) * this.pagination.value.limit
    this.search()
  }

  /** Метод получения всех сущностей */
  search = async (params?: Record<string, any>): Promise<Registry<T>> => {
    try {
      this.isLoading.value = true
      const searchParams = {
        limit: this.pagination.value.limit,
        offset: this.pagination.value.offset,
        ...params,
      }

      const response = await api.get(this.basePath, { searchParams }).json<Registry<T>>()
      this.items.value = response
      return response
    }
    catch (error) {
      handleError(error)
      return { items: [], total: 0 }
    }
    finally {
      this.isLoading.value = false
    }
  }

  getAll = async (): Promise<Registry<T>> => {
    try {
      this.isLoading.value = true

      const response = await api.get(this.basePath).json<Registry<T>>()
      this.items.value = response
      return response
    }
    catch (error) {
      handleError(error)
      return { items: [], total: 0 }
    }
    finally {
      this.isLoading.value = false
    }
  }

  /** Метод получения сущности по её ID */
  getById = async (id: number): Promise<T | null> => {
    try {
      this.isLoading.value = true
      return await api.get(`${this.basePath}/${id}`).json<T>()
    }
    catch (error) {
      handleError(error)
      return null
    }
    finally {
      this.isLoading.value = false
    }
  }

  /** Метод создания сущности */
  create = async (data: DeepPartial<T>): Promise<T | null> => {
    try {
      this.isLoading.value = true
      const response = await api.post(this.basePath, { json: data }).json<T>()

      this.toast.toast({
        title: 'Успешно',
        description: 'Запись успешно создана',
      })

      await this.search()
      return response
    }
    catch (error) {
      handleError(error)
      return null
    }
    finally {
      this.isLoading.value = false
    }
  }

  /** Метод обновления сущности */
  update = async (id: number, data: DeepPartial<T>): Promise<T | null> => {
    try {
      this.isLoading.value = true
      const response = await api.put(`${this.basePath}/${id}`, { json: data }).json<T>()

      this.toast.toast({
        title: 'Успешно',
        description: 'Запись успешно обновлена',
      })

      await this.search()
      return response
    }
    catch (error) {
      handleError(error)
      return null
    }
    finally {
      this.isLoading.value = false
    }
  }

  /** Метод удаления сущности */
  delete = async (id: number): Promise<boolean> => {
    try {
      this.isLoading.value = true
      await api.delete(`${this.basePath}/${id}`)

      this.toast.toast({
        title: 'Успешно',
        description: 'Запись успешно удалена',
      })

      await this.search()
      return true
    }
    catch (error) {
      handleError(error)
      return false
    }
    finally {
      this.isLoading.value = false
    }
  }
}
