<script setup lang="ts">
import type { AcceptableValue } from 'reka-ui'
import type { MemberRole } from '@/models/members'
import { computed, onMounted, ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectLabel, SelectSeparator, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useDictionary } from '@/composables/useDictionary'
import { requiredArrayRule, requiredRule, useFormValidation } from '@/composables/useFormValidation'
import { memberService } from '@/services/memberService'

const props = defineProps<{
  memberId: number | null
}>()

const emit = defineEmits(['saved', 'cancel'])
const isLoading = ref(false)

const validationRules = {
  firstName: [requiredRule],
  lastName: [],
  birthday: [],
  roles: [requiredArrayRule],
  tg: [requiredRule],
  id: [],
}

// Используем composable для валидации формы
const { values, errors, touched, validate, handleBlur, isValid, resetForm } = useFormValidation({
  firstName: '',
  lastName: '',
  tg: '',
  roles: ['UNSUBSCRIBER'] as MemberRole[],
  id: props.memberId,
  birthday: '',
  tgId: 0,
}, validationRules)

// Загрузка данных участника при редактировании
onMounted(async () => {
  if (props.memberId) {
    isLoading.value = true
    const member = await memberService.getById(props.memberId)
    if (member) {
      values.value.firstName = member.firstName
      values.value.lastName = member.lastName
      values.value.tg = member.tg
      values.value.tgId = member.telegramID
      values.value.birthday = member.birthday
      values.value.roles = member.roles || []
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
      firstName: values.value.firstName,
      lastName: values.value.lastName,
      tg: values.value.tg,
      roles: values.value.roles,
      telegramID: values.value.tgId,
      birthday: values.value.birthday,
    })
  }
  else {
    // Создание нового участника
    await memberService.create({
      firstName: values.value.firstName,
      lastName: values.value.lastName,
      tg: values.value.tg,
      roles: values.value.roles,
      telegramID: values.value.tgId,
      birthday: values.value.birthday,
    })
  }
  resetForm()
  emit('saved')
}

function handleCancel() {
  emit('cancel')
}

const { memberRoles, memberRolesObject } = useDictionary<MemberRole>(['memberRoles'])
function handleRolesChange(value: AcceptableValue) {
  const typedValue = value as MemberRole[]
  let newValue = [...typedValue]
  if (typedValue.includes('UNSUBSCRIBER') && values.value.roles.includes('SUBSCRIBER')) {
    newValue = typedValue.filter(item => item !== 'SUBSCRIBER')
  }
  else if (typedValue.includes('SUBSCRIBER') && values.value.roles.includes('UNSUBSCRIBER')) {
    newValue = typedValue.filter(item => item !== 'UNSUBSCRIBER')
  }

  values.value.roles = newValue
}
const fixedRoles: MemberRole[] = ['MENTOR']
const userFixedRoles = computed(() => values.value.roles.filter(item => fixedRoles.includes(item)).map(item => memberRolesObject.value[item]))
</script>

<template>
  <form class="space-y-6" @submit="handleSubmit">
    <div class="space-y-4">
      <div class="space-y-2">
        <Label for="name">Имя участника</Label>
        <Input
          id="firstName"
          v-model="values.firstName"
          placeholder="Введите имя участника"
          :class="{ 'border-destructive': touched.firstName && errors.firstName }"
          @blur="handleBlur('firstName')"
        />
        <p v-if="touched.firstName && errors.firstName" class="text-sm text-destructive">
          {{ errors.firstName }}
        </p>
      </div>

      <div class="space-y-2">
        <Label for="name">Фамилия участника</Label>
        <Input
          id="lastName"
          v-model="values.lastName"
          placeholder="Введите имя участника"
          :class="{ 'border-destructive': touched.lastName && errors.lastName }"
          @blur="handleBlur('lastName')"
        />
        <p v-if="touched.firstName && errors.lastName" class="text-sm text-destructive">
          {{ errors.lastName }}
        </p>
      </div>
      <div class="space-y-2">
        <Label for="name">Дата рождения участника</Label>
        <Input
          id="birthDate"
          v-model="values.birthday"
          type="date"
          placeholder="Введите ДР"
          :class="{ 'border-destructive': touched.lastName && errors.birthday }"
          @blur="handleBlur('birthday')"
        />
        <p v-if="touched.firstName && errors.birthday" class="text-sm text-destructive">
          {{ errors.birthday }}
        </p>
      </div>

      <div class="space-y-2">
        <Label for="roles">Роли участника</Label>
        <Select id="roles" :model-value="values.roles" multiple @update:model-value="handleRolesChange">
          <SelectTrigger>
            <SelectValue placeholder="Выберите роли участника" />
          </SelectTrigger>
          <SelectContent>
            <div v-if="userFixedRoles.length">
              <SelectLabel>Автоопределяющие роли:</SelectLabel>
              <div class="px-2 py-1.5 text-sm text-muted-foreground">
                {{ userFixedRoles.join(', ') }}
              </div>

              <SelectSeparator />
            </div>
            <SelectLabel>Доступные роли:</SelectLabel>
            <SelectItem
              v-for="role in memberRoles.filter((item) => !fixedRoles.includes(item.value))"
              :key="role.value"
              :value="role.value"
            >
              {{ role.label }}
            </SelectItem>
          </SelectContent>
        </Select>
        <p v-if="touched.roles && errors.roles" class="text-sm text-destructive">
          {{ errors.roles }}
        </p>
      </div>

      <div class="space-y-2">
        <Label for="tg">Telegram</Label>
        <Input
          id="tg"
          v-model="values.tg"
          placeholder="Введите username в Telegram (например, @username)"
          :class="{ 'border-destructive': touched.tg && errors.tg }"
          @blur="handleBlur('tg')"
        />
        <p v-if="touched.tg && errors.tg" class="text-sm text-destructive">
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
