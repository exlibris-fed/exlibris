import Vue from 'vue'
import App from './App.vue'
import VueI18n from 'vue-i18n'
import VueResource from 'vue-resource'
import VueRouter from 'vue-router'
import HomePage from './Views/Home.vue'
import LoginPage from './Views/login.vue'

Vue.use(VueI18n)
Vue.use(VueRouter)
Vue.use(VueResource)
Vue.config.productionTip = false

const routes = [
  { path: '/', component: HomePage, name: 'home' },
  { path: '/login', component: LoginPage, name: 'login' }
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
