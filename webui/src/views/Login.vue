<script>
import axios from 'axios';
import api from "@/services/axios"

export default {
  data() {
    return {
      username: '',
    };
  },
  methods: {
    async dologin() {
      try {
        const response = await api.post('/session', {
          username: this.username,
        });

        // Presumendo che la risposta contenga un token con il prefisso "Bearer "
        const token = response.data.token;
        localStorage.setItem('token', token); // Salva il token nel localStorage

        // Salva in localStorage username e userID
        // Estrai l'ID utente dal token
        const userId = token.split(" ")[1];
        localStorage.setItem('username', this.username)
        localStorage.setItem('loggedInUserId', userId);
        
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
    <div className="loginform">
        <h1>Login</h1>
        <form @submit.prevent="dologin">
          <div>
              <input type="text" id="username" placeholder="Username" v-model="username">
          </div>
          <button type="submit">Login</button>
        </form>
    </div>
</template>

<style>

  .loginform {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 90vh;
  }

  form {
    margin-top: 1rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  input {
    padding: 0.2rem;
    font-size: 1rem;
    margin-left: 1rem;
  }

  form button {
    padding: 0.2rem 1rem;
    font-size: 1rem;
    margin-left: 1rem;
    cursor: pointer;
    border-radius: 10px;
  }

</style>
