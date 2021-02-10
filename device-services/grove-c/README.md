# grove-c
This example contains the instructions and scripts to connect the Grove Sensor Kit to a Raspberry Pi 3 and use EdgeX to get data from the kit's sensors and command the kit's devices.

This example uses the EdgeX Geneva release.

This repository contains:
1. GrovePIStarterKit.pdf - Guide that describe installation, setup and use of Grove Sensors on Raspberry Pi 3.
2. docker-compose-demo-grove.yml - compose file to download Edgex Core and device grove service images based on arm64 for Raspberry PI.
3. nodered_flow.json - Node-RED flow available to read measurements from Grove Sensors.
   In the current version, readings from SoundSensor, LightIntensity and RotarySensor are available.

