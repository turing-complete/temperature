package system

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-eslab/format/tgff"
)

func loadTGFF(path string) (*Platform, *Application, error) {
	r, err := tgff.ParseFile(path)
	if err != nil {
		return nil, nil, err
	}

	p, err := loadPlatform(r.Tables)
	if err != nil {
		return nil, nil, err
	}

	a, err := loadApplication(r.Graphs)
	if err != nil {
		return nil, nil, err
	}

	return p, a, nil
}

func loadPlatform(tables []tgff.Table) (*Platform, error) {
	size := len(tables)

	if size == 0 {
		return nil, errors.New("need at least one table")
	}

	p := &Platform{
		Cores: make([]Core, size),
	}

	var err error

	for _, table := range tables {
		i := table.ID

		if i >= uint16(size) {
			return nil, errors.New("unknown table indexing scheme")
		}

		p.Cores[i], err = loadCore(table)

		if err != nil {
			return nil, err
		}
	}

	rows := len(p.Cores[0].Time)

	for i := 1; i < size; i++ {
		if rows != len(p.Cores[i].Time) {
			return nil, errors.New("inconsistent table data")
		}
	}

	return p, nil
}

func loadApplication(graphs []tgff.Graph) (*Application, error) {
	if len(graphs) != 1 {
		return nil, errors.New("need exactly one task graph")
	}

	size := uint16(len(graphs[0].Tasks))

	a := &Application{
		Tasks: make([]Task, size),
	}

	for _, task := range graphs[0].Tasks {
		i, err := extractTaskID(task.Name, size)

		if err != nil {
			return nil, err
		}

		a.Tasks[i].ID = i
		a.Tasks[i].Type = task.Type
	}

	for _, arc := range graphs[0].Arcs {
		i, err := extractTaskID(arc.From, size)

		if err != nil {
			return nil, err
		}

		j, err := extractTaskID(arc.To, size)

		if err != nil {
			return nil, err
		}

		a.Tasks[i].Children = append(a.Tasks[i].Children, uint16(j))
		a.Tasks[j].Parents = append(a.Tasks[j].Parents, uint16(i))
	}

	return a, nil
}

func loadCore(table tgff.Table) (Core, error) {
	c := Core{
		ID: table.ID,
	}

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

func extractTaskID(name string, total uint16) (uint16, error) {
	if !strings.HasPrefix(name, "t0_") {
		return 0, errors.New("unknown task naming scheme")
	}

	id, err := strconv.ParseInt(name[3:], 10, 0)

	if err != nil || id < 0 || uint16(id) >= total {
		return 0, errors.New("unknown task indexing scheme")
	}

	return uint16(id), nil
}
