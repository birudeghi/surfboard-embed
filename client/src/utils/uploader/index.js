import EventEmitter from 'events';
import get from 'lodash/get';
import avrdude from './avrgirl';

const asyncTimeout = (timeout) => new Promise(
    (resolve) => setTimeout(
        () => resolve(timeout), timeout
    )
);

class Uploader extends EventEmitter {
    constructor() {
        super();
        this.toolMap = {
            avrdude,
        };
    }

    async waitForClose(serial, count = 0) {
        if (!serial.serial.isOpen || count > 50) return null;
        await asyncTimeout(100);
        return this.waitForClose(serial, count + 1);
    }

    async upload(hex, serial, board, config) {
        const existBaud = serial.baud;
        
        const protocol = `${board?.props?.upload?.tool?.split(':').pop()}.${board?.props?.upload?.protocol}`;
        console.log("protocol: ", protocol)
        const uploader = get(this.toolMap, protocol);
        
        if (!uploader) throw new Error("Board not currently supported");

        await serial.setMute(true);
        await serial.disconnect();
        console.log("serial disconnected for flashing");
        
        await uploader(hex, board, serial.serial, {
            ...config
        });
        
        console.log("upload completed");
        await this.waitForClose(serial);
        await serial.setBaud(existBaud);
        await serial.connect();
        await serial.setMute(false);
    }
}

export default new Uploader();