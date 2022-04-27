// Copyright 2022 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fsm

import "fmt"

func ExampleFSM_SetEvent() {
	const (
		StateFoo = State("StateFoo")
		StateBar = State("StateBar")
	)

	const (
		EventFoo = Event("EventFoo")
		EventBar = Event("EventBar")
	)

	fsm := New()
	fsm.SetInitial(StateBar)

	Source(StateFoo).WithTarget(StateBar).WithEvent(EventBar).Add(fsm)
	Source(StateBar).WithTarget(StateFoo).WithEvent(EventFoo).
		WithAction(func(fsm *FSM, data interface{}) (transition bool) {
			fsm.SetEvent(EventBar, nil)
			return true
		}).
		Add(fsm)

	fmt.Println(fsm.Current())          // StateBar
	err := fsm.SendEvent(EventFoo, nil) // StateBar -> StateFoo -> StateBar
	fmt.Println(fsm.Current())          // StateBar
	fmt.Println(err)

	// Output:
	// StateBar
	// StateBar
	// <nil>
}

func ExampleFSM() {
	const (
		StateFoo = State("StateFoo")
		StateBar = State("StateBar")
	)

	const (
		EventFoo = Event("EventFoo")
		EventBar = Event("EventBar")
	)

	fsm := New()
	fsm.SetInitial(StateFoo)

	/// Add the state transitions.
	var barCount int
	Source(StateFoo).WithTarget(StateBar).WithEvent(EventBar).Add(fsm) // No Action
	Target(StateFoo).WithSource(StateBar).WithEvent(EventFoo).
		WithAction(func(fsm *FSM, data interface{}) (transition bool) { // Set Action
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
	fsm.OnEnter(func(s State) { fmt.Printf("OnEnter: %s\n", s) })
	fsm.OnExit(func(s State) { fmt.Printf("OnExit: %s\n", s) })

	fsm.OnEnterState(StateFoo, func(s State) { fmt.Printf("OnEnterState: %s\n", s) })
	fsm.OnEnterState(StateBar, func(s State) { fmt.Printf("OnEnterState: %s\n", s) })

	fsm.OnExitState(StateFoo, func(s State) { fmt.Printf("OnExitState: %s\n", s) })
	fsm.OnExitState(StateBar, func(s State) { fmt.Printf("OnExitState: %s\n", s) })

	fsm.OnTransition(func(last, current State) {
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
