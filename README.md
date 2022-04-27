# Go FSM [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-fsm)](https://pkg.go.dev/github.com/xgfone/go-fsm) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-fsm/master/LICENSE) [![Build Status](https://github.com/xgfone/go-fsm/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-fsm/actions/workflows/go.yml)

The package `fsm` provides a simple Non-Hierarchical [Finite State Machine](https://en.wikipedia.org/wiki/Finite-state_machine) based on the event. Support Go `1.5+`.


## Install
```shell
$ go get -u github.com/xgfone/go-fsm
```


## Example
```go
package main

import (
	"fmt"

	gofsm "github.com/xgfone/go-fsm"
)

func main() {
	const (
		StateFoo = gofsm.State("StateFoo")
		StateBar = gofsm.State("StateBar")
	)

	const (
		EventFoo = gofsm.Event("EventFoo")
		EventBar = gofsm.Event("EventBar")
	)

	fsm := gofsm.New()
	fsm.SetInitial(StateFoo)

	/// Add the state transitions.
	var barCount int
	gofsm.Source(StateFoo).WithTarget(StateBar).WithEvent(EventBar).Add(fsm) // No Action
	gofsm.Target(StateFoo).WithSource(StateBar).WithEvent(EventFoo).
		WithAction(func(fsm *gofsm.FSM, data interface{}) (transition bool) { // Set Action
			// TODO: do business something

			// Here as the example, after trigger the event EventFoo two twice,
			// transition the state to the target.
			if barCount > 0 {
				barCount = 0
				transition = true
			} else {
				barCount++
			}
			return
		}).
		Add(fsm)

	/// Set the listener of the state change.
	fsm.OnEnter(func(s gofsm.State) { fmt.Printf("OnEnter: %s\n", s) })
	fsm.OnExit(func(s gofsm.State) { fmt.Printf("OnExit: %s\n", s) })

	fsm.OnEnterState(StateFoo, func(s gofsm.State) { fmt.Printf("OnEnterState: %s\n", s) })
	fsm.OnEnterState(StateBar, func(s gofsm.State) { fmt.Printf("OnEnterState: %s\n", s) })

	fsm.OnExitState(StateFoo, func(s gofsm.State) { fmt.Printf("OnExitState: %s\n", s) })
	fsm.OnExitState(StateBar, func(s gofsm.State) { fmt.Printf("OnExitState: %s\n", s) })

	fsm.OnTransition(func(last, current gofsm.State) {
		fmt.Printf("OnTransition: %s -> %s\n", last, current)
	})

	/// Send the events to the state machine
	fmt.Println("------ Transition ------")
	last := fsm.Current()
	err := fsm.SendEvent(EventBar, nil)
	fmt.Printf("Event: %s, State: %s -> %s, Result: %v\n\n", EventBar, last, fsm.Current(), err)

	fmt.Println("------ Transition ------")
	last = fsm.Current()
	err = fsm.SendEvent(EventFoo, nil)
	fmt.Printf("Event: %s, State: %s -> %s, Result: %v\n\n", EventFoo, last, fsm.Current(), err)

	fmt.Println("------ Transition ------")
	last = fsm.Current()
	err = fsm.SendEvent(EventFoo, nil)
	fmt.Printf("Event: %s, State: %s -> %s, Result: %v\n\n", EventFoo, last, fsm.Current(), err)

	fmt.Println("------ Transition ------")
	last = fsm.Current()
	err = fsm.SendEvent(EventFoo, nil)
	fmt.Printf("Event: %s, State: %s -> %s, Result: %v\n\n", EventFoo, last, fsm.Current(), err)

	fmt.Println("------ Transition ------")
	last = fsm.Current()
	err = fsm.SendEvent(EventBar, nil)
	fmt.Printf("Event: %s, State: %s -> %s, Result: %v\n\n", EventBar, last, fsm.Current(), err)

	// Print the Graphviz visualizer.
	fmt.Println("------ Graphviz ------")
	fmt.Println(fsm.VisualizeGraphviz())

	// Print the Mermaid FlowChart visualizer.
	fmt.Println("------ Mermaid FlowChart ------")
	fmt.Println(fsm.VisualizeMermaidFlowChart("#aaaaaa", "#ff0000"))

	// Print the Mermaid StateDiagram visualizer.
	fmt.Println("------ Mermaid StateDiagram ------")
	fmt.Println(fsm.VisualizeMermaidStateDiagram())

	// Output:
	// ------ Transition ------
	// OnExitState: StateFoo
	// OnExit: StateFoo
	// OnEnterState: StateBar
	// OnEnter: StateBar
	// OnTransition: StateFoo -> StateBar
	// Event: EventBar, State: StateFoo -> StateBar, Result: <nil>
	//
	// ------ Transition ------
	// Event: EventFoo, State: StateBar -> StateBar, Result: source state 'StateBar' transition for the event 'EventFoo' is suspended
	//
	// ------ Transition ------
	// OnExitState: StateBar
	// OnExit: StateBar
	// OnEnterState: StateFoo
	// OnEnter: StateFoo
	// OnTransition: StateBar -> StateFoo
	// Event: EventFoo, State: StateBar -> StateFoo, Result: <nil>
	//
	// ------ Transition ------
	// Event: EventFoo, State: StateFoo -> StateFoo, Result: no transition for the event 'EventFoo'
	//
	// ------ Transition ------
	// OnExitState: StateFoo
	// OnExit: StateFoo
	// OnEnterState: StateBar
	// OnEnter: StateBar
	// OnTransition: StateFoo -> StateBar
	// Event: EventBar, State: StateFoo -> StateBar, Result: <nil>
	//
	// ------ Graphviz ------
	// digraph fsm {
	//     "StateFoo" -> "StateBar" [ label = "EventBar" ];
	//     "StateBar" -> "StateFoo" [ label = "EventFoo" ];
	//
	//     "StateBar";
	//     "StateFoo";
	// }
	//
	// ------ Mermaid FlowChart ------
	// graph LR
	//     id0[StateBar]
	//     id1[StateFoo]
	//
	//     id0 --> |EventFoo| id1
	//     id1 --> |EventBar| id0
	//
	//     style id1 fill:#aaaaaa
	//     style id0 fill:#ff0000
	//
	// ------ Mermaid StateDiagram ------
	// stateDiagram-v2
	//     [*] --> StateFoo
	//     StateBar --> StateFoo: EventFoo
	//     StateFoo --> StateBar: EventBar
	//
}
```
