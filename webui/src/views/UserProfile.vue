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
      isBanned: false
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
    await this.loadFollowAndBanStatus();
  },
  watch: {
    '$route.params.userId': 'loadUserProfile'
  },
  methods: {
    async loadUserProfile() {
      try {
        const userId = this.$route.params.userId;
        const response = await api.get(`/users/${userId}/profile`, {
          headers: {
            Authorization: localStorage.getItem('token')
          }
        });
        this.userProfile = response.data;
      } catch (error) {
        console.error('Error loading user profile:', error);
      }
    },
    async loadFollowAndBanStatus() {
      const userId = this.$route.params.userId;
      const loggedInUserId = localStorage.getItem('loggedInUserId');

      try {
        const followResponse = await api.get(`/users/${loggedInUserId}/follows/${userId}`, {
          headers: {
            Authorization: localStorage.getItem('token')
          }
        });
        this.isFollowing = followResponse.data.isFollowed;

        const banResponse = await api.get(`/users/${userId}/bans/${loggedInUserId}`, {
          headers: {
            Authorization: localStorage.getItem('token')
          }
        });
        this.isBanned = banResponse.data.isBanned; // Assicurati che isBanned venga impostato correttamente
      } catch (error) {
        console.error('Error loading follow and ban status:', error);
      }
    },
    async toggleFollow() {
      const userId = this.userProfile.user.id;
      const loggedInUserId = localStorage.getItem('loggedInUserId');

      try {
        if (this.isFollowing) {
          await api.delete(`/users/${loggedInUserId}/follows/${userId}`, {
            headers: {
              Authorization: localStorage.getItem('token')
            }
          });
          this.isFollowing = false;
          this.userProfile.numFollowers--;
        } else {
          await api.post(`/users/${loggedInUserId}/follows/${userId}`, {}, {
            headers: {
              Authorization: localStorage.getItem('token')
            }
          });
          this.isFollowing = true;
          this.userProfile.numFollowers++;
        }
      } catch (error) {
        console.error('Error nel toggle follow:', error);
      }
    },
    async toggleBan() {
      const userId = this.userProfile.user.id;
      const loggedInUserId = localStorage.getItem('loggedInUserId');

      try {
        if (this.isBanned) {
          await api.delete(`/users/${loggedInUserId}/bans/${userId}`, {
            headers: {
              Authorization: localStorage.getItem('token')
            }
          });
          this.isBanned = false;
        } else {
          await api.post(`/users/${loggedInUserId}/bans/${userId}`, {}, {
            headers: {
              Authorization: localStorage.getItem('token')
            }
          });
          this.isBanned = true;
        }
      } catch (error) {
        console.error('Error toggling ban:', error);
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
