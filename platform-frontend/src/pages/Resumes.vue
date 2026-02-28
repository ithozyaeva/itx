<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import ResumeUploadForm from '@/components/Resume/ResumeUploadForm.vue'
</script>

<template>
  <div class="min-h-screen pt-20 pb-10">
    <div class="container mx-auto px-4 max-w-4xl">
      <Typography variant="h2" as="h1" class="mb-6">
        Мои резюме
      </Typography>
      <ResumeUploadForm />
    </div>
  </div>
</template>
