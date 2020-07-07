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
  methods: {
    onInput: function (search) {
      this.$emit('termChange', search)
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
  transition: .5s ease-in-out;
}

input:focus {
  transform: scale(1.25);
}
</style>
