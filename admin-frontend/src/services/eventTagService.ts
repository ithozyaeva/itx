import type { EventTag } from '@/models/events.ts'
import { BaseService } from '@/services/api/baseService'

class EventTagService extends BaseService<EventTag> {
  constructor() {
    super('eventTags')
  }
}

export const eventTagService = new EventTagService()
