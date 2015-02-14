package system

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestLoadTGFF(t *testing.T) {
	platform, application, err := loadTGFF(findFixture("002_040"))

	assert.Success(err, t)
	assert.Equal(len(platform.Cores), 2, t)
	assert.Equal(len(application.Tasks), 40, t)

	tasks := []struct {
		children []uint
	}{
		{[]uint{1}},
		{[]uint{2, 3, 23, 32}},
		{[]uint{4, 10, 11, 12}},
		{[]uint{9, 10, 26}},
		{[]uint{5, 6, 7, 8}},

		{[]uint{10, 11, 27}},
		{nil},
		{[]uint{9}},
		{[]uint{9}},
		{[]uint{13, 14, 15}},

		{[]uint{12}},
		{[]uint{12, 31}},
		{[]uint{17, 18, 24}},
		{[]uint{16, 25}},
		{[]uint{36}},

		{[]uint{19, 20, 23, 29}},
		{[]uint{31}},
		{[]uint{22, 35}},
		{[]uint{21, 23}},
		{[]uint{20, 24}},

		{[]uint{24, 26}},
		{[]uint{27}},
		{[]uint{28, 29}},
		{[]uint{31}},
		{[]uint{30}},

		{[]uint{37, 38}},
		{nil},
		{[]uint{33, 34}},
		{[]uint{29}},
		{nil},

		{nil},
		{[]uint{36}},
		{nil},
		{nil},
		{nil},

		{nil},
		{nil},
		{nil},
		{[]uint{39}},
		{nil},
	}

	for i, task := range tasks {
		assert.Equal(application.Tasks[i].Children, task.children, t)
	}
}

func TestExtractTaskNumber(t *testing.T) {
	scenarios := []struct {
		name    string
		total   uint
		result  uint
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
		_, _, err := loadTGFF(path)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func findFixture(name string) string {
	return path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
}
