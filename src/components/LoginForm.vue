<template>
  <div id="login">
    <h1>{{ $t('form.login') }}</h1>
    <b-alert
      v-if="errorMessage"
      show
      variant="danger"
      class="text-center"
    >
      <div
        v-if="errorMessage === 'badPassword'"
      >
        {{ $t('errors.badPassword') }}
      </div>
      <div v-else-if="errorMessage === 'notVerified'">
        <p>{{ $t('errors.notVerified') }}</p>
        <b-button
          right
          @click="resendVerificationEmail"
        >
          {{ $t('form.resendVerificationEmail') }}
        </b-button>
      </div>
      <div v-else>
        {{ $t('errors.unknown') }}
      </div>
    </b-alert>
    <b-row
      class="mb-3"
    >
      <b-form
        class="col"
        inline
        @submit="login"
      >
        <label
          class="sr-only"
          for="username"
        >
          {{ $t('form.username') }}
        </label>
        <b-input
          id="username"
          v-model="username"
          type="text"
          class="mb-2 mr-sm-2 mb-sm-0"
          :placeholder="$t('form.username')"
          required
        />

        <label
          class="sr-only"
          for="password"
        >
          {{ $t('form.password') }}
        </label>
        <b-form-input
          id="password"
          v-model="password"
          type="password"
          class="mb-2 mr-sm-2 mb-sm-0"
          :placeholder="$t('form.password')"
          required
        />

        <b-button
          type="submit"
          class="mb-2 mr-sm-2 mb-sm-0"
        >
          {{ $t('form.send' ) }}
        </b-button>
      </b-form>
    </b-row>

    <b-row>
      <b-col>
        <b-button
          :to="{name: 'register'}"
        >
          {{ $t('register') }}
        </b-button>
      </b-col>
    </b-row>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'LoginForm',
  props: {
    error: {
      type: String,
      default: ''
    },
    bounceto: {
      type: String,
      default: '/'
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
    login (e) {
      e && e.preventDefault()
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
          this.$root.$emit('login')
          this.$router.push(this.bounceto)
        })
        .catch(error => {
          if (error.response && error.response.status === 401) {
            this.errorMessage = 'badPassword'
            return
          }
          if (error.response && error.response.status === 403) {
            this.errorMessage = 'notVerified'
            return
          }
          this.errorMessage = 'unknown'
          console.error(error)
        })
    },
    resendVerificationEmail () {
      this.$emit('resendVerificationEmail', this.username)
    }
  },
  i18n: {
    messages: {
      en: {
        form: {
          login: 'Login',
          username: 'Username',
          password: 'Password',
          send: 'Send',
          resendVerificationEmail: 'Resend verification email'
        },
        register: 'Register',
        errors: {
          badPassword: 'Invalid username/password combination',
          notVerified: 'Your account has not been verified',
          unknown: 'An error occurred during the request'
        }
      }
    }
  }
}
</script>
