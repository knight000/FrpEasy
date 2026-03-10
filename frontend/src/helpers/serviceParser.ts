import TOML from 'smol-toml'
import type { Service } from '../stores/preset'

export function parseAdvancedConfigToBasic(service: Service): boolean {
  if (!service.advanced_config) return false

  if (service.advanced_config.includes('{{') && service.advanced_config.includes('}}')) {
    return false
  }

  try {
    const parsed = TOML.parse(service.advanced_config) as any

    if (!Array.isArray(parsed.proxies) || parsed.proxies.length === 0) {
      return false
    }

    const proxy = parsed.proxies[0]

    service.protocol = 'TCP'
    service.local_ip = '127.0.0.1'
    service.local_port = 0
    service.remote_port = 0
    service.use_encryption = false
    service.use_compression = false

    if (proxy.name && typeof proxy.name === 'string') {
      service.name = proxy.name
    }
    if (proxy.type && typeof proxy.type === 'string') {
      const typeMap: Record<string, string> = {
        tcp: 'TCP',
        udp: 'UDP',
        http: 'HTTP',
        https: 'HTTPS',
      }
      const lowerType = proxy.type.toLowerCase()
      service.protocol = typeMap[lowerType] || proxy.type.toUpperCase()
    }
    if (proxy.localIP !== undefined) {
      service.local_ip = String(proxy.localIP)
    }
    if (proxy.localIp !== undefined) {
      service.local_ip = String(proxy.localIp)
    }
    if (proxy.localPort !== undefined) {
      service.local_port = Number(proxy.localPort)
    }
    if (proxy.remotePort !== undefined) {
      service.remote_port = Number(proxy.remotePort)
    }
    if (proxy.transport?.useEncryption !== undefined) {
      service.use_encryption = Boolean(proxy.transport.useEncryption)
    }
    if (proxy.transport?.useCompression !== undefined) {
      service.use_compression = Boolean(proxy.transport.useCompression)
    }

    return true
  } catch (e) {
    console.warn('Failed to parse advanced config:', e)
    return false
  }
}

export function createDefaultService(): Service {
  return {
    id: '',
    name: '',
    protocol: 'TCP',
    local_ip: '127.0.0.1',
    local_port: 0,
    remote_port: 0,
    use_encryption: false,
    use_compression: false,
    advanced_config: '',
    is_advanced: false,
  }
}
