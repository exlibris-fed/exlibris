<template>
  <Home :axios="axiosInstance" />
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
      axiosInstance: () => {}
    }
  },
  created: function () {
    if (!this.authToken) {
      this.$router.push({ name: 'login', props: { error: 'Log in to view this page' } })
    }
    console.log('process.env is:')
    console.log(process.env)
    this.axiosInstance = axios.create({
      baseURL: process.env.VUE_APP_API_ORIGIN,
      headers: {
        Authorization: 'Bearer ' + localStorage.getItem('auth')
      }
    })
    this.axiosInstance.interceptors.response.use((r) => r, this.logOutOnError)
  },
  methods: {
    logOutOnError: function (error) {
      if (error && error.response && error.response.status && error.response.status === 401) {
        this.$router.push({ name: 'logout' })
      }
      return Promise.reject(error)
    }
  }
}
</script>
