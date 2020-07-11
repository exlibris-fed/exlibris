<template>
  <div id="container">
    <b-alert
      v-if="action === 'register'"
      show
      variant="primary"
    >
      {{ $t('actions.register') }}
    </b-alert>
    <LoginForm
      :bounceto="bounceto"
      @resendVerificationEmail="resendVerificationEmail"
    />
  </div>
</template>

<script>
import LoginForm from '../components/LoginForm.vue'

export default {
  name: 'LoginPage',
  components: {
    LoginForm
  },
  props: {
    axios: {
      type: Function,
      required: true
    },
    bounceto: {
      type: String,
      default: '/'
    },
    action: {
      type: String,
      default: null
    }
  },
  // this blindly assumes that having *an* auth token means a user is logged in.
  // if the auth token isn't valid, it will log them out after the api returns 401.
  created: function () {
    if (localStorage.getItem('auth') !== null) {
      this.$router.push(this.bounceto)
    }
  },
  methods: {
    resendVerificationEmail (user) {
      const self = this
      this.axios.post('/verify/resend/' + user)
        .then(function () {
          self.$bvToast.toast(self.$t('verificationEmailResent.body'), {
            title: self.$t('verificationEmailResent.title'),
            solid: true,
            variant: 'info',
            autoHideDelay: 5000,
            appendToast: true
          })
        })
        .catch(function (error) {
          self.$bvToast.toast(error.message, {
            title: self.$t('error'),
            solid: true,
            variant: 'danger',
            autoHideDelay: 5000,
            appendToast: true
          })
        })
    }
  },
  i18n: {
    messages: {
      en: {
        verificationEmailResent: {
          title: 'Email sent',
          body: 'Check your email for the code'
        },
        error: 'Error',
        actions: {
          register: 'Thank you for registering! Check your email for a verification code.'
        }
      }
    }
  }
}
</script>
