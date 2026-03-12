import { computed, ref } from 'vue'

export interface OnboardingStep {
  target: string
  title: string
  description: string
  placement: 'top' | 'bottom' | 'left' | 'right'
}

const STORAGE_KEY = 'onboarding_completed'

const isActive = ref(false)
const currentStepIndex = ref(0)

const steps: OnboardingStep[] = [
  {
    target: '[data-onboarding="sidebar"]',
    title: 'Навигация',
    description: 'Здесь находится главное меню платформы. Переходите между разделами одним нажатием.',
    placement: 'right',
  },
  {
    target: '[data-onboarding="dashboard"]',
    title: 'Ваш дашборд',
    description: 'На главной странице собрана вся ключевая информация: события, задания, достижения.',
    placement: 'bottom',
  },
  {
    target: '[data-onboarding="events"]',
    title: 'События',
    description: 'Следите за мероприятиями сообщества, записывайтесь и участвуйте.',
    placement: 'right',
  },
  {
    target: '[data-onboarding="points"]',
    title: 'Баллы и награды',
    description: 'Зарабатывайте баллы за активность, участвуйте в розыгрышах и сезонных рейтингах.',
    placement: 'right',
  },
  {
    target: '[data-onboarding="profile"]',
    title: 'Ваш профиль',
    description: 'Настройте профиль, укажите навыки и контакты, чтобы другие участники могли вас найти.',
    placement: 'top',
  },
]

export function useOnboarding() {
  const isCompleted = computed(() => localStorage.getItem(STORAGE_KEY) === 'true')
  const currentStep = computed(() => steps[currentStepIndex.value])
  const totalSteps = steps.length

  function start() {
    if (isCompleted.value)
      return
    currentStepIndex.value = 0
    isActive.value = true
  }

  function nextStep() {
    if (currentStepIndex.value < steps.length - 1) {
      currentStepIndex.value++
    }
    else {
      complete()
    }
  }

  function prevStep() {
    if (currentStepIndex.value > 0) {
      currentStepIndex.value--
    }
  }

  function skip() {
    complete()
  }

  function complete() {
    isActive.value = false
    localStorage.setItem(STORAGE_KEY, 'true')
  }

  function reset() {
    localStorage.removeItem(STORAGE_KEY)
    isActive.value = false
    currentStepIndex.value = 0
  }

  return {
    isActive,
    isCompleted,
    currentStep,
    currentStepIndex,
    totalSteps,
    start,
    nextStep,
    prevStep,
    skip,
    complete,
    reset,
  }
}
