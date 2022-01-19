import axios from 'axios';

const instance = axios.create({
    baseURL: "http://localhost:3030",
    timeout: 10000
})
/**
 * Calls GET /account-details to obtain accountUid, categoryUid, and savingsGoalUid.
 */
export const getAccountDetails = async () => {
    return await instance.get("/account-details")
        .then((response) => {
            return response.data;
        })
};


/**
 * Calls GET /round-up to obtain round up amount.
 * @param {String} accId, the accountUid
 * @param {String} catId, the categoryUid 
 */
export const roundUp = async (accId, catId) => {
    return await instance.get(`/round-up/account/${accId}/category/${catId}`)
        .then((response) => {
            return response.data.roundUpAmount;
        })
};


/**
 * Calls POST /transfer to transfer amount into savings goal.
 * @param {string} accId, the accountUid 
 * @param {string} saveId, the savingsGoalUid 
 */
export const transfer = async (amount, accId, saveId) => {
    return await instance.post(`/transfer/account/${accId}/goal/${saveId}`, {
        amount: amount
    })
        .then((response) => {
            return response.data;
        })
};