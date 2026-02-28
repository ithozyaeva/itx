<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { LogOut, Menu, User } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import ThemeToggle from '@/components/ui/theme-toggle.vue'
import { useAuth } from '@/composables/useAuth'
import { useSidebar } from '@/composables/useSidebar'

const { logout } = useAuth()
const { toggleMobileSidebar } = useSidebar()
</script>

<template>
  <header class="sticky top-0 z-30 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 border-b border-border h-14 lg:h-16 flex items-center px-4 lg:px-6 justify-between">
    <div class="flex items-center gap-3">
      <Button
        variant="ghost"
        size="icon"
        class="lg:hidden"
        @click="toggleMobileSidebar"
      >
        <Menu class="h-5 w-5" />
      </Button>
      <Typography variant="h4" as="h1" class="text-base lg:text-lg">
        IT-ХОЗЯЕВА
      </Typography>
    </div>

    <div class="flex items-center space-x-2 lg:space-x-4">
      <ThemeToggle />
      <User class="h-5 w-5 hidden sm:block" />

      <Button size="sm" class="lg:size-default" @click="logout">
        <LogOut class="h-4 w-4 lg:mr-2" />
        <span class="hidden lg:inline">Выйти</span>
      </Button>
    </div>
  </header>
</template>
