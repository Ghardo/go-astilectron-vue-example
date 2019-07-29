package strip

import (
	"encoding/hex"
	"log"
	"time"

	serial "go.bug.st/serial.v1"
	"go.bug.st/serial.v1/enumerator"
)

type StripRGBData struct {
	Start      byte
	End        byte
	Red        byte
	Green      byte
	Blue       byte
	Brightness byte
}

func (s StripRGBData) AsByteArray() []byte {
	data := make([]byte, 7)
	data[0] = 'c'
	data[1] = s.Start
	data[2] = s.End
	data[3] = s.Red
	data[4] = s.Green
	data[5] = s.Blue
	data[6] = s.Brightness

	return data
}

type DeviceData struct {
	Id           int    `json:"id"`
	Vendor       string `json:"vendor"`
	Product      string `json:"product"`
	SerialNumber string `json:"serialnumber"`
	Port         string `json:"port"`
	Baud         int    `json:"baud"`
}

type StripColorData struct {
	Color      string `json:"color"`
	Brightness int    `json:"brightness"`
	Power      bool   `json:"power"`
}

func FindConnectedDevices() (deviceList []DeviceData, err error) {
	var port serial.Port
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, nil
	}
	for _, comport := range ports {
		var device DeviceData

		if comport.IsUSB {
			device.Port = comport.Name
			device.Vendor = comport.VID
			device.Product = comport.PID
			device.SerialNumber = comport.SerialNumber
			device.Baud = 9600
			port, err = Open(device.Port, device.Baud)
			if err == nil {
				time.Sleep(2 * time.Second)
				port.SetDTR(true)
				port.Write([]byte("i"))
				response, _ := Recieve(port, 7)
				if validateResonse(response) {
					device.Id = int(response[2])
					deviceList = append(deviceList, device)
				}
			}
			port.Close()
		}
	}
	return deviceList, nil
}

func Open(port string, baud int) (resource serial.Port, err error) {
	mode := &serial.Mode{
		BaudRate: baud,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	resource, err = serial.Open(port, mode)
	if err != nil {
		return nil, err
	}

	if resource == nil {
		return nil, nil
	}

	return resource, nil
}

func Connect(port string, baud int) (resource serial.Port, err error) {
	resource, err = Open(port, baud)
	if err == nil {
		time.Sleep(2 * time.Second)
		resource.SetDTR(true)
		resource.Write([]byte("i"))
		response, _ := Recieve(resource, 7)
		if validateResonse(response) {
			return resource, nil
		}

		return nil, nil
	} else {
		return nil, err
	}
}

func Recieve(resource serial.Port, count int) (response []byte, err error) {
	buff := make([]byte, count)
	var c int
	var n int
	c = 0

	for {
		n, err = resource.Read(buff)

		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		response = append(response[:], buff[:n]...)

		c = c + n
		if c >= count {
			break
		}
	}

	return response, nil
}

func validateResonse(response []byte) bool {

	if response[0] != 43 {
		return false
	}
	if response[1] != 51 {
		return false
	}
	if response[2] <= 100 {
		return false
	}
	if response[3] != 1 {
		return false
	}
	if response[4] != 0 {
		return false
	}
	if response[5] != 0 {
		return false
	}

	return true
}

func Change(resource serial.Port, data StripColorData) error {
	var brightness byte
	if !data.Power {
		brightness = 0x00
	} else {
		brightness = byte(float64(data.Brightness) / 100.0 * 255.0)
	}

	color := hexColorToByteArray(data.Color)
	stripdata := StripRGBData{0x00, 0x2c, color[0], color[1], color[2], brightness}
	_, err := resource.Write(stripdata.AsByteArray())
	return err
}

func hexToByte(h string) byte {
	data, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return data[0]
}

func hexColorToByteArray(hex string) []byte {
	data := make([]byte, 3)
	data[0] = hexToByte(hex[1:3])
	data[1] = hexToByte(hex[3:5])
	data[2] = hexToByte(hex[5:7])
	return data
}
