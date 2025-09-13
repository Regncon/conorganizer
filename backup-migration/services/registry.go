package services

import "fyne.io/fyne/v2"

type Registry struct {
	Config string
	App    fyne.App
}

func NewRegistry() *Registry {
	registry := &Registry{}
	return registry
}
