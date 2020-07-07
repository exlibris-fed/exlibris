<template>
  <b-row>
    <BookGridCard
      v-for="book in books"
      :key="book.id"
      :book="book"
      @read="readBook"
    />
  </b-row>
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
  },
  methods: {
    readBook: function (book) {
      this.$emit('read', book)
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
