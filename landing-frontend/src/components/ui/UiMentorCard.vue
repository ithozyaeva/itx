<script setup lang="ts">
import { nextTick, onMounted, ref } from 'vue'
import AvatarPlaceholderIcon from './AvatarPlaceholderIcon.vue'
import UiButton from './UiButton.vue'
import UiLabel from './UiLabel.vue'
import UiTypography from './UiTypography.vue'

const props = withDefaults(defineProps<{
  avatar: string
  name: string
  position: string
  labels?: string[]
  description?: string
  link?: string
}>(), {
  labels: () => [],
  description: '',
  link: '#',
})

const containerRef = ref<HTMLElement | null>(null)
const visibleLabels = ref<string[]>([])
const hasMore = ref(false)
const imgError = ref(false)

function handleError() {
  imgError.value = true
}

function handleLoad(e: Event) {
  const target = e.target as HTMLImageElement
  if (target.naturalWidth <= 1 && target.naturalHeight <= 1) {
    imgError.value = true
  }
}

onMounted(async () => {
  if (!containerRef.value || props.labels.length === 0)
    return

  visibleLabels.value = [...props.labels]
  hasMore.value = false
  await nextTick()

  const label = containerRef.value.querySelector('.ui-label')
  if (!label)
    return

  const maxHeight = 2 * (label as HTMLElement).offsetHeight + 8
  let currentHeight = containerRef.value.offsetHeight

  while (currentHeight > maxHeight && visibleLabels.value.length > 1) {
    visibleLabels.value.pop()
    hasMore.value = true
    await nextTick()
    currentHeight = containerRef.value!.offsetHeight
  }
})
</script>

<template>
  <div class="mentor-card">
    <div>
      <div class="top-section">
        <div class="avatar">
          <img
            v-if="!imgError"
            :src="avatar"
            alt="Mentor avatar"
            fetchpriority="high"
            width="160"
            height="160"
            decoding="async"
            @error="handleError"
            @load="handleLoad"
          >
          <AvatarPlaceholderIcon v-else />
        </div>
        <div>
          <UiTypography
            variant="h4"
            class="name"
          >
            {{ name }}
          </UiTypography>
          <UiTypography
            variant="title"
            class="position"
          >
            {{ position.length ? position : 'Не указано' }}
          </UiTypography>
        </div>
      </div>
      <div class="bottom-section">
        <div
          ref="containerRef"
          class="labels"
        >
          <UiLabel
            v-for="label in visibleLabels"
            :key="label"
          >
            {{ label }}
          </UiLabel>
          <UiLabel v-if="hasMore">
            ...
          </UiLabel>
        </div>
        <UiTypography
          variant="body-xs"
          class="description"
        >
          {{ description }}
        </UiTypography>
      </div>
    </div>
    <UiButton
      as="a"
      :href="link"
      target="_blank"
      rel="noopener noreferrer"
      variant="stroke"
      class="button"
    >
      Подробнее
    </UiButton>
  </div>
</template>

<style scoped>
.mentor-card {
  box-sizing: border-box;
  background-color: hsl(151 5% 10% / 0.85);
  border: 1px solid hsl(151 5% 18%);
  border-radius: 2px;
  flex-direction: column;
  justify-content: space-between;
  align-items: center;
  gap: 25px;
  width: 100%;
  padding: 24px 24px 24px;
  display: flex;
  position: relative;
  transition: border-color 250ms ease, transform 250ms ease, background-color 250ms ease;
}

.mentor-card::before,
.mentor-card::after {
  content: '';
  position: absolute;
  width: 10px;
  height: 10px;
  border-color: var(--color-green-700);
  pointer-events: none;
  opacity: 0;
  transition: opacity 250ms ease;
}
.mentor-card::before {
  top: -1px; left: -1px;
  border-top: 2px solid;
  border-left: 2px solid;
}
.mentor-card::after {
  bottom: -1px; right: -1px;
  border-bottom: 2px solid;
  border-right: 2px solid;
}
.mentor-card:hover {
  border-color: hsl(151 60% 54% / 0.5);
  background-color: hsl(151 5% 12% / 0.95);
}
.mentor-card:hover::before,
.mentor-card:hover::after {
  opacity: 1;
}

.top-section {
  text-align: center;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  display: flex;
}

.avatar {
  border-radius: 2px;
  flex-shrink: 0;
  width: 106px;
  height: 106px;
  overflow: hidden;
  border: 1px solid hsl(151 5% 22%);
  outline: 1px solid hsl(151 60% 54% / 0.2);
  outline-offset: 3px;
}

.avatar img {
  object-fit: cover;
  width: 100%;
  height: 100%;
}

.name {
  color: var(--color-green-700);
  margin-bottom: 4px;
}

.position {
  color: var(--color-white);
}

.bottom-section {
  flex-direction: column;
  gap: 20px;
  display: flex;
}

.labels {
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px 5px;
  max-width: 100%;
  display: flex;
}

.description {
  text-align: center;
  color: var(--color-grey);
  overflow-wrap: break-word;
  word-break: break-word;
  white-space: normal;
}

.button {
  width: 100%;
  min-height: 50px;
}
</style>
