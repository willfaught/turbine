// TODO
package turbine

import (
	"fmt"

	"github.com/willfaught/forklift"
)

func Run(path, template string) error {
	_, err := forklift.LoadPackage(path)
	if err != nil {
		return fmt.Errorf("cannot load package: %v", err)
	}
	// p.TypesInfo.
	return nil
}
