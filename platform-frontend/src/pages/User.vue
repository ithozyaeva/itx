<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import ContactsForm from '@/components/Profile/ContactsForm.vue'
import MemberProfileCard from '@/components/Profile/MemberProfileForm.vue'
import MentorInfoForm from '@/components/Profile/MentorInfoForm.vue'
import NotificationSettingsForm from '@/components/Profile/NotificationSettingsForm.vue'
import ProfTagsForm from '@/components/Profile/ProfTagsForm.vue'
import ServicesForm from '@/components/Profile/ServicesForm.vue'
import { isUserMentor, useUser } from '@/composables/useUser'
import { profileService } from '@/services/profile'

const user = useUser()
const isLoading = ref(true)

onMounted(async () => {
  await profileService.getMe()
  isLoading.value = false
})

const isMentor = isUserMentor()
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
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
