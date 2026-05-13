import type { Block } from './types'

const AMP_RE = /&/g
const QUOT_RE = /"/g
const LT_RE = /</g
const GT_RE = />/g

function escapeAttr(value: string): string {
  return value.replace(AMP_RE, '&amp;').replace(QUOT_RE, '&quot;').replace(LT_RE, '&lt;').replace(GT_RE, '&gt;')
}

export function renderBlocksToHtml(blocks: Block[]): string {
  const parts: string[] = []
  for (const block of blocks) {
    switch (block.type) {
      case 'h2': {
        const id = block.id ? ` id="${escapeAttr(block.id)}"` : ''
        parts.push(`<h2${id}>${block.text}</h2>`)
        break
      }
      case 'h3': {
        const id = block.id ? ` id="${escapeAttr(block.id)}"` : ''
        parts.push(`<h3${id}>${block.text}</h3>`)
        break
      }
      case 'p':
        parts.push(`<p>${block.html}</p>`)
        break
      case 'ul':
        parts.push(`<ul>${block.items.map(i => `<li>${i}</li>`).join('')}</ul>`)
        break
      case 'ol':
        parts.push(`<ol>${block.items.map(i => `<li>${i}</li>`).join('')}</ol>`)
        break
      case 'note':
        parts.push(`<aside class="article-note">${block.html}</aside>`)
        break
      case 'cta': {
        const rel = block.external ? ' rel="noopener"' : ''
        parts.push(`<p class="article-cta"><a href="${escapeAttr(block.href)}"${rel}>${block.label} →</a></p>`)
        break
      }
    }
  }
  return parts.join('\n')
}
