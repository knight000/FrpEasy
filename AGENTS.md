# FrpEasy - Agent Guidelines

## Project Overview

FrpEasy is a Wails v2.11.0 desktop application for managing frp (Fast Reverse Proxy) client configurations. Go backend with Vue 3 + Vuetify 3 + TypeScript frontend.

## Build Commands

```bash
# Development with hot reload
wails dev

# Build complete application
wails build

# Generate Wails bindings (after adding Go methods to app.go)
wails generate module

# Frontend only
cd frontend && npm run build

# Go backend only
go build ./...
```

## Lint & Type Check

```bash
# Frontend type check (required before commits)
cd frontend && npm run type-check

# Go vet
go vet ./...
```

## Test Commands

```bash
# Run all Go tests
go test ./...

# Run tests in a specific package
go test ./internal/frpc

# Run a single test function
go test -v ./internal/frpc -run TestFunctionName

# Run with coverage
go test -cover ./...

# Frontend tests (when added)
cd frontend && npm run test
```

## Project Structure

```
FrpEasy/
├── app.go                 # Wails app bindings - ALL exported methods available to frontend
├── main.go               # Entry point, window config, lifecycle hooks
├── internal/
│   ├── frpc/             # frp client management
│   │   ├── config.go     # TOML config generation
│   │   ├── downloader.go # frpc binary download
│   │   ├── manager.go    # Process lifecycle management
│   │   └── parser.go     # TOML/INI config parsing
│   └── models/           # Shared data models
│       └── types.go      # Server, Service, Preset, LogEntry, etc.
├── frontend/
│   ├── src/
│   │   ├── App.vue           # Main component
│   │   ├── main.ts           # Vue entry
│   │   ├── components/       # Vue components
│   │   ├── stores/preset.ts  # Pinia store (main state)
│   │   └── plugins/          # Vuetify config
│   └── wailsjs/              # AUTO-GENERATED - DO NOT EDIT
└── build/bin/                # Output directory
```

## Go Code Style

### Imports (strictly ordered)
```go
import (
    "context"
    "fmt"
    "os"

    "frpeasy/internal/frpc"
    "frpeasy/internal/models"

    "github.com/wailsapp/wails/v2/pkg/runtime"
)
```

### Naming
- PascalCase: exported functions, types, constants
- camelCase: unexported functions, variables
- Consistent acronyms: `ID`, `HTTP`, `TOML`, `JSON`

### Error Handling
```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to write config: %w", err)
}

// Simple logging in app.go
fmt.Println("Failed to create directory:", err)
```

### Models
```go
// Define in internal/models/types.go
// Always include json tags
type Server struct {
    ID      string       `json:"id"`
    Name    string       `json:"name"`
    Status  ServerStatus `json:"status"`
}

// Use typed constants for enums
type ServerStatus string
const (
    StatusOnline  ServerStatus = "online"
    StatusOffline ServerStatus = "offline"
)
```

### Wails Patterns
```go
// Exported methods on App are available to frontend
func (a *App) StartServer(presetID, serverID string, server models.Server) error {
    return a.manager.Start(presetID, serverID, &server)
}

// Emit events for async updates
runtime.EventsEmit(a.ctx, "download:progress", progress)

// File dialogs
files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{...})
```

## TypeScript/Vue Code Style

### Imports (strictly ordered)
```typescript
import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { StartServer, StopServer } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
import type { LogEntry } from '@/stores/preset'
```

### Components
```vue
<script setup lang="ts">
// Props with types
const props = defineProps<{
  logs: LogEntry[]
  title?: string
}>()

// Emits with types
const emit = defineEmits<{
  clear: []
  update: [value: string]
}>()
</script>
```

### Pinia Store
```typescript
export const usePresetStore = defineStore('preset', () => {
  const presets = ref<Preset[]>([])
  
  // Console log prefix for debugging
  function addPreset(name: string) {
    console.log('[AddPreset] Creating:', name)
  }
  
  return { presets, addPreset }
})
```

### Wails Model Conversion
```typescript
// MUST use models.Type.createFrom() when passing to Go
const serverModel = models.Server.createFrom({
  id: server.id,
  name: server.name,
  // ...other fields
})
await StartServer(presetId, serverId, serverModel, servicesModels)
```

### Vue Reactivity Pattern
```vue
<!-- Use :model-value + @update:model-value instead of v-model for objects -->
<v-switch
  :model-value="server.enabled"
  @update:model-value="toggleServer(server.id)"
/>
```

### Wails Event Listeners
```typescript
// Setup in store or component
EventsOn('download:progress', (progress: DownloadProgress) => {
  downloadProgress.value = progress
})

// Always cleanup
EventsOff('download:progress')
```

## FRP Configuration Format

frp v0.61.1 TOML format:
```toml
serverAddr = "example.com"
serverPort = 7000

[auth]
token = "xxx"

[log]
to = "console"
level = "info"

[[proxies]]
name = "ssh"
type = "tcp"
localIP = "127.0.0.1"
localPort = 22
remotePort = 6000
transport.useEncryption = true
transport.useCompression = true
```

**Critical**: `useEncryption`/`useCompression` are in `[[proxies]]` as `transport.useEncryption`

## Development Workflow

### Adding Backend Functionality
1. Add method to `app.go`
2. Run `wails generate module`
3. Import and use in `frontend/src/stores/preset.ts`
4. Update UI

### Adding Data Models
1. Define in `internal/models/types.go`
2. Run `wails generate module`
3. Create matching TypeScript interface in stores

### Before Committing
1. Run `cd frontend && npm run type-check`
2. Run `go vet ./...`
3. Test with `wails dev`
