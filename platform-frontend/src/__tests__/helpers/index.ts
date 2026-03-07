import { createApp } from 'vue'

export function withSetup<T>(composable: () => T): { result: T, app: ReturnType<typeof createApp> } {
  let result!: T
  const app = createApp({
    setup() {
      result = composable()
      return () => {}
    },
  })
  app.mount(document.createElement('div'))
  return { result, app }
}
