<script lang="ts" setup>
import type { CommunityEvent } from '@/models/events'
import { computed, onMounted, ref } from 'vue'
import EventHosts from '@/components/forms/event/EventHosts.vue'
import EventTags from '@/components/forms/event/EventTags.vue'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import { toast } from '@/components/ui/toast'
import { useDictionary } from '@/composables/useDictionary'
import { requiredRule, useFormValidation } from '@/composables/useFormValidation'
import { toDatetimeLocal } from '@/lib/utils'
import { eventsService } from '@/services/eventsService'

const props = defineProps<{
  eventId?: number
}>()

const emit = defineEmits(['saved', 'cancel'])
const isLoading = ref(false)

const validationRules = {
  title: [requiredRule],
  description: [],
  date: [requiredRule],
  timezone: [requiredRule],
  placeType: [requiredRule],
  eventType: [requiredRule],
  place: [requiredRule],
  open: [requiredRule],
  hosts: [requiredRule],
  customPlaceType: [],
  videoLink: [],
  id: [],
  eventTags: [],
}

// Популярные таймзоны в формате UTC
const timezones = [
  { value: 'UTC', label: 'UTC (UTC+0)' },
  { value: 'UTC+1', label: 'UTC+1' },
  { value: 'UTC+2', label: 'UTC+2' },
  { value: 'UTC+3', label: 'UTC+3 (Москва)' },
  { value: 'UTC+4', label: 'UTC+4' },
  { value: 'UTC+5', label: 'UTC+5' },
  { value: 'UTC+6', label: 'UTC+6' },
  { value: 'UTC+7', label: 'UTC+7' },
  { value: 'UTC+8', label: 'UTC+8' },
  { value: 'UTC+9', label: 'UTC+9' },
  { value: 'UTC+10', label: 'UTC+10' },
  { value: 'UTC+11', label: 'UTC+11' },
  { value: 'UTC+12', label: 'UTC+12' },
  { value: 'UTC-1', label: 'UTC-1' },
  { value: 'UTC-2', label: 'UTC-2' },
  { value: 'UTC-3', label: 'UTC-3' },
  { value: 'UTC-4', label: 'UTC-4' },
  { value: 'UTC-5', label: 'UTC-5' },
  { value: 'UTC-6', label: 'UTC-6' },
  { value: 'UTC-7', label: 'UTC-7' },
  { value: 'UTC-8', label: 'UTC-8' },
  { value: 'UTC-9', label: 'UTC-9' },
  { value: 'UTC-10', label: 'UTC-10' },
  { value: 'UTC-11', label: 'UTC-11' },
  { value: 'UTC-12', label: 'UTC-12' },
]

// Периоды повторения
const repeatPeriods = [
  { value: 'DAILY', label: 'Ежедневно' },
  { value: 'WEEKLY', label: 'Еженедельно' },
  { value: 'MONTHLY', label: 'Ежемесячно' },
  { value: 'YEARLY', label: 'Ежегодно' },
]

// Типы окончания повторений
const repeatEndTypes = [
  { value: 'never', label: 'Без окончания' },
  { value: 'date', label: 'По дате' },
]

// Используем composable для валидации формы
const { values, errors, touched, validate, isValid } = useFormValidation<CommunityEvent>({
  date: toDatetimeLocal(new Date().toISOString()),
  description: '',
  title: '',
  timezone: 'UTC',
  placeType: 'ONLINE',
  eventType: 'Созвон',
  place: '',
  open: false,
  hosts: [],
  id: 0,
  customPlaceType: '',
  videoLink: '',
  eventTags: [],
  isRepeating: false,
  repeatPeriod: undefined,
  repeatInterval: 1,
  repeatEndDate: undefined,
}, validationRules)

const repeatEndType = ref<'never' | 'date'>('never')

const repeatEndTypeModel = computed({
  get: () => repeatEndType.value,
  set: (value: string | null) => {
    const type = (value as 'never' | 'date') || 'never'
    repeatEndType.value = type
    if (type === 'never') {
      values.value.repeatEndDate = undefined
    }
  },
})

// Загрузка данных участника при редактировании
onMounted(async () => {
  if (props.eventId) {
    const event = await eventsService.getById(props.eventId)
    if (event) {
      values.value = {
        ...event,
        date: toDatetimeLocal(event.date),
        repeatEndDate: event.repeatEndDate ? toDatetimeLocal(event.repeatEndDate) : undefined,
      }
      // Устанавливаем тип окончания повторений
      if (event.repeatEndDate) {
        repeatEndType.value = 'date'
      }
      else {
        repeatEndType.value = 'never'
      }
    }
  }
})

// Обработчик отправки формы
async function handleSubmit(e: Event) {
  e.preventDefault()

  if (!validate()) {
    return
  }

  try {
    const eventData: any = {
      ...values.value,
      date: new Date(values.value.date).toISOString(),
    }

    // Очищаем поля повторений, если событие не повторяющееся
    if (!values.value.isRepeating) {
      eventData.repeatPeriod = null
      eventData.repeatInterval = null
      eventData.repeatEndDate = null
    }
    else {
      // Очищаем неиспользуемые поля окончания
      if (repeatEndType.value === 'never') {
        eventData.repeatEndDate = null
      }
      else if (repeatEndType.value === 'date' && values.value.repeatEndDate) {
        eventData.repeatEndDate = new Date(values.value.repeatEndDate).toISOString()
      }
    }

    if (props.eventId) {
      await eventsService.update(props.eventId, eventData)
    }
    else {
      await eventsService.create(eventData)
    }

    emit('saved')
  }
  catch (err) {
    console.error(err)
    toast({
      variant: 'destructive',
      title: 'Не удалось сохранить событие.',
    })
  }
}

