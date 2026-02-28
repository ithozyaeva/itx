<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { Camera, Edit, Loader2 } from 'lucide-vue-next'
import { reactive, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { useUser } from '@/composables/useUser'
import { profileService } from '@/services/profile'

const user = useUser()
const isEdit = ref<boolean>(false)
const isUploadingAvatar = ref(false)

const editedUser = reactive({
  firstName: user.value?.firstName,
  lastName: user.value?.lastName,
  birthday: user.value?.birthday,
  bio: user.value?.bio ?? '',
})

function handleSubmit() {
  profileService.updateMe(editedUser)
  isEdit.value = false
}

const avatarSrc = ref(user.value?.avatarUrl || `https://t.me/i/userpic/160/${user.value?.tg}.jpg`)

async function handleAvatarUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file)
    return

  isUploadingAvatar.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    const updated = await profileService.uploadAvatar(formData)
    if (updated?.avatarUrl) {
      avatarSrc.value = updated.avatarUrl
    }
  }
  finally {
    isUploadingAvatar.value = false
  }
}
</script>

<template>
  <div class="p-6 md:p-8 bg-card backdrop-blur-lg border border-border shadow-lg rounded-3xl">
    <div class="flex relative flex-col items-center space-y-6">
      <Edit
        class="absolute right-0 top-0 cursor-pointer text-muted-foreground hover:text-foreground"
        @click="isEdit = !isEdit"
      />
      <div class="flex justify-between">
        <div
          class="relative w-32 h-32 rounded-full border-4 border-border shadow-md overflow-hidden flex items-center justify-center bg-accent/20 group-hover:scale-105 transition-transform"
        >
          <img
            :src="avatarSrc"
            class="w-full h-full object-cover"
          >
          <label
            v-if="isEdit"
            class="absolute inset-0 flex items-center justify-center bg-black/40 cursor-pointer opacity-0 hover:opacity-100 transition-opacity"
          >
            <Loader2 v-if="isUploadingAvatar" class="h-6 w-6 animate-spin text-white" />
            <Camera v-else class="h-6 w-6 text-white" />
            <input
              type="file"
              accept="image/jpeg,image/png,image/webp"
              class="hidden"
              @change="handleAvatarUpload"
            >
          </label>
        </div>
      </div>
      <div class="text-center">
        <Typography v-if="!isEdit" variant="h2" as="h1">
          {{ user?.firstName }} {{ user?.lastName }}
        </Typography>
        <template v-else>
          <input
            v-model="editedUser.firstName"
            type="text"
            placeholder="Имя"
            class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
          >
          <input
            v-model="editedUser.lastName"
            type="text"
            placeholder="Фамилия"
            class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring mt-2"
          >
        </template>
        <p class="text-muted-foreground mb-4 mt-2">
          {{ user?.tg }}
        </p>
        <p v-if="!isEdit && !!user?.bio" class="text-sm text-muted-foreground mb-4">
          {{ user?.bio }}
        </p>
        <template v-if="isEdit">
          <textarea
            v-model="editedUser.bio"
            placeholder="О себе..."
            rows="3"
            class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring resize-none"
          />
        </template>
        <p v-if="!isEdit && !!user?.birthday" class="text-muted-foreground mb-4">
          Дата рождения: {{ new Date(user?.birthday).toLocaleDateString() }}
        </p>
        <template v-else-if="isEdit">
          <input
            v-model="editedUser.birthday"
            type="date"
            class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
          >
        </template>
        <Button
          v-if="isEdit"
          class="mt-5 px-4 py-2 cursor-pointer transition duration-300"
          @click="handleSubmit"
        >
          Сохранить изменения
        </Button>
      </div>
    </div>
  </div>
</template>
