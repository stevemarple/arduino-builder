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

package test

import (
	"arduino.cc/builder"
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

type ByToolIDAndVersion []*types.Tool

func (s ByToolIDAndVersion) Len() int {
	return len(s)
}
func (s ByToolIDAndVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByToolIDAndVersion) Less(i, j int) bool {
	if s[i].Name == s[j].Name {
		return s[i].Version < s[j].Version
	}
	return s[i].Name < s[j].Name
}

func TestLoadTools(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "tools_builtin"}

	loader := builder.ToolsLoader{}
	err := loader.Run(context)
	NoError(t, err)

	tools := context[constants.CTX_TOOLS].([]*types.Tool)
	require.Equal(t, 7, len(tools))

	sort.Sort(ByToolIDAndVersion(tools))

	idx := 0
	require.Equal(t, "arm-none-eabi-gcc", tools[idx].Name)
	require.Equal(t, "4.8.3-2014q1", tools[idx].Version)
	require.Equal(t, Abs(t, "./downloaded_tools/arm-none-eabi-gcc/4.8.3-2014q1"), tools[idx].Folder)
	idx++
	require.Equal(t, "avr-gcc", tools[idx].Name)
	require.Equal(t, "4.8.1-arduino5", tools[idx].Version)
	require.Equal(t, Abs(t, "./downloaded_tools/avr-gcc/4.8.1-arduino5"), tools[idx].Folder)
	idx++
	require.Equal(t, "avrdude", tools[idx].Name)
	require.Equal(t, "6.0.1-arduino5", tools[idx].Version)
	require.Equal(t, Abs(t, "./downloaded_tools/avrdude/6.0.1-arduino5"), tools[idx].Folder)
	idx++
	require.Equal(t, "bossac", tools[idx].Name)
	require.Equal(t, "1.5-arduino", tools[idx].Version)
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.5-arduino"), tools[idx].Folder)
	idx++
	require.Equal(t, "bossac", tools[idx].Name)
	require.Equal(t, "1.6.1-arduino", tools[idx].Version)
	require.Equal(t, Abs(t, "./downloaded_tools/bossac/1.6.1-arduino"), tools[idx].Folder)
}

func TestLoadToolsWithBoardManagerFolderStructure(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_board_manager_stuff"}

	loader := builder.ToolsLoader{}
	err := loader.Run(context)
	NoError(t, err)

	tools := context[constants.CTX_TOOLS].([]*types.Tool)
	require.Equal(t, 2, len(tools))

	sort.Sort(ByToolIDAndVersion(tools))

	require.Equal(t, "arm-none-eabi-gcc", tools[0].Name)
	require.Equal(t, "4.8.3-2014q1", tools[0].Version)
	require.Equal(t, Abs(t, "./downloaded_board_manager_stuff/RFduino/tools/arm-none-eabi-gcc/4.8.3-2014q1"), tools[0].Folder)

	require.Equal(t, "openocd", tools[1].Name)
	require.Equal(t, "0.9.0-arduino", tools[1].Version)
	require.Equal(t, Abs(t, "./downloaded_board_manager_stuff/arduino/tools/openocd/0.9.0-arduino"), tools[1].Folder)
}

func TestLoadLotsOfTools(t *testing.T) {
	DownloadCoresAndToolsAndLibraries(t)

	context := make(map[string]interface{})
	context[constants.CTX_TOOLS_FOLDERS] = []string{"downloaded_tools", "tools_builtin", "downloaded_board_manager_stuff"}

	loader := builder.ToolsLoader{}
	err := loader.Run(context)
	NoError(t, err)

	tools := context[constants.CTX_TOOLS].([]*types.Tool)
	require.Equal(t, 8, len(tools))

	require.Equal(t, "arm-none-eabi-gcc", tools[0].Name)
	require.Equal(t, "4.8.3-2014q1", tools[0].Version)

	require.Equal(t, "openocd", tools[7].Name)
	require.Equal(t, "0.9.0-arduino", tools[7].Version)
	require.Equal(t, Abs(t, "./downloaded_board_manager_stuff/arduino/tools/openocd/0.9.0-arduino"), tools[7].Folder)
}
