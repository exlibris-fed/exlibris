<template>
  <b-row>
    <Book
      v-for="book in books"
      :key="book.id"
      type="card"
      :book="book"
      :axios="axios"
    />
  </b-row>
</template>

<script>
import Book from './Book.vue'

export default {
  name: 'BookGrid',
  components: {
    Book
  },
  props: {
    books: {
      type: Array,
      required: true
    },
    axios: {
      type: Function,
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
