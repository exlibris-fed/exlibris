<template>
  <div id='search-wrapper'>
    <input @input="onInput" :placeholder="$t('search')" size="40"/>
  </div>
</template>

<script>
export default {
  name: 'SearchBar',
  methods: {
    onInput: function(event) {
      this.$emit('termChange', event.target.value);
    }
  },
  i18n: {
    messages: {
      en: {
        search: 'Search books...'
      }
    }
  }
};
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
