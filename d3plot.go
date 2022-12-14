package dynareadout

/*
#include <stdlib.h>
#include "dynareadout/src/d3plot.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

type D3plot struct {
	handle C.d3plot_file
}

type D3plotPart struct {
	SolidIDs      []uint64
	ThickShellIDs []uint64
	BeamIDs       []uint64
	ShellIDs      []uint64
}

func D3plotOpen(fileName string) (plotFile D3plot, err error) {
	fileNameC := C.CString(fileName)

	plotFile.handle = C.d3plot_open(fileNameC)
	C.free(unsafe.Pointer(fileNameC))

	if plotFile.handle.error_string != nil {
		err = errors.New(C.GoString(plotFile.handle.error_string))
		C.d3plot_close(&plotFile.handle)
	}

	return
}

func (plotFile D3plot) Close() {
	C.d3plot_close(&plotFile.handle)
}

func (plotFile D3plot) ReadNodeIDs() ([]uint64, error) {
	var numIds C.size_t
	dataC := C.d3plot_read_node_ids(&plotFile.handle, &numIds)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numIds == 0 {
		return []uint64{}, nil
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numIds)

	return data, nil
}

func (plotFile D3plot) ReadSolidElementIDs() ([]uint64, error) {
	var numIds C.size_t
	dataC := C.d3plot_read_solid_element_ids(&plotFile.handle, &numIds)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numIds == 0 {
		return []uint64{}, nil
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numIds)

	return data, nil
}

func (plotFile D3plot) ReadBeamElementIDs() ([]uint64, error) {
	var numIds C.size_t
	dataC := C.d3plot_read_beam_element_ids(&plotFile.handle, &numIds)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numIds == 0 {
		return []uint64{}, nil
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numIds)

	return data, nil
}

func (plotFile D3plot) ReadShellElementIDs() ([]uint64, error) {
	var numIds C.size_t
	dataC := C.d3plot_read_shell_element_ids(&plotFile.handle, &numIds)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numIds == 0 {
		return []uint64{}, nil
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numIds)

	return data, nil
}

func (plotFile D3plot) ReadThickShellElementIDs() ([]uint64, error) {
	var numIds C.size_t
	dataC := C.d3plot_read_thick_shell_element_ids(&plotFile.handle, &numIds)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numIds == 0 {
		return []uint64{}, nil
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numIds)

	return data, nil
}

func (plotFile D3plot) ReadAllElementIDs() ([]uint64, error) {
	var numIds C.size_t
	dataC := C.d3plot_read_all_element_ids(&plotFile.handle, &numIds)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numIds == 0 {
		return []uint64{}, nil
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numIds)

	return data, nil
}

func (plotFile D3plot) ReadPartIDs() ([]uint64, error) {
	var numIds C.size_t
	dataC := C.d3plot_read_part_ids(&plotFile.handle, &numIds)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numIds == 0 {
		return []uint64{}, nil
	}

	data := carrToSlice[C.d3_word, uint64](dataC, numIds)

	return data, nil
}

func (plotFile D3plot) ReadPartTitles() ([]string, error) {
	var numTitles C.size_t
	dataC := C.d3plot_read_part_titles(&plotFile.handle, &numTitles)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numTitles == 0 {
		return []string{}, nil
	}

	titles := make([]string, numTitles)
	for i := range titles {
		titleC := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
		titles[i] = C.GoString(titleC)
		C.free(unsafe.Pointer(titleC))
	}
	C.free(unsafe.Pointer(dataC))

	return titles, nil
}

func (plotFile D3plot) ReadNodeCoordinates(state uint64) ([][3]float64, error) {
	var numNodes C.size_t
	dataC := C.d3plot_read_node_coordinates(&plotFile.handle, C.size_t(state), &numNodes)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numNodes == 0 {
		return [][3]float64{}, nil
	}

	coords := make([][3]float64, numNodes)
	for i := range coords {
		nodePtr := carrIdxPtr(dataC, i*3)

		coords[i][0] = float64(*nodePtr)
		coords[i][1] = float64(carrIdx(nodePtr, 1))
		coords[i][2] = float64(carrIdx(nodePtr, 2))
	}
	C.free(unsafe.Pointer(dataC))

	return coords, nil
}

func (plotFile D3plot) ReadNodeVelocity(state uint64) ([][3]float64, error) {
	var numNodes C.size_t
	dataC := C.d3plot_read_node_velocity(&plotFile.handle, C.size_t(state), &numNodes)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numNodes == 0 {
		return [][3]float64{}, nil
	}

	vel := make([][3]float64, numNodes)
	for i := range vel {
		nodePtr := carrIdxPtr(dataC, i*3)

		vel[i][0] = float64(*nodePtr)
		vel[i][1] = float64(carrIdx(nodePtr, 1))
		vel[i][2] = float64(carrIdx(nodePtr, 2))
	}
	C.free(unsafe.Pointer(dataC))

	return vel, nil
}

func (plotFile D3plot) ReadNodeAcceleration(state uint64) ([][3]float64, error) {
	var numNodes C.size_t
	dataC := C.d3plot_read_node_acceleration(&plotFile.handle, C.size_t(state), &numNodes)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numNodes == 0 {
		return [][3]float64{}, nil
	}

	acc := make([][3]float64, numNodes)
	for i := range acc {
		nodePtr := carrIdxPtr(dataC, i*3)

		acc[i][0] = float64(*nodePtr)
		acc[i][1] = float64(carrIdx(nodePtr, 1))
		acc[i][2] = float64(carrIdx(nodePtr, 2))
	}
	C.free(unsafe.Pointer(dataC))

	return acc, nil
}

func (plotFile D3plot) ReadTime(state uint64) (float64, error) {
	timeC := C.d3plot_read_time(&plotFile.handle, C.size_t(state))
	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return float64(timeC), err
	}

	return float64(timeC), nil
}

func (plotFile D3plot) ReadSolidsState(state uint64) ([]C.d3plot_solid, error) {
	var numSolids C.size_t
	dataC := C.d3plot_read_solids_state(&plotFile.handle, C.size_t(state), &numSolids)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numSolids == 0 {
		return []C.d3plot_solid{}, nil
	}

	solids := make([]C.d3plot_solid, numSolids)
	for i := range solids {
		solids[i] = *(*C.d3plot_solid)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return solids, nil
}

func (plotFile D3plot) ReadThickShellsState(state uint64) ([]C.d3plot_thick_shell, error) {
	var numThickShells C.size_t
	dataC := C.d3plot_read_thick_shells_state(&plotFile.handle, C.size_t(state), &numThickShells)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numThickShells == 0 {
		return []C.d3plot_thick_shell{}, nil
	}

	thickShells := make([]C.d3plot_thick_shell, numThickShells)
	for i := range thickShells {
		thickShells[i] = *(*C.d3plot_thick_shell)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return thickShells, nil
}

func (plotFile D3plot) ReadBeamsState(state uint64) ([]C.d3plot_beam, error) {
	var numBeams C.size_t
	dataC := C.d3plot_read_beams_state(&plotFile.handle, C.size_t(state), &numBeams)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numBeams == 0 {
		return []C.d3plot_beam{}, nil
	}

	beams := make([]C.d3plot_beam, numBeams)
	for i := range beams {
		beams[i] = *(*C.d3plot_beam)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return beams, nil
}

func (plotFile D3plot) ReadShellsState(state uint64) ([]C.d3plot_shell, error) {
	var numShells C.size_t
	dataC := C.d3plot_read_shells_state(&plotFile.handle, C.size_t(state), &numShells)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numShells == 0 {
		return []C.d3plot_shell{}, nil
	}

	shells := make([]C.d3plot_shell, numShells)
	for i := range shells {
		shells[i] = *(*C.d3plot_shell)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return shells, nil
}

func (plotFile D3plot) ReadSolidElements() ([]C.d3plot_solid_con, error) {
	var numSolids C.size_t
	dataC := C.d3plot_read_solid_elements(&plotFile.handle, &numSolids)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numSolids == 0 {
		return []C.d3plot_solid_con{}, nil
	}

	solids := make([]C.d3plot_solid_con, numSolids)
	for i := range solids {
		solids[i] = *(*C.d3plot_solid_con)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return solids, nil
}

func (plotFile D3plot) ReadThickShellElements() ([]C.d3plot_thick_shell_con, error) {
	var numThickShells C.size_t
	dataC := C.d3plot_read_thick_shell_elements(&plotFile.handle, &numThickShells)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numThickShells == 0 {
		return []C.d3plot_thick_shell_con{}, nil
	}

	thickShells := make([]C.d3plot_thick_shell_con, numThickShells)
	for i := range thickShells {
		thickShells[i] = *(*C.d3plot_thick_shell_con)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return thickShells, nil
}

func (plotFile D3plot) ReadBeamElements() ([]C.d3plot_beam_con, error) {
	var numBeams C.size_t
	dataC := C.d3plot_read_beam_elements(&plotFile.handle, &numBeams)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numBeams == 0 {
		return []C.d3plot_beam_con{}, nil
	}

	beams := make([]C.d3plot_beam_con, numBeams)
	for i := range beams {
		beams[i] = *(*C.d3plot_beam_con)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return beams, nil
}

func (plotFile D3plot) ReadShellElements() ([]C.d3plot_shell_con, error) {
	var numShells C.size_t
	dataC := C.d3plot_read_shell_elements(&plotFile.handle, &numShells)

	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return nil, err
	}

	if numShells == 0 {
		return []C.d3plot_shell_con{}, nil
	}

	shells := make([]C.d3plot_shell_con, numShells)
	for i := range shells {
		shells[i] = *(*C.d3plot_shell_con)(unsafe.Pointer(uintptr(unsafe.Pointer(dataC)) + uintptr(i)*unsafe.Sizeof(*dataC)))
	}
	C.free(unsafe.Pointer(dataC))

	return shells, nil
}

func (plotFile D3plot) ReadTitle() (string, error) {
	titleC := C.d3plot_read_title(&plotFile.handle)
	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return "", err
	}

	title := C.GoString(titleC)
	C.free(unsafe.Pointer(titleC))
	return title, nil
}

func (plotFile D3plot) ReadPart(partIndex uint64) (D3plotPart, error) {
	partC := C.d3plot_read_part(&plotFile.handle, C.size_t(partIndex))
	if plotFile.handle.error_string != nil {
		err := errors.New(C.GoString(plotFile.handle.error_string))
		return D3plotPart{}, err
	}

	var part D3plotPart

	if partC.num_solids != 0 {
		part.SolidIDs = carrToSlice[C.d3_word, uint64](partC.solid_ids, partC.num_solids)
	}
	if partC.num_thick_shells != 0 {
		part.ThickShellIDs = carrToSlice[C.d3_word, uint64](partC.thick_shell_ids, partC.num_thick_shells)
	}
	if partC.num_beams != 0 {
		part.BeamIDs = carrToSlice[C.d3_word, uint64](partC.beam_ids, partC.num_beams)
	}
	if partC.num_shells != 0 {
		part.ShellIDs = carrToSlice[C.d3_word, uint64](partC.shell_ids, partC.num_shells)
	}

	C.d3plot_free_part(&partC)

	return part, nil
}

func (plotFile D3plot) NumTimeSteps() uint64 {
	return uint64(plotFile.handle.num_states)
}

func D3plotIndexForID(id uint64, IDs []uint64) uint64 {
	return uint64(C.d3plot_index_for_id(C.d3_word(id), (*C.d3_word)(&IDs[0]), C.size_t(len(IDs))))
}
