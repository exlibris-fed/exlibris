<template>
  <div id="search-wrapper">
    <input
      placeholder="Search Books..."
      size="40"
      @input="onInput"
    >
  </div>
</template>

<script>
export default {
  name: 'SearchBar',
  methods: {
    onInput: function (event) {
      this.$emit('termChange', event.target.value)
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
