//go:build dev

package layouts

func shouldRenderDevReload() bool {
	return true
}

func devReloadInit() string {
	return "@get('/reload', {retryMaxCount: Infinity, retryInterval: 20, retryMaxWaitMs: 200})"
}
