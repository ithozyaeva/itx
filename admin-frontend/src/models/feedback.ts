export interface Feedback {
  id: number
  userId: number | null
  userFirstName: string | null
  userLastName: string | null
  userUsername: string | null
  score: number
  comment: string | null
  createdAt: string
}
