# dynareadout

An Ansi C library for parsing binary output files of LS Dyna (d3plot, binout) with bindings for go.

## Examples

### Binout

```go
package main

import (
	"fmt"
	dro "github.com/PucklaJ/dynareadout_go"
	"os"
)

func main() {
	// This library also supports opening multiple binout files at once by globing them
	binFile, err := dro.BinoutOpen("simulation/binout*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open binout:", err)
		return
	}
	defer binFile.Close()

	// Print the children of the binout
	children := binFile.GetChildren("/")
	for i, child := range children {
		fmt.Printf("Child %d: %s\n", i, child)
	}

	// Read some data. The library implements read functions for multiple types
	nodeIds, err := binFile.ReadInt32("/nodout/metadata/ids")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read node ids:", err)
		return
	}

	for i, nid := range nodeIds {
		fmt.Printf("Node ID %d: %d\n", i, nid)
	}
}
```

### D3plot

```go
package main

import (
	"fmt"
	dro "github.com/PucklaJ/dynareadout_go"
	"os"
)

func main() {
	// Just give it the first d3plot file and it opens all of them
	plotFile, err := dro.D3plotOpen("simulation/d3plot")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open d3plot:", err)
		return
	}

	// Read the title
	title, _ := plotFile.ReadTitle()
	fmt.Println("Title:", title)

	// Read node ids
	nodeIds, _ := plotFile.ReadNodeIDs()
	fmt.Println("Nodes:", len(nodeIds))
	for i, nid := range nodeIds {
		fmt.Printf("Node %d: %d\n", i, nid)
	}

	// Read node coordinates of time step 10
	nodeCoords, _ := plotFile.ReadNodeCoordinates(10)
	for i, c := range nodeCoords {
		fmt.Printf("Node Coords %d: (%.2f, %.2f, %.2f)\n", i, c[0], c[1], c[2])
	}
}
```

## Other Languages

This library is also available for [C, C++](https://github.com/PucklaJ/dynareadout) and [python](https://github.com/PucklaJ/dynareadout_python).
