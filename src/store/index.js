import Vue from 'vue'
import Vuex from 'vuex'
import UserApi from '@/server/user'
import {SetToken, getToken, RemoveToken} from '@/utils/auth'

Vue.use(Vuex)
const User = {
  state: {
    token: getToken()
  },
  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    }
  },
  actions: {
    doLogin ({commit}, user) {
      return new Promise((resolve, reject) => {
        UserApi.Login(user).then(response => {
          let repData = response.data
          if (repData.code === 0) {
            let token = repData.data.token
            commit('SET_TOKEN', token)
            SetToken(token)
            resolve()
          } else {
            reject(repData.message)
          }
        }).catch(err => {
          reject(err)
        })
      })
    },
    doLogout ({commit}) {
      return new Promise((resolve, reject) => {
        UserApi.Logout().then(response => {
          let resp = response.data
          if (resp.code === 0) {
            commit('SET_TOKEN', '')
            RemoveToken()
            resolve()
          } else {
            reject(resp.message)
          }
        }).catch(err => {
          reject(err)
        })
      })
    },
    doTest () {
      return new Promise((resolve, reject) => {
        UserApi.Test().then(response => {
          if (response.data.code === 0) {
            let resData = response.data.data
            resolve(resData)
          } else {
            reject(response.data.message)
          }
        }).catch(error => {
          reject(error)
        })
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
