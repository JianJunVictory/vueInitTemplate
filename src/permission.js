import router from './router'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css'// progress bar style
import {getToken} from '@/utils/auth'

NProgress.configure({ showSpinner: false })// NProgress Configuration

router.beforeEach((to, from, next) => {
  NProgress.start() // start progress bar
  if (to.meta.requireAuth) {
    if (getToken()) { // login status
      if (to.fullPath === '/login') {
        next('/')
      } else {
        next()
      }
    } else { // not login status
      if (to.name === 'ActivePage') {
        next()
      } else if (to.fullPath === '/login') {
        next()
      } else {
        next('/login')
      }
    }
  } else {
    next()
  }
  NProgress.done()
})

router.afterEach(() => {
  NProgress.done()
})
