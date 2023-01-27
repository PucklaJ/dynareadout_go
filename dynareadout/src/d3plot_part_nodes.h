/***********************************************************************************
 *                         This file is part of dynareadout
 *                    https://github.com/PucklaJ/dynareadout
 ***********************************************************************************
 * Copyright (c) 2022 Jonas Pucher
 *
 * This software is provided 'as-is', without any express or implied warranty.
 * In no event will the authors be held liable for any damages arising from the
 * use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not claim
 * that you wrote the original software. If you use this software in a product,
 * an acknowledgment in the product documentation would be appreciated but is
 * not required.
 *
 * 2. Altered source versions must be plainly marked as such, and must not be
 * misrepresented as being the original software.
 *
 * 3. This notice may not be removed or altered from any source distribution.
 ************************************************************************************/

#ifndef D3PLOT_PART_NODES_H
#define D3PLOT_PART_NODES_H

#include "d3_defines.h"
#include <stddef.h>

#ifndef D3PLOT_H
struct d3plot_file;
#endif

#ifdef __cplusplus
extern "C" {
#endif

/* Returns an array containing all node ids that are inside of the part.
 * The return value needs to be deallocated by free. This functions takes a
 * d3plot_part_get_node_ids_params struct. You can set the values of the struct
 * to optimize the functions performance. If you set params to NULL all data
 * will be retrieved, allocated and deallocated inside this one function call*/
d3_word *d3plot_part_get_node_ids(d3plot_file *plot_file,
                                  const d3plot_part *part,
                                  size_t *num_part_node_ids,
                                  d3plot_part_get_node_ids_params *params);
/* The same as d3plot_part_get_node_ids, but it returns indices instead of ids.
 * Those indices can be used to index into the node_ids array returned by
 * d3plot_read_node_ids. If you set params to NULL all data
 * will be retrieved, allocated and deallocated inside this one function call*/
d3_word *d3plot_part_get_node_indices(d3plot_file *plot_file,
                                      const d3plot_part *part,
                                      size_t *num_part_node_indices,
                                      d3plot_part_get_node_ids_params *params);

#ifdef __cplusplus
}
#endif

#endif