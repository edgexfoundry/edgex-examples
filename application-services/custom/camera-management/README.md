# Camera Management Example App Service
Camera management example Edgex application service to auto discover and connect to nearby ONVIF based cameras, 
control cameras via commands, create inference pipelines for the camera video streams and publish inference
results to MQTT broker.

This app uses [EdgeX Core Services][edgex-core-services], [Edgex Onvif device service][device-onvif-camera] and [Edge Video Analytics Microservice][evam].

## Steps for running this example:
We expect you'll be running in a relatively modern Linux environment
with `docker`, `docker-compose`, and `make` installed.

1. Get the [EdgeX Core Services][edgex-core-services] and Edgex device ONVIF service running by referring to docs from
   [Edgex Onvif device service][device-onvif-camera].
2. Get [Edge Video Analytics Microservice][evam] running for inference.
```shell
# Run this once to download edge-video-analytics into the edge-video-analytics sub-folder, 
# download models, and patch pipelines
make install-edge-video-analytics

# Run the EVAM services (in another terminal)
make run-edge-video-analytics
# ...
# Leave this running. If needed to stop
make stop-edge-video-analytics
```
3. Configure Camera Credentials

   Option 1: Modify the `configuration.toml` file
   ```toml
   [Writable.InsecureSecrets.CameraCredentials]
   path = "CameraCredentials"
     [Writable.InsecureSecrets.CameraCredentials.Secrets]
     username = ""
     password = ""
   ```
   
   Option 2: Export environment variable overrides
   ```shell
   export WRITABLE_INSECURESECRETS_CAMERACREDENTIALS_SECRETS_USERNAME=<username>
   export WRITABLE_INSECURESECRETS_CAMERACREDENTIALS_SECRETS_PASSWORD=<passowrd>
   ```

4. Build and run the example application service. Web UI is used to view cameras, select models 
   and start inference pipelines for camera video streams and also view inference results streams.
```shell
# Build the app. 
make build-app

# Run the app.
make run-app
# ...
# Open your browser to http://localhost:59750
# ...
# Ctrl-C to stop it
```

### Development and Testing of UI
```shell
# Build the production web-ui into the web-ui/dist folder
# This is what is served by the app service on port 59750
make web-ui

# Serve the Web-UI in hot-reload mode on port 4200
make serve-ui
# ...
# Open your browser to http://localhost:4200
# ...
# Ctrl-C to stop it
```

[edgex-core-services]: https://github.com/edgexfoundry/edgex-go
[device-onvif-camera]: https://github.com/edgexfoundry-holding/device-onvif-camera
[evam]: https://www.intel.com/content/www/us/en/developer/articles/technical/video-analytics-service.html
