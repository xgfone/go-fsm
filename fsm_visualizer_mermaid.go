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

// VisualizeMermaidStateDiagram outputs a visualization of a FSM
// in MermaidStateDiagram format.
//
// See http://mermaid-js.github.io/mermaid/#/stateDiagram
func (f *FSM) VisualizeMermaidStateDiagram() string {
	var buf bytes.Buffer
	buf.Grow(256)

	transitions := cloneAndSortTransitions(f.Transitions())

	buf.WriteString("stateDiagram-v2\n")
	buf.WriteString(fmt.Sprintln(`    [*] -->`, f.Current()))
	for _, t := range transitions {
		fmt.Fprintf(&buf, `    %s --> %s: %s`+"\n", t.Source, t.Target, t.Event)
	}

	return buf.String()
}

// VisualizeMermaidFlowChart outputs a visualization of a FSM
// in MermaidFlowChart format.
//
// See http://mermaid-js.github.io/mermaid/#/flowchart
func (f *FSM) VisualizeMermaidFlowChart(currentStateRGB string) string {
	var buf bytes.Buffer
	buf.Grow(256)

	transitions := cloneAndSortTransitions(f.Transitions())
	states := getAllSortedStatesFromTransitions(transitions)
	stateIDs := make(map[State]string, len(transitions))
	for i, state := range states {
		stateIDs[state] = fmt.Sprintf("id%d", i)
	}

	writeFlowChartGraphType(&buf)
	writeFlowChartStates(&buf, states, stateIDs)
	writeFlowChartTransitions(&buf, transitions, states, stateIDs)
	writeFlowChartHighlight(&buf, stateIDs[f.Current()], currentStateRGB)

	return buf.String()
}

func writeFlowChartGraphType(buf *bytes.Buffer) {
	buf.WriteString("graph LR\n")
}

func writeFlowChartStates(buf *bytes.Buffer, states []State, ids map[State]string) {
	for _, state := range states {
		fmt.Fprintf(buf, `    %s[%s]`+"\n", ids[state], state)
	}
	buf.WriteString("\n")
}

func writeFlowChartTransitions(buf *bytes.Buffer, transitions []Transition,
	states []State, ids map[State]string) {

	for _, t := range transitions {
		fmt.Fprintf(buf, `    %s --> |%s| %s`+"\n", ids[t.Source], t.Event, ids[t.Target])
	}
	buf.WriteString("\n")
}

func writeFlowChartHighlight(buf *bytes.Buffer, id, rgb string) {
	if id != "" && rgb != "" {
		fmt.Fprintf(buf, `    style %s fill:%s`+"\n", id, rgb)
	}
}
