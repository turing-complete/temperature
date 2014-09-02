package system

import (
	"fmt"
	"path"
	"testing"

	"github.com/goesd/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestNewFromTGFF(t *testing.T) {
	app, platform, err := NewFromTGFF(findFixture("002_040"))

	assert.Success(err, t)
	assert.Equal(len(app.Tasks), 40, t)
	assert.Equal(len(platform.Cores), 2, t)
}

func findFixture(name string) string {
	return path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
}
