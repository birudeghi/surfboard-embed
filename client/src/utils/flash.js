import { compile } from './clientApi';
import Uploader from './uploader';

/**
 * flashInfo data structure:
 * {
 *     board,
 *     libs,
 *     files,
 *     flags    
 * }
 * @param {*} serial 
 * @param {*} flashInfo 
 * @returns 
 */
export const flash = async (serial, flashInfo) => {
    let res;
    let hex;

    try {
        res = await compile(flashInfo);
        hex = res.hex;
        await Uploader.upload(hex, serial, flashInfo.board, flashInfo.flags);

    } catch (err) {
        console.error(err);
        throw new Error("Flashing failed: Code couldn't be flashed due to error.");
    }

    return res.log;
}