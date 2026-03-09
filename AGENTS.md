# FrpEasy - Agent Guidelines

## Project Overview

FrpEasy is a Wails v2.11.0 desktop application for managing frp (Fast Reverse Proxy) client configurations. It features a Go backend with Vue 3 + Vuetify 3 + TypeScript frontend.

## Build Commands

```bash
# Build the complete application
wails build

# Build for development with hot reload
wails dev

# Generate Wails bindings (after adding Go methods)
wails generate module

# Build frontend only
cd frontend && npm run build

# Type check frontend only
cd frontend && npm run type-check

# Build Go backend only
go build ./...
```

## Test Commands

No test files currently exist in this project. When adding tests:

```bash
# Run all Go tests
go test ./...

# Run tests in a specific package
go test ./internal/frpc

# Run a single test file
go test -v ./internal/frpc -run TestFunctionName

# Run frontend tests (when added)
cd frontend && npm run test
```

## Project Structure

```
FrpEasy/
├── app.go                 # Wails app bindings - all exported methods are available to frontend
├── main.go               # Application entry point and Wails configuration
├── internal/
│   ├── frpc/             # frp client management
│   │   ├── config.go     # TOML config generation
│   │   ├── downloader.go # frpc binary download with multiple sources
│   │   ├── manager.go    # Process lifecycle management
│   │   └── parser.go     # TOML/INI config file parsing
│   └── models/           # Data models shared between Go and TypeScript
│       └── types.go      # Server, Service, Preset, LogEntry structs
├── frontend/
│   ├── src/
│   │   ├── App.vue           # Main application component
│   │   ├── main.ts           # Vue app entry point
│   │   ├── components/       # Vue components
│   │   ├── stores/           # Pinia stores
│   │   │   └── preset.ts     # Main store for presets, servers, services
│   │   └── plugins/          # Vue plugins (Vuetify, etc.)
│   └── wailsjs/              # Auto-generated Wails bindings - DO NOT EDIT
└── build/bin/                # Built application output
```

## Go Code Style

### Imports
- Standard library imports first, separated by blank line
- Third-party imports second, separated by blank line
- Local imports (frpeasy/internal/*) last
- Use explicit import aliases only when necessary

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

### Naming Conventions
- PascalCase for exported functions, types, and constants
- camelCase for unexported functions and variables
- Acronyms should be consistent (e.g., `ID`, `HTTP`, `TOML`)
- Interface names: verb or noun describing behavior (e.g., `ProcessManager`)

### Error Handling
- Return errors as the last return value
- Wrap errors with context using `fmt.Errorf("operation failed: %w", err)`
- Never ignore errors - handle or propagate them
- Use `fmt.Println` for simple logging in app.go

### Structs and Types
- Define types in `internal/models/types.go` for data shared with frontend
- Add `json` tags for all fields that need serialization
- Use typed constants for enums (e.g., `ServerStatus`, `ServiceProtocol`)

```go
type Server struct {
    ID             string       `json:"id"`
    Name           string       `json:"name"`
    Status         ServerStatus `json:"status"`
}
```

### Exported Methods for Frontend
- All exported methods on `App` struct are automatically available to frontend
- Use Wails runtime for dialogs: `runtime.OpenFileDialog`, `runtime.SaveFileDialog`
- Emit events for async updates: `runtime.EventsEmit(ctx, "event:name", data)`

## TypeScript/Vue Code Style

### Imports
- Vue imports first: `import { ref, computed } from 'vue'`
- Third-party imports second
- Local imports with `@/` alias third
- Wails bindings last

```typescript
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'
import { StartServer, StopServer } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
```

### Components
- Use `<script setup lang="ts">` syntax
- Define props with `defineProps<{}>()`
- Define emits with `defineEmits<{}>()`
- Use Vuetify components with `v-` prefix

### Stores (Pinia)
- Use Composition API style with `defineStore('name', () => {})`
- Export interfaces for types
- Add `console.log` with `[FunctionName]` prefix for debugging

```typescript
export const usePresetStore = defineStore('preset', () => {
  const presets = ref<Preset[]>([])
  
  function addPreset(name: string) {
    console.log('[AddPreset] Creating:', name)
  }
  
  return { presets, addPreset }
})
```

### Wails Model Usage
- Use `models.Type.createFrom(source)` when passing data to Go backend
- This ensures proper serialization for Wails bindings

```typescript
const serverModel = models.Server.createFrom({
  id: server.id,
  name: server.name,
  // ...other fields
})
await StartServer(presetId, serverId, serverModel, servicesModels)
```

### Vue Reactivity
- Use `:model-value` + `@update:model-value` instead of `v-model` for switches
- Prevents direct mutation issues with complex objects

```vue
<v-switch
  :model-value="server.enabled"
  @update:model-value="toggleServer(server.id)"
/>
```

## FRP Configuration Format

The application uses frp v0.61.1 TOML format:

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

**Important**: 
- `useEncryption` and `useCompression` belong in `[[proxies]]` as `transport.useEncryption`
- Log config uses `[log]` table format, not dotted notation

## Important Files to Update

When adding new backend functionality:
1. Add method to `app.go`
2. Run `wails generate module` to update bindings
3. Import and use in `frontend/src/stores/preset.ts`
4. Update UI in `frontend/src/App.vue` or components

When adding new data models:
1. Define in `internal/models/types.go`
2. Run `wails generate module`
3. Create matching TypeScript interface in `stores/preset.ts`
