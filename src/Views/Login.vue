<template>
  <div id="container">
    <LoginForm
      :bounceto="bounceto"
      @login="handleLogin"
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
    bounceto: {
      type: String,
      default: '/'
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
    handleLogin () {
      this.$emit('login')
    }
  }
}
</script>
