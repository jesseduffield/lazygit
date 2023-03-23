package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
)

type controllerCommon struct {
	c       *helpers.HelperCommon
	helpers *helpers.Helpers
}

func NewControllerCommon(
	c *helpers.HelperCommon,
	helpers *helpers.Helpers,
) *controllerCommon {
	return &controllerCommon{
		c:       c,
		helpers: helpers,
	}
}
