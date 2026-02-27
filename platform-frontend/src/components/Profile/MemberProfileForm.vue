<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { useUser } from '@/composables/useUser'
import { profileService } from '@/services/profile'
import { Typography } from 'itx-ui-kit'
import { Edit } from 'lucide-vue-next'
import { reactive, ref } from 'vue'

const user = useUser()
const isEdit = ref<boolean>(false)

const editedUser = reactive({
  firstName: user.value?.firstName,
  lastName: user.value?.lastName,
  birthday: user.value?.birthday,
})

function handleSubmit() {
  profileService.updateMe(editedUser)
  isEdit.value = false
}
</script>

<template>
  <div class="p-6 md:p-8 bg-card backdrop-blur-lg border border-border shadow-lg rounded-3xl">
    <div class="flex relative flex-col items-center space-y-6">
      <Edit class="absolute right-0 top-0 cursor-pointer text-muted-foreground hover:text-foreground" @click="isEdit = !isEdit" />
      <div class="flex justify-between">
        <div
          class="w-32 h-32 rounded-full border-4 border-border shadow-md overflow-hidden flex items-center justify-center bg-accent/20 group-hover:scale-105 transition-transform"
        >
          <img :src="`https://t.me/i/userpic/160/${user.tg}.jpg`">
        </div>
      </div>
      <div class="text-center">
        <Typography v-if="!isEdit" variant="h2" as="h1">
          {{ user?.firstName }} {{ user?.lastName }}
        </Typography>
        <template v-else>
          <input v-model="editedUser.firstName" type="text" placeholder="Имя" class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring">
          <input v-model="editedUser.lastName" type="text" placeholder="Фамилия" class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring mt-2">
        </template>
        <p class="text-muted-foreground mb-4 mt-2">
          {{ user?.tg }}
        </p>
        <p v-if="!isEdit && !!user.birthday" class="text-muted-foreground mb-4">
          Дата рождения: {{ new Date(user?.birthday).toLocaleDateString() }}
        </p>
        <template v-else-if="isEdit">
          <input v-model="editedUser.birthday" type="date" class="w-full px-4 py-2 border border-input rounded-xl bg-transparent focus:outline-none focus:ring-2 focus:ring-ring">
        </template>
        <Button v-if="isEdit" class="mt-5 px-4 py-2 cursor-pointer transition duration-300" @click="handleSubmit">
          Сохранить изменения
        </Button>
      </div>
    </div>
  </div>
</template>
