# UART Device Service
## Overview
This Device Service is a reference example of a Device Service developed with **[Go Device Service SDK](https://github.com/edgexfoundry/device-sdk-go)**.

- Function:
  - It is developed for universal serial device, such as USB to TTL serial port, rs232 interface and rs485 interface device.  It provides REST API interfaces to communicate with serial sensors or configure them.
- Characteristics:
  - Universal serial port
- Physical interface: TTL, RS232, RS485, USB to TTL
- Driver protocol: UART



## Usage
- This Device Service have to run with other EdgeX Core Services, such as Core Metadata, Core Data, and Core Command.
- After starting the service, the linux uart device will be opened and configured(defalut: /dev/ttyUSB0 with 115200bps). User can choose another uart device by sending "uartconfig" command.
- When uart device is opend,  the data sent from another sensors will be stored and they can be read by "gethex" or "getstring" command. User can also sends down link data to those sensors by "sendhex" or  "sendstring" command.



## REST API

| Method | Core Command | parameters                         | Description                                                  | Response                                       |
| ------ | ------------ | ---------------------------------- | ------------------------------------------------------------ | ---------------------------------------------- |
| get    | gethex       |                                    | Get serial device data, output in hex string format          | "{"rxbuf hex":"32333437383536"}"               |
| get    | getstring    |                                    | Get serial device data, output in ascii string format        | "{"rxbuf string":"aadd33"}"                    |
| put    | sendhex      | {"sendhex":<txbuf>}                | Send data to the serial device, input in hex string format   | 200 ok                                         |
| put    | sendstring   | {"sendstring":<txbuf>}             | Send data to the serial device, input in ascii string format | 200 ok                                         |
| get    | uartconfig   |                                    | Get serial port configuration parameters                     | "{"baud":"9600","device path":"/dev/ttyUSB5"}" |
| put    | uartconfig   | {"path":<path>,<br/>"baud":<baud>} | Configure serial port parameters                             | 200 ok                                         |



## License
[Apache-2.0](LICENSE)
