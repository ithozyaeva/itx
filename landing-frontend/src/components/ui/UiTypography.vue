<script setup lang="ts">
defineProps<{
  as?: string
  variant?: string
}>()
</script>

<template>
  <component
    :is="as || 'p'"
    :class="variant || 'body-l'"
  >
    <slot />
  </component>
</template>

<style scoped>
.h1 {
  font-family: var(--h1-font), sans-serif;
  font-weight: var(--h1-weight);
  font-size: var(--h1-size);
  line-height: var(--h1-line-height);
  letter-spacing: var(--h1-letter-spacing);
  text-transform: var(--h1-transform);
}

.h2 {
  font-family: var(--h2-font), sans-serif;
  font-weight: var(--h2-weight);
  font-size: var(--h2-size);
  line-height: var(--h2-line-height);
  letter-spacing: var(--h2-letter-spacing);
  text-transform: var(--h2-transform);
}

.h3 {
  font-family: var(--h3-font), sans-serif;
  font-weight: var(--h3-weight);
  font-size: var(--h3-size);
  line-height: var(--h3-line-height);
  letter-spacing: var(--h3-letter-spacing);
  text-transform: var(--h3-transform);
}

.h4 {
  font-family: var(--h4-font), sans-serif;
  font-weight: var(--h4-weight);
  font-size: var(--h4-size);
  line-height: var(--h4-line-height);
  letter-spacing: var(--h4-letter-spacing);
}

.price {
  font-family: var(--price-font), sans-serif;
  font-weight: var(--price-weight);
  font-size: var(--price-size);
  line-height: var(--price-line-height);
  letter-spacing: var(--price-letter-spacing);
}

.body-xl {
  font-family: var(--body-xl-font), sans-serif;
  font-weight: var(--body-xl-weight);
  font-size: var(--body-xl-size);
  line-height: var(--body-xl-line-height);
  letter-spacing: var(--body-xl-letter-spacing);
}

.body-l {
  font-family: var(--body-l-font), sans-serif;
  font-weight: var(--body-l-weight);
  font-size: var(--body-l-size);
  line-height: var(--body-l-line-height);
  letter-spacing: var(--body-l-letter-spacing);
}

.body-m {
  font-family: var(--body-m-font), sans-serif;
  font-weight: var(--body-m-weight);
  font-size: var(--body-m-size);
  line-height: var(--body-m-line-height);
  letter-spacing: var(--body-m-letter-spacing);
}

.body-s {
  font-family: var(--body-s-font), sans-serif;
  font-weight: var(--body-s-weight);
  font-size: var(--body-s-size);
  line-height: var(--body-s-line-height);
  letter-spacing: var(--body-s-letter-spacing);
}

.body-xs {
  font-family: var(--body-xs-font), sans-serif;
  font-weight: var(--body-xs-weight);
  font-size: var(--body-xs-size);
  line-height: var(--body-xs-line-height);
  letter-spacing: var(--body-xs-letter-spacing);
}

.title {
  font-family: var(--title-font), sans-serif;
  font-weight: var(--title-weight);
  font-size: var(--title-size);
  line-height: var(--title-line-height);
  letter-spacing: var(--title-letter-spacing);
}

.label {
  font-family: var(--label-font), sans-serif;
  font-weight: var(--label-weight);
  font-size: var(--label-size);
  line-height: var(--label-line-height);
  letter-spacing: var(--label-letter-spacing);
}

.button-text {
  font-family: var(--button-font), sans-serif;
  font-weight: var(--button-weight);
  font-size: var(--button-size);
  line-height: var(--button-line-height);
  letter-spacing: var(--button-letter-spacing);
  text-transform: var(--button-transform);
}

.name-text {
  font-family: var(--name-text-font), sans-serif;
  font-weight: var(--name-text-weight);
  font-size: var(--name-text-size);
  line-height: var(--name-text-line-height);
  letter-spacing: var(--name-text-letter-spacing);
}

.date {
  font-family: var(--date-font), sans-serif;
  font-weight: var(--date-weight);
  font-size: var(--date-size);
  line-height: var(--date-line-height);
  letter-spacing: var(--date-letter-spacing);
}
</style>
