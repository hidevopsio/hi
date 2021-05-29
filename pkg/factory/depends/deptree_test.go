// Copyright (c) 2016 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer
//    in this position and unchanged.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package depends

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"testing"
)

func TestDepTree(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	//
	// A working dependency graph
	//
	nodeA := NewNode(0, "A")
	nodeB := NewNode(1, "B")
	nodeC := NewNode(2, "C", nodeA)
	nodeD := NewNode(3, "D", nodeB)
	nodeE := NewNode(4, "E", nodeC, nodeD)
	nodeF := NewNode(5, "F", nodeA, nodeB)
	nodeG := NewNode(6, "G", nodeE, nodeF)
	nodeH := NewNode(7, "H", nodeG)
	nodeI := NewNode(8, "I", nodeA)
	nodeJ := NewNode(8, "J", nodeB)
	nodeK := NewNode(10, "K")

	var workingGraph Graph
	workingGraph = append(workingGraph, nodeA, nodeB, nodeC, nodeD, nodeE, nodeF, nodeG, nodeH, nodeI, nodeJ, nodeK)

	fmt.Printf(">>> A working dependency graph\n")
	displayDependencyGraph("workingGraph", workingGraph, log.Debug)

	resolved, err := resolveGraph(workingGraph)
	assert.Equal(t, nil, err)
	if err != nil {
		log.Errorf("Failed to resolve dependency graph: %s\n", err)
	} else {
		log.Debugf("The dependency graph resolved successfully")
	}
	displayDependencyGraph("resolved", resolved, log.Debug)

	//
	// A broken dependency graph with circular dependency
	//
	nodeA = NewNode(11, "A", nodeI)

	var brokenGraph Graph
	brokenGraph = append(brokenGraph, nodeA, nodeB, nodeC, nodeD, nodeE, nodeF, nodeG, nodeH, nodeI, nodeJ, nodeK)

	fmt.Printf(">>> A broken dependency graph with circular dependency\n")
	displayDependencyGraph("brokenGraph", brokenGraph, log.Debug)

	resolved, err = resolveGraph(brokenGraph)
	assert.Equal(t, ErrCircularDependency, err)
	if err != nil {
		log.Errorf("Failed to resolve dependency graph: %s\n", err)
	} else {
		log.Debugf("The dependency graph resolved successfully")
	}
	displayDependencyGraph("resolved", resolved, log.Debug)
}
