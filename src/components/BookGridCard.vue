<template>
  <div class="card">
    <img id="book-cover" :src="book.covers.medium" class="card-img-top p-4">
    <div class="card-body d-flex flex-column">
      <h5 class="card-title"> {{ book.title }} </h5>
      <p class="card-text">Author: {{ book.authors[0] }} </p>
      <div class="buttons">
        <button valign="bottom" type="button" class="btn btn-outline-success p-2">I've Read This</button>
        <button valign="bottom" type="button" class="btn btn-outline-warning mt-auto p-2">Write A Review</button>
      </div>  
    </div>
  </div> 
</template>

<script>
export default { 
  name: 'BookGridCard',
  props: ['book']
};
</script>

<style>
.card {
  --width: 100%;
  --time: 0.7s;
  display: flex;
  flex-direction: column;
  flex: 0 1 calc(25% - 1em);
  margin: 1rem .25em;

  color: white;
  background: #222;
  transition: ease-in .4s;
  animation-duration: 3s;
}

.card:hover {
  -webkit-filter: drop-shadow(0 20px 20px  rgba(0, 0, 0, 0.70));
  filter: drop-shadow(0 20px 20px rgba(0, 0, 0, 0.70));
  box-shadow: none;
  color: #222;
  background: #999;
}

.card-body {
  display: flex;
  flex-direction: column;
  margin-bottom: auto;
}

.card-text {
  margin-bottom: 0px; 
}

.buttons {
  display: flex;
  flex-direction: column;
  margin-top: auto;
}

button {
  margin-top: auto;
  margin-bottom: 3px;
}

#book-cover {
  margin: 0px;
  border: solid;
  border-color: white;
}

btn-outline-success:focus {
  color: green;
}
</style>
