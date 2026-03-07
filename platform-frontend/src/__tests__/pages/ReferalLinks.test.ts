import { describe, expect, it } from 'vitest'

describe('ReferalLinks logic', () => {
  const ITEMS_PER_PAGE = 10

  interface Link {
    id: number
    [key: string]: unknown
  }

  function handleLinkUpdated(links: Link[], updatedLink: Link) {
    const result = [...links]
    const index = result.findIndex(link => link.id === updatedLink.id)
    if (index !== -1) {
      result[index] = updatedLink
    }
    return result
  }

  function handleLinkDeleted(links: Link[], deletedLinkId: number) {
    const result = [...links]
    const index = result.findIndex(link => link.id === deletedLinkId)
    if (index !== -1) {
      result.splice(index, 1)
    }
    return result
  }

  describe('ITEMS_PER_PAGE', () => {
    it('equals 10', () => {
      expect(ITEMS_PER_PAGE).toBe(10)
    })
  })

  describe('handleLinkUpdated', () => {
    it('updates an existing link by id', () => {
      const links: Link[] = [
        { id: 1, url: 'https://a.com' },
        { id: 2, url: 'https://b.com' },
      ]
      const updated: Link = { id: 2, url: 'https://updated.com' }
      const result = handleLinkUpdated(links, updated)
      expect(result[1]).toEqual({ id: 2, url: 'https://updated.com' })
    })

    it('does not modify the array if id is not found', () => {
      const links: Link[] = [
        { id: 1, url: 'https://a.com' },
      ]
      const updated: Link = { id: 99, url: 'https://new.com' }
      const result = handleLinkUpdated(links, updated)
      expect(result).toEqual(links)
    })

    it('preserves other links unchanged', () => {
      const links: Link[] = [
        { id: 1, url: 'https://a.com' },
        { id: 2, url: 'https://b.com' },
        { id: 3, url: 'https://c.com' },
      ]
      const updated: Link = { id: 2, url: 'https://updated.com' }
      const result = handleLinkUpdated(links, updated)
      expect(result[0]).toEqual({ id: 1, url: 'https://a.com' })
      expect(result[2]).toEqual({ id: 3, url: 'https://c.com' })
      expect(result).toHaveLength(3)
    })
  })

  describe('handleLinkDeleted', () => {
    it('removes a link by id', () => {
      const links: Link[] = [
        { id: 1, url: 'https://a.com' },
        { id: 2, url: 'https://b.com' },
      ]
      const result = handleLinkDeleted(links, 1)
      expect(result).toHaveLength(1)
      expect(result[0].id).toBe(2)
    })

    it('does not modify the array if id is not found', () => {
      const links: Link[] = [
        { id: 1, url: 'https://a.com' },
      ]
      const result = handleLinkDeleted(links, 99)
      expect(result).toHaveLength(1)
      expect(result).toEqual(links)
    })

    it('handles deleting the last remaining link', () => {
      const links: Link[] = [{ id: 1, url: 'https://a.com' }]
      const result = handleLinkDeleted(links, 1)
      expect(result).toHaveLength(0)
    })

    it('handles empty array', () => {
      const result = handleLinkDeleted([], 1)
      expect(result).toHaveLength(0)
    })
  })
})
