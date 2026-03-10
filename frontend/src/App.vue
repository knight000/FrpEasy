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
      @importToml="handleImportToml"
      @exportTomlPreset="handleExportTomlPreset"
      @exportToml="handleExportToml"
      @mergePresets="openMergeDialog"
      @about="openAboutDialog"
    />

    <v-main class="bg-grey-lighten-4">
      <v-tabs v-model="currentTab" align-tabs="start" bg-color="white" color="primary">
        <v-tab value="home">
          <v-icon start size="small">mdi-home</v-icon>
          主页
        </v-tab>
        <v-tab v-for="server in store.activePreset?.servers || []" :key="server.id" :value="server.id">
          <div :class="['status-dot-sm', getStatusDotClass(server.status)]" />
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
                            <div :class="['status-dot', getStatusDotClass(server.status)]" />
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
                          <td>
                            <div class="d-flex align-center">
                              <v-chip
                                v-if="service.is_advanced"
                                color="purple"
                                size="x-small"
                                class="mr-1"
                              >
                                高级
                              </v-chip>
                              <span class="font-weight-medium">{{ service.name }}</span>
                            </div>
                          </td>
                          <td>
                            <v-chip :color="getProtocolColor(service.protocol)" size="x-small">
                              {{ service.protocol }}
                            </v-chip>
                          </td>
                          <td>{{ service.local_ip }}:{{ service.local_port }}</td>
                          <td>{{ displayPort(service) }}</td>
                          <td>
                            <v-icon :color="service.use_encryption ? 'success' : 'grey'" size="small">
                              {{ service.use_encryption ? 'mdi-check-circle' : 'mdi-close-circle' }}
                            </v-icon>
                          </td>
                          <td>
                            <v-icon :color="service.use_compression ? 'success' : 'grey'" size="small">
                              {{ service.use_compression ? 'mdi-check-circle' : 'mdi-close-circle' }}
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
                :fullscreen="logFullscreen"
                @clear="handleClearLogs(currentServer.id)"
                @toggle-fullscreen="logFullscreen = !logFullscreen"
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

    <ConfirmDialog
      v-model="deletePresetDialog"
      title="确认删除预设"
      :message="presetToDelete ? `确定要删除预设 &quot;${presetToDelete.name}&quot; 吗？此操作不可撤销。` : ''"
      confirm-text="删除"
      confirm-color="error"
      @confirm="doDeletePreset"
    />

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

    <v-dialog v-model="servicesDialog" max-width="950">
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
        <v-card-text class="pa-4" style="max-height: 70vh; overflow-y: auto;">
          <template v-if="editingServices.length > 0">
            <v-expansion-panels v-model="expandedServicePanel" multiple>
              <v-expansion-panel
                v-for="(service, i) in editingServices"
                :key="i"
                :value="i"
              >
                <v-expansion-panel-title>
                  <div class="d-flex align-center" style="width: 100%;">
                    <v-chip
                      v-if="service.is_advanced"
                      color="purple"
                      size="x-small"
                      class="mr-2"
                    >
                      高级
                    </v-chip>
                    <span class="font-weight-medium">{{ service.name || '未命名' }}</span>
                    <v-chip :color="getProtocolColor(service.protocol)" size="x-small" class="ml-2">
                      {{ service.protocol }}
                    </v-chip>
                    <span class="text-grey text-caption ml-2">
                      {{ service.local_ip }}:{{ service.local_port }} → {{ displayPort(service) }}
                    </span>
                    <v-spacer />
                    <v-btn
                      icon="mdi-delete"
                      size="small"
                      variant="text"
                      color="error"
                      @click.stop="editingServices.splice(i, 1); expandedServicePanel = null"
                    />
                  </div>
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="d-flex align-center mb-3">
                    <v-btn-toggle
                      v-model="service.is_advanced"
                      :mandatory="true"
                      density="compact"
                      color="primary"
                    >
                      <v-btn :value="false" size="small">
                        <v-icon start>mdi-form-textbox</v-icon>
                        基础
                      </v-btn>
                      <v-btn :value="true" size="small">
                        <v-icon start>mdi-code-tags</v-icon>
                        高级
                      </v-btn>
                    </v-btn-toggle>
                    <v-chip v-if="service.is_advanced" color="warning" size="x-small" class="ml-2">
                      高级模式：直接编辑 TOML
                    </v-chip>
                  </div>

                  <template v-if="!service.is_advanced">
                    <v-row>
                      <v-col cols="6">
                        <v-text-field
                          v-model="service.name"
                          label="名称"
                          density="compact"
                          variant="outlined"
                          hide-details
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-select
                          v-model="service.protocol"
                          :items="['TCP', 'UDP', 'HTTP', 'HTTPS']"
                          label="协议"
                          density="compact"
                          variant="outlined"
                          hide-details
                        />
                      </v-col>
                      <v-col cols="4">
                        <v-text-field
                          v-model="service.local_ip"
                          label="本地IP"
                          density="compact"
                          variant="outlined"
                          hide-details
                        />
                      </v-col>
                      <v-col cols="4">
                        <v-text-field
                          v-model.number="service.local_port"
                          label="本地端口"
                          type="number"
                          density="compact"
                          variant="outlined"
                          hide-details
                        />
                      </v-col>
                      <v-col cols="4">
                        <v-text-field
                          v-model.number="service.remote_port"
                          label="远程端口"
                          type="number"
                          density="compact"
                          variant="outlined"
                          hide-details
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-checkbox
                          v-model="service.use_encryption"
                          label="启用加密"
                          density="compact"
                          hide-details
                          color="primary"
                        />
                      </v-col>
                      <v-col cols="6">
                        <v-checkbox
                          v-model="service.use_compression"
                          label="启用压缩"
                          density="compact"
                          hide-details
                          color="primary"
                        />
                      </v-col>
                    </v-row>
                  </template>

                  <template v-else>
                    <div class="text-caption text-grey mb-2">
                      直接编辑 [[proxies]] 中的配置内容（不包含 [[proxies]] 标题）
                    </div>
                    <v-textarea
                      v-model="service.advanced_config"
                      label="TOML 配置"
                      rows="10"
                      auto-grow
                      density="compact"
                      variant="outlined"
                      :hint="getAdvancedConfigHint(service)"
                      persistent-hint
                      style="font-family: monospace;"
                    />
                    <v-btn
                      size="small"
                      variant="text"
                      color="primary"
                      class="mt-2"
                      @click="generateDefaultAdvancedConfig(service)"
                    >
                      从基础配置生成
                    </v-btn>
                  </template>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
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

    <ConfirmDialog
      v-model="deleteDialog"
      title="确认删除"
      message="确定要删除此服务器吗？此操作不可撤销。"
      confirm-text="删除"
      confirm-color="error"
      @confirm="doDeleteServer"
    />

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
              {{ store.downloadProgress?.total_bytes ? Math.round(store.downloadProgress.total_bytes / 1024 / 1024 * 100) / 100 : 0 }} MB
            </div>
            <div class="text-center text-caption text-grey mt-2">
              下载源: {{ getSourceLabel(downloadSource) }}
            </div>
          </template>
          <template v-else-if="store.versionFetchError">
            <v-alert type="warning" variant="tonal" class="mb-4">
              {{ store.versionFetchError }}
            </v-alert>
            <div class="text-center">
              是否使用默认版本 v0.61.1 下载？
            </div>
            <DownloadSourceSelect v-model="downloadSource" label="切换下载源" class="mt-4" />
          </template>
          <template v-else-if="store.downloadProgress?.is_error">
            <v-alert type="error" variant="tonal" class="mb-4">
              {{ store.downloadProgress?.error_message || '下载失败' }}
            </v-alert>
            <DownloadSourceSelect v-model="downloadSource" label="切换下载源" class="mb-4" />
          </template>
          <template v-else-if="store.downloadProgress?.is_complete">
            <div class="text-center pa-4">
              <v-icon size="64" color="success">mdi-check-circle</v-icon>
              <div class="text-h6 mt-4">下载完成</div>
              <div class="text-grey mt-2">frpc v{{ store.downloadedVersion }} 已安装</div>
            </div>
          </template>
          <template v-else>
            <div class="text-center pa-4">
              <v-icon size="64" color="primary">mdi-download</v-icon>
              <div class="text-h6 mt-4">需要下载 frpc 客户端</div>
              <div class="text-grey mt-2">点击下方按钮开始下载</div>
            </div>
            <DownloadSourceSelect v-model="downloadSource" label="选择下载源" class="mt-4" />
          </template>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <template v-if="store.versionFetchError">
            <v-btn variant="text" @click="downloadDialog = false">取消</v-btn>
            <v-btn color="primary" variant="flat" @click="startDownloadWithDefault">
              下载默认版本
            </v-btn>
          </template>
          <template v-else-if="store.downloadProgress?.is_complete">
            <v-btn color="primary" variant="flat" @click="closeDownloadDialog">
              完成
            </v-btn>
          </template>
          <template v-else>
            <v-btn
              v-if="!store.isDownloading"
              variant="text"
              @click="downloadDialog = false"
            >
              取消
            </v-btn>
            <v-btn
              v-if="store.downloadProgress?.is_error"
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
          </template>
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

    <v-dialog v-model="aboutDialog" max-width="400">
      <v-card>
        <v-card-title class="d-flex align-center">
          <v-icon class="mr-2" color="primary">mdi-information</v-icon>
          关于 FrpEasy
        </v-card-title>
        <v-divider />
        <v-card-text class="pa-4">
          <div class="text-center mb-4">
            <v-icon size="48" color="primary">mdi-swap-horizontal-circle</v-icon>
            <div class="text-h5 mt-2">FrpEasy</div>
            <div class="text-caption text-grey">v{{ appVersion }}</div>
          </div>

          <v-card variant="outlined" class="mb-4">
            <v-card-text class="pa-3">
              <div class="text-body-2 mb-2">frp 客户端配置管理工具</div>
              <div class="d-flex justify-space-between align-center">
                <span class="text-grey">frpc 版本:</span>
                <span v-if="currentFrpcVersion">{{ currentFrpcVersion }}</span>
                <span v-else class="text-grey">未安装</span>
              </div>
              <div class="d-flex justify-space-between align-center mt-1">
                <span class="text-grey">最新版本:</span>
                <span v-if="checkingFrpcUpdate">
                  <v-progress-circular indeterminate size="14" />
                </span>
                <span v-else-if="latestFrpcVersion">{{ latestFrpcVersion }}</span>
                <span v-else class="text-grey">-</span>
              </div>
            </v-card-text>
          </v-card>

          <v-btn
            block
            color="primary"
            variant="tonal"
            class="mb-2"
            :loading="checkingFrpcUpdate"
            @click="checkFrpcUpdate"
          >
            <v-icon start>mdi-update</v-icon>
            检查 frpc 更新
          </v-btn>

          <v-btn
            v-if="hasFrpcUpdate"
            block
            color="success"
            variant="tonal"
            class="mb-2"
            @click="openDownloadFrpc"
          >
            <v-icon start>mdi-download</v-icon>
            下载 frpc 更新
          </v-btn>

          <v-btn
            block
            color="grey"
            variant="tonal"
            @click="checkAppUpdate"
          >
            <v-icon start>mdi-cellphone-arrow-down</v-icon>
            检查 FrpEasy 更新
          </v-btn>
        </v-card-text>
        <v-divider />
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="aboutDialog = false">关闭</v-btn>
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
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import DownloadSourceSelect from '@/components/DownloadSourceSelect.vue'
import TOML from 'smol-toml'
import { getStatusDotClass, getStatusChipColor, getStatusText } from '@/composables/useStatus'
import { useSnackbar } from '@/composables/useSnackbar'
import { useDownloadSource } from '@/composables/useDownloadSource'
import { parseAdvancedConfigToBasic } from '@/helpers/serviceParser'

