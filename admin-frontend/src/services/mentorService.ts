import type { Mentor } from '@/models/mentors'
import { BaseService } from '@/services/api/baseService'

class MentorService extends BaseService<Mentor> {
  constructor() {
    super('mentors')
  }
}

export const mentorService = new MentorService()
