// @flow
import React, { Component } from 'react'
import { Provider } from 'react-redux'

import ReduxToastr from 'react-redux-toastr'
import 'react-redux-toastr/lib/css/react-redux-toastr.min.css'

import '../assets/react-toolbox/theme.css'
import ThemeProvider from 'react-toolbox/lib/ThemeProvider'
import theme from '../assets/react-toolbox/theme.js'

import store from '../store'
import { AppLayout } from '../layouts'

import './App.css'

export default class App extends Component {
  render() {
    return (
      <ThemeProvider theme={theme}>
        <Provider store={store}>
          <div className="App">
            <AppLayout />
            <ReduxToastr preventDuplicates position="top-right" />
          </div>
        </Provider>
      </ThemeProvider>
    )
  }
}
