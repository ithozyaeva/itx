<script setup lang="ts">
import type { AIMaterial, AIMaterialContentType, AIMaterialKind, CreateAIMaterialRequest } from '@/models/aiMaterial'
import { Loader2 } from 'lucide-vue-next'
import { computed, reactive, watch } from 'vue'
import FormField from '@/components/common/FormField.vue'
import {
  Dialog,
  DialogFooter,
  DialogHeader,
  DialogScrollContent,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  AI_MATERIAL_CONTENT_TYPE_OPTIONS,
  AI_MATERIAL_KIND_OPTIONS,
  AI_MATERIAL_LIMITS,
} from '@/models/aiMaterial'
import AIMaterialTagsInput from './AIMaterialTagsInput.vue'

const props = defineProps<{
  open: boolean
  initial?: AIMaterial | null
  isSubmitting?: boolean
}>()

const emit = defineEmits<{
  'update:open': [v: boolean]
  'submit': [v: CreateAIMaterialRequest]
}>()

const form = reactive<CreateAIMaterialRequest>({
  title: '',
  summary: '',
  contentType: 'prompt',
  materialKind: 'prompt',
  promptBody: '',
  externalUrl: '',
  agentConfig: '',
  tags: [],
})

const errors = reactive<Record<string, string>>({})

const HTTP_URL_RE = /^https?:\/\//i

watch(
  () => props.open,
  (open) => {
    if (open) {
      if (props.initial) {
        form.title = props.initial.title
        form.summary = props.initial.summary
        form.contentType = props.initial.contentType
        form.materialKind = props.initial.materialKind
        form.promptBody = props.initial.promptBody
        form.externalUrl = props.initial.externalUrl
        form.agentConfig = props.initial.agentConfig
        form.tags = [...props.initial.tags]
      }
      else {
        form.title = ''
        form.summary = ''
        form.contentType = 'prompt'
        form.materialKind = 'prompt'
        form.promptBody = ''
        form.externalUrl = ''
        form.agentConfig = ''
        form.tags = []
      }
      Object.keys(errors).forEach(k => delete errors[k])
    }
  },
)

const isEdit = computed(() => !!props.initial)

function setContentType(t: AIMaterialContentType) {
  form.contentType = t
}

function setKind(k: AIMaterialKind) {
  form.materialKind = k
}

function validate(): boolean {
  Object.keys(errors).forEach(k => delete errors[k])

  const title = form.title.trim()
  if ([...title].length < AI_MATERIAL_LIMITS.titleMin)
    errors.title = `Минимум ${AI_MATERIAL_LIMITS.titleMin} символа`
  else if ([...title].length > AI_MATERIAL_LIMITS.titleMax)
    errors.title = `Максимум ${AI_MATERIAL_LIMITS.titleMax} символов`

  const summary = form.summary.trim()
  if ([...summary].length < AI_MATERIAL_LIMITS.summaryMin)
    errors.summary = `Минимум ${AI_MATERIAL_LIMITS.summaryMin} символов — расскажите кому и зачем`
  else if ([...summary].length > AI_MATERIAL_LIMITS.summaryMax)
    errors.summary = `Максимум ${AI_MATERIAL_LIMITS.summaryMax} символов`

  if (form.contentType === 'prompt' && !form.promptBody?.trim())
    errors.promptBody = 'Содержимое промта обязательно'
  if (form.contentType === 'link') {
    const url = form.externalUrl?.trim() ?? ''
    if (!url)
      errors.externalUrl = 'Ссылка обязательна'
    else if (!HTTP_URL_RE.test(url))
      errors.externalUrl = 'Ссылка должна начинаться с http:// или https://'
  }
  if (form.contentType === 'agent' && !form.agentConfig?.trim())
    errors.agentConfig = 'Конфиг агента обязателен'

  return Object.keys(errors).length === 0
}

function submit() {
  if (!validate())
    return
  emit('submit', {
    title: form.title.trim(),
    summary: form.summary.trim(),
    contentType: form.contentType,
    materialKind: form.materialKind,
    promptBody: form.promptBody?.trim() ?? '',
    externalUrl: form.externalUrl?.trim() ?? '',
    agentConfig: form.agentConfig?.trim() ?? '',
    tags: form.tags,
  })
}
</script>

