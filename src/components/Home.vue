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
      :axios="axios"
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
