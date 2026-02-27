<script setup lang="ts">
import ContactsForm from '@/components/Profile/ContactsForm.vue'
import MemberProfileCard from '@/components/Profile/MemberProfileForm.vue'
import MentorInfoForm from '@/components/Profile/MentorInfoForm.vue'
import ProfTagsForm from '@/components/Profile/ProfTagsForm.vue'
import ServicesForm from '@/components/Profile/ServicesForm.vue'
import { isUserMentor, useUser } from '@/composables/useUser'
import { profileSiervice } from '@/services/profile'
import { onMounted } from 'vue'

const user = useUser()

onMounted(() => {
  profileSiervice.getMe()
})

const isMentor = isUserMentor()
</script>

<template>
  <div class="min-h-screen pt-20 pb-10">
    <div
      class="container"
      :class="isMentor ? 'grid grid-cols-1 md:grid-cols-2 gap-4' : 'mx-auto px-4 max-w-3xl flex flex-col gap-4'"
    >
      <MemberProfileCard v-if="user" />
      <div v-if="isMentor" class="grid grid-cols-1 gap-4">
        <ProfTagsForm />
        <MentorInfoForm />
      </div>
      <ServicesForm v-if="isMentor" />
      <ContactsForm v-if="isMentor" />
    </div>
  </div>
</template>
