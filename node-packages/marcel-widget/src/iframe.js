const decorate = iframe => {
  if (!iframe) return null
  const wiframe = iframe.contentWindow
  const diframe = iframe.contentDocument
  return {
    window: wiframe,
    document: diframe,
    iframe,
    addMessageListener: onMessage =>
      window.addEventListener('message', event => {
        if (event.source === wiframe) onMessage(event.data)
      }),
    postMessage: message => wiframe.postMessage(message, '*'),
  }
}

export default decorate
