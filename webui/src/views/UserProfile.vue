<template>
  <div>
    <h1>{{ userProfile.user.username }}</h1>
    <p>Numero di follower: {{ userProfile.numFollowers }}</p>
    <p>Numero di seguaci: {{ userProfile.numFollows }}</p>
    <p>Numero di foto caricate: {{ userProfile.numPhotos }}</p>
    <h2>Foto caricate:</h2>
    <ul>
      <li v-for="photo in userProfile.Photos" :key="photo.id">
        <img :src="photo.url" alt="Foto">
      </li>
    </ul>
  </div>
</template>


<script>
  import axios from 'axios'; // Importa l'istanza personalizzata di Axios
  
  export default {
  async mounted() {
    try {
      // Otteniamo il token memorizzato nel localStorage
      const token = localStorage.getItem('token');

      // Verifichiamo se il token è presente
      if (!token) {
        console.error('Token non trovato nel localStorage');
        return;
      }

      // Estraiamo l'ID dell'utente dal token
      const userId = token.split(' ')[1]; // Rimuoviamo il prefisso 'Bearer' e otteniamo direttamente l'ID

      // Effettuiamo una richiesta GET al backend per ottenere il profilo dell'utente
      const response = await axios.get(`/users/${userId}/profile`);

      // Se la richiesta ha successo, aggiorniamo lo stato del profilo utente nel componente Vue
      this.userProfile = response.data;
    } catch (error) {
      // Se si verifica un errore durante la richiesta, lo gestiamo qui
      console.error('Errore durante il recupero del profilo utente:', error);
    }
  },
  data() {
    return {
      userProfile: null // Inizializziamo userProfile come null, verrà aggiornato con i dati ricevuti dalla richiesta GET
    };
  }
};

</script>

<style>
</style>