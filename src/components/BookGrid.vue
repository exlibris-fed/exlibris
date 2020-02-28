<template>
  <div class="container-fluid">
    <div id="grid-wrapper" class="row row-cols-4">
      <BookGridCard 
        v-for="book in books" 
        :key="book.id"
        :book="book"
      >
      </BookGridCard>
    </div>
  </div>
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

#grid-wrapper {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
}

BookGridCard {
  margin: 5px;
}

BookGridCard:hover {
  -webkit-filter: drop-shadow(0 20px 20px  rgba(0, 0, 0, 0.70));
  filter: drop-shadow(0 20px 20px rgba(0, 0, 0, 0.70));
}

</style>
 