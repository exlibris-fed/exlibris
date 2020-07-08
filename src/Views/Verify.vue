<template>
  <div class="container">
    <b-alert
      v-if="errorMessage"
      variant="danger"
      show
    >
      {{ errorMessage }}
    </b-alert>

    <div
      v-if="success"
    >
      {{ $t('success') }}
      <router-link
        :to="{name: 'login'}"
      >
        {{ $t('login') }}
      </router-link>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Verify',

  data: function () {
    return {
      errorMessage: null,
      success: false
    }
  },

  mounted: function () {
    const self = this
    const key = this.$route.params.key
    axios.get(process.env.VUE_APP_API_ORIGIN + '/verify/' + key)
      .then(function () { self.success = true })
      .catch(function (e) {
        if (e && e.response && e.response.status && e.response.status === 404) {
          self.errorMessage = 'Unknown verification code'
        } else {
          console.error(e)
          self.errorMessage = e.message || 'An unknown error occurred'
        }
      })
  },

  i18n: {
    messages: {
      en: {
        success: 'Your account had been verified. You may now',
        login: 'log in'
      }
    }
  }
}
</script>
