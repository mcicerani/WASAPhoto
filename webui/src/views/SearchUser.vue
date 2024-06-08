<script>

import axios from 'axios';

export default {
  data() {
    return {
      users: [],
      search: '',
    };
  },
  methods: {
    async searchUsers() {
      try {
        const response = await axios.get(`/api/users?search=${this.search}`);
        this.users = response.data;
      } catch (error) {
        console.error('Errore durante la ricerca degli utenti:', error);
      }
    }
  }
};

</script>

<template>
    <div>
        <h2>Ricerca Utenti</h2>
        <input type="text" v-model="search" placeholder="Cerca utenti">
        <button @click="searchUsers">Cerca</button>
        <ul>
            <li v-for="user in users" :key="user.id">
                <router-link :to="`/users/${user.id}/profile`">{{ user.username }}</router-link>
            </li>
        </ul>
    </div>
</template>
