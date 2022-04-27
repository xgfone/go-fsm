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

import "sync"

var pool = sync.Pool{New: func() interface{} { return New() }}

// Acquire returns a FMS from the pool, which should be put into the pool later.
func Acquire() *FSM { return pool.Get().(*FSM) }

// Release only puts the fsm into the pool.
func Release(fsm *FSM) { pool.Put(fsm) }
