export type Block
  = | { type: 'h2', text: string, id?: string }
    | { type: 'h3', text: string, id?: string }
    | { type: 'p', html: string }
    | { type: 'ul', items: string[] }
    | { type: 'ol', items: string[] }
    | { type: 'note', html: string }
    | { type: 'cta', href: string, label: string, external?: boolean }

export interface ArticleFaq {
  q: string
  a: string
}

export interface Article {
  slug: string
  title: string
  h1: string
  breadcrumb: string
  lead: string
  description: string
  publishedAt: string
  updatedAt?: string
  excerpt: string
  tags?: string[]
  faq?: ArticleFaq[]
  body: Block[]
}
