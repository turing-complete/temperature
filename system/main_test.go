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

func TestLoadFromTGFF(t *testing.T) {
	app, platform, err := LoadFromTGFF(findFixture("002_040"))

	assert.Success(err, t)
	assert.Equal(len(app.Tasks), 40, t)
	assert.Equal(len(platform.Cores), 2, t)

	tasks := []struct{
		children []uint
	}{
		{[]uint{1}},
		{[]uint{2, 3, 32}},
		{[]uint{9, 10, 26}},
		{[]uint{5, 6, 7, 8}},
	}

	for i, task := range tasks {
		assert.Equal(len(app.Tasks[i].Children), len(task.children), t)
	}
}

func TestExtractTaskNumber(t *testing.T) {
	scenarios := []struct{
		name    string
		total   uint32
		result  uint32
		success bool
	}{
		{"t0_0", 50, 0, true},
		{"t0_42", 50, 42, true},
		{"t0_42", 43, 42, true},
		{"t0_42", 42, 0, false},
		{"t1_42", 50, 0, false},
		{"t0_-2", 50, 0, false},
	}

	for _, s := range scenarios {
		result, err := extractTaskNumber(s.name, s.total)

		if s.success {
			assert.Success(err, t)
		} else {
			assert.Failure(err, t)
		}

		assert.Equal(result, s.result, t)
	}
}

func BenchmarkLoadFromTGFF(b *testing.B) {
	path := findFixture("002_040")

	for i := 0; i < b.N; i++ {
		LoadFromTGFF(path)
	}
}

func findFixture(name string) string {
	return path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
}
