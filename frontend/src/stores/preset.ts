import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import {
  IsFrpcDownloaded,
  DownloadFrpc,
  GetFrpcVersion,
  StartServer,
  StopServer,
  ExportToml,
  ImportFrpFiles,
  ExportPresetToml,
  ImportPresetFromToml,
  ExportPresetAsTomlBatch,
  SaveAppConfig,
  LoadAppConfig,
  GetAppVersion,
  GetLatestFrpcVersion,
  GetCurrentFrpcVersion,
  CompareFrpcVersions
} from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
import TOML from 'smol-toml'
import { createServerModel, createServiceModels } from '../helpers/modelConverters'
import { toSerializableServer, toSerializableService } from '../helpers/serializers'

export interface LogEntry {
  id: string
  timestamp: number
  message: string
  type: string
}

export interface Server {
  id: string
  name: string
  address: string
  port: number
  token: string
  status: string
  enabled: boolean
  logs: LogEntry[]
  uptime: number
}

export interface Service {
  id: string
  name: string
  protocol: string
  local_ip: string
  local_port: number
  remote_port: number
  use_encryption: boolean
  use_compression: boolean
  advanced_config?: string
  is_advanced?: boolean
}

export interface Preset {
  id: string
  name: string
  servers: Server[]
  services: Service[]
}

export interface DownloadProgress {
  total_bytes: number
  downloaded: number
  percentage: number
  is_complete: boolean
  is_error: boolean
  error_message?: string
  downloaded_version?: string
  version_fetch_error?: string
}

export type DownloadSource = 'github' | 'ghproxy'

const CLIPBOARD_KEY = 'frpeasy_clipboard'

const generateId = () => Math.random().toString(36).substring(2, 9)

const createDefaultServer = (name: string): Server => ({
  id: generateId(),
  name,
  address: '',
  port: 7000,
  token: '',
  status: 'offline',
  enabled: false,
  logs: [],
  uptime: 0,
})

const defaultPresets: Preset[] = [
  {
    id: generateId(),
    name: '默认预设',
    servers: [createDefaultServer('主服务器')],
    services: [],
  },
]

const loadFromStorage = async (): Promise<Preset[]> => {
  try {
    const stored = await LoadAppConfig()
    console.log('[LoadFromStorage] LoadAppConfig returned:', stored ? 'data' : 'empty')
    if (stored) {
      const parsed = TOML.parse(stored) as any
      console.log('[LoadFromStorage] Parsed presets count:', parsed.presets?.length || 0)
      if (parsed.presets && Array.isArray(parsed.presets)) {
        return parsed.presets.map((p: Preset) => ({
          ...p,
          servers: (p.servers || []).map((s: Server) => ({
            ...s,
            logs: [],
            uptime: 0,
            status: 'offline',
          })),
          services: p.services || [],
        }))
      }
    }
  } catch (e) {
    console.error('Failed to load presets from config:', e)
  }
  return []
}

async function saveToStorage(presets: Preset[]) {
  try {
    const toSave = presets.map((p) => ({
      id: p.id,
      name: p.name,
      servers: (p.servers || []).map(toSerializableServer),
      services: (p.services || []).map(toSerializableService),
    }))
    await SaveAppConfig(TOML.stringify({ presets: toSave }))
  } catch (e) {
    console.error('Failed to save presets to storage:', e)
  }
}

