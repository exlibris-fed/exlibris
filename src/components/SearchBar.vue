<template>
  <div>
    <input v-on:input="onInput" size="40"/>
  </div>
</template>

<script>
export default {
  name: 'SearchBar',
  methods: {
    onInput: function(event) {
      console.log(event.target.value);
    }
  }
};
</script>

<style>
</style>
 