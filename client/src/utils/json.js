import { Type } from 'js-binary';

const embedSchema = new Type({
    compileInfo: {
        fqbn: 'string',
        files: [
            {
                content: 'string',
                name: 'string'
            }
        ],
        flags: {
            verbose: 'boolean',
            preferLocal: 'boolean'
        },
        libs: [
            {
                name: 'string',
                url: 'string',
                version: 'string'
            }
        ]
    },
    embedInfo: {
        selectedBoard: 'string',
        compatibleBoards: [
            {
                fqbn: 'string',
                name: 'string',
                version: 'string',
                propertiesId: 'string',
                baud: "number"
            }
        ],
        files: [
            {
                content: 'string',
                name: 'string'
            }
        ],
        libs: [
            {
                name: 'string',
                url: 'string',
                version: 'string'
            }
        ],
        appId: 'string'
    }
});

//no need account management; this can be stored anonymously
const serverSchema = new Type({
    appId: 'string',
    email: 'string',
    compatibleBoards: [
        {
            fqbn: 'string',
            name: 'string',
            
            version: 'string',
            propertiesId: 'string',
            baud: "number"
        }
    ],
    libs: [
        {
            name: 'string',
            url: 'string',
            version: 'string'
        }
    ],
    files: [
        {
            content: 'string',
            name: 'string'
        }
    ] 
});


