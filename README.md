# argocd-playground

Demonstrating how to setup Argo CD on a k3s cluster using arkade and k3d. 

> This is just a learning playground

# Prerequisites
You need to install Docker on your machine and you need to register for a Docker Hub account as your Docker images will be stored there

# Install k3d
k3d is a little helper to run k3s in docker, where k3s is the lightweight Kubernetes distribution by Rancher. It actually removes millions of lines of code from k8s. If you just need a learning playground, k3s is definitely your choice.

Check out [k3d Github Page](https://github.com/rancher/k3d#get) to see the installation guide.

> When creating a cluster, ``k3d`` utilises ``kubectl`` and ``kubectl`` is not part of ``k3d``. If you don't have ``kubectl``, please install and set up [here](https://kubernetes.io/docs/tasks/tools/install-kubectl/). 

Once you've installed ``k3d`` and ``kubectl``, run
```
k3d create -n argocd-playground
```

We need to make ``kubectl`` to use the kubeconfig for that cluster.
```
export KUBECONFIG="$(k3d get-kubeconfig --name='argocd-playground')"
```

# Install arkade 
Moving on to [arkade](https://github.com/alexellis/arkade), it provides a simple Golang CLI with strongly-typed flags to install charts and apps to your cluster in one command. Originally, the codebase is derived from [k3sup](https://github.com/alexellis/k3sup) which I've contributed last month. 

```
curl -sLS https://dl.get-arkade.dev | sudo sh
```

Once you've installed it, you should see the following
```
New version of arkade installed to /usr/local/bin
            _             _
  __ _ _ __| | ____ _  __| | ___
 / _` | '__| |/ / _` |/ _` |/ _ \
| (_| | |  |   < (_| | (_| |  __/
 \__,_|_|  |_|\_\__,_|\__,_|\___|

Get Kubernetes apps the easy way

Version: 0.2.2
Git Commit: 9063b6eb16deae5978805f71b0e749828c815490
```

Install Argo CD via arkade. You can use an alias ``ark`` or ``arkade``. 
```
ark install argocd
```

You should see the following info
```
Using kubeconfig: /Users/wingkwong/.config/k3d/argocd-playground/kubeconfig.yaml
Node architecture: "amd64"
=======================================================================
= ArgoCD has been installed                                           =
=======================================================================


# Get the ArgoCD CLI

brew tap argoproj/tap
brew install argoproj/tap/argocd

# Or download via https://github.com/argoproj/argo-cd/releases/latest

# Username is "admin", get the password

kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2

# Port-forward

kubectl port-forward svc/argocd-server -n argocd 8081:443 &

http://localhost:8081

# Get started with ArgoCD at
# https://argoproj.github.io/argo-cd/#quick-start

Thanks for using arkade!
```

Follow the step to enable port forwarding 
```
kubectl port-forward svc/argocd-server -n argocd 8081:443 &
```

```
Forwarding from [::1]:8081 -> 8080
```

Open your browser and browse ``http://localhost:8080/``. You should see the Argo CD UI.
![image](https://user-images.githubusercontent.com/35857179/77913084-3bc3cc00-72c6-11ea-8175-572f46bfa626.png)

As stated in the console info upon the completion of installation, the username is ``admin`` and you can get hte password by running 

```
kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2
```

> If you want to check out the info, you can run ``ark info argocd``.

After logging in, you should see the application page.
![image](https://user-images.githubusercontent.com/35857179/77913544-053a8100-72c7-11ea-8047-7c2b5dc3b493.png)

Set your application name. Use the project ``default`` and choose the sync policy to ``Manual``.
![image](https://user-images.githubusercontent.com/35857179/77914661-d6250f00-72c8-11ea-9554-38afc52f7fdf.png)

Connect your repository to Argo CD. Select the revision and the path where your manifests files are located.
![image](https://user-images.githubusercontent.com/35857179/77918790-f35cdc00-72ce-11ea-93dc-488f50f947e1.png)

Set the cluster to ``https://kubernetes.default.svc`` with ``default`` namespace.
![image](https://user-images.githubusercontent.com/35857179/77914835-0ec4e880-72c9-11ea-8833-18e60096172b.png)

Click ``Create``. Then you should see there is an application on the portal.
![image](https://user-images.githubusercontent.com/35857179/77915215-a6c2d200-72c9-11ea-8683-06c5a8eb7c54.png)

You can also switch it to the list view
![image](https://user-images.githubusercontent.com/35857179/77915232-af1b0d00-72c9-11ea-89b0-1cbffc3b923c.png)

or summary view
![image](https://user-images.githubusercontent.com/35857179/77915240-b6421b00-72c9-11ea-8be3-c9e500d35e90.png)

Here is my application
```go
package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", Handler)

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}
```

Let's add ``deployment.yaml``
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: argocd-playground
spec:
  replicas: 1
  revisionHistoryLimit: 3
  selector:
    matchLabels:
      app: argocd-playground
  template:
    metadata:
      labels:
        app: argocd-playground
    spec:
      containers:
      - image: wingkwong/argocd-playground:v1
        name: argocd-playground
        ports:
        - containerPort: 8888
```

and ``service.yaml``

```yml
apiVersion: v1
kind: Service
metadata:
  name: argocd-playground
spec:
  ports:
  - port: 8888
    targetPort: 8888
  selector:
    app: argocd-playground
```

Once you've pushed your commit, Argo CD detects changes under ``manifests``. It updates the status to ``OutOfSync``.
![image](https://user-images.githubusercontent.com/35857179/77918199-20f55580-72ce-11ea-8784-d365b8af6b31.png)

Let's sync. 
![image](https://user-images.githubusercontent.com/35857179/77918456-7e89a200-72ce-11ea-88dd-625b97c64e3c.png)

Enable port forwarding 
```
kubectl port-forward svc/argocd-playground 8888:8888
```

Verify v1 in the browser
```
http://localhost:8888/
```

You should see
```
{"status":"ok"}
```

Update the application

![image](https://user-images.githubusercontent.com/35857179/78024972-d20df580-738b-11ea-9c5f-245c4277c2ff.png)

Build and push the docker image to docker hub. Then update the image tag to v2 in ``deployment.yaml``.

```yaml
- image: wingkwong/argocd-playground:v2
```

Go back to Argo CD UI, the status becomes ``OutofSync``.

![image](https://user-images.githubusercontent.com/35857179/78025927-688ee680-738d-11ea-81d9-4d6d5c05c34c.png)

Click ``SYNC``

A new pod is being created, while the original one is still here.

![image](https://user-images.githubusercontent.com/35857179/78025985-83f9f180-738d-11ea-9791-556a3d50851f.png)

Once it is ready, the original one will be deleted.
![image](https://user-images.githubusercontent.com/35857179/78026005-89efd280-738d-11ea-9804-abbd78b56358.png)

You should see the below error
```
E0331 20:24:00.727018   61938 portforward.go:400] an error occurred forwarding 8888 -> 8888: error forwarding port 8888 to pod 0f8b6902adcdbfdcde17a17bc1d182db8c4c849ba50ef369d90969e1349797b5, uid : failed to find sandbox "0f8b6902adcdbfdcde17a17bc1d182db8c4c849ba50ef369d90969e1349797b5" in store: does not exist
```

We should stop port forwarding before redeploying a different version. Let's kill it and do it again. 

```
kubectl port-forward svc/argocd-playground 8888:8888
```

Go to 
```
http://localhost:8888/
```

Now you can see the new changes
```
{"status":"ok", "message": "hello-world"}
```

# Clean up
```
k3d delete -n argocd-playground
```

# Compare with FluxCD
Argo CD allows users to sync in an application level instead of a repository level by setting the Path. It supports different templating such as kustomize, helm, ksonnet, jsonnet, etc. With an UI portal, users can simply manage the application there. However, it cannot monitor a docker repository and deploy from the repository. The docker image needs to be manually updated for each updates.  

# Useful links
- [Argo CD](https://argoproj.github.io/argo-cd/)
- [arkade](https://github.com/alexellis/arkade#get-arkade)
- [k3d](https://github.com/rancher/k3d)