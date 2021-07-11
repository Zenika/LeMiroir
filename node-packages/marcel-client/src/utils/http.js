const isRequestError = status => error => (error.status || error.response && error.response.status) === status
export const isUnauthorized = isRequestError(401)
export const isNotFound = isRequestError(404)

export const isJson = response => response.headers.get('Content-Type').includes('application/json')

export const queryParams = () => {
  const params = new URL(location.href).searchParams
  return new Proxy({}, {
    get: (_, prop) => params.get(prop)
  })
}

export const request = async (url, options) => {
  const response = await fetch(url, {
    ...options,
    credentials: 'include',
  })

  if (!response.ok) {
    const error = new Error('Fetch error')
    error.response = response
    throw error
  }

  if (isJson(response)) return response.json()

  return response
}

export const requestWithBody = async (url, body, options) => request(url, {
  body: body && JSON.stringify(body),
  ...options,
  headers: {
    'Content-Type': 'application/json',
    ...options.headers
  },
})
