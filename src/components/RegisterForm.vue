<template>
  <div id="register">
    <h1>{{ $t('form.register') }}</h1>
    <div v-if="error">
      {{ error }}
    </div>
    <b-form @submit="register">
      <b-form-group
        id="username-group"
        :label="$t('form.username')"
        label-for="username"
      >
        <b-form-input
          id="username"
          v-model="username"
          type="text"
          required
        />
      </b-form-group>
      <b-form-group
        id="password-group"
        :label="$t('form.password')"
        label-for="password"
      >
        <b-form-input
          v-model="password"
          type="password"
          required
        />
      </b-form-group>
      <b-form-group
        id="password-confirm-group"
        :label="$t('form.confirmPassword')"
        label-for="password-confirm"
      >
        <b-form-input
          v-model="confirmPassword"
          type="password"
          required
        />
      </b-form-group>
      <b-form-group
        id="email-group"
        :label="$t('form.email')"
        label-for="email"
      >
        <b-form-input
          v-model="email"
          type="email"
          required
        />
      </b-form-group>
      <b-form-group
        id="display-name-group"
        :label="$t('form.displayName')"
        label-for="display-name"
      >
        <b-form-input
          v-model="displayName"
          type="text"
          required
        />
      </b-form-group>
      <b-button type="submit">
        {{ $t('form.send' ) }}
      </b-button>
    </b-form>
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
