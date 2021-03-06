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

// Package fsm provides a simple Non-Hierarchical Finite State Machine
// based on the event.
package fsm

import (
	"fmt"
	"sort"
)

// Event is the event type.
type Event string

// State is the state type.
type State string

// Action is a function that is called when the state is transitioned.
type Action func(fsm *FSM, data interface{}) (transition bool)

// TransitionError is an transition error.
type TransitionError struct {
	Event Event

	// If empty, the error represents no transition to support the event.
	// Or, represents the state transition is suspended by Action.
	Source State
	Target State
}

// IsSuspended reports whether the error is that the state transition
// is suspended by Action.
func IsSuspended(err error) bool {
	if te, ok := err.(TransitionError); ok {
		return te.IsSuspended()
	}
	return false
}

// IsNoTransition reports whether the error is that there is no state transition
// to support the event.
func IsNoTransition(err error) bool {
	if te, ok := err.(TransitionError); ok {
		return te.IsNoTransition()
	}
	return false
}

// IsSuspended reports whether the state transition is suspended by Action.
func (e TransitionError) IsSuspended() bool { return len(e.Source) > 0 }

// IsNoTransition reports whether there is no state transition to support the event.
func (e TransitionError) IsNoTransition() bool { return len(e.Source) > 0 }

func (e TransitionError) Error() string {
	if e.Source == "" {
		return fmt.Sprintf("no transition for the event '%s'", e.Event)
	}

	const s = "source state '%s' transition for the event '%s' is suspended"
	return fmt.Sprintf(s, e.Source, e.Event)
}

// Transition represents the state transition based on the input event
// from source to target.
type Transition struct {
	Event  Event
	Source State
	Target State

	// If Action is nil, transition the state from source to target directly.
	// Or, call it before transitioning the state and transition the state
	// from source to target only if returning true.
	Action Action
}

// NewTransition returns a Transition.
func NewTransition(source, target State, event Event, action Action) Transition {
	return Transition{Event: event, Source: source, Target: target, Action: action}
}

// Source returns a new Transition with the source state.
func Source(source State) Transition { return Transition{Source: source} }

// Target returns a new Transition with the target state.
func Target(target State) Transition { return Transition{Target: target} }

// WithSource returns a new Transition with the source state.
func (t Transition) WithSource(source State) Transition {
	t.Source = source
	return t
}

// WithTarget returns a new Transition with the target state.
func (t Transition) WithTarget(target State) Transition {
	t.Target = target
	return t
}

// WithEvent returns a new Transition with the event.
func (t Transition) WithEvent(event Event) Transition {
	t.Event = event
	return t
}

// WithAction returns a new Transition with the action state.
func (t Transition) WithAction(action Action) Transition {
	t.Action = action
	return t
}

// Add is a handy proxy method to add the current transition into the given FSM.
func (t Transition) Add(fsm *FSM) { fsm.AddTransitions(t) }

// FSM is a finite state machine without thread-safe.
type FSM struct {
	exit        func(State)
	enter       func(State)
	transition  func(last, current State)
	exitStates  map[State]func(State)
	enterStates map[State]func(State)
	transitions []Transition

	current State
	event   Event
	data    interface{}
}

// New creates a new finite state machine.
func New() *FSM {
	return &FSM{
		enterStates: make(map[State]func(State), 16),
		exitStates:  make(map[State]func(State), 16),
	}
}

// Reset resets the machine to the initial state.
func (f *FSM) Reset() {
	for key := range f.exitStates {
		delete(f.exitStates, key)
	}
	for key := range f.enterStates {
		delete(f.enterStates, key)
	}

	*f = FSM{exitStates: f.exitStates, enterStates: f.enterStates}
}

// SetCurrent resets the current state to current.
func (f *FSM) SetCurrent(current State) {
	if current == "" {
		panic("the current state must not be empty")
	}
	f.current = current
}

// Current returns the current state.
func (f *FSM) Current() State { return f.current }

// States returns all the states.
func (f *FSM) States() (states []State) {
	transitions := f.Transitions()
	states = make([]State, 0, len(transitions))
	for _, t := range transitions {
		if !hasState(states, t.Source) {
			states = append(states, t.Source)
		}
		if !hasState(states, t.Target) {
			states = append(states, t.Target)
		}
	}
	return
}

// Events returns all the events.
func (f *FSM) Events() (events []Event) {
	transitions := f.Transitions()
	events = make([]Event, 0, len(transitions))
	for _, t := range transitions {
		if !hasEvent(events, t.Event) {
			events = append(events, t.Event)
		}
	}
	return
}

// Terminations returns all the termination states.
func (f *FSM) Terminations() (states []State) {
	var sources, targets []State
	for _, t := range f.Transitions() {
		if !hasState(sources, t.Source) {
			sources = append(sources, t.Source)
		}
		if !hasState(targets, t.Target) {
			targets = append(targets, t.Target)
		}
	}

	for _, state := range targets {
		if !hasState(sources, state) {
			states = append(states, state)
		}
	}

	return
}

