#!/bin/bash
rm -rf /tmp/consensusd-temp

consensusd --devnet --appdir=/tmp/consensusd-temp --profile=6061 --loglevel=debug &
CONSENSUSD_PID=$!
CONSENSUSD_KILLED=0
function killConsensusdIfNotKilled() {
    if [ $CONSENSUSD_KILLED -eq 0 ]; then
      kill $CONSENSUSD_PID
    fi
}
trap "killConsensusdIfNotKilled" EXIT

sleep 1

application-level-garbage --devnet -alocalhost:42611 -b blocks.dat --profile=7000
TEST_EXIT_CODE=$?

kill $CONSENSUSD_PID

wait $CONSENSUSD_PID
CONSENSUSD_KILLED=1
CONSENSUSD_EXIT_CODE=$?

echo "Exit code: $TEST_EXIT_CODE"
echo "Consensusd exit code: $CONSENSUSD_EXIT_CODE"

if [ $TEST_EXIT_CODE -eq 0 ] && [ $CONSENSUSD_EXIT_CODE -eq 0 ]; then
  echo "application-level-garbage test: PASSED"
  exit 0
fi
echo "application-level-garbage test: FAILED"
exit 1
