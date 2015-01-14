package temperature

import (
	"path"
)

const (
	fixturePath = "fixtures"
)

func findFixture(name string) string {
	return path.Join(fixturePath, name)
}
