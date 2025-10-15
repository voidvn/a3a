#!/bin/sh
# wait-for-it.sh

set -e

# Get host and port from first argument
hostport=$1
host=$(echo $hostport | cut -d: -f1)
port=$(echo $hostport | cut -d: -f2 -s)

# Default port if not specified
if [ -z "$port" ]; then
    port=5432
fi

shift
cmd="/app/s4s-backend"

# Use netcat (nc) to check if the port is open
until nc -z "$host" "$port"; do
  >&2 echo "Postgres is unavailable at $host:$port - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd
