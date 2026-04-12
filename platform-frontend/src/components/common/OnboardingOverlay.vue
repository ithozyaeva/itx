<script setup lang="ts">
import { ChevronLeft, ChevronRight, X } from 'lucide-vue-next'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { useOnboarding } from '@/composables/useOnboarding'

const {
  isActive,
  currentStep,
  currentStepIndex,
  totalSteps,
  nextStep,
  prevStep,
  skip,
} = useOnboarding()

const tooltipStyle = ref<Record<string, string>>({})
const highlightStyle = ref<Record<string, string>>({})
const arrowClass = ref('')

let retryTimer: ReturnType<typeof setTimeout> | null = null

function positionTooltip(retries = 5) {
  if (retryTimer) {
    clearTimeout(retryTimer)
    retryTimer = null
  }
  if (!currentStep.value)
    return

  const el = document.querySelector(currentStep.value.target)
  if (!el) {
    // Элемент ещё не отрендерен — повторяем с задержкой
    if (retries > 0) {
      retryTimer = setTimeout(positionTooltip, 200, retries - 1)
      return
    }
    // После всех попыток — пропускаем шаг
    if (currentStepIndex.value < totalSteps - 1)
      nextStep()
    else
      skip()
    return
  }

  const rect = el.getBoundingClientRect()

  // На мобилке элемент может быть за экраном (скрытый сайдбар) — пропускаем
  if (rect.right <= 0 || rect.left >= window.innerWidth) {
    if (currentStepIndex.value < totalSteps - 1)
      nextStep()
    else
      skip()
    return
  }
  const padding = 8
  const tooltipGap = 12

  highlightStyle.value = {
    top: `${rect.top - padding}px`,
    left: `${rect.left - padding}px`,
    width: `${rect.width + padding * 2}px`,
    height: `${rect.height + padding * 2}px`,
  }

  const placement = currentStep.value.placement
  const style: Record<string, string> = {}

  if (placement === 'right') {
    style.top = `${rect.top}px`
    style.left = `${rect.right + tooltipGap}px`
    arrowClass.value = 'arrow-left'
  }
  else if (placement === 'left') {
    style.top = `${rect.top}px`
    style.right = `${window.innerWidth - rect.left + tooltipGap}px`
    arrowClass.value = 'arrow-right'
  }
  else if (placement === 'bottom') {
    style.top = `${rect.bottom + tooltipGap}px`
    style.left = `${rect.left}px`
    arrowClass.value = 'arrow-top'
  }
  else {
    style.bottom = `${window.innerHeight - rect.top + tooltipGap}px`
    style.left = `${rect.left}px`
    arrowClass.value = 'arrow-bottom'
  }

  tooltipStyle.value = style
}

function repositionOnResize() {
  positionTooltip(0)
}

watch([currentStepIndex, isActive], () => {
  if (isActive.value) {
    nextTick(() => positionTooltip())
  }
})

onMounted(() => {
  window.addEventListener('resize', repositionOnResize)
  if (isActive.value) {
    nextTick(() => positionTooltip())
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', repositionOnResize)
  if (retryTimer)
    clearTimeout(retryTimer)
})
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isActive && currentStep"
      class="fixed inset-0 z-[100]"
    >
      <!-- Overlay -->
      <div
        class="absolute inset-0 bg-black/60 transition-opacity"
        @click="skip"
      />

      <!-- Highlight cutout -->
      <div
        class="absolute rounded-xl ring-2 ring-primary bg-transparent z-[101] pointer-events-none transition-all duration-300"
        :style="highlightStyle"
      />

      <!-- Tooltip -->
      <div
        class="absolute z-[102] w-80 bg-card border border-border rounded-sm p-5 shadow-xl transition-all duration-300"
        :class="arrowClass"
        :style="tooltipStyle"
      >
        <button
          class="absolute top-3 right-3 text-muted-foreground hover:text-foreground transition-colors"
          aria-label="Закрыть тур"
          @click="skip"
        >
          <X class="h-4 w-4" />
        </button>

        <h3 class="text-sm font-semibold mb-1">
          {{ currentStep.title }}
        </h3>
        <p class="text-sm text-muted-foreground mb-4">
          {{ currentStep.description }}
        </p>

        <!-- Progress dots -->
        <div class="flex items-center justify-between">
          <div class="flex gap-1.5">
            <div
              v-for="i in totalSteps"
              :key="i"
              class="h-1.5 w-1.5 rounded-full transition-colors"
              :class="i - 1 === currentStepIndex ? 'bg-accent' : 'bg-muted'"
            />
          </div>

          <div class="flex gap-2">
            <Button
              v-if="currentStepIndex > 0"
              variant="ghost"
              size="sm"
              @click="prevStep"
            >
              <ChevronLeft class="h-4 w-4" />
            </Button>
            <Button
              size="sm"
              @click="nextStep"
            >
              {{ currentStepIndex === totalSteps - 1 ? 'Готово' : 'Далее' }}
              <ChevronRight
                v-if="currentStepIndex < totalSteps - 1"
                class="h-4 w-4 ml-1"
              />
            </Button>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.arrow-left::before {
  content: '';
  position: absolute;
  left: -6px;
  top: 16px;
  width: 12px;
  height: 12px;
  background: hsl(var(--card));
  border-left: 1px solid hsl(var(--border));
  border-bottom: 1px solid hsl(var(--border));
  transform: rotate(45deg);
}

.arrow-right::before {
  content: '';
  position: absolute;
  right: -6px;
  top: 16px;
  width: 12px;
  height: 12px;
  background: hsl(var(--card));
  border-right: 1px solid hsl(var(--border));
  border-top: 1px solid hsl(var(--border));
  transform: rotate(45deg);
}

.arrow-top::before {
  content: '';
  position: absolute;
  top: -6px;
  left: 24px;
  width: 12px;
  height: 12px;
  background: hsl(var(--card));
  border-left: 1px solid hsl(var(--border));
  border-top: 1px solid hsl(var(--border));
  transform: rotate(45deg);
}

.arrow-bottom::before {
  content: '';
  position: absolute;
  bottom: -6px;
  left: 24px;
  width: 12px;
  height: 12px;
  background: hsl(var(--card));
  border-right: 1px solid hsl(var(--border));
  border-bottom: 1px solid hsl(var(--border));
  transform: rotate(45deg);
}
</style>
