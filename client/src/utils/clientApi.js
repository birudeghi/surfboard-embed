import axios from 'axios';

const appInstance = axios.create({
    baseURL: "http://localhost:8080",
    timeout: 10000
})

const duinoInstance = axios.create({
    baseURL: "http://localhost:3030",
    timeout: 10000
})

/**
 * GET call to URL provided that returns the code to the embed.
 */
export const getInfo = async (appId) => {
    return await appInstance.get(`/app/info/${appId}`)
        .then((response) => {
            return response.data;
        })
        .catch((error) => {
            return error.response;
        })
};

/**
 * Calls POST /compile to compile sketch.
 * @param {String} appId, the appId for the embed
 */
export const compile = async (compileInfo) => {

    const req = {
        fqbn: compileInfo.board.fqbn,
        files: compileInfo.files,
        libs: compileInfo.libs,
        flags: {
            verbose: false,
            preferLocal: false
        }
    }
    return await duinoInstance.post("/v3/compile", req)
        .then((response) => {
            return response.data;
        })
};