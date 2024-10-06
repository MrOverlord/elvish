package edit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"src.elv.sh/pkg/cli/modes"
	"src.elv.sh/pkg/cli/term"
	"src.elv.sh/pkg/etk"
	"src.elv.sh/pkg/etk/comps"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/eval/vals"
	"src.elv.sh/pkg/parse"
	"src.elv.sh/pkg/ui"
)

// TODO: Binding
func startMinibuf(c etk.Context) {
	pushAddon(c, withFinish(
		etk.WithInit(comps.TextArea, "prompt", addonPromptText(" MINIBUF ")),
		func(c etk.Context) {
			code := etk.BindState(c, "buffer", comps.TextBuffer{}).Get().Content
			src := parse.Source{Name: "[minibuf]", Code: code}
			notifyPort, cleanup := makeNotifyPort(c)
			defer cleanup()
			ports := []*eval.Port{eval.DummyInputPort, notifyPort, notifyPort}
			err := c.Frame().Evaler.Eval(src, eval.EvalCfg{Ports: ports})
			if err != nil {
				c.AddMsg(modes.ErrorText(err))
			}
		},
	), true)
}

// TODO: Is this the correct abstraction??
//
// This feels a bit like solving the same problem as (etk.Context).WithBinding,
// just from "outside" rather than "inside"?
//
// - WithBinding makes it possible to override from multiple levels higher, but
// can't compose multiple overrides
//
// - This allows composing multiple overrides, but only one level
func withFinish(f etk.Comp, finishFn func(etk.Context)) etk.Comp {
	return func(c etk.Context) (etk.View, etk.React) {
		v, r := f(c)
		return v, func(e term.Event) etk.Reaction {
			reaction := r(e)
			if reaction == etk.Finish {
				finishFn(c)
			}
			return reaction
		}
	}
}

func withAfterReact(f etk.Comp, afterFn func(etk.Context, etk.Reaction) etk.Reaction) etk.Comp {
	return func(c etk.Context) (etk.View, etk.React) {
		v, r := f(c)
		return v, func(e term.Event) etk.Reaction {
			reaction := r(e)
			return afterFn(c, reaction)
		}
	}
}

// Duplicate with pkg/edit/key_binding.go
func makeNotifyPort(c etk.Context) (*eval.Port, func()) {
	ch := make(chan any)
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		// Relay value outputs
		for v := range ch {
			notifyf(c, "[value out] %s", vals.ReprPlain(v))
		}
		wg.Done()
	}()
	go func() {
		// Relay byte outputs
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if line != "" {
					notifyf(c, "[bytes out] %s", line)
				}
				if err != io.EOF {
					notifyf(c, "[bytes error] %s", err)
				}
				break
			}
			notifyf(c, "[bytes out] %s", line[:len(line)-1])
		}
		r.Close()
		wg.Done()
	}()
	port := &eval.Port{Chan: ch, File: w}
	cleanup := func() {
		close(ch)
		w.Close()
		wg.Wait()
	}
	return port, cleanup
}

func notifyf(c etk.Context, format string, args ...any) {
	c.AddMsg(ui.T(fmt.Sprintf(format, args...)))
}
