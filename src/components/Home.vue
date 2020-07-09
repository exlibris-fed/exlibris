<template>
  <div>
    <b-row>
      <b-col>
        <h2>{{ $t('callToAction') }}</h2>
      </b-col>
    </b-row>
    <b-row
      class="mb-5"
    >
      <b-col>
        <SearchBar @termChange="onTermChange" />
      </b-col>
    </b-row>
    <div
      v-if="loading"
      class="text-center"
    >
      <b-spinner />
    </div>
    <BookGrid
      v-else
      :books="books"
      @read="read"
    />
  </div>
</template>

<script>
import SearchBar from '../components/SearchBar.vue'
import BookGrid from '../components/BookGrid.vue'

export default {
  name: 'HomePage',
  components: {
    SearchBar,
    BookGrid
  },
  props: {
    axios: {
      type: Function,
      required: true
    }
  },
  data () {
    return {
      books: [],
      lastRead: {},
      loading: false
    }
  },
  methods: {
    onTermChange (searchTerm) {
      this.loading = true
      this.axios.get('/book', {
        params: {
          title: searchTerm
        }
      }).then(response => {
        // only have 16 rows of 3
        this.loading = false
        this.books = response.data.slice(0, 48)
      })
    },

    read (book) {
      const self = this
      const id = book.id.split('/')[2]
      this.lastRead = book
      this.axios.post('/book/' + id + '/read')
        .then(self.successToast)
        .catch(self.errorToast)
    },

    successToast () {
      self.$bvToast.toast(this.$t('readSuccess.message', { title: this.lastRead.title }), {
        title: this.$t('readSuccess.title'),
        solid: true,
        variant: 'info',
        autoHideDelay: 5000,
        appendToast: true
      })
    },

    errorToast (error) {
      self.$bvToast.toast(error.message, {
        title: self.$t('error'),
        solid: true,
        variant: 'danger',
        autoHideDelay: 5000,
        appendToast: true
      })
    }
  },
  i18n: {
    messages: {
      en: {
        callToAction: 'What have you read lately?',
        error: 'Error',
        readSuccess: {
          title: 'Book Read',
          message: '{title} has been added to your feed'
        }
      }
    }
  }
}
</script>
