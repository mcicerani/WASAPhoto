<template>
  <div>
    <h1>{{ userProfile.user.username }}</h1>
    
    <!-- Link per cambiare l'username (solo per il proprio profilo) -->
    <p v-if="isOwnProfile">              
      <RouterLink :to="`/users/${userProfile.user.id}/profile/edit`" class="nav-link">
        Cambia Username
      </RouterLink>
    </p>
    
    <div>
      <p>Followers: {{ userProfile.numFollowers }}</p>
      <p>Follows: {{ userProfile.numFollowing }}</p>
      <p>Foto: {{ userProfile.numPhotos }}</p>
    </div>
    
    <div>
      <h2>Foto:</h2>
      <ul>
        <li v-for="photo in userProfile.Photos" :key="photo.id">
          <img :src="photo.url" :alt="`Photo ${photo.id}`">
        </li>
      </ul>
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
        followers: [],
        numFollowers: 0,
        follows: [],
        numFollowing: 0,
        Photos: [],
        numPhotos: 0
      }
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
      console.log(`Il token è: ${token}`);

      const userId = token.split(' ')[1];
      console.log(`User ID from token: ${userId}`);

      console.log("Getting user profile data");
      const response = await api.get(`/users/${userId}/profile`, {
        headers: {
          Authorization: token
        }
      });
      console.log('User profile response:', response.data);
      this.userProfile = response.data;  // Assegna la risposta alla variabile userProfile
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
};
</script>

<style>
/* Stili CSS opzionali */
</style>
