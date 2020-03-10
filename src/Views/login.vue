<template>
  <div id="container">
    <div class="hello">
      <h1>{{ $t('form.login') }}</h1>
      <div v-if="errorMessage">
        {{ errorMessage }}
      </div>
      <input
        v-model="username"
        type="text"
        :placeholder="$t('form.username')"
      >
      <input
        v-model="password"
        type="password"
        :placeholder="$t('form.password')"
      >
      <button @click="login()">
        {{ $t('form.send' ) }}
      </button>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'LoginPage',
  props: {
    error: {
      type: String,
      default: ''
    },
    bounceto: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      username: undefined,
      password: undefined,
      errorMessage: this.error
    }
  },
  methods: {
    login () {
      axios.post(process.env.VUE_APP_API_ORIGIN + '/authenticate', {
        username: this.username,
        password: this.password
      })
        .then(response => {
          if (!response || !response.data || !response.data.bearer) {
            this.errorMessage = this.$t('errors.badPassword')
            return
          }
          this.errorMessage = ''
          localStorage.setItem('auth', response.data.bearer)
          this.$router.push(this.bounceto || '/')
        })
        .catch(error => {
          if (error.response && error.response.status === 401) {
            this.errorMessage = 'Invalid username/password combination'
            return
          }
          this.error = 'An error occurred during the request' // this sucks as well
          console.error(error)
        })
    }
  },
  i18n: {
    messages: {
      en: {
        form: {
          login: 'Login',
          username: 'Username',
          password: 'Password',
          send: 'Send'
        },
        errors: {
          badPassword: 'Invalid username/password combination'
        }
      }
    }
  }
}
</script>

<style scoped>
  h1, h2 {
    text-align: center;
    font-family: 'Literata';
    font-weight: normal;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }

  textarea {
    width: 600px;
    height: 200px;
  }

  #container {
    display: flex;
    justify-content: center;
  }

  .hello {
    justify-content: center;
  }
</style>
