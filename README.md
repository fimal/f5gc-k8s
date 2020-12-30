# Free5gc Stage3 for Kubernetes
The free5GC is an open-source project for 5th generation (5G) mobile core networks. The ultimate goal of this project is to implement the 3GPP Release 15 (R15) and Release 16 (R16) 5G core network (5GC).

This repository aims to install dockerized free5gc stage3 and ride on kubernetes.
The main target of this project to investigate 5G attack surface and integrate Radware KWAF Security solution

Following attacks demonstraited:
 - IMSI cracking - Getting hold of target’s IMSI is the first step of a set of attacks such as location sharing, MMS/SMS hijacking, battery depletion, etc. The attacker can “discover” portions of the target’s IMSI and then can brute force guess the whole value.
Data enriched from N1/N2 interfaces supplied towards SBI and is analyzed by Radware KWAF to detect the attack
 - Resource exhaustion – a UE generates a high rate of authorization requests saturating the AUSF.
SUCI data enriched from N1 supplied to SBI and tracked by Radware KWAF to detect and mitigate the attack.
 - Attacks on the SBI based control plane - Attacker can infiltrate the control plane and gathers additional data / breaks control plane.
Data is inspected by Radware KWAF and detected/mitigated using Positive and Negative Security.

gNB/UE simulators used from project: https://github.com/aligungr/UERANSIM and https://github.com/hhorai/gnbsim

NOTE: Can work at version `5.0.0-23-generic` to run UPF(GTP5G Module.)

![Image of 5G](https://github.com/fimal/f5gc-k8s/blob/main/5G.png)

### Installing and testing

## Download the related files.
```
$ git clone https://github.com/fimal/f5gc-k8s.git
$ cd f5gc-k8s
```

## Install K8S cluster kubeadm
```
$ sudo kubeadm init --pod-network-cidr=192.167.0.0/16 --service-cidr=10.96.0.0/16
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

## Install K8S cluster with Calico CNI
```
$ kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml		# installing calico cni
$ kubectl create -f https://docs.projectcalico.org/manifests/custom-resources.yaml
Note: Before creating this manifest, read its contents and make sure its settings are correct for your environment. For example, you may need to change the default IP pool CIDR to match your pod network CIDR.
```

## Install Multus
```
$ git clone https://github.com/intel/multus-cni.git && cd multus-cni
$ cat ./images/multus-daemonset.yml | kubectl apply -f -
```

## Build images.
```
$ cd images
$ make build-base	# building base image
$ make build-amf    # building amf image
```

## Install Radware KWAF
Download HELMs and images from Radware site.
https://www.radware.com/products/kubernetes-waf1 .
Follow instruction in Installation guide

### Build IC (IMSI Cracking deployment)
 ```
 $ cd images/f5gc-ic
 $ make build-container
 ```

## Run Manifest
```
$ cd manifests
$ kubectl apply -k f5gc-nrf # First SBI function should be DB and NRF
or
$ ./run apply
```

## Apply KWAF Configuration
```
kubectl apply kwaf/
```

## Optionally Build AUSF with HELM
* Run manifest section already include in script installing ausf-nf with kwaf
 ```
 $ helm install ausf -n f5gc free5gc-apps/ -f waas/values.yaml -f free5gc-apps/values.yaml -f free5gc-apps/applications/ausf.yaml --debug
 ```

## NAT rules - in order to route traffic to SBI and Internet
```
$ kubectl exec -it -n f5gc f5gc-upf-67db5cb59d-jnpfp -c f5gc-upf -- bash
$
$ iptables -t nat -A POSTROUTING -s 60.60.0.0/16 -d 10.96.0.0/12 -j SNAT --to-source {upf-pod-ip-address} #add rule to SBI NF
$ iptables -t nat -v -L POSTROUTING -n --line-number             # list all nat rules
$ iptables -t nat -D POSTROUTING 1                               # delete rule line 1
```

## Tuning MTU - each k8s node need to configure mtu and update mtu on tunneling interfaces (N3 - GTPU)
```
 $ ip link set dev ens224 mtu 9000
```
