const proxy = require('http-proxy-middleware')

module.exports = function (app) {
  app.use(proxy('/api', { target: 'http://localhost:8090/', ws: true }))
  app.use(proxy('/uris', { target: 'http://localhost:8090/' }))
}