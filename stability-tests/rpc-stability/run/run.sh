#!/bin/bash
rm -rf /tmp/consensusd-temp

consensusd --devnet --appdir=/tmp/consensusd-temp --profile=6061 --loglevel=debug &
CONSENSUSD_PID=$!

sleep 1

rpc-stability --devnet -p commands.json --profile=7000
TEST_EXIT_CODE=$?

kill $CONSENSUSD_PID

wait $CONSENSUSD_PID
CONSENSUSD_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Consensusd exit code: $CONSENSUSD_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $CONSENSUSD_EXIT_CODE -eq 0 ]; then
  echo "rpc-stability test: PASSED"
  exit 0
fi
echo "rpc-stability test: FAILED"
exit 1
