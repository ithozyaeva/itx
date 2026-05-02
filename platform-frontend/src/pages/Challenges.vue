<script setup lang="ts">
// Stub-страница: заполняется на этапе 5 (челленджи).
import { Sword } from 'lucide-vue-next'
import { Typography } from '@/components/ui/typography'
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8 max-w-4xl">
    <div class="font-mono text-[11px] text-muted-foreground/60 tracking-wider mb-2">
      ~/community/challenges
    </div>
    <div class="flex items-center justify-between mb-6">
      <Typography variant="h2" as="h1">
        Челленджи
      </Typography>
    </div>
    <div class="rounded-sm border border-dashed bg-muted/20 p-12 text-center">
      <Sword class="h-10 w-10 text-muted-foreground/50 mx-auto mb-3" />
      <p class="text-sm text-muted-foreground">
        Скоро здесь появятся еженедельные и ежемесячные челленджи.
      </p>
    </div>
  </div>
</template>
