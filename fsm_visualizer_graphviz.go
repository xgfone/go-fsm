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

import (
	"bytes"
	"fmt"
)

// VisualizeGraphviz outputs a visualization of a FSM in Graphviz format.
func (f *FSM) VisualizeGraphviz() string {
	transitions := cloneAndSortTransitions(f.Transitions())

	var buf bytes.Buffer
	buf.Grow(256)

	writeHeaderLine(&buf)
	writeTransitions(&buf, f.Initial(), transitions)
	writeStates(&buf, transitions)
	writeFooter(&buf)

	return buf.String()
}

func writeHeaderLine(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`digraph fsm {`))
	buf.WriteString("\n")
}

func writeTransitions(buf *bytes.Buffer, initial State, transitions []Transition) {
	// make sure the initial state is at top
	for _, t := range transitions {
		if t.Source == initial {
			fmt.Fprintf(buf, `    "%s" -> "%s" [ label = "%s" ];`+"\n", t.Source, t.Target, t.Event)
		}
	}

	for _, t := range transitions {
		if t.Source != initial {
			fmt.Fprintf(buf, `    "%s" -> "%s" [ label = "%s" ];`+"\n", t.Source, t.Target, t.Event)
		}
	}

	buf.WriteString("\n")
}

func writeStates(buf *bytes.Buffer, transitions []Transition) {
	states := getAllSortedStatesFromTransitions(transitions)
	for _, s := range states {
		buf.WriteString(fmt.Sprintf(`    "%s";`+"\n", s))
	}
}

func writeFooter(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintln("}"))
}

func hasState(ss []State, s State) bool {
	for _, _s := range ss {
		if s == _s {
			return true
		}
	}
	return false
}
