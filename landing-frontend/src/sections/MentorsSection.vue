<script setup lang="ts">
import type { Mentor } from '@/components/ui/MentorsCard.vue'
import { Button, MentorCard, TagsGroup, Typography } from 'itx-ui-kit'

import { AlertCircle } from 'lucide-vue-next'
import { computed, onMounted, ref, watch } from 'vue'
import { useYandexMetrika } from 'yandex-metrika-vue3'

function reloadPage() {
  window.location.reload()
}

const mentors = ref<Mentor[]>([])
const selectedSpecialization = ref<string[]>(['Все'])
const specializations = ref<string[]>([])

const loading = ref(true)
const error = ref<string | null>(null)

const visibleCount = ref(6)
const yandexMetrika = useYandexMetrika()

function handleMentorClick(mentorName: string, link: string) {
  yandexMetrika.reachGoal('mentor_details_click', {
    mentor: mentorName,
  } as any)
  yandexMetrika.extLink(link, { title: `Ментор ${mentorName}` })
}

async function fetchMentors() {
  try {
    loading.value = true
    const response = await fetch(`/api/mentors`, {
      headers: {
        Accept: 'application/json',
      },
    })

    if (!response.ok) {
      throw new Error('Не удалось загрузить данные о менторах')
    }

    const data = await response.json()
    mentors.value = data.items ?? []
    // Собираем все уникальные теги
    const allTags = new Set<string>()
    data.items?.forEach((mentor: Mentor) => {
      mentor.profTags?.forEach((tag) => {
        allTags.add(tag.title)
      })
    })

    specializations.value = ['Все', ...Array.from(allTags)]
  }
  catch (err) {
    console.error('Ошибка при загрузке менторов:', err)
    error.value = 'Не удалось загрузить данные о менторах. Пожалуйста, попробуйте позже.'
  }
  finally {
    loading.value = false
  }
}

function showMore() {
  visibleCount.value += 6
}

watch(selectedSpecialization, (val, oldVal) => {
  if (val.includes('Все') && !oldVal.includes('Все')) {
    selectedSpecialization.value = ['Все']
    return
  }

  if (val.length === 0) {
    selectedSpecialization.value = ['Все']
    return
  }

  if (val.length > 1 && val.includes('Все')) {
    selectedSpecialization.value = val.filter(t => t !== 'Все')
  }
})

const filteredMentors = computed(() => {
  const tags = selectedSpecialization.value

  if (tags.includes('Все') || tags.length === 0) {
    return mentors.value
  }

  return mentors.value.filter(m =>
    m.profTags?.some(t => tags.includes(t.title)),
  )
})

const displayedMentors = computed(() =>
  filteredMentors.value.slice(0, visibleCount.value).map((mentor: Mentor) => ({
    id: mentor.id,
    avatar: `https://t.me/i/userpic/160/${mentor.tg}.jpg`,
    name: `${mentor.firstName} ${mentor.lastName}`,
    position: mentor.occupation,
    description: mentor.experience,
    labels: mentor.profTags.map(tag => tag.title),
    link: mentor.contacts.find(c => c.type === 1)?.link ?? '#',
  })),
)

onMounted(fetchMentors)
</script>

<template>
  <section
    id="mentors"
    class="w-full pt-16 md:pt-24 lg:pt-32"
  >
    <div class="container px-6 md:px-10">
      <div
        v-if="error && !mentors.length"
        class="text-center py-10"
      >
        <div class="flex justify-center mb-4">
          <AlertCircle class="h-12 w-12 text-red-500" />
        </div>
        <Typography
          variant="body-m"
          class="text-red-500 mb-2"
        >
          {{ error }}
        </Typography>
        <Button
          class="mt-4"
          @click="reloadPage"
        >
          Повторить
        </Button>
      </div>

      <div
        v-else
        class="flex flex-col gap-9"
      >
        <div class="flex flex-col items-center justify-center gap-8">
          <div class="flex flex-col items-center justify-center gap-5">
            <Typography
              variant="h2"
              as="h2"
              class="text-accent"
            >
              База менторов
            </Typography>
            <Typography
              variant="body-xl"
              as="p"
              class="text-center"
            >
              Эксперты нашего сообщества, которые готовы поделиться экспертизой
            </Typography>
          </div>
          <TagsGroup
            v-model="selectedSpecialization"
            :tags="specializations"
            :multiple="true"
            class="justify-center"
          />
        </div>

        <div
          v-if="loading"
          class="text-center py-16"
        >
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-accent mx-auto" />
          <Typography
            variant="body-l"
            as="p"
            class="mt-4 text-muted-foreground"
          >
            Загрузка менторов...
          </Typography>
        </div>

        <Typography
          v-else-if="filteredMentors.length === 0"
          variant="body-l"
          as="p"
          class="text-center py-16 text-muted-foreground"
        >
          Менторы не найдены
        </Typography>

        <div
          v-else
          class="flex flex-col items-center justify-center gap-9 lg:gap-12"
        >
          <div class="grid w-full grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-5">
            <div
              v-for="mentor in displayedMentors"
              :key="mentor.id"
              class="contents"
              @click="handleMentorClick(mentor.name, mentor.link)"
            >
              <MentorCard
                :avatar="mentor.avatar"
                :name="mentor.name"
                :position="mentor.position"
                :description="mentor.description"
                :labels="mentor.labels"
                :link="mentor.link"
              />
            </div>
          </div>
          <Button
            v-if="visibleCount < filteredMentors.length"
            variant="filled"
            as="button"
            class="w-fit self-center"
            @click="showMore"
          >
            Показать больше
          </Button>
        </div>
      </div>
    </div>
  </section>
</template>
