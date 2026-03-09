#!/ 创建修复脚本
cat > fix_app.sed << 'EOF'
        <template v-if="store.activePreset">
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
                              @update:model-value="handleToggleServer(store.activePreset!.id, server.id)"
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

            <div v-if="singleServer" class="mt-4">
              <LogConsole
                :logs="singleServer.logs"
                :title="`${singleServer.name} - 控制台输出`"
                @clear="handleClearLogs(singleServer.id)"
              />
            </div>
          </template>

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
                  </v-card-t
