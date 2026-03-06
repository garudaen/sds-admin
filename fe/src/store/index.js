import { createStore } from 'vuex'
import app from './modules/app'
import domains from './modules/domains'

export default createStore({
  modules: {
    app,
    domains,
  },
})
