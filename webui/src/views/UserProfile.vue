<template>
    <div>
      <h1>User Profile</h1>
      <p>Username: {{ user.username }}</p>
    </div>
  </template>
  
  <script>
  import axios from 'axios'; // Importa l'istanza personalizzata di Axios
  
  export default {
    data() {
      return {
        user: {}
      };
    },
    mounted() {
      this.fetchUserProfile();
    },
    methods: {
      async fetchUserProfile() {
        try {
          const userId = this.$route.params.userId;
          const response = await axios.get(`/users/${userId}/profile`);
          this.user = response.data;
        } catch (error) {
          console.error('Error fetching user profile:', error);
        }
      }
    }
  };
  </script>
  