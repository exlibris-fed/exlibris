<template>
  <div id="app">
    <div id="title">
      <h1>Ex-Libris</h1>
      <img src="../public/greek-column.svg" alt="greek pillar" height="50" width="50" />
    </div>
    <div id="seachbar">
      <SearchBar @termChange="onTermChange"></SearchBar>
    </div>
    <div id="discover">
      <h1>Discover</h1>
    </div>
    <div id="bookgrid-container">
    <BookGrid :books="books"></BookGrid>
    </div>
  </div>
</template>

<script>
import SearchBar from './components/SearchBar.vue';
import BookGrid from './components/BookGrid.vue';
import axios from 'axios'

export default {
  name:'App',
  components: {
    SearchBar,
    BookGrid
  },
  data() {
    return { books: [] };
  },
  methods: { 	
    onTermChange(searchTerm) {
      axios.get('http://openlibrary.org/search.json?', {
        params: {
          q: searchTerm
        }
      }).then(response => { 
          this.books = response.data.docs;
      });
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
  display: inline-grid;
  grid-template-columns: 100px 100px 100px 100px;
  grid-template-rows: auto;
}

h1 {
  margin: 5px;
  font-family: 'Literata'; 
}

img {
  margin: 5px;
}

</style>
