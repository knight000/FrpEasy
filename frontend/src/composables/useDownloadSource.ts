const DOWNLOAD_SOURCES = [
  { value: 'ghproxy' as const, label: 'GHProxy (国内加速)' },
  { value: 'github' as const, label: 'GitHub (直连)' },
]

export type DownloadSourceValue = typeof DOWNLOAD_SOURCES[number]['value']

export function useDownloadSource() {
  function getSourceLabel(source: string): string {
    const item = DOWNLOAD_SOURCES.find(s => s.value === source)
    return item?.label || source
  }

  return {
    downloadSources: DOWNLOAD_SOURCES,
    getSourceLabel,
  }
}