function handleCancel() {
  emit('cancel')
}

const { placeTypes } = useDictionary(['placeTypes'])
</script>

<template>
  <form class="space-y-6" @submit="handleSubmit">
    <div class="space-y-4">
      <div class="space-y-2">
        <Label for="title">Название встречи</Label>
        <Input id="title" v-model="values.title" type="text" placeholder="Название встречи" />
        <p v-if="touched.title && errors.title" class="text-sm text-destructive">
          {{ errors.title }}
        </p>
      </div>
      <div class="space-y-2">
        <Label for="description">Описание встречи</Label>
        <Textarea id="description" v-model="values.description" placeholder="Описание встречи" />
      </div>
      <div class="space-y-2">
        <Label for="even">Теги встречи</Label>
        <EventTags v-model:event-tags="values.eventTags" />
      </div>
      <div class="space-y-2">
        <Label for="date">Дата встречи</Label>
        <Input id="date" v-model="values.date" type="datetime-local" placeholder="Введите дату встречи" />
        <p v-if="touched.date && errors.date" class="text-sm text-destructive">
          {{ errors.date }}
        </p>
      </div>
      <div class="space-y-2">
        <Label for="timezone">Таймзона</Label>
        <Select id="timezone" v-model="values.timezone">
          <SelectTrigger>
            <SelectValue placeholder="Выберите таймзону" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="tz in timezones" :key="tz.value" :value="tz.value">
              {{ tz.label }}
            </SelectItem>
          </SelectContent>
        </Select>
        <p v-if="touched.timezone && errors.timezone" class="text-sm text-red-500">
          {{ errors.timezone }}
        </p>
      </div>
      <div class="space-y-2">
        <Label for="placeType">Тип проведения встречи</Label>
        <Select id="placeType" v-model="values.placeType">
          <SelectTrigger>
            <SelectValue placeholder="Тип проведения встречи" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem v-for="option in placeTypes" :key="option.value" :value="option.value">
              {{ option.label }}
            </SelectItem>
          </SelectContent>
        </Select>
        <p v-if="touched.placeType && errors.placeType" class="text-sm text-destructive">
          {{ errors.placeType }}
        </p>
      </div>
      <div v-if="values.placeType !== 'ONLINE'" class="space-y-2">
        <Label for="customPlaceType">Тип проведения встречи (дополнительно)</Label>
        <Input id="customPlaceType" v-model="values.customPlaceType" type="text" placeholder="Укажите город" />
      </div>
      <div class="space-y-2">
        <Label for="place">Место проведения встречи</Label>
        <Textarea id="place" v-model="values.place" placeholder="Ссылка/адрес" />
        <p v-if="touched.place && errors.place" class="text-sm text-destructive">
          {{ errors.place }}
        </p>
      </div>
      <div class="space-y-2">
        <Label for="eventType">Тип встречи</Label>
        <Input
          id="eventType" v-model="values.eventType" type="text"
          placeholder="Тип встречи: созвон/митап/воркшоп и тд"
        />
        <p v-if="touched.eventType && errors.eventType" class="text-sm text-destructive">
          {{ errors.eventType }}
        </p>
      </div>
      <div class="flex items-center space-x-2">
        <Label for="open">Открытая встреча (не только для участников сообщества)</Label>
        <Checkbox id="open" v-model="values.open" />
      </div>
      <div class="space-y-2">
        <Label for="videoLink">Ссылка на запись</Label>
        <Input
          id="videoLink" v-model="values.videoLink" type="text"
          placeholder="Ссылка на ютуб/облако"
        />
      </div>
      <div class="space-y-2">
        <Label for="hosts">Спикеры встречи</Label>
        <EventHosts v-model:hosts="values.hosts" />
      </div>
      <div class="space-y-4 border-t pt-4">
        <div class="flex items-center space-x-2">
          <Checkbox id="isRepeating" v-model="values.isRepeating" :checked="values.isRepeating" />
          <Label for="isRepeating">Повторяющееся событие</Label>
        </div>
        <div v-if="values.isRepeating" class="space-y-4 pl-6 border-l-2">
          <div class="space-y-2">
            <Label for="repeatPeriod">Период повторения</Label>
            <Select id="repeatPeriod" v-model="values.repeatPeriod">
              <SelectTrigger>
                <SelectValue placeholder="Выберите период" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="period in repeatPeriods" :key="period.value" :value="period.value">
                  {{ period.label }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="space-y-2">
            <Label for="repeatInterval">Интервал</Label>
            <Input
              id="repeatInterval"
              v-model.number="values.repeatInterval"
              type="number"
              min="1"
              placeholder="Каждые N периодов"
            />
            <p class="text-xs text-muted-foreground">
              Например: 2 = каждые 2 дня/недели/месяца/года
            </p>
          </div>
          <div class="space-y-2">
            <Label>Окончание повторений</Label>
            <Select v-model="repeatEndTypeModel">
              <SelectTrigger>
                <SelectValue placeholder="Выберите тип окончания" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="endType in repeatEndTypes" :key="endType.value" :value="endType.value">
                  {{ endType.label }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div v-if="repeatEndType === 'date'" class="space-y-2">
            <Label for="repeatEndDate">Дата окончания</Label>
            <Input
              id="repeatEndDate"
              v-model="values.repeatEndDate"
              type="datetime-local"
              placeholder="Дата окончания повторений"
            />
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-end space-x-2">
      <Button type="button" variant="outline" @click="handleCancel">
        Отмена
      </Button>
      <Button type="submit" :disabled="isLoading || !isValid">
        {{ isLoading ? 'Сохранение...' : props.eventId ? 'Сохранить' : 'Добавить' }}
      </Button>
    </div>
  </form>
</template>
