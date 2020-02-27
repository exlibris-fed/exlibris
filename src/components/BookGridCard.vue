<template>
  <li>
    {{ book.title_suggest }}
  </li>
</template>

<script>
export default {
  name: 'BookGridCard',
  props: ['book']
};
</script>
