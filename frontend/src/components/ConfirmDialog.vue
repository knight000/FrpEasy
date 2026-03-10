<template>
  <v-dialog v-model="dialog" max-width="400">
    <v-card>
      <v-card-title>{{ title }}</v-card-title>
      <v-card-text v-if="message">{{ message }}</v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="cancel">{{ cancelText }}</v-btn>
        <v-btn :color="confirmColor" variant="flat" @click="confirm">{{ confirmText }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = withDefaults(defineProps<{
  modelValue: boolean
  title: string
  message?: string
  confirmText?: string
  cancelText?: string
  confirmColor?: string
}>(), {
  confirmText: '确定',
  cancelText: '取消',
  confirmColor: 'primary'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

const dialog = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  dialog.value = val
})

watch(dialog, (val) => {
  emit('update:modelValue', val)
})

function cancel() {
  dialog.value = false
  emit('cancel')
}

function confirm() {
  dialog.value = false
  emit('confirm')
}
</script>
