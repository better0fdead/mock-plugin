package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/better0fdead/mock-plugin/parser"

	"github.com/sirupsen/logrus"
)

const Version = "0.0.1"

//go:embed about.md
var about []byte

// description is a plugin description used for 'tg plugin list' command
const description = "template plugin"

// source is a plugin source URL used for 'tg plugin update' command
const source = "github.com/better0fdead/mock-plugin"

type Description struct {
	Desc    string `json:"Desc,omitempty"`
	Version string `json:"Version,omitempty"`
}

type PluginCtx struct {
	// Plugin flags
	// Add your plugin flags here
	Version string
	// Flags for tg commands
	Help        bool
	Doc         bool
	Source      bool
	Description bool
}

type SendCtx struct {
	Pr    []byte            `json:"Pr,omitempty"`
	Flags map[string]string `json:"flags,omitempty"`
}

func setPluginCtx(flags map[string]string) PluginCtx {
	// Plugin flags.
	// Add your plugin flags here.
	pluginCtx := PluginCtx{
		Version: flags["Version"],
	}

	// Flags for tg usage.
	if _, exists := flags["help"]; exists {
		pluginCtx.Help = true
	}
	if _, exists := flags["h"]; exists {
		pluginCtx.Help = true
	}
	if _, exists := flags["doc"]; exists {
		pluginCtx.Doc = true
	}
	if _, exists := flags["source"]; exists {
		pluginCtx.Source = true
	}
	if _, exists := flags["desc"]; exists {
		pluginCtx.Description = true
	}
	return pluginCtx
}

func DesirializeData(jsonData []byte) (PluginCtx, parser.PackageInfo, error) {
	var req SendCtx
	log := logrus.WithTime(time.Now())
	err := json.Unmarshal(jsonData, &req)
	if err != nil {
		log.Infof("error deserializing sent Data: %s", err.Error())
		return PluginCtx{}, parser.PackageInfo{}, err
	}

	pluginCtx := setPluginCtx(req.Flags)

	var parsedPackage parser.PackageInfo
	if len(req.Pr) > 0 {
		err = json.Unmarshal(req.Pr, &parsedPackage)
		if err != nil {
			log.Infof("error deserializing parser data: %s", err.Error())
			return PluginCtx{}, parser.PackageInfo{}, err
		}
		for j, service := range parsedPackage.Services {
			for i, method := range service.Methods {
				var parametrs []parser.FieldPkgInfo
				parametrs = append(parametrs, parser.FieldPkgInfo{Name: "ctx", Kind: "context.Context"})
				parsedPackage.Types["context.Context"] = parser.TypeInfo{Name: "context.Context", IsScalar: true, Pkg: "context"}
				parametrs = append(parametrs, method.Parameters...)
				returns := append(method.Returns, parser.FieldPkgInfo{Name: "err", Kind: "error", IsScalar: true})
				parsedPackage.Services[j].Methods[i].Parameters = parametrs
				parsedPackage.Services[j].Methods[i].Returns = returns
			}
		}
	}

	return pluginCtx, parsedPackage, err
}

// Generate generates whatever your plugin whants to do.
// pluginCtx contains your plugin flags.
// parsedPackage contains parsed ast tree of services.
func Generate(pluginCtx PluginCtx, parsedPackage parser.PackageInfo) error {
	var err error

	return err
}

// help returns help message for your plugin.
func help() string {
	helpMsg := description + "\n"

	return helpMsg
}

func main() {
	socket, err := net.Listen("unix", "./plugin.sock")
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove("./plugin.sock")
		os.Exit(1)
	}()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()
			buf := make([]byte, 50000)

			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			pluginCtx, parsedPackage, err := DesirializeData(buf[:n])
			if err != nil {
				conn.Write([]byte(err.Error()))
				return
			}
			if pluginCtx.Help {
				conn.Write([]byte(help()))
				return
			}
			if pluginCtx.Doc {
				conn.Write(about)
				return
			}
			if pluginCtx.Source {
				conn.Write([]byte(source))
				return
			}
			if pluginCtx.Description {
				description := Description{Desc: description, Version: Version}
				mDescription, err := json.Marshal(description)
				if err != nil {
					conn.Write([]byte(err.Error()))
					return
				}
				conn.Write(mDescription)
				return
			}

			err = Generate(pluginCtx, parsedPackage)
			if err != nil {
				conn.Write([]byte(err.Error()))
			} else {
				conn.Write([]byte("done"))
			}

		}(conn)
	}
}
