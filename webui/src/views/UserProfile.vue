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
          <button @click="toggleBan" class="btn btn-danger">{{ isBanned ? 'Unban' : 'Ban' }}</button>
          <!-- Pulsante per il toggle follow -->
          <button @click="toggleFollow" class="btn btn-primary">{{ isFollowing ? 'Unfollow' : 'Follow' }}</button>
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
      isBanned: false,
      isBannedByLoggedInUser: false
    };
  },
  computed: {
    isOwnProfile() {
      const loggedInUserId = localStorage.getItem('loggedInUserId');
      return loggedInUserId && parseInt(loggedInUserId) === this.userProfile.user.id;
    }
  },
  async mounted() {
    await this.loadUserProfile();
  },
  watch: {
    '$route.params.userId': 'loadUserProfile'
  },
  methods: {
    async loadUserProfile() {
      try {
        const userId = this.$route.params.userId; // Ottieni l'ID utente dai parametri della route
        console.log(`User ID from route: ${userId}`);

        console.log("Getting user profile data");
        const response = await api.get(`/users/${userId}/profile`, {
          headers: {
            Authorization: localStorage.getItem('token')
          }
        });
        console.log('User profile response:', response.data);
        this.userProfile = response.data;

        // Verifica se l'utente loggato sta seguendo questo profilo
        const loggedInUserId = localStorage.getItem('loggedInUserId');
        if (loggedInUserId) {
          const follows = await this.fetchFollows(loggedInUserId);
          const bans = await this.fetchBans(loggedInUserId);

          this.isFollowing = follows.includes(userId);
          this.isBannedByLoggedInUser = bans.includes(userId);
          this.isBanned = await this.checkIfBanned(userId, loggedInUserId);
        }

      } catch (error) {
        console.error('Errore nel caricamento del profilo:', error);
      }
    },
    async fetchFollows(userId) {
      try {
        const response = await api.get(`/users/${userId}/follows`);
        return response.data.map(user => user.id);
      } catch (error) {
        console.error('Errore nel recupero degli utenti seguiti:', error);
        return [];
      }
    },
    async fetchBans(userId) {
      try {
        const response = await api.get(`/users/${userId}/bans`);
        return response.data.map(user => user.id);
      } catch (error) {
        console.error('Errore nel recupero degli utenti bannati:', error);
        return [];
      }
    },
    async checkIfBanned(profileUserId, loggedInUserId) {
      try {
        const response = await api.get(`/users/${loggedInUserId}/bans/${profileUserId}`);
        return response.data.isBanned;
      } catch (error) {
        console.error('Errore nel recupero dello stato di ban:', error);
        return false;
      }
    },
    async toggleFollow() {
      try {
        const userId = this.userProfile.user.id;
        const loggedInUserId = localStorage.getItem('loggedInUserId');
        if (this.isFollowing) {
          // Unfollow
          await api.delete(`/users/${loggedInUserId}/follows/${userId}`);
        } else {
          // Follow
          await api.post(`/users/${loggedInUserId}/follows/${userId}`);
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
        const loggedInUserId = localStorage.getItem('loggedInUserId');

        // Carica i seguaci e i bannati dell'utente loggato
        const follows = await this.fetchFollows(loggedInUserId);
        const bans = await this.fetchBans(loggedInUserId);

        // Verifica se l'utente loggato è già bannato o seguito
        const isAlreadyFollowing = follows.includes(userId);
        const isAlreadyBanned = bans.includes(userId);

        // Verifica se l'utente visualizzato è già bannato dall'utente loggato
        if (isAlreadyBanned) {
          console.log('Utente già bannato da te, impossibile eseguire l\'azione.');
          return;
        }

        // Verifica se l'utente loggato è già seguito dall'utente visualizzato
        if (isAlreadyFollowing) {
          console.log('Utente già seguito da te, impossibile eseguire l\'azione.');
          return;
        }

        // Se non ci sono problemi, esegui l'azione di ban/unban
        if (this.isBanned) {
          // Unban
          await api.delete(`/users/${loggedInUserId}/bans/${userId}`);
        } else {
          // Ban
          await api.post(`/users/${loggedInUserId}/bans/${userId}`);
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
