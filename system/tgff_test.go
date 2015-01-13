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
		children []uint16
	}{
		{[]uint16{1}},
		{[]uint16{2, 3, 23, 32}},
		{[]uint16{4, 10, 11, 12}},
		{[]uint16{9, 10, 26}},
		{[]uint16{5, 6, 7, 8}},

		{[]uint16{10, 11, 27}},
		{nil},
		{[]uint16{9}},
		{[]uint16{9}},
		{[]uint16{13, 14, 15}},

		{[]uint16{12}},
		{[]uint16{12, 31}},
		{[]uint16{17, 18, 24}},
		{[]uint16{16, 25}},
		{[]uint16{36}},

		{[]uint16{19, 20, 23, 29}},
		{[]uint16{31}},
		{[]uint16{22, 35}},
		{[]uint16{21, 23}},
		{[]uint16{20, 24}},

		{[]uint16{24, 26}},
		{[]uint16{27}},
		{[]uint16{28, 29}},
		{[]uint16{31}},
		{[]uint16{30}},

		{[]uint16{37, 38}},
		{nil},
		{[]uint16{33, 34}},
		{[]uint16{29}},
		{nil},

		{nil},
		{[]uint16{36}},
		{nil},
		{nil},
		{nil},

		{nil},
		{nil},
		{nil},
		{[]uint16{39}},
		{nil},
	}

	for i, task := range tasks {
		assert.Equal(application.Tasks[i].Children, task.children, t)
	}
}

func TestExtractTaskNumber(t *testing.T) {
	scenarios := []struct {
		name    string
		total   uint16
		result  uint16
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
