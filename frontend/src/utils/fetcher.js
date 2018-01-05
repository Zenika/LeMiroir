export const fetcher = baseUrl => {
  const request = (url, options) =>
    fetch(baseUrl + url, { ...options, credentials: 'include' })
      .then(response => {
        if (~~(response.status / 100) !== 2) throw response
        return response
      })
      .then(res => res.json())

  const requestWithBody = (url, body, options = {}) =>
    request(url, {
      headers: body ? { 'Content-Type': 'application/json' } : {},
      body: body ? JSON.stringify(body) : null,
      ...options,
    })

  return {
    get: url => request(url),
    post: (url, body) => requestWithBody(url, body, { method: 'POST' }),
    put: (url, body) => requestWithBody(url, body, { method: 'PUT' }),
    del: url => request(url, { method: 'DELETE' }),
  }
}

let backendFetcherInstance
export const backendFetcher = config => {
  if (!backendFetcherInstance) {
    const baseUrl = `${config.ssl ? 's' : ''}://${config.urls.backend}`
    backendFetcherInstance = fetcher('http' + baseUrl)
    backendFetcherInstance.ws = clientId => new WebSocket(`ws${baseUrl}/clients/${clientId}/ws`)
  }
  return backendFetcherInstance
}

export const localFetcher = fetcher('./')
