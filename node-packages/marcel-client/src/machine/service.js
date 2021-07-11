import { interpret } from 'robot3'
import { writable } from 'svelte/store'
import machine, { initialContext } from './machine'

const { subscribe, set } = writable({ state: machine.current, context: initialContext, machine })

const service = interpret(
  machine,
  ({ machine, context }) => {
    // eslint-disable-next-line no-console
    if (IS_DEV_MODE) console.debug('changed state', machine.current, context)
    set({ state: machine.current, context, machine })
  }
)

const send = (...args) => {
  // eslint-disable-next-line no-console
  if (IS_DEV_MODE) console.debug('event sent', ...args)
  return service.send(...args)
}

export default { subscribe, send }
