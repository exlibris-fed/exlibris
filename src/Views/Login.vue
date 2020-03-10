<template>
  <div id="container">
    <LoginForm />
  </div>
</template>

<script>
import LoginForm from '../components/LoginForm.vue'

export default {
  name: 'LoginPage',
  components: {
    LoginForm
  }
}
</script>

<style scoped>
  #container {
    display: flex;
    justify-content: center;
  }
</style>
