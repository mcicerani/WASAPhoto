<script>

import axios from 'axios';
import api from "@/services/axios"


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
        const response = await api.get(`/api/users?search=${this.search}`);
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
      <form className="searchform">
          <h1>Ricerca Utenti</h1>
          <input type="text" v-model="search" placeholder="Cerca utenti">
          <button @click="searchUsers">Cerca</button>
        </form>
        <ul>
            <li v-for="user in users" :key="user.id">
                <router-link :to="`/users/${user.id}/profile`">{{ user.username }}</router-link>
            </li>
        </ul>
    </div>
</template>


<style>
  .searchform {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
</style>
