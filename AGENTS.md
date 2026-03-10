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

## Go Code Style

### Imports (strictly ordered)
```go
import (
    // Standard library
    "context"
    "fmt"
    "os"

    // Internal packages
    "frpeasy/internal/config"
    "frpeasy/internal/frpc"
    "frpeasy/internal/models"

    // External packages
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

### Windows Subprocess - Hide Console Window
```go
func hideWindow(cmd *exec.Cmd) {
    if runtime.GOOS == "windows" {
        cmd.SysProcAttr = &syscall.SysProcAttr{
            HideWindow:    true,
            CreationFlags: 0x08000000,
        }
    }
}
```
Apply to ALL subprocess commands (frpc, taskkill, powershell, etc.) to prevent console windows.

### UUID Generation
Use `github.com/google/uuid` for ID generation to prevent collisions:
```go
id := uuid.New().String()[:8]
```

## TypeScript/Vue Code Style

### Imports (strictly ordered)
```typescript
// Vue core
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { defineStore } from 'pinia'

// Wails runtime
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { StartServer, StopServer } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'

// External libraries
import TOML from 'smol-toml'

// Local imports
import type { LogEntry } from '@/stores/preset'
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

### Context Menu Implementation (Wails/WebView compatible)
Avoid `v-menu` with absolute positioning in WebView. Use fixed positioning:
```vue
<v-card
  v-if="showContextMenu"
  class="context-menu"
  :style="{ left: x + 'px', top: y + 'px' }"
  @click.stop
>
  <v-list density="compact" bg-color="#2d2d2d">
    <v-list-item @click="handleAction">...</v-list-item>
  </v-list>
</v-card>
```
```css
.context-menu { position: fixed; z-index: 1000; }
```
Save selection before menu opens (selection clears when menu appears):
```typescript
function onContextMenu(e: MouseEvent) {
  selectedText.value = window.getSelection()?.toString() || ''
  contextMenuX.value = e.clientX
  contextMenuY.value = e.clientY
  showContextMenu.value = true
}
```

### Global Event Listeners
Always clean up in `onUnmounted`:
```typescript
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  document.addEventListener('keydown', handleKeydown)
})
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  document.removeEventListener('keydown', handleKeydown)
})
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

## File Naming Conventions

| Type | Pattern | Example |
|------|---------|---------|
| Preset export | `预设-{name}.toml` | `预设-生产环境.toml` |
| frp config export | `frpc-{serverName}.toml` | `frpc-主服务器.toml` |
| Config storage | `config.toml` | `{exe_dir}/frpeasy/config.toml` |

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

## Key Discoveries

1. **Windows Console Fix**: `exec.Command` shows console unless `SysProcAttr` with `HideWindow: true` is set
2. **TOML for Configs**: Use TOML (snake_case) for config files, JSON (camelCase) for frontend communication
3. **v-menu Positioning**: Fails in Wails WebView - use fixed positioning instead
4. **Selection Loss**: Menu opening clears text selection - save selection before showing menu
