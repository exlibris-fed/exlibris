<template>
  <Book
    v-if="book"
    :book="book"
    :axios="axios"
  />
  <b-spinner
    v-else
  />
</template>

<script>
import Book from '../components/Book.vue'

export default {
  name: 'BookDetail',
  components: {
    Book
  },
  props: {
    axios: {
      type: Function,
      required: true
    }
  },
  data: function () {
    return {
      book: null
    }
  },
  created () {
    const self = this
    this.axios.get('/book/' + this.$route.params.book)
      .then(function (response) { self.book = response.data })
      .catch(function (error) {
        self.$bvToast.toast(error.message, {
          title: self.$t('error'),
          solid: true,
          variant: 'danger',
          autoHideDelay: 5000,
          appendToast: true
        })
      })
  }
}
</script>
