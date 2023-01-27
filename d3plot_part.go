package dynareadout

/*
#include "dynareadout/src/d3plot.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

type D3plotPart struct {
	SolidIDs      []uint64
	ThickShellIDs []uint64
	BeamIDs       []uint64
	ShellIDs      []uint64
}

func (part D3plotPart) GetNodeIDs(plotFile D3plot) ([]uint64, error) {
	var cPart C.d3plot_part

	if len(part.SolidIDs) != 0 {
		cPart.solid_ids = (*C.d3_word)(unsafe.Pointer(&part.SolidIDs[0]))
		cPart.num_solids = C.size_t(len(part.SolidIDs))
	}

	if len(part.ThickShellIDs) != 0 {
		cPart.thick_shell_ids = (*C.d3_word)(unsafe.Pointer(&part.ThickShellIDs[0]))
		cPart.num_thick_shells = C.size_t(len(part.ThickShellIDs))
	}

	if len(part.BeamIDs) != 0 {
		cPart.beam_ids = (*C.d3_word)(unsafe.Pointer(&part.BeamIDs[0]))
		cPart.num_beams = C.size_t(len(part.BeamIDs))
	}

	if len(part.ShellIDs) != 0 {
		cPart.shell_ids = (*C.d3_word)(unsafe.Pointer(&part.ShellIDs[0]))
		cPart.num_shells = C.size_t(len(part.ShellIDs))
	}

	var numPartNodeIDs C.size_t
	dataC := C.d3plot_part_get_node_ids(&plotFile.handle, &cPart, &numPartNodeIDs, nil)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numPartNodeIDs)

	return data, nil
}

func (part D3plotPart) GetNodeIndices(plotFile D3plot) ([]uint64, error) {
	var cPart C.d3plot_part

	if len(part.SolidIDs) != 0 {
		cPart.solid_ids = (*C.d3_word)(unsafe.Pointer(&part.SolidIDs[0]))
		cPart.num_solids = C.size_t(len(part.SolidIDs))
	}

	if len(part.ThickShellIDs) != 0 {
		cPart.thick_shell_ids = (*C.d3_word)(unsafe.Pointer(&part.ThickShellIDs[0]))
		cPart.num_thick_shells = C.size_t(len(part.ThickShellIDs))
	}

	if len(part.BeamIDs) != 0 {
		cPart.beam_ids = (*C.d3_word)(unsafe.Pointer(&part.BeamIDs[0]))
		cPart.num_beams = C.size_t(len(part.BeamIDs))
	}

	if len(part.ShellIDs) != 0 {
		cPart.shell_ids = (*C.d3_word)(unsafe.Pointer(&part.ShellIDs[0]))
		cPart.num_shells = C.size_t(len(part.ShellIDs))
	}

	var numPartNodeIDs C.size_t
	dataC := C.d3plot_part_get_node_indices(&plotFile.handle, &cPart, &numPartNodeIDs, nil)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numPartNodeIDs)

	return data, nil
}
