import { transition, reduce } from 'robot3'

export const isEmpty = (ctx, event) => !event.data

export const storeData = (key) => (ctx, { data }) => ({ ...ctx, [key]: data })

export const storeError = message => (ctx, { error }) => ({ ...ctx, error: message || error && error.message })

// eslint-disable-next-line no-unused-vars
export const storeDataWithoutError = (key) => ({ error, ...ctx }, { data }) => ({ ...ctx, [key]: data })

// eslint-disable-next-line no-unused-vars
export const removeError = () => ({ error, ...ctx }) => ctx

export const transitionWithData = (type, state, contextKey, ...args) => (
  transition(type, state, reduce(storeDataWithoutError(contextKey)), ...args)
)