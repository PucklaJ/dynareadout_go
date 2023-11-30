package dynareadout

/*
#cgo CFLAGS: -ansi
#include <errno.h>
#include <stdlib.h>
#include "dynareadout/src/key.h"
#include "dynareadout/src/include_transform.h"
#include "header.h"
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"sync"
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

type KeyFileParseCallback func(KeyParseInfo, string, *Card, int)

type KeyFileParseConfig struct {
	ParseIncludes          bool
	IgnoreNotFoundIncludes bool
	ExtraIncludePaths      []string
}

type KeyFileWarning struct {
	warningString string
}

type KeyParseInfo struct {
	handle C.key_parse_info_t
}

type IncludeTransform struct {
	handle C.include_transform_t
}

type DefineTransformation struct {
	handle C.define_transformation_t
}

type TransformOption struct {
	Name       string
	Parameters [7]float64
}

func (w *KeyFileWarning) Error() string {
	return w.warningString
}

func DefaultKeyFileParseConfig() KeyFileParseConfig {
	return KeyFileParseConfig{
		ParseIncludes:          true,
		IgnoreNotFoundIncludes: false,
	}
}

func KeyFileParse(fileName string, parseConfig KeyFileParseConfig) (Keywords, *KeyFileWarning, error) {
	var keywords Keywords
	var warning *KeyFileWarning
	var errorString *C.char
	var warningString *C.char

	cParseConfig := parseConfig.toC()
	fileNameC := C.CString(fileName)

	keywords.handle = C.key_file_parse(fileNameC, &keywords.numKeywords, &cParseConfig, &errorString, &warningString)
	C.free(unsafe.Pointer(fileNameC))

	if cParseConfig.extra_include_paths != nil {
		for i := C.size_t(0); i < cParseConfig.num_extra_include_paths; i++ {
			ptr := (**C.char)(unsafe.Add(unsafe.Pointer(cParseConfig.extra_include_paths), uintptr(i)*unsafe.Sizeof(*cParseConfig.extra_include_paths)))
			C.free(unsafe.Pointer(*ptr))
		}
		C.free(unsafe.Pointer(cParseConfig.extra_include_paths))
	}

	if warningString != nil {
		warning = new(KeyFileWarning)
		warning.warningString = C.GoString(warningString)
		C.free(unsafe.Pointer(warningString))
	}

	if errorString != nil {
		err := errors.New(C.GoString(errorString))
		C.free(unsafe.Pointer(errorString))

		return keywords, warning, err
	}

	return keywords, warning, nil
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

func (c Card) ParseGetTypeWidth(valueWidth int) int {
	return int(C.card_parse_get_type_width(c.handle, C.uint8_t(valueWidth)))
}

func (c Card) TryParseInt(value *int) {
	valueC := C.int64_t(*value)
	C._card_try_parse_int(c.handle, &valueC)
	*value = int(valueC)
}

func (c Card) TryParseFloat64(value *float64) {
	valueC := C.double(*value)
	C._card_try_parse_float64(c.handle, &valueC)
	*value = float64(valueC)
}

func (i KeyParseInfo) FileName() string {
	return C.GoString(i.handle.file_name)
}

func (i KeyParseInfo) LineNumber() int {
	return int(i.handle.line_number)
}

func (i KeyParseInfo) IncludePaths() []string {
	rv := make([]string, i.handle.num_include_paths)
	for j := C.size_t(0); j < i.handle.num_include_paths; j++ {
		pathC := (*C.char)(unsafe.Add(unsafe.Pointer(i.handle.include_paths), uintptr(j)*unsafe.Sizeof(*i.handle.include_paths)))
		rv = append(rv, C.GoString(pathC))
	}
	return rv
}

func (i KeyParseInfo) RootFolder() string {
	return C.GoString(i.handle.root_folder)
}

func (c KeyFileParseConfig) toC() (cfg C.key_parse_config_t) {
	if c.ParseIncludes {
		cfg.parse_includes = 1
	}
	if c.IgnoreNotFoundIncludes {
		cfg.ignore_not_found_includes = 1
	}
	if len(c.ExtraIncludePaths) != 0 {
		cfg.extra_include_paths = (**C.char)(C.malloc(C.size_t(uintptr(len(c.ExtraIncludePaths)) * unsafe.Sizeof(*cfg.extra_include_paths))))
		if cfg.extra_include_paths != nil {
			cfg.num_extra_include_paths = C.size_t(len(c.ExtraIncludePaths))
			for i, v := range c.ExtraIncludePaths {
				ptr := (**C.char)(unsafe.Add(unsafe.Pointer(cfg.extra_include_paths), uintptr(i)*unsafe.Sizeof(*cfg.extra_include_paths)))
				*ptr = C.CString(v)
			}
		}
	}
	return
}

var keyFileCallbacks map[uintptr]KeyFileParseCallback
var keyFileCallbacksMtx sync.Mutex

func KeyFileParseWithCallback(fileName string, callback KeyFileParseCallback, parseConfig KeyFileParseConfig) (*KeyFileWarning, error) {
	fileNameC := C.CString(fileName)
	var errorString *C.char
	var warningString *C.char
	var warning *KeyFileWarning

	cParseConfig := parseConfig.toC()

	keyFileCallbacksMtx.Lock()
	if keyFileCallbacks == nil {
		keyFileCallbacks = make(map[uintptr]KeyFileParseCallback)
	}
	var callbackIndex uintptr
	for callbackIndex = 0; ; callbackIndex++ {
		if _, ok := keyFileCallbacks[callbackIndex]; !ok {
			break
		}
	}
	keyFileCallbacks[callbackIndex] = callback
	keyFileCallbacksMtx.Unlock()
	defer func() {
		keyFileCallbacksMtx.Lock()
		defer keyFileCallbacksMtx.Unlock()
		delete(keyFileCallbacks, callbackIndex)
	}()

	C.key_file_parse_with_callback(fileNameC,
		C.key_file_callback(C.keyFileParseGoCallback),
		&cParseConfig,
		&errorString,
		&warningString,
		unsafe.Pointer(callbackIndex),
		nil,
	)
	C.free(unsafe.Pointer(fileNameC))

	if cParseConfig.extra_include_paths != nil {
		for i := C.size_t(0); i < cParseConfig.num_extra_include_paths; i++ {
			ptr := (**C.char)(unsafe.Add(unsafe.Pointer(cParseConfig.extra_include_paths), uintptr(i)*unsafe.Sizeof(*cParseConfig.extra_include_paths)))
			C.free(unsafe.Pointer(*ptr))
		}
		C.free(unsafe.Pointer(cParseConfig.extra_include_paths))
	}

	if warningString != nil {
		warning = new(KeyFileWarning)
		warning.warningString = C.GoString(warningString)
		C.free(unsafe.Pointer(warningString))
	}

	if errorString != nil {
		errStr := C.GoString(errorString)
		C.free(unsafe.Pointer(errorString))
		return warning, errors.New(errStr)
	}

	return warning, nil
}

//export keyFileParseGoCallback
func keyFileParseGoCallback(infoC C.key_parse_info_t, keywordNameC *C.char, cardC *C.card_t, cardIndexC C.size_t, userData unsafe.Pointer) {
	callbackIndex := uintptr(userData)
	keyFileCallbacksMtx.Lock()
	defer keyFileCallbacksMtx.Unlock()
	callback := keyFileCallbacks[callbackIndex]

	info := KeyParseInfo{infoC}
	keywordName := C.GoString(keywordNameC)
	var card *Card
	if cardC != nil {
		card = new(Card)
		card.handle = cardC
	}
	var cardIndex int
	if cardIndexC == C.size_t(math.MaxUint64) {
		cardIndex = math.MaxInt
	} else {
		cardIndex = int(cardIndexC)
	}

	callback(info, keywordName, card, cardIndex)
}

func KeyParseIncludeTransform(kw Keyword) IncludeTransform {
	handle := C.key_parse_include_transform(kw.handle)
	return IncludeTransform{handle}
}

func (it *IncludeTransform) Free() {
	C.key_free_include_transform(&it.handle)
}

func (it *IncludeTransform) ParseCard(card Card, cardIndex int) {
	C.key_parse_include_transform_card(&it.handle, card.handle, C.size_t(cardIndex))
}

func (it IncludeTransform) FileName() string {
	return C.GoString(it.handle.file_name)
}

func (it IncludeTransform) Idnoff() int {
	return int(it.handle.idnoff)
}

func (it IncludeTransform) Ideoff() int {
	return int(it.handle.ideoff)
}

func (it IncludeTransform) Idpoff() int {
	return int(it.handle.idpoff)
}

func (it IncludeTransform) Idmoff() int {
	return int(it.handle.idmoff)
}

func (it IncludeTransform) Idsoff() int {
	return int(it.handle.idsoff)
}

func (it IncludeTransform) Idfoff() int {
	return int(it.handle.iddoff)
}

func (it IncludeTransform) idroff() int {
	return int(it.handle.idroff)
}

func (it IncludeTransform) Prefix() string {
	return C.GoString(it.handle.prefix)
}

func (it IncludeTransform) Suffix() string {
	return C.GoString(it.handle.suffix)
}

func (it IncludeTransform) Fctmas() float64 {
	return float64(it.handle.fctmas)
}

func (it IncludeTransform) Fcttim() float64 {
	return float64(it.handle.fcttim)
}

func (it IncludeTransform) Fctlen() float64 {
	return float64(it.handle.fctlen)
}

func (it IncludeTransform) Fcttem() string {
	return C.GoString(it.handle.fcttem)
}

func (it IncludeTransform) Incout1() int {
	return int(it.handle.incout1)
}

func (it IncludeTransform) Tranid() int {
	return int(it.handle.tranid)
}

func KeyParseDefineTransformation(kw Keyword, isTitle bool) DefineTransformation {
	var isTitleC C.int
	if isTitle {
		isTitleC = 1
	}

	handle := C.key_parse_define_transformation(kw.handle, isTitleC)
	return DefineTransformation{handle}
}

func (dt *DefineTransformation) Free() {
	C.key_free_define_transformation(&dt.handle)
}

func (dt *DefineTransformation) ParseCard(card Card, cardIndex int, isTitle bool) {
	var isTitleC C.int
	if isTitle {
		isTitleC = 1
	}
	C.key_parse_define_transformation_card(&dt.handle, card.handle, C.size_t(cardIndex), isTitleC)
}

func (dt DefineTransformation) Tranid() int {
	return int(dt.handle.tranid)
}

func (dt DefineTransformation) Title() string {
	if dt.handle.title == nil {
		return ""
	}
	return C.GoString(dt.handle.title)
}

func (dt DefineTransformation) Options() []TransformOption {
	rv := make([]TransformOption, dt.handle.num_options)
	for i := 0; i < len(rv); i++ {
		optC := *(*C.transformation_option_t)(unsafe.Add(unsafe.Pointer(dt.handle.options), unsafe.Sizeof(*dt.handle.options)*uintptr(i)))

		rv[i].Name = C.GoString(optC.name)
		for j := 0; j < 7; j++ {
			rv[i].Parameters[j] = float64(optC.parameters[j])
		}
	}
	return rv
}
