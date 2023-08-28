package dynareadout

/*
#cgo CFLAGS: -DTHREAD_SAFE
#cgo LDFLAGS: -lm
#include "dynareadout/src/binary_search.c"
#include "dynareadout/src/binout_directory.c"
#include "dynareadout/src/binout_glob.c"
#include "dynareadout/src/binout_read.c"
#include "dynareadout/src/binout.c"
#include "dynareadout/src/d3_buffer.c"
#include "dynareadout/src/d3plot_data.c"
#include "dynareadout/src/d3plot_part_nodes.c"
#include "dynareadout/src/d3plot_state.c"
#include "dynareadout/src/d3plot.c"
#include "dynareadout/src/extra_string.c"
#include "dynareadout/src/key.c"
#include "dynareadout/src/line.c"
#include "dynareadout/src/multi_file.c"
#include "dynareadout/src/path_view.c"
#include "dynareadout/src/path.c"
#include "dynareadout/src/sync.c"
#include "header.h"

int get_errno() { return errno; }
*/
import "C"
