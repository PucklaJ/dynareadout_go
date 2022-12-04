package dynareadout

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBindings(t *testing.T) {
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

	assert.Equal(t, "metadata", children[0])
	for i, child := range children[1:] {
		assert.Equal(t, fmt.Sprintf("d%06d", i+1), child)
	}

	children = binFile.GetChildren("/nodout/metadata/")
	if !assert.Len(t, children, 7) {
		return
	}

	assert.Equal(t, "title", children[0])
	assert.Equal(t, "version", children[1])
	assert.Equal(t, "revision", children[2])
	assert.Equal(t, "date", children[3])
	assert.Equal(t, "legend", children[4])
	assert.Equal(t, "legend_ids", children[5])
	assert.Equal(t, "ids", children[6])

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
