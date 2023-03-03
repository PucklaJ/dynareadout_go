package dynareadout

/*
#include <errno.h>
#include <stdlib.h>
#include "dynareadout/src/key.h"

int get_errno() { return errno; }
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

const (
	CardParseInt    = C.CARD_PARSE_INT
	CardParseFloat  = C.CARD_PARSE_FLOAT
	CardParseString = C.CARD_PARSE_STRING

	DefaultValueWidth = 10
)

type Keywords struct {
	handle      *C.keyword_t
	numKeywords C.size_t
}

type Keyword struct {
	handle *C.keyword_t
}

type Card struct {
	handle *C.card_t
}

func KeyFileParse(fileName string, parseIncludes bool) (Keywords, error) {
	var keywords Keywords
	var errorString *C.char
	var parseIncludesC C.int

	if parseIncludes {
		parseIncludesC = 1
	} else {
		parseIncludesC = 0
	}

	fileNameC := C.CString(fileName)

	keywords.handle = C.key_file_parse(fileNameC, &keywords.numKeywords, parseIncludesC, &errorString)
	C.free(unsafe.Pointer(fileNameC))

	if errorString != nil {
		err := errors.New(C.GoString(errorString))
		C.free(unsafe.Pointer(errorString))

		return keywords, err
	}

	return keywords, nil
}

func (k *Keywords) Free() {
	C.key_file_free(k.handle, k.numKeywords)
}

func (k Keywords) Len() int {
	return int(k.numKeywords)
}

func (k *Keywords) Get(name string, index int) (Keyword, error) {
	var keyword Keyword
	nameC := C.CString(name)

	keyword.handle = C.key_file_get(k.handle, k.numKeywords, nameC, C.size_t(index))
	C.free(unsafe.Pointer(nameC))
	if keyword.handle == nil {
		return keyword, fmt.Errorf("Could not find keyword \"%s\" with index %d", name, index)
	}

	return keyword, nil
}

func (k *Keywords) GetSlice(name string) ([]Keyword, error) {
	var sliceSize C.size_t
	var keywordC *C.keyword_t
	nameC := C.CString(name)

	keywordC = C.key_file_get_slice(k.handle, k.numKeywords, nameC, &sliceSize)
	C.free(unsafe.Pointer(nameC))
	if keywordC == nil {
		return nil, fmt.Errorf("Could not find keyword \"%s\"", name)
	}

	keywords := make([]Keyword, sliceSize)

	for i := 0; i < int(sliceSize); i++ {
		sliceElement := (*C.keyword_t)(unsafe.Pointer(uintptr(unsafe.Pointer(keywordC)) + (uintptr(i) * unsafe.Sizeof(*keywordC))))
		keywords[i].handle = sliceElement
	}

	return keywords, nil
}

func (k Keyword) Len() int {
	return int(k.handle.num_cards)
}

func (k *Keyword) Get(index int) Card {
	if index < 0 || index >= k.Len() {
		panic("Index out of range")
	}

	var card Card
	card.handle = (*C.card_t)(unsafe.Pointer(uintptr(unsafe.Pointer(k.handle.cards)) + (uintptr(index) * unsafe.Sizeof(*k.handle.cards))))
	return card
}

func (k *Keyword) GetSlice() []Card {
	cards := make([]Card, k.Len())

	for i := 0; i < len(cards); i++ {
		cards[i].handle = (*C.card_t)(unsafe.Pointer(uintptr(unsafe.Pointer(k.handle.cards)) + (uintptr(i) * unsafe.Sizeof(*k.handle.cards))))
	}

	return cards
}

func (c *Card) ParseBegin(valueWidth int) {
	C.card_parse_begin(c.handle, C.uint8_t(valueWidth))
}

func (c *Card) ParseNext() {
	C.card_parse_next(c.handle)
}

func (c *Card) ParseNextWidth(valueWidth int) {
	C.card_parse_next_width(c.handle, C.uint8_t(valueWidth))
}

func (c Card) ParseDone() bool {
	return C.card_parse_done(c.handle) != 0
}

func (c Card) currentString() string {
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	for i := c.handle.current_index; i < c.handle.current_index+c.handle.value_width; i++ {
		b := *(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(c.handle.string)) + uintptr(i)))
		buffer.WriteByte(byte(b))
	}
	return buffer.String()
}

func (c Card) currentStringWidth(valueWidth int) string {
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	for i := c.handle.current_index; i < c.handle.current_index+C.uint8_t(valueWidth); i++ {
		b := *(*C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(c.handle.string)) + uintptr(i)))
		buffer.WriteByte(byte(b))
	}
	return buffer.String()
}

func (c Card) ParseInt() (int, error) {
	intC := C.card_parse_int(c.handle)
	if C.get_errno() != 0 {
		errStr := c.currentString()

		return 0, fmt.Errorf("Failed to parse \"%s\" as int", errStr)
	}
	return int(intC), nil
}

func (c Card) ParseIntWidth(valueWidth int) (int, error) {
	intC := C.card_parse_int_width(c.handle, C.uint8_t(valueWidth))
	if C.get_errno() != 0 {
		errStr := c.currentStringWidth(valueWidth)

		return 0, fmt.Errorf("Failed to parse \"%s\" as int", errStr)
	}
	return int(intC), nil
}

func (c Card) ParseFloat32() (float32, error) {
	floatC := C.card_parse_float32(c.handle)
	if C.get_errno() != 0 {
		errStr := c.currentString()

		return 0, fmt.Errorf("Failed to parse \"%s\" as float32", errStr)
	}
	return float32(floatC), nil
}

func (c Card) ParseFloat32Width(valueWidth int) (float32, error) {
	floatC := C.card_parse_float32_width(c.handle, C.uint8_t(valueWidth))
	if C.get_errno() != 0 {
		errStr := c.currentStringWidth(valueWidth)

		return 0, fmt.Errorf("Failed to parse \"%s\" as float32", errStr)
	}
	return float32(floatC), nil
}

func (c Card) ParseFloat64() (float64, error) {
	floatC := C.card_parse_float64(c.handle)
	if C.get_errno() != 0 {
		errStr := c.currentString()

		return 0, fmt.Errorf("Failed to parse \"%s\" as float64", errStr)
	}
	return float64(floatC), nil
}

func (c Card) ParseFloat64Width(valueWidth int) (float64, error) {
	floatC := C.card_parse_float64_width(c.handle, C.uint8_t(valueWidth))
	if C.get_errno() != 0 {
		errStr := c.currentStringWidth(valueWidth)

		return 0, fmt.Errorf("Failed to parse \"%s\" as float64", errStr)
	}
	return float64(floatC), nil
}

func (c Card) ParseString() string {
	strC := C.card_parse_string(c.handle)
	str := C.GoString(strC)
	C.free(unsafe.Pointer(strC))
	return str
}

func (c Card) ParseStringWidth(valueWidth int) string {
	strC := C.card_parse_string_width(c.handle, C.uint8_t(valueWidth))
	str := C.GoString(strC)
	C.free(unsafe.Pointer(strC))
	return str
}

func (c Card) ParseStringNoTrim() string {
	strC := C.card_parse_string_no_trim(c.handle)
	str := C.GoString(strC)
	C.free(unsafe.Pointer(strC))
	return str
}

func (c Card) ParseStringWidthNoTrim(valueWidth int) string {
	strC := C.card_parse_string_width_no_trim(c.handle, C.uint8_t(valueWidth))
	str := C.GoString(strC)
	C.free(unsafe.Pointer(strC))
	return str
}

func (c Card) ParseWhole() string {
	strC := C.card_parse_whole(c.handle)
	str := C.GoString(strC)
	C.free(unsafe.Pointer(strC))
	return str
}

func (c Card) ParseWholeNoTrim() string {
	strC := C.card_parse_whole_no_trim(c.handle)
	str := C.GoString(strC)
	C.free(unsafe.Pointer(strC))
	return str
}

func (c Card) ParseGetType() int {
	return int(C.card_parse_get_type(c.handle))
}

func(c Card) ParseGetTypeWidth(valueWidth int) int {
	return int(C.card_parse_get_type_width(c.handle, C.uint8_t(valueWidth)))
}