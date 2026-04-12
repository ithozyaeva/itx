<script setup lang="ts">
import { LogOut, Menu, Terminal } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import ThemeToggle from '@/components/ui/theme-toggle.vue'
import { useAuth } from '@/composables/useAuth'
import { useSidebar } from '@/composables/useSidebar'

const { logout } = useAuth()
const { toggleMobileSidebar } = useSidebar()
</script>

<template>
  <header class="sticky top-0 z-30 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 border-b border-border h-14 flex items-center px-4 lg:px-6 justify-between">
    <div class="flex items-center gap-3">
      <Button
        variant="ghost"
        size="icon"
        class="lg:hidden"
        @click="toggleMobileSidebar"
      >
        <Menu class="h-5 w-5" />
      </Button>
      <div class="hidden lg:flex items-center gap-2 font-mono text-xs text-muted-foreground">
        <Terminal class="h-3.5 w-3.5" />
        <span class="tracking-wider">admin@itx:~$</span>
        <span class="text-accent animate-pulse-glow">_</span>
      </div>
    </div>

    <div class="flex items-center gap-2">
      <span class="hidden sm:flex items-center gap-1.5 font-mono text-[10px] text-muted-foreground uppercase tracking-widest mr-2">
        <span class="w-1.5 h-1.5 rounded-full bg-accent shadow-[0_0_6px_hsl(var(--accent))]" />
        online
      </span>
      <ThemeToggle />
      <Button
        variant="ghost"
        size="sm"
        class="text-muted-foreground hover:text-foreground"
        @click="logout"
      >
        <LogOut class="h-4 w-4 lg:mr-2" />
        <span class="hidden lg:inline font-mono text-xs">exit</span>
      </Button>
    </div>
  </header>
</template>
