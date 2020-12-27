#/bin/bash
function test_blocking {
    echo "Verify Enforcer blocks ${SERVICE_NAME}"
    result=`wget --spider -S ${SERVICE_NAME}/cmd.exe 2>&1 | grep "  HTTP/" | awk '{print $2}'`
    echo "$result"
    if [ $result == 403 ];then echo "Passed, HTTP Code $result" && exit 0;else echo "Failed, HTTP Code is $result" && exit 1;fi
}
echo "running as user: ${UID} for namespace ${NAMESPACE}"
test_blocking