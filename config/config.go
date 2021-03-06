/*
      Package config holds network configuration file and secret.
*/
package config

import "os"

var (
    Network = &NetworkConfig{
      RdbAddress: "localhost:28015",
      RDB: "courses",
    }
    Local = &LocalConfig{
      DefaultFN: os.Getenv("GOPATH") + "/src/github.com/ahermida/Goberon/courses.html",
      CatFN: os.Getenv("GOPATH") + "/src/github.com/ahermida/Goberon/catalog.html",
    }
)

type NetworkConfig struct {
    RdbAddress string //rethinkDB port
    RDB string        //rethink database name
}

type LocalConfig struct {
    DefaultFN string //default filename
    CatFN string     //filename for catalog (if it's set)
}
