<script setup lang="ts">
import type { PublicTier } from '@/services/subscriptions'
import { Check, Coins, Crown, Share2 } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { openLink } from '@/composables/useTelegramWebApp'
import { creditsService } from '@/services/credits'
import { handleError } from '@/services/errorService'
import { subscriptionsService } from '@/services/subscriptions'

const tiers = ref<PublicTier[]>([])
const balance = ref<number | null>(null)
const loading = ref(true)
const purchasing = ref(false)
const confirmTier = ref<PublicTier | null>(null)

const { toast } = useToast()

const botUsername = import.meta.env.VITE_TELEGRAM_BOT_NAME
const botSubLink = botUsername ? `https://t.me/${botUsername}?start=sub` : ''

onMounted(async () => {
  try {
    const [tiersResp, creditsResp] = await Promise.all([
      subscriptionsService.getPublicTiers(),
      creditsService.getMine().catch(() => null),
    ])
    tiers.value = tiersResp
    if (creditsResp)
      balance.value = creditsResp.balance
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

const sortedTiers = computed(() => [...tiers.value].sort((a, b) => a.level - b.level))

function canBuy(tier: PublicTier) {
  return tier.price_credits != null && balance.value != null && balance.value >= tier.price_credits
}

// Минимальная цена среди тарифов с price_credits — для блока «не хватает».
// Если ни один тариф нельзя купить за credits, блок не показываем.
const minCreditsPrice = computed(() => {
  const prices = tiers.value
    .map(t => t.price_credits)
    .filter((p): p is number => p != null && p > 0)
  return prices.length ? Math.min(...prices) : null
})

const showEarnCTA = computed(() => {
  if (balance.value == null || minCreditsPrice.value == null)
    return false
  return balance.value < minCreditsPrice.value
})

async function confirmPurchase() {
  if (!confirmTier.value)
    return
  const slug = confirmTier.value.slug
  const tierName = confirmTier.value.name
  purchasing.value = true
  try {
    const result = await creditsService.purchaseTier(slug)
    balance.value = result.balance_left
    const expires = new Date(result.expires_at).toLocaleDateString('ru-RU')
    toast({
      title: 'Подписка активирована',
      description: `«${tierName}» активна до ${expires}. Бот пришлёт инвайты в чаты в течение получаса.`,
    })
    confirmTier.value = null
  }
  catch (e) {
    handleError(e)
  }
  finally {
    purchasing.value = false
  }
}
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/tariffs
    </div>
    <div class="flex items-end justify-between mb-2 gap-4 flex-wrap">
      <Typography
        variant="h2"
        as="h1"
      >
        Тарифы подписки
      </Typography>
      <RouterLink
        v-if="balance != null"
        to="/credits"
        class="flex items-center gap-2.5 rounded-sm border border-border bg-card px-3 py-2 hover:border-accent/40 transition-colors"
      >
        <Coins class="h-5 w-5 text-accent" />
        <div>
          <p class="text-lg font-bold leading-none">
            {{ formatPrice(balance) }}
          </p>
          <p class="text-[11px] text-muted-foreground mt-0.5">
            кредитов
          </p>
        </div>
      </RouterLink>
    </div>
    <p class="text-muted-foreground mb-8 max-w-2xl">
      Подписка открывает доступ к закрытым чатам, событиям, бирже заданий и менторам.
      Можно оплатить через Boosty или реферальными кредитами.
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
        v-for="tier in sortedTiers"
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
        <div class="flex items-baseline gap-1 mb-1">
          <span class="text-3xl font-bold">{{ formatPrice(tier.price) }}</span>
          <span class="text-muted-foreground">₽/мес</span>
        </div>
        <div
          v-if="tier.price_credits != null"
          class="flex items-center gap-1.5 mb-4 text-sm text-muted-foreground"
        >
          <Coins class="w-3.5 h-3.5 text-accent/70" />
          <span>или {{ formatPrice(tier.price_credits) }} кредитов / 30 дней</span>
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
        <div class="space-y-2">
          <a
            v-if="tier.boosty_url"
            :href="tier.boosty_url"
            class="block text-center px-4 py-3 rounded-sm bg-accent text-accent-foreground font-medium hover:bg-accent/90 transition-colors"
            @click.prevent="openLink(tier.boosty_url)"
          >
            Оформить на Boosty →
          </a>
          <Button
            v-if="tier.price_credits != null"
            :disabled="!canBuy(tier)"
            variant="outline"
            class="w-full py-3 h-auto"
            @click="confirmTier = tier"
          >
            <Coins class="w-4 h-4 mr-2" />
            <span v-if="balance == null">
              Загрузка баланса…
            </span>
            <span v-else-if="canBuy(tier)">
              Купить за {{ formatPrice(tier.price_credits) }} кр.
            </span>
            <span v-else>
              Не хватает {{ formatPrice(tier.price_credits - balance) }} кр.
            </span>
          </Button>
        </div>
      </div>
    </div>

    <!-- Если баланса не хватает даже на самый дешёвый тариф — мягко
         подсказываем, как набрать кредиты. Не показываем юзерам с balance=null
         (auth/network проблема) и тем, у кого баланс уже достаточен. -->
    <RouterLink
      v-if="showEarnCTA"
      to="/referals"
      class="mt-10 flex items-center gap-3 p-5 rounded-sm border border-accent/30 bg-accent/[0.04] hover:border-accent/50 transition-colors group"
    >
      <Share2 class="w-5 h-5 text-accent shrink-0" />
      <div class="flex-1 min-w-0">
        <p class="text-sm font-medium">
          Не хватает кредитов? Получите за приглашённых
        </p>
        <p class="text-xs text-muted-foreground">
          Создайте реферальную ссылку — за каждую конверсию +30 кр., +50% разово при покупке подписки рефералом, +20% каждый месяц.
        </p>
      </div>
      <span class="text-xs text-accent shrink-0 group-hover:underline">
        создать →
      </span>
    </RouterLink>

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
        class="inline-block px-4 py-2 rounded-sm border border-accent/50 text-accent hover:bg-accent/10 transition-colors text-sm"
        @click.prevent="openLink(botSubLink)"
      >
        Открыть бота → /sub
      </a>
    </div>

    <Dialog
      :open="confirmTier !== null"
      @update:open="(v: boolean) => { if (!v) confirmTier = null }"
    >
      <DialogContent class="sm:max-w-[450px]">
        <DialogHeader>
          <DialogTitle>Купить «{{ confirmTier?.name }}»?</DialogTitle>
          <DialogDescription>
            Спишем {{ formatPrice(confirmTier?.price_credits ?? 0) }} кредитов и активируем тариф на 30 дней.
            <br>
            Если уже была активна подписка — срок продлится на 30 дней от текущей даты истечения.
            <br>
            <br>
            Доступ в чаты появится в течение получаса (как только пройдёт ближайший автопроверка бота).
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button
            variant="outline"
            :disabled="purchasing"
            @click="confirmTier = null"
          >
            Отмена
          </Button>
          <Button
            :disabled="purchasing"
            @click="confirmPurchase"
          >
            <span v-if="purchasing">Покупка…</span>
            <span v-else>Подтвердить</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
