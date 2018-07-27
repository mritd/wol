package config

import (
	"bytes"
	"text/template"

	"fmt"

	"strings"

	"regexp"

	"github.com/mritd/wol/pkg/utils"
	"github.com/mritd/wol/pkg/wol"
	"github.com/spf13/viper"
)

func Add(m wol.Machine) {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	for _, t := range machines {
		if m.Name == t.Name {
			utils.Exit("machine name is used!", 1)
		}
	}
	machines = append(machines, m)
	viper.Set("machines", machines)
	utils.CheckAndExit(viper.WriteConfig())
}

func Del(m wol.Machine) {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	for i, t := range machines {
		if m.Name == t.Name {
			machines = append(machines[:i], machines[i+1:]...)
		}
	}
	viper.Set("machines", machines)
	utils.CheckAndExit(viper.WriteConfig())
}

func List() {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	tpl := `Name            Mac
---------------------------------
{{range . }}{{ .Name | ListLayout }}{{ .Mac | ToUpper }}
{{end}}`
	t := template.New("")
	t.Funcs(map[string]interface{}{
		"ListLayout": ListLayout,
		"ToUpper":    strings.ToUpper,
	})
	t.Parse(tpl)
	var buf bytes.Buffer
	utils.CheckAndExit(t.Execute(&buf, machines))
	fmt.Println(buf.String())
}

func ListLayout(name string) string {
	if len(name) < 16 {
		return fmt.Sprintf("%-16s", name)
	} else {
		return fmt.Sprintf("%-16s", utils.ShortenString(name, 16))
	}
}

func FindMac(str string) string {
	reg := regexp.MustCompile(`^([0-9a-fA-F]{2}[:-]){5}([0-9a-fA-F]{2})$`)
	if reg.MatchString(str) {
		return str
	}
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	for _, m := range machines {
		if str == m.Name {
			return m.Mac
		}
	}
	return ""
}
