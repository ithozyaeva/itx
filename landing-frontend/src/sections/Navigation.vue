<script setup lang="ts">
import { MenuItem } from 'itx-ui-kit'

const menuItems = [
  { label: 'менторы', href: '#mentors' },
  { label: 'встречи', href: '#meets' },
  { label: 'тарифы', href: '#tariffs' },
  { label: 'вопросы', href: '#FAQ' },
]
</script>

<template>
  <nav class="flex gap-2 lg:gap-10">
    <MenuItem
      v-for="item in menuItems"
      :key="item.href"
      :href="item.href"
      variant="default"
      as="a"
      class="capitalize border-4 border-transparent"
    >
      {{ item.label }}
    </MenuItem>
  </nav>
</template>
