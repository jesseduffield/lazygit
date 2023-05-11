package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/controllers/helpers"
)

type ControllerCommon struct {
	*helpers.HelperCommon
	IGetHelpers
}

type IGetHelpers interface {
	Helpers() *helpers.Helpers
}

func NewControllerCommon(
	c *helpers.HelperCommon,
	IGetHelpers IGetHelpers,
) *ControllerCommon {
	return &ControllerCommon{
		HelperCommon: c,
		IGetHelpers:  IGetHelpers,
	}
}
