//go:build !dev

package layouts

func shouldRenderDevReload() bool {
	return false
}

func devReloadInit() string {
	return ""
}