<template>
  <Dialog :open="open" @update:open="(v: boolean) => emit('update:open', v)">
    <DialogScrollContent class="max-w-2xl">
      <DialogHeader>
        <DialogTitle>{{ isEdit ? 'Редактировать материал' : 'Новый материал' }}</DialogTitle>
      </DialogHeader>

      <form class="space-y-4" @submit.prevent="submit">
        <FormField label="Название" :error="errors.title" required html-for="aim-title">
          <input
            id="aim-title"
            v-model="form.title"
            type="text"
            :maxlength="AI_MATERIAL_LIMITS.titleMax"
            class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="Например: Промт для код-ревью на TypeScript"
          >
        </FormField>

        <FormField label="Описание" :error="errors.summary" required html-for="aim-summary">
          <textarea
            id="aim-summary"
            v-model="form.summary"
            :maxlength="AI_MATERIAL_LIMITS.summaryMax"
            class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-20 resize-none"
            placeholder="Что это, кому подходит, какую задачу решает"
          />
          <p class="mt-1 text-xs text-muted-foreground">
            {{ form.summary.length }} / {{ AI_MATERIAL_LIMITS.summaryMax }}
          </p>
        </FormField>

        <FormField label="Категория" required>
          <div class="flex flex-wrap gap-2">
            <button
              v-for="opt in AI_MATERIAL_KIND_OPTIONS"
              :key="opt.value"
              type="button"
              class="px-3 py-1.5 rounded-sm text-sm font-medium border transition-colors"
              :class="form.materialKind === opt.value
                ? 'bg-primary text-primary-foreground border-primary'
                : 'bg-card border-border text-muted-foreground hover:text-foreground'"
              @click="setKind(opt.value)"
            >
              {{ opt.label }}
            </button>
          </div>
        </FormField>

        <FormField label="Тип контента" required>
          <div class="grid sm:grid-cols-3 gap-2">
            <button
              v-for="opt in AI_MATERIAL_CONTENT_TYPE_OPTIONS"
              :key="opt.value"
              type="button"
              class="text-left rounded-sm border p-3 transition-colors"
              :class="form.contentType === opt.value
                ? 'bg-primary/5 border-primary text-foreground'
                : 'bg-card border-border text-muted-foreground hover:text-foreground hover:border-accent'"
              @click="setContentType(opt.value)"
            >
              <div class="font-medium text-sm mb-0.5">
                {{ opt.label }}
              </div>
              <div class="text-xs">
                {{ opt.description }}
              </div>
            </button>
          </div>
        </FormField>

        <FormField
          v-if="form.contentType === 'prompt'"
          label="Текст промта"
          :error="errors.promptBody"
          required
          html-for="aim-prompt"
        >
          <textarea
            id="aim-prompt"
            v-model="form.promptBody"
            :maxlength="AI_MATERIAL_LIMITS.promptMax"
            class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary min-h-40 resize-y"
            placeholder="Текст промта или системной инструкции..."
          />
        </FormField>

        <FormField
          v-if="form.contentType === 'link'"
          label="Ссылка"
          :error="errors.externalUrl"
          required
          html-for="aim-url"
        >
          <input
            id="aim-url"
            v-model="form.externalUrl"
            type="url"
            :maxlength="AI_MATERIAL_LIMITS.urlMax"
            class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
            placeholder="https://github.com/..."
          >
        </FormField>

        <FormField
          v-if="form.contentType === 'agent'"
          label="Конфиг агента"
          :error="errors.agentConfig"
          required
          html-for="aim-agent"
        >
          <textarea
            id="aim-agent"
            v-model="form.agentConfig"
            :maxlength="AI_MATERIAL_LIMITS.agentMax"
            class="w-full rounded-sm border border-border bg-background px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-primary min-h-40 resize-y"
            placeholder="JSON или YAML конфиг агента"
          />
        </FormField>

        <FormField label="Теги">
          <AIMaterialTagsInput v-model="form.tags" />
        </FormField>

        <DialogFooter>
          <button
            type="submit"
            class="px-4 py-2 rounded-sm bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
            :disabled="isSubmitting"
          >
            <Loader2 v-if="isSubmitting" class="h-4 w-4 animate-spin inline mr-1" />
            {{ isEdit ? 'Сохранить' : 'Опубликовать' }}
          </button>
        </DialogFooter>
      </form>
    </DialogScrollContent>
  </Dialog>
</template>
