//@flow
import { range } from 'lodash'

export default range(20).map(i => ({
  name: `Plugin ${i}`,
  elementName: `plugin-${i}`,
  icon: 'picture_in_picture_alt',
  props: {
    prop1: { name: 'prop1', type: 'string', value: 'hello world !' },
    prop2: { name: 'prop2', type: 'number', value: 42 },
    prop3: { name: 'prop3', type: 'boolean', value: true },
  },
}))
