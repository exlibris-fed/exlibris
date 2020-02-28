<template>
  <div id='search-wrapper'>
    <input @input="onInput" placeholder="Search Books..." size="40"/>
  </div>
</template>

<script>
export default {
  name: 'SearchBar',
  methods: {
    onInput: function(event) {
      this.$emit('termChange', event.target.value);
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
 