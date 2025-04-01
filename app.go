package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

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

type dirNode struct {
	Files []string `json:"files"`
	Subdirs map[string]*dirNode `json:"subdirs"`
}

func NewDirNode() *dirNode {
	return &dirNode{
		Files: []string{},
		Subdirs: map[string]*dirNode{},
	}
}

func (a *App) scanDirs() (*dirNode, error) {
	root := NewDirNode()	
	err := filepath.Walk(a.Config.Workspace, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(a.Config.Workspace, path)
		if rel == "." {
			return nil
		}

		parent := root
		parts := strings.Split(rel, string(os.PathSeparator))

		// set parent dir to the current directory of where the file at
		for i := 0; i < len(parts)-1; i++ {
			dir := parts[i]
			if _, exists := parent.Subdirs[dir]; !exists {
				parent.Subdirs[dir] = NewDirNode()
			} 
			parent = parent.Subdirs[dir]
		}

		if info.IsDir() {
			parent.Subdirs[info.Name()] = NewDirNode()
		} else {
			parent.Files = append(parent.Files, rel)
		}

		return nil
	})
	return root, err
}

func (a *App) GetSnippetPaths() (*dirNode, error) {
	node, err := a.scanDirs()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (a *App) ReadFile(p string) (string, error) {
	b, err := os.ReadFile(filepath.Join(a.Config.Workspace, p))
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "", err
	}
	str := string(b)
	return str, err
}

func getUniquePath(path string) string {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)
	dir := filepath.Dir(path)

	counter := 1
	newPath := path

	for {
		newName := fmt.Sprintf("%s(%d)%s", name, counter, ext)
		newPath = filepath.Join(dir, newName)

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			break 
		}
		counter++
	}

	return newPath
}

func (a *App) Move(src string, dst string) error {
	srcPath := filepath.Join(a.Config.Workspace, src)
	dstPath := filepath.Join(a.Config.Workspace, dst)

	fileInfo, err := os.Stat(dstPath) 
	if err != nil {
		if os.IsNotExist(err) {
			dstParent := filepath.Dir(dstPath)
			if err := os.MkdirAll(dstParent, os.ModePerm); err != nil {
				runtime.LogError(a.ctx, "failed to create destination directory: "+err.Error())
				return fmt.Errorf("failed to create destination directory: %v", err.Error())
			}
		}
		return err	
	} 

	if fileInfo.IsDir() {
		dstPath = filepath.Join(dstPath, filepath.Base(srcPath))
	}

if _, err := os.Stat(dstPath); err == nil {
		dstPath = getUniquePath(dstPath)
	}

	err = os.Rename(srcPath, dstPath)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err
	} 
	return nil
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
	return p, nil
}

func (a *App) UpdateFile(p string, s string) error {
	f, err := os.OpenFile(filepath.Join(a.Config.Workspace, p), os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	// Write new content
	if _, err := f.Write([]byte(s)); err != nil {
		return err
	}

	return nil
}