<script>
export default {
  name: 'LogoutPage',
  props: {
    bounceto: {
      type: String,
      default: '/'
    }
  },
  created: function () {
    localStorage.removeItem('auth')
    this.$root.$emit('logout')
    this.$router.push(this.bounceto)
  },
  // we don't actually want to render anything. this is essentially a redirect.
  render: function () {
    return ''
  }
}
</script>
