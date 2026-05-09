<script setup lang="ts">
import type { ReferralCabinet } from '@/services/referral'
import { Check, Coins, Copy, Share2, TrendingUp, Users } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { formatShortDate } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { referralService } from '@/services/referral'

const data = ref<ReferralCabinet | null>(null)
const loading = ref(true)
const copying = ref(false)

const { toast } = useToast()

function formatNumber(n: number) {
  return n.toLocaleString('ru-RU')
}

const conversionRate = computed(() => {
  if (!data.value || data.value.invitedTotal === 0)
    return 0
  return Math.round((Number(data.value.withActiveSub) / Number(data.value.invitedTotal)) * 100)
})

// Короткое представление для UI: t.me/<bot>?…/ref_<code> может быть длинным
// и в `<code>` с break-all расплывается на 3 строки. Показываем «t.me/…/ref_X»,
// в title — full URL, в clipboard — full URL.
const DEEPLINK_RE = /^https?:\/\/(t\.me\/)[^?]+\?start=(ref_\w+)$/
const shortDeeplink = computed(() => {
  const url = data.value?.deeplink
  if (!url)
    return ''
  const match = url.match(DEEPLINK_RE)
  return match ? `${match[1]}…/${match[2]}` : url
})

async function copyDeeplink() {
  if (!data.value?.deeplink)
    return
  copying.value = true
  try {
    await navigator.clipboard.writeText(data.value.deeplink)
    toast({ title: 'Ссылка скопирована', description: 'За каждого приглашённого по ней — кредиты на ваш баланс.' })
  }
  catch {
    toast({ title: 'Не удалось скопировать', description: data.value.deeplink })
  }
  finally {
    copying.value = false
  }
}

function shareNative() {
  if (!data.value?.deeplink)
    return
  // Web Share API — на iOS/Android открывает нативный share-sheet с Telegram, Whatsapp и т.п.
  if (navigator.share) {
    navigator.share({
      title: 'IT-X — сообщество',
      text: 'Присоединяйся к IT-X — закрытые чаты, события, биржа заданий и менторы.',
      url: data.value.deeplink,
    }).catch(() => {})
  }
  else {
    copyDeeplink()
  }
}

onMounted(async () => {
  try {
    data.value = await referralService.getCabinet()
  }
  catch (e) {
    handleError(e)
  }
  finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-3xl">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/referral
    </div>
    <Typography variant="h2" as="h1" class="mb-3">
      Пригласить в IT-X
    </Typography>
    <p class="text-muted-foreground mb-6 max-w-2xl">
      Поделитесь персональной ссылкой. За каждого приглашённого вам начисляются кредиты,
      которыми можно оплатить свою подписку.
    </p>

    <div
      v-if="loading"
      class="text-center py-16 text-muted-foreground"
    >
      Загрузка…
    </div>

    <template v-else-if="data">
      <!-- Персональная ссылка -->
      <div class="rounded-sm border bg-card p-5 mb-6">
        <div class="flex items-center gap-2 mb-3 text-sm text-muted-foreground">
          <Share2 class="w-4 h-4" />
          <span>Ваша персональная ссылка</span>
        </div>
        <div class="flex items-center gap-2 mb-4">
          <code
            class="flex-1 px-3 py-2 rounded-sm bg-muted/50 text-sm font-mono truncate"
            :title="data.deeplink || ''"
          >
            {{ shortDeeplink || '—' }}
          </code>
          <Button
            :disabled="!data.deeplink || copying"
            variant="outline"
            size="icon"
            @click="copyDeeplink"
          >
            <Copy class="w-4 h-4" />
          </Button>
        </div>
        <Button
          :disabled="!data.deeplink"
          class="w-full"
          @click="shareNative"
        >
          <Share2 class="w-4 h-4 mr-2" />
          Поделиться ссылкой
        </Button>
      </div>

      <!-- Статистика -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 mb-6">
        <div class="rounded-sm border bg-card p-4">
          <div class="flex items-center gap-2 text-xs text-muted-foreground mb-2">
            <Users class="w-3.5 h-3.5" /> Приглашено
          </div>
          <div class="text-2xl font-bold">
            {{ formatNumber(data.invitedTotal) }}
          </div>
        </div>
        <div class="rounded-sm border bg-card p-4">
          <div class="flex items-center gap-2 text-xs text-muted-foreground mb-2">
            <Check class="w-3.5 h-3.5" /> С активной подпиской
          </div>
          <div class="text-2xl font-bold">
            {{ formatNumber(data.withActiveSub) }}
            <span
              v-if="data.invitedTotal > 0"
              class="text-sm text-muted-foreground font-normal"
            >
              · {{ conversionRate }}%
            </span>
          </div>
        </div>
        <div class="rounded-sm border bg-card p-4">
          <div class="flex items-center gap-2 text-xs text-muted-foreground mb-2">
            <TrendingUp class="w-3.5 h-3.5" /> Заработано всего
          </div>
          <div class="text-2xl font-bold">
            {{ formatNumber(data.totalEarned) }} <span class="text-sm text-muted-foreground font-normal">кр.</span>
          </div>
        </div>
      </div>

      <!-- Балансовая карточка -->
      <RouterLink
        to="/credits"
        class="flex items-center gap-3 p-4 mb-6 rounded-sm border bg-card hover:border-accent/40 transition-colors group"
      >
        <div class="rounded-sm bg-accent/10 p-2.5">
          <Coins class="w-5 h-5 text-accent" />
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm text-muted-foreground">
            Доступный баланс
          </div>
          <div class="text-xl font-bold">
            {{ formatNumber(data.balance) }} кр.
          </div>
        </div>
        <span class="text-xs text-muted-foreground group-hover:text-foreground transition-colors">
          история и покупка →
        </span>
      </RouterLink>

      <!-- Список приглашённых -->
      <Typography variant="h4" as="h2" class="mb-3">
        Последние приглашённые
      </Typography>

      <div
        v-if="data.recentInvitees.length === 0"
        class="rounded-sm border border-dashed bg-card/50 p-8 text-center text-muted-foreground"
      >
        Пока никого. Поделитесь ссылкой выше.
      </div>

      <div
        v-else
        class="space-y-2"
      >
        <div
          v-for="invitee in data.recentInvitees"
          :key="invitee.id"
          class="flex items-center gap-3 p-3 bg-card border border-border rounded-sm"
        >
          <div class="w-9 h-9 rounded-full bg-muted overflow-hidden flex items-center justify-center text-sm font-medium shrink-0">
            <img
              v-if="invitee.avatarUrl"
              :src="invitee.avatarUrl"
              :alt="invitee.firstName"
              class="w-full h-full object-cover"
            >
            <span v-else>{{ (invitee.firstName || invitee.tg || '?').charAt(0).toUpperCase() }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium truncate">
              {{ [invitee.firstName, invitee.lastName].filter(Boolean).join(' ') || `@${invitee.tg}` }}
            </div>
            <div class="text-xs text-muted-foreground">
              {{ formatShortDate(invitee.joinedAt) }}
              <span
                v-if="invitee.tg"
                class="ml-1"
              >· @{{ invitee.tg }}</span>
            </div>
          </div>
          <span
            class="text-xs px-2 py-1 rounded-sm shrink-0"
            :class="invitee.hasActiveSub
              ? 'bg-green-500/10 text-green-600 dark:text-green-400'
              : 'bg-muted text-muted-foreground'"
          >
            {{ invitee.hasActiveSub ? 'Подписан' : 'Без подписки' }}
          </span>
        </div>
      </div>
    </template>
  </div>
</template>
