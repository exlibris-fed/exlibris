<template>
  <div id="app">
    <div id="title">
      <h1>Ex-Libris</h1>
      <img src="../public/greek-column.svg" alt="greek pillar" height="50" width="50" />
    </div>
    <div id="seachbar">
      <SearchBar></SearchBar>
    </div>
    <div id="content">
      <h1>Discover</h1>
    </div>
  </div>
</template>

<script>
import SearchBar from './components/SearchBar.vue';

export default {
  name:'App',
  components: {
    SearchBar
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

h1 {
  margin: 5px;
  font-family: 'Literata'; 
}

img {
  margin: 5px;
}

</style>
