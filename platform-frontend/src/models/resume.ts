export type WorkFormat = 'REMOTE' | 'HYBRID' | 'OFFICE' | ''

export interface Resume {
  id: number
  tgId: number
  fileName: string
  filePath: string
  workExperience?: string
  desiredPosition?: string
  workFormat?: WorkFormat
  createdAt: string
  updatedAt: string
  parsedConfidence?: number
}

export interface ParsedResumeData {
  workExperience?: string
  desiredPosition?: string
  workFormat?: WorkFormat
  confidence?: number
}
