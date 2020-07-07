<template>
  <div id="app">
    <div id="seachbar">
      <SearchBar @termChange="onTermChange" />
    </div>
    <div id="discover">
      <h2>{{ $t('callToAction') }}</h2>
    </div>
    <div id="bookgrid-container">
      <BookGrid :books="books" />
    </div>
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
      type: Object,
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
        this.books = response.data
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

<style>

#title {
  margin: 20px;
  display: flex;
  justify-content: center;
}

#seachbar {
  display: flex;
  justify-content: center;
}

BookGird {
  display: grid;
}

h1 {
  margin: 5px;
  font-family: 'Literata';
}

h2 {
  margin: 5px;
  font-family: 'Literata';
  text-align: center;
}

img {
  margin: 5px;
}

</style>
