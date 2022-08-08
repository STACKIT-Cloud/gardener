// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package botanist_test

import (
	"context"
	"net"
	"path/filepath"

	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cr "github.com/gardener/gardener/pkg/chartrenderer"
	"github.com/gardener/gardener/pkg/client/kubernetes"
	"github.com/gardener/gardener/pkg/client/kubernetes/fake"
	"github.com/gardener/gardener/pkg/features"
	"github.com/gardener/gardener/pkg/gardenlet/apis/config"
	gardenletfeatures "github.com/gardener/gardener/pkg/gardenlet/features"
	"github.com/gardener/gardener/pkg/operation"
	. "github.com/gardener/gardener/pkg/operation/botanist"
	mockcoredns "github.com/gardener/gardener/pkg/operation/botanist/component/coredns/mock"
	mocketcd "github.com/gardener/gardener/pkg/operation/botanist/component/etcd/mock"
	mockhvpa "github.com/gardener/gardener/pkg/operation/botanist/component/hvpa/mock"
	mockkubeapiserver "github.com/gardener/gardener/pkg/operation/botanist/component/kubeapiserver/mock"
	mockkubecontrollermanager "github.com/gardener/gardener/pkg/operation/botanist/component/kubecontrollermanager/mock"
	mockkubeproxy "github.com/gardener/gardener/pkg/operation/botanist/component/kubeproxy/mock"
	mockkubescheduler "github.com/gardener/gardener/pkg/operation/botanist/component/kubescheduler/mock"
	mockkubestatemetrics "github.com/gardener/gardener/pkg/operation/botanist/component/kubestatemetrics/mock"
	mockresourcemanager "github.com/gardener/gardener/pkg/operation/botanist/component/resourcemanager/mock"
	mockvpa "github.com/gardener/gardener/pkg/operation/botanist/component/vpa/mock"
	mockvpnseedserver "github.com/gardener/gardener/pkg/operation/botanist/component/vpnseedserver/mock"
	mockvpnshoot "github.com/gardener/gardener/pkg/operation/botanist/component/vpnshoot/mock"
	gardenpkg "github.com/gardener/gardener/pkg/operation/garden"
	seedpkg "github.com/gardener/gardener/pkg/operation/seed"
	shootpkg "github.com/gardener/gardener/pkg/operation/shoot"
	"github.com/gardener/gardener/pkg/utils/imagevector"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	secretsmanager "github.com/gardener/gardener/pkg/utils/secrets/manager"
	fakesecretsmanager "github.com/gardener/gardener/pkg/utils/secrets/manager/fake"
	"github.com/gardener/gardener/pkg/utils/test"
	. "github.com/gardener/gardener/pkg/utils/test/matchers"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Monitoring", func() {
	var (
		ctrl *gomock.Controller

		fakeGardenClient client.Client
		k8sGardenClient  kubernetes.Interface

		fakeSeedClient client.Client
		k8sSeedClient  kubernetes.Interface

		chartApplier kubernetes.ChartApplier
		sm           secretsmanager.Interface

		mockHVPA                  *mockhvpa.MockInterface
		mockEtcdMain              *mocketcd.MockInterface
		mockEtcdEvents            *mocketcd.MockInterface
		mockKubeAPIServer         *mockkubeapiserver.MockInterface
		mockKubeScheduler         *mockkubescheduler.MockInterface
		mockKubeControllerManager *mockkubecontrollermanager.MockInterface
		mockKubeStateMetrics      *mockkubestatemetrics.MockInterface
		mockCoreDNS               *mockcoredns.MockInterface
		mockKubeProxy             *mockkubeproxy.MockInterface
		mockVPNSeedServer         *mockvpnseedserver.MockInterface
		mockVPNShoot              *mockvpnshoot.MockInterface
		mockResourceManager       *mockresourcemanager.MockInterface
		mockVPA                   *mockvpa.MockInterface

		botanist *Botanist

		ctx              = context.TODO()
		projectNamespace = "garden-foo"
		seedNamespace    = "shoot--foo--bar"
		shootName        = "bar"

		shootPurposeEvaluation = gardencorev1beta1.ShootPurposeEvaluation
		shootPurposeTesting    = gardencorev1beta1.ShootPurposeTesting
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		fakeGardenClient = fakeclient.NewClientBuilder().WithScheme(kubernetes.GardenScheme).Build()
		fakeSeedClient = fakeclient.NewClientBuilder().WithScheme(kubernetes.SeedScheme).Build()

		mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{corev1.SchemeGroupVersion, appsv1.SchemeGroupVersion})
		renderer := cr.NewWithServerVersion(&version.Info{GitVersion: "1.2.3"})
		chartApplier = kubernetes.NewChartApplier(renderer, kubernetes.NewApplier(fakeSeedClient, mapper))

		k8sGardenClient = fake.NewClientSetBuilder().
			WithClient(fakeGardenClient).
			Build()
		k8sSeedClient = fake.NewClientSetBuilder().
			WithClient(fakeSeedClient).
			WithChartApplier(chartApplier).
			WithRESTConfig(&rest.Config{}).
			Build()
		sm = fakesecretsmanager.New(fakeSeedClient, seedNamespace)

		mockHVPA = mockhvpa.NewMockInterface(ctrl)
		mockEtcdMain = mocketcd.NewMockInterface(ctrl)
		mockEtcdEvents = mocketcd.NewMockInterface(ctrl)
		mockKubeAPIServer = mockkubeapiserver.NewMockInterface(ctrl)
		mockKubeScheduler = mockkubescheduler.NewMockInterface(ctrl)
		mockKubeControllerManager = mockkubecontrollermanager.NewMockInterface(ctrl)
		mockKubeStateMetrics = mockkubestatemetrics.NewMockInterface(ctrl)
		mockCoreDNS = mockcoredns.NewMockInterface(ctrl)
		mockKubeProxy = mockkubeproxy.NewMockInterface(ctrl)
		mockVPNSeedServer = mockvpnseedserver.NewMockInterface(ctrl)
		mockVPNShoot = mockvpnshoot.NewMockInterface(ctrl)
		mockResourceManager = mockresourcemanager.NewMockInterface(ctrl)
		mockVPA = mockvpa.NewMockInterface(ctrl)

		botanist = &Botanist{
			Operation: &operation.Operation{
				K8sGardenClient: k8sGardenClient,
				K8sSeedClient:   k8sSeedClient,
				SecretsManager:  sm,
				Config:          &config.GardenletConfiguration{},
				Garden: &gardenpkg.Garden{
					Project: &gardencorev1beta1.Project{},
				},
				Seed:                &seedpkg.Seed{},
				SeedNamespaceObject: &corev1.Namespace{},
				Shoot: &shootpkg.Shoot{
					SeedNamespace: seedNamespace,
					Networks: &shootpkg.Networks{
						Pods:     []net.IPNet{},
						Services: []net.IPNet{},
					},
					Components: &shootpkg.Components{
						ControlPlane: &shootpkg.ControlPlane{
							EtcdMain:              mockEtcdMain,
							EtcdEvents:            mockEtcdEvents,
							KubeAPIServer:         mockKubeAPIServer,
							KubeScheduler:         mockKubeScheduler,
							KubeControllerManager: mockKubeControllerManager,
							KubeStateMetrics:      mockKubeStateMetrics,
							ResourceManager:       mockResourceManager,
							VerticalPodAutoscaler: mockVPA,
							VPNSeedServer:         mockVPNSeedServer,
						},
						SystemComponents: &shootpkg.SystemComponents{
							CoreDNS:   mockCoreDNS,
							KubeProxy: mockKubeProxy,
							VPNShoot:  mockVPNShoot,
						},
						HVPA: mockHVPA,
					},
					ReversedVPNEnabled: true,
				},
				ImageVector: imagevector.ImageVector{
					{Name: "grafana"},
					{Name: "prometheus"},
					{Name: "configmap-reloader"},
					{Name: "blackbox-exporter"},
					{Name: "kube-state-metrics"},
				},
			},
		}

		botanist.Seed.SetInfo(&gardencorev1beta1.Seed{
			Status: gardencorev1beta1.SeedStatus{
				KubernetesVersion: pointer.String("1.2.3"),
			},
		})

		botanist.Shoot.SetInfo(&gardencorev1beta1.Shoot{
			ObjectMeta: metav1.ObjectMeta{
				Name:      shootName,
				Namespace: projectNamespace,
			},
			Spec: gardencorev1beta1.ShootSpec{
				Purpose: &shootPurposeEvaluation,
			},
			Status: gardencorev1beta1.ShootStatus{
				TechnicalID: seedNamespace,
			},
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#DeploySeedMonitoring", func() {
		BeforeEach(func() {
			DeferCleanup(test.WithFeatureGate(gardenletfeatures.FeatureGate, features.HVPA, true))

			mockHVPA.EXPECT().ScrapeConfigs()
			mockHVPA.EXPECT().AlertingRules()
			mockEtcdMain.EXPECT().ScrapeConfigs()
			mockEtcdMain.EXPECT().AlertingRules()
			mockEtcdEvents.EXPECT().ScrapeConfigs()
			mockEtcdEvents.EXPECT().AlertingRules()
			mockKubeAPIServer.EXPECT().ScrapeConfigs()
			mockKubeAPIServer.EXPECT().AlertingRules()
			mockKubeScheduler.EXPECT().ScrapeConfigs()
			mockKubeScheduler.EXPECT().AlertingRules()
			mockKubeControllerManager.EXPECT().ScrapeConfigs()
			mockKubeControllerManager.EXPECT().AlertingRules()
			mockKubeStateMetrics.EXPECT().ScrapeConfigs()
			mockKubeStateMetrics.EXPECT().AlertingRules()
			mockCoreDNS.EXPECT().ScrapeConfigs()
			mockCoreDNS.EXPECT().AlertingRules()
			mockKubeProxy.EXPECT().ScrapeConfigs()
			mockKubeProxy.EXPECT().AlertingRules()
			mockVPNSeedServer.EXPECT().ScrapeConfigs()
			mockVPNSeedServer.EXPECT().AlertingRules()
			mockVPNShoot.EXPECT().ScrapeConfigs()
			mockVPNShoot.EXPECT().AlertingRules()
			mockResourceManager.EXPECT().ScrapeConfigs()
			mockResourceManager.EXPECT().AlertingRules()
			mockVPA.EXPECT().ScrapeConfigs()
			mockVPA.EXPECT().AlertingRules()
		})

		It("should delete the legacy ingress secrets", func() {
			By("creating secrets managed outside of this function for whose secretsmanager.Get() will be called")
			Expect(fakeSeedClient.Create(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "ca", Namespace: seedNamespace}})).To(Succeed())
			Expect(fakeSeedClient.Create(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "ca-etcd", Namespace: seedNamespace}})).To(Succeed())
			Expect(fakeSeedClient.Create(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "etcd-client", Namespace: seedNamespace}})).To(Succeed())
			Expect(fakeSeedClient.Create(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "observability-ingress", Namespace: seedNamespace}})).To(Succeed())
			Expect(fakeSeedClient.Create(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "observability-ingress-users", Namespace: seedNamespace}})).To(Succeed())
			Expect(fakeSeedClient.Create(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "generic-token-kubeconfig", Namespace: seedNamespace}})).To(Succeed())

			defer test.WithVar(&ChartsPath, filepath.Join("..", "..", "..", "charts"))()

			legacySecret1 := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: seedNamespace, Name: "prometheus-basic-auth"}}
			Expect(fakeSeedClient.Create(ctx, legacySecret1)).To(Succeed())

			legacySecret2 := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: seedNamespace, Name: "alertmanager-basic-auth"}}
			Expect(fakeSeedClient.Create(ctx, legacySecret2)).To(Succeed())

			Expect(botanist.DeploySeedMonitoring(ctx)).To(Succeed())

			Expect(fakeSeedClient.Get(ctx, client.ObjectKeyFromObject(legacySecret1), &corev1.Secret{})).To(BeNotFoundError())
			Expect(fakeSeedClient.Get(ctx, client.ObjectKeyFromObject(legacySecret2), &corev1.Secret{})).To(BeNotFoundError())
		})
	})

	Describe("#DeploySeedGrafana", func() {
		It("should generate two ingress secrets", func() {
			defer test.WithVar(&ChartsPath, filepath.Join("..", "..", "..", "charts"))()

			Expect(botanist.DeploySeedGrafana(ctx)).To(Succeed())

			secretList := &corev1.SecretList{}
			Expect(fakeSeedClient.List(ctx, secretList, client.InNamespace(seedNamespace), client.MatchingLabels{
				"name":       "observability-ingress",
				"managed-by": "secrets-manager",
			})).To(Succeed())
			Expect(secretList.Items).To(HaveLen(1))
			Expect(secretList.Items[0].Labels).To(HaveKeyWithValue("persist", "true"))

			Expect(fakeSeedClient.List(ctx, secretList, client.InNamespace(seedNamespace), client.MatchingLabels{
				"name":       "observability-ingress-users",
				"managed-by": "secrets-manager",
			})).To(Succeed())
			Expect(secretList.Items).To(HaveLen(1))
			Expect(secretList.Items[0].Labels).To(HaveKeyWithValue("persist", "true"))
		})

		It("should delete the legacy ingress secrets", func() {
			defer test.WithVar(&ChartsPath, filepath.Join("..", "..", "..", "charts"))()

			legacySecret1 := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: seedNamespace, Name: "monitoring-ingress-credentials"}}
			Expect(fakeSeedClient.Create(ctx, legacySecret1)).To(Succeed())

			legacySecret2 := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: seedNamespace, Name: "monitoring-ingress-credentials-users"}}
			Expect(fakeSeedClient.Create(ctx, legacySecret2)).To(Succeed())

			legacySecret3 := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: seedNamespace, Name: "grafana-users-basic-auth"}}
			Expect(fakeSeedClient.Create(ctx, legacySecret3)).To(Succeed())

			legacySecret4 := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: seedNamespace, Name: "grafana-operators-basic-auth"}}
			Expect(fakeSeedClient.Create(ctx, legacySecret4)).To(Succeed())

			Expect(botanist.DeploySeedGrafana(ctx)).To(Succeed())

			Expect(fakeSeedClient.Get(ctx, client.ObjectKeyFromObject(legacySecret1), &corev1.Secret{})).To(BeNotFoundError())
			Expect(fakeSeedClient.Get(ctx, client.ObjectKeyFromObject(legacySecret2), &corev1.Secret{})).To(BeNotFoundError())
			Expect(fakeSeedClient.Get(ctx, client.ObjectKeyFromObject(legacySecret3), &corev1.Secret{})).To(BeNotFoundError())
			Expect(fakeSeedClient.Get(ctx, client.ObjectKeyFromObject(legacySecret4), &corev1.Secret{})).To(BeNotFoundError())
		})

		It("should sync the ingress credentials for the users observability to the garden project namespace", func() {
			defer test.WithVar(&ChartsPath, filepath.Join("..", "..", "..", "charts"))()

			Expect(fakeGardenClient.Get(ctx, kutil.Key(projectNamespace, shootName+".monitoring"), &corev1.Secret{})).To(BeNotFoundError())

			Expect(botanist.DeploySeedGrafana(ctx)).To(Succeed())

			secret := &corev1.Secret{}
			Expect(fakeGardenClient.Get(ctx, kutil.Key(projectNamespace, shootName+".monitoring"), secret)).To(Succeed())
			Expect(secret.Annotations).To(HaveKeyWithValue("url", "https://gu-foo--bar."))
			Expect(secret.Labels).To(HaveKeyWithValue("gardener.cloud/role", "monitoring"))
			Expect(secret.Data).To(And(HaveKey("username"), HaveKey("password"), HaveKey("auth"), HaveKey("basic_auth.csv")))
		})

		It("should cleanup the secrets when shoot purpose is changed", func() {
			defer test.WithVar(&ChartsPath, filepath.Join("..", "..", "..", "charts"))()

			Expect(fakeGardenClient.Get(ctx, kutil.Key(projectNamespace, shootName+".monitoring"), &corev1.Secret{})).To(BeNotFoundError())

			Expect(botanist.DeploySeedGrafana(ctx)).To(Succeed())
			Expect(fakeGardenClient.Get(ctx, kutil.Key(projectNamespace, shootName+".monitoring"), &corev1.Secret{})).To(Succeed())
			Expect(*botanist.Shoot.GetInfo().Spec.Purpose == shootPurposeEvaluation).To(BeTrue())

			botanist.Shoot.Purpose = shootPurposeTesting
			Expect(botanist.DeploySeedGrafana(ctx)).To(Succeed())

			Expect(fakeGardenClient.Get(ctx, kutil.Key(projectNamespace, shootName+".monitoring"), &corev1.Secret{})).To(BeNotFoundError())
		})
	})
})
