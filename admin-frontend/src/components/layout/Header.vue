<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { LogOut, User } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import ThemeToggle from '@/components/ui/theme-toggle.vue'
import { useAuth } from '@/composables/useAuth'

const { logout } = useAuth()
</script>

<template>
  <header class="sticky top-0 z-40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 border-b border-border h-16 flex items-center px-6 justify-between">
    <div>
      <Typography variant="h4" as="h1">
        IT-ХОЗЯЕВА
      </Typography>
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
