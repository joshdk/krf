package krf.joshdk.github.com

matched if {
    input.kind == "Service"
    input.metadata.name == "my-service"
}

matched if {
    input.kind == "Pod"
    input.metadata.name == "test-pod"
}