// Transitions returns all the transitions.
func (f *FSM) Transitions() []Transition { return f.transitions }

// GetTransition returns the transition by the source state and the input event.
func (f *FSM) GetTransition(source State, event Event) (transition Transition, ok bool) {
	if index := f.indexTransition(source, event); index > -1 {
		transition, ok = f.transitions[index], true
	}
	return
}

// AddTransitions appends a set of transitions to transfer the state.
//
// Notice: the current implementation requires that the source, target
// and event must be set.
func (f *FSM) AddTransitions(transitions ...Transition) {
	for _, t := range transitions {
		if t.Source == "" || t.Target == "" || t.Event == "" {
			panic("invalid state transition: source, target, or event is empty")
		}
	}

	for _, t := range transitions {
		if index := f.indexTransition(t.Source, t.Event); index > -1 {
			f.transitions[index] = t
		} else {
			f.transitions = append(f.transitions, t)
		}
	}
}

func (f *FSM) indexTransition(source State, event Event) (index int) {
	for i, t := range f.transitions {
		if t.Source == source && t.Event == event {
			return i
		}
	}
	return -1
}

// OnEnter sets a function that will be called when entering any state.
func (f *FSM) OnEnter(fn func(State)) { f.enter = fn }

// OnExit sets a function that will be called when exiting any state.
func (f *FSM) OnExit(fn func(State)) { f.exit = fn }

// OnEnterState sets a function that will be called when entering a specific state.
func (f *FSM) OnEnterState(state State, fn func(State)) { f.enterStates[state] = fn }

// OnExitState sets a function that will be called when exiting a specific state.
func (f *FSM) OnExitState(state State, fn func(State)) { f.exitStates[state] = fn }

// OnTransition sets a function that will be called
// when the state is transferred from last to current.
func (f *FSM) OnTransition(fn func(last, current State)) { f.transition = fn }

// TestEvent reports whether the event can trigger the state transition.
func (f *FSM) TestEvent(event Event) bool {
	return f.indexTransition(f.Current(), event) > -1
}

// SetEvent sets the event with the data as the new input to continue
// to transition the state after finishing to transition the last state,
// which is used in the transition action because SendEvent cannot be used.
func (f *FSM) SetEvent(event Event, data interface{}) {
	f.event, f.data = event, data
}

// SendEvent sends an Event to the state machine, applying at most one transition.
func (f *FSM) SendEvent(event Event, data interface{}) (err error) {
	if event == "" {
		panic("FSM: the event must not be empty")
	}

	for {
		f.SetEvent("", nil)
		err = f.sendEvent(event, data)
		if f.event == "" || (err != nil && !IsSuspended(err)) {
			break
		}
		event, data = f.event, f.data
	}

	return
}

func (f *FSM) sendEvent(event Event, data interface{}) error {
	current := f.Current()
	for _, t := range f.Transitions() {
		if t.Source == current && t.Event == event {
			if t.Action != nil && !t.Action(f, data) {
				// Transition is suspended.
				return TransitionError{Event: event, Source: t.Source, Target: t.Target}
			}

			if fn, ok := f.exitStates[current]; ok {
				fn(current)
			}
			if f.exit != nil {
				f.exit(current)
			}

			f.SetCurrent(t.Target)

			if fn, ok := f.enterStates[t.Target]; ok {
				fn(t.Target)
			}
			if f.enter != nil {
				f.enter(t.Target)
			}

			if f.transition != nil {
				f.transition(current, t.Target)
			}

			return nil
		}
	}

	return TransitionError{Event: event} // No Transition
}

type sortedTransitions []Transition

func (ts sortedTransitions) Len() int      { return len(ts) }
func (ts sortedTransitions) Swap(i, j int) { ts[i], ts[j] = ts[j], ts[i] }
func (ts sortedTransitions) Less(i, j int) bool {
	if ts[i].Source == ts[j].Source {
		return ts[i].Event < ts[j].Event
	}
	return ts[i].Source < ts[j].Source
}

type sortedStates []State

func (ss sortedStates) Len() int           { return len(ss) }
func (ss sortedStates) Swap(i, j int)      { ss[i], ss[j] = ss[j], ss[i] }
func (ss sortedStates) Less(i, j int) bool { return ss[i] < ss[j] }

func cloneAndSortTransitions(ts []Transition) []Transition {
	transitions := make(sortedTransitions, len(ts))
	copy(transitions, ts)
	sort.Sort(transitions)
	return transitions
}

func getAllSortedStatesFromTransitions(transitions []Transition) []State {
	states := make(sortedStates, 0, len(transitions))
	for _, t := range transitions {
		if !hasState(states, t.Source) {
			states = append(states, t.Source)
		}
		if !hasState(states, t.Target) {
			states = append(states, t.Target)
		}
	}

	sort.Stable(states)
	return states
}

func hasState(ss []State, s State) bool {
	for _, _s := range ss {
		if s == _s {
			return true
		}
	}
	return false
}

func hasEvent(es []Event, e Event) bool {
	for _, _e := range es {
		if e == _e {
			return true
		}
	}
	return false
}
