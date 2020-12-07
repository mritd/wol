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
		Copyright:            "Copyright (c) 2020 mritd, All rights reserved.",
		EnableBashCompletion: true,
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
				Usage:   "device name",
			},
			&cli.StringFlag{
				Name:    "mac",
				Aliases: []string{"m"},
				Usage:   "device mac address",
			},
		},
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			addCmd(),
			delCmd(),
			wakeCmd(),
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
		Usage: "add device",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "device name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "mac",
				Aliases:  []string{"m"},
				Usage:    "device mac address",
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

			d := &Device{
				Name:               c.String("name"),
				Mac:                c.String("mac"),
				BroadcastInterface: c.String("interface"),
				BroadcastIP:        c.String("ip"),
				Port:               c.Int("port"),
			}
			return cfg.AddDevice(d)
		},
	}
}

func delCmd() *cli.Command {
	return &cli.Command{
		Name:         "del",
		Usage:        "del device",
		BashComplete: bashComplete,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}
			dev := c.Args().First()

			var cfg WolConfig
			err := cfg.LoadFrom(c.String("config"))
			if err != nil {
				return err
			}

			d := &Device{
				Name: dev,
				Mac:  dev,
			}
			return cfg.DelDevice(d)
		},
	}
}

func printCmd() *cli.Command {
	return &cli.Command{
		Name:  "print",
		Usage: "print devices",
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

func wakeCmd() *cli.Command {
	return &cli.Command{
		Name:         "wake",
		Usage:        "wake device",
		BashComplete: bashComplete,
		Action: func(c *cli.Context) error {
			if c.NArg() != 1 {
				return cli.ShowAppHelp(c)
			}
			dev := c.Args().First()

			var cfg WolConfig
			err := cfg.LoadFrom(c.String("config"))
			if err != nil {
				return err
			}

			_, device := cfg.FindDevice(&Device{
				Name: dev,
				Mac:  dev,
			})
			if device == nil {
				return fmt.Errorf("not found device [%s]", dev)
			}
			return device.Wake()
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

func bashComplete(c *cli.Context) {
	if c.NArg() > 0 {
		return
	}

	var cfg WolConfig
	err := cfg.LoadFrom(c.String("config"))
	if err != nil {
		logger.Error(err)
		return
	}
	for _, dev := range cfg.Devices {
		fmt.Println(dev.Name)
	}
}
