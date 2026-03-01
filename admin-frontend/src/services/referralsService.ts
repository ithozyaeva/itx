import type { ReferralLink } from '@/models/referrals'
import { BaseService } from '@/services/api/baseService'

class ReferralsService extends BaseService<ReferralLink> {
  constructor() {
    super('admin-referals')
  }
}

export const referralsService = new ReferralsService()
