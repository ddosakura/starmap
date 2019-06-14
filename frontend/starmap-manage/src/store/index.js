import Vue from 'vue'
import Vuex from 'vuex'

import { moduleAuth } from './modules/auth'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    auth: moduleAuth
  }
})
