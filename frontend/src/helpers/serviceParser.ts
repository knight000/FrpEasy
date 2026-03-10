import TOML from 'smol-toml'
import type { Service } from '../stores/preset'

const FRPEASY_PREFIX = "#FRPEASY#"

interface TemplateDisplayInfo {
	namePattern?: string
	protocol?: string
	portRange?: string
}

function parseFrpeasyPrefix(config: string): TemplateDisplayInfo | null {
	const firstLine = config.split('\n')[0]
	if (!firstLine.startsWith(FRPEASY_PREFIX)) {
	 return null
	}

	const data = firstLine.substring(FRPEASY_PREFIX.length)
	const parts = data.split('#')
	const info: TemplateDisplayInfo = {}

	for (const part of parts) {
		const [key, value] = part.split('=')
	 if (key === 'name') {
      info.namePattern = value
    } else if (key === 'protocol') {
      info.protocol = value
    } else if (key === 'ports') {
      info.portRange = value
    }
  }

  return info
}

export function parseAdvancedConfigToBasic(service: Service): void {
  if (!service.is_advanced || !service.advanced_config) return

  if (service.advanced_config.startsWith(FRPEASY_PREFIX)) {
    const info = parseFrpeasyPrefix(service.advanced_config)
    if (info) {
      if (info.namePattern) service.name = info.namePattern
      if (info.protocol) service.protocol = info.protocol
      if (info.portRange) service._portRange = info.portRange
    }
    return
  }

  if (service.advanced_config.includes('{{') && service.advanced_config.includes('}}')) {
    service.name = "模板服务"
    return
  }

  try {
    const parsed = TOML.parse(service.advanced_config) as any

    service.protocol = 'TCP'
    service.local_ip = ''
    service.local_port = 0
    service.remote_port = 0
    service.use_encryption = false
    service.use_compression = false

    if (parsed.name && typeof parsed.name === 'string') {
      service.name = parsed.name
    }
    if (parsed.type && typeof parsed.type === 'string') {
      const typeMap: Record<string, string> = {
        tcp: 'TCP',
        udp: 'UDP',
        http: 'HTTP',
        https: 'HTTPS',
      }
      const lowerType = parsed.type.toLowerCase()
      service.protocol = typeMap[lowerType] || parsed.type.toUpperCase()
    }
    if (parsed.localIP !== undefined) {
      service.local_ip = String(parsed.localIP)
    }
    if (parsed.localIp !== undefined) {
      service.local_ip = String(parsed.localIp)
    }
    if (parsed.localPort !== undefined) {
      service.local_port = Number(parsed.localPort)
    }
    if (parsed.remotePort !== undefined) {
      service.remote_port = Number(parsed.remotePort)
    }
    if (parsed.transport?.useEncryption !== undefined) {
      service.use_encryption = Boolean(parsed.transport.useEncryption)
    }
    if (parsed.transport?.useCompression !== undefined) {
      service.use_compression = Boolean(parsed.transport.useCompression)
    }
  } catch (e) {
    console.warn('Failed to parse advanced config:', e)
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
