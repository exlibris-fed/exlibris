<template>
  <div id="container">
    <b-alert
      v-if="error"
      show
      variant="danger"
    >
      {{ error }}
    </b-alert>
    <h1>{{ $t('register') }}</h1>
    <RegisterForm
      @register="register"
    />
  </div>
</template>

<script>
import RegisterForm from '../components/RegisterForm.vue'

export default {
  name: 'Register',
  components: {
    RegisterForm
  },
  props: {
    axios: {
      type: Function,
      required: true
    }
  },
  data: function () {
    return {
      error: null
    }
  },
  methods: {
    register (data) {
      const self = this
      this.axios.post('/register', {
        username: data.username,
        password: data.password,
        email: data.email,
        display_name: data.displayName
      })
        .then(function () {
          self.error = undefined
          self.$router.push({
            name: 'login',
            query: {
              action: 'register'
            }
          })
        })
        .catch(function (error) {
          if (error.response && error.response.status === 409) {
            self.error = self.$t('error.duplicate')
          } else {
            self.error = self.$t('error.unknown')
            console.error(error)
          }
        })
    }
  },
  i18n: {
    messages: {
      en: {
        register: 'Register',
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
