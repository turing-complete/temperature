package system

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ready-steady/format/tgff"
)

func loadTGFF(path string) (*Platform, *Application, error) {
	result, err := tgff.ParseFile(path)
	if err != nil {
		return nil, nil, err
	}

	platform, err := loadPlatform(result.Tables)
	if err != nil {
		return nil, nil, err
	}

	application, err := loadApplication(result.Graphs)
	if err != nil {
		return nil, nil, err
	}

	return platform, application, nil
}

func loadPlatform(tables []tgff.Table) (*Platform, error) {
	size := len(tables)

	if size == 0 {
		return nil, errors.New("need at least one table")
	}

	platform := &Platform{
		Cores: make([]Core, size),
	}

	var err error

	for _, table := range tables {
		i := table.ID

		if i >= uint16(size) {
			return nil, errors.New("unknown table indexing scheme")
		}

		platform.Cores[i], err = loadCore(table)

		if err != nil {
			return nil, err
		}
	}

	rows := len(platform.Cores[0].Time)

	for i := 1; i < size; i++ {
		if rows != len(platform.Cores[i].Time) {
			return nil, errors.New("inconsistent table data")
		}
	}

	return platform, nil
}

func loadApplication(graphs []tgff.Graph) (*Application, error) {
	if len(graphs) != 1 {
		return nil, errors.New("need exactly one task graph")
	}

	size := uint16(len(graphs[0].Tasks))

	application := &Application{
		Tasks: make([]Task, size),
	}

	for _, task := range graphs[0].Tasks {
		i, err := extractTaskID(task.Name, size)

		if err != nil {
			return nil, err
		}

		application.Tasks[i].ID = i
		application.Tasks[i].Type = task.Type
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

		application.Tasks[i].Children = append(application.Tasks[i].Children, uint16(j))
		application.Tasks[j].Parents = append(application.Tasks[j].Parents, uint16(i))
	}

	return application, nil
}

func loadCore(table tgff.Table) (Core, error) {
	core := Core{
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
		return core, errors.New("need columns named type, time, and power")
	}

	size := len(tycol.Data)

	for i := 0; i < size; i++ {
		if int(tycol.Data[i]) != i {
			return core, errors.New("data should be sorted by type")
		}
	}

	core.Time = make([]float64, size)
	copy(core.Time, tmcol.Data)

	core.Power = make([]float64, size)
	copy(core.Power, pwcol.Data)

	return core, nil
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
