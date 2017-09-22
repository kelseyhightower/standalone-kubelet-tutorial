# Standalone Kubelet Tutorial

This tutorial will guide you through running the Kubernetes [Kubelet](https://kubernetes.io/docs/admin/kubelet/) in standalone mode on [Container Linux](https://coreos.com/why). You will also deploy an application using a [static pod](https://kubernetes.io/docs/tasks/administer-cluster/static-pod/), test it, then upgrade the application.

## Compute Instance

Create the `standalone-kubelet` compute instance:

```
gcloud compute instances create standalone-kubelet \
  --async \
  --boot-disk-size 200 \
  --can-ip-forward \
  --image-family coreos-stable \
  --image-project coreos-cloud \
  --machine-type n1-standard-1 \
  --tags standalone-kubelet
```

Allow HTTP traffic to the `standalone-kubelet` instance:

```
gcloud compute firewall-rules create allow-standalone-kubelet \
  --allow tcp:80 \
  --target-tags standalone-kubelet
```

## Install the Standalone Kubelet

SSH into the `standalone-kubelet` compute instance:

```
gcloud compute ssh standalone-kubelet
```

Download the Kubelet systemd unit file:

```
wget -q --show-progress --https-only --timestamping \
  https://raw.githubusercontent.com/kelseyhightower/standalone-kubelet-tutorial/master/kubelet.service
```

Move the `kubelet.service` unit file to the systemd configuration directory:

```
sudo mv kubelet.service /etc/systemd/system/
```

Start the `kubelet` service:

```
sudo systemctl daemon-reload
```

```
sudo systemctl enable kubelet
```

```
sudo systemctl start kubelet
```

### Verification

It will take a few minutes for `kubelet` container to download and initialize. Verify the `kubelet` is running:

```
sudo systemctl status kubelet
```

## Static Pods

In this section you will deploy an application that responds to HTTP requests with its running config and version. The application configuration will be initialized before the HTTP service starts by using an [init container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/). Once the pod has started the application configuration will be updated every 30 seconds by a configuration sidecar.

SSH into the `standalone-kubelet` compute instance:

```
gcloud compute ssh standalone-kubelet
```

Verify no container are running:

```
sudo docker ps
```

Verify no container images are installed:

```
sudo docker images
```

Create the kubelet manifests directory:

```
sudo mkdir -p /etc/kubernetes/manifests
```

Download the `app-v0.1.0.yaml` pod manifest:

```
wget -q --show-progress --https-only --timestamping \
  https://raw.githubusercontent.com/kelseyhightower/standalone-kubelet-tutorial/master/pods/app-v0.1.0.yaml
```

Move the `app-v0.1.0.yaml` pod manifest to the kubelet manifest directory:

```
sudo mv app-v0.1.0.yaml /etc/kubernetes/manifests/app.yaml
```

> Notice the `app-v0.1.0.yaml` pod manifest is being renamed to `app.yaml`. This prevents our application from being deployed twice. Each pod must have a unique `metadata.name`.

List the installed container images:

```
sudo docker images
```
```
REPOSITORY                             TAG                 IMAGE ID            CREATED             SIZE
gcr.io/hightowerlabs/app               0.1.0               c7d7002a0776        23 minutes ago      6.325 MB
gcr.io/hightowerlabs/configurator      0.1.0               164e54187008        4 hours ago         2.346 MB
gcr.io/google_containers/pause-amd64   3.0                 99e59f495ffa        16 months ago       746.9 kB
```

List the running containers:

```
docker ps
```

> You should see three containers running which represent the `app` pod. Docker does not understand pods so the containers are listed as individual containers following the Kubernetes naming convention. 

At this point the `app` pod is up and running on port 80 in the host namespace.

```
curl http://127.0.0.1
```

```
version: 0.1.0
hostname: standalone-kubelet.c.hightowerlabs.internal
key: 1506065606
```

Wait about 30 seconds and make another HTTP request:

```
curl http://127.0.0.1
```

```
version: 0.1.0
hostname: standalone-kubelet.c.hightowerlabs.internal
key: 1506065636
```

> Notice the `key` field has changed.

### Testing Remote Access

The `app` pod is listening on `0.0.0.0:80` in the host network and is accessible via the external IP of the `standalone-kubelet` compute instance.

Get the external IP of the `standalone-kubelet` instance:

```
EXTERNAL_IP=$(gcloud compute instances describe standalone-kubelet \
  --format 'value(networkInterfaces[0].accessConfigs[0].natIP)')
```

Make and HTTP request to using the external IP:

```
curl http://${EXTERNAL_IP}
```

```
version: 0.1.0
hostname: standalone-kubelet.c.hightowerlabs.internal
key: 1506065696
```

## Updating Static Pods

```
gcloud compute ssh standalone-kubelet
```

Download the `app-v0.2.0.yaml` pod manifest:

```
wget -q --show-progress --https-only --timestamping \
  https://raw.githubusercontent.com/kelseyhightower/standalone-kubelet-tutorial/master/pods/app-v0.2.0.yaml
```

Move the `app-v0.2.0.yaml` pod manifest to the kubelet manifest directory:

```
sudo mv app-v0.2.0.yaml /etc/kubernetes/manifests/app.yaml
```

> Notice the `app-v0.2.0.yaml` is being renamed to `app.yaml`. This overwrites the current pod manifest and will force the kubelet upgrade the `app` pod.

List the installed container images:

```
docker images
```
```
REPOSITORY                             TAG                 IMAGE ID            CREATED             SIZE
gcr.io/hightowerlabs/app               0.1.0               c7d7002a0776        30 minutes ago      6.325 MB
gcr.io/hightowerlabs/app               0.2.0               3028d4a68eb1        31 minutes ago      6.325 MB
gcr.io/hightowerlabs/configurator      0.1.0               164e54187008        4 hours ago         2.346 MB
gcr.io/google_containers/pause-amd64   3.0                 99e59f495ffa        16 months ago       746.9 kB
```

> Notice the `gcr.io/hightowerlabs/app:0.2.0` image has been added to the local repository.

At this point `app` version `0.2.0` is up and running.

```
curl http://127.0.0.1
```

```
version: 0.2.0
hostname: standalone-kubelet.c.hightowerlabs.internal
key: 1506065786
```

## Cleanup

```
gcloud -q compute instances delete standalone-kubelet
```

```
gcloud -q compute firewall-rules delete allow-standalone-kubelet
```
