package system

import (
	"errors"

	"github.com/goesd/format/tgff"
)

type Graph interface {
}

type Table interface {
}

func NewFromTGFF(path string) (*App, *Platform, error) {
	result, err := tgff.ParseFile(path)

	if err != nil {
		return nil, nil, err
	}

	if len(result.Graphs) != 1 {
		return nil, nil, errors.New("need to have exactly one task graph")
	}

	tables := make([]Table, len(result.Tables))
	for i := range result.Tables {
		tables[i] = Table(result.Tables[i])
	}

	return New(result.Graphs[0], tables...)
}

func New(graph Graph, tables ...Table) (*App, *Platform, error) {
	app, err := NewApp(graph)

	if err != nil {
		return nil, nil, err
	}

	platform, err := NewPlatform(tables...)

	if err != nil {
		return nil, nil, err
	}

	return app, platform, nil
}
