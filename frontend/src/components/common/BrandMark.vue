<script setup lang="ts">
withDefaults(
  defineProps<{
    size?: string;
  }>(),
  {
    size: '1em',
  },
);
</script>

<template>
  <span
    class="brand-mark"
    :style="{ width: size, height: size }"
    aria-hidden="true"
  ></span>
</template>

<style scoped>
.brand-mark {
  display: inline-block;
  flex: 0 0 auto;
  background-color: currentColor;
  -webkit-mask-image: url('/assets/logo.svg');
  mask-image: url('/assets/logo.svg');
  -webkit-mask-position: center;
  mask-position: center;
  -webkit-mask-repeat: no-repeat;
  mask-repeat: no-repeat;
  -webkit-mask-size: contain;
  mask-size: contain;
}
</style>
