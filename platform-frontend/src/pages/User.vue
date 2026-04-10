<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { BarChart3, Bell, FileText, Folder, Loader2, MessageSquare, RefreshCw, ShoppingBag } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import ContactsForm from '@/components/Profile/ContactsForm.vue'
import MemberProfileCard from '@/components/Profile/MemberProfileForm.vue'
import MentorInfoForm from '@/components/Profile/MentorInfoForm.vue'
import NotificationSettingsForm from '@/components/Profile/NotificationSettingsForm.vue'
import ProfTagsForm from '@/components/Profile/ProfTagsForm.vue'
import ServicesForm from '@/components/Profile/ServicesForm.vue'
import { isUserMentor, useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { profileService } from '@/services/profile'

const user = useUser()
const isLoading = ref(true)
const loadError = ref(false)

onMounted(async () => {
  try {
    await profileService.getMe()
  }
  catch (err) {
    handleError(err)
    loadError.value = true
  }
  finally {
    isLoading.value = false
  }
})

const isMentor = isUserMentor()

const quickLinks = [
  { title: 'Моя статистика', path: '/my-stats', icon: BarChart3 },
  { title: 'Рефералки', path: '/referals', icon: Folder },
  { title: 'Резюме', path: '/resumes', icon: FileText },
  { title: 'Мои отзывы', path: '/my-reviews', icon: MessageSquare },
  { title: 'Уведомления', path: '/notifications', icon: Bell },
  { title: 'Барахолка', path: '/marketplace', icon: ShoppingBag },
]
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <Typography variant="h2" as="h1" class="mb-6">
      Мой профиль
    </Typography>

    <!-- Quick links hub -->
    <div class="grid grid-cols-3 sm:grid-cols-6 gap-2 sm:gap-3 mb-6">
      <router-link
        v-for="link in quickLinks"
        :key="link.path"
        :to="link.path"
        class="flex flex-col items-center gap-1.5 rounded-2xl border bg-card p-3 sm:p-4 hover:bg-muted/50 hover:border-accent/30 transition-all text-center"
      >
        <component :is="link.icon" class="h-5 w-5 text-muted-foreground" />
        <span class="text-[10px] sm:text-xs text-muted-foreground leading-tight">{{ link.title }}</span>
      </router-link>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>
    <div
      v-else-if="loadError"
      class="flex flex-col items-center justify-center gap-3 py-12 text-muted-foreground"
    >
      <p>Не удалось загрузить профиль</p>
      <button
        class="flex items-center gap-2 px-4 py-2 rounded-xl border hover:bg-muted transition-colors text-sm"
        @click="loadError = false; isLoading = true; profileService.getMe().then(() => { isLoading = false }).catch((err) => { handleError(err); loadError = true; isLoading = false })"
      >
        <RefreshCw class="h-4 w-4" />
        Повторить
      </button>
    </div>
    <div
      v-else-if="isMentor"
      class="grid grid-cols-1 md:grid-cols-2 gap-4"
    >
      <MemberProfileCard v-if="user" />
      <div class="flex flex-col gap-4">
        <ProfTagsForm />
        <MentorInfoForm />
      </div>
      <ServicesForm />
      <ContactsForm />
      <NotificationSettingsForm />
    </div>
    <div
      v-else
      class="max-w-2xl mx-auto flex flex-col gap-4"
    >
      <MemberProfileCard v-if="user" />
      <NotificationSettingsForm />
    </div>
  </div>
</template>
