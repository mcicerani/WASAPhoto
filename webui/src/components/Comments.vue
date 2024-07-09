<template>
    <div class="mt-4">
      <div v-for="(commentResponse, index) in sortedComments" :key="commentResponse.id" class="row mb-2 align-items-start">
        <div class="col-md-8 align-self-center">
          <p><strong>{{ getUserName(commentResponse.user_id) }}:</strong> {{ commentResponse.text }}</p>
        </div>
        <div class="col-md-4 d-flex justify-content-end">
          <button v-if="commentResponse.user_id == loggedInUserId" @click="deleteComment(commentResponse.id)" class="btn btn-sm btn-danger align-self-center">Delete</button>
        </div>
      </div>
      <div class="row">
        <div class="col-md-12">
          <input type="text" v-model="comment" placeholder="Add a comment" class="form-control" />
        </div>
        <div class="col-md-12 text-center mt-2">
          <button @click="addComment" class="btn btn-primary">Submit</button>
        </div>
      </div>
    </div>
</template>
        
<script>
import api from "@/services/axios";
import { reactive } from 'vue';

export default {
    props: {
    userId: {
        type: String,
        required: true
        },
    photoId: {
        type: String,
        required: true
        }
    },
    data() {
      return {
        comments: [],
        comment: '',
        loggedInUserId: localStorage.getItem("loggedInUserId"), // Retrieve logged in user ID from localStorage
        token: localStorage.getItem("token"), // Retrieve token from localStorage
        userProfiles: reactive({}), // Usa reactive per inizializzare userProfiles
        };
    },
    mounted() {
    // Recupera l'ID dell'utente loggato da localStorage
    this.loggedInUserId = localStorage.getItem("loggedInUserId");
    },
    computed: {
      sortedComments() {
        return this.comments.slice().sort((a, b) => b.timestamp.localeCompare(a.timestamp));
      }
    },
    created() {
      this.fetchComments();
      this.fetchUserProfiles(); // chiamata al metodo all'avvio del componente
    },
    methods: {
        async fetchComments() {
            try {
                console.log('Fetching comments...');
                const response = await api.get(`/users/${this.userId}/photos/${this.photoId}/comments`, {
                headers: {
                    Authorization: localStorage.getItem("token")
                }
                });
                this.comments = response.data;
                console.log('Comments fetched:', this.comments);
                this.fetchUserProfiles();
            } catch (error) {
                console.error('Error fetching comments:', error);
            }
        },
        async fetchUserProfiles() {
            const userIds = [...new Set(this.comments.map(comment => comment.user_id))];
            for (const userId of userIds) {
                if (!this.userProfiles[userId]) {
                try {
                    console.log(`Fetching profile for user ${userId}...`);
                    const response = await api.get(`/users/${userId}/profile`, {
                    headers: {
                        Authorization: localStorage.getItem("token")
                    }
                    });
                    this.userProfiles[userId] = response.data.user.username;
                    console.log(`Profile fetched for user ${userId}:`, response.data.user.username);
                } catch (error) {
                    console.error('Error fetching user profile:', error);
                }
                }
            }
            console.log('Updated user profiles:', this.userProfiles); // Verifica userProfiles dopo il fetching
        },
        async addComment() {
            try {
                const response = await api.post(`/users/${this.userId}/photos/${this.photoId}/comments`,
                    `comment=${encodeURIComponent(this.comment)}`,
                    {
                        headers: {
                            Authorization: localStorage.getItem("token")
                        },
                    }
                );
                const newComment = response.data; // Utilizza response.data per accedere ai dati del nuovo commento
                this.comments.push(newComment); // Aggiungi il nuovo commento all'array di commenti nel frontend
                this.comment = ''; // Resetta il campo del commento dopo averlo aggiunto
                console.log('Comment added:', newComment);
                this.fetchUserProfiles(); // Esegui il fetch dei profili utente come necessario
                } catch (error) {
                console.error('Error adding comment:', error);
                }
            },
            async deleteComment(commentId) {
                try {
                    await api.delete(`/users/${this.userId}/photos/${this.photoId}/comments/${commentId}`, {
                    headers: {
                        Authorization: this.token
                    }
                    });
                    this.comments = this.comments.filter(comment => comment.id !== commentId);
                    console.log(`Comment ${commentId} deleted.`);
                } catch (error) {
                    console.error('Error deleting comment:', error);
                }
            },
            getUserName(userId) {
                return this.userProfiles[userId] || 'Unknown';
            },
        }
    };
</script>
  
<style>
    .mt-4{
        margin:auto;
        width: 95%;
    }

    .form-control{
        width: 95%;
        margin:auto;
    }

</style>
  