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
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"mynewt.apache.org/newt/newtmgr/config"
	"mynewt.apache.org/newt/newtmgr/protocol"
	"mynewt.apache.org/newt/newtmgr/transport"
)

func runtestRunCmd(cmd *cobra.Command, args []string) {
	cpm, err := config.NewConnProfileMgr()
	if err != nil {
		nmUsage(cmd, err)
	}

	profile, err := cpm.GetConnProfile(ConnProfileName)
	if err != nil {
		nmUsage(cmd, err)
	}

	conn, err := transport.NewConnWithTimeout(profile, time.Second*1)
	if err != nil {
		nmUsage(cmd, err)
	}
	defer conn.Close()

	runner, err := protocol.NewCmdRunner(conn)
	if err != nil {
		nmUsage(cmd, err)
	}

	runtest, err := protocol.RunTest()
	if err != nil {
		nmUsage(cmd, err)
	}

	nmr, err := runtest.EncodeWriteRequest()
	if err != nil {
		nmUsage(cmd, err)
	}

	if err := runner.WriteReq(nmr); err != nil {
		nmUsage(cmd, err)
	}

	rsp, err := runner.ReadResp()
	if err == nil {
		cRsp, err := protocol.DecodeRunTestResponse(rsp.Data)
		if err != nil {
			nmUsage(cmd, err)
		}
		if cRsp.Err != 0 {
			fmt.Printf("Failed, error:%d\n", cRsp.Err)
		}
	}
	fmt.Println("Done")
}

func runtestCmd() *cobra.Command {
	runtestEx := "   runtest "

	runtestCmd := &cobra.Command{
		Use:     "runtest ",
		Short:   "Initiate named test on remote endpoint using newtmgr (named test not yet supported)",
		Example: runtestEx,
		Run:     runtestRunCmd,
	}

	return runtestCmd
}
