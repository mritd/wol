package config

import (
	"bytes"
	"text/template"

	"fmt"

	"github.com/mritd/wol/pkg/utils"
	"github.com/mritd/wol/pkg/wol"
	"github.com/spf13/viper"
)

func Add(name, mac, broadcastInterface, broadcastIP string, port int) {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	machines = append(machines, wol.Machine{
		Name:               name,
		Mac:                mac,
		BroadcastIP:        broadcastIP,
		BroadcastInterface: broadcastInterface,
		Port:               port,
	})
	viper.Set("machines", machines)
	utils.CheckAndExit(viper.WriteConfig())
}

func Del(name string) {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	for i, m := range machines {
		if name == m.Name {
			machines = append(machines[:i], machines[i+1:]...)
		}
	}
	viper.Set("machines", machines)
	utils.CheckAndExit(viper.WriteConfig())
}

func List() {
	var machines []wol.Machine
	utils.CheckAndExit(viper.UnmarshalKey("machines", &machines))
	tpl := `  Name           Mac
              ------------------------------
{{range .machines }}{{.Name}}   {{.Mac}}{{end}}`
	t, err := template.New("").Parse(tpl)
	utils.CheckAndExit(err)
	var buf bytes.Buffer
	utils.CheckAndExit(t.Execute(&buf, machines))
	fmt.Println(buf.String())
}
