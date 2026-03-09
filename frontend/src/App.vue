<template>
  <v-app>
    <PresetSidebar
      :active-preset-id="store.activePresetId"
      :presets="store.presets"
      @create="openCreatePreset"
      @select="store.setActivePreset"
      @toggleServer="handleToggleServer"
      @copy="handleCopyPreset"
      @paste="handlePastePreset"
      @delete="handleDeletePreset"
      @importFrp="handleImportFrp"
      @importJson="handleImportJson"
      @exportJson="handleExportJson"
      @exportToml="handleExportToml"
      @mergePresets="openMergeDialog"
    />

    <v-main class="bg-grey-lighten-4">
      <v-tabs v-model="currentTab" align-tabs="start" bg-color="white" color="primary">
        <v-tab value="home">
          <v-icon start size="small">mdi-home</v-icon>
          主页
        </v-tab>
        <v-tab v-for="server in store.activePreset?.servers || []" :key="server.id" :value="server.id">
          <div :class="['status-dot-sm', getStatusColor(server.status)]" />
          {{ server.name }}
        </v-tab>
        <v-btn icon="mdi-plus" variant="text" class="ml-2" @click="openAddServer" />
      </v-tabs>

      <v-container class="pa-4 main-container" fluid>
        <template v-if="store.activePreset">
          <!-- 主页视图 -->
          <template v-if="currentTab === 'home'">
            <v-row>
              <v-col cols="12" md="8">
                <v-card class="mb-4" variant="outlined">
                  <v-card-title class="d-flex align-center py-3">
                    <v-icon class="mr-2">mdi-server-network</v-icon>
                    服务器列表
                    <v-chip size="small" class="ml-2">{{ store.activePreset?.servers?.length || 0 }}</v-chip>
                    <v-spacer />
                    <v-btn size="small" color="primary" variant="tonal" @click="openAddServer">
                      <v-icon start>mdi-plus</v-icon>
                      添加服务器
                    </v-btn>
                  </v-card-title>
                  <v-divider />
                  <v-card-text class="pa-0">
                    <v-table density="compact">
                      <thead>
                        <tr>
                          <th style="width: 40px"></th>
                          <th>名称</th>
                          <th>地址</th>
                          <th>状态</th>
                          <th style="width: 120px">操作</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr
                          v-for="server in store.activePreset?.servers || []"
                          :key="server.id"
                          class="server-row"
                          @click="currentTab = server.id"
                        >
                          <td>
                            <div :class="['status-dot', getStatusColor(server.status)]" />
                          </td>
                          <td>
                            <div class="font-weight-medium">{{ server.name }}</div>
                          </td>
                          <td>{{ server.address }}:{{ server.port }}</td>
                          <td>
                            <v-chip :color="getStatusChipColor(server.status)" size="x-small">
                              {{ getStatusText(server.status) }}
                            </v-chip>
                          </td>
                          <td>
                            <v-switch
                              :model-value="server.enabled"
                              color="success"
                              density="compact"
                              hide-details
                              inline
                              @click.stop
                              @update:model-value="toggleServer(server.id)"
                            />
                          </td>
                        </tr>
                      </tbody>
                    </v-table>
                  </v-card-text>
                </v-card>

                <v-card variant="outlined">
                  <v-card-title class="d-flex align-center py-3">
                    <v-icon class="mr-2">mdi-swap-horizontal</v-icon>
                    服务规则
                    <v-chip size="small" class="ml-2">{{ store.activePreset?.services?.length || 0 }}</v-chip>
                    <v-spacer />
                    <v-btn size="small" color="primary" variant="tonal" @click="openEditServices">
                      <v-icon start>mdi-pencil</v-icon>
                      编辑规则
                    </v-btn>
                  </v-card-title>
                  <v-divider />
                  <v-card-text class="pa-0">
                    <v-table v-if="(store.activePreset?.services?.length || 0) > 0" density="compact">
                      <thead>
                        <tr>
                          <th>名称</th>
                          <th>协议</th>
                          <th>本地地址</th>
                          <th>远程端口</th>
                          <th>加密</th>
                          <th>压缩</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="service in store.activePreset?.services || []" :key="service.id">
                          <td class="font-weight-medium">{{ service.name }}</td>
                          <td>
                            <v-chip :color="getProtocolColor(service.protocol)" size="x-small">
                              {{ service.protocol }}
                            </v-chip>
                          </td>
                          <td>{{ service.localIp }}:{{ service.localPort }}</td>
                          <td>{{ service.remotePort }}</td>
                          <td>
                            <v-icon :color="service.useEncryption ? 'success' : 'grey'" size="small">
                              {{ service.useEncryption ? 'mdi-check-circle' : 'mdi-close-circle' }}
                            </v-icon>
                          </td>
                          <td>
                            <v-icon :color="service.useCompression ? 'success' : 'grey'" size="small">
                              {{ service.useCompression ? 'mdi-check-circle' : 'mdi-close-circle' }}
                            </v-icon>
                          </td>
                        </tr>
                      </tbody>
                    </v-table>
                    <div v-else class="pa-4 text-center text-grey">
                      <v-icon size="32" class="mb-2">mdi-information-outline</v-icon>
                      <div>暂未配置服务规则</div>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" md="4">
                <v-card variant="outlined" class="mb-4">
                  <v-card-title class="py-3">
                    <v-icon class="mr-2">mdi-chart-box</v-icon>
                    状态概览
                  </v-card-title>
                  <v-divider />
                  <v-card-text>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">服务器总数</span>
                      <span class="font-weight-bold">{{ store.activePreset?.servers?.length || 0 }}</span>
                    </div>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">在线服务器</span>
                      <span class="font-weight-bold text-success">{{ onlineCount }}</span>
                    </div>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">服务规则</span>
                      <span class="font-weight-bold">{{ store.activePreset?.services?.length || 0 }}</span>
                    </div>
                  </v-card-text>
                </v-card>

                <v-card variant="outlined">
                  <v-card-title class="py-3">
                    <v-icon class="mr-2">mdi-cog</v-icon>
                    快捷操作
                  </v-card-title>
                  <v-divider />
                  <v-card-text>
                    <v-btn block color="primary" variant="tonal" class="mb-2" @click="startAllServers">
                      <v-icon start>mdi-play</v-icon>
                      启动全部
                    </v-btn>
                    <v-btn block color="error" variant="tonal" class="mb-2" @click="stopAllServers">
                      <v-icon start>mdi-stop</v-icon>
                      停止全部
                    </v-btn>
                    <v-btn
                      v-if="store.presets.length > 1"
                      block
                      color="error"
                      variant="text"
                      @click="handleDeletePreset(store.activePresetId!)"
                    >
                      <v-icon start>mdi-delete</v-icon>
                      删除预设
                    </v-btn>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- 单服务器时显示控制台 -->
            <div v-if="singleServer" class="mt-4">
              <LogConsole
                :logs="singleServer.logs"
                :title="`${singleServer.name} - 控制台输出`"
                @clear="handleClearLogs(singleServer.id)"
              />
            </div>
          </template>

          <!-- 服务器详情视图 -->
          <template v-else-if="currentServer">
            <v-row>
              <v-col cols="12" md="6">
                <v-card variant="outlined">
                  <v-card-title class="d-flex align-center py-3">
                    <v-icon class="mr-2">mdi-server</v-icon>
                    服务器配置
                    <v-spacer />
                    <v-btn size="small" variant="text" @click="openEditServer(currentServer)">
                      <v-icon>mdi-pencil</v-icon>
                    </v-btn>
                  </v-card-title>
                  <v-divider />
                  <v-card-text>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">名称</span>
                      <span class="font-weight-medium">{{ currentServer.name }}</span>
                    </div>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">地址</span>
                      <span>{{ currentServer.address }}:{{ currentServer.port }}</span>
                    </div>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">状态</span>
                      <v-chip :color="getStatusChipColor(currentServer.status)" size="small">
                        {{ getStatusText(currentServer.status) }}
                      </v-chip>
                    </div>
                  </v-card-text>
                  <v-divider />
                  <v-card-actions class="pa-4">
                    <v-switch
                      :model-value="currentServer.enabled"
                      color="success"
                      label="启用连接"
                      @update:model-value="toggleServer(currentServer.id)"
                    />
                    <v-spacer />
                    <v-btn
                      v-if="(store.activePreset?.servers?.length || 0) > 1"
                      color="error"
                      variant="text"
                      @click="confirmDeleteServer(currentServer.id)"
                    >
                      <v-icon start>mdi-delete</v-icon>
                      删除
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-col>

              <v-col cols="12" md="6">
                <v-card variant="outlined">
                  <v-card-title class="py-3">
                    <v-icon class="mr-2">mdi-information</v-icon>
                    运行信息
                  </v-card-title>
                  <v-divider />
                  <v-card-text>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">运行时间</span>
                      <span class="font-weight-medium">{{ formatUptime(currentServer.uptime) }}</span>
                    </div>
                    <div class="d-flex justify-space-between mb-3">
                      <span class="text-grey">日志条数</span>
                      <span>{{ currentServer.logs.length }}</span>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <div class="mt-4">
              <LogConsole
                :logs="currentServer.logs"
                :title="`${currentServer.name} - 控制台输出`"
                @clear="handleClearLogs(currentServer.id)"
              />
            </div>
          </template>
        </template>

        <v-card v-else class="pa-6 text-center">
          <v-icon color="grey" size="64">mdi-folder-plus</v-icon>
          <div class="text-h6 mt-4">暂无预设</div>
          <v-btn class="mt-4" color="primary" @click="openCreatePreset">
            <v-icon start>mdi-plus</v-icon>
            创建预设
          </v-btn>
        </v-card>
      </v-container>
    </v-main>

    <v-dialog v-model="createDialog" max-width="450">
      <v-card>
        <v-card-title>创建新预设</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newName"
            label="预设名称"
            variant="outlined"
            class="mb-3"
            @keyup.enter="createPreset"
          />
          <v-select
            v-model="copyFromPresetId"
            :items="[{ title: '空白预设', value: '' }, ...store.presets.map(p => ({ title: p.name, value: p.id }))]"
            label="复制自"
            variant="outlined"
            clearable
            hint="选择已有预设进行复制，或创建空白预设"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="createDialog = false">取消</v-btn>
          <v-btn color="primary" variant="flat" @click="createPreset">创建</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="deletePresetDialog" max-width="400">
      <v-card>
        <v-card-title>确认删除预设</v-card-title>
        <v-card-text>
          确定要删除预设 "{{ presetToDelete?.name }}" 吗？此操作不可撤销。
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deletePresetDialog = false">取消</v-btn>
          <v-btn color="error" variant="flat" @click="doDeletePreset">删除</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="serverDialog" max-width="600">
      <v-card v-if="editingServer">
        <v-card-title>{{ isNewServer ? '添加服务器' : '编辑服务器' }}</v-card-title>
        <v-divider />
        <v-card-text class="pa-4">
          <v-text-field v-model="editingServer.name" label="服务器名称" variant="outlined" class="mb-3" />
          <v-text-field v-model="editingServer.address" label="服务器地址" variant="outlined" class="mb-3" />
          <v-row>
            <v-col cols="6">
              <v-text-field v-model.number="editingServer.port" label="端口" type="number" variant="outlined" />
            </v-col>
            <v-col cols="6">
              <v-text-field v-model="editingServer.token" label="令牌" type="password" variant="outlined" />
            </v-col>
          </v-row>
        </v-card-text>
        <v-divider />
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="text" @click="serverDialog = false">取消</v-btn>
          <v-btn color="primary" variant="flat" @click="saveServer">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="servicesDialog" max-width="1000">
      <v-card>
        <v-card-title class="d-flex align-center">
          编辑服务规则
          <v-spacer />
          <v-btn size="small" color="primary" variant="tonal" @click="addService">
            <v-icon start>mdi-plus</v-icon>
            添加规则
          </v-btn>
        </v-card-title>
        <v-divider />
        <v-card-text class="pa-4">
          <template v-if="editingServices.length > 0">
            <v-simple-table dense class="services-table">
              <thead>
                <tr>
                  <th style="width: 180px">名称</th>
                  <th style="width: 140px">协议</th>
                  <th style="width: 160px">本地IP</th>
                  <th style="width: 100px">本地端口</th>
                  <th style="width: 100px">远程端口</th>
                  <th style="width: 60px" class="text-center">加密</th>
                  <th style="width: 60px" class="text-center">压缩</th>
                  <th style="width: 40px"></th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(service, i) in editingServices" :key="i">
                  <td><v-text-field v-model="service.name" density="compact" hide-details variant="outlined" /></td>
                  <td><v-select v-model="service.protocol" :items="['TCP', 'UDP', 'HTTP', 'HTTPS']" density="compact" hide-details variant="outlined" /></td>
                  <td><v-text-field v-model="service.localIp" density="compact" hide-details variant="outlined" /></td>
                  <td><v-text-field v-model.number="service.localPort" density="compact" hide-details type="number" variant="outlined" /></td>
                  <td><v-text-field v-model.number="service.remotePort" density="compact" hide-details type="number" variant="outlined" /></td>
                  <td class="text-center"><v-checkbox v-model="service.useEncryption" density="compact" hide-details color="primary" /></td>
                  <td class="text-center"><v-checkbox v-model="service.useCompression" density="compact" hide-details color="primary" /></td>
                  <td class="text-center"><v-btn icon="mdi-delete" size="small" variant="text" color="error" @click="editingServices.splice(i, 1)" /></td>
                </tr>
              </tbody>
            </v-simple-table>
          </template>
          <div v-else class="text-center py-8 text-grey">暂未配置服务规则</div>
        </v-card-text>
        <v-divider />
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="text" @click="servicesDialog = false">取消</v-btn>
          <v-btn color="primary" variant="flat" @click="saveServices">保存</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title>确认删除</v-card-title>
        <v-card-text>确定要删除此服务器吗？此操作不可撤销。</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialog = false">取消</v-btn>
          <v-btn color="error" variant="flat" @click="doDeleteServer">删除</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="downloadDialog" max-width="500" persistent>
      <v-card>
        <v-card-title>下载 frpc</v-card-title>
        <v-card-text>
          <template v-if="store.isDownloading">
            <div class="text-center mb-4">
              <v-progress-circular
                :model-value="store.downloadProgress?.percentage || 0"
                size="100"
                width="8"
                color="primary"
              >
                {{ Math.round(store.downloadProgress?.percentage || 0) }}%
              </v-progress-circular>
            </div>
            <div class="text-center text-grey">
              正在下载 frpc {{ store.downloadProgress?.downloaded ? Math.round(store.downloadProgress.downloaded / 1024 / 1024 * 100) / 100 : 0 }} MB / 
              {{ store.downloadProgress?.totalBytes ? Math.round(store.downloadProgress.totalBytes / 1024 / 1024 * 100) / 100 : 0 }} MB
            </div>
            <div class="text-center text-caption text-grey mt-2">
              下载源: {{ getSourceLabel(downloadSource) }}
            </div>
          </template>
          <template v-else-if="store.downloadProgress?.isError">
            <v-alert type="error" variant="tonal" class="mb-4">
              {{ store.downloadProgress?.errorMessage || '下载失败' }}
            </v-alert>
            <div class="mb-4">
              <div class="text-subtitle-2 mb-2">切换下载源:</div>
              <v-select
                v-model="downloadSource"
                :items="downloadSources"
                item-title="label"
                item-value="value"
                density="compact"
                variant="outlined"
                hide-details
              />
            </div>
          </template>
          <template v-else>
            <div class="text-center pa-4">
              <v-icon size="64" color="primary">mdi-download</v-icon>
              <div class="text-h6 mt-4">需要下载 frpc 客户端</div>
              <div class="text-grey mt-2">点击下方按钮开始下载</div>
            </div>
            <div class="mt-4">
              <div class="text-subtitle-2 mb-2">选择下载源:</div>
              <v-select
                v-model="downloadSource"
                :items="downloadSources"
                item-title="label"
                item-value="value"
                density="compact"
                variant="outlined"
                hide-details
              />
            </div>
          </template>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            v-if="!store.isDownloading"
            variant="text"
            @click="downloadDialog = false"
          >
            取消
          </v-btn>
          <v-btn
            v-if="store.downloadProgress?.isError"
            color="primary"
            variant="flat"
            @click="startDownload"
          >
            重试
          </v-btn>
          <v-btn
            v-else-if="!store.isDownloading"
            color="primary"
            variant="flat"
            @click="startDownload"
          >
            开始下载
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="mergeDialog" max-width="500">
      <v-card>
        <v-card-title>合并预设</v-card-title>
        <v-divider />
        <v-card-text class="pa-4">
          <v-text-field
            v-model="mergedPresetName"
            label="新预设名称"
            variant="outlined"
            class="mb-4"
            placeholder="合并后的预设名称"
          />
          <div class="text-subtitle-2 mb-2">选择要合并的预设:</div>
          <v-list v-if="store.presets.length >= 2" class="border rounded" density="compact">
            <v-list-item
              v-for="preset in store.presets"
              :key="preset.id"
              @click="toggleMergeSelection(preset.id)"
            >
              <template #prepend>
                <v-checkbox
                  :model-value="selectedMergeIds.includes(preset.id)"
                  hide-details
                  density="compact"
                  @click.stop
                  @update:model-value="toggleMergeSelection(preset.id)"
                />
              </template>
              <v-list-item-title>{{ preset.name }}</v-list-item-title>
              <v-list-item-subtitle>
                {{ preset.servers?.length || 0 }} 服务器, {{ preset.services?.length || 0 }} 服务
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
          <div v-else class="text-center py-4 text-grey">
            需要至少 2 个预设才能合并
          </div>
          <v-alert v-if="selectedMergeIds.length >= 2" type="info" variant="tonal" class="mt-4" density="compact">
            将合并 {{ selectedMergeIds.length }} 个预设
            ({{ getMergedServerCount }} 服务器, {{ getMergedServiceCount }} 服务)
          </v-alert>
        </v-card-text>
        <v-divider />
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn variant="text" @click="mergeDialog = false">取消</v-btn>
          <v-btn
            color="primary"
            variant="flat"
            :disabled="selectedMergeIds.length < 2 || !mergedPresetName.trim()"
            @click="doMergePresets"
          >
            合并
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" :timeout="3000">
      {{ snackbar.message }}
    </v-snackbar>
  </v-app>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { usePresetStore, type Server, type Service } from '@/stores/preset'
