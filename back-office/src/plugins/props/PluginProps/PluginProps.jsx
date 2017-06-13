// @flow
import React from 'react'

import { SearchField } from '../../../common'
import PluginProp from '../PluginProp'
import type { Plugin } from '../../type'

import './PluginProps.css'

const PluginProps = (props: {
  plugin?: Plugin,
  filter: string,
  changeFilter: string => void,
}) => {
  const { plugin, filter, changeFilter } = props

  if (!plugin) return <div className="PluginsProps" />

  const { name, props: pluginProps } = plugin

  return (
    <div className="PluginProps">
      <h2>{name}</h2>
      <SearchField label="Search Prop" value={filter} onChange={changeFilter} />
      {pluginProps.map(p => <PluginProp prop={p} key={p.name} />)}
    </div>
  )
}

export default PluginProps
