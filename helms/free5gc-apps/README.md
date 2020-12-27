# Project Title
This is helm for free5gc app secured by Radware KWAF
## Getting Started
git clone https://github.com/fimal/f5gc-k8s.git

### HELM TEMPLATE GENERATOR FOR Free5Gc
```
helm template f5gc-ausf -n f5gc free5gc-apps/ -f  waas/values.yaml -f free5gc-apps/applications/ausf.yaml --set global.kwafns=kwaf --output-dir helm3/hacks/target
```

### HELM INSTALL EXAMPLE FOR Free5Gc
### Note: please specify correct NAMESPACE where KWAF controller deployed - helm install "--set" or values.yaml
```
helm install f5gc-ausf -n test helm3/sample-app/ -f helm3/waas/values.yaml -f helm3/sample-app/applications/ausf.yaml --set global.kwafns=kwaf --create-namespace --debug --wait

```

## HELM TEST
```
helm test <<HELM-RELEASE-NAME>> -n test
```

## HELM DELETE
```
helm uninstall <<HELM-RELEASE-NAME>> -n test
```
