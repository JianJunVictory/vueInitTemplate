import router from './router'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css'// progress bar style

NProgress.configure({ showSpinner: false })// NProgress Configuration

router.beforeEach((to, from, next) => {
  NProgress.start() // start progress bar
  if (to.path === '/') {
    console.log('aaaaaaaaaaaaaaaaaaaaa')
    console.log('aaaaaVVVVVVVVVVVVVVVVaaaaaaa')
  }
  next()
  NProgress.done() // if current page is login will not trigger afterEach hook, so manually handle it
})

router.afterEach(() => {
  NProgress.done() // finish progress bar
})
