package config

import (
	"bytes"
	"text/template"

	"fmt"

	"strings"

	"github.com/mritd/wol/pkg/utils"
	"github.com/mritd/wol/pkg/wol"
	"github.com/spf13/viper"
)

func Add(m wol.Machine) {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	machines = append(machines, m)
	viper.Set("machines", machines)
	utils.CheckAndExit(viper.WriteConfig())
}

func Del(m wol.Machine) {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	for i, tmpm := range machines {
		if m.Name == tmpm.Name {
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
		return fmt.Sprintf("%-16s", utils.ShortenString(name, 8))
	}
}
