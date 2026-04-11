<script setup lang="ts">
import { Camera, Edit, Loader2, Star } from 'lucide-vue-next'
import { onMounted, reactive, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { useToast } from '@/components/ui/toast'
import { Typography } from '@/components/ui/typography'
import { useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { pointsService } from '@/services/points'
import { profileService } from '@/services/profile'

const { toast } = useToast()

const user = useUser()
const isEdit = ref<boolean>(false)
const isUploadingAvatar = ref(false)
const pointsBalance = ref<number | null>(null)

onMounted(async () => {
  try {
    const summary = await pointsService.getMyPoints()
    pointsBalance.value = summary.balance
  }
  catch (err) {
    handleError(err)
  }
})

const editedUser = reactive({
  firstName: user.value?.firstName,
  lastName: user.value?.lastName,
  birthday: user.value?.birthday,
  bio: user.value?.bio ?? '',
  grade: user.value?.grade ?? '',
  company: user.value?.company ?? '',
  tg: user.value?.tg ?? '',
})

const isSaving = ref(false)

async function handleSubmit() {
  isSaving.value = true
  try {
    await profileService.updateMe(editedUser)
    toast({ title: 'Профиль обновлён' })
    isEdit.value = false
  }
  catch (err) {
    handleError(err)
  }
  finally {
    isSaving.value = false
  }
}

const avatarSrc = ref(user.value?.avatarUrl || (user.value?.tg ? `https://t.me/i/userpic/160/${user.value.tg}.jpg` : ''))
const avatarError = ref(false)

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
  catch (err) {
    handleError(err)
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
            v-if="!avatarError"
            :src="avatarSrc"
            :alt="user?.firstName ?? 'Аватар'"
            class="w-full h-full object-cover"
            @error="avatarError = true"
          >
          <span
            v-else
            class="text-2xl font-bold text-muted-foreground"
          >{{ user?.firstName?.[0] }}{{ user?.lastName?.[0] }}</span>
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
          <div class="space-y-3 w-full">
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
              class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
            >
          </div>
        </template>
        <template v-if="!isEdit">
          <p
            v-if="user?.tg"
            class="text-sm text-muted-foreground mt-1"
          >
            @{{ user?.tg }}
          </p>
        </template>
        <template v-else>
          <input
            v-model="editedUser.tg"
            type="text"
            placeholder="Username в Telegram"
            class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring mt-3"
          >
        </template>
        <div
          v-if="pointsBalance !== null"
          class="flex items-center justify-center gap-1.5 mt-3"
        >
          <Star class="h-4 w-4 text-yellow-500" />
          <span class="font-medium text-yellow-500">{{ pointsBalance }} баллов</span>
        </div>
        <template v-if="!isEdit">
          <p
            v-if="user?.grade || user?.company"
            class="text-sm text-muted-foreground mt-2"
          >
            {{ [user?.grade, user?.company].filter(Boolean).join(' · ') }}
          </p>
          <p
            v-if="user?.bio"
            class="text-sm text-muted-foreground mt-2"
          >
            {{ user?.bio }}
          </p>
          <p
            v-if="user?.birthday"
            class="text-sm text-muted-foreground mt-2"
          >
            Дата рождения: {{ new Date(user?.birthday).toLocaleDateString() }}
          </p>
        </template>
        <template v-if="isEdit">
          <div class="flex gap-3 mt-3">
            <input
              v-model="editedUser.grade"
              type="text"
              placeholder="Грейд (Junior, Middle, Senior...)"
              class="w-1/2 px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
            >
            <input
              v-model="editedUser.company"
              type="text"
              placeholder="Место работы"
              class="w-1/2 px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
            >
          </div>
          <textarea
            v-model="editedUser.bio"
            placeholder="О себе..."
            rows="3"
            class="w-full mt-3 px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring resize-none"
          />
          <input
            v-model="editedUser.birthday"
            type="date"
            placeholder="Дата рождения"
            title="Дата рождения"
            class="w-full mt-3 px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring"
          >
        </template>
        <Button
          v-if="isEdit"
          class="mt-5 px-4 py-2 cursor-pointer transition duration-300"
          :disabled="isSaving"
          @click="handleSubmit"
        >
          <Loader2 v-if="isSaving" class="h-4 w-4 animate-spin mr-2" />
          Сохранить изменения
        </Button>
      </div>
    </div>
  </div>
</template>
