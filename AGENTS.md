# FrpEasy - Agent Guidelines

## Project Overview

FrpEasy is a Wails v2.11.0 desktop application for managing frp (Fast Reverse Proxy) client configurations. Go backend with Vue 3 + Vuetify 3 + TypeScript frontend.

## Build Commands

```bash
wails dev                    # Development with hot reload
wails build                  # Build complete application
wails generate module        # Generate Wails bindings (after adding Go methods to app.go)
cd frontend && npm run build # Frontend only
go build ./...               # Go backend only
```

## Lint & Test Commands

```bash
cd frontend && npm run type-check  # Frontend type check (required before commits)
go vet ./...                       # Go lint

go test ./...                                   # Run all Go tests
go test ./internal/frpc                         # Run tests in specific package
go test -v ./internal/frpc -run TestFunctionName # Run single test
go test -cover ./...                            # Run with coverage
```

## Project Structure

```
FrpEasy/
├── app.go                      # Wails app bindings - ALL exported methods available to frontend
├── main.go                     # Entry point, window config, lifecycle hooks
├── internal/
│   ├── config/config.go        # App configuration (TOML load/save)
│   ├── frpc/                   # frp client management
│   │   ├── config.go           # TOML config generation for frpc
│   │   ├── downloader.go       # frpc binary download
│   │   ├── manager.go          # Process lifecycle management
│   │   └── parser.go           # TOML/INI config parsing
│   └── models/types.go         # Shared data models (Server, Service, Preset, LogEntry)
├── frontend/
│   ├── src/
│   │   ├── App.vue             # Main component
│   │   ├── components/         # Vue components
│   │   └── stores/preset.ts    # Pinia store (main state, Wails bindings)
│   └── wailsjs/                # AUTO-GENERATED - DO NOT EDIT
├── .github/workflows/release.yml  # GitHub Actions auto-release
└── {exe_dir}/frpeasy/          # Runtime data directory
    ├── config.toml             # App configuration
    ├── bin/                    # frpc binary
    └── configs/                # frpc runtime configs
```

## Go Code Style

### Imports (strictly ordered)
```go
import (
    // 1. Standard library
    "context"
    "fmt"
    "os"

    // 2. Internal packages
    "frpeasy/internal/config"
    "frpeasy/internal/frpc"
    "frpeasy/internal/models"

    // 3. External packages
    "github.com/wailsapp/wails/v2/pkg/runtime"
    "github.com/google/uuid"
)
```

### Naming & Conventions
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
```go
id := uuid.New().String()[:8]
```

## TypeScript/Vue Code Style

### Imports (strictly ordered)
```typescript
// 1. Vue core
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { defineStore } from 'pinia'

// 2. Wails runtime (import in stores, not components)
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { StartServer, StopServer } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'

// 3. External libraries
import TOML from 'smol-toml'

// 4. Local imports
import type { LogEntry } from '@/stores/preset'
```

### Key Patterns

**Async Storage Functions** - always use `await`:
```typescript
await saveToStorage(presets.value)
const loaded = await loadFromStorage()
```

**Wails Model Conversion**:
```typescript
const serverModel = models.Server.createFrom({ id: server.id, name: server.name })
await StartServer(presetId, serverId, serverModel, servicesModels)
```

**Vue Reactivity**:
```vue
<v-switch
  :model-value="server.enabled"
  @update:model-value="toggleServer(server.id)"
/>
```

**Context Menu (Wails/WebView compatible)** - avoid `v-menu` with absolute positioning:
```vue
<v-card v-if="showContextMenu" class="context-menu"
  :style="{ left: x + 'px', top: y + 'px' }" @click.stop>
  <v-list density="compact" bg-color="#2d2d2d">
    <v-list-item @click="handleAction">...</v-list-item>
  </v-list>
</v-card>
```
```css
.context-menu { position: fixed; z-index: 1000; }
```

**Save selection before context menu** (selection clears when menu appears):
```typescript
function onContextMenu(e: MouseEvent) {
  selectedText.value = window.getSelection()?.toString() || ''
  contextMenuX.value = e.clientX
  contextMenuY.value = e.clientY
  showContextMenu.value = true
}
```

**Global Event Listeners** - always clean up:
```typescript
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
```

## FRP Configuration Format (frp v0.61.1)

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

## Field Naming

All fields use `snake_case` (e.g., `local_ip`, `local_port`, `use_encryption`, `preset_id`, `server_id`)

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
1. `cd frontend && npm run type-check`
2. `go vet ./...`
3. Test with `wails dev`

### Creating Release
```bash
git tag v1.0.0
git push origin v1.0.0
```
GitHub Actions will auto-build and create release with Windows exe.

## Key Discoveries

1. **Windows Console Fix**: `exec.Command` shows console unless `SysProcAttr` with `HideWindow: true` is set
2. **TOML for Configs**: Use TOML (snake_case) for config files, JSON (camelCase) for frontend communication
3. **v-menu Positioning**: Fails in Wails WebView - use `position: fixed` with v-card instead
4. **Selection Loss**: Menu opening clears text selection - save selection before showing menu
5. **Wails Import Path**: TypeScript cannot resolve `../../wailsjs/go/main/App` in Vue components - import in Pinia store instead
6. **smol-toml**: `TOML.parse()` returns `TomlTable` type - use `as any` to avoid TypeScript errors
7. **v-checkbox Centering**: Use `.v-checkbox .v-selection-control { justify-content: center; }` with `vertical-align: middle` on td
