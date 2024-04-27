#!/bin/bash

APPDIR=/tmp/consensusd-temp
CONSENSUSD_RPC_PORT=29587

rm -rf "${APPDIR}"

consensusd --simnet --appdir="${APPDIR}" --rpclisten=0.0.0.0:"${CONSENSUSD_RPC_PORT}" --profile=6061 &
CONSENSUSD_PID=$!

sleep 1

RUN_STABILITY_TESTS=true go test ../ -v -timeout 86400s -- --rpc-address=127.0.0.1:"${CONSENSUSD_RPC_PORT}" --profile=7000
TEST_EXIT_CODE=$?

kill $CONSENSUSD_PID

wait $CONSENSUSD_PID
CONSENSUSD_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Consensusd exit code: $CONSENSUSD_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $CONSENSUSD_EXIT_CODE -eq 0 ]; then
  echo "mempool-limits test: PASSED"
  exit 0
fi
echo "mempool-limits test: FAILED"
exit 1