const store = usePresetStore()
const { snackbar, showSnackbar } = useSnackbar()
const { getSourceLabel } = useDownloadSource()

const currentTab = ref<string>('home')
const createDialog = ref(false)
const newName = ref('')
const copyFromPresetId = ref<string>('')
const serverDialog = ref(false)
const servicesDialog = ref(false)
const deleteDialog = ref(false)
const deletePresetDialog = ref(false)
const downloadDialog = ref(false)
const downloadSource = ref<'github' | 'ghproxy'>('ghproxy')

const serverToDelete = ref<string | null>(null)
const presetToDelete = ref<{id: string, name: string} | null>(null)
const isNewServer = ref(false)
const editingServer = ref<{name: string, address: string, port: number, token: string} | null>(null)
const editingServices = ref<Service[]>([])
const expandedServicePanel = ref<number | number[] | null>(null)

const mergeDialog = ref(false)
const mergedPresetName = ref('')
const selectedMergeIds = ref<string[]>([])

const logFullscreen = ref(false)

const aboutDialog = ref(false)
const appVersion = ref('1.0.0')
const currentFrpcVersion = ref('')
const latestFrpcVersion = ref('')
const checkingFrpcUpdate = ref(false)
const hasFrpcUpdate = ref(false)

onMounted(async () => {
  await store.initPresets()
  await store.initFrpc()
  store.setupLogListener()
  await store.autoStartEnabledServers()
})

