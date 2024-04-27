#!/bin/bash
rm -rf /tmp/consensusd-temp

consensusd --simnet --appdir=/tmp/consensusd-temp --profile=6061 &
CONSENSUSD_PID=$!

sleep 1

orphans --simnet -alocalhost:42511 -n20 --profile=7000
TEST_EXIT_CODE=$?

kill $CONSENSUSD_PID

wait $CONSENSUSD_PID
CONSENSUSD_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Consensusd exit code: $CONSENSUSD_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $CONSENSUSD_EXIT_CODE -eq 0 ]; then
  echo "orphans test: PASSED"
  exit 0
fi
echo "orphans test: FAILED"
exit 1
