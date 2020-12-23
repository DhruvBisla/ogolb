package init

import (
	"os"
)

func init() {}

func Setup() error {
	err := os.Mkdir("content", 0755)
	if err != nil {
		return err
	}
	err = os.Mkdir("static", 0755)
	if err != nil {
		return err
	}
	err = os.Mkdir("templates", 0755)
	if err != nil {
		return err
	}
	return err

}
