package dynareadout

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinout(t *testing.T) {
	binFile, err := BinoutOpen("dynareadout/test_data/binout0*")
	if !assert.Nil(t, err) {
		return
	}
	defer binFile.Close()

	children := binFile.GetChildren("/")
	if !assert.Len(t, children, 2) {
		return
	}

	assert.Equal(t, children[0], "nodout")
	assert.Equal(t, children[1], "rcforc")

	children = binFile.GetChildren("/nodout")
	if !assert.Len(t, children, 602) {
		return
	}

	assert.Equal(t, "metadata", children[len(children)-1])
	for i, child := range children[:len(children)-1] {
		assert.Equal(t, fmt.Sprintf("d%06d", i+1), child)
	}

	timesteps, err := binFile.GetNumTimesteps("/nodout")
	assert.Nil(t, err)
	assert.Equal(t, uint64(601), timesteps)
	_, err = binFile.GetNumTimesteps("/nodout/schinken")
	assert.NotNil(t, err)

	children = binFile.GetChildren("/nodout/metadata/")
	if !assert.Len(t, children, 7) {
		return
	}

	assert.Equal(t, "date", children[0])
	assert.Equal(t, "ids", children[1])
	assert.Equal(t, "legend", children[2])
	assert.Equal(t, "legend_ids", children[3])
	assert.Equal(t, "revision", children[4])
	assert.Equal(t, "title", children[5])
	assert.Equal(t, "version", children[6])

	if !assert.True(t, binFile.VariableExists("/nodout/metadata/legend")) {
		return
	}
	if !assert.Equal(t, uint64(BinoutTypeInt8), binFile.GetTypeID("/nodout/metadata/legend")) {
		return
	}

	legend, err := binFile.ReadString("/nodout/metadata/legend")
	if !assert.Nil(t, err) {
		return
	}

	assert.Len(t, legend, 80)
	assert.Equal(t, "History_node_1                                                                  ", legend)

	if !assert.True(t, binFile.VariableExists("/nodout/metadata/ids")) {
		return
	}
	if !assert.Equal(t, uint64(BinoutTypeInt64), binFile.GetTypeID("/nodout/metadata/ids")) {
		return
	}

	nodeIds, err := binFile.ReadInt64("/nodout/metadata/ids")
	if !assert.Nil(t, err) {
		return
	}
	assert.Len(t, nodeIds, 1)

	if !assert.True(t, binFile.VariableExists("/rcforc/metadata/title")) {
		return
	}
	if !assert.Equal(t, uint64(BinoutTypeInt8), binFile.GetTypeID("/rcforc/metadata/title")) {
		return
	}

	title, err := binFile.ReadString("/rcforc/metadata/title")
	if !assert.Nil(t, err) {
		return
	}
	assert.Len(t, title, 80)

	assert.Equal(t, "Pouch_macro_37Ah                                                                ", title)
}

func TestD3plot(t *testing.T) {
	plotFile, err := D3plotOpen("dynareadout/test_data/d3plot")
	if !assert.Nil(t, err) {
		return
	}
	defer plotFile.Close()

	title, err := plotFile.ReadTitle()
	assert.Nil(t, err)
	assert.Equal(t, "Pouch_macro_37Ah                        ", title)

	// TODO: Read Run Time

	if !assert.Equal(t, uint64(102), plotFile.NumTimeSteps()) {
		return
	}

	nodeIds, err := plotFile.ReadNodeIDs()
	if !assert.Nil(t, err) || !assert.Len(t, nodeIds, 114893) {
		return
	}

	assert.Equal(t, uint64(84285019), nodeIds[59530])
	assert.Equal(t, uint64(10), nodeIds[0])
	assert.Equal(t, uint64(84340381), nodeIds[114892])
	assert.Equal(t, uint64(2852), nodeIds[2458])

	elementIDs, err := plotFile.ReadAllElementIDs()
	if !assert.Nil(t, err) || !assert.Len(t, elementIDs, 133456) {
		return
	}

	assert.Equal(t, uint64(1), elementIDs[0])
	assert.Equal(t, uint64(2), elementIDs[1])
	assert.Equal(t, uint64(3), elementIDs[2])
	assert.Equal(t, uint64(4), elementIDs[3])
	assert.Equal(t, uint64(72044862), elementIDs[133318])

	timeValue, err := plotFile.ReadTime(0)
	assert.Nil(t, err)
	assert.Greater(t, 1e-6, math.Abs(timeValue-0.0))
	timeValue, err = plotFile.ReadTime(10)
	assert.Nil(t, err)
	assert.Greater(t, 1e-6, math.Abs(timeValue-0.999915))
	timeValue, err = plotFile.ReadTime(19)
	assert.Nil(t, err)
	assert.Greater(t, 1e-6, math.Abs(timeValue-1.899986))
}
