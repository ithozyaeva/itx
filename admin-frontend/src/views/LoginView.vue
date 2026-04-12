<script setup lang="ts">
import { Shield } from 'lucide-vue-next'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import TelegramAuth from '@/components/TelegramAuth.vue'
import { isAuthenticated } from '@/services/authService'

const router = useRouter()

function handleTelegramAuthSuccess() {
  if (isAuthenticated.value) {
    router.push('/dashboard')
  }
}

onMounted(() => {
  if (isAuthenticated.value) {
    router.push('/dashboard')
  }
})
</script>

<template>
  <div class="flex h-screen w-full items-center justify-center relative overflow-hidden">
    <div class="absolute inset-0 dark:opacity-100 opacity-0 transition-opacity">
      <div class="absolute inset-0" style="background-image: linear-gradient(to right, hsl(151 5% 15% / 0.15) 1px, transparent 1px), linear-gradient(to bottom, hsl(151 5% 15% / 0.15) 1px, transparent 1px); background-size: 48px 48px;" />
    </div>

    <div class="relative w-full max-w-sm mx-4">
      <div class="terminal-card bg-card p-6 space-y-6">
        <div class="space-y-3">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-sm bg-accent/10 border border-accent/20 flex items-center justify-center">
              <Shield class="w-5 h-5 text-accent" />
            </div>
            <div>
              <div class="font-mono text-[10px] uppercase tracking-wider text-muted-foreground">
                // authentication
              </div>
              <h1 class="text-lg font-semibold">
                Админ-панель
              </h1>
            </div>
          </div>
          <div class="font-mono text-xs text-muted-foreground border-t border-border pt-3">
            <span class="text-accent">$</span> sudo login --method telegram
          </div>
        </div>

        <TelegramAuth @auth-success="handleTelegramAuthSuccess" />

        <div class="flex items-center justify-between font-mono text-[10px] text-muted-foreground/50 pt-2 border-t border-border">
          <span>IT-X ADMIN v2.0</span>
          <span class="flex items-center gap-1">
            <span class="w-1.5 h-1.5 rounded-full bg-accent shadow-[0_0_4px_hsl(var(--accent))]" />
            secure
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
