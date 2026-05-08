<script setup lang="ts">
import type { ReferrerInfo } from '@/services/referrer'
import { Sparkles } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { isUserSubscribed } from '@/composables/useUser'
import { referrerService } from '@/services/referrer'

const referrer = ref<ReferrerInfo | null>(null)
const open = ref(false)

function fullName(a: ReferrerInfo['author']) {
  return [a.firstName, a.lastName].filter(Boolean).join(' ').trim() || `@${a.tg}`
}

async function close() {
  open.value = false
  try {
    await referrerService.markSeen()
  }
  catch {
    // best-effort: если ack упал, при следующем заходе модалка покажется снова — приемлемо
  }
}

onMounted(async () => {
  // Подписчики не нуждаются в призыве «оформить подписку» — баннер
  // нацелен только на UNSUBSCRIBER. Если юзер уже подписан — атрибуция
  // уже сработала в auth, smb-баннер не нужен.
  if (isUserSubscribed().value)
    return
  try {
    const info = await referrerService.getMine()
    if (info && !info.seenAt) {
      referrer.value = info
      open.value = true
    }
  }
  catch {
    // network/auth fail — silent
  }
})
</script>

<template>
  <Dialog
    :open="open"
    @update:open="(v: boolean) => { if (!v) close() }"
  >
    <DialogContent class="sm:max-w-[460px]">
      <DialogHeader>
        <div class="flex items-center gap-3 mb-2">
          <div class="rounded-full bg-accent/15 p-2">
            <Sparkles class="w-5 h-5 text-accent" />
          </div>
          <DialogTitle>Вас пригласили в IT-X</DialogTitle>
        </div>
        <DialogDescription>
          <span v-if="referrer">
            <strong class="text-foreground">{{ fullName(referrer.author) }}</strong>
            пригласил вас в сообщество. Подписка открывает доступ к закрытым чатам, событиям, бирже заданий и менторам — оформите тариф, чтобы начать.
          </span>
        </DialogDescription>
      </DialogHeader>

      <div
        v-if="referrer"
        class="flex items-center gap-3 rounded-sm border bg-card p-3 my-2"
      >
        <div class="w-10 h-10 rounded-full bg-muted overflow-hidden flex items-center justify-center text-sm font-medium">
          <img
            v-if="referrer.author.avatarUrl"
            :src="referrer.author.avatarUrl"
            :alt="fullName(referrer.author)"
            class="w-full h-full object-cover"
          >
          <span v-else>{{ fullName(referrer.author).charAt(0).toUpperCase() }}</span>
        </div>
        <div class="flex-1 min-w-0">
          <div class="text-sm font-medium truncate">
            {{ fullName(referrer.author) }}
          </div>
          <div
            v-if="referrer.author.tg"
            class="text-xs text-muted-foreground"
          >
            @{{ referrer.author.tg }}
          </div>
        </div>
      </div>

      <DialogFooter class="gap-2 sm:gap-2">
        <Button
          variant="outline"
          @click="close"
        >
          Позже
        </Button>
        <RouterLink
          to="/tariffs"
          class="inline-flex items-center justify-center px-4 py-2 rounded-sm bg-accent text-accent-foreground text-sm font-medium hover:bg-accent/90 transition-colors"
          @click="close"
        >
          Перейти к тарифам →
        </RouterLink>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
