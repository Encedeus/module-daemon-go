package crater

import (
	"github.com/Encedeus/module-daemon-go/module"
)

type Crater struct {
	Name        string
	Variants    []*Variant
	Description string
	Module      *module.Module
}

type Variant struct {
	Name        string
	Description string
}
