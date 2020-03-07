<template>
  <div id="app">
    <div id="seachbar">
      <SearchBar @termChange="onTermChange"></SearchBar>
    </div>
    <div id="discover">
      <h2>{{ $t('callToAction') }}</h2>
    </div>
    <div id="bookgrid-container">
    <BookGrid :books="books"></BookGrid>
    </div>
  </div>
</template>

<script>
import SearchBar from '../components/SearchBar.vue';
import BookGrid from '../components/BookGrid.vue';
import axios from 'axios'

export default {
  name:'HomePage',
  components: {
    SearchBar,
    BookGrid
  },
  data() {
    return { books: [] };
  },
  methods: {
    onTermChange(searchTerm) {
      axios.get(process.env.VUE_APP_API_ORIGIN+'/book', {
        params: {
          title: searchTerm
        }
      }).then(response => {
          this.books = response.data;
      });
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