export const usePresetStore = defineStore('preset', () => {
  const presets = ref<Preset[]>([])
  const activePresetId = ref<string | null>(null)
  const activeServerId = ref<string | null>(null)

  const frpcDownloaded = ref(false)
  const frpcVersion = ref('')
  const downloadProgress = ref<DownloadProgress | null>(null)
  const isDownloading = ref(false)
  const downloadedVersion = ref('')
  const versionFetchError = ref('')

  const activePreset = computed(() =>
    presets.value.find((p) => p.id === activePresetId.value)
  )

  const activeServer = computed(() => {
    if (!activePreset.value) return null
    return (
      activePreset.value.servers.find((s) => s.id === activeServerId.value) ||
      null
    )
  })

  async function initPresets() {
    console.log('[InitPresets] Loading presets...')
    const loaded = await loadFromStorage()
    console.log('[InitPresets] Loaded presets count:', loaded.length)
    
    if (loaded.length > 0) {
      presets.value = loaded
      activePresetId.value = loaded[0].id
      activeServerId.value = loaded[0].servers?.[0]?.id || null
      console.log('[InitPresets] Loaded from config file')
      return
    }
    
    console.log('[InitPresets] Creating default presets')
    presets.value = defaultPresets
    activePresetId.value = defaultPresets[0].id
    activeServerId.value = defaultPresets[0].servers?.[0]?.id || null
    await saveToStorage(defaultPresets)
    console.log('[InitPresets] Default presets created and saved')
  }

  async function initFrpc() {
    try {
      frpcDownloaded.value = await IsFrpcDownloaded()
      console.log('[InitFrpc] Frpc downloaded:', frpcDownloaded.value)
      if (frpcDownloaded.value) {
        frpcVersion.value = await GetFrpcVersion()
        console.log('[InitFrpc] Frpc version:', frpcVersion.value)
      }
    } catch (e) {
      console.error('[InitFrpc] Failed to check frpc status:', e)
    }
  }

  async function startDownloadFrpc(source: DownloadSource = 'ghproxy', useDefaultVersion = false) {
    if (isDownloading.value) return
    isDownloading.value = true
    downloadProgress.value = null
    if (!useDefaultVersion) {
      downloadedVersion.value = ''
    }
    versionFetchError.value = ''
    console.log('[DownloadFrpc] Starting download with source:', source, 'useDefault:', useDefaultVersion)

    EventsOn('download:progress', (progress: DownloadProgress) => {
      downloadProgress.value = progress
      console.log('[DownloadFrpc] Progress:', progress.percentage.toFixed(1) + '%')
      if (progress.is_complete) {
        isDownloading.value = false
        downloadedVersion.value = progress.downloaded_version || ''
        console.log('[DownloadFrpc] Download completed, version:', downloadedVersion.value)
        frpcDownloaded.value = true
        initFrpc()
      }
      if (progress.version_fetch_error) {
        isDownloading.value = false
        versionFetchError.value = progress.version_fetch_error
        console.error('[DownloadFrpc] Version fetch error:', progress.version_fetch_error)
      }
      if (progress.is_error && !progress.version_fetch_error) {
        isDownloading.value = false
        console.error('[DownloadFrpc] Download failed:', progress.error_message)
      }
    })

    try {
      await DownloadFrpc(source, useDefaultVersion)
    } catch (e) {
      console.error('[DownloadFrpc] Download failed:', e)
      isDownloading.value = false
    }
  }

  function setupLogListener() {
    EventsOn('server:log', (data: { preset_id: string; server_id: string; log: LogEntry }) => {
      console.log('[ServerLog]', data.log.message)
      const preset = presets.value.find((p) => p.id === data.preset_id)
      if (!preset) return
      const server = preset.servers.find((s) => s.id === data.server_id)
      if (!server) return

      server.logs.push(data.log)
      if (server.logs.length > 100) {
        server.logs.shift()
      }
    })
  }

  function cleanupListeners() {
    EventsOff('download:progress')
    EventsOff('server:log')
  }

  function setActivePreset(id: string) {
    activePresetId.value = id
    const preset = presets.value.find((p) => p.id === id)
    if (preset && preset.servers.length > 0) {
      activeServerId.value = preset.servers[0].id
    } else {
      activeServerId.value = null
    }
  }

  function setActiveServer(id: string) {
    activeServerId.value = id
  }

  async function toggleServer(presetId: string, serverId: string) {
    console.log('[ToggleServer] Starting...', { presetId, serverId })
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset) {
      console.error('[ToggleServer] Preset not found:', presetId)
      return
    }

    if (!frpcDownloaded.value) {
      console.error('[ToggleServer] frpc not downloaded')
      return
    }

    const server = preset.servers.find((s) => s.id === serverId)
    if (!server) {
      console.error('[ToggleServer] Server not found:', serverId)
      return
    }

    server.enabled = !server.enabled
    console.log('[ToggleServer] Server enabled:', server.enabled)

    if (server.enabled) {
      server.status = 'connecting'
      server.logs = []
      server.uptime = 0

      try {
        console.log('[ToggleServer] Starting server...', { 
          server: server.name, 
          address: server.address, 
          port: server.port 
        })
        const serverModel = createServerModel(server)
        const servicesModels = createServiceModels(preset.services)
        console.log('[ToggleServer] Calling StartServer with services:', servicesModels.length)
        await StartServer(presetId, serverId, serverModel, servicesModels)
        server.status = 'online'
        console.log('[ToggleServer] Server started successfully')
      } catch (e) {
        console.error('[ToggleServer] Failed to start server:', e)
        server.status = 'error'
        server.enabled = false
      }
    } else {
      try {
        console.log('[ToggleServer] Stopping server...')
        await StopServer(presetId, serverId)
        console.log('[ToggleServer] Server stopped successfully')
      } catch (e) {
        console.error('[ToggleServer] Failed to stop server:', e)
      }
    server.status = 'offline'
      server.logs = []
      server.uptime = 0
    }

    await saveToStorage(presets.value)
  }

  async function addPreset(name: string, sourcePreset?: Preset) {
    console.log('[AddPreset] Creating preset:', name, { copyFrom: sourcePreset?.name })
    const newPreset: Preset = sourcePreset
      ? {
          id: generateId(),
          name,
          servers: sourcePreset.servers.map((s) => ({
            ...s,
            id: generateId(),
            status: 'offline',
            enabled: false,
            logs: [],
            uptime: 0,
          })),
          services: sourcePreset.services.map((s) => ({
            ...s,
            id: generateId(),
          })),
        }
      : {
          id: generateId(),
          name,
          servers: [createDefaultServer('主服务器')],
          services: [],
        }
    presets.value.push(newPreset)
    await saveToStorage(presets.value)
    activePresetId.value = newPreset.id
    activeServerId.value = newPreset.servers[0]?.id || null
    console.log('[AddPreset] Preset created:', newPreset.id)
  }

  async function deletePreset(id: string) {
    console.log('[DeletePreset] Deleting preset:', id)
    const index = presets.value.findIndex((p) => p.id === id)
    if (index === -1) return
    if (presets.value.length <= 1) return

    const preset = presets.value[index]
    for (const server of preset.servers) {
      if (server.enabled) {
        try {
          await StopServer(id, server.id)
        } catch (e) {
          console.error('[DeletePreset] Failed to stop server:', e)
        }
      }
    }

    presets.value.splice(index, 1)
    await saveToStorage(presets.value)

    if (activePresetId.value === id) {
      activePresetId.value = presets.value[0]?.id || null
      activeServerId.value = presets.value[0]?.servers?.[0]?.id || null
    }
    console.log('[DeletePreset] Preset deleted')
  }

  function copyPreset(id: string) {
    const preset = presets.value.find((p) => p.id === id)
    if (!preset) return
    const clipboardData = {
      name: preset.name,
      servers: preset.servers.map((s) => ({
        name: s.name,
        address: s.address,
        port: s.port,
        token: s.token,
      })),
      services: preset.services.map((s) => ({
        name: s.name,
        protocol: s.protocol,
        local_ip: s.local_ip,
        local_port: s.local_port,
        remote_port: s.remote_port,
        use_encryption: s.use_encryption,
        use_compression: s.use_compression,
      })),
    }
    localStorage.setItem(CLIPBOARD_KEY, JSON.stringify(clipboardData))
  }

  function hasClipboard(): boolean {
    return localStorage.getItem(CLIPBOARD_KEY) !== null
  }

  function pastePreset(name: string): boolean {
    const stored = localStorage.getItem(CLIPBOARD_KEY)
    if (!stored) return false
    try {
      const data = JSON.parse(stored)
      addPreset(name, data)
      return true
    } catch {
      return false
    }
  }

  async function addServer(presetId: string, serverData?: Partial<Server>) {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset) return

    const newServer: Server = {
      ...createDefaultServer(`服务器 ${preset.servers.length + 1}`),
      ...serverData,
      id: generateId(),
    } as Server
    preset.servers.push(newServer)
    await saveToStorage(presets.value)
    activeServerId.value = newServer.id
  }

  async function updateServer(
    presetId: string,
    serverId: string,
    updates: Partial<Server>
  ) {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset) return

    const server = preset.servers.find((s) => s.id === serverId)
    if (!server) return

    Object.assign(server, updates)
    await saveToStorage(presets.value)
  }

  async function deleteServer(presetId: string, serverId: string) {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset || preset.servers.length <= 1) return

    const server = preset.servers.find((s) => s.id === serverId)
    if (server?.enabled) {
      try {
        await StopServer(presetId, serverId)
      } catch (e) {
        console.error('Failed to stop server:', e)
      }
    }

    const index = preset.servers.findIndex((s) => s.id === serverId)
    if (index === -1) return

    preset.servers.splice(index, 1)
    await saveToStorage(presets.value)

    if (activeServerId.value === serverId) {
      activeServerId.value = preset.servers[0]?.id || null
    }
  }

  function clearLogs(presetId: string, serverId: string) {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset) return

    const server = preset.servers.find((s) => s.id === serverId)
    if (!server) return

    server.logs = []
  }

  async function updatePreset(id: string, updates: Partial<Preset>) {
    const preset = presets.value.find((p) => p.id === id)
    if (preset) {
      Object.assign(preset, updates)
      await saveToStorage(presets.value)
    }
  }

  async function updatePresetName(id: string, name: string) {
    await updatePreset(id, { name })
  }

  function exportPresetTomlContent(presetId: string): string {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset) return ''

    const tomlObj = {
      name: preset.name,
      servers: preset.servers.map((s) => ({
        name: s.name,
        address: s.address,
        port: s.port,
        token: s.token,
      })),
      services: preset.services.map((s) => ({
        name: s.name,
        protocol: s.protocol,
        local_ip: s.local_ip,
        local_port: s.local_port,
        remote_port: s.remote_port,
        use_encryption: s.use_encryption,
        use_compression: s.use_compression,
      })),
    }
    return TOML.stringify(tomlObj)
  }

  async function exportAsToml(presetId: string, serverId: string): Promise<string> {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset) return ''

    const server = preset.servers.find((s) => s.id === serverId)
    if (!server) return ''

    try {
      const serverModel = models.Server.createFrom(server)
      const servicesModels = preset.services.map(s => models.Service.createFrom(s))
      return await ExportToml(serverModel, servicesModels)
    } catch (e) {
      console.error('Failed to export TOML:', e)
      return ''
    }
  }

  async function importFrpFiles(): Promise<models.ImportResult[]> {
    try {
      return await ImportFrpFiles()
    } catch (e) {
      console.error('[ImportFrpFiles] Failed:', e)
      return [{ error: String(e) } as models.ImportResult]
    }
  }

  async function exportPresetTomlPreset(presetId: string): Promise<string> {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset) return ''
    
    const tomlContent = exportPresetTomlContent(presetId)
    if (!tomlContent) return ''
    
    try {
      return await ExportPresetToml(preset.name, tomlContent)
    } catch (e) {
      console.error('[ExportPresetTomlPreset] Failed:', e)
      return ''
    }
  }

  async function importPresetToml(): Promise<Preset | null> {
    try {
      const tomlStr = await ImportPresetFromToml()
      if (!tomlStr) return null
      
      const data = TOML.parse(tomlStr) as any
      return {
        id: generateId(),
        name: data.name || '导入的预设',
        servers: (data.servers || []).map((s: any) => ({
          ...createDefaultServer(s.name || '服务器'),
          id: generateId(),
          name: s.name,
          address: s.address,
          port: s.port,
          token: s.token,
          status: 'offline' as const,
          enabled: false,
          logs: [],
          uptime: 0,
        })),
        services: (data.services || []).map((s: any) => ({
          id: generateId(),
          name: s.name,
          protocol: s.protocol,
          local_ip: s.local_ip,
          local_port: s.local_port,
          remote_port: s.remote_port,
          use_encryption: s.use_encryption ?? false,
          use_compression: s.use_compression ?? false,
        })),
      }
    } catch (e) {
      console.error('[ImportPresetToml] Failed:', e)
      return null
    }
  }

  async function exportPresetToml(presetId: string): Promise<string> {
    const preset = presets.value.find((p) => p.id === presetId)
    if (!preset || preset.servers.length === 0) return ''

    try {
      const serversJson = JSON.stringify(preset.servers.map(s => ({
        name: s.name,
        address: s.address,
        port: s.port,
        token: s.token,
      })))

      const servicesJson = JSON.stringify(preset.services.map(s => ({
        name: s.name,
        protocol: s.protocol,
        local_ip: s.local_ip,
        local_port: s.local_port,
        remote_port: s.remote_port,
        use_encryption: s.use_encryption,
        use_compression: s.use_compression,
      })))

      return await ExportPresetAsTomlBatch(serversJson, servicesJson, preset.name)
    } catch (e) {
      console.error('[ExportPresetToml] Failed:', e)
      return ''
    }
  }

  async function addImportedPreset(preset: Preset) {
    presets.value.push(preset)
    await saveToStorage(presets.value)
    activePresetId.value = preset.id
    activeServerId.value = preset.servers[0]?.id || null
  }

  async function mergePresets(presetIds: string[], newName: string): Promise<boolean> {
    console.log('[MergePresets] Merging presets:', presetIds, 'into:', newName)
    
    if (presetIds.length < 2) {
      console.error('[MergePresets] Need at least 2 presets to merge')
      return false
    }

    const mergedPreset: Preset = {
      id: generateId(),
      name: newName,
      servers: [],
      services: [],
    }

    const serviceKeySet = new Set<string>()

    for (const presetId of presetIds) {
      const preset = presets.value.find(p => p.id === presetId)
      if (!preset) continue

      for (const server of preset.servers) {
        mergedPreset.servers.push({
          ...server,
          id: generateId(),
          name: `${preset.name} - ${server.name}`,
          status: 'offline',
          enabled: false,
          logs: [],
          uptime: 0,
        })
      }

      for (const service of preset.services) {
        const serviceKey = `${service.protocol}-${service.local_ip}-${service.local_port}-${service.remote_port}`
        
        if (serviceKeySet.has(serviceKey)) {
          console.log('[MergePresets] Skipping duplicate service:', service.name, serviceKey)
          continue
        }
        
        serviceKeySet.add(serviceKey)
        
        mergedPreset.services.push({
          ...service,
          id: generateId(),
        })
      }
    }

    if (mergedPreset.servers.length === 0) {
      console.error('[MergePresets] No servers found in selected presets')
      return false
    }

    presets.value.push(mergedPreset)
    await saveToStorage(presets.value)
    activePresetId.value = mergedPreset.id
    activeServerId.value = mergedPreset.servers[0]?.id || null
    
    console.log('[MergePresets] Merged preset created:', mergedPreset.id)
    return true
  }

  async function autoStartEnabledServers() {
    if (!frpcDownloaded.value) return
    
    for (const preset of presets.value) {
      for (const server of preset.servers) {
        if (server.enabled) {
          console.log('[AutoStart] Starting server:', server.name)
          try {
            server.status = 'connecting'
            server.logs = []
            const serverModel = createServerModel(server)
            const servicesModels = createServiceModels(preset.services)
            await StartServer(preset.id, server.id, serverModel, servicesModels)
            server.status = 'online'
            console.log('[AutoStart] Server started:', server.name)
          } catch (e) {
            console.error('[AutoStart] Failed to start server:', server.name, e)
            server.status = 'error'
            server.enabled = false
          }
        }
      }
    }
  }

  async function getAppVersion(): Promise<string> {
    return await GetAppVersion()
  }

  async function getLatestFrpcVersion(): Promise<string> {
    return await GetLatestFrpcVersion()
  }

  async function getCurrentFrpcVersion(): Promise<string> {
    return await GetCurrentFrpcVersion()
  }

  async function compareFrpcVersions(v1: string, v2: string): Promise<number> {
    return await CompareFrpcVersions(v1, v2)
  }

  return {
    presets,
    activePresetId,
    activePreset,
    activeServer,
    frpcDownloaded,
    frpcVersion,
    downloadProgress,
    isDownloading,
    initFrpc,
    startDownloadFrpc,
    setupLogListener,
    cleanupListeners,
    setActivePreset,
    setActiveServer,
    toggleServer,
    addPreset,
    deletePreset,
    copyPreset,
    pastePreset,
    hasClipboard,
    addServer,
    updateServer,
    deleteServer,
    clearLogs,
    updatePreset,
    updatePresetName,
    exportPresetTomlContent,
    exportAsToml,
    importFrpFiles,
    exportPresetTomlPreset,
    importPresetToml,
    exportPresetToml,
    addImportedPreset,
    mergePresets,
    autoStartEnabledServers,
    initPresets,
    getAppVersion,
    getLatestFrpcVersion,
    getCurrentFrpcVersion,
    compareFrpcVersions,
    downloadedVersion,
    versionFetchError,
  }
})