import PresetSidebar from '@/components/PresetSidebar.vue'
import LogConsole from '@/components/LogConsole.vue'

const store = usePresetStore()

const currentTab = ref<string>('home')
const createDialog = ref(false)
const newName = ref('')
const copyFromPresetId = ref<string>('')
const serverDialog = ref(false)
const servicesDialog = ref(false)
const deleteDialog = ref(false)
const deletePresetDialog = ref(false)
const downloadDialog = ref(false)
const downloadSource = ref<'github' | 'ghproxy' | 'fastgit' | 'moeyy'>('ghproxy')
const downloadSources = [
  { value: 'ghproxy', label: 'GHProxy (国内加速)' },
  { value: 'moeyy', label: 'Moeyy (国内加速)' },
  { value: 'fastgit', label: 'FastGit (国内加速)' },
  { value: 'github', label: 'GitHub (直连)' },
]

function getSourceLabel(source: string) {
  const item = downloadSources.find(s => s.value === source)
  return item?.label || source
}
const serverToDelete = ref<string | null>(null)
const presetToDelete = ref<{id: string, name: string} | null>(null)
const isNewServer = ref(false)
const editingServer = ref<{name: string, address: string, port: number, token: string} | null>(null)
const editingServices = ref<Service[]>([])
const snackbar = ref({ show: false, message: '', color: 'success' })

