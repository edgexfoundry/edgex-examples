package driver

import (
	"os"
	"io/ioutil"
	"errors"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

type GPIODev struct {
	lc     logger.LoggingClient
	gpio int
}

func NewGPIODev(lc logger.LoggingClient) *GPIODev {
	return &GPIODev{lc: lc, gpio: -1}
}


func (dev *GPIODev) ExportGPIO(gpio int) error {
	err := exportgpio(gpio)
	if err == nil {
		dev.gpio = gpio
	}
	return err
}

func (dev *GPIODev) UnexportGPIO(gpio int) error {
	err := unexportgpio(gpio)
	if err == nil {
		dev.gpio = -1
	}
	return err
}

func (dev *GPIODev) SetDirection(direction string) error {
	if dev.gpio == -1 {
		return errors.New("Please export gpio first")
	}
	return setgpiodirection(dev.gpio, direction)
}

func (dev *GPIODev) GetDirection() (string,error) {
	if dev.gpio == -1 {
		return "", errors.New("Please export gpio first")
	}
	direction, err := getgpiodirection(dev.gpio)
	if err != nil {
		return "", err
	}  else {
		res, _ := json.Marshal(map[string]interface{}{"gpio": dev.gpio, "direction": direction})
		return string(res), err
	}
}

func (dev *GPIODev) SetGPIO(value int) error {
	if dev.gpio == -1 {
		return errors.New("Please export gpio first")
	}
	direction, err := getgpiodirection(dev.gpio)
	if err != nil {
		return err
	}
	if strings.Contains(direction, "in")  {
		return errors.New("Can not set the gpio which is input state")
	}

	return setgpiovalue(dev.gpio, value)
}

func (dev *GPIODev) GetGPIO() ( string, error ) {

	if dev.gpio == -1 {
		return "", errors.New("Please export gpio first")
	}
	gpiovalue, err := getgpiovalue(dev.gpio)
	if err != nil {
		return "", err
	}  else {
		res, _ := json.Marshal(map[string]interface{}{"gpio": dev.gpio, "value": gpiovalue})
		return string(res), err
	}
}


func exportgpio(gpioNum int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	} else {
		return ioutil.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
	}
}

func unexportgpio(gpioNum int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return ioutil.WriteFile("/sys/class/gpio/unexport", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
	} else {
		return nil
	}
}

func setgpiodirection(gpioNum int, direction string) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var way string
		if direction == "in" {
			way = "in"
		} else {
			way = "out"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum), []byte(way), 0644)
	} else {
		return errors.New("Please export gpio first")
	}
}

func getgpiodirection(gpioNum int) (string, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		direction, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum))
		if err != nil {
			return "", err
		} else {
			return strings.Replace(string(direction), "\n", "", -1), err
		}
	} else {
		return "", errors.New("Please export gpio first")
	}
}

func setgpiovalue(gpioNum int, value int) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var tmp string
		if value == 0 {
			tmp = "0"
		} else {
			tmp = "1"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum), []byte(tmp), 0644)
	} else {
		return errors.New("Please export gpio first")
	}
}

func getgpiovalue(gpioNum int) (int, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ret, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum))
		if err != nil {
			return 0, err
		} else {
			value, _ := strconv.Atoi(strings.Replace(string(ret), "\n", "", -1))
			return value, err
		}
	} else {
		return -1, errors.New("Please export gpio first")
	}
}
