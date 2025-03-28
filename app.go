package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Config struct {
	Workspace string
}

type App struct {
	ctx context.Context
	Config Config
}

func NewApp() *App {
	return &App{
		Config: Config{},
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	err := a.setup()
	if err != nil {
		runtime.LogError(ctx, err.Error())
	}
}

func (a *App) setup() error {
	p, err := a.createWorkspace()
	if err != nil {
		return err
	}
	a.Config.Workspace = p
	return nil
}


func (a *App) createWorkspace() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := path.Join(dir, ".copet")
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			return "", nil
		}
	} else if err != nil {
			return "", nil
	}
	return p, nil
}

func (a *App) GetSnippetPaths() ([]string, error) {
	paths:= []string{}
	p := a.Config.Workspace
	if len(p) == 0 {
		err := "workspace path is not yet defined"
		runtime.LogError(a.ctx, err)		
		return paths, errors.New(err)
	}

	target := path.Join(p)
	dir, err := os.Open(target)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("error while opening directory %v, %v\n", p, err))
		return paths, err
	}

	entries, err := dir.ReadDir(0)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("error while reading directory %v, %v\n", p, err))
		return paths, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		paths = append(paths, filepath.Join(dir.Name(), entry.Name()))
	}

	return paths, nil
}

func (a *App) ReadFile(p string) (string, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}
	str := string(b)
	return str, err
}

func (a *App) CreateFile(p string) (string, error) {
	target := filepath.Join(a.Config.Workspace, p)
	dir := filepath.Dir(target)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err.Error())
	}

	f, err := os.Create(target)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err.Error())
	}
	defer f.Close()
	return target, nil
}