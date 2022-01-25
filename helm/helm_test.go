package helm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/thallesfreitaszup/helm-module/helm"
	"github.com/thallesfreitaszup/helm-module/helm/mocks"
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
		It("should return the correct manifests", func() {
			h := helm.New(source, mockGetter, options, dst)
			mockGetter.On("Get").Return(nil)
			manifests, err := h.Render()
			Expect(err).To(BeNil())
			Expect(len(manifests)).To(Equal(2))
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
			h := helm.New(source, mockGetter, options, dst)
			mockGetter.On("Get").Return(nil)
			manifests, err := h.Render()
			expectedError := "template: fake-app/templates/deployment.yaml:4:18: executing \"fake-app/templates/deployment.yaml\" at <.Xalues.xpto>: nil pointer evaluating interface {}.xpto"
			Expect(err.Error()).To(Equal(expectedError))
			Expect(len(manifests)).To(Equal(0))
		})
	})

})
