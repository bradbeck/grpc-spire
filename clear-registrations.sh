#!/bin/bash

# Clear all workload registrations.

set -euo pipefail

function fetch_entries() {
  kubectl exec -n spire deployment/spire-server -- \
    /opt/spire/bin/spire-server entry show
}

function delete_entry() {
  kubectl exec -n spire deployment/spire-server -- \
    /opt/spire/bin/spire-server entry delete -entryID "$1"
}

echo "Fetching existing SPIRE entries"
_entries=$(fetch_entries | grep 'Entry ID' | sed 's/^Entry ID.*: //')
echo "$_entries"

for _entry in $_entries; do
  echo "Deleting: $_entry"
  delete_entry "$_entry"
done

echo "Done!"
