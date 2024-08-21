package resource

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func ResourceUpdate(cmd *cobra.Command, args []string) error {
	namespace, _ := cmd.Flags().GetString("namespace")
	pvc, _ := cmd.Flags().GetString("pvc")
	fmt.Printf("Mounting %s in %s\n", pvc, namespace)

	// Generate a random string for the container name
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	containerName := fmt.Sprintf("pvcmount-%d", rng.Intn(100000))

	// Create the Pod object
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      containerName,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name: "pvc-storage",
					VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvc,
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name:  containerName,
					Image: "busybox:latest",
					VolumeMounts: []v1.VolumeMount{
						{
							MountPath: "/mnt",
							Name:      "pvc-storage",
						},
					},
					Command: []string{"/bin/sh"},
					TTY:     true,
					Stdin:   true,
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}

	// Load the kubeconfig from the default location
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	// Create the Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create the Pod in the specified namespace
	_, err = clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	// Wait for the Pod to be in the Running state
	for {
		p, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if p.Status.Phase == v1.PodRunning {
			break
		}
		fmt.Println("Pod is not running yet, retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	// Schedule the deletion of the Pod after the shell session ends
	defer func() {
		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("Failed to delete pod: %v\n", err)
		}
	}()

	// Exec into the Pod
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(namespace).
		SubResource("exec").
		Param("container", containerName).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true").
		Param("command", "/bin/sh").
		Param("command", "-c").
		Param("command", "cd /mnt && exec /bin/sh")

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		return err
	}

	return nil
}
