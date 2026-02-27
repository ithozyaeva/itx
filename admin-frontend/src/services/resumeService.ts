import type { Resume, ResumeFilter } from '@/models/resume'
import api from '@/lib/api'
import { cleanParams } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { BaseService } from './api/baseService'

class ResumeService extends BaseService<Resume> {
  constructor() {
    super('resumes')
  }

  async searchWithFilters(filter?: ResumeFilter) {
    return this.search(filter)
  }

  async downloadArchive(filter?: ResumeFilter) {
    try {
      const response = await api.get('resumes/download', {
        searchParams: cleanParams(filter ?? {}),
      })
      const blob = await response.blob()
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = 'resumes.zip'
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)
    }
    catch (error) {
      handleError(error)
    }
  }
}

export const resumeService = new ResumeService()
