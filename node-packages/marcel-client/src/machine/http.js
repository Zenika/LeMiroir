import { invoke } from 'robot3'
import * as http from '../utils/http'

const makeUrl = ({ config: { apiURI } }, url) => {
  const base = apiURI.endsWith('/') ? apiURI.slice(0, apiURI.length - 1) : apiURI
  return `${base}${url}`
}

const makeRequest = (method) => (url, options) => ctx => (
  http.request(makeUrl(ctx, url), { ...method, options })
)

const makeRequestWithBody = (method) => (url, body, options) => ctx => (
  http.requestWithBody(makeUrl(ctx, url), body, { method, ...options })
)

export const get = makeRequest('GET')
export const post = makeRequestWithBody('POST')
export const put = makeRequestWithBody('PUT')

const makeHttpAction = action => (ctx, event) => {
  const httpRequest = action(ctx, event)
  if (!httpRequest) throw new Error('The first parameter of invokeHttp should return the result of an http function (get, post or put)')
  return httpRequest(ctx)
}

export const invokeHttp = (action, ...transitions) => invoke(
  makeHttpAction(action),
  ...transitions
)

export const isUnauthorized = (ctx, { error }) => http.isUnauthorized(error)
