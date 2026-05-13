import type { Article } from './types'
import { article as cursorVsClaudeCodeVsWindsurf } from './cursor-vs-claude-code-vs-windsurf'
import { article as itSoobshchestvaRossii } from './it-soobshchestva-rossii'
import { article as kakStatAiInzhenerom } from './kak-stat-ai-inzhenerom'
import { article as kakVybratMentoraVIt } from './kak-vybrat-mentora-v-it'
import { article as podgotovkaKSobesedovaniyuSLlm } from './podgotovka-k-sobesedovaniyu-s-llm'
import { article as vaybkodingNaPraktike } from './vaybkoding-na-praktike'

export const articles: Article[] = [
  cursorVsClaudeCodeVsWindsurf,
  vaybkodingNaPraktike,
  kakStatAiInzhenerom,
  podgotovkaKSobesedovaniyuSLlm,
  kakVybratMentoraVIt,
  itSoobshchestvaRossii,
]

export function getArticleBySlug(slug: string): Article | undefined {
  return articles.find(a => a.slug === slug)
}

export { renderBlocksToHtml } from './render'
export type { Article, ArticleFaq, Block } from './types'
