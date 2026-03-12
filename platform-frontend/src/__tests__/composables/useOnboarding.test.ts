import { beforeEach, describe, expect, it, vi } from 'vitest'

describe('useOnboarding', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.resetModules()
  })

  async function getUseOnboarding() {
    const mod = await import('@/composables/useOnboarding')
    return mod.useOnboarding
  }

  it('isCompleted is false when localStorage has no key', async () => {
    const useOnboarding = await getUseOnboarding()
    const { isCompleted } = useOnboarding()

    expect(isCompleted.value).toBe(false)
  })

  it('isCompleted is true when localStorage has onboarding_completed = true', async () => {
    localStorage.setItem('onboarding_completed', 'true')

    const useOnboarding = await getUseOnboarding()
    const { isCompleted } = useOnboarding()

    expect(isCompleted.value).toBe(true)
  })

  it('start sets isActive to true and currentStepIndex to 0', async () => {
    const useOnboarding = await getUseOnboarding()
    const { isActive, currentStepIndex, start } = useOnboarding()

    start()

    expect(isActive.value).toBe(true)
    expect(currentStepIndex.value).toBe(0)
  })

  it('start does nothing when already completed', async () => {
    localStorage.setItem('onboarding_completed', 'true')

    const useOnboarding = await getUseOnboarding()
    const { isActive, start } = useOnboarding()

    start()

    expect(isActive.value).toBe(false)
  })

  it('nextStep increments currentStepIndex', async () => {
    const useOnboarding = await getUseOnboarding()
    const { currentStepIndex, start, nextStep } = useOnboarding()

    start()
    nextStep()

    expect(currentStepIndex.value).toBe(1)
  })

  it('nextStep calls complete when on last step', async () => {
    const useOnboarding = await getUseOnboarding()
    const { isActive, currentStepIndex, start, nextStep, totalSteps } = useOnboarding()

    start()

    // Navigate to last step
    for (let i = 0; i < totalSteps - 1; i++) {
      nextStep()
    }

    expect(currentStepIndex.value).toBe(totalSteps - 1)

    // One more nextStep should call complete
    nextStep()

    expect(isActive.value).toBe(false)
    expect(localStorage.getItem('onboarding_completed')).toBe('true')
  })

  it('prevStep decrements currentStepIndex', async () => {
    const useOnboarding = await getUseOnboarding()
    const { currentStepIndex, start, nextStep, prevStep } = useOnboarding()

    start()
    nextStep()
    nextStep()
    expect(currentStepIndex.value).toBe(2)

    prevStep()
    expect(currentStepIndex.value).toBe(1)
  })

  it('prevStep does nothing when on step 0', async () => {
    const useOnboarding = await getUseOnboarding()
    const { currentStepIndex, start, prevStep } = useOnboarding()

    start()
    expect(currentStepIndex.value).toBe(0)

    prevStep()
    expect(currentStepIndex.value).toBe(0)
  })

  it('skip calls complete', async () => {
    const useOnboarding = await getUseOnboarding()
    const { isActive, start, skip } = useOnboarding()

    start()
    expect(isActive.value).toBe(true)

    skip()

    expect(isActive.value).toBe(false)
    expect(localStorage.getItem('onboarding_completed')).toBe('true')
  })

  it('complete sets isActive to false and writes to localStorage', async () => {
    const useOnboarding = await getUseOnboarding()
    const { isActive, start, complete } = useOnboarding()

    start()
    complete()

    expect(isActive.value).toBe(false)
    expect(localStorage.getItem('onboarding_completed')).toBe('true')
  })

  it('reset clears localStorage and resets state', async () => {
    const useOnboarding = await getUseOnboarding()
    const { isActive, currentStepIndex, isCompleted, start, nextStep, complete, reset } = useOnboarding()

    start()
    nextStep()
    complete()

    expect(isCompleted.value).toBe(true)

    reset()

    expect(isActive.value).toBe(false)
    expect(currentStepIndex.value).toBe(0)
    expect(localStorage.getItem('onboarding_completed')).toBeNull()
    expect(isCompleted.value).toBe(false)
  })

  it('currentStep returns the correct step object', async () => {
    const useOnboarding = await getUseOnboarding()
    const { currentStep, start, nextStep } = useOnboarding()

    start()

    expect(currentStep.value.target).toBe('[data-onboarding="sidebar"]')
    expect(currentStep.value.title).toBe('Навигация')

    nextStep()

    expect(currentStep.value.target).toBe('[data-onboarding="dashboard"]')
    expect(currentStep.value.title).toBe('Ваш дашборд')
  })

  it('totalSteps returns 5', async () => {
    const useOnboarding = await getUseOnboarding()
    const { totalSteps } = useOnboarding()

    expect(totalSteps).toBe(5)
  })
})
