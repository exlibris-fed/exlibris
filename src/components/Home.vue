<template>
  <div>
    <b-row>
      <b-col>
        <SearchBar @termChange="onTermChange" />
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <h2>{{ $t('callToAction') }}</h2>
      </b-col>
    </b-row>
    <BookGrid
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
    return { books: [] }
  },
  methods: {
    onTermChange (searchTerm) {
      this.axios.get('/book', {
        params: {
          title: searchTerm
        }
      }).then(response => {
        // only have 16 rows of 3
        this.books = response.data.slice(0, 48)
      })
    },

    read (book) {
      const id = book.id.split('/')[1]
      this.axios.post('/book/' + id + '/read')
        .then(r => console.log(r)) // TODO
        .error(r => console.error(r)) // also TODO
    }
  },
  i18n: {
    messages: {
      en: {
        callToAction: 'What have you read lately?'
      }
    }
  }
}
</script>
