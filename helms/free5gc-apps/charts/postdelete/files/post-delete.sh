#!/usr/bin/env bash
set -e
echo "WAAS Sample App Postdelete starting"
echo "Print kubectl client and server versions"
kubectl version
function delete_secure_configuration {
    #securebeat
    CONFIG_INFIX=$1
    set +e
    kubectl delete -n ${NAMESPACE} configmap ${NAMES_PREFIX}${CONFIG_INFIX}-ca-config-${APP_NAME}
    kubectl delete -n ${NAMESPACE} secret ${NAMES_PREFIX}${CONFIG_INFIX}-client-secret-${APP_NAME}
    set -e
}
delete_secure_configuration securebeat
echo "WAAS Sample App Postdelete complete"