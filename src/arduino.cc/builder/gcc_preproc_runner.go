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
	"arduino.cc/builder/builder_utils"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/i18n"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"path/filepath"
	"strings"
)

type GCCPreprocRunner struct{}

func (s *GCCPreprocRunner) Run(context map[string]interface{}) error {
	buildProperties := utils.GetMapStringStringOrDefault(context, constants.CTX_BUILD_PROPERTIES)
	properties := utils.MergeMapsOfStrings(make(map[string]string), buildProperties)

	sketchBuildPath := context[constants.CTX_SKETCH_BUILD_PATH].(string)
	sketch := context[constants.CTX_SKETCH].(*types.Sketch)
	properties[constants.BUILD_PROPERTIES_SOURCE_FILE] = filepath.Join(sketchBuildPath, filepath.Base(sketch.MainFile.Name)+".cpp")

	includes := context[constants.CTX_INCLUDE_FOLDERS].([]string)
	includes = utils.Map(includes, utils.WrapWithHyphenI)
	properties[constants.BUILD_PROPERTIES_INCLUDES] = strings.Join(includes, constants.SPACE)
	builder_utils.RemoveHyphenMDDFlagFromGCCCommandLine(properties)

	verbose := context[constants.CTX_VERBOSE].(bool)
	logger := context[constants.CTX_LOGGER].(i18n.Logger)
	output, err := builder_utils.ExecRecipe(properties, constants.RECIPE_PREPROC_MACROS, true, verbose, false, logger)
	if err != nil {
		return utils.WrapError(err)
	}

	context[constants.CTX_GCC_MINUS_E_SOURCE] = string(output)

	return nil
}
