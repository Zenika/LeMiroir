#!/usr/bin/env node

const fs = require('fs')
const path = require('path')
const { execSync } = require('child_process')
const { getPluginInfo, fatal } = require('./utils')
const files = require('./files')
const commands = require('./commands')

if (process.argv.length < 3) fatal('A plugin name should be given : $ yarn create marcel-plugin my-plugin')

const plugin = getPluginInfo(process.argv[2])

// Initialize directory
if (fs.existsSync(plugin.path)) fatal(`Can't create plugin : directory ${plugin.path} already exists`)
fs.mkdirSync(plugin.path)

// Write files
files.forEach((file => {
  const filePath = path.resolve(plugin.path, file.path)
  fs.mkdirSync(path.dirname(filePath), { recursive: true })
  fs.writeFileSync(filePath, file.content(plugin))
}))

// Run commands
commands.forEach(({ command, cwd }) => {
  execSync(
    command(plugin),
    { cwd: cwd ? path.resolve(plugin.path, cwd) : plugin.path }
  )
})

console.info(`Your plugin has been generated.

Start a development server with:

cd ${plugin.eltName}/frontend
yarn dev

Then visit http://localhost:1234/ and start hacking !`)
