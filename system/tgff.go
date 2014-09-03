package system

import (
	"errors"
	"strconv"
	"strings"

	"github.com/goesd/format/tgff"
)

func LoadFromTGFF(path string) (App, Platform, error) {
	r, err := tgff.ParseFile(path)

	if err != nil {
		return App{}, Platform{}, err
	}

	a, err := loadApp(r.Graphs)

	if err != nil {
		return App{}, Platform{}, err
	}

	p, err := loadPlatform(r.Tables)

	if err != nil {
		return App{}, Platform{}, err
	}

	return a, p, nil
}

func loadApp(graphs []tgff.Graph) (App, error) {
	a := App{}

	if len(graphs) != 1 {
		return a, errors.New("need exactly one task graph")
	}

	size := uint32(len(graphs[0].Tasks))

	a.Tasks = make([]Task, size)

	for _, task := range graphs[0].Tasks {
		i, err := extractTaskID(task.Name, size)

		if err != nil {
			return a, err
		}

		a.Tasks[i].ID = i
		a.Tasks[i].Type = task.Type
	}

	for _, arc := range graphs[0].Arcs {
		i, err := extractTaskID(arc.From, size)

		if err != nil {
			return a, err
		}

		j, err := extractTaskID(arc.To, size)

		if err != nil {
			return a, err
		}

		a.Tasks[i].Children = append(a.Tasks[i].Children, &a.Tasks[j])
	}

	return a, nil
}

func loadPlatform(tables []tgff.Table) (Platform, error) {
	p := Platform{}

	if len(tables) == 0 {
		return p, errors.New("need at least one table")
	}

	p.Cores = make([]Core, len(tables))

	var err error

	for i := range tables {
		p.Cores[i], err = loadCore(tables[i])

		if err != nil {
			return p, err
		}

		if i == 0 {
			continue
		}

		if len(p.Cores[i-1].Time) != len(p.Cores[i].Time) {
			return p, errors.New("inconsistent table data")
		}
	}

	return p, nil
}

func loadCore(table tgff.Table) (Core, error) {
	c := Core{}

	var tycol, tmcol, pwcol *tgff.Column

	for i := range table.Columns {
		col := &table.Columns[i]

		name := strings.ToLower(col.Name)

		if strings.Index(name, "type") >= 0 {
			tycol = col
		} else if strings.Index(name, "time") >= 0 {
			tmcol = col
		} else if strings.Index(name, "power") >= 0 {
			pwcol = col
		}

		if tycol != nil && tmcol != nil && pwcol != nil {
			break
		}
	}

	if tycol == nil || tmcol == nil || pwcol == nil {
		return c, errors.New("need columns named type, time, and power")
	}

	size := len(tycol.Data)

	for i := 0; i < size; i++ {
		if int(tycol.Data[i]) != i {
			return c, errors.New("data should be sorted by type")
		}
	}

	c.Time = make([]float64, size)
	copy(c.Time, tmcol.Data)

	c.Power = make([]float64, size)
	copy(c.Power, pwcol.Data)

	return c, nil
}

func extractTaskID(name string, total uint32) (uint32, error) {
	if !strings.HasPrefix(name, "t0_") {
		return 0, errors.New("unknown task naming scheme")
	}

	id, err := strconv.ParseInt(name[3:], 10, 0)

	if err != nil || id < 0 || uint32(id) >= total {
		return 0, errors.New("unknown task indexing scheme")
	}

	return uint32(id), nil
}
