<template>
  <div id="app">
    <Home />
  </div>
</template>

<script>
import Home from '../components/Home.vue'

export default {
  name: 'HomePage',
  components: {
    Home
  },
  data: function () {
    return {
      authToken: localStorage.getItem('auth')
    }
  },
  created: function () {
    if (!this.authToken) {
      this.$router.push({ name: 'login', props: { error: 'Log in to view this page' } })
    }
  }
}
</script>
