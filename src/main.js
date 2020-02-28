import Vue from 'vue'
import App from './App.vue'
import VueResource from 'vue-resource';
import VueRouter from 'vue-router';
import HomePage from './Views/Home.vue'
import LoginPage from './Views/login.vue'

Vue.use(VueRouter);
Vue.use(VueResource);
Vue.config.productionTip = false

const routes = [
  { path: '/home', component: HomePage},
  { path: '/', component: LoginPage }
];

const router = new VueRouter({
  routes,
  mode: 'history'
});

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
