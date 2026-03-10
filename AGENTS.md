# FrpEasy - Agent Guidelines

## Project Overview

FrpEasy is a Wails v2.11.0 desktop application for managing frp (Fast Reverse Proxy) client configurations. Go backend with Vue 3 + Vuetify 3 + TypeScript frontend.

## Build/Lint/Test Commands

```bash
# Build
wails dev                    # Development with hot reload
wails build                  # Build complete application
wails generate module        # Generate Wails bindings (after modifying app.go or models)
cd frontend && npm run build # Frontend only

# Lint
cd frontend && npm run type-check  # Frontend type check (required before commits)
go vet ./...                       # Go lint

# Test
go test ./...                                    # Run all Go tests
go test ./internal/frpc                          # Run tests in specific package
go test -v ./internal/frpc -run TestFunctionName # Run single test with verbose
go test -cover ./...                             # Run with coverage
```

## Project Structure

```
FrpEasy/
├── app.go                      # Wails bindings - ALL exported methods available to frontend
├── main.go                     # Entry point, window config, lifecycle hooks
├── internal/
│   ├── config/config.go        # App config (TOML load/save)
│   ├── frpc/                   # frp client: config gen, download, process mgmt, parsing
│   └── models/types.go         # Data models (Server, Service, Preset, LogEntry)
├── frontend/
│   ├── src/
│   │   ├── App.vue             # Main component
│   │   ├── components/         # Vue components
│   │   ├── composables/        # Vue composables (useStatus, useSnackbar)
│   │   ├── helpers/            # Helper functions (modelConverters, serializers)
│   │   └── stores/preset.ts    # Pinia store (state + Wails bindings)
│   └── wailsjs/                # AUTO-GENERATED - DO NOT EDIT
└── {exe_dir}/frpeasy/          # Runtime data: config.toml, bin/, configs/
```

## Go Code Style

### Imports (strictly ordered with blank lines between groups)
```go
import (
	"context"
	"fmt"
	"os"

	"frpeasy/internal/config"
	"frpeasy/internal/frpc"
	"frpeasy/internal/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/google/uuid"
)
```

### Naming & Formatting
- PascalCase: exported functions, types, constants (`StartServer`, `ServerStatus`)
- camelCase: unexported functions, variables (`hideWindow`, `configPath`)
- Struct tags: `snake_case` for JSON/TOML (`json:"local_ip"`)

### Error Handling
```go
if err != nil {
    return fmt.Errorf("failed to write config: %w", err)
}
```

### Windows Subprocess - Always Hide Console
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
Apply to ALL subprocess commands (frpc, taskkill, powershell, etc.).

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

// 2. Wails runtime (import in stores, NOT in Vue components)
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { StartServer, StopServer } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'

// 3. External libraries
import TOML from 'smol-toml'

// 4. Local imports
import type { LogEntry } from '@/stores/preset'
import { createServiceModels, createServerModel } from '@/helpers/modelConverters'
import { toSerializableService, toSerializableServer } from '@/helpers/serializers'
import { useStatus } from '@/composables/useStatus'
import { useSnackbar } from '@/composables/useSnackbar'
```

### Key Patterns
- **Wails Model Conversion**: Use helpers from `@/helpers/modelConverters`
- **Data Serialization**: Use helpers from `@/helpers/serializers`
- **Status Display**: Use `useStatus` composable for consistent status mapping
- **Snackbar**: Use `useSnackbar` composable for notifications
- **Async Storage**: Always use `await` with storage functions
- **Event Listeners**: Always clean up in `onUnmounted()` - use `EventsOff()` for Wails events
- **Context Menu**: `v-menu` with absolute positioning fails in Wails WebView - use `position: fixed` with v-card

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
```

**Critical**: `useEncryption`/`useCompression` are inside `[[proxies]]` as `transport.useEncryption`.

## Naming Conventions

### Field Naming
All fields use `snake_case` in JSON/TOML (e.g., `local_ip`, `local_port`, `use_encryption`, `is_advanced`).

### File Naming
| Type | Pattern | Example |
|------|---------|---------|
| Preset export | `预设-{name}.toml` | `预设-生产环境.toml` |
| frp config export | `frpc-{serverName}.toml` | `frpc-主服务器.toml` |

## Development Workflow

### Adding Backend Functionality
1. Add exported method to `app.go`
2. Run `wails generate module`
3. Import and use in `frontend/src/stores/preset.ts`

### Adding/Modifying Data Models
1. Update `internal/models/types.go`
2. Run `wails generate module`
3. Update TypeScript interface in `stores/preset.ts`

### Before Committing
1. `cd frontend && npm run type-check`
2. `go vet ./...`
3. Test with `wails dev`

### Creating Release
```bash
git tag v1.0.0
git push origin v1.0.0
```
GitHub Actions auto-builds and creates release.

## Helper Files

