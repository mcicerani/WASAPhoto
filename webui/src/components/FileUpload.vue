<template>
    <!-- Bottone per caricare la foto -->
    <input type="file" style="display: none" ref="fileInput" @change="uploadPhotoDirectly">
    <button class="btn btn-primary" @click="$refs.fileInput.click()">Carica Foto</button>
  </template>
  
  <script>
  import api from "@/services/axios";
  
  export default {
    data() {
      return {
        selectedFile: null,
        message: '',
        success: false,
      };
    },
    methods: {
      async uploadPhotoDirectly(event) {
        // Ottieni il file selezionato dall'evento
        this.selectedFile = event.target.files[0];
  
        // Verifica se è stato selezionato un file
        if (!this.selectedFile) {
          this.message = 'Please select a file.';
          this.success = false;
          return;
        }
  
        const userId = localStorage.getItem("loggedInUserId");  
        const token = localStorage.getItem("token");
        const formData = new FormData();
        formData.append('image', this.selectedFile);
  
        try {
          const response = await api.post(`/users/${userId}/photos`, formData, {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': token
            }
          });
  
          this.success = true;
          console.log(response.data);
          // Esegui qualche azione dopo il caricamento, come ricaricare la pagina
          window.location.reload();
        } catch (error) {
          alert('Error uploading photo: ' + error.response.data.message);
          this.success = false;
          console.error(error);
        }
      }
    }
  };
  </script>
  
  <style scoped>
  .upload-photo {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  
  /* Stili per il bottone */
  button {
    padding: 2rem;
    background-color: #007bff;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }
  
  button:hover {
    background-color: #0056b3;
  }
  </style>
  