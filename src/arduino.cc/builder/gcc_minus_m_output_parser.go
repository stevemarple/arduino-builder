/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/utils"
	"strings"
)

type GCCMinusMOutputParser struct{}

func (s *GCCMinusMOutputParser) Run(context map[string]interface{}) error {
	output := context[constants.CTX_GCC_MINUS_M_OUTPUT].(string)

	rows := strings.Split(output, "\n")
	includes := make([]string, 0)
	if len(rows) > 2 {
		for _, row := range rows[2:] {
			if !strings.HasPrefix(row, constants.SPACE) {
				row = strings.TrimSpace(row)
				if row != constants.EMPTY_STRING {
					row = strings.TrimSuffix(row, ":")
					row = strings.Replace(row, "\\ ", " ", -1)
					includes = append(includes, row)
				}
			}
		}
	}

	if !utils.MapHas(context, constants.CTX_INCLUDES) {
		context[constants.CTX_INCLUDES] = includes
		return nil
	}

	previousIncludes := utils.SliceToMapStringBool(context[constants.CTX_INCLUDES].([]string), true)
	currentIncludes := utils.SliceToMapStringBool(includes, true)

	mergedIncludes := utils.MergeMapsOfStringBool(previousIncludes, currentIncludes)

	context[constants.CTX_INCLUDES] = utils.KeysOfMapOfStringBool(mergedIncludes)

	return nil
}
