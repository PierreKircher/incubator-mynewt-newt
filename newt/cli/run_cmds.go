/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cli

import (
	"os"

	"github.com/spf13/cobra"

	"mynewt.apache.org/newt/newt/newtutil"
	"mynewt.apache.org/newt/util"
)

func runRunCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		NewtUsage(cmd, util.NewNewtError("Must specify target"))
	}

	TryGetProject()

	b, err := TargetBuilderForTargetOrUnittest(args[0])
	if err != nil {
		NewtUsage(cmd, err)
	}

	testPkg := b.GetTestPkg()
	if testPkg != nil {
		b.InjectSetting("TESTUTIL_SYSTEM_ASSERT", "1")
		if err := b.SelfTestCreateExe(); err != nil {
			NewtUsage(nil, err)
		}
		if err := b.SelfTestDebug(); err != nil {
			NewtUsage(nil, err)
		}
	} else {
		if err := b.Build(); err != nil {
			NewtUsage(nil, err)
		}

		/*
		 * Run create-image if version number is specified. If no version
		 * number, remove .img which would'be been created. This so that
		 * download script will barf if it needs an image for this type of
		 * target, instead of downloading an older version.
		 */
		if len(args) > 1 {
			_, _, err = b.CreateImages(args[1], "", 0)
			if err != nil {
				NewtUsage(cmd, err)
			}
		} else {
			os.Remove(b.AppBuilder.AppImgPath())

			if b.LoaderBuilder != nil {
				os.Remove(b.LoaderBuilder.AppImgPath())
			}
		}

		if err := b.Load(extraJtagCmd); err != nil {
			NewtUsage(nil, err)
		}

		if err := b.Debug(extraJtagCmd, true, noGDB_flag); err != nil {
			NewtUsage(nil, err)
		}
	}
}

func AddRunCommands(cmd *cobra.Command) {
	runHelpText := "Same as running\n" +
		" - build <target>\n" +
		" - create-image <target> <version>\n" +
		" - load <target>\n" +
		" - debug <target>\n\n" +
		"Note if version number is omitted, create-image step is skipped\n"
	runHelpEx := "  newt run <target-name> [<version>]\n"

	runCmd := &cobra.Command{
		Use:     "run",
		Short:   "build/create-image/download/debug <target>",
		Long:    runHelpText,
		Example: runHelpEx,
		Run:     runRunCmd,
	}

	runCmd.PersistentFlags().StringVarP(&extraJtagCmd, "extrajtagcmd", "", "",
		"Extra commands to send to JTAG software")
	runCmd.PersistentFlags().BoolVarP(&noGDB_flag, "noGDB", "n", false,
		"Do not start GDB from command line")
	runCmd.PersistentFlags().BoolVarP(&newtutil.NewtForce,
		"force", "f", false,
		"Ignore flash overflow errors during image creation")

	cmd.AddCommand(runCmd)
	AddTabCompleteFn(runCmd, func() []string {
		return append(targetList(), unittestList()...)
	})
}
