import { describe, expect, it } from 'vitest'
import type { Raffle, RaffleCreateRequest } from '@/services/raffleService'

describe('raffleService types', () => {
  it('Raffle type has correct shape', () => {
    const raffle: Raffle = {
      id: 1,
      title: 'Test Raffle',
      description: 'A test raffle',
      prize: 'Prize',
      ticketCost: 10,
      maxTickets: 100,
      endsAt: '2026-04-01T00:00:00Z',
      status: 'ACTIVE',
      winnerId: null,
      createdAt: '2026-03-01T00:00:00Z',
    }

    expect(raffle.id).toBe(1)
    expect(raffle.title).toBe('Test Raffle')
    expect(raffle.status).toBe('ACTIVE')
    expect(raffle.winnerId).toBeNull()
  })

  it('Raffle status can be FINISHED', () => {
    const raffle: Raffle = {
      id: 2,
      title: 'Finished Raffle',
      description: 'Done',
      prize: 'Prize 2',
      ticketCost: 5,
      maxTickets: 50,
      endsAt: '2026-02-01T00:00:00Z',
      status: 'FINISHED',
      winnerId: 42,
      createdAt: '2026-01-01T00:00:00Z',
    }

    expect(raffle.status).toBe('FINISHED')
    expect(raffle.winnerId).toBe(42)
  })

  it('RaffleCreateRequest has required fields', () => {
    const request: RaffleCreateRequest = {
      title: 'New Raffle',
      description: 'Description',
      prize: 'A prize',
      ticketCost: 15,
      maxTickets: 200,
      endsAt: '2026-05-01T00:00:00Z',
    }

    expect(request.title).toBe('New Raffle')
    expect(request.ticketCost).toBe(15)
    expect(request.maxTickets).toBe(200)
  })
})
