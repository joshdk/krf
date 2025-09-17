#!/bin/sh

apiversion="$1"
kind="$2"
name="$3"
namespace="$4"

if [ "$RESOURCE_APIVERSION" != "$apiversion" ]; then
	exit 1
elif [ "$RESOURCE_KIND" != "$kind" ]; then
	exit 1
elif [ "$RESOURCE_NAME" != "$name" ]; then
	exit 1
elif [ "$RESOURCE_NAMESPACE" != "$namespace" ]; then
	exit 1
else
	exit 0
fi

#desired="${1:-1}"

#exec jq --exit-status --arg desired $desired '.spec.replicas == ($desired | tonumber)'
