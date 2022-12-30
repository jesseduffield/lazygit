package controllers

import (
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/snake"
)

type SnakeController struct {
	baseController
	*controllerCommon

	getGame func() *snake.Game
}

var _ types.IController = &SnakeController{}

func NewSnakeController(
	common *controllerCommon,
	getGame func() *snake.Game,
) *SnakeController {
	return &SnakeController{
		baseController:   baseController{},
		controllerCommon: common,
		getGame:          getGame,
	}
}

func (self *SnakeController) GetKeybindings(opts types.KeybindingsOpts) []*types.Binding {
	bindings := []*types.Binding{
		{
			Key:     opts.GetKey(opts.Config.Universal.NextItem),
			Handler: self.SetDirection(snake.Down),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.PrevItem),
			Handler: self.SetDirection(snake.Up),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.PrevBlock),
			Handler: self.SetDirection(snake.Left),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.NextBlock),
			Handler: self.SetDirection(snake.Right),
		},
		{
			Key:     opts.GetKey(opts.Config.Universal.Return),
			Handler: self.Escape,
		},
	}

	return bindings
}

func (self *SnakeController) Context() types.Context {
	return self.contexts.Snake
}

func (self *SnakeController) SetDirection(direction snake.Direction) func() error {
	return func() error {
		self.getGame().SetDirection(direction)
		return nil
	}
}

func (self *SnakeController) Escape() error {
	return self.c.PushContext(self.contexts.Submodules)
}
