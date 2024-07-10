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

        // Check if Photos exists and is an array
        if (response.data.Photos && Array.isArray(response.data.Photos)) {
          // Sort the photos by timestamp (assuming timestamp is a string in format YYYYMMDDHHMMSS)
          response.data.Photos.sort((a, b) => {
            if (a.timestamp > b.timestamp) return -1; // Descending order
            if (a.timestamp < b.timestamp) return 1;
            return 0;
          });

          this.userStream = response.data.Photos; // Assign the sorted photos to userStream
        } else {
          this.userStream = []; // Assign an empty array if there are no photos
        }
      } catch (error) {
        console.error(error);
        // Handle the error if necessary
      }
    },
    async toggleLike(photo) {
      const loggedInUserId = localStorage.getItem('loggedInUserId');
      const token = localStorage.getItem('token');

      try {
        // Fetch the current list of likes for the photo
        const likesResponse = await api.get(`/users/${photo.user_id}/photos/${photo.id}/likes`, {
          headers: {
            Authorization: token
          }
        });

        const existingLikes = likesResponse.data || [];
        const existingLike = existingLikes.find(like => like.user_id == loggedInUserId); // Ensure comparison with integer

        if (existingLike) {
          // Unlike the photo
          await api.delete(`/users/${photo.user_id}/photos/${photo.id}/likes/${existingLike.id}`, {
            headers: {
              Authorization: token
            }
          });
          photo.isLiked = false;
          photo.likes--;
          console.log('Photo unliked successfully');
        } else {
          // Like the photo
          await api.post(`/users/${photo.user_id}/photos/${photo.id}/likes`, {}, {
            headers: {
              Authorization: token
            }
          });
          photo.isLiked = true;
          photo.likes++;
          console.log('Photo liked successfully');
        }

        // Refetch the likes to update the likeId and isLiked properties
        const updatedLikesResponse = await api.get(`/users/${photo.user_id}/photos/${photo.id}/likes`, {
          headers: {
            Authorization: token
          }
        });

        const updatedLikes = updatedLikesResponse.data || [];
        const updatedLike = updatedLikes.find(like => like.user_id == loggedInUserId);

        // Update photo properties
        photo.isLiked = !!updatedLike;
        photo.likeId = updatedLike ? updatedLike.id : null;

      } catch (error) {
        console.error('Error toggling like:', error);
      }
    },
    formatTimestamp(timestamp) {
      if (!timestamp || timestamp.length !== 14) {
        return ''; // Handle invalid cases, such as missing or incorrect format timestamp
      }
      // Extract individual parts of the timestamp
      const year = timestamp.substring(0, 4);
      const month = timestamp.substring(4, 6);
      const day = timestamp.substring(6, 8);
      const hours = timestamp.substring(8, 10);
      const minutes = timestamp.substring(10, 12);
      const seconds = timestamp.substring(12, 14);
      
      // Return the formatted string in the desired format
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    }
  }
};
</script>