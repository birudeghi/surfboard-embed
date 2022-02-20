import React, { useState, useMemo } from 'react';
import { oneDark } from './oneDark';
import { cpp } from '@codemirror/lang-cpp';
import { basicSetup } from '@codemirror/basic-setup';
import CodeMirror from 'rodemirror';
import './editor.scss';

const Editor = props => {
    const { code, onUpdate } = props;

    const extensions = useMemo(
      () => [basicSetup, oneDark, cpp()],
      [],
    )

    const handleValue = code => {
      onUpdate(code);
    }
    
    return (
      <CodeMirror
        value={code}
        onUpdate={(v) => {
          if (v.docChanged) {
            handleValue(v.state.doc.toString())
          }
        }}
        extensions={extensions}
      />
    )
    
 
}

export default Editor;

