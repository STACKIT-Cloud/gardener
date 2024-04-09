// Copyright 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package token_test

import (
	"context"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	resourcesv1alpha1 "github.com/gardener/gardener/pkg/apis/resources/v1alpha1"
	"github.com/gardener/gardener/pkg/nodeagent/apis/config"
	"github.com/gardener/gardener/pkg/nodeagent/controller/token"
	"github.com/gardener/gardener/pkg/utils"
)

var _ = Describe("Token controller tests", func() {
	var (
		testFS afero.Afero

		accessToken1, accessToken2 = []byte("access-token-1"), []byte("access-token-2")
		path1, path2               = "/some/path", "/some/other/path"
		secret1, secret2           *corev1.Secret
		syncPeriod                 = time.Second
	)

	BeforeEach(func() {
		secret1Name, err := utils.GenerateRandomString(64)
		Expect(err).NotTo(HaveOccurred())
		secret2Name, err := utils.GenerateRandomString(64)
		Expect(err).NotTo(HaveOccurred())

		secret1 = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(secret1Name),
				Namespace: metav1.NamespaceSystem,
			},
			Data: map[string][]byte{resourcesv1alpha1.DataKeyToken: accessToken1},
		}
		secret2 = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strings.ToLower(secret2Name),
				Namespace: metav1.NamespaceSystem,
			},
			Data: map[string][]byte{resourcesv1alpha1.DataKeyToken: accessToken2},
		}

		By("Setup manager")
		mgr, err := manager.New(restConfig, manager.Options{
			Metrics: metricsserver.Options{BindAddress: "0"},
			Cache: cache.Options{
				DefaultNamespaces: map[string]cache.Config{metav1.NamespaceSystem: {}},
			},
		})
		Expect(err).NotTo(HaveOccurred())

		By("Register controller")
		testFS = afero.Afero{Fs: afero.NewMemMapFs()}
		Expect((&token.Reconciler{
			FS: testFS,
			Config: config.TokenControllerConfig{
				SyncConfigs: []config.TokenSecretSyncConfig{
					{
						SecretName: secret1.Name,
						Path:       path1,
					},
					{
						SecretName: secret2.Name,
						Path:       path2,
					},
				},
				SyncPeriod: &metav1.Duration{Duration: syncPeriod},
			},
		}).AddToManager(mgr)).To(Succeed())

		By("Start manager")
		mgrContext, mgrCancel := context.WithCancel(ctx)

		go func() {
			defer GinkgoRecover()
			Expect(mgr.Start(mgrContext)).To(Succeed())
		}()

		DeferCleanup(func() {
			By("Stop manager")
			mgrCancel()
		})
	})

	JustBeforeEach(func() {
		By("Create access token secrets")
		Expect(testClient.Create(ctx, secret1)).To(Succeed())
		Expect(testClient.Create(ctx, secret2)).To(Succeed())
		DeferCleanup(func() {
			Expect(testClient.Delete(ctx, secret1)).To(Succeed())
			Expect(testClient.Delete(ctx, secret2)).To(Succeed())
		})
	})

	It("should write the tokens to the local file system", func() {
		Eventually(func(g Gomega) {
			token1OnDisk, err := afero.ReadFile(testFS, path1)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(token1OnDisk).To(Equal(accessToken1))

			token2OnDisk, err := afero.ReadFile(testFS, path2)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(token2OnDisk).To(Equal(accessToken2))
		}).Should(Succeed())
	})

	It("should update the tokens on the local file system after the sync period", func() {
		Eventually(func(g Gomega) {
			token1OnDisk, err := afero.ReadFile(testFS, path1)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(token1OnDisk).To(Equal(accessToken1))

			token2OnDisk, err := afero.ReadFile(testFS, path2)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(token2OnDisk).To(Equal(accessToken2))
		}).Should(Succeed())

		By("Update tokens in secret data")
		newToken1 := []byte("new-token1")
		patch := client.MergeFrom(secret1.DeepCopy())
		secret1.Data[resourcesv1alpha1.DataKeyToken] = newToken1
		Expect(testClient.Patch(ctx, secret1, patch)).To(Succeed())

		newToken2 := []byte("new-token1")
		patch = client.MergeFrom(secret2.DeepCopy())
		secret2.Data[resourcesv1alpha1.DataKeyToken] = newToken2
		Expect(testClient.Patch(ctx, secret2, patch)).To(Succeed())

		Eventually(func(g Gomega) {
			token1OnDisk, err := afero.ReadFile(testFS, path1)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(token1OnDisk).To(Equal(newToken1))

			token2OnDisk, err := afero.ReadFile(testFS, path2)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(token2OnDisk).To(Equal(newToken2))
		}).WithTimeout(2 * syncPeriod).Should(Succeed())
	})

	Context("unrelated secret", func() {
		BeforeEach(func() {
			secret1.Name = "some-other-secret"
			secret2.Name = "yet-another-secret"
		})

		It("should do nothing because the secret is unrelated", func() {
			Consistently(func(g Gomega) {
				token1OnDisk, err := afero.ReadFile(testFS, path1)
				g.Expect(token1OnDisk).To(BeEmpty())
				g.Expect(err).To(MatchError(ContainSubstring("file does not exist")))

				token2OnDisk, err := afero.ReadFile(testFS, path2)
				g.Expect(token2OnDisk).To(BeEmpty())
				g.Expect(err).To(MatchError(ContainSubstring("file does not exist")))
			}).Should(Succeed())
		})
	})
})
