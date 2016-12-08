package actions

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"

	"go.polydawn.net/reppl/model"
)

func Init(c *cli.Context) error {
	p := model.Project{}

	if _, err := os.Stat(".reppl"); os.IsNotExist(err) {
		p.Init()
		p.WriteFile(".reppl")
		fmt.Println("created new project file")
		return nil
	} else {
		return errors.New("project file aready exists")
	}
}
