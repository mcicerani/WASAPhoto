<template>
    <div className="usernameForm">
      <h1>Cambia Username</h1>
      <form @submit.prevent="changeUsername">
        <div>
          <label for="newUsername">Nuovo Username:</label>
          <input type="text" id="newUsername" v-model="newUsername">
        </div>
        <button type="submit">Salva</button>
      </form>
    </div>
  </template>
  
  <script>
  import api from "@/services/axios";
  
  export default {
    data() {
      return {
        newUsername: ''
      };
    },
    methods: {
      async changeUsername() {
        try {
          const token = localStorage.getItem('token');
          if (!token) {
            console.error('Token non trovato nel localStorage');
            return;
          }
          const userId = token.split(' ')[1];
          
          const response = await api.put(`/users/${userId}/profile/edit`, {
            username: this.newUsername
          }, {
            headers: {
              Authorization: token
            }
          });
          
          console.log('Username cambiato con successo:', response.data);
          alert('Username cambiato con successo!');
          // Puoi aggiungere il reindirizzamento alla pagina del profilo o a un'altra pagina di conferma
        } catch (error) {
          if (error.response) {
            console.error('Errore nella risposta API:', error.response.status);
          } else if (error.request) {
            console.error('Errore nella richiesta:', error.request);
          } else {
            console.error('Errore generale:', error.message);
          }
        }
      }
    }
  };
  </script>
  
  <style>

  .usernameForm {
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
    