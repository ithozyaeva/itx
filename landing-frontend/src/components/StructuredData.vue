<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue'

interface Props {
  data: object | object[]
}

const props = defineProps<Props>()

function createJsonLdScripts() {
  const dataArray = Array.isArray(props.data) ? props.data : [props.data]

  return dataArray.map((item) => {
    const script = document.createElement('script')
    script.type = 'application/ld+json'
    script.textContent = JSON.stringify(item)
    document.head.appendChild(script)
    return script
  })
}

let scripts: HTMLScriptElement[] = []

onMounted(() => {
  scripts = createJsonLdScripts()
})

watch(() => props.data, () => {
  scripts.forEach(script => script.remove())
  scripts = createJsonLdScripts()
}, { deep: true })

onUnmounted(() => {
  scripts.forEach(script => script.remove())
})
</script>

<template>
  <div style="display: none" />
</template>
