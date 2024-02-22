package module

import (
	"errors"
	"fmt"
	protoapi "github.com/Encedeus/module-daemon-go/proto/go"
	"slices"
)

var (
	ErrUnsupportedCrater  = errors.New("unsupported crater")
	ErrUnsupportedVariant = errors.New("unsupported crater variant")
)

type Crater struct {
	Id          string
	Name        string
	Variants    []*Variant
	Description string
	Provider    *Module
}

type ServerRunningState int

const (
	STARTING ServerRunningState = iota
	RUNNING
	RESTARTING
	STOPPING
	STOPPED
)

type Variant struct {
	Id                string
	Name              string
	Description       string
	DataDirectoryPath string
	GetConsoleLogs    func(c *Crater, m *Module, s *protoapi.Server) []byte
	StartServer       func(c *Crater, m *Module, s *protoapi.Server) error
	StopServer        func(c *Crater, m *Module, s *protoapi.Server) error
	RestartServer     func(c *Crater, m *Module, s *protoapi.Server) error
	GetRunningState   func(c *Crater, m *Module, s *protoapi.Server) ServerRunningState
	CreateServer      func(opts *protoapi.ServersCreateRequest) (*protoapi.ServersCreateResponse, error)
}

func (m *Module) RegisterCrater(c Crater) {
	m.Craters = append(m.Craters, &c)
}

type CraterHandler struct {
	RegisteredCraters *[]*Crater
}

func HasVariant(variant string, craters []*Crater) (bool, *Variant) {
	for _, c := range craters {
		for _, v := range c.Variants {
			if v.Name == variant {
				return true, v
			}
		}
	}

	return false, nil
}

func HasCrater(id string, craters []*Crater) bool {
	/*	for _, c := range craters {
		fmt.Printf("%v %v", c.Id, id)
		if c.Id == id {
			return true
		}
	}*/
	return slices.ContainsFunc(craters, func(crater *Crater) bool {
		fmt.Printf("Crater: %v %v\n", crater.Id, id)
		return crater.Id == id
	})

	//return false
}

func (ch *CraterHandler) CreateServer(opts *protoapi.ServersCreateRequest) (*protoapi.ServersCreateResponse, error) {
	fmt.Printf("Create req: %+v\n", *opts)
	supportsCrater := HasCrater(opts.Crater, *ch.RegisteredCraters)
	if !supportsCrater {
		return nil, ErrUnsupportedCrater
	}

	supportsVariant, variant := HasVariant(opts.CraterVariant, *ch.RegisteredCraters)
	if !supportsVariant {
		return nil, ErrUnsupportedVariant
	}

	resp, err := variant.CreateServer(opts)
	if err != nil {
		return nil, err
	}

	return resp, err
}
