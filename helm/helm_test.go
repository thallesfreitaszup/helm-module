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

package helm_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/hashicorp/go-getter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thallesfreitaszup/helm-module/helm"
	"github.com/thallesfreitaszup/helm-module/helm/mocks"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"path/filepath"
)

var _ = Describe("Helm", func() {
	var source, dst string
	var mockGetter *mocks.ManifestGetter
	var options helm.Options
	BeforeEach(func() {
		source = "repo.example/org/repo"
		dst = "./fake-chart/fake-app"
		mockGetter = new(mocks.ManifestGetter)
		options = helm.Options{}
	})
	Context("when the source is a valid chart", func() {
		It("should download from source and return the correct manifests", func() {
			randBytes := make([]byte, 16)
			rand.Read(randBytes)
			dst := filepath.Join(os.TempDir(), hex.EncodeToString(randBytes))
			pwd, _ := os.Getwd()
			client := getter.Client{
				Src:  "./fake-chart/fake-app",
				Dst:  dst,
				Mode: getter.ClientModeAny,
				Pwd:  pwd,
				Ctx:  context.TODO(),
			}
			h := helm.New(client.Src, &client, options, client.Dst)
			manifests, err := h.Render()
			Expect(err).To(BeNil())
			Expect(len(manifests)).To(Equal(2))
		})

		It("should return the correct manifests", func() {
			h := helm.New(source, mockGetter, options, dst)
			mockGetter.On("Get").Return(nil)
			manifests, err := h.Render()
			Expect(err).To(BeNil())
			Expect(len(manifests)).To(Equal(2))
		})
	})

	Context("when fails to download chart", func() {
		It("should return error", func() {
			expectedError := "failed to download repo"
			h := helm.New(source, mockGetter, options, dst)
			mockGetter.On("Get").Return(errors.New(expectedError))
			manifests, err := h.Render()
			Expect(err.Error()).To(Equal(expectedError))
			Expect(len(manifests)).To(Equal(0))
		})
	})

	Context("when there cached manifests should return it", func() {
		It("should return error", func() {
			mockCache := new(mocks.Cache)
			options.Cache = mockCache
			h := helm.New(source, mockGetter, options, dst)
			mockCache.On("GetManifests", source).Return(getUnstructuredManifests(), nil)
			manifests, err := h.Render()
			mockGetter.On("Get").Times(0)
			Expect(err).To(BeNil())
			Expect(len(manifests)).To(Equal(1))
			Expect(manifests).To(Equal(getUnstructuredManifests()))
		})
	})

	Context("when fails to load chart", func() {
		It("should return error", func() {
			dst = "./wrong-path"
			h := helm.New(source, mockGetter, options, dst)
			mockGetter.On("Get").Return(nil)
			manifests, err := h.Render()
			Expect(err).To(Not(BeNil()))
			Expect(len(manifests)).To(Equal(0))
		})
	})

	Context("when there errors on chart schema", func() {
		It("should return error", func() {
			dst = "./fake-chart/fake-app-with-schema"
			h := helm.New(source, mockGetter, options, dst)
			mockGetter.On("Get").Return(nil)
			manifests, err := h.Render()
			errorSubstring := "values don't meet the specifications of the schema(s) in the following chart(s):"
			Expect(err.Error()).To(ContainSubstring(errorSubstring))
			Expect(len(manifests)).To(Equal(0))
		})
	})

	Context("when fails to render manifest", func() {
		It("should return error", func() {
			dst = "./fake-chart/fake-app-invalid"
			expectedErr := "template: fake-app/templates/deployment.yaml:4:18: executing \"fake-app/templates/deployment.yaml\" at <.Xalues.xpto>: nil pointer evaluating interface {}.xpto"
			h := helm.New(source, mockGetter, options, dst)
			mockGetter.On("Get").Return(nil)
			manifests, err := h.Render()
			Expect(err.Error()).To(Equal(expectedErr))
			Expect(len(manifests)).To(Equal(0))

		})
	})

})

func getUnstructuredManifests() []unstructured.Unstructured {
	unstructuredManifest := unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "fake-deployment",
			},
		},
	}
	return []unstructured.Unstructured{unstructuredManifest}
}
