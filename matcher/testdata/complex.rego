package krf.joshdk.github.com

included if {
    input.kind == "Service"
    input.metadata.name == "my-service"
}

excluded if {
    input.spec.ports[_].port == 8080
}

excluded if {
    input.spec.ports[_].port == 8443
}

included if {
    input.kind == "Pod"
    input.metadata.name == "test-pod"
}

matched if {
    included
    not excluded
}
