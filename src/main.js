import Vue from 'vue'
import App from './App.vue'
import VueI18n from 'vue-i18n'
import VueResource from 'vue-resource'
import VueRouter from 'vue-router'
import { BootstrapVue, IconsPlugin } from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import HomePage from './Views/Home.vue'
import LoginPage from './Views/Login.vue'
import LogoutPage from './Views/Logout.vue'
import Register from './Views/Register.vue'
import Verify from './Views/Verify.vue'
import ProfilePage from './Views/ProfilePage.vue'
import BookDetail from './Views/BookDetail.vue'

Vue.use(VueI18n)
Vue.use(VueRouter)
Vue.use(VueResource)
Vue.use(BootstrapVue)
Vue.use(IconsPlugin)
Vue.config.productionTip = false

const routes = [
  { path: '/', component: HomePage, name: 'home' },
  { path: '/user/:user', component: ProfilePage, name: 'profile', alias: '/@:user' },
  { path: '/book/:book', component: BookDetail, name: 'book' },
  { path: '/login', component: LoginPage, name: 'login' },
  { path: '/logout', component: LogoutPage, name: 'logout' },
  { path: '/register', component: Register, name: 'register' },
  { path: '/verify/:key', component: Verify, name: 'verify' }
]

const router = new VueRouter({
  routes,
  mode: 'history'
})

const i18n = new VueI18n({
  locale: 'en'
})

new Vue({
  router,
  i18n,
  render: h => h(App)
}).$mount('#app')
