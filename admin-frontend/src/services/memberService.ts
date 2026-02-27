import type { Member } from '@/models/members'
import { BaseService } from '@/services/api/baseService'

class MemberService extends BaseService<Member> {
  constructor() {
    super('members')
  }
}

export const memberService = new MemberService()
