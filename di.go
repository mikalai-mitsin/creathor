package main

import (
	"github.com/018bf/creathor/internal/configs"
	"github.com/018bf/creathor/internal/generators/containers"
	"os"
	"path"
)

func CreateDI(data *configs.Project) error {
	directories := []string{
		path.Join(destinationPath, "internal", "containers"),
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	files := []*Template{
		{
			SourcePath:      "templates/internal/containers/fx.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "containers", "fx.go"),
			Name:            "Uber FX DI container",
		},
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	fx := containers.FxContainer{Project: data}
	if err := fx.SyncFxModule(); err != nil {
		return err
	}
	if err := fx.SyncMigrateContainer(); err != nil {
		return err
	}
	if data.GRPCEnabled {
		if err := fx.SyncGrpcContainer(); err != nil {
			return err
		}
	}
	if data.GatewayEnabled {
		if err := fx.SyncGatewayContainer(); err != nil {
			return err
		}
	}
	if data.RESTEnabled {
		if err := fx.SyncRestContainer(); err != nil {
			return err
		}
	}
	return nil
}
