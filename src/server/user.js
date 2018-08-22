import axios from 'axios'

function Login (user) {
  return axios({
    method: 'POST',
    url: '/login',
    data: user
  })
}

export default {
  Login: Login
}
