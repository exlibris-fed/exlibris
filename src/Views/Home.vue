<template>
  <div id="app">
    <Home :axios="axiosInstance" />
  </div>
</template>

<script>
import Home from '../components/Home.vue'
import axios from 'axios'

export default {
  name: 'HomePage',
  components: {
    Home
  },
  data: function () {
    return {
      authToken: localStorage.getItem('auth'),
      axiosInstance: axios.create({
        baseURL: process.env.VUE_APP_API_ORIGIN,
        headers: {
          Authorization: 'Bearer ' + localStorage.getItem('auth')
        }
      })
    }
  },
  created: function () {
    if (!this.authToken) {
      this.$router.push({ name: 'login', props: { error: 'Log in to view this page' } })
    }
  }
}
</script>
