import * as config from '../config.json';


const request = async (url, options) => {
    return fetch(url, {
        ...options,
        headers: {
            ...config.defaultHeaders,
            // ...options.headers,
        }
    }).then(resp => {
        if (!resp.ok) {
            throw new Error(`got response ${resp.status} and with test ${resp.statusText}`)
        }
        return resp.json()
    })
        .then(resp => resp)
        .catch(error => {
            console.error("The error is", error)
            return error;
        })
}

export const useFetch = () => {

    const get = async (url, params = null, options = {}) => {
        console.log(url, config.defaultHeaders);

        if (!url) {
            throw new SyntaxError("Url is required but not given")
        }

        return request(url, {
            ...options,
            // params,
        })
    }

    const post = async (url, params = {}, options = {}) => {

        if (!url) {
            throw new SyntaxError("Url is required but not given")
        }

        return request(url, {
            ...options,
            method: 'POST',
            body: params,
        })
    }

    const put = async (url, params = {}, options = {}) => {

        if (!url) {
            throw new SyntaxError("Url is required but not given")
        }

        return request(url, {
            ...options,
            method: 'PUT',
            body: params,
        })
    }

    return {
        get,
        post,
        put,
    }
}