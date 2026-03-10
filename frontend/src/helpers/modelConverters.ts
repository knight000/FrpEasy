import { models } from '../../wailsjs/go/models'
import type { Server, Service } from '../stores/preset'

export function createServerModel(server: Server): models.Server {
	return models.Server.createFrom({
		id: server.id,
		name: server.name,
		address: server.address,
		port: server.port,
		token: server.token,
		status: server.status,
		enabled: server.enabled,
		logs: server.logs,
		uptime: server.uptime,
	})
}

export function createServiceModel(service: Service): models.Service {
	return models.Service.createFrom({
		id: service.id,
		name: service.name,
		protocol: service.protocol,
		local_ip: service.local_ip,
		local_port: service.local_port,
		remote_port: service.remote_port,
		use_encryption: service.use_encryption,
		use_compression: service.use_compression,
		advanced_config: service.advanced_config,
		is_advanced: service.is_advanced,
	})
}

export function createServiceModels(services: Service[]): models.Service[] {
	return services.map(createServiceModel)
}
