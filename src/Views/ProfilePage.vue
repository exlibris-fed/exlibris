<template>
  <b-container>
    <b-spinner v-if="user == null" />
    <Profile
      v-else
      :user="user"
      :feed="feed"
    />
  </b-container>
</template>

<script>
import Profile from '../components/Profile.vue'

export default {
  name: 'ProfilePage',

  components: {
    Profile
  },

  props: {
    axios: {
      type: Function,
      required: true
    },
    user: {
      type: Object,
      default: null
    }
  },

  data: function () {
    return {
      feed: null
    }
  },

  created () {
    const self = this
    this.axios.get('/book/read')
      .then(function (response) { self.feed = response.data })
      .catch(function (error) {
        self.$bvToast.toast(error.message, {
          title: self.$t('error'),
          solid: true,
          variant: 'danger',
          autoHideDelay: 5000,
          appendToast: true
        })
      })
  },

  i18n: {
    messages: {
      en: {
        error: 'Error'
      }
    }
  }
}
</script>
