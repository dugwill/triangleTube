package triangleTube

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type zone struct {
	Name   string   `json:"name"`
	Status bool     `json:"status"`
	Pin    rpio.Pin `json:"gpioPin"`
}

func initGpio() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Zone monitor initialized")
}

// AddZone creates a new zone and adds it to the zone map
// the zone map is created if necessary
func (t *TriangleTube) AddZone(name string, pin uint8) (err error) {
	if t.ZoneMap == nil {
		t.ZoneMap = make(map[string]*zone)
	}

	if _, ok := t.ZoneMap[name]; ok {
		return fmt.Errorf("zone exists")
	}

	z := zone{
		Name:   name,
		Status: false,
		Pin:    rpio.Pin(pin),
	}

	t.ZoneMap[name] = &z

	go t.monitorZone(name)

	return err
}

// RemoveZone removes the name zone from the zone map
func (t *TriangleTube) RemoveZone(name string) (err error) {
	if t.ZoneMap != nil {
		delete(t.ZoneMap, name)
	}
	return nil
}

// UpdateZone removes the name zone from the zone map
func (t *TriangleTube) UpdateZone(name string) (err error) {
	if t.ZoneMap != nil {
		delete(t.ZoneMap, name)
	}
	return nil
}

func (t *TriangleTube) monitorZone(n string) {
	if t.ZoneMap[n].Pin == 0 {
		fmt.Println("Zone input not set: ", t.ZoneMap[n].Name)
		return
	}
	for {
		if t.ZoneMap[n].Pin.Read() == rpio.Low {
			fmt.Printf("%s On: \n", t.ZoneMap[n].Name)
			t.ZoneMap[n].Status = true
		} else {
			t.ZoneMap[n].Status = false

		}
		time.Sleep(time.Second)
	}
}
