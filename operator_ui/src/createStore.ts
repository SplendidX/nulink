import { Middleware } from 'redux'
import { createStore } from '@nulink/redux'
import thunkMiddleware from 'redux-thunk'
import { createExplorerConnectionMiddleware } from './middleware'
import reducer from './reducers'

const middleware: Middleware[] = [
  thunkMiddleware,
  createExplorerConnectionMiddleware(),
]

export default () => createStore(reducer, middleware)
