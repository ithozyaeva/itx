<script setup lang="ts">
import type { CreditsSummary } from '@/services/credits'
import { Coins, Inbox, Share2, Sparkles, Users } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import { Button } from '@/components/ui/button'
import { Typography } from '@/components/ui/typography'
import { formatShortDate } from '@/lib/utils'
import { creditsService } from '@/services/credits'
import { handleError } from '@/services/errorService'

const PAGE_SIZE = 15
const data = ref<CreditsSummary | null>(null)
const loading = ref(true)
const visibleCount = ref(PAGE_SIZE)

const visibleTransactions = computed(() => data.value?.transactions.slice(0, visibleCount.value) ?? [])
const hasMore = computed(() => data.value ? visibleCount.value < data.value.transactions.length : false)

const reasonLabels: Record<string, string> = {
  referal_conversion: 'Конверсия по реф-ссылке',
  referral_purchase_first: 'Реферал впервые оформил подписку',
  referral_purchase_recurring: 'Реферал — ежемесячная активность',
  subscription_purchase: 'Покупка подписки',
  admin_manual: 'Начислено админом',
}

function formatNumber(n: number) {
  return n.toLocaleString('ru-RU')
}

onMounted(async () => {
  try {
    data.value = await creditsService.getMine()
  }
  catch (e) {
    handleError(e)
  }
  finally {
    loading.value = false
  }
})

function loadMore() {
  visibleCount.value += PAGE_SIZE
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-3xl">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/credits
    </div>
    <Typography variant="h2" as="h1" class="mb-4">
      Реферальные кредиты
    </Typography>
    <p class="text-muted-foreground mb-6 max-w-2xl">
      Тратятся на покупку подписки. Начисляются за приглашённых: единоразовая
      выплата при первой покупке реферала и ежемесячная — пока он активен.
    </p>

    <div
      v-if="loading"
      class="text-center py-16 text-muted-foreground"
    >
      Загрузка…
    </div>

    <template v-else-if="data">
      <div class="rounded-sm border bg-card p-5 flex items-center gap-4 mb-6">
        <div class="rounded-sm bg-accent/10 p-3">
          <Coins class="w-7 h-7 text-accent" />
        </div>
        <div class="flex-1">
          <div class="text-sm text-muted-foreground">
            Текущий баланс
          </div>
          <div class="text-3xl font-bold">
            {{ formatNumber(data.balance) }}
          </div>
        </div>
        <RouterLink
          to="/tariffs"
          class="px-4 py-2 rounded-sm bg-accent text-accent-foreground text-sm font-medium hover:bg-accent/90 transition-colors"
        >
          Купить подписку →
        </RouterLink>
      </div>

      <!-- Onboarding-блок «как заработать»: показываем тем, у кого нет
           конверсий ИЛИ пустой баланс — иначе они не поймут, как пополнить.
           Реф-программа физически реализована через /referals (ссылки на
           вакансии), но связь между «креди ты» и «ссылки на вакансии»
           неочевидна без явного объяснения. -->
      <div
        v-if="data.transactions.length === 0 || data.balance === 0"
        class="rounded-sm border border-accent/30 bg-accent/[0.04] p-5 mb-6"
      >
        <div class="flex items-start gap-3 mb-3">
          <Sparkles class="w-5 h-5 text-accent shrink-0 mt-0.5" />
          <div>
            <h3 class="font-medium mb-1">
              Как зарабатывать кредиты
            </h3>
            <p class="text-sm text-muted-foreground">
              Создавайте реферальные ссылки на вакансии в своей компании. За каждого приглашённого участника + за его покупку подписки вам начисляются кредиты, которыми можно оплатить свою подписку.
            </p>
          </div>
        </div>
        <ul class="space-y-2 mb-4 text-sm pl-8">
          <li class="flex items-center gap-2">
            <Users class="w-4 h-4 text-muted-foreground shrink-0" />
            <span><strong>30 кр.</strong> — за каждого, кто пришёл по вашей ссылке</span>
          </li>
          <li class="flex items-center gap-2">
            <Coins class="w-4 h-4 text-muted-foreground shrink-0" />
            <span><strong>50%</strong> — разово, когда реферал оформил подписку</span>
          </li>
          <li class="flex items-center gap-2">
            <Coins class="w-4 h-4 text-muted-foreground shrink-0" />
            <span><strong>20%</strong> — каждый месяц, пока реферал активен</span>
          </li>
        </ul>
        <RouterLink
          to="/referals"
          class="inline-flex items-center gap-2 px-4 py-2 rounded-sm bg-accent text-accent-foreground text-sm font-medium hover:bg-accent/90 transition-colors"
        >
          <Share2 class="w-4 h-4" />
          Создать реферальную ссылку
        </RouterLink>
      </div>

      <Typography variant="h4" as="h2" class="mb-3">
        История транзакций
      </Typography>

      <EmptyState
        v-if="data.transactions.length === 0"
        :icon="Inbox"
        variant="dashed"
        title="Пока нет транзакций"
        description="Создайте реферальную ссылку — за каждую конверсию +30 кредитов."
      />

      <div
        v-else
        class="space-y-2"
      >
        <div
          v-for="tx in visibleTransactions"
          :key="tx.id"
          class="flex items-center gap-4 p-4 bg-card border border-border rounded-sm"
        >
          <div class="flex-1 min-w-0">
            <div class="font-medium text-sm truncate">
              {{ tx.description || reasonLabels[tx.reason] || tx.reason }}
            </div>
            <div class="text-xs text-muted-foreground mt-0.5">
              {{ reasonLabels[tx.reason] ?? tx.reason }}
            </div>
          </div>
          <div class="shrink-0 text-right">
            <div
              class="font-bold"
              :class="tx.amount > 0 ? 'text-green-500' : 'text-red-500'"
            >
              {{ tx.amount > 0 ? '+' : '' }}{{ formatNumber(tx.amount) }}
            </div>
            <div class="text-xs text-muted-foreground">
              {{ formatShortDate(tx.createdAt) }}
            </div>
          </div>
        </div>
        <div
          v-if="hasMore"
          class="mt-4 flex justify-center"
        >
          <Button
            variant="outline"
            @click="loadMore"
          >
            Показать ещё
          </Button>
        </div>
      </div>
    </template>
  </div>
</template>
