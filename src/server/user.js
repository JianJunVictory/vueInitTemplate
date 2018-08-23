import axios from 'axios'
import store from '@/store'
import {getToken} from '@/utils/auth'

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

function Login (user) {
  return axios({
    method: 'POST',
    url: '/login',
    data: user
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
  Test: GetDB
}
