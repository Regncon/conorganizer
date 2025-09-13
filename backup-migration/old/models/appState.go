package models

import (
	"fmt"

	"fyne.io/fyne/v2/data/binding"
)

type AppError struct {
	Error   binding.Bool
	Message binding.String
}

type DBState struct {
	IsLoaded binding.Bool
	IsValid  binding.Bool
	Path     binding.String
}

type ENVstate struct {
	Path     binding.String
	IsLoaded binding.Bool
	IsValid  binding.Bool
}

type PrefixState struct {
	Value    binding.String
	IsUnique binding.Bool
}

type AppState struct {
	DB       DBState
	ENV      ENVstate
	Prefix   PrefixState
	OnChange binding.String
}

func NewAppState() *AppState {
	appState := &AppState{
		DB: DBState{
			IsLoaded: binding.NewBool(),
			IsValid:  binding.NewBool(),
			Path:     binding.NewString(),
		},
		ENV: ENVstate{
			IsLoaded: binding.NewBool(),
			IsValid:  binding.NewBool(),
			Path:     binding.NewString(),
		},
		Prefix: PrefixState{
			Value:    binding.NewString(),
			IsUnique: binding.NewBool(),
		},
		OnChange: binding.NewString(),
	}

	// Manage states
	appState.OnChange.AddListener(binding.NewDataListener(func() {
		fmt.Println("app state changed")
	}))
	appState.DB.Path.AddListener(binding.NewDataListener(func() {
		dbPath, _ := appState.DB.Path.Get()
		if dbPath == "" {
			appState.DB.IsLoaded.Set(false)
		} else {
			appState.OnChange.Set(dbPath)
			appState.DB.IsLoaded.Set(true)
		}
	}))
	appState.ENV.Path.AddListener(binding.NewDataListener(func() {
		envPath, _ := appState.ENV.Path.Get()
		if envPath == "" {
			appState.ENV.IsLoaded.Set(false)
		} else {
			appState.OnChange.Set(envPath)
			appState.ENV.IsLoaded.Set(true)
		}
	}))
	appState.Prefix.Value.AddListener(binding.NewDataListener(func() {

	}))

	return appState
}
