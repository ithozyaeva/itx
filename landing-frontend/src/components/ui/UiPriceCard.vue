<script setup lang="ts">
import { computed } from 'vue'
import UiButton from './UiButton.vue'
import UiTypography from './UiTypography.vue'

const props = defineProps<{
  name: string
  price: number
  oldPrice?: number
  features: string[]
  link: string
  variant?: 'default' | 'highlighted'
}>()

const discount = computed(() => {
  if (!props.oldPrice)
    return null
  const old = +props.oldPrice
  const current = +props.price
  if (old <= 0 || current >= old)
    return null
  return Math.round((old - current) * 100 / old)
})
</script>

<template>
  <div
    class="price-card"
    :class="variant || 'default'"
  >
    <div class="content">
      <div class="header">
        <UiTypography
          variant="h3"
          as="h3"
          class="name"
        >
          {{ name }}
        </UiTypography>
        <div class="price-container">
          <UiTypography
            v-if="discount"
            as="p"
            variant="title"
            class="discount"
          >
            &minus;{{ discount }}%
          </UiTypography>
        </div>
      </div>
      <ul class="features-container">
        <li
          v-for="(feature, index) in features"
          :key="index"
          class="feature"
        >
          <UiTypography variant="body-s">
            {{ feature }}
          </UiTypography>
        </li>
      </ul>
    </div>
    <div class="footer">
      <div class="price-container">
        <div class="prices">
          <UiTypography
            v-if="oldPrice"
            variant="price"
            class="old-price"
          >
            {{ oldPrice }}&thinsp;&#8381;
          </UiTypography>
          <UiTypography
            variant="price"
            class="price"
          >
            {{ price }}&thinsp;&#8381;
          </UiTypography>
        </div>
        <UiTypography
          variant="label"
          class="period"
        >
          в месяц
        </UiTypography>
      </div>
      <UiButton
        as="a"
        :href="link"
        target="_blank"
        variant="filled"
        class="button"
      >
        Подписаться
      </UiButton>
    </div>
  </div>
</template>

<style scoped>
.price-card {
  box-sizing: border-box;
  border-radius: var(--radius-card);
  flex-direction: column;
  justify-content: space-between;
  gap: 32px;
  min-width: 300px;
  min-height: 360px;
  padding: 28px 24px 24px;
  display: flex;
}

.content {
  text-align: center;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 24px;
  display: flex;
}

.header {
  text-align: center;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  display: flex;
}

.features-container {
  text-align: start;
  flex-direction: column;
  align-items: start;
  gap: 8px;
  max-width: 400px;
  padding-left: 20px;
  list-style: outside;
  display: flex;
}

.price-container {
  flex-direction: column;
  align-items: center;
  display: flex;
}

.prices {
  align-items: center;
  gap: 6px;
  margin: 12px 0 8px;
  display: flex;
}

.prices .old-price {
  text-decoration: line-through;
}

.discount {
  border-radius: var(--radius-default);
  border: 1px solid;
  padding: 4px 8px;
}

.footer {
  flex-direction: column;
  align-items: center;
  gap: 20px;
  display: flex;
}

.footer .button {
  justify-content: center;
  align-items: center;
  width: fit-content;
  display: flex;
}

/* Default variant */
.price-card.default {
  background-color: var(--color-green-black-500);
}

.price-card.default .name {
  color: var(--color-green-700);
}

.price-card.default .discount {
  color: var(--color-white);
  border-color: var(--color-white);
}

.price-card.default .price {
  color: var(--color-white);
}

.price-card.default .old-price {
  color: var(--color-white);
  opacity: 0.1;
}

.price-card.default .period {
  color: var(--color-white);
  opacity: 0.4;
}

.price-card.default .feature {
  color: var(--color-white);
}

/* Highlighted variant */
.price-card.highlighted {
  background-color: var(--color-green-700);
}

.price-card.highlighted .name {
  color: var(--color-green-black-700);
}

.price-card.highlighted .discount {
  color: var(--color-green-black-700);
  border-color: var(--color-green-black-700);
}

.price-card.highlighted .price {
  color: var(--color-green-black-700);
}

.price-card.highlighted .old-price {
  color: var(--color-green-black-700);
  opacity: 0.3;
}

.price-card.highlighted .period {
  color: var(--color-green-black-700);
  opacity: 0.4;
}

.price-card.highlighted .feature {
  color: var(--color-green-black-700);
}
</style>
