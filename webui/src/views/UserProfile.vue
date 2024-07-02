<template>
  <div class="container mt-5">
    <div class="row">
      <div class="col-md-4 offset-md-4">
        <h1 class="text-center">{{ userProfile.user.username }}</h1>
        <!-- Link per cambiare l'username (solo per il proprio profilo) -->
        <p v-if="isOwnProfile" class="text-center">              
          <RouterLink :to="`/users/${userProfile.user.id}/profile/edit`" class="nav-link">
            Cambia Username
          </RouterLink>
        </p>
        <p v-else class="text-center">
          <!-- Pulsante per il toggle ban -->
          <button @click="toggleBan(userProfile.user.id)" class="btn btn-danger">{{ isBanned ? 'Unban' : 'Ban' }}</button>
          <!-- Pulsante per il toggle follow -->
          <button @click="toggleFollow(userProfile.user.id)" class="btn btn-primary">{{ isFollowing ? 'Unfollow' : 'Follow' }}</button>
        </p>        
      </div>
    </div>

    <div class="row mt-4">
      <div class="col-md-4">
        <div>
          <p class="text-center">Followers</p>
          <p class="text-center">{{ userProfile.numFollowers }}</p>
        </div>
      </div>
      <div class="col-md-4">
        <div>
          <p class="text-center">Follows</p>
          <p class="text-center">{{ userProfile.numFollowing }}</p>
        </div>
      </div>
      <div class="col-md-4">
        <div>
          <p class="text-center">Foto</p>
          <p class="text-center">{{ userProfile.numPhotos }}</p>
        </div>
      </div>
    </div>

    <div class="row mt-4">
      <div class="col-md-12">
        <h2 class="text-center">Foto:</h2>
        <ul>
          <li v-for="photo in userProfile.Photos" :key="photo.id">
            <img :src="photo.url" :alt="`Photo ${photo.id}`">
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import api from "@/services/axios";
import { RouterLink } from "vue-router";

export default {
  components: {
    RouterLink
  },
  data() {
    return {
      userProfile: {
        user: {
          id: 0,
          username: ''
        },
        numFollowers: 0,
        numFollowing: 0,
        numPhotos: 0,
        Photos: []
      },
      isFollowing: false,
      isBanned: false
    };
  },
  computed: {
    isOwnProfile() {
      const token = localStorage.getItem('token');
      if (token) {
        const userId = token.split(' ')[1];
        return parseInt(userId) === this.userProfile.user.id;
      }
      return false;
    }
  },
    async mounted() {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        console.error('Token non trovato nel localStorage');
        return;
      }

      // Recupera l'ID dell'utente destinatario dall'URL usando $route.params
      const userId = this.$route.params.userId;
      console.log(`User ID from route params: ${userId}`);

      // Esegui la chiamata API per ottenere il profilo dell'utente destinatario
      const response = await api.get(`/users/${userId}/profile`, {
        headers: {
          Authorization: token
        }
      });
      console.log('User profile response:', response.data);
      this.userProfile = response.data;

      // Verifica se l'utente loggato sta seguendo questo profilo
      if (this.isOwnProfile) {
        this.isFollowing = true; // Se è il proprio profilo, assume che si stia seguendo
      } else {
        this.isFollowing = await this.checkIfFollowing(userId);
      }

      // Verifica se l'utente loggato è bannato dal profilo
      this.isBanned = await this.checkIfBanned(userId);

    } catch (error) {
      console.error('Errore nel caricamento del profilo:', error);
    }
  },
  methods: {
    async checkIfFollowing(userId) {
      try {
        const response = await api.get(`/users/${userId}/follows/${this.userProfile.user.id}`);
        return response.data.isFollowing;
      } catch (error) {
        console.error('Errore nel recupero dello stato di follow:', error);
        return false;
      }
    },
    async checkIfBanned(userId) {
      try {
        const response = await api.get(`/users/${userId}/bans/${this.userProfile.user.id}`);
        return response.data.isBanned;
      } catch (error) {
        console.error('Errore nel recupero dello stato di ban:', error);
        return false;
      }
    },
    async toggleFollow() {
      try {
        const userId = this.userProfile.user.id;
        if (this.isFollowing) {
          // Unfollow
          await api.delete(`/users/${userId}/follows/${userId}`);
        } else {
          // Follow
          await api.post(`/users/${userId}/follows/${userId}`);
        }
        // Aggiorna lo stato di isFollowing dopo l'azione
        this.isFollowing = !this.isFollowing;
      } catch (error) {
        console.error('Errore nel toggle follow:', error);
      }
    },
    async toggleBan() {
      try {
        const userId = this.userProfile.user.id;
        if (this.isBanned) {
          // Unban
          await api.delete(`/users/${userId}/bans/${userId}`);
        } else {
          // Ban
          await api.post(`/users/${userId}/bans/${userId}`);
        }
        // Aggiorna lo stato di isBanned dopo l'azione
        this.isBanned = !this.isBanned;
      } catch (error) {
        console.error('Errore nel toggle ban:', error);
      }
    }
  }
};
</script>

<style>
  h1 {
    text-transform: capitalize;
  }
</style>
