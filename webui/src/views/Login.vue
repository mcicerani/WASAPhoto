<script>
import axios from 'axios';

export default {
  data() {
    return {
      username: '',
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('/api/login', {
          username: this.username,
        });

        // Presumendo che la risposta contenga un token con il prefisso "Bearer "
        const token = response.data.token;
        localStorage.setItem('token', token); // Salva il token nel localStorage

        // Estrai l'ID utente dal token
        const userId = token.split(" ")[1];

        // Reindirizza al profilo utente
        this.$router.push(`/users/${userId}/profile`);
      } catch (error) {
        console.error('Errore di login:', error);
      }
    }
  }
};
</script>

<template>
    <div>
        <h2>Login</h2>
        <form @submit.prevent="login">
            <div>
                <label for="username">Username</label>
                <input type="text" id="username" placeholder="Username" v-model="username">
            </div>
            <button type="submit">Login</button>
        </form>
    </div>
</template>