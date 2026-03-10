import { ref } from 'vue'

interface SnackbarState {
	show: boolean
	message: string
	color: string
}

const snackbar = ref<SnackbarState>({
	show: false,
	message: '',
	color: 'success',
})

type SnackbarColor = 'success' | 'error' | 'info' | 'warning'

export function useSnackbar() {
	function showSnackbar(message: string, color: SnackbarColor = 'success') {
		snackbar.value = { show: true, message, color }
	}

	function showSuccess(message: string) {
		showSnackbar(message, 'success')
	}

	function showError(message: string) {
		showSnackbar(message, 'error')
	}

	function showInfo(message: string) {
		showSnackbar(message, 'info')
	}

	function showWarning(message: string) {
		showSnackbar(message, 'warning')
	}

	return {
		snackbar,
		showSnackbar,
		showSuccess,
		showError,
		showInfo,
		showWarning,
	}
}
