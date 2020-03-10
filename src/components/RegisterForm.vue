<template>
  <div id="register">
    <h1>{{ $t('form.register') }}</h1>
    <div v-if="error">
      {{ error }}
    </div>
    <form @submit="register">
      <div>
        <input
          v-model="username"
          type="text"
          required
          :placeholder="$t('form.username')"
        >
      </div>
      <div>
        <input
          v-model="password"
          type="password"
          required
          :placeholder="$t('form.password')"
        >
      </div>
      <div>
        <input
          v-model="confirmPassword"
          type="password"
          required
          :placeholder="$t('form.confirmPassword')"
        >
      </div>
      <div>
        <input
          v-model="email"
          type="email"
          required
          :placeholder="$t('form.email')"
        >
      </div>
      <div>
        <input
          v-model="displayName"
          type="text"
          required
          :placeholder="$t('form.displayName')"
        >
      </div>
      <div>
        <button type="submit">
          {{ $t('form.send' ) }}
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'RegisterForm',
  props: {
    bounceto: {
      type: String,
      default: ''
    }
  },
  data () {
    return {
      username: undefined,
      password: undefined,
      confirmPassword: undefined,
      email: undefined,
      displayName: undefined,
      error: undefined
    }
  },
  methods: {
    register (e) {
      e.preventDefault()
      if (this.password !== this.confirmPassword) {
        this.error = this.$t('error.mismatchedPassword')
        return
      }

      axios.post(process.env.VUE_APP_API_ORIGIN + '/register', {
        username: this.username,
        password: this.password,
        email: this.email,
        display_name: this.displayName
      })
        .then(response => {
          if (!response || response.status !== 201) {
            this.error = this.$t('error.unknown')
            return
          }
          this.error = ''
          // TODO display a "hey cool thanks" message. and also email a confirmation token.
          this.$router.push(this.bounceto || { name: 'login' })
        })
        .catch(error => {
          if (error.response && error.response.status === 409) {
            this.error = this.$t('error.duplicate')
          } else {
            this.error = this.$t('error.unknown')
            console.error(error)
          }
        })
    }
  },
  i18n: {
    messages: {
      en: {
        form: {
          register: 'Register',
          username: 'Username',
          password: 'Password',
          confirmPassword: 'Password (again)',
          email: 'Email address',
          displayName: 'Display Name',
          send: 'Send'
        },
        error: {
          mismatchedPassword: 'Passwords do not match',
          duplicate: 'That username is already taken',
          unknown: 'An unknown error occurred'
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
</style>
