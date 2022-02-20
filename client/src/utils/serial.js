import SerialPort from 'avrgirl-arduino/lib/browser-serialport';
import EventEmitter from 'events';

class SurfSerial extends EventEmitter {
    constructor() {
        super();
        this.baud = 9600
        this.connected = false
        this.port = null //stub for browser-serialport
        this.serial = null

        this._handleData = (buff) => {
            this.emit('data', buff);
        }
    }

    setBaud(baud) {
        this.baud = baud;
    }

    setupSerial() {
        this.serial = new SerialPort(this.port, {
            baudRate: this.baud,
            autoOpen: false
        });
        this.serial.on('data', this._handleData);
    }

    openSerial() {
        return new Promise((resolve, reject) => {
            this.serial.open(function (error) {
                if (error) {
                    console.log("Failed to open: " + error);
                    reject("Failed to open: " + error);
                } else {
                    console.log("Port open.")
                    resolve(true);
                }
            });
        })
    }
    
    readSerial() {
        return new Promise((resolve, reject) => {
            this.serial.on('data', function(data) {
                console.log("data received: ", data);
                resolve(data);
            })
        })
    }

    closeSerial() {
        if (this.serial.isOpen) {
            this.serial.close(function (error) {
                if (error) {
                    console.log("Failed to close: " + error)
                    return false;
                } else {
                    if (!this.serial.isOpen) {
                        return true;
                    }
                }
            })
        }
    }
    
}

export default SurfSerial;