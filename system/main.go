package system

import (
	"errors"
	"strings"

	"github.com/goesd/format/tgff"
)

func LoadFromTGFF(path string) (App, Platform, error) {
	r, err := tgff.ParseFile(path)

	if err != nil {
		return App{}, Platform{}, err
	}

	a, err := loadAppFromTGFF(r.Graphs)

	if err != nil {
		return App{}, Platform{}, err
	}

	p, err := loadPlatformFromTGFF(r.Tables)

	if err != nil {
		return App{}, Platform{}, err
	}

	return a, p, nil
}

func loadAppFromTGFF(graphs []tgff.Graph) (App, error) {
	if len(graphs) != 1 {
		return App{}, errors.New("need exactly one task graph")
	}

	a := App{
		Tasks: make([]Task, len(graphs[0].Tasks)),
	}

	for i, t := range graphs[0].Tasks {
		a.Tasks[i].Type = t.Type
	}

	return a, nil
}

func loadPlatformFromTGFF(tables []tgff.Table) (Platform, error) {
	if len(tables) == 0 {
		return Platform{}, errors.New("need at least one table")
	}

	p := Platform{
		Cores: make([]Core, len(tables)),
	}

	var err error

	for i := range tables {
		p.Cores[i], err = loadCoreFromTGFF(tables[i])

		if err != nil {
			return Platform{}, err
		}

		if i == 0 {
			continue
		}

		if len(p.Cores[i-1].Time) != len(p.Cores[i].Time) {
			return Platform{}, errors.New("inconsistent table data")
		}
	}

	return p, nil
}

func loadCoreFromTGFF(table tgff.Table) (Core, error) {
	var tycol, tmcol, pwcol *tgff.Column

	for i := range table.Columns {
		c := &table.Columns[i]

		name := strings.ToLower(c.Name)

		if strings.Index(name, "type") >= 0 {
			tycol = c
		} else if strings.Index(name, "time") >= 0 {
			tmcol = c
		} else if strings.Index(name, "power") >= 0 {
			pwcol = c
		}

		if tycol != nil && tmcol != nil && pwcol != nil {
			break
		}
	}

	if tycol == nil || tmcol == nil || pwcol == nil {
		return Core{}, errors.New("need columns named type, time, and power")
	}

	size := len(tycol.Data)

	for i := 0; i < size; i++ {
		if int(tycol.Data[i]) != i {
			return Core{}, errors.New("data should be sorted by type")
		}
	}

	time := make([]float64, size)
	copy(time, tmcol.Data)

	power := make([]float64, size)
	copy(power, pwcol.Data)

	return Core{Time: time, Power: power}, nil
}
