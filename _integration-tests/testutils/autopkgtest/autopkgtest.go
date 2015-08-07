// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package autopkgtest

import (
	"fmt"
	"path/filepath"
	"strconv"

	"launchpad.net/snappy/_integration-tests/testutils"
	"launchpad.net/snappy/_integration-tests/testutils/tpl"
)

const (
	controlTpl    = "_integration-tests/data/tpl/control"
	dataOutputDir = "_integration-tests/data/output/"
	adtrunTpl     = "adt-run -B --setup-commands 'touch /run/autopkgtest_no_reboot.stamp' --override-control %s --built-tree %s --output-dir %s %s"
)

var (
	controlFile = filepath.Join(dataOutputDir, "control")
	// dependency aliasing
	execCommand      = testutils.ExecCommand
	prepareTargetDir = testutils.PrepareTargetDir
	tplExecute       = tpl.Execute
)

// Autopkgtest is the type that knows how to call adt-run
type Autopkgtest struct {
	sourceCodePath      string // location of the source code on the host
	testArtifactsPath   string // location of the test artifacts on the host
	testFilter          string
	integrationTestName string
}

// NewAutopkgtest is the Autopkgtest constructor
func NewAutopkgtest(sourceCodePath, testArtifactsPath, testFilter, integrationTestName string) *Autopkgtest {
	return &Autopkgtest{
		sourceCodePath:      sourceCodePath,
		testArtifactsPath:   testArtifactsPath,
		testFilter:          testFilter,
		integrationTestName: integrationTestName}
}

// AdtRunLocal starts a kvm running the image passed as argument and runs the
// autopkgtests using it as the testbed.
func (a *Autopkgtest) AdtRunLocal(imgPath string) {
	// Run the tests on the latest rolling edge image.
	a.adtRun(kvmSSHOptions(imgPath))
}

// AdtRunRemote runs the autopkgtests using a remote machine as the testbed.
func (a *Autopkgtest) AdtRunRemote(testbedIP string, testbedPort int) {
	execCommand("ssh-copy-id", "-p", strconv.Itoa(testbedPort),
		"ubuntu@"+testbedIP)
	a.adtRun(remoteTestbedSSHOptions(testbedIP, testbedPort))
}

func (a *Autopkgtest) adtRun(testbedOptions string) {
	a.createControlFile()

	fmt.Println("Calling adt-run...")
	outputDir := filepath.Join(a.testArtifactsPath, "output")
	prepareTargetDir(outputDir)

	cmdStr := fmt.Sprintf(adtrunTpl,
		controlFile, a.sourceCodePath, outputDir, testbedOptions)

	execCommand(cmdStr)
}

func (a *Autopkgtest) createControlFile() error {
	return tplExecute(controlTpl, controlFile,
		struct {
			Filter, Test string
		}{
			a.integrationTestName, a.testFilter})
}
