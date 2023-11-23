package dynareadout

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinout(t *testing.T) {
	binFile, err := BinoutOpen("test_data/binout0*")
	if !assert.Nil(t, err) {
		return
	}
	defer binFile.Close()

	children := binFile.GetChildren("/")
	if !assert.Len(t, children, 5) {
		return
	}

	assert.Equal(t, "bndout", children[0])
	assert.Equal(t, "glstat", children[1])
	assert.Equal(t, "nodout", children[2])
	assert.Equal(t, "rcforc", children[3])
	assert.Equal(t, "sleout", children[4])

	children = binFile.GetChildren("/nodout")
	if !assert.Len(t, children, 14999) {
		return
	}

	assert.Equal(t, "metadata", children[len(children)-1])
	for i, child := range children[:len(children)-1] {
		assert.Equal(t, fmt.Sprintf("d%06d", i+1), child)
	}

	timesteps, err := binFile.GetNumTimesteps("/nodout")
	assert.Nil(t, err)
	assert.Equal(t, uint64(14998), timesteps)
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
	assert.Equal(t, "                                                                                ", legend)

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

	assert.Equal(t, "LS-DYNA keyword deck by LS-PrePost                                              ", title)

	realPath, typeID, timed, err := binFile.SimplePathToReal("nodout/x_displacement")
	assert.Nil(t, err)
	assert.Equal(t, "/nodout/x_displacement", realPath)
	assert.Equal(t, BinoutTypeFloat32, typeID)
	assert.True(t, timed)

	yDisp, err := binFile.ReadTimedFloat32("/nodout/y_displacement")
	assert.Nil(t, err)
	assert.Len(t, yDisp, 14998)
	assert.Len(t, yDisp[0], 1)
}

func TestD3plot(t *testing.T) {
	plotFile, err := D3plotOpen("test_data/d3plot_files/d3plot")
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

	part, err := plotFile.ReadPart(1)
	assert.Nil(t, err)
	assert.Equal(t, 10, part.LenShellIDs())
	part.Free()

	part, err = plotFile.ReadPartByID(71000063, nil)
	assert.Nil(t, err)

	assert.Equal(t, 7368, part.LenShellIDs())
	partNodeIDs, err := part.GetNodeIDs(plotFile)
	assert.Nil(t, err)
	assert.Len(t, partNodeIDs, 7370)

	partNodeIndices, err := part.GetNodeIndices(plotFile)
	assert.Nil(t, err)
	assert.Len(t, partNodeIndices, len(partNodeIDs))
	part.Free()

	for i, ind := range partNodeIndices {
		assert.True(t, partNodeIDs[i] == nodeIds[ind], i)
	}
}

func TestKeyFile(t *testing.T) {
	keywords, warn, err := KeyFileParse("test_data/key_file.k", DefaultKeyFileParseConfig())
	assert.Nil(t, warn)
	if !assert.Nil(t, err) {
		return
	}
	defer keywords.Free()

	parts, err := keywords.GetSlice("PART")
	if !assert.Nil(t, err) {
		return
	}

	if !assert.Len(t, parts, 2) {
		return
	}

	title := parts[0].GetSlice()[0]
	assert.Equal(t, "Cube", title.ParseWhole())

	title = parts[1].GetSlice()[0]
	assert.Equal(t, "Ground", title.ParseWhole())

	card := parts[0].Get(1)

	card.ParseBegin(DefaultValueWidth)
	v, err := card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 71000063, v)
	card.ParseNext()
	v, err = card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 71000063, v)
	card.ParseNext()
	v, err = card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 6, v)
	card.ParseNext()
	v, err = card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 0, v)
	card.ParseNext()
	v, err = card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 0, v)
	card.ParseNext()
	v, err = card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 0, v)
	card.ParseNext()
	v, err = card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 0, v)
	card.ParseNext()
	_, err = card.ParseInt()
	assert.NotNil(t, err)
	card.ParseNext()
	assert.True(t, card.ParseDone())

	keyword, err := keywords.Get("SET_NODE_LIST_TITLE", 0)
	if !assert.Nil(t, err) {
		return
	}

	card = keyword.Get(1)

	card.ParseBegin(DefaultValueWidth)
	assert.Equal(t, CardParseInt, card.ParseGetType())
	v, err = card.ParseInt()
	assert.Nil(t, err)
	assert.Equal(t, 1, v)
	card.ParseNext()
	assert.Equal(t, CardParseFloat, card.ParseGetType())
	vf, err := card.ParseFloat64()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, vf)
	card.ParseNext()
	assert.Equal(t, CardParseFloat, card.ParseGetType())
	vf, err = card.ParseFloat64()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, vf)
	card.ParseNext()
	assert.Equal(t, CardParseFloat, card.ParseGetType())
	vf, err = card.ParseFloat64()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, vf)
	card.ParseNext()
	assert.Equal(t, CardParseFloat, card.ParseGetType())
	vf, err = card.ParseFloat64()
	assert.Nil(t, err)
	assert.Equal(t, 0.0, vf)
	card.ParseNext()
	assert.Equal(t, CardParseString, card.ParseGetType())
	vStr := card.ParseString()
	assert.Equal(t, "MECH", vStr)
	card.ParseNext()
	assert.True(t, card.ParseDone())
}

func TestKeyFileParseWithCallback(t *testing.T) {
	warn, err := KeyFileParseWithCallback("test_data/key_file.k",
		func(fileName string, lineNumber int, keywordName string, card *Card, cardIndex int) {
			var cardString string
			if card != nil {
				cardString = fmt.Sprint("Card: ", card.ParseWholeNoTrim(), " CardIndex: ", cardIndex)
			}

			fmt.Println("Filename:", fileName, "Line:", lineNumber, "Keyword:", keywordName, cardString)
		}, DefaultKeyFileParseConfig())
	assert.Nil(t, warn)
	if !assert.Nil(t, err) {
		return
	}
}
