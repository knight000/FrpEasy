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
```

## Project Structure

```
FrpEasy/
├── app.go                 # Wails app bindings - ALL exported methods available to frontend
├── main.go               # Entry point, window config, lifecycle hooks
├── internal/
│   ├── config/           # App configuration (TOML)
│   │   └── config.go     # Config load/save, TOML <-> JSON conversion
│   ├── frpc/             # frp client management
│   │   ├── config.go     # TOML config generation for frpc
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
├── build/bin/                # Output directory
└── {exe_dir}/frpeasy/        # Runtime data
    ├── config.toml           # App configuration (presets, servers, services)
    ├── bin/                  # frpc binary
    └── configs/              # frpc runtime configs
```

## Configuration Storage

- **App Config**: `{exe_dir}/frpeasy/config.toml` - Stores presets, servers, services, enabled state
- **Clipboard**: `localStorage` - Temporary copy/paste data for presets
- **Auto-start**: Servers with `enabled: true` auto-start on app launch

## Go Code Style

### Imports (strictly ordered)
```go
import (
    "context"
    "fmt"
    "os"

    "frpeasy/internal/config"
    "frpeasy/internal/frpc"
    "frpeasy/internal/models"

    "github.com/wailsapp/wails/v2/pkg/runtime"
)
```

### Naming
- PascalCase: exported functions, types, constants
- camelCase: unexported functions, variables
- Consistent acronyms: `ID`, `HTTP`, `TOML`, `JSON`, `IP`

### Error Handling
```go
if err != nil {
    return fmt.Errorf("failed to write config: %w", err)
}
```

### Config Module Usage
```go
// Load config
appConfig, err := config.LoadConfig(filepath.Join(a.dataDir, "config.toml"))

// Save config
config.SaveConfig(path, appConfig)

// Convert for frontend
jsonStr := config.ToJSON(appConfig)

// Parse from frontend
appConfig, err := config.FromJSON(jsonStr)
```

## TypeScript/Vue Code Style

### Imports (strictly ordered)
```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { StartServer, StopServer, SaveAppConfig, LoadAppConfig } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
```

### Async Storage Functions
All storage functions are async - always use `await`:
```typescript
await saveToStorage(presets.value)
const loaded = await loadFromStorage()
```

### Wails Model Conversion
```typescript
const serverModel = models.Server.createFrom({
  id: server.id,
  name: server.name,
  // ...other fields
})
await StartServer(presetId, serverId, serverModel, servicesModels)
```

### Vue Reactivity Pattern
```vue
<v-switch
  :model-value="server.enabled"
  @update:model-value="toggleServer(server.id)"
/>
```

## FRP Configuration Format

frp v0.61.1 TOML format:
```toml
serverAddr = "example.com"
serverPort = 7000

[auth]
token = "xxx"

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
