import type { ProfTag } from '@/models/mentors'
import { BaseService } from '@/services/api/baseService'

class ProfTagService extends BaseService<ProfTag> {
  constructor() {
    super('profTags')
  }
}

export const profTagService = new ProfTagService()
