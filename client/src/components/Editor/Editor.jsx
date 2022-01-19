import React, { useEffect, useRef } from 'react';
import { EditorView } from '@codemirror/view';
import { EditorState } from '@codemirror/state';
import { cpp } from '@codemirror/lang-cpp';
import { basicSetup } from '@codemirror/basic-setup';
import './editor.scss';

const Editor = props => {
    const { code, pxHeight } = props;
    const editorRef = useRef();
    
    useEffect(() => {
        if (editorRef.current === null) return;

        let myTheme = EditorView.theme({
            "&": {
                color: "#93a1a1",
                backgroundColor: "#FDF6E3",
                maxHeight: `${pxHeight}px`
              },

              ".cm-content": {
                caretColor: "#FFFFFF",
                fontFamily: "Source Code Pro",
                fontSize: "13px",
                minHeight: `${pxHeight}px`
              },
              "&.cm-focused .cm-cursor": {
                borderLeftColor: "#002B36"
              },
              "&.cm-focused .cm-selectionBackground, ::selection": {
                backgroundColor: "#93a1a1"
              },
              ".cm-gutters": {
                backgroundColor: "#FDF6E3",
                color: "#93a1a1",
                minHeight: `${pxHeight}px`
              },

              ".cm-activeLine": {
                backgroundColor: "#EEE8D5"
              },

              ".cm-activeLineGutter": {
                backgroundColor: "transparent"
              },

              ".cm-scroller": {
                  overflowY: "auto"
              }
            }, {dark: false});
        
        const state = EditorState.create({
            doc: code,
            extensions: [basicSetup, cpp(), myTheme]
        });
        
        const myView = new EditorView({
            state: state,
            parent: editorRef.current
            })
        
        

        return () => {
            myView.destroy();
        };
    }, [code]);
    

    return (
        <div ref={editorRef} className="editor-section" />
    )
}

export default Editor;