import axios from 'axios'
import store from '@/store'
import Router from '@/router'
import {getToken, RemoveToken} from '@/utils/auth'

axios.interceptors.request.use(config => {
  if (store.state.User.token) {
    config.headers['authorization'] = 'bearer ' + getToken()
  }
  return config
}, error => {
  // Do something with request error
  console.log(error) // for debug
  Promise.reject(error)
})

axios.interceptors.response.use(response => {
  let rep = response.data
  if (rep.code === -2001) {
    RemoveToken()
    Router.replace({path: '/login'})
    return response
  } else {
    return response
  }
}, error => {
  console.log('err' + error) // for debug
  return Promise.reject(error)
})
function Login (user) {
  return axios({
    method: 'POST',
    url: '/login',
    data: user
  })
}
function Active (token) {
  return axios({
    method: 'POST',
    url: '/active',
    data: {'token': token}
  })
}
function Logout () {
  return axios({
    method: 'POST',
    url: '/logout'
  })
}
function GetDB () {
  return axios({
    method: 'POST',
    url: '/test'
  })
}
export default {
  Login: Login,
  Logout: Logout,
  Test: GetDB,
  Active: Active
}
