package cache

import (
	"fmt"

	"github.com/bashidogames/gdvm"
)

type Clear struct{}

func (c *Clear) Run(app *gdvm.App) error {
	err := app.Cache.Clear()
	if err != nil {
		return fmt.Errorf("cannot clear cache: %w", err)
	}

	return nil
}

type Cache struct {
	Clear Clear `cmd:"" help:"Clear the cache"`
}
