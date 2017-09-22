# Standalone Kubelet Tutorial

This tutorial will guide you through running the Kubernetes Kubelet in standalone mode.

## Compute Instance

Create the `kubelet` compute instance:

```
gcloud compute instances create kubelet \
  --async \
  --boot-disk-size 200 \
  --can-ip-forward \
  --image-family coreos-stable \
  --image-project coreos-cloud \
  --machine-type n1-standard-2
```

## Install a Standalone Kubelet

SSH into the `kubelet` compute instance:

```
gcloud compute ssh kubelet
```

Create the Kubelet systemd unit file:

```
cat > kubelet.service <<EOF
[Service]
Environment=KUBELET_IMAGE_TAG=v1.7.6_coreos.0
Environment="RKT_RUN_ARGS=--uuid-file-save=/var/run/kubelet-pod.uuid \
  --volume=resolv,kind=host,source=/etc/resolv.conf \
  --mount volume=resolv,target=/etc/resolv.conf"
ExecStartPre=-/usr/bin/rkt rm --uuid-file=/var/run/kubelet-pod.uuid
ExecStart=/usr/lib/coreos/kubelet-wrapper \
  --allow-privileged \
  --file-check-frequency 30s \
  --max-pods 10 \
  --minimum-image-ttl-duration 300s \
  --pod-manifest-path=/etc/kubernetes/manifests \
  --sync-frequency 30s
ExecStop=-/usr/bin/rkt stop --uuid-file=/var/run/kubelet-pod.uuid
Restart=always
RestartSec=10s

[Install]
WantedBy=multi-user.target
EOF
```

Move the `kubelet.service` unit file to the system configuration directory:

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

Verify the `kubelet` is running:

```
sudo systemctl status kubelet
```

## Static Pods

In this section you will run an example application that responds to HTTP request with its running config and version.

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

> Notice the `app-v0.1.0.yaml` is being renamed to `app.yaml`. This prevents our application from being deployed twice. Each pod must have a unique `metadata.name`.

List the installed images:

```
sudo docker images
```
```
REPOSITORY                             TAG                 IMAGE ID            CREATED             SIZE
gcr.io/hightowerlabs/app               0.1.0               8444c1627aa1        9 minutes ago       6.325 MB
gcr.io/hightowerlabs/configurator      0.1.0               164e54187008        About an hour ago   2.346 MB
gcr.io/google_containers/pause-amd64   3.0                 99e59f495ffa        16 months ago       746.9 kB
```

List the running containers:

```
docker ps
```

At this point the `app` pod is up and running on port 80 in the host namespace.

```
curl http://127.0.0.1
```

### Testing Remote Access

```
KUBELET_EXTERNAL_IP=$(gcloud compute instances describe kubelet \
  --format 'value(networkInterfaces[0].accessConfigs[0].natIP)')
```

```
curl http://${KUBELET_EXTERNAL_IP}
```

## Updating Static Pods

```
gcloud compute ssh kubelet
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

> Notice the `app-v0.2.0.yaml` is being renamed to `app.yaml`. This overwrites the current pod manifest and will force the kubelet upgrade the app pod.
