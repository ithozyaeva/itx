<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { requiredRule, useFormValidation } from '@/composables/useFormValidation'
import { memberService } from '@/services/memberService'
import { onMounted, ref } from 'vue'

const props = defineProps<{
  memberId: number | null
}>()

const emit = defineEmits(['saved', 'cancel'])
const isLoading = ref(false)

const validationRules = {
  name: [requiredRule],
  tg: [requiredRule],
  id: [],
}

// Используем composable для валидации формы
const { values, errors, touched, validate, handleBlur, isValid, resetForm } = useFormValidation({
  name: '',
  tg: '',
  id: props.memberId,
}, validationRules)

// Загрузка данных участника при редактировании
onMounted(async () => {
  if (props.memberId) {
    isLoading.value = true
    const member = await memberService.getById(props.memberId)
    if (member) {
      values.value.name = member.name
      values.value.tg = member.tg
    }
    isLoading.value = false
  }
})

// Обработчик отправки формы
async function handleSubmit(e: Event) {
  e.preventDefault()

  if (!validate()) {
    return
  }

  if (props.memberId) {
    // Обновление существующего участника
    await memberService.update(props.memberId, {
      id: props.memberId,
      name: values.value.name,
      tg: values.value.tg,
    })
  }
  else {
    // Создание нового участника
    await memberService.create({
      name: values.value.name,
      tg: values.value.tg,
    })
  }
  resetForm()
  emit('saved')
}

function handleCancel() {
  emit('cancel')
}
</script>

<template>
  <form class="space-y-6" @submit="handleSubmit">
    <div class="space-y-4">
      <div class="space-y-2">
        <Label for="name">Имя участника</Label>
        <Input
          id="name"
          v-model="values.name"
          placeholder="Введите имя участника"
          :class="{ 'border-red-500': touched.name && errors.name }"
          @blur="handleBlur('name')"
        />
        <p v-if="touched.name && errors.name" class="text-sm text-red-500">
          {{ errors.name }}
        </p>
      </div>

      <div class="space-y-2">
        <Label for="tg">Telegram</Label>
        <Input
          id="tg"
          v-model="values.tg"
          placeholder="Введите username в Telegram (например, @username)"
          :class="{ 'border-red-500': touched.tg && errors.tg }"
          @blur="handleBlur('tg')"
        />
        <p v-if="touched.tg && errors.tg" class="text-sm text-red-500">
          {{ errors.tg }}
        </p>
      </div>
    </div>

    <div class="flex justify-end space-x-2">
      <Button type="button" variant="outline" @click="handleCancel">
        Отмена
      </Button>
      <Button type="submit" :disabled="isLoading || !isValid">
        {{ isLoading ? 'Сохранение...' : props.memberId ? 'Сохранить' : 'Добавить' }}
      </Button>
    </div>
  </form>
</template>
