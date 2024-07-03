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
        if (!token) {
          console.error('Token non trovato nel localStorage');
          return;
        }

        const response = await api.get(`/users?username=${this.search}`, {
          headers: {
            Authorization: `${token}`,
          },
        });
        const user = response.data;
        if (user) {
          this.userProfile = user;
          this.$router.push(`/users/${user.id}/profile`);
        } else {
          this.userProfile = null;
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
