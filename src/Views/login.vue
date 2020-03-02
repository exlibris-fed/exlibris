<template>
  <div id="container">
    <div class="hello">
        <h1>Login</h1>
        <input type="text" v-model="username" placeholder="Username" />
        <input type="password" v-model="password" placeholder="Password" />
        <button v-on:click="login(username, password)"><router-link to="/home">Send</router-link></button>
        <br />
        <br />
        <div>{{ response }}</div>
    </div>
  </div>
</template>

<script>
import axios from "axios";

  export default {
    name: 'LoginPage',
    data () {
      return {
        firstname: "",
        lastname: "",
      }
    },
  methods: {
    login(username, password) {
      axios.post(process.env.VUE_APP_API_ORIGIN+'/authenticate', {
        "username": `${username}`,
        "password": `${password}`
      }).then(response => console.log(response))
    }
  }
}
</script>

<style scoped>
  h1, h2 {
    text-align: center;
    font-family: 'Literata';
    font-weight: normal;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    display: inline-block;
    margin: 0 10px;
  }

  a {
    color: #42b983;
  }

  textarea {
    width: 600px;
    height: 200px;
  }

  #container {
    display: flex;
    justify-content: center;
  }

  .hello {
    justify-content: center;
  }
</style>