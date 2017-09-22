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

Download the `app` pod configuration file:

```
wget 
```
