export function SetToken (token) {
  sessionStorage.setItem('authToken', token)
}
export function RemoveToken () {
  sessionStorage.removeItem('authToken')
}

export function getToken () {
  return sessionStorage.getItem('authToken')
}
