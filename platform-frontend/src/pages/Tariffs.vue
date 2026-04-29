<script setup lang="ts">
import type { PublicTier } from '@/services/subscriptions'
import { Check, Crown } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { Typography } from '@/components/ui/typography'
import { handleError } from '@/services/errorService'
import { subscriptionsService } from '@/services/subscriptions'

const tiers = ref<PublicTier[]>([])
const loading = ref(true)

const botUsername = import.meta.env.VITE_TELEGRAM_BOT_NAME
const botSubLink = botUsername ? `https://t.me/${botUsername}?start=sub` : ''

onMounted(async () => {
  try {
    tiers.value = await subscriptionsService.getPublicTiers()
  }
  catch (e) {
    handleError(e)
  }
  finally {
    loading.value = false
  }
})

function formatPrice(price: number) {
  return price.toLocaleString('ru-RU')
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/tariffs
    </div>
    <Typography
      variant="h2"
      as="h1"
      class="mb-2"
    >
      Тарифы подписки
    </Typography>
    <p class="text-muted-foreground mb-8 max-w-2xl">
      Подписка открывает доступ к закрытым чатам, событиям, бирже заданий и менторам.
      Оплата через Boosty — отмена в любой момент.
    </p>

    <div
      v-if="loading"
      class="text-center py-16 text-muted-foreground"
    >
      Загрузка тарифов…
    </div>

    <div
      v-else-if="tiers.length === 0"
      class="text-center py-16 text-muted-foreground"
    >
      Тарифы временно недоступны. Попробуйте позже или напишите админу.
    </div>

    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6"
    >
      <div
        v-for="tier in tiers"
        :key="tier.id"
        class="relative flex flex-col rounded-sm border bg-card p-6 transition-colors hover:border-accent/40"
      >
        <div class="flex items-center justify-between mb-3">
          <span class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground/60">
            [ уровень.{{ String(tier.level).padStart(2, '0') }} ]
          </span>
          <Crown
            v-if="tier.level >= 3"
            class="w-4 h-4 text-accent/60"
          />
        </div>
        <h3 class="font-display uppercase text-2xl mb-2">
          {{ tier.name }}
        </h3>
        <div class="flex items-baseline gap-1 mb-4">
          <span class="text-3xl font-bold">{{ formatPrice(tier.price) }}</span>
          <span class="text-muted-foreground">₽/мес</span>
        </div>
        <p
          v-if="tier.description"
          class="text-sm text-muted-foreground mb-5"
        >
          {{ tier.description }}
        </p>
        <ul class="flex-1 space-y-2 mb-6">
          <li
            v-for="(feature, i) in tier.features"
            :key="i"
            class="flex items-start gap-2 text-sm"
          >
            <Check
              class="w-4 h-4 mt-0.5 shrink-0 text-accent"
              stroke-width="3"
            />
            <span>{{ feature }}</span>
          </li>
        </ul>
        <a
          :href="tier.boosty_url"
          target="_blank"
          rel="noopener noreferrer"
          class="block text-center px-4 py-3 rounded-sm bg-accent text-accent-foreground font-medium hover:bg-accent/90 transition-colors"
        >
          Оформить на Boosty →
        </a>
      </div>
    </div>

    <div
      v-if="botSubLink"
      class="mt-10 p-6 rounded-sm border border-dashed border-accent/30 bg-accent/[0.03]"
    >
      <h3 class="font-medium mb-2">
        Уже оплатил?
      </h3>
      <p class="text-sm text-muted-foreground mb-4">
        Бот проверит подписку и выдаст инвайты в чаты по вашему тиру.
      </p>
      <a
        :href="botSubLink"
        target="_blank"
        rel="noopener"
        class="inline-block px-4 py-2 rounded-sm border border-accent/50 text-accent hover:bg-accent/10 transition-colors text-sm"
      >
        Открыть бота → /sub
      </a>
    </div>
  </div>
</template>
