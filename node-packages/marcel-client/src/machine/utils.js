import { reduce } from 'robot3'

export const isEmpty = (ctx, event) => !event.data

export const storeData = (key) => (ctx, { data }) => ({ ...ctx, [key]: data })

export const storeError = message => (ctx, { error }) => {
  console.error(message, error)
  return { ...ctx, error: message || error && error.message }
}

// eslint-disable-next-line no-unused-vars
export const storeDataWithoutError = (key) => reduce(({ error, ...ctx }, { data }) => {
  return ({ ...ctx, [key]: data })
})

// eslint-disable-next-line no-unused-vars
export const removeError = () => ({ error, ...ctx }) => ctx
