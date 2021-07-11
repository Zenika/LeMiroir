import { request, queryParams } from '../utils/http'
import * as http from './http'

export const loadConfig = () => request('./config.json')

export const refreshLogin = () => http.post('/auth/login')

export const login = (ctx, { login, password }) => http.post('/auth/login', { login, password })

export const createClient = (ctx, { name, mediaID }) => {
  const params = queryParams()
  return http.post('/clients/', { name: name || params.name, mediaID: Number(mediaID || params.mediaID) })
}

export const loadClient = () => http.get(`/clients/${localStorage.clientID}/`)

export const updateClient = ({ client }) => {
  const { name, mediaID } = queryParams()
  return http.put('/clients/', { ...client, name, mediaID: Number(mediaID) })
}

export const loadMedia = ({ client }) => http.get(`/medias/${client.mediaID}/`)

export async function connectClient({ client, config: { websocketURI } }) {
  const base = websocketURI.endsWith('/') ? websocketURI : `${websocketURI}/`
  const connection = new WebSocket(`${base}clients/${client.id}/ws`)

  const close = () => connection.close()

  return new Promise((resolve, reject) => {
    try {
      connection.addEventListener('open', () => resolve({ close }))
      connection.addEventListener('close', () => this.send('clientConnectionClosed'))
      connection.addEventListener('error', error => this.send({ type: 'clientConnectionError', error }))
      connection.addEventListener('message', message => {
        if (message.data === 'update') this.send('clientUpdated')
      })
    } catch (error) {
      reject(error)
    }
  })
}
