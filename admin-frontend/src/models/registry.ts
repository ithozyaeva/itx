export interface Pagination {
  offset: number
  limit: number
}

export interface Registry<T> {
  items: T[]
  total: number
}
