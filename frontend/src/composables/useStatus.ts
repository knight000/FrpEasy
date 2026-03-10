import type { Server } from '../stores/preset'

type ServerStatus = Server['status']

interface StatusInfo {
	dotClass: string
	chipColor: string
	text: string
}

const STATUS_MAP: Record<ServerStatus, StatusInfo> = {
	online: { dotClass: 'bg-success', chipColor: 'success', text: '在线' },
	connecting: { dotClass: 'bg-warning', chipColor: 'warning', text: '连接中' },
	error: { dotClass: 'bg-error', chipColor: 'error', text: '错误' },
	offline: { dotClass: 'bg-grey', chipColor: 'grey', text: '离线' },
}

const DEFAULT_STATUS: StatusInfo = STATUS_MAP.offline

export function useStatus(status: ServerStatus): StatusInfo {
	return STATUS_MAP[status] || DEFAULT_STATUS
}

export function getStatusDotClass(status: ServerStatus): string {
	return STATUS_MAP[status]?.dotClass || DEFAULT_STATUS.dotClass
}

export function getStatusChipColor(status: ServerStatus): string {
	return STATUS_MAP[status]?.chipColor || DEFAULT_STATUS.chipColor
}

export function getStatusText(status: ServerStatus): string {
	return STATUS_MAP[status]?.text || DEFAULT_STATUS.text
}
