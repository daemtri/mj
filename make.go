// +build ignore

package main

import (
	"fmt"
	"github.com/duanqy/builder"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	fmt.Println("GOPATH", dir)
	builder.SetEnv("GOPATH", dir)
	group := builder.Group{
		{Name: "mj_serer", Source: "src/server.go", Version: "v1.0.0", BinPath: "bin"},
		{Name: "mj_robot", Source: "src/robot.go", Version: "v1.0.0", BinPath: "bin"},
	}
	group.Build()
}
