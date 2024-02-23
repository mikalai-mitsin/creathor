package main

import (
	"os"
	"path"

	"github.com/018bf/creathor/internal/generators/pkg/containers"

	"github.com/018bf/creathor/internal/configs"
)

func CreateDI(data *configs.Project) error {
	directories := []string{
		path.Join(destinationPath, "internal", "pkg", "containers"),
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	fx := containers.NewFxContainer(data)
	if err := fx.Sync(); err != nil {
		return err
	}
	return nil
}
