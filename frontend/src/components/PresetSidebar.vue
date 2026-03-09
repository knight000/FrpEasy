<template>
  <v-navigation-drawer permanent :width="280" elevation="1">
    <v-toolbar color="primary" density="compact" class="text-white">
      <v-toolbar-title class="text-subtitle-1 font-weight-bold">
        FrpEasy 预设
      </v-toolbar-title>
      <v-spacer />
      <v-tooltip text="新建预设" location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="props"
            icon="mdi-plus"
            variant="text"
            @click="emit('create')"
          />
        </template>
      </v-tooltip>
    </v-toolbar>

    <v-list class="pa-2" nav>
      <v-list-group
        v-for="preset in presets"
        :key="preset.id"
        :value="preset.id"
      >
        <template #activator="{ props }">
          <v-list-item
            v-bind="props"
            :active="activePresetId === preset.id"
            class="mb-1 rounded-lg"
            @click="emit('select', preset.id)"
            @contextmenu.prevent="showContextMenu($event, preset)"
          >
            <template #prepend>
              <v-icon :color="getOnlineCount(preset) === (preset.servers?.length || 0) ? 'success' : getOnlineCount(preset) > 0 ? 'warning' : 'grey'">
                mdi-folder
              </v-icon>
            </template>
            <v-list-item-title class="font-weight-medium">
              {{ preset.name }}
            </v-list-item-title>
            <v-list-item-subtitle>
              <v-chip :color="getOnlineCount(preset) > 0 ? 'success' : 'grey'" size="x-small">
                {{ getOnlineCount(preset) }}/{{ preset.servers?.length || 0 }} 在线
              </v-chip>
              <span class="ml-1">{{ preset.services?.length || 0 }} 规则</span>
            </v-list-item-subtitle>
          </v-list-item>
        </template>

        <v-list-item
          v-for="server in preset.servers"
          :key="server.id"
          class="mb-1 rounded-lg"
          density="compact"
          lines="one"
        >
          <template #prepend>
            <div :class="['status-dot', getStatusColor(server.status)]" />
          </template>
          <v-list-item-title class="text-body-2">
            {{ server.name }}
          </v-list-item-title>
          <v-list-item-subtitle class="text-caption">
            {{ server.address }}:{{ server.port }}
          </v-list-item-subtitle>
          <template #append>
            <v-switch
              :model-value="server.enabled"
              color="success"
              density="compact"
              hide-details
              @click.stop
              @update:model-value="emit('toggleServer', preset.id, server.id)"
            />
          </template>
        </v-list-item>
      </v-list-group>
    </v-list>

    <template #append>
      <div class="pa-3">
        <v-btn block color="primary" variant="tonal" class="mb-2" @click="emit('create')">
          <v-icon start>mdi-plus</v-icon>
          新建预设
        </v-btn>
        <v-btn block color="purple" variant="tonal" class="mb-2" @click="emit('mergePresets')">
          <v-icon start>mdi-merge</v-icon>
          合并预设
        </v-btn>
        <v-btn block color="success" variant="tonal" class="mb-2" @click="emit('importFrp')">
          <v-icon start>mdi-file-import</v-icon>
          导入 frp 配置
        </v-btn>
        <v-btn block color="info" variant="tonal" @click="emit('importToml')">
          <v-icon start>mdi-file-document-outline</v-icon>
          导入 FrpEasy 预设
        </v-btn>
      </div>
    </template>

    <v-menu
      v-model="contextMenu.show"
      :position-x="contextMenu.x"
      :position-y="contextMenu.y"
      absolute
      offset-y
    >
      <v-list density="compact">
        <v-list-item @click="doCopy">
          <template #prepend>
            <v-icon>mdi-content-copy</v-icon>
          </template>
          <v-list-item-title>复制</v-list-item-title>
        </v-list-item>
        <v-list-item :disabled="!hasClipboardData" @click="doPaste">
          <template #prepend>
            <v-icon>mdi-content-paste</v-icon>
          </template>
          <v-list-item-title>粘贴</v-list-item-title>
        </v-list-item>
        <v-divider />
        <v-list-item @click="doExportTomlPreset">
          <template #prepend>
            <v-icon>mdi-export</v-icon>
          </template>
          <v-list-item-title>导出 FrpEasy 预设</v-list-item-title>
        </v-list-item>
        <v-list-item @click="doExportToml">
          <template #prepend>
            <v-icon>mdi-file-document-outline</v-icon>
          </template>
          <v-list-item-title>导出 frp TOML</v-list-item-title>
        </v-list-item>
        <v-divider />
        <v-list-item @click="doDelete">
          <template #prepend>
            <v-icon color="error">mdi-delete</v-icon>
          </template>
          <v-list-item-title class="text-error">删除</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { Preset, Server } from '@/stores/preset'
import { usePresetStore } from '@/stores/preset'

defineProps<{
  presets: Preset[]
  activePresetId: string | null
}>()

const emit = defineEmits<{
  (e: 'select', id: string): void
  (e: 'create'): void
  (e: 'toggleServer', presetId: string, serverId: string): void
  (e: 'delete', id: string): void
  (e: 'copy', id: string): void
  (e: 'paste'): void
  (e: 'importFrp'): void
  (e: 'importToml'): void
  (e: 'exportTomlPreset', id: string): void
  (e: 'exportToml', id: string): void
  (e: 'mergePresets'): void
}>()

const store = usePresetStore()
const hasClipboardData = ref(false)

const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  preset: null as Preset | null,
})

function checkClipboard() {
  hasClipboardData.value = store.hasClipboard()
}

function showContextMenu(e: MouseEvent, preset: Preset) {
  contextMenu.value = {
    show: false,
    x: e.clientX,
    y: e.clientY,
    preset,
  }
  checkClipboard()
  setTimeout(() => {
    contextMenu.value.show = true
  }, 0)
}

function doCopy() {
  if (contextMenu.value.preset) {
    emit('copy', contextMenu.value.preset.id)
    checkClipboard()
  }
}

function doPaste() {
  if (hasClipboardData.value) {
    emit('paste')
  }
}

function doDelete() {
  if (contextMenu.value.preset) {
    emit('delete', contextMenu.value.preset.id)
  }
}

function doExportTomlPreset() {
  if (contextMenu.value.preset) {
    emit('exportTomlPreset', contextMenu.value.preset.id)
  }
}

function doExportToml() {
  if (contextMenu.value.preset) {
    emit('exportToml', contextMenu.value.preset.id)
  }
}

function getStatusColor(status: Server['status']) {
  switch (status) {
    case 'online': return 'bg-success'
    case 'connecting': return 'bg-warning'
    case 'error': return 'bg-error'
    default: return 'bg-grey'
  }
}

function getOnlineCount(preset: Preset): number {
  return preset.servers.filter(s => s.status === 'online').length
}

onMounted(() => {
  checkClipboard()
  window.addEventListener('storage', checkClipboard)
})

onUnmounted(() => {
  window.removeEventListener('storage', checkClipboard)
})
</script>

<style scoped>
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
  margin-right: 8px;
}

.v-list-group .v-list-item {
  padding-inline-start: 20px !important;
}
</style>
