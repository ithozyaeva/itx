<script setup lang="ts">
import UiTag from './UiTag.vue'

const props = withDefaults(defineProps<{
  tags?: string[]
  modelValue?: string | string[]
  multiple?: boolean
}>(), {
  tags: () => [],
  modelValue: undefined,
  multiple: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string | string[]]
}>()

function handleClick(tag: string) {
  if (props.multiple) {
    const current = Array.isArray(props.modelValue) ? [...props.modelValue] : []
    const index = current.indexOf(tag)
    if (index > -1) {
      current.splice(index, 1)
    }
    else {
      current.push(tag)
    }
    emit('update:modelValue', current)
  }
  else {
    emit('update:modelValue', tag)
  }
}

function isActive(tag: string): boolean {
  if (props.multiple && Array.isArray(props.modelValue)) {
    return props.modelValue.includes(tag)
  }
  return props.modelValue === tag
}
</script>

<template>
  <div class="tags-group">
    <UiTag
      v-for="tag in tags"
      :key="tag"
      :variant="isActive(tag) ? 'active' : 'default'"
      @click="handleClick(tag)"
    >
      {{ tag }}
    </UiTag>
  </div>
</template>

<style scoped>
.tags-group {
  flex-wrap: wrap;
  gap: 10px;
  display: flex;
}
</style>
