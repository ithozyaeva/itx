<script setup lang="ts">
import { Button } from '@/components/ui/button'
import ThemeToggle from '@/components/ui/theme-toggle.vue'
import { useAuth } from '@/composables/useAuth'
import { LogOut, User } from 'lucide-vue-next'

const { logout } = useAuth()
</script>

<template>
  <header class="bg-card border-b border-border h-16 flex items-center px-6 justify-between">
    <div>
      <h1 class="text-xl font-bold">
        IT-ХОЗЯЕВА
      </h1>
    </div>

    <div class="flex items-center space-x-4">
      <ThemeToggle />
      <User class="h-5 w-5" />

      <Button @click="logout">
        <LogOut class="mr-2 h-4 w-4" />
        <span>Выйти</span>
      </Button>
    </div>
  </header>
</template>
