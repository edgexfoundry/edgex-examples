package driver

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/tarm/serial"
)

type UartDev struct {
	lc     logger.LoggingClient
	config *serial.Config
	port *serial.Port
	rxbuf []byte
	rxlen int
}

func NewUartDev(devicePath string, baud int, lc logger.LoggingClient) *UartDev {
	config := &serial.Config{
		Name:        devicePath,
		Baud:        baud,
		Size:        8,
		StopBits:    serial.Stop1,
		Parity:      'N',
		ReadTimeout: 10 * time.Millisecond,
	}

	return &UartDev{lc: lc, config: config}
}


func (dev *UartDev) Listen() error {
	var err error
	dev.port, err = serial.OpenPort(dev.config)
	if err != nil {
		fmt.Printf("serial open port %s fail\n", dev.config.Name)
		return err
	}
	defer dev.port.Close()
	fmt.Printf("serial open port %s ok\n", dev.config.Name)

	dev.rxlen = 0
	dev.rxbuf = make([]byte, 1024)
	tmpbuff := make([]byte, 1024)
	for {
		n, err := dev.port.Read(tmpbuff)
		if err == nil && n > 0 && dev.rxlen < 1024 {
			for i := 0; i < n; i++ {
				dev.rxbuf[dev.rxlen+i] = tmpbuff[i]
			}
			dev.rxlen = dev.rxlen + n
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (dev *UartDev) GetHex() (string, error) {
	fmt.Printf("GetHex\n")
	res, _ := json.Marshal(map[string]interface{}{"rxbuf hex": hex.EncodeToString(dev.rxbuf[:dev.rxlen])})
	dev.rxlen = 0
	return string(res), nil
}

func (dev *UartDev) GetString() (string, error) {
	fmt.Printf("GetString\n")
	res, _ := json.Marshal(map[string]interface{}{"rxbuf string": string(dev.rxbuf[:dev.rxlen])})
	dev.rxlen = 0
	return string(res), nil
}

func (dev *UartDev) SendHex(txbuf []byte) error {
	fmt.Printf("SendHex\n")
	_, err := dev.port.Write(txbuf)
	return err
}

func (dev *UartDev) SendString(txbuf string) error {
	fmt.Printf("SendString\n")
	_, err := dev.port.Write( []byte(txbuf))
	return err
}

func (dev *UartDev) GetUartConfig() (string, error) {
	fmt.Printf("GetUartConfig\n")
	res, _ := json.Marshal(map[string]interface{}{"device path": dev.config.Name, "baud": strconv.Itoa(dev.config.Baud)})
	return string(res), nil
}

func (dev *UartDev) UartConfig(path string, baud int) error {
	var err error
	fmt.Printf("UartConfig\n")

	dev.config.Name = path
	dev.config.Baud = baud

	dev.port.Close()
	time.Sleep(10 * time.Millisecond)
	dev.port, err = serial.OpenPort(dev.config)
	if err != nil {
		fmt.Printf("serial open port %s fail\n", dev.config.Name)
		return err
	}
	return nil
}
