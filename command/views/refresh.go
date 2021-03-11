package views

import (
	"fmt"

	"github.com/hashicorp/terraform/command/arguments"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/tfdiags"
)

// The Refresh view is used for the refresh command.
type Refresh interface {
	Outputs(outputValues map[string]*states.OutputValue)

	Operation() Operation
	Hooks() []terraform.Hook

	Diagnostics(diags tfdiags.Diagnostics)
	HelpPrompt()
}

// NewRefresh returns an initialized Refresh implementation for the given ViewType.
func NewRefresh(vt arguments.ViewType, runningInAutomation bool, view *View) Refresh {
	switch vt {
	case arguments.ViewHuman:
		return &RefreshHuman{
			view:         view,
			inAutomation: runningInAutomation,
			countHook:    &countHook{},
		}
	default:
		panic(fmt.Sprintf("unknown view type %v", vt))
	}
}

// The RefreshHuman implementation renders human-readable text logs, suitable for
// a scrolling terminal.
type RefreshHuman struct {
	view *View

	inAutomation bool

	countHook *countHook
}

var _ Refresh = (*RefreshHuman)(nil)

func (v *RefreshHuman) Outputs(outputValues map[string]*states.OutputValue) {
	if len(outputValues) > 0 {
		v.view.streams.Print(v.view.colorize.Color("[reset][bold][green]\nOutputs:\n\n"))
		NewOutput(arguments.ViewHuman, v.view).Output("", outputValues)
	}
}

func (v *RefreshHuman) Operation() Operation {
	return NewOperation(arguments.ViewHuman, v.inAutomation, v.view)
}

func (v *RefreshHuman) Hooks() []terraform.Hook {
	return []terraform.Hook{
		v.countHook,
		NewUiHook(v.view),
	}
}

func (v *RefreshHuman) Diagnostics(diags tfdiags.Diagnostics) {
	v.view.Diagnostics(diags)
}

func (v *RefreshHuman) HelpPrompt() {
	v.view.HelpPrompt("refresh")
}