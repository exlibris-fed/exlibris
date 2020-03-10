<template>
  <div id="container">
    <RegisterForm />
  </div>
</template>

<script>
import RegisterForm from '../components/RegisterForm.vue'

export default {
  name: 'Register',
  components: {
    RegisterForm
  }
}
</script>

<style scoped>
  #container {
    display: flex;
    justify-content: center;
  }
</style>
