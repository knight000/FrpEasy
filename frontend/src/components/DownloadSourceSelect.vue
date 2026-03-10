<template>
  <div>
    <div class="text-subtitle-2 mb-2">{{ label }}</div>
    <v-select
      :model-value="modelValue"
      :items="downloadSources"
      item-title="label"
      item-value="value"
      density="compact"
      variant="outlined"
      hide-details
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { useDownloadSource, type DownloadSourceValue } from '@/composables/useDownloadSource'

defineProps<{
  modelValue: DownloadSourceValue
  label?: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: DownloadSourceValue): void
}>()

const { downloadSources } = useDownloadSource()
</script>
