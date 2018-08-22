import Vue from 'vue'
import Vuex from 'vuex'
import UserApi from '../server/user'

Vue.use(Vuex)
const User = {
  actions: {
    doLogin ({commit}, user) {
      return UserApi.Login(user).then(response => {
        return response
      }).catch(err => {
        console.log(err)
        return err
      })
    }
  }
}
const moduleA = {
  state: {
    message: 'moduleA info'
  },
  mutations: {
    modify (state) {
      if (state.message === 'moduleA info') {
        state.message = '模块A中的信息'
      } else {
        state.message = 'moduleA info'
      }
    }
  },
  actions: {
    Amodify ({ commit }) {
      commit('modify')
    }
  }
}

const moduleB = {
  state: {
    info: 'moduleB info'
  },
  mutations: {
    change (state) {
      if (state.info === 'moduleB info') {
        state.info = '模块B中的信息'
      } else {
        state.info = 'moduleB info'
      }
    }
  },

  actions: {
    Bmodify ({ commit }) {
      commit('change')
    }
  }
}

const store = new Vuex.Store({
  modules: {
    User,
    a: moduleA,
    b: moduleB
  }
})

export default store