const mergeDialog = ref(false)
const mergedPresetName = ref('')
const selectedMergeIds = ref<string[]>([])

onMounted(async () => {
  await store.initFrpc()
  store.setupLogListener()
})

onUnmounted(() => {
  store.cleanupListeners()
})

const currentServer = computed(() => {
  if (currentTab.value === 'home') return null
  return store.activePreset?.servers?.find((s) => s.id === currentTab.value) || null
})

const singleServer = computed(() => {
  const servers = store.activePreset?.servers
  return servers?.length === 1 ? servers[0] : null
})

const onlineCount = computed(() => {
  return store.activePreset?.servers?.filter((s) => s.status === 'online').length || 0
})

const getMergedServerCount = computed(() => {
  return selectedMergeIds.value.reduce((sum, id) => {
    const preset = store.presets.find(p => p.id === id)
    return sum + (preset?.servers?.length || 0)
  }, 0)
})

const getMergedServiceCount = computed(() => {
  return selectedMergeIds.value.reduce((sum, id) => {
    const preset = store.presets.find(p => p.id === id)
    return sum + (preset?.services?.length || 0)
  }, 0)
})

watch(
  () => store.activePresetId,
  () => {
    currentTab.value = 'home'
  }
)

function getStatusColor(status: Server['status']) {
  const map: Record<Server['status'], string> = {
    online: 'bg-success',
    connecting: 'bg-warning',
    error: 'bg-error',
    offline: 'bg-grey'
  }
  return map[status] || 'bg-grey'
}

