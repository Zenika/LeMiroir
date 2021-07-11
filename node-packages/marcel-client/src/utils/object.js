export const withoutNil = obj => {
  if (!obj || typeof obj !== 'object') return obj
  return Object.fromEntries(Object.entries(obj).filter(([, value]) => value != null))
}
