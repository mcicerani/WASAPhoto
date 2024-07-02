<template>
  <div v-if="UserProfile.User">
    <h1>{{ UserProfile.User.Username }}</h1>
    <p>Numero di follower: {{ UserProfile.NumFollowers }}</p>
    <p>Numero di seguaci: {{ UserProfile.NumFollowing }}</p>
    <p>Numero di foto caricate: {{ UserProfile.NumPhotos }}</p>
    <h2>Foto caricate:</h2>
    <ul>
      <li v-for="photo in UserProfile.Photos" :key="photo.id">
        <img :src="photo.url" :alt="`Photo ${photo.id}`">
      </li>
    </ul>
  </div>
  <div v-else>
    <p>Caricamento...</p>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      UserProfile: {
        User: {
          ID: 0,
          Username: ''
        },
        Followers: [],
        NumFollowers: 0,
        Follows: [],
        NumFollowing: 0,
        Photos: [],
        NumPhotos: 0
      }
    };
  },
  async mounted() {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        console.error('Token non trovato nel localStorage');
        return;
      }
      console.log(`Il token è: ${token}`)

      const userId = token.split(' ')[1];
      console.log(`User ID from token: ${userId}`);

      try {
        const response = await axios.get(`/users/${userId}/profile`, {
          headers: {
            Authorization: token
          }
        });
        console.log('User profile response:', response.data);
        this.UserProfile = response.data;
        } catch (error) {
        if (error.response) {
          console.error('Errore nella risposta API:', error.response.status);
        } else if (error.request) {
          console.error('Errore nella richiesta:', error.request);
        } else {
          console.error('Errore generale:', error.message);
        }
        }
      this.UserProfile = response.data;  // Assegna la risposta alla variabile UserProfile
    } catch (error) {
      console.error('Errore durante il recupero del profilo utente:', error);
      // Gestisci l'errore nel modo appropriato (ad esempio, mostrando un messaggio all'utente)
    }
  }
};
</script>

<style>
/* Stili CSS opzionali */
</style>
