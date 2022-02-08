package timer

import (
	"cathub.me/go-gin-examples/domain/user"
	"context"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func Start() {
	c := cron.New(cron.WithSeconds())

	// Execute every 30 minutes
	wrap(c, "0 */30 * * * ?", func() {
		user.GetUserService().ClearUnActivationUser(context.Background())
	})

	c.Start()
}

func wrap(c *cron.Cron, spec string, cmd func()) {
	_, err := c.AddFunc(spec, cmd)
	if err != nil {
		log.Err(err).Msg("Add clearUnActivationUser() timer task failed")
	}
}
