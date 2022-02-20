import React from 'react';
import Surfboard from './components/Surfboard';
import SurfSerial from './utils/SerialImpl';

const serial = new SurfSerial();

function App() {
  return (
      <div className="App">
        <Surfboard appId="6205267583c793f23ee6b706" serial={serial}/>
      </div>
  );
}

export default App;