onUnmounted(() => {
  store.cleanupListeners()
})

const currentServer = computed(() => {
  if (currentTab.value === 'home') return null
  return store.activePreset?.servers?.find((s) => s.id === currentTab.value) || null
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

function getProtocolColor(protocol: Service['protocol']) {
  const colors: Record<string, string> = {
    TCP: 'primary',
    UDP: 'success',
    HTTP: 'warning',
    HTTPS: 'info',
  }
  return colors[protocol] || 'grey'
}

function displayPort(service: Service): string {
  if (service.display_ports) {
    return service.display_ports
  }
  return String(service.remote_port)
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
    showSnackbar(sourcePreset ? '预设已复制创建' : '预设已创建')
  }
}

function handleCopyPreset(id: string) {
  store.copyPreset(id)
  showSnackbar('预设已复制到剪贴板')
}

function handlePastePreset() {
  if (store.pastePreset('新建预设 (副本)')) {
    showSnackbar('预设已粘贴')
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
    showSnackbar('预设已删除')
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

function startDownloadWithDefault() {
    store.startDownloadFrpc(downloadSource.value, true)
}

async function closeDownloadDialog() {
  downloadDialog.value = false
  if (aboutDialog.value) {
    currentFrpcVersion.value = await store.getCurrentFrpcVersion()
    hasFrpcUpdate.value = false
    latestFrpcVersion.value = ''
  }
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
  showSnackbar('已启动全部服务器')
}

function stopAllServers() {
  if (!store.activePreset) return
  store.activePreset.servers.forEach((s) => {
    if (s.enabled) store.toggleServer(store.activePreset!.id, s.id)
  })
  showSnackbar('已停止全部服务器')
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
    showSnackbar('服务器已添加')
  } else if (currentServer.value) {
    store.updateServer(store.activePreset.id, currentServer.value.id, editingServer.value)
    showSnackbar('配置已保存')
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
    showSnackbar('服务器已删除')
  }
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
    local_port: 8080,
    remote_port: 9000,
    local_ip: '127.0.0.1',
    use_encryption: false,
    use_compression: false,
    advanced_config: '',
    is_advanced: false
  })
  expandedServicePanel.value = editingServices.value.length - 1
}

