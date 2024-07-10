<template>
    <div class="row mt-4">
      <div class="col-md-12">
        <ul class="listaFoto">
          <li v-for="(photo, index) in userStream" :key="photo.id" :class="{'new-row': index % 1 === 0}">
            <div class="card">
              <div class="card-body">
                <img class="text-center" :src="'data:image/jpeg;base64,' + photo.image_data" alt="User Photo">
                <div class="text-center likes">
                  <button @click="toggleLike(photo)" type="button" class="like-button btn btn-primary btn-sm align-self-center" :class="{'liked': photo.isLiked}" data-toggle="button" aria-pressed="false" autocomplete="off">
                    <svg class="feather">
                      <use href="/feather-sprite-v4.29.0.svg#heart"/>
                    </svg>
                    {{ photo.likes }}
                  </button>
                </div>
                <div class="text-center time">
                  <p>
                    <svg class="feather">
                      <use href="/feather-sprite-v4.29.0.svg#clock"/>
                    </svg>
                    {{ formatTimestamp(photo.timestamp) }}
                  </p>
                </div>
                <div class="text-center comments">
                  <p>
                    <svg class="feather">
                      <use href="/feather-sprite-v4.29.0.svg#message-square"/>
                    </svg>
                    {{ photo.comments }}
                  </p>
                </div>  
              </div>
              <!-- Pass userId and photoId props to Comments component -->
              <Comments :userId="photo.user_id" :photoId="photo.id" />
            </div>
          </li>
        </ul>
      </div>
    </div>
  </template>
  
  <script>
  import api from "@/services/axios";
  import Comments from "@/components/Comments.vue";
  
  export default {
    components: {
      Comments
    },
    data() {
      return {
        userStream: [] // Initialize userStream as an empty array
      };
    },
    mounted() {
      this.fetchUserStream(); // Fetch user stream photos when the component is mounted
    },
    methods: {
      async fetchUserStream() {
        try {
          const loggedInUserId = localStorage.getItem('loggedInUserId');
          const response = await api.get(`/users/${loggedInUserId}/stream`, {
            headers: {
              Authorization: localStorage.getItem("token")
            }
          });
          
          // Ordina le foto in base al timestamp (assumendo che timestamp sia in formato stringa YYYYMMDDHHMMSS)
          response.data.Photos.sort((a, b) => {
            if (a.timestamp > b.timestamp) return -1; // Ordine decrescente
            if (a.timestamp < b.timestamp) return 1;
            return 0;
          });
  
          this.userStream = response.data.Photos; // Assegna le foto ordinate a userStream
        } catch (error) {
          console.error(error);
          // Gestisci l'errore se necessario
        }
      },
      async toggleLike(photo) {
        const loggedInUserId = localStorage.getItem('loggedInUserId');
        const token = localStorage.getItem('token');

        try {
          // Fetch the current list of likes for the photo
          const likesResponse = await api.get(`/users/${loggedInUserId}/photos/${photo.id}/likes`, {
            headers: {
              Authorization: token
            }
          });

          const existingLikes = likesResponse.data || [];
          const existingLike = existingLikes.find(like => like.user_id == loggedInUserId);

          if (existingLike) {
            // Unlike the photo
            await api.delete(`/users/${loggedInUserId}/photos/${photo.id}/likes/${existingLike.id}`, {
              headers: {
                Authorization: token
              }
            });
            photo.isLiked = false;
            photo.likes--;
            console.log('Photo unliked successfully');
          } else {
            // Like the photo
            const response = await api.post(`/users/${loggedInUserId}/photos/${photo.id}/likes`, {}, {
              headers: {
                Authorization: token
              }
            });
            photo.isLiked = true;
            photo.likes++;
            console.log('Photo liked successfully');
          }
        } catch (error) {
          console.error('Error toggling like:', error);
        }
      },
      formatTimestamp(timestamp) {
        if (!timestamp || timestamp.length !== 14) {
          return ''; // Gestione di casi non validi, ad esempio timestamp mancante o formato non corretto
        }
        // Estrai le singole parti del timestamp
        const year = timestamp.substring(0, 4);
        const month = timestamp.substring(4, 6);
        const day = timestamp.substring(6, 8);
        const hours = timestamp.substring(8, 10);
        const minutes = timestamp.substring(10, 12);
        const seconds = timestamp.substring(12, 14);
        
        // Restituisci la stringa formattata nel formato desiderato
        return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
      }
    }
  };
  </script>
  