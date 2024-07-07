<template>
  <div class="container mt-5">
    <div class="row">
      <div class="col-md-4 offset-md-4">
        <h1 class="text-center">{{ userProfile.user.username }}</h1>
        <p v-if="isOwnProfile" class="text-center">
          <RouterLink :to="`/users/${userProfile.user.id}/profile/edit`" class="nav-link">
            Cambia Username
          </RouterLink>
        </p>
        <p v-else class="text-center">
          <button @click="toggleBan" class="btn btn-danger">{{ isBanned ? 'Unban' : 'Ban' }}</button>
          <button @click="toggleFollow" class="btn btn-primary">{{ isFollowing ? 'Unfollow' : 'Follow' }}</button>
        </p>
      </div>
    </div>
    <div class="row mt-4">
      <div class="col-md-4">
        <div>
          <h3 class="text-center">Followers</h3>
          <p class="text-center">{{ userProfile.numFollowers }}</p>
        </div>
      </div>
      <div class="col-md-4">
        <div>
          <h3 class="text-center">Seguiti</h3>
          <p class="text-center">{{ userProfile.numFollowing }}</p>
        </div>
      </div>
      <div class="col-md-4">
        <div>
          <h3 class="text-center">Post</h3>
          <p class="text-center">{{ userProfile.numPhotos }}</p>
        </div>
      </div>
    </div>

    <div class="row mt-4">
      <div class="col-md-12">
        <ul class="listaFoto">
          <li v-for="(photo, index) in userProfile.Photos" :key="photo.id" :class="{'new-row': index % 1 === 0}">
            <div class="card">
              <img class="text-center" :src="'data:image/jpeg;base64,' + photo.image_data" alt="User Photo">
              <div class="text-center likes">
                <p>
                  <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#heart"/>
                  </svg>
                  {{ photo.likes }}
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
          </li>
        </ul>
      </div>
    </div>
    <div class="row mt-4">
      <div class="col-md-12 text-center">
        <FileUpload @photo-uploaded="addPhotoToProfile" />
      </div>
    </div>
  </div>
</template>

<script>
import api from "@/services/axios";
import { RouterLink } from "vue-router";
import FileUpload from "@/components/FileUpload.vue";

export default {
  components: {
    RouterLink,
    FileUpload
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

        // Fetch likes and comments for each photo
        await Promise.all(this.userProfile.Photos.map(photo => this.loadPhotoDetails(photo)));
      } catch (error) {
        console.error('Error loading user profile:', error);
      }
    },
    async loadPhotoDetails(photo) {
      try {
        const userId = this.userProfile.user.id;
        const likesResponse = await api.get(`/users/${userId}/photos/${photo.id}/likes`, {
          headers: {
            Authorization: localStorage.getItem('token')
          }
        });
        const commentsResponse = await api.get(`/users/${userId}/photos/${photo.id}/comments`, {
          headers: {
            Authorization: localStorage.getItem('token')
          }
        });

        photo.likes = likesResponse.data ? likesResponse.data.length : 0;
        photo.comments = commentsResponse.data ? commentsResponse.data.length : 0;
      } catch (error) {
        console.error('Error loading photo details:', error);
        photo.likes = 0;
        photo.comments = 0;
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
        this.isBanned = banResponse.data.isBanned;
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
    },
    addPhotoToProfile(photo) {
      this.userProfile.Photos.push(photo);
      this.userProfile.numPhotos++;
    }
  }
};
</script>

<style>
  h1 {
    text-transform: capitalize;
  }
  .btn {
    margin: 1rem;
    background-color: white;
    transition: 0.5s all;
    width: 7rem;
  }
  .btn-danger {
    border: solid 1px red;
    color: red;
  }
  .btn-danger:hover {
    background-color: red;
    color: white;
  }
  .btn-primary {
    border: solid 1px blue;
    color: blue;
  }
  .btn-primary:hover {
    background-color: blue;
    color: white;
  }

  .listaFoto {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    padding: 0;
    list-style-type: none;
  }

  .listaFoto li {
    width: 75%;
    margin: 2rem auto;
    box-sizing: border-box;
  }

  .listaFoto li.new-row {
    clear: both;
  }

  .card{
    width:100%;
    display: grid;
    grid-template-columns: auto auto auto;
    grid-template-rows: auto auto auto 3rem;
    box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
    transition: 0.3s;
  }

  .card:hover{
    box-shadow: 0 8px 16px 0 rgba(0,0,0,0.2)
  }

  .card img{
    margin: 2rem auto;
    width: 90%;
    grid-column-start: 1;
    grid-column-end: 4;
    grid-row-start: 1;
    grid-row-end: 4;
  }

  .comments{
    grid-column-start: 1;
    grid-column-end: 2;
    grid-row-start: 4;
    grid-row-end: 5 ;
  }

  .likes{
    grid-column-start: 3;
    grid-column-end: 4;
    grid-row-start: 4;
    grid-row-end: 5 ;
  }


</style>