function generateDefaultAdvancedConfig(service: Service) {
  const lines: string[] = []
  lines.push('[[proxies]]')
  lines.push(`name = "${service.name}"`)
  lines.push(`type = "${service.protocol.toLowerCase()}"`)
  lines.push(`localIP = "${service.local_ip}"`)
  lines.push(`localPort = ${service.local_port}`)
  lines.push(`remotePort = ${service.remote_port}`)
  if (service.use_encryption) {
    lines.push('transport.useEncryption = true')
  }
  if (service.use_compression) {
    lines.push('transport.useCompression = true')
  }
  service.advanced_config = lines.join('\n')
}

function getAdvancedConfigHint(service: Service): string {
  return `name = "${service.name}"\ntype = "${service.protocol.toLowerCase()}"\nlocalIP = "${service.local_ip}"\nlocalPort = ${service.local_port}\nremotePort = ${service.remote_port}`
}

function saveServices() {
  if (!store.activePreset) return
  
  for (const service of editingServices.value) {
    if (service.is_advanced && service.advanced_config) {
      parseAdvancedConfigToBasic(service)
    }
  }
  
  store.updatePreset(store.activePreset.id, { services: JSON.parse(JSON.stringify(editingServices.value)) })
  servicesDialog.value = false
  showSnackbar('服务规则已保存')
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
      showSnackbar(`成功导入 ${successCount} 个预设`)
    }
    if (errorMessages.length > 0) {
      showSnackbar(errorMessages[0], 'error')
    }
  } catch (e) {
    showSnackbar('导入失败', 'error')
  }
}

