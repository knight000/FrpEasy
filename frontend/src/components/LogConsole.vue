<template>
  <v-card class="log-card" theme="dark">
    <v-toolbar color="#2d2d2d" density="compact" flat>
      <v-toolbar-title class="text-body-2 text-grey-lighten-1">
        {{ title || '控制台输出' }}
      </v-toolbar-title>
      <v-chip v-if="autoScroll" color="success" size="x-small" class="mr-2">
        自动滚动
      </v-chip>
      <v-spacer />
      <v-tooltip text="切换自动滚动" location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            density="compact"
            :icon="autoScroll ? 'mdi-lock' : 'mdi-lock-open-variant'"
            size="small"
            variant="text"
            @click="autoScroll = !autoScroll"
          />
        </template>
      </v-tooltip>
      <v-tooltip text="复制日志" location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            density="compact"
            icon="mdi-content-copy"
            size="small"
            variant="text"
            @click="copyLogs"
          />
        </template>
      </v-tooltip>
      <v-tooltip text="清空日志" location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            density="compact"
            icon="mdi-delete-outline"
            size="small"
            variant="text"
            @click="emit('clear')"
          />
        </template>
      </v-tooltip>
    </v-toolbar>
    <div ref="scrollContainer" class="log-console" @scroll="onScroll">
      <div v-for="item in logs" :key="item.id" class="log-line">
        <span class="log-time">[{{ formatTime(item.timestamp) }}]</span>
        <span
          :class="{
            'log-error': item.type === 'error',
            'log-warn': item.type === 'warn',
            'log-info': item.type === 'info',
          }"
          class="log-type"
        >
          {{ item.type.toUpperCase() }}
        </span>
        <span class="log-message">{{ item.message }}</span>
      </div>
      <div v-if="logs.length === 0" class="log-empty">
        暂无日志
      </div>
    </div>
  </v-card>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue'
import type { LogEntry } from '@/stores/preset'

const props = defineProps<{
  logs: LogEntry[]
  title?: string
}>()

const emit = defineEmits<{
  clear: []
}>()

const autoScroll = ref(true)
const scrollContainer = ref<HTMLDivElement | null>(null)
let isUserScrolling = false
let scrollTimer: ReturnType<typeof setTimeout> | null = null

function onScroll() {
  if (!scrollContainer.value) return
  const el = scrollContainer.value
  const isAtBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 30
  
  if (!isAtBottom) {
    isUserScrolling = true
    if (scrollTimer) clearTimeout(scrollTimer)
    scrollTimer = setTimeout(() => {
      isUserScrolling = false
    }, 2000)
  } else {
    isUserScrolling = false
  }
}

function scrollToBottom() {
  if (!scrollContainer.value) return
  const el = scrollContainer.value
  el.scrollTop = el.scrollHeight
}

function doAutoScroll() {
  if (autoScroll.value && !isUserScrolling) {
    nextTick(() => {
      scrollToBottom()
    })
  }
}

watch(
  () => props.logs.length,
  () => {
    doAutoScroll()
  }
)

onMounted(() => {
  scrollToBottom()
})

onUnmounted(() => {
  if (scrollTimer) clearTimeout(scrollTimer)
})

function formatTime(timestamp: number): string {
  const d = new Date(timestamp)
  return d.toLocaleTimeString('zh-CN', { hour12: false })
}

function copyLogs() {
  const text = props.logs
    .map((l) => `[${formatTime(l.timestamp)}] ${l.type.toUpperCase()} ${l.message}`)
    .join('\n')
  navigator.clipboard.writeText(text)
}
</script>

<style scoped>
.log-card {
  display: flex;
  flex-direction: column;
  height: 280px;
  background-color: #1e1e1e;
}

.log-console {
  height: calc(100% - 48px);
  overflow-y: auto;
  overflow-x: hidden;
  background-color: #1e1e1e;
  color: #d4d4d4;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.85rem;
  padding: 8px;
}

.log-console::-webkit-scrollbar {
  width: 8px;
}

.log-console::-webkit-scrollbar-track {
  background: #2d2d2d;
  border-radius: 4px;
}

.log-console::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 4px;
}

.log-console::-webkit-scrollbar-thumb:hover {
  background: #666;
}

.log-line {
  padding: 2px 8px;
  white-space: pre-wrap;
  word-break: break-all;
  line-height: 1.5;
}

.log-time {
  color: #888;
}

.log-type {
  margin: 0 8px;
  font-weight: 500;
}

.log-error {
  color: #f44336;
}

.log-warn {
  color: #ff9800;
}

.log-info {
  color: #4caf50;
}

.log-debug {
  color: #2196f3;
}

.log-message {
  color: #d4d4d4;
}

.log-empty {
  text-align: center;
  color: #888;
  padding: 32px;
}
</style>
