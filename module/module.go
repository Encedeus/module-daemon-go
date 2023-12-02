package module

type Manifest struct {
    Name             string   `hcl:"name"`
    Authors          []string `hcl:"authors"`
    Verison          string   `hcl:"version"`
    FrontendMainFile string   `hcl:"frontend_main"`
    // BackendMainFile  string   `hcl:"backend_main"`
}

type Module struct {
    Port     Port
    Manifest Manifest
}

type Port uint16

type Configuration struct {
    Port     Port
    Manifest Manifest
}
