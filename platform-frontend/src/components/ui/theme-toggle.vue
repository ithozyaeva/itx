<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { useColorMode } from '@vueuse/core'
import { Laptop, Moon, Sun } from 'lucide-vue-next'

// Используем useColorMode из @vueuse/core и устанавливаем светлую тему по умолчанию
const mode = useColorMode({
  initialValue: 'light',
})
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="outline" size="icon">
        <Sun v-if="mode === 'light'" class="h-[1.2rem] w-[1.2rem]" />
        <Moon v-else-if="mode === 'dark'" class="h-[1.2rem] w-[1.2rem]" />
        <Laptop v-else class="h-[1.2rem] w-[1.2rem]" />
        <span class="sr-only">Переключить тему</span>
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuItem @click="mode = 'light'">
        <Sun class="mr-2 h-4 w-4" />
        <span>Светлая</span>
      </DropdownMenuItem>
      <DropdownMenuItem @click="mode = 'dark'">
        <Moon class="mr-2 h-4 w-4" />
        <span>Темная</span>
      </DropdownMenuItem>
      <DropdownMenuItem @click="mode = 'auto'">
        <Laptop class="mr-2 h-4 w-4" />
        <span>Системная</span>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
