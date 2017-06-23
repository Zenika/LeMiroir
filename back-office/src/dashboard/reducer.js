//@flow
import type { Reducer } from 'redux'
import { actions } from './actions'
import { mapValues } from 'lodash'
import { set, update, unset, chain } from 'immutadot'
import uuid from 'uuid/v4'
import type {
  DashboardAction,
  DashboardState,
  LayoutMap,
  PluginInstanceMap,
} from './type'

import mockedData from '../mocked-data/dashboards'
const intialState = mockedData

const updatePlugins = (layout: LayoutMap) => (plugins: PluginInstanceMap) => {
  return mapValues(plugins, plugin => {
    if (!layout[plugin.instanceId])
      throw new Error('Plugin instance not found in layout')

    const { x, y, w: columns, h: rows } = layout[plugin.instanceId]
    return { ...plugin, x, y, columns, rows }
  })
}

const dashboard: Reducer<DashboardState, DashboardAction> = (
  state = intialState,
  action,
) => {
  switch (action.type) {
    case actions.SELECT_PLUGIN: {
      return { ...state, selectedPlugin: action.payload.instanceId }
    }
    case actions.SELECT_DASHBOARD: {
      return { ...state, selectedDashboard: action.payload.dashboardId }
    }
    case actions.UNSELECT_DASHBOARD: {
      return { ...state, selectedDashboard: null }
    }
    case actions.DELETE_DASHBOARD: {
      return unset(state, `dashboards.${action.payload.dashboardId}`)
    }
    case actions.ADD_DASHBOARD: {
      const id = uuid()
      return chain(state)
        .set(`dashboards.${id}`, {
          id,
          name: 'Dashboard',
          description: '',
          cols: 20,
          rows: 20,
          ratio: 16 / 9,
          plugins: [],
        })
        .set('selectedDashboard', id)
        .value()
    }
    case actions.ADD_PLUGIN: {
      const instanceId = uuid()
      return set(state, `dashboard.plugins.${instanceId}`, {
        ...action.payload.plugin,
        x: 0,
        y: 0,
        columns: 1,
        rows: 1,
        instanceId,
      })
    }
    case actions.DELETE_PLUGIN: {
      return state.selectedPlugin
        ? unset(state, `dashboard.plugins.${state.selectedPlugin}`)
        : state
    }
    case actions.CHANGE_PROP: {
      const { instanceId, prop, value } = action.payload
      return set(
        state,
        `dashboard.plugins.${instanceId}.props.${prop.name}.value`,
        value,
      )
    }
    case actions.SAVE_LAYOUT: {
      const { layout } = action.payload
      return update(state, 'dashboard.plugins', updatePlugins(layout))
    }
    case actions.UPDATE_CONFIG: {
      const { property, value } = action.payload
      const parsedValue = !isNaN(value) ? parseFloat(value) : value
      return set(state, `dashboard.${property}`, parsedValue)
    }
    default:
      return state
  }
}

export default dashboard