async function handleImportToml() {
  try {
    const preset = await store.importPresetToml()
    if (preset) {
      store.addImportedPreset(preset)
      showSnackbar('预设已导入')
    }
  } catch (e) {
    showSnackbar('导入失败', 'error')
  }
}

async function handleExportTomlPreset(presetId: string) {
  try {
    const path = await store.exportPresetTomlPreset(presetId)
    if (path) {
      showSnackbar('预设已导出')
    }
  } catch (e) {
    showSnackbar('导出失败', 'error')
  }
}

async function handleExportToml(presetId: string) {
  try {
    const path = await store.exportPresetToml(presetId)
    if (path) {
      showSnackbar('frp 配置已导出到目录')
    }
  } catch (e) {
    showSnackbar('导出失败', 'error')
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

async function doMergePresets() {
  if (selectedMergeIds.value.length < 2 || !mergedPresetName.value.trim()) return
  
  const success = await store.mergePresets(selectedMergeIds.value, mergedPresetName.value.trim())
  
  if (success) {
    mergeDialog.value = false
    showSnackbar(`已合并 ${selectedMergeIds.value.length} 个预设`)
  } else {
    showSnackbar('合并失败', 'error')
  }
}

async function openAboutDialog() {
  aboutDialog.value = true
  appVersion.value = await store.getAppVersion()
  currentFrpcVersion.value = await store.getCurrentFrpcVersion()
  latestFrpcVersion.value = ''
}

async function checkFrpcUpdate() {
  checkingFrpcUpdate.value = true
  try {
    latestFrpcVersion.value = await store.getLatestFrpcVersion()
    const compareResult = await store.compareFrpcVersions(latestFrpcVersion.value, currentFrpcVersion.value)
    hasFrpcUpdate.value = compareResult > 0
    if (hasFrpcUpdate.value) {
      showSnackbar(`发现新版本 frpc ${latestFrpcVersion.value}`, 'info')
    } else if (latestFrpcVersion.value) {
      showSnackbar('frpc 已是最新版本')
    } else {
      showSnackbar('无法获取版本信息', 'warning')
    }
  } catch {
    showSnackbar('检查更新失败', 'error')
  } finally {
    checkingFrpcUpdate.value = false
  }
}

function openDownloadFrpc() {
  aboutDialog.value = false
  downloadDialog.value = true
}

function checkAppUpdate() {
  showSnackbar('此功能暂未开放，请关注 GitHub 发布页', 'info')
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
.services-table td.text-center {
  vertical-align: middle;
}
.services-table td.text-center .v-checkbox .v-selection-control {
  justify-content: center;
}
</style>