function getStatusChipColor(status: Server['status']) {
  const map: Record<Server['status'], string> = {
    online: 'success',
    connecting: 'warning',
    error: 'error',
    offline: 'grey'
  }
  return map[status] || 'grey'
}

function getStatusText(status: Server['status']) {
  const map: Record<Server['status'], string> = {
    online: '在线',
    connecting: '连接中',
    error: '错误',
    offline: '离线'
  }
  return map[status] || '离线'
}

function getProtocolColor(protocol: Service['protocol']) {
  return { TCP: 'primary', UDP: 'success', HTTP: 'warning', HTTPS: 'info' }[protocol] || 'grey'
}

function formatUptime(seconds: number): string {
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分${seconds % 60}秒`
  const hours = Math.floor(seconds / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  return `${hours}时${mins}分`
}

function openCreatePreset() {
  newName.value = ''
  copyFromPresetId.value = ''
  createDialog.value = true
}

function createPreset() {
  if (newName.value.trim()) {
    const sourcePreset = copyFromPresetId.value
      ? store.presets.find((p) => p.id === copyFromPresetId.value)
      : undefined
    store.addPreset(newName.value.trim(), sourcePreset)
    createDialog.value = false
    snackbar.value = {
      show: true,
      message: sourcePreset ? '预设已复制创建' : '预设已创建',
      color: 'success'
    }
  }
}

function handleCopyPreset(id: string) {
  store.copyPreset(id)
  snackbar.value = { show: true, message: '预设已复制到剪贴板', color: 'success' }
}

function handlePastePreset() {
  if (store.pastePreset('新建预设 (副本)')) {
    snackbar.value = { show: true, message: '预设已粘贴', color: 'success' }
  }
}

function handleDeletePreset(id: string) {
  const preset = store.presets.find((p) => p.id === id)
  if (preset) {
    presetToDelete.value = { id: preset.id, name: preset.name }
    deletePresetDialog.value = true
  }
}

function doDeletePreset() {
  if (presetToDelete.value) {
    store.deletePreset(presetToDelete.value.id)
    deletePresetDialog.value = false
    snackbar.value = { show: true, message: '预设已删除', color: 'success' }
    presetToDelete.value = null
  }
}

function handleToggleServer(presetId: string, serverId: string) {
  if (!store.frpcDownloaded) {
    downloadDialog.value = true
    return
  }
  store.toggleServer(presetId, serverId)
}

function toggleServer(serverId: string) {
  if (!store.frpcDownloaded) {
    downloadDialog.value = true
    return
  }
  if (store.activePreset) store.toggleServer(store.activePreset.id, serverId)
}

function startDownload() {
  store.startDownloadFrpc(downloadSource.value)
}

function handleClearLogs(serverId: string) {
  if (store.activePreset) {
    store.clearLogs(store.activePreset.id, serverId)
  }
}

function startAllServers() {
  if (!store.activePreset) return
  store.activePreset.servers.forEach((s) => {
    if (!s.enabled) store.toggleServer(store.activePreset!.id, s.id)
  })
  snackbar.value = { show: true, message: '已启动全部服务器', color: 'success' }
}

function stopAllServers() {
  if (!store.activePreset) return
  store.activePreset.servers.forEach((s) => {
    if (s.enabled) store.toggleServer(store.activePreset!.id, s.id)
  })
  snackbar.value = { show: true, message: '已停止全部服务器', color: 'success' }
}

function openAddServer() {
  if (!store.activePreset) return
  isNewServer.value = true
  editingServer.value = {
    name: `服务器 ${store.activePreset.servers.length + 1}`,
    address: '',
    port: 7000,
    token: '',
  }
  serverDialog.value = true
}

function openEditServer(server: Server) {
  isNewServer.value = false
  editingServer.value = {
    name: server.name,
    address: server.address,
    port: server.port,
    token: server.token,
  }
  serverDialog.value = true
}

function saveServer() {
  if (!store.activePreset || !editingServer.value) return
  if (isNewServer.value) {
    store.addServer(store.activePreset.id, editingServer.value)
    snackbar.value = { show: true, message: '服务器已添加', color: 'success' }
  } else if (currentServer.value) {
    store.updateServer(store.activePreset.id, currentServer.value.id, editingServer.value)
    snackbar.value = { show: true, message: '配置已保存', color: 'success' }
  }
  serverDialog.value = false
}

function confirmDeleteServer(serverId: string) {
  serverToDelete.value = serverId
  deleteDialog.value = true
}

function doDeleteServer() {
  if (store.activePreset && serverToDelete.value) {
    store.deleteServer(store.activePreset.id, serverToDelete.value)
    currentTab.value = 'home'
    snackbar.value = { show: true, message: '服务器已删除', color: 'success' }
  }
  deleteDialog.value = false
  serverToDelete.value = null
}

function openEditServices() {
  if (!store.activePreset) return
  editingServices.value = JSON.parse(JSON.stringify(store.activePreset.services))
  servicesDialog.value = true
}

function addService() {
  editingServices.value.push({
    id: Math.random().toString(36).substring(2, 9),
    name: `服务 ${editingServices.value.length + 1}`,
    protocol: 'TCP',
    localPort: 8080,
    remotePort: 9000,
    localIp: '127.0.0.1',
    useEncryption: false,
    useCompression: false
  })
}

function saveServices() {
  if (!store.activePreset) return
  store.updatePreset(store.activePreset.id, { services: JSON.parse(JSON.stringify(editingServices.value)) })
  servicesDialog.value = false
  snackbar.value = { show: true, message: '服务规则已保存', color: 'success' }
}

async function handleImportFrp() {
  try {
    const results = await store.importFrpFiles()
    let successCount = 0
    let errorMessages: string[] = []

    for (const result of results) {
      if (result.preset) {
        store.addImportedPreset(result.preset)
        successCount++
      } else if (result.error) {
        errorMessages.push(result.error)
      }
    }

    if (successCount > 0) {
      snackbar.value = { show: true, message: `成功导入 ${successCount} 个预设`, color: 'success' }
    }
    if (errorMessages.length > 0) {
      snackbar.value = { show: true, message: errorMessages[0], color: 'error' }
    }
  } catch (e) {
    snackbar.value = { show: true, message: '导入失败', color: 'error' }
  }
}

async function handleImportJson() {
  try {
    const preset = await store.importPresetJson()
    if (preset) {
      store.addImportedPreset(preset)
      snackbar.value = { show: true, message: '预设已导入', color: 'success' }
    }
  } catch (e) {
    snackbar.value = { show: true, message: '导入失败', color: 'error' }
  }
}

async function handleExportJson(presetId: string) {
  try {
    const path = await store.exportPresetJson(presetId)
    if (path) {
      snackbar.value = { show: true, message: '预设已导出', color: 'success' }
    }
  } catch (e) {
    snackbar.value = { show: true, message: '导出失败', color: 'error' }
  }
}

async function handleExportToml(presetId: string) {
  const preset = store.presets.find(p => p.id === presetId)
  if (!preset || !preset.servers[0]) return

  try {
    const path = await store.exportPresetToml(presetId, preset.servers[0].id)
    if (path) {
      snackbar.value = { show: true, message: 'frp 配置已导出', color: 'success' }
    }
  } catch (e) {
    snackbar.value = { show: true, message: '导出失败', color: 'error' }
  }
}

function openMergeDialog() {
  mergedPresetName.value = '合并预设'
  selectedMergeIds.value = []
  mergeDialog.value = true
}

function toggleMergeSelection(presetId: string) {
  const index = selectedMergeIds.value.indexOf(presetId)
  if (index === -1) {
    selectedMergeIds.value.push(presetId)
  } else {
    selectedMergeIds.value.splice(index, 1)
  }
}

function doMergePresets() {
  if (selectedMergeIds.value.length < 2 || !mergedPresetName.value.trim()) return
  
  const success = store.mergePresets(selectedMergeIds.value, mergedPresetName.value.trim())
  
  if (success) {
    mergeDialog.value = false
    snackbar.value = { 
      show: true, 
      message: `已合并 ${selectedMergeIds.value.length} 个预设`, 
      color: 'success' 
    }
  } else {
    snackbar.value = { show: true, message: '合并失败', color: 'error' }
  }
}
</script>

<style>
body {
  font-family: 'Roboto', sans-serif;
  margin: 0;
}
html, body, .v-application {
  overflow: hidden !important;
}
.main-container {
  height: calc(100vh - 48px);
  overflow-y: auto;
}
.status-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
}
.status-dot-sm {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
  margin-right: 6px;
}
.server-row {
  cursor: pointer;
}
.server-row:hover {
  background-color: rgba(0, 0, 0, 0.04);
}
.services-table th {
  font-size: 13px;
  font-weight: 500;
  color: #666 !important;
}
.services-table td {
  padding: 4px 8px;
}
.services-table .v-field__input {
  font-size: 13px;
}
.services-table .v-select__selection {
  font-size: 13px;
}
</style>
