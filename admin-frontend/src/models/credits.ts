export interface AdminCreditTransaction {
  id: number
  memberId: number
  memberFirstName: string
  memberLastName: string
  memberUsername: string
  amount: number
  reason: string
  sourceType: string
  description: string
  createdAt: string
}

export interface AdminAwardCreditsRequest {
  memberId: number
  amount: number
  description: string
}
