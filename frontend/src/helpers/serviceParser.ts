import { NormalizeService } from '../../wailsjs/go/main/App'
import { models } from '../../wailsjs/go/models'
import type { Service } from '../stores/preset'

export function createServiceInput(partial?: Partial<Service>): models.ServiceInput {
  return {
    name: partial?.name ?? '',
    protocol: partial?.protocol ?? 'TCP',
    local_ip: partial?.local_ip ?? '',
    local_port: partial?.local_port ?? 0,
    remote_port: partial?.remote_port ?? 0,
    use_encryption: partial?.use_encryption ?? false,
    use_compression: partial?.use_compression ?? false,
    advanced_config: partial?.advanced_config ?? '',
    is_advanced: partial?.is_advanced ?? false,
  }
}

export function toService(model: models.Service): Service {
  return {
    id: model.id,
    name: model.name,
    protocol: model.protocol,
    local_ip: model.local_ip,
    local_port: model.local_port,
    remote_port: model.remote_port,
    use_encryption: model.use_encryption,
    use_compression: model.use_compression,
    advanced_config: model.advanced_config,
    is_advanced: model.is_advanced,
    display_ports: model.display_ports ?? '',
    display_local_ports: model.display_local_ports ?? '',
  }
}

export async function createService(partial?: Partial<Service>): Promise<Service> {
  const input = createServiceInput(partial)
  const model = await NormalizeService(input)
  return toService(model)
}
