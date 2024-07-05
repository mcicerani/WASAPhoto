<template>
  <div class="searchform">
    <form @submit.prevent="searchUsers">
      <h1>Ricerca Utenti</h1>
      <input type="text" v-model="search" placeholder="Cerca utenti" required>
      <button type="submit">Cerca</button>
    </form>
  </div>
</template>

<script>
import api from "@/services/axios";

export default {
  data() {
    return {
      search: '',
      userProfile: null,
      searchExecuted: false,
    };
  },
  methods: {
    async searchUsers() {
      try {
        const token = localStorage.getItem('token');
        const loggedInUserId = localStorage.getItem('loggedInUserId');

        if (!token) {
          console.log('Token non trovato nel localStorage');
          return;
        }

        // Search for the user by username
        const response = await api.get(`/users?username=${this.search}`, {
          headers: {
            Authorization: `${token}`,
          },
        });
        const user = response.data;
        if (user) {
          this.userProfile = user;

          // Check if the logged-in user is banned by the searched user
          const isBannedResponse = await api.get(`/users/${loggedInUserId}/bans/${user.id}`, {
            headers: {
              Authorization: `${token}`,
            },
          });

          if (isBannedResponse.data.isBanned) {
            alert(`Impossibile visualizzare il profilo cercato! Sei stato bannato da ${user.username}!`);
          } else {
            this.$router.push(`/users/${user.id}/profile`);
          }
        } else {
          this.userProfile = null;
          alert("Nessun Utente trovato con questo nome");
        }
        this.searchExecuted = true;
      } catch (error) {
        if (error.response && error.response.status === 401) {
          console.error('Unauthorized: ', error.response.data);
        } else {
          console.error('Errore durante la ricerca degli utenti:', error);
        }
        this.userProfile = null;
        this.searchExecuted = true;
        alert("Nessun Utente trovato con questo nome")
      }
    },
  },
};
</script>

<style>
.searchform {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 90vh;
}
</style>
