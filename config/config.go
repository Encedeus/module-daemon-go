package config

import (
	"errors"
	"fmt"
	_ "github.com/joho/godotenv"
)

var (
	InvalidCompilerOptimizationLevel = errors.New("invalid wasm compiler optimization level")
)

// const DefaultLocation = "/etc/encedeus"
const DefaultLocation = "./"

type Configuration struct {
	StorageLocationPath string                `hcl:"storage_location"`
	Server              ServerConfiguration   `hcl:"server,block"`
	DB                  DatabaseConfiguration `hcl:"database,block"`
	Auth                AuthConfiguration     `hcl:"auth,block"`
	CDN                 CDNConfiguration      `hcl:"cdn,block"`
	Modules             ModulesConfiguration  `hcl:"modules,block"`
	Skyhook             SkyhookConfiguration  `hcl:"skyhook,block"`
}

/*func NewConfiguration() Configuration {
	secretAccess := string(security.RandomBytes(24))
	secretRefresh := string(security.RandomBytes(24))

	config := Configuration{
		Server: ServerConfiguration{
			Host: "localhost",
			AddPort: 8080,
		},
		DB: DatabaseConfiguration{
			Host:     "localhost",
			AddPort:     5432,
			User:     "postgres",
			DBName:   "PanelDB",
			Password: "root",
		},
		Auth: AuthConfiguration{
			JWTSecretAccess:  secretAccess,
			JWTSecretRefresh: secretRefresh,
		},
		CDN: CDNConfiguration{
			Directory: "./pfp",
		},
		Modules: ModulesConfiguration{
			ModulesDirectory:              "./modules",
			WasmCompilerOptimizationLevel: LPerf2,
		},
		Skyhook: SkyhookConfiguration{
			DefaultPort: 8000,
			MinFreeRAM:  2147000000,
		},
	}

	return config
}*/

type ServerConfiguration struct {
	Host string `hcl:"host"`
	Port int    `hcl:"port"`
}

type WasmCompilerOptimizationLevel int

const (
	LPerf0 WasmCompilerOptimizationLevel = iota
	LPerf1
	LPerf2
	LPerf3
	LSize1
	LSize2
)

/*func (col WasmCompilerOptimizationLevel) GetLibraryOptimizationLevel() wasmedge.CompilerOptimizationLevel {
	optLevel := wasmedge.CompilerOptLevel_O2

	switch col {
	case LPerf0:
		optLevel = wasmedge.CompilerOptLevel_O0
	case LPerf1:
		optLevel = wasmedge.CompilerOptLevel_O1
	case LPerf2:
		optLevel = wasmedge.CompilerOptLevel_O2
	case LPerf3:
		optLevel = wasmedge.CompilerOptLevel_O3
	case LSize1:
		optLevel = wasmedge.CompilerOptLevel_Os
	case LSize2:
		optLevel = wasmedge.CompilerOptLevel_Oz
	}

	return optLevel
}*/

type ModulesConfiguration struct {
	// Relative to the main "encedeus" directory
	ModulesDirectory string `hcl:"modules_dir"`
	// Optimization level for the WasmEdge compiler
	WasmCompilerOptimizationLevel WasmCompilerOptimizationLevel `hcl:"compiler_optimization_level"`
}

type DatabaseConfiguration struct {
	Host     string `hcl:"host"`
	Port     int    `hcl:"port"`
	User     string `hcl:"user"`
	DBName   string `hcl:"name"`
	Password string `hcl:"password"`
}

type AuthConfiguration struct {
	JWTSecretAccess  string `hcl:"jwt_secret_access"`
	JWTSecretRefresh string `hcl:"jwt_secret_refresh"`
}

type CDNConfiguration struct {
	Directory string `hcl:"dir"`
}

type SkyhookConfiguration struct {
	DefaultPort         uint16 `hcl:"default_port"`
	MinFreeRAM          uint64 `hcl:"min_free_ram"`
	MinFreeDisk         uint64 `hcl:"min_free_disk"`
	MinFreeLogicalCores uint32 `hcl:"min_free_logical_cores"`
}

func (s *ServerConfiguration) URI() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
