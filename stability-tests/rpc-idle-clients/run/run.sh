#!/bin/bash
rm -rf /tmp/consensusd-temp

NUM_CLIENTS=128
consensusd --devnet --appdir=/tmp/consensusd-temp --profile=6061 --rpcmaxwebsockets=$NUM_CLIENTS &
CONSENSUSD_PID=$!
CONSENSUSD_KILLED=0
function killConsensusdIfNotKilled() {
  if [ $CONSENSUSD_KILLED -eq 0 ]; then
    kill $CONSENSUSD_PID
  fi
}
trap "killConsensusdIfNotKilled" EXIT

sleep 1

rpc-idle-clients --devnet --profile=7000 -n=$NUM_CLIENTS
TEST_EXIT_CODE=$?

kill $CONSENSUSD_PID

wait $CONSENSUSD_PID
CONSENSUSD_EXIT_CODE=$?
CONSENSUSD_KILLED=1

echo "Exit code: $TEST_EXIT_CODE"
echo "Consensusd exit code: $CONSENSUSD_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $CONSENSUSD_EXIT_CODE -eq 0 ]; then
  echo "rpc-idle-clients test: PASSED"
  exit 0
fi
echo "rpc-idle-clients test: FAILED"
exit 1
