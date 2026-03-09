<template>
  <v-card :class="['log-card', { 'log-card--fullscreen': fullscreen }]" theme="dark">
    <v-toolbar color="#2d2d2d" density="compact" flat>
      <v-toolbar-title class="text-body-2 text-grey-lighten-1">
        {{ title || '控制台输出' }}
      </v-toolbar-title>
      <v-chip v-if="autoScroll" color="success" size="x-small" class="mr-2">
        自动滚动
      </v-chip>
      <v-spacer />
      
      <v-text-field
        v-if="fullscreen"
        v-model="searchQuery"
        density="compact"
        variant="outlined"
        hide-details
        placeholder="搜索日志..."
        prepend-inner-icon="mdi-magnify"
        clearable
        class="search-input mr-2"
      />
      
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
      <v-tooltip v-if="!fullscreen" text="全屏" location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            density="compact"
            icon="mdi-fullscreen"
            size="small"
            variant="text"
            @click="emit('toggleFullscreen')"
          />
        </template>
      </v-tooltip>
      <v-tooltip v-else text="退出全屏" location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            density="compact"
            icon="mdi-fullscreen-exit"
            size="small"
            variant="text"
            @click="emit('toggleFullscreen')"
          />
        </template>
      </v-tooltip>
    </v-toolbar>
    <div ref="scrollContainer" class="log-console" @scroll="onScroll" @contextmenu.prevent="onContextMenu">
      <div v-for="item in filteredLogs" :key="item.id" class="log-line">
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
      <div v-if="filteredLogs.length === 0" class="log-empty">
        {{ searchQuery ? '无匹配结果' : '暂无日志' }}
      </div>
    </div>

    <v-card
      v-if="showContextMenu"
      class="context-menu"
      :style="{ left: contextMenuX + 'px', top: contextMenuY + 'px' }"
      @click.stop
    >
      <v-list density="compact" bg-color="#2d2d2d">
        <v-list-item @click="handleCopySelected">
          <template #prepend>
            <v-icon size="small">mdi-content-copy</v-icon>
          </template>
          <v-list-item-title>复制选中</v-list-item-title>
        </v-list-item>
        <v-list-item @click="handleCopyAll">
          <template #prepend>
            <v-icon size="small">mdi-content-multiple</v-icon>
          </template>
          <v-list-item-title>复制全部</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-card>
  </v-card>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import type { LogEntry } from '@/stores/preset'

const props = defineProps<{
  logs: LogEntry[]
  title?: string
  fullscreen?: boolean
}>()

const emit = defineEmits<{
  clear: []
  toggleFullscreen: []
}>()

const autoScroll = ref(true)
const scrollContainer = ref<HTMLDivElement | null>(null)
const searchQuery = ref('')
const showContextMenu = ref(false)
const contextMenuX = ref(0)
const contextMenuY = ref(0)
const selectedText = ref('')
let isUserScrolling = false
let scrollTimer: ReturnType<typeof setTimeout> | null = null

const filteredLogs = computed(() => {
  if (!searchQuery.value) {
    return props.logs
  }
  const query = searchQuery.value.toLowerCase()
  return props.logs.filter(log => 
    log.message.toLowerCase().includes(query) ||
    log.type.toLowerCase().includes(query)
  )
})

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
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  if (scrollTimer) clearTimeout(scrollTimer)
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})

function formatTime(timestamp: number): string {
  const d = new Date(timestamp)
  return d.toLocaleTimeString('zh-CN', { hour12: false })
}

function copyLogs() {
  const logsToCopy = searchQuery.value ? filteredLogs.value : props.logs
  const text = logsToCopy
    .map((l) => `[${formatTime(l.timestamp)}] ${l.type.toUpperCase()} ${l.message}`)
    .join('\n')
  navigator.clipboard.writeText(text)
}

function onContextMenu(e: MouseEvent) {
  const selection = window.getSelection()
  selectedText.value = selection ? selection.toString() : ''
  contextMenuX.value = e.clientX
  contextMenuY.value = e.clientY
  showContextMenu.value = true
}

function copySelectedText() {
  if (selectedText.value) {
    navigator.clipboard.writeText(selectedText.value)
  }
}

function copyAllLogs() {
  copyLogs()
}

function handleCopySelected() {
  copySelectedText()
  showContextMenu.value = false
}

function handleCopyAll() {
  copyAllLogs()
  showContextMenu.value = false
}

function handleClickOutside(e: MouseEvent) {
  if (!showContextMenu.value) return
  const target = e.target as HTMLElement
  const menuEl = document.querySelector('.context-menu')
  if (menuEl && !menuEl.contains(target)) {
    showContextMenu.value = false
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && showContextMenu.value) {
    showContextMenu.value = false
  }
}
</script>

<style scoped>
.log-card {
  display: flex;
  flex-direction: column;
  height: 280px;
  background-color: #1e1e1e;
  transition: all 0.3s ease;
}

.log-card--fullscreen {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 280px;
  height: 100vh;
  z-index: 100;
  border-radius: 0;
}

.log-card--fullscreen .log-console {
  height: calc(100vh - 48px);
}

.search-input {
  max-width: 200px;
}

.log-console {
  flex: 1;
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

.context-menu {
  position: fixed;
  z-index: 1000;
  min-width: 120px;
}
</style>
