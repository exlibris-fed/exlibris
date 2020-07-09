<template>
  <b-container>
    <b-spinner v-if="user == null" />
    <Profile
      v-else
      :user="user"
    />
  </b-container>
</template>

<script>
import axios from 'axios'
import Profile from '../components/Profile.vue'

export default {
  name: 'ProfilePage',

  components: {
    Profile
  },

  data: function () {
    return {
      axios: null,
      user: null
    }
  },

  created: function () {
    const self = this
    const config = {
      baseURL: process.env.VUE_APP_API_ORIGIN
    }
    if (localStorage.getItem('auth')) {
      config.headers = {
        Authorization: 'Bearer ' + localStorage.getItem('auth')
      }
    }
    this.axios = axios.create(config)

    this.axios.get('/user/' + this.$route.params.user)
      .then(function (response) { self.user = response.data })
      .catch(r => console.error(r))
  }
}
</script>
