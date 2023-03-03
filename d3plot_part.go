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
	handle C.d3plot_part
}

func (part D3plotPart) SolidID(index int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(part.handle.solid_ids)) + uintptr(index)*unsafe.Sizeof(*part.handle.solid_ids)))
}

func (part D3plotPart) BeamID(index int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(part.handle.beam_ids)) + uintptr(index)*unsafe.Sizeof(*part.handle.beam_ids)))
}

func (part D3plotPart) ShellID(index int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(part.handle.shell_ids)) + uintptr(index)*unsafe.Sizeof(*part.handle.shell_ids)))
}

func (part D3plotPart) ThickShellID(index int) uint64 {
	return *(*uint64)(unsafe.Pointer(uintptr(unsafe.Pointer(part.handle.thick_shell_ids)) + uintptr(index)*unsafe.Sizeof(*part.handle.thick_shell_ids)))
}

func (part D3plotPart) LenSolidIDs() int {
	return int(part.handle.num_solids)
}

func (part D3plotPart) LenBeamIDs() int {
	return int(part.handle.num_beams)
}

func (part D3plotPart) LenShellIDs() int {
	return int(part.handle.num_shells)
}

func (part D3plotPart) LenThickShellIDs() int {
	return int(part.handle.num_thick_shells)
}

func (part D3plotPart) Free() {
	C.d3plot_free_part(&part.handle)
}

func (part D3plotPart) GetNodeIDs(plotFile D3plot) ([]uint64, error) {
	var numPartNodeIDs C.size_t
	dataC := C.d3plot_part_get_node_ids2(&plotFile.handle, &part.handle, &numPartNodeIDs, nil, 0, nil, 0, nil, 0, nil, 0, nil, 0, nil, nil, nil, nil)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numPartNodeIDs)

	return data, nil
}

func (part D3plotPart) GetNodeIndices(plotFile D3plot) ([]uint64, error) {
	var numPartNodeIDs C.size_t
	dataC := C.d3plot_part_get_node_indices2(&plotFile.handle, &part.handle, &numPartNodeIDs, nil, 0, nil, 0, nil, 0, nil, 0, nil, nil, nil, nil)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numPartNodeIDs)

	return data, nil
}
