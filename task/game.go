package task

import (
	"github.com/robfig/cron"
)

func InitTask() {

}

func GameTask() {
	c := cron.New()
	spec := "*/1 * * * * *"
	c.AddFunc(spec, func() {

	})
	c.Start()
}
