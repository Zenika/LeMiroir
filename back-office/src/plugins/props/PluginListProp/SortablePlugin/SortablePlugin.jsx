//@flow
import React from 'react'
import { SortableHandle } from 'react-sortable-hoc'
import ListItem from 'react-toolbox/lib/list/ListItem'
import FontIcon from 'react-toolbox/lib/font_icon/FontIcon'
import type { Plugin } from '../../../type'

import './SortablePlugin.css'

const DragHandle = SortableHandle(() =>
  <FontIcon value="menu" className="grab" />,
)

class SortablePlugin extends React.Component {
  props: {
    plugin: Plugin,
    onDelete: Plugin => void,
  }

  onDelete = () => {
    this.props.onDelete(this.props.plugin)
  }

  render() {
    return (
      <ListItem
        caption={this.props.plugin.name}
        ripple={false}
        leftIcon={<DragHandle />}
        rightIcon={
          <FontIcon
            value="delete"
            style={{ color: 'red', cursor: 'pointer' }}
            onClick={this.onDelete}
          />
        }
      />
    )
  }
}

export default SortablePlugin
