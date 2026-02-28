export interface AdminPointTransaction {
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

export interface AdminAwardRequest {
  memberId: number
  amount: number
  description: string
}
