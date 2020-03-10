<template>
  <div class="container-fluid">
    <div
      id="grid-wrapper"
      class="row row-cols-6"
    >
      <BookGridCard
        v-for="book in books"
        :key="book.id"
        :book="book"
      />
    </div>
  </div>
</template>

<script>
import BookGridCard from './BookGridCard.vue'

export default {
  name: 'BookGrid',
  components: {
    BookGridCard
  },
  props: {
    books: {
      type: Array,
      required: true
    }
  }
}
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
