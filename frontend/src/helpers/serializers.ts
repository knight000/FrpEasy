import type { Server, Service, Preset } from '../stores/preset'

export function toSerializableServer(server: Server) {
	return {
		id: server.id,
		name: server.name,
		address: server.address,
		port: server.port,
		token: server.token,
		enabled: server.enabled,
	}
}

export function toSerializableService(service: Service) {
	if (service.is_advanced) {
		return {
			id: service.id,
			name: service.name,
			is_advanced: service.is_advanced,
			advanced_config: service.advanced_config,
		}
	}
	return {
		id: service.id,
		name: service.name,
		protocol: service.protocol,
		local_ip: service.local_ip,
		local_port: service.local_port,
		remote_port: service.remote_port,
		use_encryption: service.use_encryption,
		use_compression: service.use_compression,
	}
}

export function toSerializablePreset(preset: Preset) {
	return {
		id: preset.id,
		name: preset.name,
		servers: preset.servers.map(toSerializableServer),
		services: preset.services.map(toSerializableService),
	}
}
