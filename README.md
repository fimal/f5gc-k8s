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

### Installing and testing

* Download the related files.

```
$ git clone https://github.com/fimal/f5gc-k8s.git
$ cd f5gc-k8s
```

* install K8S cluster with Calico CNI and MULTUS
```
$ kubectl apply -f clusters/cni/calico/ 		# installing calico cni
$ kubectl apply -f clusters/cni/multus/			# installing multus
```

* Build images.

```
$ cd images
$ make build-base	# building base image
$ make build-amf    # building amf image
$ make build-container

```
 * Build IC (IMSI Cracking)
 * NAT rules

```
$ iptables -t nat -A POSTROUTING -s 60.60.0.0/16 -d 10.98.0.0/16 -j SNAT --to-source 172.20.104.54 #add rule
$ iptables -t nat -v -L POSTROUTING -n --line-number             # list all nat rules
$ iptables -t nat -D POSTROUTING 1                               # delete rule line 1

```
 * Tuning MTU - k8s node configure mtu

```
 $ p link set dev ens224 mtu 9000
```
