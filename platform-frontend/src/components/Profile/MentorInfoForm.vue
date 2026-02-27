<script setup lang="ts">
import type { Mentor } from '@/models/profile'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { useUser } from '@/composables/useUser'
import { profileSiervice } from '@/services/profile'
import { Edit } from 'lucide-vue-next'
import { reactive, ref, watchEffect } from 'vue'

const user = useUser<Mentor>()
const isEdit = ref<boolean>(false)

const editedUser = reactive({
  occupation: user.value?.occupation,
  experience: user.value?.experience,
})

watchEffect(() => {
  editedUser.occupation = user.value?.occupation
  editedUser.experience = user.value?.experience
})

function handleSubmit() {
  profileSiervice.updateMentorInfo(editedUser)
  isEdit.value = false
}
</script>

<template>
  <div class="p-6 md:p-8 bg-card backdrop-blur-lg border border-border shadow-lg rounded-3xl w-full">
    <div class="flex relative flex-col items-center space-y-6">
      <Edit class="absolute right-0 top-0 cursor-pointer text-muted-foreground hover:text-foreground" @click="isEdit = !isEdit" />
      <div class="space-y-2 w-full mt-2">
        <Label for="occupation">Специализация</Label>
        <Input
          id="occupation"
          v-model="editedUser.occupation"
          placeholder="фронт/бек/другое"
          :readonly="!isEdit"
          max-length="255"
        />
      </div>
      <div class="space-y-2 w-full">
        <Label for="experience">Опыт работы</Label>
        <Textarea
          id="experience"
          v-model="editedUser.experience"
          placeholder="Кто Вы, где работали, кем. Желательно вкратце"
          rows="4"
          :readonly="!isEdit"
          max-length="255"
        />
      </div>
      <Button v-if="isEdit" class="mt-1 px-4 py-2 cursor-pointer transition duration-300" @click="handleSubmit">
        Сохранить изменения
      </Button>
    </div>
  </div>
</template>
