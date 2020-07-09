<template>
  <div>
    <b-navbar
      sticky
      variant="light"
      toggleable="md"
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
      class="mt-3"
    >
      <router-view
        :axios="axios"
        :user="user"
        @login="handleLogin"
        @logout="handleLogout"
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
    const self = this
    const config = {
      baseURL: process.env.VUE_APP_API_ORIGIN
    }
    if (window.localStorage.getItem('auth') && window.localStorage.getItem('auth') != null) {
      const jwt = localStorage.getItem('auth')
      this.authenticated = true
      config.headers = {
        Authorization: 'Bearer ' + jwt
      }
      const payload = this.decodeJWT(jwt)
      if (payload && payload.kid) {
        this.axios = axios.create(config)
        this.axios.get('/user/' + payload.kid)
          .then(function (response) { self.user = response.data })
          .catch(r => console.error(r))
      }
    }
  },
  methods: {
    handleLogin () {
      this.authenticated = true
    },
    handleLogout () {
      this.authenticated = false
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
