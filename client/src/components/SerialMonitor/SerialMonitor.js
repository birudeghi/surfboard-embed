import { Terminal } from 'xterm';

//TODO create a class just like serialPort
const terminal = new Terminal();
terminal.open(document.getElementById("surfboard-xterm"));

// TODO newline
export const terminalWrite = string => {
    terminal.write(string.replace(/\r\n/g, '\n').replace(/\n/g, '\r\n'));
}