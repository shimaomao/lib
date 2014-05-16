#!/bin/bash
AWK_COLLECT_TEST_MODULES=$(mktemp /tmp/speedland-apps.awk.XXXXXX)
TEST_LOG=$(mktemp /tmp/speedland-apps.test.XXXXXX)
cat <<-EOS >> $AWK_COLLECT_TEST_MODULES
BEGIN {
  FS = "/";
}
{
  i=1
  dir=""
  while(i<NF)
  {
    dir = dir "" \$(i) "/"
    i++
  }
  print substr(dir, 1)
}
EOS
trap "Cleanup Test Environment..." 0 1 2 3 15
trap "rm $AWK_COLLECT_TEST_MODULES" 0 1 2 3 15
trap "rm $TEST_LOG" 0 1 2 3 15
trap "rm -fr /tmp/appengine-aetest*" 0 1 2 3 15
trap "rm -fr /tmp/speedland-apps.*" 0 1 2 3 15
trap "ps aux | grep '/tmp/tmp\*-go-bin/_go_app' | awk '{print $2}' | xargs kill 2>/dev/null"  0 1 2 3 15
trap "ps aux | grep appengine-aetest | grep -v grep | awk '{print $2}' | xargs kill 2>/dev/null"  0 1 2 3 15

SRCH_DIR="."
if [[ "$TARGET_DIR" != "" ]]; then
  SRCH_DIR=$TARGET_DIR
fi
test_dir=$(find $SRCH_DIR -name '*_test.go' | awk -f $AWK_COLLECT_TEST_MODULES | uniq)

EXIT_CODE=0
TEST_COUNT=0
for t in ${test_dir}
do
  c=$(grep -R "func Test" ${t} | wc -l)
  TEST_COUNT=$(expr $TEST_COUNT + $c)
  goapp test -v ${t} ${TEST_FLAGS} >> ${TEST_LOG}
  if [ $? -ne 0 ]; then
    EXIT_CODE=1
  fi
done

echo ===================== Test Result =====================
cat ${TEST_LOG}
echo ===================== Test Stats =====================
PASSED=$(cat ${TEST_LOG} | grep -- "^---" | grep PASS | wc -l | awk '{print $1}')
FAILED=$(cat ${TEST_LOG} | grep -- "^---" | grep FAIL | wc -l | awk '{print $1}')
SKIPPED=$(expr $TEST_COUNT - $PASSED - $FAILED)
echo -e " PASSED:\t${PASSED}/${TEST_COUNT}"
echo -e " FAILED:\t${FAILED}/${TEST_COUNT}"
echo -e "SKIPPED:\t${SKIPPED}/${TEST_COUNT}"
exit $EXIT_CODE

# echo ===================== Test Result =====================
# cat $(TEST_LOG)
# make cleanproc
# echo ===================== Test Stats =====================
# echo "PASS: $$(cat $(TEST_LOG) | grep -- "^---" | grep PASS | wc -l)"
# echo "FAIL: $$(cat $(TEST_LOG) | grep -- "^---" | grep FAIL | wc -l)"
# make cleantmp
