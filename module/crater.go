package module

import (
	"errors"
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
	//GetConsoleLogs    func(srv protoapi.Server) []byte
	StartServer   func(srv protoapi.Server) error
	StopServer    func(srv protoapi.Server) error
	RestartServer func(srv protoapi.Server) error
	//GetRunningState   func(c *Crater, m *Module, s *protoapi.Server) ServerRunningState
	CreateServer func(opts *protoapi.ServersCreateRequest) (*protoapi.ServersCreateResponse, error)
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
			if v.Id == variant {
				return true, v
			}
		}
	}

	return false, nil
}

func HasCrater(id string, craters []*Crater) bool {
	return slices.ContainsFunc(craters, func(crater *Crater) bool {
		return crater.Id == id
	})
}

func (ch *CraterHandler) CreateServer(opts *protoapi.ServersCreateRequest) (*protoapi.ServersCreateResponse, error) {
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
