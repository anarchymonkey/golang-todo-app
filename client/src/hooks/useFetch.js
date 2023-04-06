import { useCallback, useState } from "react"
import * as config from "../config.json"


const useFetch = () => {
	const [loading, setLoading] = useState(false)
	const [error, setError] = useState(null)

	const request = async (url, options) => {
		setLoading(true)
		try {
			const resp = await fetch(url, {
				...options,
				headers: {
					...config.defaultHeaders,
					// ...options.headers,
				}
			})

			if (!resp.ok) {
				throw new Error(`got response ${resp.status} and with test ${resp.statusText}`)
			}
            
			return resp.json()
		} catch (err) {
			setError(err)
		} finally {
			setLoading(false)
		}
	}

	const get = useCallback(async (url, params = null, options = {}) => {
		console.log(url, config.defaultHeaders)

		if (!url) {
			throw new SyntaxError("Url is required but not given")
		}

		return request(url, {
			...options,
			params,
		})
	}, [request])

	const post = useCallback(async (url, params = {}, options = {}) => {

		if (!url) {
			throw new SyntaxError("Url is required but not given")
		}

		return request(url, {
			...options,
			method: "POST",
			body: JSON.stringify(params),
		})
	}, [request])

	const put = useCallback(async (url, params = {}, options = {}) => {

		if (!url) {
			throw new SyntaxError("Url is required but not given")
		}
		console.log(params)

		return request(url, {
			...options,
			method: "PUT",
			body: JSON.stringify(params),
		})
	}, [request])

	const deleteRequest = useCallback(async (url, params = {}, options = {}) => {

		if (!url) {
			throw new SyntaxError("Url is required but not given")
		}
		console.log(params)

		return request(url, {
			...options,
			method: "DELETE",
			body: JSON.stringify(params),
		})
	}, [request])

	return {
		loading,
		error,
		get,
		post,
		put,
		deleteRequest,
	}
}

export default useFetch