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

func TestLoadTGFF(t *testing.T) {
	plat, app, err := LoadTGFF(findFixture("002_040"))

	assert.Success(err, t)
	assert.Equal(len(plat.Cores), 2, t)
	assert.Equal(len(app.Tasks), 40, t)

	tasks := []struct {
		children []uint32
	}{
		{[]uint32{1}},
		{[]uint32{2, 3, 23, 32}},
		{[]uint32{4, 10, 11, 12}},
		{[]uint32{9, 10, 26}},
		{[]uint32{5, 6, 7, 8}},

		{[]uint32{10, 11, 27}},
		{nil},
		{[]uint32{9}},
		{[]uint32{9}},
		{[]uint32{13, 14, 15}},

		{[]uint32{12}},
		{[]uint32{12, 31}},
		{[]uint32{17, 18, 24}},
		{[]uint32{16, 25}},
		{[]uint32{36}},

		{[]uint32{19, 20, 23, 29}},
		{[]uint32{31}},
		{[]uint32{22, 35}},
		{[]uint32{21, 23}},
		{[]uint32{20, 24}},

		{[]uint32{24, 26}},
		{[]uint32{27}},
		{[]uint32{28, 29}},
		{[]uint32{31}},
		{[]uint32{30}},

		{[]uint32{37, 38}},
		{nil},
		{[]uint32{33, 34}},
		{[]uint32{29}},
		{nil},

		{nil},
		{[]uint32{36}},
		{nil},
		{nil},
		{nil},

		{nil},
		{nil},
		{nil},
		{[]uint32{39}},
		{nil},
	}

	for i, task := range tasks {
		assert.Equal(app.Tasks[i].Children, task.children, t)
	}
}

func TestExtractTaskNumber(t *testing.T) {
	scenarios := []struct {
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
		result, err := extractTaskID(s.name, s.total)

		if s.success {
			assert.Success(err, t)
		} else {
			assert.Failure(err, t)
		}

		assert.Equal(result, s.result, t)
	}
}

func BenchmarkLoadTGFF(b *testing.B) {
	path := findFixture("002_040")

	for i := 0; i < b.N; i++ {
		LoadTGFF(path)
	}
}

func findFixture(name string) string {
	return path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
}
