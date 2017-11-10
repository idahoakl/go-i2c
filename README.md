# go-i2c
Implementation of i2c bus written in Golang.

## Utility usage
i2c-test utility is available to allow for directly interacting with devices on the i2c bus.

Bus and device address values are in base 10.

#### Write
Usage: `i2c-test <bus> <device addr> <data>`

The data should be a string of hex encoded bytes.  The leading '0x' should be omitted.

Example: `i2c-test 1 64 f3`

Output: No output on success, error message on error.

#### Read
Usage: `i2c-test <bus> <device addr> <num bytes to read>`

Read data will be output in hex format.

Example: `i2c-test 1 64 3`

Output: `0x6164DE`