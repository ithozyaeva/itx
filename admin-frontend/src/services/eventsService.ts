import type { CommunityEvent } from '@/models/events'
import { BaseService } from '@/services/api/baseService'

class EventsService extends BaseService<CommunityEvent> {
  constructor() {
    super('events')
  }
}

export const eventsService = new EventsService()
