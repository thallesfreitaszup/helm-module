/*
 * Copyright 2022 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helm

import (
	"bytes"
	"errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"strings"
)

type ManifestGetter interface {
	Get() error
}
type Helm struct {
	Source         string
	ManifestGetter ManifestGetter
	Dst            string
	Options        Options
	Decoder        runtime.Codec
}

func New(source string, getter ManifestGetter, options Options, dst string) Helm {
	return Helm{
		Source:         source,
		ManifestGetter: getter,
		Options:        options,
		Decoder:        scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...),
		Dst:            dst,
	}
}

func (h Helm) Render() ([]unstructured.Unstructured, error) {
	var unstructuredManifests []unstructured.Unstructured
	unstructuredManifests, err := h.getCachedManifests()
	if err == nil {
		return unstructuredManifests, err
	}
	err = h.ManifestGetter.Get()
	if err != nil {
		return unstructuredManifests, err
	}
	helmChart, values, err := h.getChartAndValues()

	if err != nil {
		return unstructuredManifests, err
	}
	manifestList, err := engine.Render(helmChart, values)
	if err != nil {
		return unstructuredManifests, err
	}
	err = h.toUnstructured(manifestList, &unstructuredManifests)
	if err != nil {
		return unstructuredManifests, err
	}
	return unstructuredManifests, nil
}

func (h Helm) getCachedManifests() ([]unstructured.Unstructured, error) {
	if h.Options.Cache == nil {
		return nil, errors.New("no cache defined")
	}
	return h.Options.Cache.GetManifests(h.Source)
}

func (h Helm) toUnstructured(list map[string]string, unstructuredManifests *[]unstructured.Unstructured) error {

	var unstructuredRes unstructured.Unstructured
	for name, manifest := range list {
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}
		err := runtime.DecodeInto(h.Decoder, bytes.NewBufferString(manifest).Bytes(), &unstructuredRes)
		if err != nil {
			return err
		}
		*unstructuredManifests = append(*unstructuredManifests, unstructuredRes)
	}
	return nil
}

func (h Helm) getChartAndValues() (*chart.Chart, chartutil.Values, error) {
	chart, err := loader.Load(h.Dst)
	if err != nil {
		return nil, nil, err
	}
	values, err := chartutil.ToRenderValues(chart, chart.Values, chartutil.ReleaseOptions{}, nil)
	if err != nil {
		return nil, nil, err
	}
	return chart, values, nil
}

type Auth struct {
	SSHKey      string
	BearerToken string
}

type Cache interface {
	GetManifests(source string) ([]unstructured.Unstructured, error)
}

type Options struct {
	Cache Cache
	Auth  Auth
}
