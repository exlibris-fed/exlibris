<template>
  <ul>
    <BookGridCard 
      v-for="book in books" 
      :key="book.id"
      v-bind:book="book"
    >
    </BookGridCard>
  </ul>
</template>

<script>
import BookGridCard from './BookGridCard.vue';

export default {
  name: 'BookGrid',
  components: {
    BookGridCard
  },
  props: ['books']
};
</script>

<style>
BookGridCard {
  display: inline-grid;
}
</style>
 