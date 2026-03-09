/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

interface Window {
  go?: {
    main?: {
      App?: {
        StartPreset: (presetName: string) => Promise<void>
        StopPreset: (presetName: string) => Promise<void>
        GetDataDir: () => Promise<string>
      }
    }
  }
}
