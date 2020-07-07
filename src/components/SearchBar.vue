<template>
  <b-form-input
    :placeholder="$t('search')"
    size="40"
    @input="onInput"
  />
</template>

<script>
export default {
  name: 'SearchBar',
  data () {
    return { timer: null }
  },
  methods: {
    onInput: function (search) {
      if (this.timer) {
        clearTimeout(this.timer)
        this.timer = null
      }
      this.timer = setTimeout(() => {
        this.$emit('termChange', search)
      }, 800)
    }
  },
  i18n: {
    messages: {
      en: {
        search: 'Search books...'
      }
    }
  }
}
</script>

<style>
input {
  text-align: center;
  transition: 0.5s ease-in-out;
}

input:focus {
  transform: scale(1.25);
}
</style>
