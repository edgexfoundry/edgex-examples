# GPIO Device Service
## Overview
This Device Service is a reference example of a Device Service developed with **[Go Device Service SDK](https://github.com/edgexfoundry/device-sdk-go)**.

- Function:
  - It is developed to control the system gpio through the following functions. For example, export or unexport the gpio, set gpio direction to output or input mode, set the gpio output status and get gpio status.
- Physical interface: system gpio (/sys/class/gpio)
- Driver protocol: IO



## Usage
- This Device Service have to run with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command.
- After starting the service, user can use "exportgpio" command to export a system gpio which comes from the path "/sys/class/gpio". Then configure the gpio by sending "gpiodirection" or "gpiovalue" command, and check it's configuration by getting them.



## REST API

| Method | Core Command  | parameters                | Description                                                  | Response                             |
| ------ | ------------- | ------------------------- | ------------------------------------------------------------ | ------------------------------------ |
| put    | exportgpio    | {"export":<gpionum>}      | Export a gpio from "/sys/class/gpio"<br><gpionum>: int, gpio number | 200 ok                               |
| put    | unexportgpio  | {"unexport":<gpionum>}    | Export a gpio from "/sys/class/gpio"<br/><gpionum>: int, gpio number | 200 ok                               |
| put    | gpiodirection | {"direction":<direction>} | Set direction for the exported gpio<br/><direction>: string, "in" or "out" | 200 ok                               |
| get    | gpiodirection |                           | Get direction of the exported gpio                           | "{\"direction\":\"in\",\"gpio\":65}" |
| put    | gpiovalue     | {"value":<value>}         | Set value for the exported gpio<br/><value>: int, 1 or 0     | 200 ok                               |
| get    | gpiovalue     |                           | Get value of the exported gpio                               | "{\"gpio\":65,\"value\":1}"          |





## License
[Apache-2.0](LICENSE)

