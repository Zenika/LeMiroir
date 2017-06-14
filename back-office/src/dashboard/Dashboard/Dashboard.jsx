//@flow
import React from 'react'
import Grid from '../Grid'
import type { Dashboard as DashboardT } from '../type'
import './Dashboard.css'

const Dashboard = ({ dashboard }: { dashboard: DashboardT }) => {
  const layout = dashboard.plugins.map(
    ({ x, y, columns, rows, ...instance }) => ({
      layout: { x, y, h: rows, w: columns },
      plugin: instance,
    }),
  )

  return <Grid ratio={2} rows={20} cols={20} layout={layout} />
}

export default Dashboard
