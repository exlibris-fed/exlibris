<template>
  <div>
    <b-navbar
      sticky
      variant="light"
      toggleable="md"
      class="mb-3"
    >
      <b-navbar-toggle target="nav-text-collapse" />
      <b-navbar-brand
        to="/"
      >
        <img
          class="text-right"
          src="../public/greek-column.svg"
          alt="greek pillar"
          height="50"
          width="50"
        >
        {{ $t('title') }}
      </b-navbar-brand>

      <b-collapse
        id="nav-text-collapse"
        is-nav
      >
        <b-navbar-nav
          class="ml-auto"
        >
          <HeaderProfile
            v-if="user"
            :user="user"
          />
          <b-nav-item
            v-else-if="authenticated"
            to="logout"
          >
            {{ $t.logout }}
          </b-nav-item>
          <b-nav-item
            v-else
            to="login"
          >
            {{ $t('login') }}
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-container
      id="app"
    >
      <router-view
        :axios="axios"
        :user="user"
      />
    </b-container>
  </div>
</template>

<script>
import axios from 'axios'
import HeaderProfile from './components/HeaderProfile.vue'

// TODO move header into its own component

export default {
  name: 'App',
  components: {
    HeaderProfile
  },
  data: function () {
    return {
      authenticated: false,
      axios: null,
      user: null
    }
  },
  created () {
    this.$root.$on('login', this.handleLogin)
    this.$root.$on('logout', this.handleLogout)

    if (window.localStorage.getItem('auth') && window.localStorage.getItem('auth') != null) {
      this.authenticated = true
    }
    this.axios = this.buildAxios()
    if (this.authenticated) {
      this.getAuthenticatedUser()
    }
  },
  methods: {
    buildAxios () {
      const config = {
        baseURL: process.env.VUE_APP_API_ORIGIN
      }
      if (this.authenticated) {
        config.headers = {
          Authorization: 'Bearer ' + localStorage.getItem('auth')
        }
      }
      return axios.create(config)
    },
    getAuthenticatedUser () {
      const self = this
      const payload = this.decodeJWT(localStorage.getItem('auth'))
      if (payload && payload.kid) {
        this.axios.get('/user/' + payload.kid)
          .then(function (response) { self.user = response.data })
          .catch(function (error) {
            // if your user doesn't exist, you are not logged in
            if (error && error.response && error.response.status === 404) {
              self.$router.push({ name: 'logout' })
            }
            console.error(error)
          })
      }
    },
    handleLogin () {
      this.authenticated = true
      this.axios = this.buildAxios()
      this.getAuthenticatedUser()
    },
    handleLogout () {
      this.authenticated = false
      this.axios = this.buildAxios()
      this.user = null
    },
    decodeJWT (token) {
      const base64Url = token.split('.')[1]
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
      return JSON.parse(window.atob(base64))
    }
  },
  i18n: {
    messages: {
      en: {
        title: 'exlibris',
        login: 'Log In',
        logout: 'Log Out'
      }
    }
  }
}
</script>