### modelConverters.ts
```typescript
createServerModel(server: Server): models.Server
createServiceModel(service: Service): models.Service
createServiceModels(services: Service[]): models.Service[]
```

### serializers.ts
```typescript
toSerializableServer(server: Server)
toSerializableService(service: Service)
toSerializablePreset(preset: Preset)
```

### useStatus.ts
```typescript
useStatus(status: ServerStatus): { dotClass, chipColor, text }
getStatusDotClass(status: ServerStatus): string
getStatusChipColor(status: ServerStatus): string
getStatusText(status: ServerStatus): string
```

### useSnackbar.ts
```typescript
useSnackbar(): { snackbar, showSnackbar, showSuccess, showError, showInfo, showWarning }
```

## Key Discoveries

1. **Windows Console**: `exec.Command` shows console unless `SysProcAttr` with `HideWindow: true` is set
2. **TOML vs JSON**: TOML (snake_case) for config files, JSON (snake_case via tags) for Wails IPC
3. **v-menu Positioning**: Fails in Wails WebView - use `position: fixed` with v-card
4. **Selection Loss**: Menu opening clears text selection - save before showing
5. **Wails Imports**: TypeScript cannot resolve `../../wailsjs/go/main/App` in Vue components - import in Pinia store
6. **smol-toml**: `TOML.parse()` returns `TomlTable` type - use `as any` for TypeScript
7. **Service Advanced Mode**: When `is_advanced=true` and `advanced_config` is non-empty, use it directly; otherwise use basic fields

## Code Optimization TODO

### High Priority
- [x] Unify struct definitions in config package and models package

### Medium Priority
- [x] Extract generic confirm dialog component
- [x] Optimize download source selector UI duplication

### Low Priority
- [x] Consider moving ImportResult to models package
- [ ] Optimize wails generate module post-sync workflow

## New Helper Files Created

- `frontend/src/helpers/modelConverters.ts` - Wails model conversion utilities
- `frontend/src/helpers/serializers.ts` - Data serialization utilities
- `frontend/src/helpers/serviceParser.ts` - Service config parsing (FRPEASY prefix, TOML)
- `frontend/src/composables/useStatus.ts` - Status mapping composable
- `frontend/src/composables/useSnackbar.ts` - Snackbar utilities
- `frontend/src/composables/useDownloadSource.ts` - Download source configuration
- `frontend/src/components/ConfirmDialog.vue` - Generic confirmation dialog component
- `frontend/src/components/DownloadSourceSelect.vue` - Download source selector component
- `internal/config/config.go` - Added conversion functions between config and models
- `internal/frpc/template_parser.go` - frp Go template parsing (port range)

## Completed Optimizations

- [x] Fixed duplicate `serverAddr` assignment in `parser.go`
- [x] Simplified `isBasicProxyField` function in `parser.go`
- [x] Removed nested `generateID` function in `manager.go`
- [x] Created helper files: modelConverters.ts, serializers.ts, useStatus.ts, useSnackbar.ts
- [x] Updated preset.ts to use createServiceModels helper
- [x] Updated preset.ts to use toSerializableServer/Service helper
- [x] Merged startDownloadFrpc and startDownloadFrpcWithDefault functions
- [x] Updated App.vue to use useStatus composable
- [x] Updated PresetSidebar.vue to use shared status functions
- [x] Replaced snackbar patterns with useSnackbar composable in App.vue
- [x] Unified config and models package struct definitions with conversion functions
- [x] Simplified LogConsole.vue copy function wrappers
- [x] Extracted ConfirmDialog component for reusable confirmation dialogs
- [x] Created DownloadSourceSelect component and useDownloadSource composable
- [x] Moved ImportResult to models package
- [x] Support Go template port range mapping (parseNumberRangePair)

## Go Template Support

FrpEasy supports frp's Go template syntax for port range mapping:

### Example Configuration
```toml
{{- range $_, $v := parseNumberRangePair "6000-6006,6007" "6000-6006,6007" }}
[[proxies]]
name = "tcp-{{ $v.First }}"
type = "tcp"
localPort = {{ $v.First }}
remotePort = {{ $v.Second }}
{{- end }}
```

### Implementation Details
- **Import**: Backend parses Go template and extracts display information
- **Display**: `display_ports` field shows "6000-6007" (range format)
- **Export**: Template configs are exported as-is, no metadata prefix
- **Storage**: `advanced_config` contains raw template, `display_ports` contains port range

### Files Modified
- `internal/frpc/template_parser.go` - Parse Go template blocks
- `internal/frpc/parser.go` - Extract display info (name pattern, protocol, ports)
- `internal/models/types.go` - Added `DisplayPorts` field
- `frontend/src/stores/preset.ts` - Added `display_ports` to Service interface
- `frontend/src/App.vue` - `displayPort()` uses `display_ports` field
- `frontend/src/helpers/serviceParser.ts` - Simplified (no prefix parsing)
