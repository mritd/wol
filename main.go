package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mitchellh/go-homedir"
	"github.com/mritd/logger"

	"github.com/urfave/cli/v2"
)

var (
	version   string
	buildDate string
	commitID  string
)

func main() {
	home, err := homedir.Dir()
	if err != nil {
		logger.Fatal(err)
	}
	app := &cli.App{
		Name:    "wol",
		Usage:   "Wake-on-LAN TOOL",
		Version: fmt.Sprintf("%s %s %s", version, buildDate, commitID),
		Authors: []*cli.Author{
			{
				Name:  "mritd",
				Email: "mritd@linux.com",
			},
		},
		Copyright: "Copyright (c) 2020 mritd, All rights reserved.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   filepath.Join(home, ".wol.yaml"),
				Usage:   "wol config",
				EnvVars: []string{"WOL_CONFIG"},
			},
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "machine name",
			},
			&cli.StringFlag{
				Name:    "mac",
				Aliases: []string{"m"},
				Usage:   "machine mac address",
			},
		},
		Action: func(c *cli.Context) error {
			dev := c.Args().First()
			if dev == "" {
				return cli.ShowAppHelp(c)
			}

			var cfg WolConfig
			err := cfg.LoadFrom(c.String("config"))
			if err != nil {
				return err
			}
			m := &Machine{
				Name: dev,
				Mac:  dev,
			}
			_, fm := cfg.FindMachine(m)
			if fm == nil {
				return fmt.Errorf("not found machine [%v]", m)
			}
			return fm.Wake()
		},
		Commands: []*cli.Command{
			addCmd(),
			delCmd(),
			printCmd(),
			exampleCmd(),
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logger.Error(err)
	}
}

func addCmd() *cli.Command {
	return &cli.Command{
		Name:  "add",
		Usage: "add machine",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "machine name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "mac",
				Aliases:  []string{"m"},
				Usage:    "machine mac address",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "interface",
				Aliases: []string{"i"},
				Usage:   "broadcast interface",
			},
			&cli.StringFlag{
				Name:    "addr",
				Aliases: []string{"a"},
				Usage:   "broadcast address",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "broadcast port",
			},
		},
		Action: func(c *cli.Context) error {
			// check mac address
			reg := regexp.MustCompile(`^([0-9a-fA-F]{2}[:-]){5}([0-9a-fA-F]{2})$`)
			if !reg.MatchString(c.String("mac")) {
				return fmt.Errorf("invalid mac address: %s", c.String("mac"))
			}

			var cfg WolConfig
			err := cfg.LoadFrom(c.String("config"))
			if err != nil {
				return err
			}

			m := &Machine{
				Name:               c.String("name"),
				Mac:                c.String("mac"),
				BroadcastInterface: c.String("interface"),
				BroadcastIP:        c.String("ip"),
				Port:               c.Int("port"),
			}
			return cfg.AddMachine(m)
		},
	}
}

func delCmd() *cli.Command {
	return &cli.Command{
		Name:  "del",
		Usage: "del machine",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "machine name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "mac",
				Aliases:  []string{"m"},
				Usage:    "machine mac address",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			var cfg WolConfig
			err := cfg.LoadFrom(c.String("config"))
			if err != nil {
				return err
			}
			m := &Machine{
				Name: c.String("name"),
				Mac:  c.String("mac"),
			}
			return cfg.DelMachine(m)
		},
	}
}

func printCmd() *cli.Command {
	return &cli.Command{
		Name:  "print",
		Usage: "print machines",
		Action: func(c *cli.Context) error {
			var cfg WolConfig
			err := cfg.LoadFrom(c.String("config"))
			if err != nil {
				return err
			}
			return cfg.Print()
		},
	}
}
func exampleCmd() *cli.Command {
	return &cli.Command{
		Name:  "example",
		Usage: "print example config",
		Action: func(c *cli.Context) error {
			fmt.Println(ExampleConfig())
			return nil
		},
	}
}
