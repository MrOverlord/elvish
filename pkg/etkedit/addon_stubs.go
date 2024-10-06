package edit

import (
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
)

func startInstant(c etk.Context) {
}

func startCommand(c etk.Context) {
}

func startLocation(c etk.Context) {
}

func startHistlist(c etk.Context) {
}

func startHistory(c etk.Context) {
}

func startLastcmd(c etk.Context) {
}

func etkBindingFromBindingMap(ed *Editor, m *bindingsMap) etk.Binding {
	return func(ev term.Event, c etk.Context, r etk.React) etk.Reaction {
		reaction := r(ev)
		if reaction == etk.Unused {
			handled := ed.callBinding(m, ev)
			if handled {
				return etk.Consumed
			}
		}
		return reaction
	}
}
