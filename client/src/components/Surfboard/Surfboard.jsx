import React, { useEffect, useState } from 'react';
import { Button, SIZE, SHAPE, KIND } from 'baseui/button';
import { useStyletron } from 'baseui';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { ResizableBox } from 'react-resizable';

import Editor from '../Editor';  
import Configuration from "../Configuration";
import { getInfo } from '../../utils/clientApi';
import { convertBoardApi2Client } from '../../utils/common';
import { flash } from "../../utils/flash";
import './surfboard.scss';
import 'xterm/css/xterm.css';
import _ from 'lodash';

const term = new Terminal();
const fitAddon = new FitAddon();
term.loadAddon(fitAddon);
fitAddon.fit();

const Surfboard = props => {
    const { appId, serial } = props;

    const [css] = useStyletron();

    const [selectBoard, setSelectBoard] = useState(null);
    const [boards, setBoards] = useState([]);
    const [libs, setLibs] = useState([]);
    const [files, setFiles] = useState("");
    const [fileName, setfileName] = useState("");
    const [errorMsg, setErrorMsg] = useState("");
    // const [disableFlash, setDisableFlash] = useState(true);
    const [isConnected, setIsConnected] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [isFlashed, setIsFlashed] = useState(false);
    const [termText, setTermText] = useState("");

    const handleTerminal = async () => {
        term.open(document.getElementById('surfboard-terminal'));
        await serial.serial.on('data', function(data) {
            term.write(data.toString().replace(/\r\n/g, '\n').replace(/\n/g, '\r\n'));
        });
    }

    useEffect(() => {
        let isSubscribed = true;

        const init = async () => {
            const info = await getInfo(appId);

            if (isSubscribed) {
                setFiles(info.files[0].content);
                setfileName(info.files[0].name);
                setLibs(info.libs);
                setBoards(info.compatibleBoards);
                setIsConnected(serial.connected);
            }

            await handleTerminal();
        };

        init();

        return () => isSubscribed = false;

    }, [appId]);

    const handleSelect = index => {
        const selectedBoard = boards[index];
        const flashBoard = convertBoardApi2Client(selectedBoard);
        setSelectBoard(flashBoard);
    }

    const handleFlash = async () => {
        setIsLoading(true);
        console.log("isConnected: ", isConnected);

        const flashInfo = {};
        flashInfo.board = selectBoard;
        flashInfo.libs = libs;
        flashInfo.files = [{
            content: files,
            name: `${fileName.split(".")[0]}/${fileName.split(".")[0]}.ino`
        }];
        flashInfo.flags = {
            verbose: false,
            preferLocal: false
        };

        if (isFlashed === false) {
            try {
                const log = await flash(serial, flashInfo);
                console.log('log: ', log);
                setIsFlashed(true);

            } catch (err) {
                setIsLoading(false);
                setErrorMsg(err.message);
                return;
            }

            // can also add handleReject value function like the above
        }
    }

    const handleConnect = async () => {
        setIsLoading(true);
        if (!isConnected) {
            try {
                await serial.requestDevice();
                //TODO change to serCurrentDevice for confirmation
                // Conclusion: There is supposed to be only one requestPort being invoked, but there was two.
                //isOpen is false and this.currentDevice is true when the connect() function is invoked.
                //isOpen only changes when BroswerSerialPort is instantiated. 

                // there is a way for duinoapp-client to allow avrgirl use existing serialport so that it soedn't have to requestPort again.
                
                setIsConnected(true);

            } catch (err) { 
                setErrorMsg(err.message);
                console.log('requestDevice rejected');
                setIsLoading(false);
                return;
            }
        }
        
        setErrorMsg("");
        setIsLoading(false);
    }

    const handleUpdate = code => {
        setFiles(code);
    }

    /*TODO include loading state */
    return (
        <div className="surfboard-container" id='surfboard-container' style={{height:600, width:500}}>
            <div className="surfboard-config">
                <Editor code={files} onUpdate={handleUpdate} />
            </div>
            <div className="surfboard-interactive">
                <div className="surfboard-connect-container">
                    <Configuration boards={boards} onSelect={handleSelect} disabled={!isConnected} />
                    <div className="surfboard-flash">
                        {errorMsg 
                        ? (
                            <p className="surfboard-error-msg">{errorMsg}</p>
                        )
                        : null       
                        }
                        {files ? (
                            isConnected
                            ? (
                                <Button 
                                    size={SIZE.compact}
                                    shape={SHAPE.pill}
                                    onClick={handleFlash}
                                    isLoading={isLoading}
                                    className={css({
                                        fontWeight: 600,
                                        ":enabled:hover": {background: "#1F70E9"},
                                    })}
                                    disabled={!selectBoard}
                                >
                                    Flash
                                </Button>
                            )
                            : (
                                <Button
                                    size={SIZE.compact}
                                    shape={SHAPE.pill}
                                    kind={KIND.secondary}
                                    onClick={handleConnect}
                                    isLoading={isLoading}
                                >
                                    Connect
                                </Button>
                            )
                        ) : null}
                    </div>
                </div>
                <ResizableBox className="surfboard-terminal-container" height={200} width={500}>
                    <div className="surfboard-terminal" id="surfboard-terminal" />
                </ResizableBox>
               
            </div>
        </div>
    )
}

export default Surfboard;