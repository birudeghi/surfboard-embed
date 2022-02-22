import React from 'react';
import Surfboard from './components/Surfboard';
import SurfSerial from './utils/SerialImpl';

const serial = new SurfSerial();

function App(props) {
  return (
      <div className="App">
        <Surfboard appId={props.appid} serial={serial}/>
      </div>
  );
}

export default App;
