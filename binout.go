package dynareadout

/*
#include "dynareadout/src/binout.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

const (
	BinoutTypeInt8    = 1
	BinoutTypeInt16   = 2
	BinoutTypeInt32   = 3
	BinoutTypeInt64   = 4
	BinoutTypeUint8   = 5
	BinoutTypeUint16  = 6
	BinoutTypeUint32  = 7
	BinoutTypeUint64  = 8
	BinoutTypeFloat32 = 9
	BinoutTypeFloat64 = 10
	BinoutTypeInvalid = math.MaxUint64
)

type cType interface {
	C.int8_t | C.int16_t | C.int32_t | C.int64_t | C.uint8_t | C.uint16_t | C.uint32_t | C.uint64_t | C.float | C.double
}

type goType interface {
	int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func carrIdxPtr[Tc cType](arr *Tc, idx int) *Tc {
	return (*Tc)(unsafe.Pointer(uintptr(unsafe.Pointer(arr)) + uintptr(idx)*unsafe.Sizeof(*arr)))
}

func carrIdx[Tc cType](arr *Tc, idx int) Tc {
	return *carrIdxPtr(arr, idx)
}

func carrToSlice[Tc cType, Ts goType](carr *Tc, csize C.size_t) []Ts {
	slice := make([]Ts, csize)
	for i := range slice {
		slice[i] = Ts(carrIdx(carr, i))
	}
	C.free(unsafe.Pointer(carr))
	return slice
}

type Binout struct {
	handle C.binout_file
}

func BinoutOpen(fileName string) (bin_file Binout, err error) {
	fileNameC := C.CString(fileName)

	bin_file.handle = C.binout_open(fileNameC)
	C.free(unsafe.Pointer(fileNameC))

	openErrorC := C.binout_open_error(&bin_file.handle)
	if openErrorC != nil {
		err = errors.New(C.GoString(openErrorC))
		C.free(unsafe.Pointer(openErrorC))
	}

	return
}

func (bin_file Binout) Close() {
	C.binout_close(&bin_file.handle)
}

func (bin_file Binout) ReadInt8(path string) ([]int8, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_i8(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.int8_t, int8](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadInt16(path string) ([]int16, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_i16(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.int16_t, int16](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadInt32(path string) ([]int32, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_i32(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.int32_t, int32](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadInt64(path string) ([]int64, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_i64(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.int64_t, int64](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadUint8(path string) ([]uint8, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_u8(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.uint8_t, uint8](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadUint16(path string) ([]uint16, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_u16(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.uint16_t, uint16](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadUint32(path string) ([]uint32, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_u32(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.uint32_t, uint32](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadUint64(path string) ([]uint64, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_u64(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.uint64_t, uint64](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadFloat32(path string) ([]float32, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_f32(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.float, float32](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadFloat64(path string) ([]float64, error) {
	pathC := C.CString(path)

	var dataSize C.size_t
	dataC := C.binout_read_f64(&bin_file.handle, pathC, &dataSize)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	slice := carrToSlice[C.double, float64](dataC, dataSize)

	return slice, nil
}

func (bin_file Binout) ReadString(path string) (string, error) {
	typeID := bin_file.GetTypeID(path)

	var dataC unsafe.Pointer
	var dataSize C.size_t
	pathC := C.CString(path)
	defer C.free(unsafe.Pointer(pathC))

	switch typeID {
	case BinoutTypeInt8:
		dataC = unsafe.Pointer(C.binout_read_i8(&bin_file.handle, pathC, &dataSize))
	case BinoutTypeUint8:
		dataC = unsafe.Pointer(C.binout_read_u8(&bin_file.handle, pathC, &dataSize))
	default:
		typeName := C.GoString(C._binout_get_type_name(C.uint64_t(typeID)))
		return "", fmt.Errorf("Type \"%s\" can not be converted to a string", typeName)
	}

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return "", err
	}

	str := C.GoStringN((*C.char)(dataC), C.int(dataSize))
	C.free(dataC)

	return str, nil
}

func (bin_file Binout) ReadTimedFloat32(path string) ([][]float32, error) {
	var numValues C.size_t
	var numTimesteps C.size_t
	pathC := C.CString(path)

	dataC := C.binout_read_timed_f32(&bin_file.handle, pathC, &numValues, &numTimesteps)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	data := make([][]float32, numTimesteps)
	for i := range data {
		data[i] = make([]float32, numValues)
		for j := range data[i] {
			data[i][j] = float32(carrIdx(dataC, i*int(numValues)+j))
		}
	}

	C.free(unsafe.Pointer(dataC))
	return data, nil
}

func (bin_file Binout) ReadTimedFloat64(path string) ([][]float64, error) {
	var numValues C.size_t
	var numTimesteps C.size_t
	pathC := C.CString(path)

	dataC := C.binout_read_timed_f64(&bin_file.handle, pathC, &numValues, &numTimesteps)
	C.free(unsafe.Pointer(pathC))

	if bin_file.handle.error_string != nil {
		err := errors.New(C.GoString(bin_file.handle.error_string))
		return nil, err
	}

	data := make([][]float64, numTimesteps)
	for i := range data {
		data[i] = make([]float64, numValues)
		for j := range data[i] {
			data[i][j] = float64(carrIdx(dataC, i*int(numValues)+j))
		}
	}

	C.free(unsafe.Pointer(dataC))
	return data, nil
}

func (bin_file Binout) GetTypeID(path string) uint64 {
	pathC := C.CString(path)

	typeID := C.binout_get_type_id(&bin_file.handle, pathC)
	C.free(unsafe.Pointer(pathC))

	return uint64(typeID)
}

func (bin_file Binout) GetChildren(path string) []string {
	pathC := C.CString(path)

	var numChildren C.size_t
	childrenC := C.binout_get_children(&bin_file.handle, pathC, &numChildren)

	if numChildren == 0 {
		return []string{}
	}

	children := make([]string, numChildren)
	for i := range children {
		childC := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(childrenC)) + uintptr(i)*unsafe.Sizeof(*childrenC)))
		children[i] = C.GoString(childC)
	}
	C.free(unsafe.Pointer(childrenC))

	return children
}

func (bin_file Binout) VariableExists(path string) bool {
	pathC := C.CString(path)
	defer C.free(unsafe.Pointer(pathC))

	return C.binout_variable_exists(&bin_file.handle, pathC) != 0
}

func (bin_file Binout) GetNumTimesteps(path string) (uint64, error) {
	pathC := C.CString(path)

	timesteps := C.binout_get_num_timesteps(&bin_file.handle, pathC)
	C.free(unsafe.Pointer(pathC))

	if timesteps == math.MaxUint64 {
		return 0, errors.New("The path does not exist or contains files")
	}

	return uint64(timesteps), nil
}
