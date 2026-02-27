import type { ParsedResumeData, Resume, WorkFormat } from '@/models/resume'
import { apiClient } from './api'

export interface UploadResumePayload {
  file: File
  workExperience?: string
  desiredPosition?: string
  workFormat?: WorkFormat
}

export const resumeService = {
  async upload(payload: UploadResumePayload) {
    const formData = new FormData()
    formData.append('file', payload.file)
    if (payload.workExperience)
      formData.append('workExperience', payload.workExperience)
    if (payload.desiredPosition)
      formData.append('desiredPosition', payload.desiredPosition)
    if (payload.workFormat)
      formData.append('workFormat', payload.workFormat)

    const response = await apiClient.post('resumes', {
      body: formData,
    })

    return response.json<{ resume: Resume, parsed: ParsedResumeData }>()
  },

  async listMine() {
    const response = await apiClient.get('resumes/me')
    return response.json<Resume[]>()
  },

  async update(id: number, payload: Partial<Pick<Resume, 'workExperience' | 'desiredPosition' | 'workFormat'>>) {
    const response = await apiClient.patch(`resumes/${id}`, {
      json: payload,
    })
    return response.json<Resume>()
  },

  async delete(id: number) {
    await apiClient.delete(`resumes/${id}`)
  },
}
