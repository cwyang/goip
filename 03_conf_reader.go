// Conf reader
package main

import (
	"encoding/json"

	"github.com/kylelemons/go-gypsy/yaml"

	"gopkg.in/gcfg.v1"
	
	"fmt"
	"os"
)

type confType struct {
	Foo bool
	Bar string
}

func main() {
	file, _ := os.Open("conf.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := confType{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(conf.Bar)

	/**/

	config, err := yaml.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(config.Get("foo"))
	fmt.Println(config.Get("bar"))

	/**/

	confIni := struct {
		Section struct {
			Foo bool	// should be Capitalized
			Bar string
		}
	}{}
	err2 := gcfg.ReadFileInto(&confIni, "conf.ini")
	if err2 != nil {
		fmt.Println("Error:", err2)
	}
	fmt.Println(confIni.Section.Foo)
	fmt.Println(confIni.Section.Bar)
}

