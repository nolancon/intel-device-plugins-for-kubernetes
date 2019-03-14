// Copyright 2017 Intel Corporation. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/intel/intel-device-plugins-for-kubernetes/cmd/qat_plugin/kerneldrv"
	"github.com/intel/intel-device-plugins-for-kubernetes/cmd/qat_plugin/vfiodrv"
	"github.com/intel/intel-device-plugins-for-kubernetes/pkg/debug"
	"github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin"
)

const (
	namespace = "qat.intel.com"
)

func main() {
	var plugin deviceplugin.Scanner
	var err error

	debugEnabled := flag.Bool("debug", false, "enable debug output")
	mode := flag.String("mode", "vfio", "plugin mode which can be either vfio (default) or kernel")

	dpdkDriver := flag.String("dpdk-driver", "vfio-pci", "DPDK Device driver for configuring the QAT device")
	kernelVfDrivers := flag.String("kernel-vf-drivers", "dh895xccvf,c6xxvf,c3xxxvf,d15xxvf", "Comma separated VF Device Driver of the QuickAssist Devices in the system. Devices supported: DH895xCC,C62x,C3xxx and D15xx")
	maxNumDevices := flag.Int("max-num-devices", 32, "maximum number of QAT devices to be provided to the QuickAssist device plugin")
	flag.Parse()

	if *debugEnabled {
		debug.Activate()
	}

	switch *mode {
	case "vfio":
		plugin, err = vfiodrv.NewDevicePlugin(*maxNumDevices, *kernelVfDrivers, *dpdkDriver)
	case "kernel":
		plugin = kerneldrv.NewDevicePlugin()
	default:
		err = errors.Errorf("Uknown mode: %s", *mode)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("QAT device plugin started in '%s' mode\n", *mode)
	manager := deviceplugin.NewManager(namespace, plugin)
	manager.Run()
}
