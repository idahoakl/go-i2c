// Package i2c provides low level control over the linux i2c bus.
//
// Before usage you should load the i2c-dev kenel module
//
//      sudo modprobe i2c-dev
//
// Each i2c bus can address 127 independent i2c devices, and most
// linux systems contain several buses.
package i2c

import (
	"fmt"
	"os"
	"syscall"
	"sync"
)

const (
	I2C_SLAVE = 0x0703
)

// I2C represents a connection to an i2c device.
type I2C struct {
	rc *os.File
	mtx *sync.Mutex
}

// New opens a connection to an i2c device.
func NewI2C(bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}

	this := &I2C{
		rc: f,
		mtx: &sync.Mutex{},
	}
	return this, nil
}

func (this *I2C) setAddress(addr uint8) error {
	return ioctl(this.rc.Fd(), I2C_SLAVE, uintptr(addr))
}

// Write sends buf to the remote i2c device. The interpretation of
// the message is implementation dependant.
func (this *I2C) Write(addr uint8, buf []byte) (int, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	return this.writeNoSync(addr, buf)
}

func (this *I2C) writeNoSync(addr uint8, buf []byte) (int, error) {
	if e := this.setAddress(addr); e != nil {
		return nil, e
	}

	return this.rc.Write(buf)
}

func (this *I2C) WriteByte(addr uint8, b byte) (int, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	return this.writeByteNoSync(addr, b)
}

func (this *I2C) writeByteNoSync(addr uint8, b byte) (int, error) {
	var buf [1]byte
	buf[0] = b

	return this.writeNoSync(addr, buf[:])
}


func (this *I2C) Read(addr uint8, p []byte) (int, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	return this.readNoSync(addr, p)
}

func (this *I2C) readNoSync(addr uint8, p []byte) (int, error) {
	if e := this.setAddress(addr); e != nil {
		return nil, e
	}

	return this.rc.Read(p)
}


func (this *I2C) Close() error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	return this.rc.Close()
}

// SMBus (System Management Bus) protocol over I2C.
// Read byte from i2c device register specified in reg.
func (this *I2C) ReadRegU8(addr uint8, reg byte) (byte, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	_, err := this.writeNoSync(addr, []byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 1)
	_, err = this.readNoSync(addr, buf)
	if err != nil {
		return 0, err
	}
	log.Debug("Read U8 %d from reg 0x%0X", buf[0], reg)
	return buf[0], nil
}

// SMBus (System Management Bus) protocol over I2C.
// Write byte to i2c device register specified in reg.
func (this *I2C) WriteRegU8(addr uint8, reg byte, value byte) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	buf := []byte{reg, value}
	_, err := this.writeNoSync(addr, buf)
	if err != nil {
		return err
	}
	log.Debug("Write U8 %d to reg 0x%0X", value, reg)
	return nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read unsigned big endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegU16BE(addr uint8, reg byte) (uint16, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	_, err := this.writeNoSync(addr, []byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = this.readNoSync(addr, buf)
	if err != nil {
		return 0, err
	}
	w := uint16(buf[0])<<8 + uint16(buf[1])
	log.Debug("Read U16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read unsigned little endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegU16LE(addr uint8, reg byte) (uint16, error) {
	w, err := this.ReadRegU16BE(addr, reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read signed big endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegS16BE(addr uint8, reg byte) (int16, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	_, err := this.writeNoSync(addr, []byte{reg})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 2)
	_, err = this.readNoSync(addr, buf)
	if err != nil {
		return 0, err
	}
	w := int16(buf[0])<<8 + int16(buf[1])
	log.Debug("Read S16 %d from reg 0x%0X", w, reg)
	return w, nil
}

// SMBus (System Management Bus) protocol over I2C.
// Read unsigned little endian word (16 bits) from i2c device
// starting from address specified in reg.
func (this *I2C) ReadRegS16LE(addr uint8, reg byte) (int16, error) {
	w, err := this.ReadRegS16BE(addr, reg)
	if err != nil {
		return 0, err
	}
	// exchange bytes
	w = (w&0xFF)<<8 + w>>8
	return w, nil

}

// SMBus (System Management Bus) protocol over I2C.
// Write unsigned big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegU16BE(addr uint8, reg byte, value uint16) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	buf := []byte{reg, byte((value & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := this.writeNoSync(addr, buf)
	if err != nil {
		return err
	}
	log.Debug("Write U16 %d to reg 0x%0X", value, reg)
	return nil
}

// SMBus (System Management Bus) protocol over I2C.
// Write unsigned big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegU16LE(addr uint8, reg byte, value uint16) error {
	w := (value*0xFF00)>>8 + value<<8
	return this.WriteRegU16BE(addr, reg, w)
}

// SMBus (System Management Bus) protocol over I2C.
// Write signed big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegS16BE(addr uint8, reg byte, value int16) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	buf := []byte{reg, byte((uint16(value) & 0xFF00) >> 8), byte(value & 0xFF)}
	_, err := this.writeNoSync(addr, buf)
	if err != nil {
		return err
	}
	log.Debug("Write S16 %d to reg 0x%0X", value, reg)
	return nil
}

// SMBus (System Management Bus) protocol over I2C.
// Write signed big endian word (16 bits) value to i2c device
// starting from address specified in reg.
func (this *I2C) WriteRegS16LE(addr uint8, reg byte, value int16) error {
	w := int16((uint16(value)*0xFF00)>>8) + value<<8
	return this.WriteRegS16BE(addr, reg, w)
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
