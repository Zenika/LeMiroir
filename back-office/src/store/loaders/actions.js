//@flow
import { toastr } from 'react-redux-toastr'
import { set, map } from 'lodash/fp'
import { mapPluginsToDashboard } from '../../common/utils'
import { config, backend, userBackend } from '../../api'
import type { LoadingThunkAction } from './type'
import { pluginsSelector } from '../../plugins'

export const actions = {
  LOAD_INITIAL_STARTED: 'LOADERS/LOAD_INITIAL_STARTED',
  LOAD_INITIAL_FINISHED: 'LOADERS/LOAD_INITIAL_FINISHED',
  LOAD_CONFIG_STARTED: 'LOADERS/LOAD_CONFIG',
  LOAD_CONFIG_SUCCESSED: 'LOADERS/LOAD_CONFIG_SUCCESSED',
  LOAD_CONFIG_FAILED: 'LOADERS/LOAD_CONFIG_FAILED',
  LOAD_PLUGINS_STARTED: 'LOADERS/LOAD_PLUGINS_STARTED',
  LOAD_PLUGINS_SUCCESSED: 'LOADERS/LOAD_PLUGINS_SUCESSED',
  LOAD_PLUGINS_FAILED: 'LOADERS/LOAD_PLUGINS_FAILED',
  LOAD_DASHBOARDS_STARTED: 'LOADERS/LOAD_DASHBOARDS_STARTED',
  LOAD_DASHBOARDS_SUCCESSED: 'LOADERS/LOAD_DASHBOARDS_SUCCESSED',
  LOAD_DASHBOARDS_FAILED: 'LOADERS/LOAD_DASHBOARDS_FAILED',
  LOAD_CLIENTS_STARTED: 'LOADERS/LOAD_CLIENTS_STARTED',
  LOAD_CLIENTS_SUCCESSED: 'LOADERS/LOAD_CLIENTS_SUCCESSED',
  LOAD_CLIENTS_FAILED: 'LOADERS/LOAD_CLIENTS_FAILED',
  LOAD_USERS_STARTED: 'LOADERS/LOAD_USERS_STARTED',
  LOAD_USERS_SUCCESSED: 'LOADERS/LOAD_USERS_SUCCESSED',
  LOAD_USERS_FAILED: 'LOADERS/LOAD_USERS_FAILED',
}

const loadDispatcher = dispatch => (ressource, loadPromise) => {
  const RESSOURCE = ressource.toUpperCase()
  dispatch({ type: actions[`LOAD_${RESSOURCE}_STARTED`] })

  const onSuccess = result =>
    dispatch({
      type: actions[`LOAD_${RESSOURCE}_SUCCESSED`],
      payload: { [ressource]: result },
    })

  const onFail = error => {
    dispatch({
      type: actions[`LOAD_${RESSOURCE}_FAILED`],
      payload: { error },
    })
    console.error(`Error while loading ${ressource}`, error)
    throw new Error(`Error while fetching ${ressource}`)
  }

  return loadPromise.then(onSuccess).catch(onFail)
}

const getDashboards = availablePlugins =>
  backend
    .getAllDashboards()
    .then(map(mapPluginsToDashboard(availablePlugins)))
    .then(map(set('ratio', 16 / 9)))

export const loadInitData = (): LoadingThunkAction => (dispatch, getState) => {
  dispatch({ type: actions.LOAD_INITIAL_STARTED })

  const load = loadDispatcher(dispatch)
  load('plugins', backend.getAvailablePlugins())
    .then(() => load('dashboards', getDashboards(pluginsSelector(getState()))))
    .then(() => load('clients', backend.getClients()))
    .then(() => dispatch({ type: actions.LOAD_INITIAL_FINISHED }))
    .catch((error: Error) => toastr.error('Erreur lors du chargement', error.message))
}

export const loadDashboards = (): LoadingThunkAction => (dispatch, getState) => {
  const load = loadDispatcher(dispatch)
  load('dashboards', getDashboards(pluginsSelector(getState())))
}

export const loadPlugins = (): LoadingThunkAction => (dispatch, getState) => {
  const load = loadDispatcher(dispatch)
  load('plugins', backend.getAvailablePlugins())
}

export const loadConfig = (): LoadingThunkAction => (dispatch, getState) => {
  const load = loadDispatcher(dispatch)
  load('config', config.loadConfig())
}

export const loadClients = (): LoadingThunkAction => (dispatch, getState) => {
  const load = loadDispatcher(dispatch)
  load('clients', backend.getClients())
}

export const loadUsers = (): LoadingThunkAction => (dispatch, getState) => {
  const load = loadDispatcher(dispatch)
  load('users', userBackend.getAllUsers())
}
