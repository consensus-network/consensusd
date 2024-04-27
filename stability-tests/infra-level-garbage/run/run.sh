#!/bin/bash
rm -rf /tmp/consensusd-temp

consensusd --devnet --appdir=/tmp/consensusd-temp --profile=6061 &
CONSENSUSD_PID=$!

sleep 1

infra-level-garbage --devnet -alocalhost:42611 -m messages.dat --profile=7000
TEST_EXIT_CODE=$?

kill $CONSENSUSD_PID

wait $CONSENSUSD_PID
CONSENSUSD_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Consensusd exit code: $CONSENSUSD_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $CONSENSUSD_EXIT_CODE -eq 0 ]; then
  echo "infra-level-garbage test: PASSED"
  exit 0
fi
echo "infra-level-garbage test: FAILED"
exit 1
