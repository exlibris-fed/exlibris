<template>
  <div id="container">
    <div class="hello">
        <h1>Login</h1>
        <input type="text" v-model="input.firstname" placeholder="First Name" />
        <input type="text" v-model="input.lastname" placeholder="Last Name" />
        <button v-on:click="sendData()"><router-link to="/home">Send</router-link></button>
        <br />
        <br />
        <div>{{ response }}</div>
    </div>
  </div>
</template>

<script>
    export default {
        name: 'LoginPage',
        data () {
            return {
                ip: "",
                input: {
                    firstname: "",
                    lastname: ""
                },
                response: ""
            }
        },
        mounted() {
            this.$http.get("https://httpbin.org/ip").then(result => {
                this.ip = result.body.origin;
            }, error => {
                console.error(error);
            });
        },
        methods: {
            sendData() {
                this.$http.post("https://httpbin.org/post", this.input, { headers: { "content-type": "application/json" } }).then(result => {
                    this.response = result.data;
                }, error => {
                    console.error(error);
                });
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