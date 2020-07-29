#!/bin/sh
# Start by running grpc health check
grpc_health_probe -addr=:${GRPC_PORT:-8181} || exit 1

# If data for making test RPC calls has been added to the container, run those tests too.
find healthchecks -mindepth 3 -maxdepth 3 | while read DIR; do
	METHOD="`basename "$DIR"`"
	REST="`dirname "$DIR"`"
	SERVICE="`basename "$REST"`"
	REST="`dirname "$REST"`"
	TESTNAME="`basename "$REST"`"
	TEMPFILE=`mktemp`
	echo "grpcurl --plaintext -d @ localhost:${GRPC_PORT:-8181} $SERVICE/$METHOD < $DIR/input.json > $TEMPFILE"
	if ! grpcurl --plaintext -d @ localhost:${GRPC_PORT:-8181} $SERVICE/$METHOD < $DIR/input.json > $TEMPFILE; then
		echo Test $DIR failed with error $?
		exit 1
	fi
	if ! cmp -s "$DIR/output.json" "$TEMPFILE"; then
		echo Test $DIR returned wrong result.
		echo -------------------------------------------------
		echo Expected result
		echo -------------------------------------------------
		cat $DIR/output.json

		echo -------------------------------------------------
		echo Actual result
		echo -------------------------------------------------
		cat $TEMPFILE
		exit 1
	fi
	rm $TEMPFILE
done
