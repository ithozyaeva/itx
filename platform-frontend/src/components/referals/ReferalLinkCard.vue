<script setup lang="ts">
import type { ReferalLink } from '@/models/referals'
import type { PropType } from 'vue'
import ReferalLinkForm from '@/components/referals/ReferalLinkForm.vue'
import { useDictionary } from '@/composables/useDictionary'
import { useUser } from '@/composables/useUser'
import { dateFormatter } from '@/lib/utils'
import { referalLinkService } from '@/services/referals'
import { Typography } from 'itx-ui-kit'
import { Pencil, Trash } from 'lucide-vue-next'
import { computed, ref } from 'vue'

const props = defineProps({
  link: {
    type: Object as PropType<ReferalLink>,
    required: true,
  },
})

const emit = defineEmits(['updated', 'deleted'])

const user = useUser()
const isEditing = ref(false)

const isOwner = computed(() => user.value?.id === props.link.author.id)

function startEditing() {
  isEditing.value = true
}

async function saveEdit(editedLink: Partial<ReferalLink>) {
  try {
    const response = await referalLinkService.updateLink({ ...props.link, ...editedLink })
    emit('updated', response)
    isEditing.value = false
  }
  catch (error) {
    console.error('Ошибка при сохранении реферальной ссылки:', error)
  }
}

function cancelEdit() {
  isEditing.value = false
}

async function handleDelete() {
  try {
    await referalLinkService.deleteLink(props.link.id)
    emit('deleted', props.link.id)
  }
  catch (error) {
    console.error('Ошибка при удалении реферальной ссылки:', error)
  }
}

const { gradesObject } = useDictionary(['grades'])
</script>

<template>
  <div
    data-reveal
    class="bg-card rounded-3xl border p-4 hover:shadow-md transition-shadow flex flex-col gap-1"
  >
    <!-- Режим просмотра -->
    <div v-if="!isEditing">
      <div class="flex justify-between items-start mb-3">
        <Typography variant="h4" as="h3">
          {{ link.company }}
        </Typography>
        <div class="space-x-2">
          <button v-if="isOwner" class="p-1 -mt-1 rounded hover:bg-secondary cursor-pointer" @click="startEditing">
            <Pencil :size="16" />
          </button>
          <button v-if="isOwner" class="p-1 -mt-1 rounded hover:bg-secondary cursor-pointer" @click="handleDelete">
            <Trash :size="16" />
          </button>
        </div>
      </div>
      <p class="text-sm text-muted-foreground">
        Автор: <a :href="`https://t.me/${link.author.tg}`" target="_blank" class="underline">{{ link.author.firstName }} {{ link.author.lastName }}</a>
      </p>
      <div class="space-y-1 text-sm">
        <div class="space-x-2">
          <span class="font-bold">Грейд:</span>
          <span> {{ gradesObject[link.grade] }}</span>
        </div>
        <div class="space-x-2">
          <span class="font-bold">Навыки:</span>
          <span> {{ link.profTags.map(tag => tag.title).join(', ') }}</span>
        </div>
        <div class="space-x-2">
          <span class="font-bold">Количество вакансий:</span>
          <span> {{ link.vacationsCount }}</span>
        </div>
        <div class="space-x-2">
          <span class="font-bold">Обновлено:</span>
          <span> {{ dateFormatter.format(new Date(link.updatedAt)) }}</span>
        </div>
      </div>
    </div>

    <ReferalLinkForm
      v-if="isEditing"
      :link="link"
      title="Редактировать ссылку"
      @save="saveEdit"
      @cancel="cancelEdit"
    />
  </div>
</template>
