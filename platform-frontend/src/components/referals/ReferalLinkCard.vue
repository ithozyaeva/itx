<script setup lang="ts">
import type { PropType } from 'vue'
import type { ReferalLink } from '@/models/referals'
import { Typography } from 'itx-ui-kit'
import { Loader2, Pencil, Trash } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import ReferalLinkForm from '@/components/referals/ReferalLinkForm.vue'
import { Badge } from '@/components/ui/badge'
import { useDictionary } from '@/composables/useDictionary'
import { useUser } from '@/composables/useUser'
import { dateFormatter } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { referalLinkService } from '@/services/referals'

const props = defineProps({
  link: {
    type: Object as PropType<ReferalLink>,
    required: true,
  },
})

const emit = defineEmits(['updated', 'deleted'])

const user = useUser()
const isEditing = ref(false)
const isSaving = ref(false)
const isDeleting = ref(false)

const isOwner = computed(() => user.value?.id === props.link.author.id)

function startEditing() {
  isEditing.value = true
}

async function saveEdit(editedLink: Partial<ReferalLink>) {
  isSaving.value = true
  try {
    const response = await referalLinkService.updateLink({ ...props.link, ...editedLink })
    emit('updated', response)
    isEditing.value = false
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSaving.value = false
  }
}

function cancelEdit() {
  isEditing.value = false
}

async function handleDelete() {
  isDeleting.value = true
  try {
    await referalLinkService.deleteLink(props.link.id)
    emit('deleted', props.link.id)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isDeleting.value = false
  }
}

const { gradesObject, referalLinkStatusesObject } = useDictionary(['grades', 'referalLinkStatuses'])
</script>

<template>
  <div
    data-reveal
    class="bg-card rounded-3xl border p-4 hover:shadow-md transition-shadow flex flex-col gap-1"
  >
    <!-- Режим просмотра -->
    <div v-if="!isEditing">
      <div class="flex justify-between items-start mb-3">
        <div class="flex flex-wrap items-center gap-2 min-w-0">
          <Typography variant="h4" as="h3">
            {{ link.company }}
          </Typography>
          <Badge :variant="link.status === 'active' ? 'default' : 'secondary'">
            {{ referalLinkStatusesObject[link.status] }}
          </Badge>
        </div>
        <div v-if="isOwner" class="flex items-center gap-1 shrink-0 ml-2">
          <button class="p-1.5 rounded-lg hover:bg-secondary cursor-pointer text-muted-foreground hover:text-foreground transition-colors" :disabled="isSaving" @click="startEditing">
            <Pencil :size="14" />
          </button>
          <ConfirmDialog
            title="Удалить ссылку?"
            description="Реферальная ссылка будет удалена без возможности восстановления."
            confirm-label="Удалить"
            @confirm="handleDelete"
          >
            <template #trigger>
              <button class="p-1.5 rounded-lg hover:bg-destructive/10 cursor-pointer text-muted-foreground hover:text-destructive transition-colors" :disabled="isDeleting">
                <Loader2 v-if="isDeleting" :size="14" class="animate-spin" />
                <Trash v-else :size="14" />
              </button>
            </template>
          </ConfirmDialog>
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
        <div v-if="link.expiresAt" class="space-x-2">
          <span class="font-bold">Срок действия до:</span>
          <span> {{ dateFormatter.format(new Date(link.expiresAt)) }}</span>
        </div>
        <div v-if="link.conversionsCount > 0" class="space-x-2">
          <span class="font-bold">Конверсии:</span>
          <span>{{ link.conversionsCount }}</span>
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
      :is-saving="isSaving"
      title="Редактировать ссылку"
      @save="saveEdit"
      @cancel="cancelEdit"
    />
  </div>
</template>
