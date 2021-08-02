package ingresses

import (
	"context"
	"testing"

	generictesting "github.com/loft-sh/vcluster/pkg/controllers/resources/generic/testing"
	"github.com/loft-sh/vcluster/pkg/util/locks"
	"github.com/loft-sh/vcluster/pkg/util/loghelper"
	testingutil "github.com/loft-sh/vcluster/pkg/util/testing"
	"github.com/loft-sh/vcluster/pkg/util/translate"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func newFakeSyncer(lockFactory locks.LockFactory, pClient *testingutil.FakeIndexClient, vClient *testingutil.FakeIndexClient) *syncer {
	return &syncer{
		sharedMutex:     lockFactory.GetLock("ingress-controller"),
		eventRecoder:    &testingutil.FakeEventRecorder{},
		targetNamespace: "test",
		virtualClient:   vClient,
		localClient:     pClient,
	}
}

func TestSync(t *testing.T) {
	vBaseSpec := networkingv1.IngressSpec{
		DefaultBackend: &networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: "testservice",
			},
			Resource: &corev1.TypedLocalObjectReference{
				Name: "testbackendresource",
			},
		},
		Rules: []networkingv1.IngressRule{
			{
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: "testbackendservice",
									},
									Resource: &corev1.TypedLocalObjectReference{
										Name: "testbackendresource",
									},
								},
							},
						},
					},
				},
			},
		},
		TLS: []networkingv1.IngressTLS{
			{
				SecretName: "testtlssecret",
			},
		},
	}
	pBaseSpec := networkingv1.IngressSpec{
		DefaultBackend: &networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: translate.PhysicalName("testservice", "testns"),
			},
			Resource: &corev1.TypedLocalObjectReference{
				Name: translate.PhysicalName("testbackendresource", "testns"),
			},
		},
		Rules: []networkingv1.IngressRule{
			{
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{
							{
								Backend: networkingv1.IngressBackend{
									Service: &networkingv1.IngressServiceBackend{
										Name: translate.PhysicalName("testbackendservice", "testns"),
									},
									Resource: &corev1.TypedLocalObjectReference{
										Name: translate.PhysicalName("testbackendresource", "testns"),
									},
								},
							},
						},
					},
				},
			},
		},
		TLS: []networkingv1.IngressTLS{
			{
				SecretName: translate.PhysicalName("testtlssecret", "testns"),
			},
		},
	}
	changedIngressStatus := networkingv1.IngressStatus{
		LoadBalancer: corev1.LoadBalancerStatus{
			Ingress: []corev1.LoadBalancerIngress{
				{
					IP:       "123:123:123:123",
					Hostname: "testhost",
				},
			},
		},
	}
	vObjectMeta := metav1.ObjectMeta{
		Name:        "testingress",
		Namespace:   "testns",
		ClusterName: "myvcluster",
	}
	pObjectMeta := metav1.ObjectMeta{
		Name:      translate.PhysicalName("testingress", "testns"),
		Namespace: "test",
		Labels: map[string]string{
			translate.MarkerLabel:    translate.Suffix,
			translate.NamespaceLabel: translate.NamespaceLabelValue(vObjectMeta.Namespace),
		},
	}
	baseIngress := &networkingv1.Ingress{
		ObjectMeta: vObjectMeta,
		Spec:       vBaseSpec,
	}
	createdIngress := &networkingv1.Ingress{
		ObjectMeta: pObjectMeta,
		Spec:       pBaseSpec,
	}
	updateIngress := &networkingv1.Ingress{
		ObjectMeta: vObjectMeta,
		Spec: networkingv1.IngressSpec{
			IngressClassName: stringPointer("updatedingressclass"),
		},
	}
	updatedIngress := &networkingv1.Ingress{
		ObjectMeta: pObjectMeta,
		Spec: networkingv1.IngressSpec{
			IngressClassName: stringPointer("updatedingressclass"),
		},
	}
	noUpdateIngress := &networkingv1.Ingress{
		ObjectMeta: vObjectMeta,
		Spec:       vBaseSpec,
		Status:     changedIngressStatus,
	}
	backwardUpdateIngress := &networkingv1.Ingress{
		ObjectMeta: pObjectMeta,
		Spec: networkingv1.IngressSpec{
			IngressClassName: stringPointer("backwardsupdatedingressclass"),
		},
		Status: changedIngressStatus,
	}
	backwardNoUpdateIngress := &networkingv1.Ingress{
		ObjectMeta: pObjectMeta,
		Spec:       networkingv1.IngressSpec{},
	}
	backwardUpdatedIngress := &networkingv1.Ingress{
		ObjectMeta: vObjectMeta,
		Spec: networkingv1.IngressSpec{
			DefaultBackend:   vBaseSpec.DefaultBackend,
			IngressClassName: stringPointer("backwardsupdatedingressclass"),
			Rules:            vBaseSpec.Rules,
			TLS:              vBaseSpec.TLS,
		},
		Status: changedIngressStatus,
	}
	lockFactory := locks.NewDefaultLockFactory()

	generictesting.RunTests(t, []*generictesting.SyncTest{
		{
			Name:                "Create forward",
			InitialVirtualState: []runtime.Object{baseIngress},
			ExpectedVirtualState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {baseIngress},
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {createdIngress},
			},
			Sync: func(ctx context.Context, pClient *testingutil.FakeIndexClient, vClient *testingutil.FakeIndexClient, scheme *runtime.Scheme, log loghelper.Logger) {
				syncer := newFakeSyncer(lockFactory, pClient, vClient)

				_, err := syncer.ForwardCreate(ctx, baseIngress, log)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:                 "Update forward",
			InitialVirtualState:  []runtime.Object{baseIngress},
			InitialPhysicalState: []runtime.Object{createdIngress},
			ExpectedVirtualState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {baseIngress},
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {updatedIngress},
			},
			Sync: func(ctx context.Context, pClient *testingutil.FakeIndexClient, vClient *testingutil.FakeIndexClient, scheme *runtime.Scheme, log loghelper.Logger) {
				syncer := newFakeSyncer(lockFactory, pClient, vClient)
				needed, err := syncer.ForwardUpdateNeeded(createdIngress, updateIngress)
				if err != nil {
					t.Fatal(err)
				} else if !needed {
					t.Fatal("Expected backward update to be needed")
				}

				_, err = syncer.ForwardUpdate(ctx, createdIngress, updateIngress, log)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:                 "Update forward not needed",
			InitialVirtualState:  []runtime.Object{baseIngress},
			InitialPhysicalState: []runtime.Object{createdIngress},
			ExpectedVirtualState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {baseIngress},
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {createdIngress},
			},
			Sync: func(ctx context.Context, pClient *testingutil.FakeIndexClient, vClient *testingutil.FakeIndexClient, scheme *runtime.Scheme, log loghelper.Logger) {
				syncer := newFakeSyncer(lockFactory, pClient, vClient)
				needed, err := syncer.ForwardUpdateNeeded(createdIngress, noUpdateIngress)
				if err != nil {
					t.Fatal(err)
				} else if needed {
					t.Fatal("Expected backward update to be not needed")
				}

				_, err = syncer.ForwardUpdate(ctx, createdIngress, noUpdateIngress, log)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:                 "Update backwards",
			InitialVirtualState:  []runtime.Object{baseIngress},
			InitialPhysicalState: []runtime.Object{createdIngress},
			ExpectedVirtualState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {backwardUpdatedIngress},
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {createdIngress},
			},
			Sync: func(ctx context.Context, pClient *testingutil.FakeIndexClient, vClient *testingutil.FakeIndexClient, scheme *runtime.Scheme, log loghelper.Logger) {
				syncer := newFakeSyncer(lockFactory, pClient, vClient)
				needed, err := syncer.BackwardUpdateNeeded(backwardUpdateIngress, baseIngress)
				if err != nil {
					t.Fatal(err)
				} else if !needed {
					t.Fatal("Expected backward update to be needed")
				}

				_, err = syncer.BackwardUpdate(ctx, backwardUpdateIngress, baseIngress, log)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			Name:                 "Update backwards not needed",
			InitialVirtualState:  []runtime.Object{baseIngress},
			InitialPhysicalState: []runtime.Object{createdIngress},
			ExpectedVirtualState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {baseIngress},
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				networkingv1.SchemeGroupVersion.WithKind("Ingress"): {createdIngress},
			},
			Sync: func(ctx context.Context, pClient *testingutil.FakeIndexClient, vClient *testingutil.FakeIndexClient, scheme *runtime.Scheme, log loghelper.Logger) {
				syncer := newFakeSyncer(lockFactory, pClient, vClient)
				needed, err := syncer.BackwardUpdateNeeded(backwardNoUpdateIngress, baseIngress)
				if err != nil {
					t.Fatal(err)
				} else if needed {
					t.Fatal("Expected backward update to be not needed")
				}

				_, err = syncer.BackwardUpdate(ctx, backwardNoUpdateIngress, baseIngress, log)
				if err != nil {
					t.Fatal(err)
				}
			},
		},
	})

}

func stringPointer(str string) *string {
	return &str
}
