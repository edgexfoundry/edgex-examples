# Example of security-enabled EdgeX Remote Device Service

This is an example implementation of an EdgeX remote device service with EdgeX security features enabled.

This example uses an SSH tunnel to provide network security between two EdgeX nodes,
and uses the delayed start services feature to obtain a secret store token at runtime
in order to make authenticated connections to Consul and the EdgeX secret store.

## Build and Run

Please see the detailed steps in "how-to guide" for remote device services in documentation repository here: <https://github.com/edgexfoundry/edgex-docs/blob/main/docs_src/security/Ch-RemoteDeviceServices.md>.
