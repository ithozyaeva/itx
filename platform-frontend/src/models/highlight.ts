export interface ChatHighlight {
  id: number
  chatId: number
  messageId: number
  authorTelegramId: number
  authorUsername: string
  authorFirstName: string
  messageText: string
  memberId?: number
  createdAt: string
}
