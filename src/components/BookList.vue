<template>
  <b-list-group>
    <Book
      v-for="book in feed"
      :key="book.id + book.timestamp"
      :book="book"
      type="list"
    />
  </b-list-group>
</template>

<script>
import Book from './Book.vue'

export default {
  name: 'BookList',

  components: {
    Book
  },

  props: {
    feed: {
      type: Array,
      default: () => []
    }
  }
}
</script>
