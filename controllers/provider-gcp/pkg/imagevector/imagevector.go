// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//go:generate packr2

package imagevector

import (
	"github.com/gobuffalo/packr/v2"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gardener/gardener/pkg/utils/imagevector"

	"k8s.io/apimachinery/pkg/util/runtime"
)

// ImageVector is the image vector that contains al the needed images
var ImageVector imagevector.ImageVector

func init() {
	box := packr.New("charts", "../../charts")

	imagesYaml, err := box.FindString("images.yaml")
	runtime.Must(err)

	tmpFile, err := ioutil.TempFile("", "")
	runtime.Must(err)
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	_, err = io.Copy(tmpFile, strings.NewReader(imagesYaml))
	runtime.Must(err)

	// TODO: ReadImageVector should be split into ReadImageVector and MergeImageVector
	ImageVector, err = imagevector.ReadImageVector(tmpFile.Name())
	runtime.Must(err)
}
