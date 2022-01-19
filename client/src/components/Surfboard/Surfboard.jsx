import React, { useState } from 'react';
import Editor from '../Editor';  
import './surfboard.scss';

const Surfboard = props => {
    const { link } = props;

    const [boards, setBoards] = useState([]);
    const [selectBoard, setSelectBoard] = useState({});
    const [code, setCode] = useState("");

    return (
        <div className="surfboard-container" id='surfboard-container' style={{height:600, width:500}}>
            <Editor code="/* Arduino code goes here */" pxHeight={400} />
        </div>
    )
}

export default Surfboard;